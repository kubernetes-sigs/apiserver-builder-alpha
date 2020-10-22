package filepath

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/conversion"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	genericapirequest "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/rest"
)

var ErrFileNotExists = fmt.Errorf("file doesn't exist")

func NewFilepathREST(
	groupResource schema.GroupResource,
	codec runtime.Codec,
	rootpath string,
	isNamespaced bool,
	newFunc func() runtime.Object,
	newListFunc func() runtime.Object) rest.Storage {

	objRoot := filepath.Join(rootpath, groupResource.Group, groupResource.Resource)
	if err := ensureDir(objRoot); err != nil {
		panic(fmt.Sprintf("unable to write data dir: %s", err))
	}

	// file REST
	rest := &filepathREST{
		TableConvertor: rest.NewDefaultTableConvertor(groupResource),
		codec:          codec,
		objRootPath:    objRoot,
		isNamespaced:   isNamespaced,
		newFunc:        newFunc,
		newListFunc:    newListFunc,
	}
	return rest
}

type filepathREST struct {
	rest.TableConvertor
	codec        runtime.Codec
	objRootPath  string
	isNamespaced bool
	newFunc      func() runtime.Object
	newListFunc  func() runtime.Object
}

var _ rest.Storage = &filepathREST{}
var _ rest.Creater = &filepathREST{}
var _ rest.Updater = &filepathREST{}
var _ rest.GracefulDeleter = &filepathREST{}
var _ rest.CollectionDeleter = &filepathREST{}
var _ rest.Getter = &filepathREST{}
var _ rest.Lister = &filepathREST{}
var _ rest.Scoper = &filepathREST{}

func (f *filepathREST) New() runtime.Object {
	return f.newFunc()
}

func (f *filepathREST) NewList() runtime.Object {
	return f.newListFunc()
}

func (f *filepathREST) NamespaceScoped() bool {
	return f.isNamespaced
}

func (f *filepathREST) Get(
	ctx context.Context,
	name string,
	options *metav1.GetOptions,
) (runtime.Object, error) {
	return read(f.codec, f.objectFileName(ctx, name), f.newFunc)
}

func (f *filepathREST) List(
	ctx context.Context,
	options *internalversion.ListOptions,
) (runtime.Object, error) {
	newListObj := f.NewList()
	v, err := getListPrt(newListObj)
	if err != nil {
		return nil, err
	}

	dirname := f.objectDirName(ctx)
	if err := visitDir(dirname, f.newFunc, f.codec, func(path string, obj runtime.Object) {
		appendItem(v, obj)
	}); err != nil {
		return nil, fmt.Errorf("failed walking filepath %v", dirname)
	}
	return newListObj, nil
}

var (
	ErrNamespaceNotExists = errors.New("namespace does not exist")
)

func (f *filepathREST) Create(
	ctx context.Context,
	obj runtime.Object,
	createValidation rest.ValidateObjectFunc,
	options *metav1.CreateOptions,
) (runtime.Object, error) {

	if createValidation != nil {
		if err := createValidation(ctx, obj); err != nil {
			return nil, err
		}
	}

	if f.isNamespaced {
		// ensures namespace dir
		ns, ok := genericapirequest.NamespaceFrom(ctx)
		if !ok {
			return nil, ErrNamespaceNotExists
		}
		if err := ensureDir(filepath.Join(f.objRootPath, ns)); err != nil {
			return nil, err
		}
	}

	accessor, err := meta.Accessor(obj)
	if err != nil {
		return nil, err
	}
	filename := f.objectFileName(ctx, accessor.GetName())

	if exists(filename) {
		return nil, ErrFileNotExists
	}

	if err := write(f.codec, filename, obj); err != nil {
		return nil, err
	}
	return obj, nil
}

func (f *filepathREST) Update(
	ctx context.Context,
	name string,
	objInfo rest.UpdatedObjectInfo,
	createValidation rest.ValidateObjectFunc,
	updateValidation rest.ValidateObjectUpdateFunc,
	forceAllowCreate bool,
	options *metav1.UpdateOptions,
) (runtime.Object, bool, error) {

	isCreate := false
	oldObj, err := f.Get(ctx, name, nil)
	if err != nil {
		if !forceAllowCreate {
			return nil, false, err
		}
		isCreate = true
	}

	// TODO: should not be necessary, verify Get works before creating filepath
	if f.isNamespaced {
		// ensures namespace dir
		ns, ok := genericapirequest.NamespaceFrom(ctx)
		if !ok {
			return nil, false, ErrNamespaceNotExists
		}
		if err := ensureDir(filepath.Join(f.objRootPath, ns)); err != nil {
			return nil, false, err
		}
	}

	updatedObj, err := objInfo.UpdatedObject(ctx, oldObj)
	filename := f.objectFileName(ctx, name)

	if isCreate {
		if createValidation != nil {
			if err := createValidation(ctx, updatedObj); err != nil {
				return nil, false, err
			}
		}
		if err := write(f.codec, filename, updatedObj); err != nil {
			return nil, false, err
		}
		return updatedObj, true, nil
	}

	if updateValidation != nil {
		if err := updateValidation(ctx, updatedObj, oldObj); err != nil {
			return nil, false, err
		}
	}
	if err := write(f.codec, filename, updatedObj); err != nil {
		return nil, false, err
	}
	return updatedObj, false, nil
}

func (f *filepathREST) Delete(ctx context.Context, name string, deleteValidation rest.ValidateObjectFunc, options *metav1.DeleteOptions) (runtime.Object, bool, error) {
	filename := f.objectFileName(ctx, name)
	if !exists(filename) {
		return nil, false, ErrFileNotExists
	}

	oldObj, err := f.Get(ctx, name, nil)
	if err != nil {
		return nil, false, err
	}
	if deleteValidation != nil {
		if err := deleteValidation(ctx, oldObj); err != nil {
			return nil, false, err
		}
	}

	if err := os.Remove(filename); err != nil {
		return nil, false, err
	}
	return oldObj, true, nil
}

func (f *filepathREST) DeleteCollection(
	ctx context.Context,
	deleteValidation rest.ValidateObjectFunc,
	options *metav1.DeleteOptions,
	listOptions *internalversion.ListOptions,
) (runtime.Object, error) {
	newListObj := f.NewList()
	v, err := getListPrt(newListObj)
	if err != nil {
		return nil, err
	}
	dirname := f.objectDirName(ctx)
	if err := visitDir(dirname, f.newFunc, f.codec, func(path string, obj runtime.Object) {
		os.Remove(path)
		appendItem(v, obj)
	}); err != nil {
		return nil, fmt.Errorf("failed walking filepath %v", dirname)
	}
	return newListObj, nil
}

func (f *filepathREST) objectFileName(ctx context.Context, name string) string {
	if f.isNamespaced {
		// FIXME: return error if namespace is not found
		ns, _ := genericapirequest.NamespaceFrom(ctx)
		return filepath.Join(f.objRootPath, ns, name+".json")
	}
	return filepath.Join(f.objRootPath, name+".json")
}

func (f *filepathREST) objectDirName(ctx context.Context) string {
	if f.isNamespaced {
		// FIXME: return error if namespace is not found
		ns, _ := genericapirequest.NamespaceFrom(ctx)
		return filepath.Join(f.objRootPath, ns)
	}
	return filepath.Join(f.objRootPath)
}

func write(encoder runtime.Encoder, filepath string, obj runtime.Object) error {
	buf := new(bytes.Buffer)
	if err := encoder.Encode(obj, buf); err != nil {
		return err
	}
	return ioutil.WriteFile(filepath, buf.Bytes(), 0700)
}

func read(decoder runtime.Decoder, filepath string, newFunc func() runtime.Object) (runtime.Object, error) {
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	newObj := newFunc()
	decodedObj, _, err := decoder.Decode(content, nil, newObj)
	if err != nil {
		return nil, err
	}
	return decodedObj, nil
}

func exists(filepath string) bool {
	_, err := os.Stat(filepath)
	return err == nil
}

func ensureDir(dirname string) error {
	if !exists(dirname) {
		return os.MkdirAll(dirname, 0700)
	}
	return nil
}

func visitDir(dirname string, newFunc func() runtime.Object, codec runtime.Codec, visitFunc func(string, runtime.Object)) error {
	return filepath.Walk(dirname, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(info.Name(), ".json") {
			return nil
		}
		newObj, err := read(codec, path, newFunc)
		if err != nil {
			return err
		}
		visitFunc(path, newObj)
		return nil
	})
}

func appendItem(v reflect.Value, obj runtime.Object) {
	v.Set(reflect.Append(v, reflect.ValueOf(obj).Elem()))
}

func getListPrt(listObj runtime.Object) (reflect.Value, error) {
	listPtr, err := meta.GetItemsPtr(listObj)
	if err != nil {
		return reflect.Value{}, err
	}
	v, err := conversion.EnforcePtr(listPtr)
	if err != nil || v.Kind() != reflect.Slice {
		return reflect.Value{}, fmt.Errorf("need ptr to slice: %v", err)
	}
	return v, nil
}

package filepath

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	genericapirequest "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/client-go/kubernetes/scheme"
)

func TestClusterResource(t *testing.T) {
	root := "./test-cluster-resource"
	os.Mkdir(root, 0700)
	defer os.RemoveAll(root)

	codec, ok := runtime.SerializerInfoForMediaType(scheme.Codecs.SupportedMediaTypes(), "application/json")
	require.True(t, ok)

	r := NewFilepathREST(
		schema.GroupResource{
			Group:    "g1",
			Resource: "r1",
		},
		codec.Serializer,
		root,
		false,
		func() runtime.Object { return &v1.Pod{} },
		func() runtime.Object { return &v1.PodList{} },
	)
	RESTTestFunc(t, r, "", context.TODO())
}

func TestNamespacedResource(t *testing.T) {
	root := "./test-namespaced-resource"
	namespace := "test-namespace"
	os.Mkdir(root, 0700)
	defer os.RemoveAll(root)

	codec, ok := runtime.SerializerInfoForMediaType(scheme.Codecs.SupportedMediaTypes(), "application/json")
	require.True(t, ok)

	r := NewFilepathREST(
		schema.GroupResource{
			Group:    "",
			Resource: "",
		},
		codec.Serializer,
		root,
		true,
		func() runtime.Object { return &v1.Pod{} },
		func() runtime.Object { return &v1.PodList{} },
	)

	ctx := genericapirequest.WithNamespace(context.TODO(), namespace)
	RESTTestFunc(t, r, namespace, ctx)
}

func RESTTestFunc(t *testing.T, r rest.Storage, namespace string, ctx context.Context) {
	name1 := "foo-1"
	obj1 := r.New()
	accessor1, err := meta.Accessor(obj1)
	require.NoError(t, err)
	accessor1.SetName(name1)
	accessor1.SetNamespace(namespace)

	// get should fail when no such object
	obj, err := r.(rest.Getter).Get(
		ctx,
		name1,
		nil)
	assert.Error(t, err)
	assert.Nil(t, obj)

	// get should fail when no such object
	obj, created, err := r.(rest.Updater).Update(
		ctx,
		name1,
		rest.DefaultUpdatedObjectInfo(obj1),
		nil,
		nil,
		false,
		nil)
	assert.Error(t, err)
	assert.False(t, created)
	assert.Nil(t, obj)

	// delete should fail when no such object
	oldObj, deleted, err := r.(rest.GracefulDeleter).Delete(
		ctx,
		name1,
		nil,
		nil)
	assert.Error(t, err)
	assert.False(t, deleted)
	assert.Nil(t, oldObj)

	// create should work
	obj, err = r.(rest.Creater).Create(
		ctx,
		obj1,
		nil,
		nil)
	assert.NoError(t, err)
	assert.NotNil(t, obj)

	// create should fail when already exists
	obj, err = r.(rest.Creater).Create(
		ctx,
		obj1,
		nil,
		nil)
	assert.Error(t, err)
	assert.Nil(t, obj)

	// get should work
	obj, err = r.(rest.Getter).Get(
		ctx,
		name1,
		nil)
	assert.NoError(t, err)
	assert.NotNil(t, obj)

	// update should work
	obj, created, err = r.(rest.Updater).Update(
		ctx,
		name1,
		rest.DefaultUpdatedObjectInfo(obj),
		nil,
		nil,
		false,
		nil)
	assert.NoError(t, err)
	assert.False(t, created)
	assert.NotNil(t, obj)

	name2 := "foo-2"
	obj2 := r.New()
	accessor2, err := meta.Accessor(obj2)
	require.NoError(t, err)
	accessor2.SetName(name2)
	accessor2.SetNamespace(namespace)

	// update-create should work
	obj2, created, err = r.(rest.Updater).Update(
		ctx,
		name2,
		rest.DefaultUpdatedObjectInfo(obj2),
		nil,
		nil,
		true,
		nil)
	assert.NoError(t, err)
	assert.True(t, created)
	assert.NotNil(t, obj2)

	// list should work
	objList, err := r.(rest.Lister).List(
		ctx,
		nil)
	assert.NoError(t, err)
	assert.NotNil(t, objList)
	assert.Equal(t, 2, len(objList.(*v1.PodList).Items))

	// delete should work
	obj, deleted, err = r.(rest.GracefulDeleter).Delete(
		ctx,
		name1,
		nil,
		nil)
	assert.NoError(t, err)
	assert.True(t, deleted)
	assert.NotNil(t, obj)

	// delete-collection should work
	deletedList, err := r.(rest.CollectionDeleter).DeleteCollection(
		ctx,
		nil,
		nil,
		nil)
	assert.NoError(t, err)
	assert.NotNil(t, deletedList)

	// re-list should return nothing
	objList, err = r.(rest.Lister).List(
		ctx,
		nil)
	assert.NoError(t, err)
	assert.NotNil(t, objList)
	assert.Equal(t, 0, len(objList.(*v1.PodList).Items))
}

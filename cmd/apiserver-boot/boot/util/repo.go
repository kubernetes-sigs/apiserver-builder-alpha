package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/mod/modfile"
	"k8s.io/klog"
)

var repo string

func LoadRepoFromGoPath() error {
	gopath := os.Getenv("GOPATH")
	if len(gopath) == 0 {
		return fmt.Errorf("GOPATH not defined")
	}
	goSrc := filepath.Join(gopath, "src")
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	if !strings.HasPrefix(filepath.Dir(wd), goSrc) {
		return fmt.Errorf("apiserver-boot must be run from the directory containing the go package to "+
			"bootstrap. This must be under $GOPATH/src/<package>. "+
			"\nCurrent GOPATH=%s.  \nCurrent directory=%s", gopath, wd)
	}
	repo = strings.Replace(wd, goSrc+string(filepath.Separator), "", 1)
	return nil
}

func LoadRepoFromGoMod() error {
	mod, err := ioutil.ReadFile("go.mod")
	if err != nil {
		return errors.Wrap(err, "failed reading go.mod file")
	}
	modPath := modfile.ModulePath(mod)
	if len(modPath) == 0 {
		return fmt.Errorf("failed parsing go.mod, empty module path")
	}
	repo = modPath
	return nil
}

func LoadRepoFromGoPathOrGoMod() error {
	if err := LoadRepoFromGoPath(); err != nil {
		// reading from go mod
		return LoadRepoFromGoMod()
	}
	return nil
}

func GetRepo() string {
	if len(repo) > 0 {
		return repo
	}
	if err := LoadRepoFromGoPathOrGoMod(); err != nil {
		klog.Fatal(err)
	}
	return repo
}

func SetRepo(r string) {
	repo = r
}

/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package release

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func PackageTar(goos, goarch, tooldir, vendordir string) {
	// create the new file
	fw, err := os.Create(fmt.Sprintf("%s-%s-%s-%s.tar.gz", output, version, goos, goarch))
	if err != nil {
		log.Fatalf("failed to create output file %s %v", output, err)
	}
	defer fw.Close()

	// setup gzip of tar
	gw := gzip.NewWriter(fw)
	defer gw.Close()

	// setup tar writer
	tw := tar.NewWriter(gw)
	defer tw.Close()

	// Add all of the bin files
	// Add all of the bin files
	filepath.Walk(filepath.Join(tooldir, "bin"), TarFile{
		tw,
		0555,
		tooldir,
		"",
	}.Do)
}

type TarFile struct {
	Writer *tar.Writer
	Mode   int64
	Root   string
	Parent string
}

func (t TarFile) Do(path string, info os.FileInfo, err error) error {
	if info.IsDir() {
		return nil
	}

	eval, err := filepath.EvalSymlinks(path)
	if err != nil {
		log.Fatal(err)
	}
	if eval != path {
		name := strings.Replace(path, t.Root, "", -1)
		if len(t.Parent) != 0 {
			name = filepath.Join(t.Parent, name)
		}
		linkName := strings.Replace(eval, t.Root, "", -1)
		if len(t.Parent) != 0 {
			linkName = filepath.Join(t.Parent, linkName)
		}
		hdr := &tar.Header{
			Name:     name,
			Mode:     t.Mode,
			Linkname: linkName,
		}
		if err := t.Writer.WriteHeader(hdr); err != nil {
			log.Fatalf("failed to write output for %s %v", path, err)
		}
		return nil
	}

	return t.Write(path)
}

func (t TarFile) Write(path string) error {
	// Get the relative name of the file
	name := strings.Replace(path, t.Root, "", -1)
	if len(t.Parent) != 0 {
		name = filepath.Join(t.Parent, name)
	}
	body, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("failed to read file %s %v", path, err)
	}
	if len(body) == 0 {
		return nil
	}

	hdr := &tar.Header{
		Name: name,
		Mode: t.Mode,
		Size: int64(len(body)),
	}
	if err := t.Writer.WriteHeader(hdr); err != nil {
		log.Fatalf("failed to write output for %s %v", path, err)
	}
	if _, err := t.Writer.Write(body); err != nil {
		log.Fatalf("failed to write output for %s %v", path, err)
	}
	return nil
}

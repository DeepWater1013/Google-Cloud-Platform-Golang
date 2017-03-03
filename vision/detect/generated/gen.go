// Copyright 2017 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

//+build ignore

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

func main() {
	tmpl, err := ioutil.ReadFile("generated/sample-template.go")
	if err != nil {
		log.Fatal("ReadFile:", err)
	}

	// strip build tag from the top of the file.
	start := bytes.Index(tmpl, []byte("package main"))
	tmpl = tmpl[start:]

	out := &bytes.Buffer{}

	// Header.
	fmt.Fprintln(out, header)

	normal := string(tmpl)
	normal = strings.Replace(normal, boilerplateSentinel, boilerplate, -1)
	out.WriteString(normal)

	// Don't do imports twice.
	start = bytes.Index(tmpl, []byte("\nfunc "))
	tmpl = tmpl[start:]

	gcs := string(tmpl)
	gcs = strings.Replace(gcs, boilerplateSentinel, gcsBoilerplate, -1)
	// Append suffix to function name.
	gcs = strings.Replace(gcs, "(w io.Writer", "GCS(w io.Writer", -1)
	out.WriteString(gcs)

	if err := ioutil.WriteFile("detect.go", out.Bytes(), 0640); err != nil {
		log.Fatal(err)
	}
}

const boilerplateSentinel = "\t// Boilerplate is inserted by gen.go\n"

const boilerplate = `	ctx := context.Background()

	client, err := vision.NewClient(ctx)
	if err != nil {
		return err
	}

	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	image, err := vision.NewImageFromReader(f)
	if err != nil {
		return err
	}
`
const gcsBoilerplate = `	ctx := context.Background()

	client, err := vision.NewClient(ctx)
	if err != nil {
		return err
	}

	image := vision.NewImageFromURI(file)
`

const header = `// Copyright 2017 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

//go:generate go run generated/gen.go

// DO NOT EDIT THIS FILE.
// It is generated from the source in generated/sample-template.go
`

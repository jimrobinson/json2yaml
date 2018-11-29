package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ghodss/yaml"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	if len(os.Args) == 1 {
		process("stdin", os.Stdin)
	} else {
		for _, fpath := range os.Args[1:] {
			fh, err := os.Open(fpath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "unable to open %s: %v\n", fpath, err)
				continue
			}
			defer fh.Close()
			process(fpath, fh)
			fh.Close()
		}
	}
}

func process(label string, r io.Reader) {
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s: %v\n", label, err)
		return
	}

	var out []byte
	if strings.Contains(os.Args[0], "json2yaml") {
		out, err = yaml.JSONToYAML(buf)
	} else {
		out, err = yaml.YAMLToJSON(buf)
		if err == nil {
			dst := &bytes.Buffer{}
			err = json.Indent(dst, out, "", "\t")
			if err == nil {
				out = dst.Bytes()
			} else {
				fmt.Fprintf(os.Stderr, "warning: unable to reindent json for %s: %v", label, err)
			}
		}
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "error parsing %s: %v\n", label, err)
	}

	fmt.Println(string(out))
}

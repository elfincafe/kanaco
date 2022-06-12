package kanaco

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
)

func mode(path string) string {
	basename := filepath.Base(path)
	ext := filepath.Ext(basename)
	mode := strings.ReplaceAll(strings.ReplaceAll(basename, "output.", ""), ext, "")
	return mode
}

func TestByte(t *testing.T) {
	content, err := ioutil.ReadFile("./data/input.txt")
	if err != nil {
		t.Errorf(err.Error())
	}
	paths, err := filepath.Glob("./data/output.srk.txt")
	for _, path := range paths {
		mode := mode(path)
		expect, err := ioutil.ReadFile(path)
		if err != nil {
			t.Errorf(err.Error())
		}
		result := Byte(content, mode)
		if bytes.Compare(result, expect) != 0 {
			expects := bytes.Split(expect, []byte("\n"))
			results := bytes.Split(result, []byte("\n"))
			msg := strings.Builder{}
			msg.WriteString(fmt.Sprintf("\n[%s] ---------\n", mode))
			for k, e := range expects {
				r := results[k]
				if bytes.Compare(e, r) != 0 {
					msg.WriteString("Expect: ")
					msg.Write(e)
					msg.WriteString("\nResult: ")
					msg.Write(r)
					msg.WriteString("\n")
				}
			}
			t.Errorf(msg.String())
		}
	}
}

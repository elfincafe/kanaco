package kanaco

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var (
	output = "output.*.txt"
)

func mode4Test(path string) string {
	basename := filepath.Base(path)
	ext := filepath.Ext(basename)
	tmp := strings.Split(strings.ReplaceAll(strings.ReplaceAll(basename, "output.", ""), ext, ""), ".")
	mode := strings.Builder{}
	for _, m := range tmp {
		switch m {
		case "la":
			mode.WriteString("a")
		case "lc":
			mode.WriteString("c")
		case "lh":
			mode.WriteString("h")
		case "lk":
			mode.WriteString("k")
		case "ln":
			mode.WriteString("n")
		case "lr":
			mode.WriteString("r")
		case "ls":
			mode.WriteString("s")
		case "ua":
			mode.WriteString("A")
		case "uc":
			mode.WriteString("C")
		case "uh":
			mode.WriteString("H")
		case "uk":
			mode.WriteString("K")
		case "un":
			mode.WriteString("N")
		case "ur":
			mode.WriteString("R")
		case "us":
			mode.WriteString("S")
		}
	}
	return mode.String()
}

func TestByte(t *testing.T) {
	content, err := ioutil.ReadFile("./data/input.txt")
	if err != nil {
		t.Errorf(err.Error())
	}
	paths, err := filepath.Glob("./data/" + output)
	if err != nil {
		t.Errorf(err.Error())
	}
	for _, path := range paths {
		mode := mode4Test(path)
		expect, err := ioutil.ReadFile(path)
		if err != nil {
			t.Errorf(err.Error())
		}
		result := Byte(content, mode)
		if strings.Compare(string(result), string(expect)) != 0 {
			expects := bytes.Split(expect, []byte("\n"))
			results := bytes.Split(result, []byte("\n"))
			msg := strings.Builder{}
			msg.WriteString(fmt.Sprintf("\n[%s] ---------\n", mode))
			for k, e := range expects {
				r := results[k]
				fmt.Printf("[%d] %s(%T) <-> %s\n", k, e, e, results[k])
				if strings.Compare(string(r), string(e)) != 0 {
					msg.WriteString(fmt.Sprintf("Expect(%d): ", k))
					msg.Write(e)
					msg.WriteString(fmt.Sprintf("\nResult(%d): ", k))
					msg.Write(r)
					msg.WriteString("\n")
				}
			}
			t.Errorf(msg.String())
		}
	}
}

func TestString(t *testing.T) {
	content, err := ioutil.ReadFile("./data/input.txt")
	if err != nil {
		t.Errorf(err.Error())
	}
	paths, err := filepath.Glob("./data/" + output)
	if err != nil {
		t.Errorf(err.Error())
	}
	for _, path := range paths {
		mode := mode4Test(path)
		expect, err := ioutil.ReadFile(path)
		if err != nil {
			t.Errorf(err.Error())
		}
		result := String(string(content), mode)
		if strings.Compare(string(result), string(expect)) != 0 {
			expects := bytes.Split(expect, []byte("\n"))
			results := strings.Split(result, "\n")
			msg := strings.Builder{}
			msg.WriteString(fmt.Sprintf("\n[%s] ---------\n", mode))
			for k, e := range expects {
				r := results[k]
				if strings.Compare(string(r), string(e)) != 0 {
					msg.WriteString(fmt.Sprintf("Expect(%d): ", k))
					msg.Write(e)
					msg.WriteString(fmt.Sprintf("\nResult(%d): ", k))
					msg.WriteString(r)
					msg.WriteString("\n")
				}
			}
			t.Errorf(msg.String())
		}
	}
}

func TestNewReader(t *testing.T) {
	f, _ := os.Open("./data/input.txt")
	r := NewReader(f, "a")
	tp := fmt.Sprintf("%T", r)
	if tp != "*kanaco.Reader" {
		t.Errorf("Reader is invalid (%s)", tp)
	}
}

func TestRead(t *testing.T) {
	paths, _ := filepath.Glob("./data/" + output)
	for _, path := range paths {
		expects, _ := ioutil.ReadFile(path)
		mode := mode4Test(path)
		f, _ := os.Open("./data/input.txt")
		r := NewReader(f, mode)
		results := []byte{}
		for {
			buf := make([]byte, 4096)
			_, err := r.Read(buf)
			if err == io.EOF {
				break
			}
			if err != nil {
				t.Errorf(err.Error())
				break
			}
			results = append(results, buf...)
		}
		if strings.Compare(string(results), string(expects)) == 0 {
			rLines := bytes.Split(results, []byte("\n"))
			eLines := bytes.Split(expects, []byte("\n"))
			msg := strings.Builder{}
			msg.WriteString(fmt.Sprintf("\n[%s] ---------\n", mode))
			for k, et := range eLines {
				rs := rLines[k]
				if strings.Compare(string(et), string(rs)) != 0 {
					msg.WriteString(fmt.Sprintf("Expect(%d): ", k))
					msg.Write(et)
					msg.WriteString(fmt.Sprintf("\nResult(%d): ", k))
					msg.Write(rs)
					msg.WriteString("\n")
				}
			}
			t.Errorf(msg.String())
		}
	}
}

func TestIs1Byte(t *testing.T) {
	cases := map[string]bool{
		"!": true, "a": true, "0": true, "}": true, "~": true,
		"¡": false, "¢": false, "߹": false, "ߺ": false,
		"✀": false, "ぁ": false, "ァ": false, "ｦ": false,
		"🌀": false, "💯": false, "😀": false, "🧦": false,
	}
	for input, output := range cases {
		if output != is1Byte([]byte(input)) {
			t.Errorf(`Function is1Byte returns %v for "%s".`, output, input)
		}
	}
}

func TestIs2Bytes(t *testing.T) {
	cases := map[string]bool{
		"!": false, "a": false, "0": false, "}": false, "~": false,
		"¡": true, "¢": true, "߹": true, "ߺ": true,
		"✀": false, "ぁ": false, "ァ": false, "ｦ": false,
		"🌀": false, "💯": false, "😀": false, "🧦": false,
	}
	for input, output := range cases {
		if output != is2Bytes([]byte(input)) {
			t.Errorf(`Function is2Bytes returns %v for "%s".`, output, input)
		}
	}
}

func TestIs3Bytes(t *testing.T) {
	cases := map[string]bool{
		"!": false, "a": false, "0": false, "}": false, "~": false,
		"¡": false, "¢": false, "߹": false, "ߺ": false,
		"✀": true, "ぁ": true, "ァ": true, "ｦ": true,
		"🌀": false, "💯": false, "😀": false, "🧦": false,
	}
	for input, output := range cases {
		if output != is3Bytes([]byte(input)) {
			t.Errorf(`Function is3Bytes returns %v for "%s".`, output, input)
		}
	}
}

func TestIs4Bytes(t *testing.T) {
	cases := map[string]bool{
		"!": false, "a": false, "0": false, "}": false, "~": false,
		"¡": false, "¢": false, "߹": false, "ߺ": false,
		"✀": false, "ぁ": false, "ァ": false, "ｦ": false,
		"🌀": true, "💯": true, "😀": true, "🧦": true,
	}
	for input, output := range cases {
		if output != is4Bytes([]byte(input)) {
			t.Errorf(`Function is4Bytes returns %v for "%s".`, output, input)
		}
	}
}

func TestIsVoiced(t *testing.T) {
	cases := map[string]bool{
		"ﾊ": false, "ﾊﾞ": true, "ﾊﾟ": false, "ハ": false, "バ": false,
	}
	for input, output := range cases {
		if output != isVoiced([]byte(input)) {
			t.Errorf(`Function isVoiced returns %v for "%s".`, output, input)
		}
	}
}

func TestIsSemiVoiced(t *testing.T) {
	cases := map[string]bool{
		"ﾊ": false, "ﾊﾞ": false, "ﾊﾟ": true, "ハ": false, "バ": false,
	}
	for input, output := range cases {
		if output != isSemiVoiced([]byte(input)) {
			t.Errorf(`Function isSemiVoiced returns %v for "%s".`, output, input)
		}
	}
}

func TestLowerR(t *testing.T) {
	cases := map[string]string{
		"ａ": "a", "ｚ": "z", "Ａ": "A", "Ｚ": "Z",
		"/": "/", "0": "0", ":": ":",
	}
	for input, output := range cases {
		c := extract([]byte(input))
		r := lowerR(c)
		if strings.Compare(string(r), output) != 0 {
			t.Errorf(`Function lowerR returns "%s" for "%s", expecting "%s".`, string(r), input, output)
		}
	}
}

func TestUpperR(t *testing.T) {
	cases := map[string]string{
		"a": "ａ", "z": "ｚ", "A": "Ａ", "Z": "Ｚ",
		"/": "/", "0": "0", ":": ":",
	}
	for input, output := range cases {
		c := extract([]byte(input))
		r := upperR(c)
		if strings.Compare(string(r), output) != 0 {
			t.Errorf(`Function upperR returns "%s" for "%s", expecting "%s".`, string(r), input, output)
		}
	}
}

func TestLowerN(t *testing.T) {
}

func TestUpperN(t *testing.T) {
}

func TestLowerA(t *testing.T) {
}

func TestUpperA(t *testing.T) {
}

func TestLowerS(t *testing.T) {
}

func TestUpperS(t *testing.T) {
}

func TestLowerK(t *testing.T) {
}

func TestUpperK(t *testing.T) {
}

func TestLowerH(t *testing.T) {
}

func TestUpperH(t *testing.T) {
}

func TestLowerC(t *testing.T) {
}

func TestUpperC(t *testing.T) {
}

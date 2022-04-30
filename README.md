# KanaCo

## Overview
KanaCo is the kana character converter inspired by the function mb_convert_kana in PHP.

## Install

    # go get github.com/elfincafe/kanaco

## Mode

|Mode|Description|
|-|-|
|r|Convert zenkaku alphabets to hankaku|
|R|Convert hankaku alphabets to zenkaku|
|n|Convert zenkaku numbers to hankaku|
|N|Convert hankaku numbers to zenkaku|
|a|Convert zenkaku alphabets and numbers to hankaku (U+0021 - U+007E excluding U+0022, U+0027, U+005C, U+007E)|
|A|Convert hankaku alphabets and numbers to zenkaku (U+0021 - U+007E excluding U+0022, U+0027, U+005C, U+007E)|
|s|Convert zenkaku space to hankaku (U+3000 -> U+0020)|
|S|Convert hankaku space to zenkaku (U+0020 -> U+3000)|
|k|Convert zenkaku katakana to hankaku katakana|
|K|Convert hankaku katakana to zenkaku katakana|
|h|Convert zenkaku hiragana to hankaku katakana|
|H|Convert hankaku katakana to zenkaku hiragana|
|c|Convert zenkaku katakana to hankaku hiragana|
|C|Convert zenkaku hiragana to zenkaku katakana|

## Usage
```go
package main
	
import (
    "io"
    "io/ioutil"
    "github.com/elfincafe/kanaco"
)

func main () {

    // String Style
    in1 := "123abcABC １２３ａｂｃＡＢＣ"
    out1 := kanaco.String(in1, "a")
    println(out1) // 123abcABC 123abcABC

    // Byte Style
    in2 := []byte("123abcABC １２３ａｂｃＡＢＣ")
    out2 := kanaco.Byte(in2, "RS")
    println(out2) // １２３ａｂｃＡＢＣ　１２３ａｂｃＡＢＣ

    // Reader Style
    ioutil.WriteFile("example.txt", []byte("ｶﾅｺ　ｺﾝﾊﾞｰﾀｰ　Ｖｅｒ１"), 0644)
    f, _ := os.Open("example.txt")
    reader := kanaco.NewReader(f, "Kas")
    for {
        buf := make([]byte, 4096)
        _, err := reader.Read(buf)
        if err == io.EOF {
            break
        }
        println(string(buf)) // カナコ コンバｰタｰ Ver1
    }
}
```

## License
KanaCo is distributed under The MIT License.  
https://opensource.org/licenses/mit-license.php


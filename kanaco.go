package kanaco

import (
	"bufio"
	"fmt"
	"io"
)

const (
	BufChars    int = 6
	FLT_ASIS    int = 0
	FLT_LOWER_R int = 1 << 0
	FLT_UPPER_R int = 1 << 1
	FLT_LOWER_N int = 1 << 2
	FLT_UPPER_N int = 1 << 3
	FLT_LOWER_A int = 1 << 4
	FLT_UPPER_A int = 1 << 5
	FLT_LOWER_S int = 1 << 6
	FLT_UPPER_S int = 1 << 7
	FLT_LOWER_K int = 1 << 8
	FLT_UPPER_K int = 1 << 9
	FLT_LOWER_H int = 1 << 10
	FLT_UPPER_H int = 1 << 11
	FLT_LOWER_C int = 1 << 12
	FLT_UPPER_C int = 1 << 13
)

type (
	Reader struct {
		r    *bufio.Reader
		mode string
	}
	character struct {
		val     []byte
		cval    []byte // converted value
		filters int    // FLT_LOWER_* or FLT_UPPER_*
	}
	filter func(*character)
)

func Byte(b []byte, mode string) []byte {
	if len(mode) == 0 {
		return []byte{}
	}
	filters := createFilters(mode)

	buf := make([]byte, 0, 512)
	c := new(character)
	length := len(b)
	for i := 0; i < length; i++ {
		c.init()
		extract(c, b[i:])
		conv(c, filters)
		buf = append(buf, c.cval...)
		proseed := len(c.val) - 1
		if proseed < 0 {
			proseed = 0
		}
		i += proseed
	}
	return buf
}

func String(str, mode string) string {
	return string(Byte([]byte(str), mode))
}

func NewReader(r io.Reader, mode string) *Reader {
	reader := new(Reader)
	reader.r = bufio.NewReader(r)
	reader.mode = mode
	return reader
}

func (r *Reader) Read(p []byte) (int, error) {
	line, err := r.r.ReadBytes('\n')
	if err == io.EOF {
		return 0, err
	}
	if err != nil {
		return 0, err
	}
	line = Byte(line, r.mode)
	if len(p) < len(line) {
		return 0, fmt.Errorf("buffer size is not enough")
	}
	n := copy(p, line)
	return n, nil
}

// -------------------------------------

func (c *character) init() {
	c.val = []byte{}
	c.cval = []byte{}
	c.filters = FLT_ASIS
}

func is1Byte(b []byte) bool {
	if len(b) == 0 {
		return false
	}
	if b[0]&0x80 == 0x00 {
		return true
	} else {
		return false
	}
}

func is2Bytes(b []byte) bool {
	if len(b) < 2 {
		return false
	}
	if (b[0]&0xc2 == 0xc2) && (b[1]&0x80 == 0x80) {
		return true
	} else {
		return false
	}
}

func is3Bytes(b []byte) bool {
	if len(b) < 3 {
		return false
	}
	if (b[0]&0xe0 == 0xe0) && (b[1]&0x80 == 0x80) && (b[2]&0x80 == 0x80) {
		return true
	} else {
		return false
	}
}

func is4Bytes(b []byte) bool {
	if len(b) < 4 {
		return false
	}
	if (b[0]&0xf0 == 0xf0) && (b[1]&0x80 == 0x80) && (b[2]&0x80 == 0x80) && (b[3]&0x80 == 0x80) {
		return true
	} else {
		return false
	}
}

func isVoiced(b []byte) bool {
	if len(b) < 6 {
		return false
	}
	if b[3] == 0xef && b[4] == 0xbe && b[5] == 0x9e {
		if (b[0] == 0xef) && (b[1] == 0xbd) && (b[2] > 0xb5) && (b[2] < 0xc0) { // ｶ - ｿ
			return true
		} else if b[0] == 0xef && b[1] == 0xbe && b[2] > 0x79 && b[2] < 0x85 { // ﾀ - ﾄ
			return true
		} else if b[0] == 0xef && b[1] == 0xbe && b[2] > 0x89 && b[2] < 0x8f { // ﾊ - ﾎ
			return true
		} else if b[0] == 0xef && b[1] == 0xbd && b[2] == 0xb3 { // ｳ
			return true
		}
	}
	return false
}

func isSemiVoiced(b []byte) bool {
	if len(b) < 6 {
		return false
	}
	if (b[3] == 0xef) && (b[4] == 0xbe) && (b[5] == 0x9f) {
		if (b[0] == 0xef) && (b[1] == 0xbe) && (b[2] > 0x89) && (b[2] < 0x8f) { // ﾊ - ﾎ
			return true
		}
	}
	return false
}

func lowerR(c *character) {
	if c.filters&FLT_LOWER_R != FLT_LOWER_R {
		return
	}
	// Ａ-Ｚ -> A-Z
	if c.val[2] >= 0xa1 && c.val[2] <= 0xba {
		c.cval = []byte{c.val[2] - 0x60}
	}
	// ａ-ｚ -> a-z
	if c.val[2] >= 0x81 && c.val[2] <= 0x9a {
		c.cval = []byte{c.val[2] - 0x20}
	}
}

func upperR(c *character) {
	if c.filters&FLT_UPPER_R != FLT_UPPER_R {
		return
	}
	// A-Z -> Ａ-Ｚ
	if c.val[0] >= 0x41 && c.val[0] <= 0x5a {
		c.cval = []byte{0xef, 0xbc, c.val[0] + 0x60}
	}
	// a-z -> ａ-ｚ
	if c.val[0] >= 0x61 && c.val[0] <= 0x7a {
		c.cval = []byte{0xef, 0xbd, c.val[0] + 0x20}
	}
}

func lowerN(c *character) {
	if c.filters&FLT_LOWER_N != FLT_LOWER_N {
		return
	}
	c.cval = []byte{c.val[2] - 0x60}
}

func upperN(c *character) {
	if c.filters&FLT_UPPER_N != FLT_UPPER_N {
		return
	}
	c.cval = []byte{0xef, 0xbc, c.val[0] + 0x60}
}

func lowerA(c *character) {
	if c.filters&FLT_LOWER_A != FLT_LOWER_A {
		return
	}
	if c.val[1] == 0xbc && c.val[2] >= 0x81 && c.val[2] <= 0xbf {
		c.cval = []byte{c.val[2] - 0x60}
	}
	if c.val[1] == 0xbd && c.val[2] >= 0x80 && c.val[2] <= 0x9d {
		c.cval = []byte{c.val[2] - 0x20}
	}
}

func upperA(c *character) {
	if c.filters&FLT_UPPER_A != FLT_UPPER_A {
		return
	}
	if c.val[0] >= 0x21 && c.val[0] <= 0x5f {
		c.cval = []byte{0xef, 0xbc, c.val[0] + 0x60}
	}
	if c.val[0] >= 0x60 && c.val[0] <= 0x7d {
		c.cval = []byte{0xef, 0xbd, c.val[0] + 0x20}
	}
}

func lowerS(c *character) {
	if c.filters&FLT_LOWER_S != FLT_LOWER_S {
		return
	}
	c.cval = []byte{0x20}
}

func upperS(c *character) {
	if c.filters&FLT_UPPER_S != FLT_UPPER_S {
		return
	}
	c.cval = []byte{0xe3, 0x80, 0x80}
}

func lowerK(c *character) {
	if c.filters&FLT_LOWER_K != FLT_LOWER_K {
		return
	}
	cval2 := byte(0x00)
	if c.val[1] == 0x80 {
		switch c.val[2] {
		case 0x81: /* 、 -> ､ */
			c.cval = []byte{0xef, 0xbd, 0xa4}
		case 0x82: /* 。 -> ｡ */
			c.cval = []byte{0xef, 0xbd, 0xa1}
		}
	} else if c.val[1] == 0x82 {
		switch c.val[2] {
		case 0x9b: /* ゛ ->  ﾞ */
			c.cval = []byte{0xef, 0xbe, 0x9e}
		case 0x9c: /* ゜ ->  ﾟ */
			c.cval = []byte{0xef, 0xbe, 0x9f}
		// ァ行: ァ -> ｧ, ィ -> ｨ, ゥ -> ｩ, ェ -> ｪ, ォ -> ｫ
		case 0xa1, 0xa3, 0xa5, 0xa7, 0xa9:
			cval2 = 0xa7 + (c.val[2]-0xa1)/0x02
			c.cval = []byte{0xef, 0xbd, cval2}
		// ア行: ア -> ｱ, イ -> ｲ, ウ -> ｳ, エ -> ｴ, オ -> ｵ
		case 0xa2, 0xa4, 0xa6, 0xa8, 0xaa:
			cval2 = 0xb1 + (c.val[2]-0xa2)/0x02
			c.cval = []byte{0xef, 0xbd, cval2}
		// カ行: カ -> ｶ, キ -> ｷ, ク -> ｸ, ケ -> ｹ, コ -> ｺ
		case 0xab, 0xad, 0xaf, 0xb1, 0xb3:
			cval2 = 0xb6 + (c.val[2]-0xab)/0x02
			c.cval = []byte{0xef, 0xbd, cval2}
		// ガ行: ガ -> ｶﾞ, ギ -> ｷﾞ, グ -> ｸﾞ, ゲ -> ｹﾞ, ゴ -> ｺﾞ
		case 0xac, 0xae, 0xb0, 0xb2, 0xb4:
			cval2 = 0xb6 + (c.val[2]-0xab)/0x02
			c.cval = []byte{0xef, 0xbd, cval2, 0xef, 0xbe, 0x9e}
		// サ行: サ -> ｻ, シ -> ｼ, ス -> ｽ, セ -> ｾ, ソ -> ｿ
		case 0xb5, 0xb7, 0xb9, 0xbb, 0xbd:
			cval2 = 0xbb + (c.val[2]-0xb5)/0x02
			c.cval = []byte{0xef, 0xbd, cval2}
		// ザ行: ザ -> ｻﾞ, ジ -> ｼﾞ, ズ -> ｽﾞ, ゼ -> ｾﾞ, ゾ -> ｿﾞ
		case 0xb6, 0xb8, 0xba, 0xbc, 0xbe:
			cval2 = 0xbb + (c.val[2]-0xb5)/0x02
			c.cval = []byte{0xef, 0xbd, cval2, 0xef, 0xbe, 0x9e}
			// タ行(タ): タ -> ﾀ
		case 0xbf:
			c.cval = []byte{0xef, 0xbe, 0x80}
		}
	} else if c.val[1] == 0x83 {
		switch c.val[2] {
		case 0x81:
			// タ行(チ): チ -> ﾁ
			c.cval = []byte{0xef, 0xbe, 0x81}
		// ダ行(ダ,ヂ): ダ -> ﾀﾞ, ヂ -> ﾁﾞ
		case 0x80, 0x82:
			cval2 = 0x80 + (c.val[2]-0x80)/0x02
			c.cval = []byte{0xef, 0xbe, cval2, 0xef, 0xbe, 0x9e}
		// タ行(ッ): ッ -> ｯ
		case 0x83:
			c.cval = []byte{0xef, 0xbd, 0xaf}
		// タ行(ツ,テ,ト): ツ -> ﾂ, テ -> ﾃ, ト -> ﾄ
		case 0x84, 0x86, 0x88:
			cval2 = 0x82 + (c.val[2]-0x84)/0x02
			c.cval = []byte{0xef, 0xbe, cval2}
		// ダ行(ヅ, デ, ド): ヅ -> ﾂﾞ, デ -> ﾃﾞ, ド -> ﾄﾞ
		case 0x85, 0x87, 0x89:
			cval2 = 0x82 + (c.val[2]-0x84)/0x02
			c.cval = []byte{0xef, 0xbe, cval2, 0xef, 0xbe, 0x9e}
		// ナ行: ナ -> ﾅ, ニ -> ﾆ, ヌ -> ﾇ, ネ -> ﾈ, ノ -> ﾉ
		case 0x8a, 0x8b, 0x8c, 0x8d, 0x8e:
			cval2 = c.val[2] - 0x05
			c.cval = []byte{0xef, 0xbe, cval2}
		// ハ行: ハ -> ﾊ, ヒ -> ﾋ, フ -> ﾌ, ヘ -> ﾍ, ホ -> ﾎ
		case 0x8f, 0x92, 0x95, 0x98, 0x9b:
			cval2 = 0x8a + (c.val[2]-0x8f)/0x03
			c.cval = []byte{0xef, 0xbe, cval2}
		// バ行: バ -> ﾊﾞ, ビ -> ﾋﾞ, ブ -> ﾌﾞ, ベ -> ﾍﾞ, ボ -> ﾎﾞ
		case 0x90, 0x93, 0x96, 0x99, 0x9c:
			cval2 = 0x8a + (c.val[2]-0x8f)/0x03
			c.cval = []byte{0xef, 0xbe, cval2, 0xef, 0xbe, 0x9e}
		// パ行: パ -> ﾊﾟ, ピ -> ﾋﾟ, プ -> ﾌﾟ, ペ -> ﾍﾟ, ポ -> ﾎﾟ
		case 0x91, 0x94, 0x97, 0x9a, 0x9d:
			cval2 = 0x8a + (c.val[2]-0x8f)/0x03
			c.cval = []byte{0xef, 0xbe, cval2, 0xef, 0xbe, 0x9f}
		// マ行: マ -> ﾏ, ミ -> ﾐ, ム -> ﾑ, メ -> ﾒ, モ -> ﾓ
		case 0x9e, 0x9f, 0xa0, 0xa1, 0xa2:
			cval2 = c.val[2] - 0x0f
			c.cval = []byte{0xef, 0xbe, cval2}
		// ャ行: ャ -> ｬ, ュ -> ｭ, ョ -> ｮ
		case 0xa3, 0xa5, 0xa7:
			cval2 = 0xac + (c.val[2]-0xa3)/0x02
			c.cval = []byte{0xef, 0xbd, cval2}
		// ヤ行: ヤ -> ﾔ, ユ -> ﾕ, ヨ -> ﾖ
		case 0xa4, 0xa6, 0xa8:
			cval2 = 0x94 + (c.val[2]-0xa4)/0x02
			c.cval = []byte{0xef, 0xbe, cval2}
		// ラ行: ラ -> ﾗ, リ -> ﾘ, ル -> ﾙ, レ -> ﾚ, ロ -> ﾛ
		case 0xa9, 0xaa, 0xab, 0xac, 0xad:
			cval2 = c.val[2] - 0x12
			c.cval = []byte{0xef, 0xbe, cval2}
		// ヮ -> ﾜ, ワ -> ﾜ
		case 0xae, 0xaf:
			c.cval = []byte{0xef, 0xbe, 0x9c}
		// ヰ -> ｲ
		case 0xb0:
			c.cval = []byte{0xef, 0xbd, 0xb2}
		// ヱ -> ｴ
		case 0xb1:
			c.cval = []byte{0xef, 0xbd, 0xb4}
		// ヲ -> ｦ
		case 0xb2:
			c.cval = []byte{0xef, 0xbd, 0xa6}
		// ン -> ﾝ
		case 0xb3:
			c.cval = []byte{0xef, 0xbe, 0x9d}
		// ヴ -> ｳﾞ
		case 0xb4:
			c.cval = []byte{0xef, 0xbd, 0xb3, 0xef, 0xbe, 0x9e}
		// ・ -> ･
		case 0xbb:
			c.cval = []byte{0xef, 0xbd, 0xa5}
		// ー -> ｰ
		case 0xbc:
			c.cval = []byte{0xef, 0xbd, 0xb0}
		}
	}
}

func upperK(c *character) {
	if c.filters&FLT_UPPER_K != FLT_UPPER_K {
		return
	}
	cval2 := byte(0x00)
	if c.val[1] == 0xbd {
		switch c.val[2] {
		// ｡ -> 。
		case 0xa1:
			c.cval = []byte{0xe3, 0x80, 0x82}
		// ｢ -> 「, ｣ -> 」
		case 0xa2, 0xa3:
			c.cval = []byte{0xe3, 0x80, c.val[2] - 0x16}
		// ､ -> 、
		case 0xa4:
			c.cval = []byte{0xe3, 0x80, 0x81}
		// ･ -> ・
		case 0xa5:
			c.cval = []byte{0xe3, 0x83, 0xbb}
		// ｦ -> ヲ
		case 0xa6:
			c.cval = []byte{0xe3, 0x83, 0xb2}
		// ｧ行: ｧ -> ァ, ｨ -> ィ, ｩ -> ゥ, ｪ -> ェ, ｫ -> ォ
		case 0xa7, 0xa8, 0xa9, 0xaa, 0xab:
			cval2 = 0xa1 + (c.val[2]-0xa7)*0x02
			c.cval = []byte{0xe3, 0x82, cval2}
		// ｬ行: ｬ -> ャ, ｭ -> ュ, ｮ -> ョ
		case 0xac, 0xad, 0xae:
			cval2 = 0xa3 + (c.val[2]-0xac)*0x02
			c.cval = []byte{0xe3, 0x83, cval2}
		// ｯ -> ッ
		case 0xaf:
			c.cval = []byte{0xe3, 0x83, 0x83}
		// ｰ -> ー
		case 0xb0:
			c.cval = []byte{0xe3, 0x83, 0xbc}
		// ｱ行: ｱ -> ア, ｲ -> イ, ｴ -> エ, ｵ -> オ
		case 0xb1, 0xb2, 0xb4, 0xb5:
			cval2 = 0xa2 + (c.val[2]-0xb1)*0x02
			c.cval = []byte{0xe3, 0x82, cval2}
		// ｱ行(ｳ,ｳﾞ): ｳ -> ウ, ｳﾞ -> ヴ
		case 0xb3:
			if len(c.val) > 5 && c.val[5] == 0x9e {
				c.cval = []byte{0xe3, 0x83, 0xb4}
			} else {
				c.cval = []byte{0xe3, 0x82, 0xa6}
			}
		// ｶ行: ｶ ->カ, ｷ -> キ, ｸ -> ク, ｹ -> ケ, ｺ -> コ
		// ｶﾞ行: ｶﾞ -> ガ, ｷﾞ -> ギ, ｸﾞ -> グ, ｹﾞ -> ゲ, ｺﾞ -> ゴ
		case 0xb6, 0xb7, 0xb8, 0xb9, 0xba:
			if len(c.val) > 5 && c.val[5] == 0x9e {
				cval2 = 0xac + (c.val[2]-0xb6)*0x02
			} else {
				cval2 = 0xab + (c.val[2]-0xb6)*0x02
			}
			c.cval = []byte{0xe3, 0x82, cval2}
		// ｻ行: ｻ -> サ, ｼ -> シ, ｽ -> ス, ｾ -> セ, ｿ -> ソ
		// ｻﾞ行: ｻﾞ -> ザ, ｼﾞ -> ジ, ｽﾞ -> ズ, ｾﾞ -> ゼ, ｿﾞ -> ゾ
		case 0xbb, 0xbc, 0xbd, 0xbe, 0xbf:
			if len(c.val) > 5 && c.val[5] == 0x9e {
				cval2 = 0xb6 + (c.val[2]-0xbb)*0x02
			} else {
				cval2 = 0xb5 + (c.val[2]-0xbb)*0x02
			}
			c.cval = []byte{0xe3, 0x82, cval2}
		}
	} else if c.val[1] == 0xbe {
		switch c.val[2] {
		// タ行(タ)・ダ行(ダ): ﾀ -> タ, ﾀﾞ -> ダ
		case 0x80:
			if len(c.val) > 5 && c.val[5] == 0x9e {
				c.cval = []byte{0xe3, 0x83, 0x80}
			} else {
				c.cval = []byte{0xe3, 0x82, 0xbf}
			}
		// タ行(チ)・ダ行(ヂ): ﾁ -> チ, ﾁﾞ -> ヂ
		case 0x81:
			if len(c.val) > 5 && c.val[5] == 0x9e {
				cval2 = 0x82
			} else {
				cval2 = 0x81
			}
			c.cval = []byte{0xe3, 0x83, cval2}
		// タ行(ツ,テ,ト)・ダ行(ヅ,デ,ド): ﾃ -> テ, ﾄ -> ト, ﾃﾞ -> デ, ﾄﾞ -> ド
		case 0x82, 0x83, 0x84:
			if len(c.val) > 5 && c.val[5] == 0x9e {
				cval2 = 0x85 + (c.val[2]-0x82)*0x02
			} else {
				cval2 = 0x84 + (c.val[2]-0x82)*0x02
			}
			c.cval = []byte{0xe3, 0x83, cval2}
		// ナ行: (ﾅ -> ナ, ﾆ -> ニ, ﾇ -> ヌ, ﾈ -> ネ, ﾉ -> ノ)
		case 0x85, 0x86, 0x87, 0x88, 0x89:
			cval2 = c.val[2] + 0x05
			c.cval = []byte{0xe3, 0x83, cval2}
		// ハ行: ﾊ -> ハ, ﾋ -> ヒ, ﾌ -> フ, ﾍ -> ヘ, ﾎ -> ホ
		// バ行: ﾊﾞ -> バ, ﾋﾞ -> ビ, ﾌﾞ -> ブ, ﾍﾞ -> ベ, ﾎﾞ -> ボ
		// パ行: ﾊﾟ -> パ, ﾋﾟ -> ピ, ﾌﾟ -> プ, ﾍﾟ -> ペ, ﾎﾟ -> ポ,
		case 0x8a, 0x8b, 0x8c, 0x8d, 0x8e:
			if len(c.val) > 5 && c.val[5] == 0x9e {
				cval2 = 0x8a + (c.val[2]-0x88)*0x03
			} else if len(c.val) > 5 && c.val[5] == 0x9f {
				cval2 = 0x8b + (c.val[2]-0x88)*0x03
			} else {
				cval2 = 0x89 + (c.val[2]-0x88)*0x03
			}
			c.cval = []byte{0xe3, 0x83, cval2}
		// マ行: ﾏ -> マ, ﾐ -> ミ, ﾑ -> ム, ﾒ -> メ, ﾓ -> モ
		case 0x8f, 0x90, 0x91, 0x92, 0x93:
			cval2 = c.val[2] + 0x0f
			c.cval = []byte{0xe3, 0x83, cval2}
		// ヤ行: ﾔ -> ヤ, ﾕ -> ユ, ﾖ -> ヨ
		case 0x94, 0x95, 0x96:
			cval2 = 0xa4 + (c.val[2]-0x94)*0x02
			c.cval = []byte{0xe3, 0x83, cval2}
		case 0x97, 0x98, 0x99, 0x9a, 0x9b:
			cval2 = c.val[2] + 0x12
			c.cval = []byte{0xe3, 0x83, cval2}
		// ﾜ -> ワ
		case 0x9c:
			c.cval = []byte{0xe3, 0x83, 0xaf}
		// ﾝ -> ン
		case 0x9d:
			c.cval = []byte{0xe3, 0x83, 0xb3}
		case 0x9e: /* ﾞ */
			c.cval = []byte{0xe3, 0x82, 0x9b}
		case 0x9f: /* ﾟ */
			c.cval = []byte{0xe3, 0x82, 0x9c}
		}
	}
}

func lowerH(c *character) {
	if c.filters&FLT_LOWER_H != FLT_LOWER_H {
		return
	}
	cval2 := byte(0x00)
	if c.val[1] == 0x80 {
		switch c.val[2] {
		// 、 -> ､
		case 0x81:
			c.cval = []byte{0xef, 0xbd, 0xa4}
		// 。 -> ｡
		case 0x82:
			c.cval = []byte{0xef, 0xbd, 0xa1}
		}
	} else if c.val[1] == 0x81 {
		switch c.val[2] {
		// ぁ行: ぁ -> ｧ, ぃ -> ｨ, ぅ -> ｩ, ぇ -> ｪ, ぉ -> ｫ
		case 0x81, 0x83, 0x85, 0x87, 0x89:
			cval2 := 0xa7 + (c.val[2]-0x81)/0x02
			c.cval = []byte{0xef, 0xbd, cval2}
		// あ行: あ -> ｱ, い -> ｲ, う -> ｳ, え -> ｴ, お -> ｵ
		case 0x82, 0x84, 0x86, 0x88, 0x8a:
			cval2 := 0xb1 + (c.val[2]-0x82)/0x02
			c.cval = []byte{0xef, 0xbd, cval2}
		// か行: か -> ｶ, き -> ｷ, く -> ｸ, け -> ｹ, こ -> ｺ
		case 0x8b, 0x8d, 0x8f, 0x91, 0x93:
			cval2 := 0xb6 + (c.val[2]-0x8b)/0x02
			c.cval = []byte{0xef, 0xbd, cval2}
		// が行: が -> ｶﾞ, ぎ -> ｷﾞ, ぐ -> ｸﾞ, げ -> ｹﾞ, ご -> ｺﾞ
		case 0x8c, 0x8e, 0x90, 0x92, 0x94:
			cval2 = 0xb6 + (c.val[2]-0x8c)/0x02
			c.cval = []byte{0xef, 0xbd, cval2, 0xef, 0xbe, 0x9e}
		// さ行: さ -> ｻ, し -> ｼ, す -> ｽ, せ -> ｾ, そ -> ｿ
		case 0x95, 0x97, 0x99, 0x9b, 0x9d:
			cval2 = 0xbb + (c.val[2]-0x95)/0x02
			c.cval = []byte{0xef, 0xbd, cval2}
		// ざ行: ざ -> ｻﾞ, じ -> ｼﾞ, ず -> ｽﾞ, ぜ -> ｾﾞ, ぞ -> ｿﾞ
		case 0x96, 0x98, 0x9a, 0x9c, 0x9e:
			cval2 = 0xbb + (c.val[2]-0x96)/0x02
			c.cval = []byte{0xef, 0xbd, cval2, 0xef, 0xbe, 0x9e}
		// た行(た): た -> ﾀ
		case 0x9f:
			c.cval = []byte{0xef, 0xbe, 0x80}
		case 0xa1:
			// た行(ち): ち -> ﾁ
			c.cval = []byte{0xef, 0xbe, 0x81}
		// だ行(だ,ぢ): だ -> ﾀﾞ, ぢ -> ﾁﾞ
		case 0xa0, 0xa2:
			cval2 = 0x80 + (c.val[2]-0xa0)/0x02
			c.cval = []byte{0xef, 0xbe, cval2, 0xef, 0xbe, 0x9e}
		// た行(っ): っ -> ｯ
		case 0xa3:
			c.cval = []byte{0xef, 0xbd, 0xaf}
		// た行(つ,て,と): つ -> ﾂ, て -> ﾃ, と -> ﾄ
		case 0xa4, 0xa6, 0xa8:
			cval2 = 0x82 + (c.val[2]-0xa4)/0x02
			c.cval = []byte{0xef, 0xbe, cval2}
		// だ行(づ,で,ど): づ -> ﾂﾞ, で -> ﾃﾞ, ど -> ﾄﾞ
		case 0xa5, 0xa7, 0xa9:
			cval2 = 0x82 + (c.val[2]-0xa4)/0x02
			c.cval = []byte{0xef, 0xbe, cval2, 0xef, 0xbe, 0x9e}
		// な行: な -> ﾅ, に -> ﾆ, ぬ -> ﾇ, ね -> ﾈ, の -> ﾉ
		case 0xaa, 0xab, 0xac, 0xad, 0xae:
			cval2 = c.val[2] - 0x25
			c.cval = []byte{0xef, 0xbe, cval2}

		// は行: は -> ﾊ, ひ -> ﾋ, ふ -> ﾌ, へ -> ﾍ, ほ -> ﾎ
		case 0xaf, 0xb2, 0xb5, 0xb8, 0xbb:
			cval2 = 0x8a + (c.val[2]-0xaf)/0x03
			c.cval = []byte{0xef, 0xbe, cval2}
			// ば行: ば -> ﾊﾞ, び -> ﾋﾞ, ぶ -> ﾌﾞ, べ -> ﾍﾞ, ぼ -> ﾎﾞ
		case 0xb0, 0xb3, 0xb6, 0xb9, 0xbc:
			cval2 = 0x8a + (c.val[2]-0xaf)/0x03
			c.cval = []byte{0xef, 0xbe, cval2, 0xef, 0xbe, 0x9e}
		// ぱ行: ぱ -> ﾊﾟ, ぴ -> ﾋﾟ, ぷ -> ﾌﾟ, ぺ -> ﾍﾟ, ぽ -> ﾎﾟ
		case 0xb1, 0xb4, 0xb7, 0xba, 0xbd:
			cval2 = 0x8a + (c.val[2]-0xaf)/0x03
			c.cval = []byte{0xef, 0xbe, cval2, 0xef, 0xbe, 0x9f}
		// ま行(ま,み): ま -> ﾏ, み -> ﾐ
		case 0xbe, 0xbf:
			cval2 = c.val[2] - 0x2f
			c.cval = []byte{0xef, 0xbe, cval2}
		}
	} else if c.val[1] == 0x82 {
		switch c.val[2] {
		// ま行(む,め,も): む -> ﾑ, め -> ﾒ, も -> ﾓ
		case 0x80, 0x81, 0x82:
			cval2 = c.val[2] + 0x11
			c.cval = []byte{0xef, 0xbe, cval2}
		// ゃ行: ゃ -> ｬ, ゅ -> ｭ, ょ -> ｮ
		case 0x83, 0x85, 0x87:
			cval2 = 0xac + (c.val[2]-0x83)/0x02
			c.cval = []byte{0xef, 0xbd, cval2}
		// や行: や -> ﾔ, ゆ -> ﾕ, よ -> ﾖ
		case 0x84, 0x86, 0x88:
			cval2 = 0x94 + (c.val[2]-0x84)/0x02
			c.cval = []byte{0xef, 0xbe, cval2}
		// ら行: ら -> ﾗ, り -> ﾘ, る -> ﾙ, れ -> ﾚ, ろ -> ﾛ
		case 0x89, 0x8a, 0x8b, 0x8c, 0x8d:
			cval2 = c.val[2] + 0x0e
			c.cval = []byte{0xef, 0xbe, cval2}
		// ゎ -> ﾜ, わ -> ﾜ
		case 0x8e, 0x8f:
			c.cval = []byte{0xef, 0xbe, 0x9c}
		// ゐ -> ｲ
		case 0x90:
			c.cval = []byte{0xef, 0xbd, 0xb2}
		// ゑ -> ｴ
		case 0x91:
			c.cval = []byte{0xef, 0xbd, 0xb4}
		// を -> ｦ
		case 0x92:
			c.cval = []byte{0xef, 0xbd, 0xa6}
		// ン -> ﾝ
		case 0x93:
			c.cval = []byte{0xef, 0xbe, 0x9d}
		// ゛ -> ﾞ
		case 0x9b:
			c.cval = []byte{0xef, 0xbe, 0x9e}
		// ゜ -> ﾟ
		case 0x9c:
			c.cval = []byte{0xef, 0xbe, 0x9f}

		}
	} else if c.val[1] == 0x83 {
		switch c.val[2] {
		// ・ -> ･
		case 0xbb:
			c.cval = []byte{0xef, 0xbd, 0xa5}
		// ー -> ｰ
		case 0xbc:
			c.cval = []byte{0xef, 0xbd, 0xb0}
		}
	}
}

func upperH(c *character) {
	if c.filters&FLT_UPPER_H != FLT_UPPER_H {
		return
	}
	cval2 := byte(0x00)
	if c.val[1] == 0xbd {
		switch c.val[2] {
		// ｡ -> 。
		case 0xa1:
			c.cval = []byte{0xe3, 0x80, 0x82}
		// ｢ -> 「, ｣ -> 」
		case 0xa2, 0xa3:
			cval2 = c.val[2] - 0x16
			c.cval = []byte{0xe3, 0x80, cval2}
		// ､ -> 、
		case 0xa4:
			c.cval = []byte{0xe3, 0x80, 0x81}
		// ･ -> ・
		case 0xa5:
			c.cval = []byte{0xe3, 0x83, 0xbb}
		// ｦ -> を
		case 0xa6:
			c.cval = []byte{0xe3, 0x82, 0x92}
		// ｧ行: ｧ -> ぁ, ｨ -> ぃ, ｩ -> ぅ, ｪ -> ぇ, ｫ -> ぉ
		case 0xa7, 0xa8, 0xa9, 0xaa, 0xab:
			cval2 = 0x81 + (c.val[2]-0xa7)*0x02
			c.cval = []byte{0xe3, 0x81, cval2}
		// ｬ行: ｬ -> ゃ, ｭ -> ゅ, ｮ -> ょ
		case 0xac, 0xad, 0xae:
			cval2 = 0x83 + (c.val[2]-0xac)*0x02
			c.cval = []byte{0xe3, 0x82, cval2}
		// ﾀ行(ｯ): ｯ -> っ
		case 0xaf:
			c.cval = []byte{0xe3, 0x81, 0xa3}
		// ｰ -> ー
		case 0xb0: /* ｰ */
			c.cval = []byte{0xe3, 0x83, 0xbc}
		// ｱ行: ｱ -> あ, ｲ -> い, ｳ -> う, ｴ -> え, ｵ -> お
		case 0xb1, 0xb2, 0xb4, 0xb5:
			cval2 = 0x82 + (c.val[2]-0xb1)*0x02
			c.cval = []byte{0xe3, 0x81, cval2}
		case 0xb3:
			c.cval = []byte{0xe3, 0x81, 0x86}
			if len(c.val) > 5 && c.val[5] == 0x9e {
				c.cval = append(c.cval, []byte{0xe3, 0x82, 0x9b}...)
			}
		// ｶ行: ｶ -> か, ｷ -> き, ｸ -> く, ｹ -> け, ｺ -> こ
		// ｶﾞ行: ｶﾞ -> が, ｷﾞ -> ぎ, ｸﾞ -> ぐ, ｹﾞ -> げ, ｺﾞ -> ご
		case 0xb6, 0xb7, 0xb8, 0xb9, 0xba:
			if len(c.val) > 5 && c.val[5] == 0x9e {
				cval2 = 0x8c + (c.val[2]-0xb6)*0x02
			} else {
				cval2 = 0x8b + (c.val[2]-0xb6)*0x02
			}
			c.cval = []byte{0xe3, 0x81, cval2}
		// ｻ行: ｻ -> さ, ｼ -> し, ｽ -> す, ｾ -> せ, ｿ -> そ
		// ｻﾞ行: ｻﾞ -> ざ, ｼﾞ -> じ, ｽﾞ -> ず, ｾﾞ -> ぜ, ｿﾞ -> ぞ
		case 0xbb, 0xbc, 0xbd, 0xbe, 0xbf:
			if len(c.val) > 5 && c.val[5] == 0x9e {
				cval2 = 0x96 + (c.val[2]-0xbb)*0x02
			} else {
				cval2 = 0x95 + (c.val[2]-0xbb)*0x02
			}
			c.cval = []byte{0xe3, 0x81, cval2}
		}
	} else if c.val[1] == 0xbe {
		switch c.val[2] {
		// ﾀ行(ﾀ,ﾁ): ﾀ -> た, ﾁ -> ち
		// ﾀﾞ行(ﾀﾞ,ﾁﾞ): ﾀﾞ -> だ, ﾁﾞ -> ぢ
		case 0x80, 0x81:
			if len(c.val) > 5 && c.val[5] == 0x9e {
				cval2 = 0xa0 + (c.val[2]-0x80)*0x02
			} else {
				cval2 = 0x9f + (c.val[2]-0x80)*0x02
			}
			c.cval = []byte{0xe3, 0x81, cval2}
		// ﾀ行(ﾂ,ﾃ,ﾄ): ﾂ -> つ, ﾃ -> て, ﾄ -> と
		// ﾀﾞ行(ﾂﾞ,ﾃﾞ,ﾄﾞ): ﾂﾞ -> づ, ﾃﾞ -> で, ﾄﾞ -> ど
		case 0x82, 0x83, 0x84:
			if len(c.val) > 5 && c.val[5] == 0x9e {
				cval2 = 0xa5 + (c.val[2]-0x82)*0x02
			} else {
				cval2 = 0xa4 + (c.val[2]-0x82)*0x02
			}
			c.cval = []byte{0xe3, 0x81, cval2}
		// ﾅ行: ﾅ -> な, ﾆ -> に, ﾇ -> ぬ, ﾈ -> ね, ﾉ -> の
		case 0x85, 0x86, 0x87, 0x88, 0x89:
			cval2 = c.val[2] + 0x25
			c.cval = []byte{0xe3, 0x81, cval2}
			// ﾊ行: ﾊ -> は, ﾋ -> ひ, ﾌ -> ふ, ﾍ -> へ, ﾎ -> ほ
			// ﾊﾞ行: ﾊﾞ -> ば, ﾋﾞ -> び, ﾌﾞ -> ぶ, ﾍﾞ -> べ, ﾎﾞ -> ぼ
			// ﾊﾟ行: ﾊﾟ -> ぱ, ﾋﾟ -> ぴ, ﾌﾟ -> ぷ, ﾍﾟ -> ぺ, ﾎﾟ -> ぽ
		case 0x8a, 0x8b, 0x8c, 0x8d, 0x8e:
			if len(c.val) > 5 && c.val[5] == 0x9e {
				cval2 = 0xb0 + (c.val[2]-0x8a)*0x03
			} else if len(c.val) > 5 && c.val[5] == 0x9f {
				cval2 = 0xb1 + (c.val[2]-0x8a)*0x03
			} else {
				cval2 = 0xaf + (c.val[2]-0x8a)*0x03
			}
			c.cval = []byte{0xe3, 0x81, cval2}
		// ﾏ行: ﾏ -> ま, ﾐ -> み
		case 0x8f, 0x90:
			cval2 = c.val[2] + 0x2f
			c.cval = []byte{0xe3, 0x81, cval2}
		// ﾏ行: ﾑ -> む, ﾒ -> め, ﾓ -> も
		case 0x91, 0x92, 0x93:
			cval2 = c.val[2] - 0x11
			c.cval = []byte{0xe3, 0x82, cval2}
		// ﾔ行: ﾔ -> や, ﾕ -> ゆ, ﾖ -> よ
		case 0x94, 0x95, 0x96:
			cval2 = 0x84 + (c.val[2]-0x94)*0x02
			c.cval = []byte{0xe3, 0x82, cval2}
		// ﾗ行: ﾗ -> ら, ﾘ -> り, ﾙ -> る, ﾚ -> れ, ﾛ -> ろ
		case 0x97, 0x98, 0x99, 0x9a, 0x9b:
			cval2 = c.val[2] - 0x0e
			c.cval = []byte{0xe3, 0x82, cval2}
		// ﾜ行: ﾜ -> わ
		case 0x9c:
			c.cval = []byte{0xe3, 0x82, 0x8f}
		// ﾝ -> ん
		case 0x9d:
			c.cval = []byte{0xe3, 0x82, 0x93}
		// ﾞ -> ゛
		case 0x9e:
			c.cval = []byte{0xe3, 0x82, 0x9b}
		// ﾟ -> ゜
		case 0x9f:
			c.cval = []byte{0xe3, 0x82, 0x9c}
		}
	}
}

func lowerC(c *character) {
	if c.filters&FLT_LOWER_C != FLT_LOWER_C {
		return
	}
	switch c.val[1] {
	case 0x82: // ァ - タ
		if c.val[2] >= 0xa1 && c.val[2] <= 0xbf {
			c.cval = []byte{0xe3, 0x81, c.val[2] - 0x20}
		}
	case 0x83:
		if c.val[2] >= 0x80 && c.val[2] <= 0x9f { // ダ - ミ
			c.cval = []byte{0xe3, 0x81, c.val[2] + 0x20}
		} else if c.val[2] >= 0xa0 && c.val[2] <= 0xb3 { // ム - ン
			c.cval = []byte{0xe3, 0x82, c.val[2] - 0x20}
		} else if c.val[2] >= 0xbd && c.val[2] <= 0xbe { // ヽヾ
			c.cval = []byte{0xe3, 0x82, c.val[2] - 0x20}
		}
	}
}

func upperC(c *character) {
	if c.filters&FLT_UPPER_C != FLT_UPPER_C {
		return
	}
	switch c.val[1] {
	case 0x81:
		if c.val[2] >= 0x81 && c.val[2] <= 0x9f { // ぁ - た
			c.cval = []byte{0xe3, 0x82, c.val[2] + 0x20}
		} else if c.val[2] >= 0xa0 && c.val[2] <= 0xbf { // だ - み
			c.cval = []byte{0xe3, 0x83, c.val[2] - 0x20}
		}
	case 0x82:
		if c.val[2] >= 0x80 && c.val[2] <= 0x93 { // む - ん
			c.cval = []byte{0xe3, 0x83, c.val[2] + 0x20}
		} else if c.val[2] >= 0x9d && c.val[2] <= 0x9e { // ゝゞ
			c.cval = []byte{0xe3, 0x83, c.val[2] + 0x20}
		}
	}
}

func asis(c *character) {
	c.cval = make([]byte, len(c.val))
	copy(c.cval, c.val)
}

func extract(c *character, s []byte) {
	length := 1
	if is1Byte(s) {
		c0 := s[0]
		if c0 == 0x20 { // Space
			c.filters = FLT_UPPER_S
		} else if c0 >= 0x30 && c0 <= 0x39 { // 0 - 9
			c.filters = FLT_UPPER_A | FLT_UPPER_N
		} else if c0 >= 0x41 && c0 <= 0x5a { // A - Z
			c.filters = FLT_UPPER_A | FLT_UPPER_R
		} else if c0 >= 0x61 && c0 <= 0x7a { // a - z
			c.filters = FLT_UPPER_A | FLT_UPPER_R
		} else if c0 >= 0x21 && c0 <= 0x7d && c0 != 0x22 && c0 != 0x27 && c0 != 0x5c {
			c.filters = FLT_UPPER_A
		}
	} else if is3Bytes(s) {
		length = 3
		c0, c1, c2 := s[0], s[1], s[2]
		if c0 == 0xef {
			if c1 == 0xbc {
				if c2 >= 0x90 && c2 <= 0x99 { // ０ - ９
					c.filters = FLT_LOWER_A | FLT_LOWER_N
				} else if c2 >= 0xa1 && c2 <= 0xba { // Ａ - Ｚ
					c.filters = FLT_LOWER_A | FLT_LOWER_R
				} else if c2 != 0x82 && c2 != 0x87 && c2 != 0xbc { // except ＂ ＇ ＼
					c.filters = FLT_LOWER_A
				}
			} else if c1 == 0xbd {
				if c2 >= 0x81 && c2 <= 0x9a { // ａ - ｚ
					c.filters = FLT_LOWER_A | FLT_LOWER_R
				} else if c2 >= 0x80 && c2 <= 0x9d { // ｀ - ｝
					c.filters = FLT_LOWER_A
				} else if c2 >= 0xa1 && c2 <= 0xbf { // ｡ ｢ ｣ ､ ･ ｦ - ｿ
					if isVoiced(s) { // voiced
						length = 6
						if c2 == 0xb3 {
							c.filters = FLT_UPPER_H | FLT_UPPER_K
						} else {
							c.filters = FLT_UPPER_H | FLT_UPPER_K
						}
					} else {
						c.filters = FLT_UPPER_H | FLT_UPPER_K
					}
				}
			} else if c1 == 0xbe {
				if c2 >= 0x80 && c2 <= 0x84 { // ﾀ - ﾄ
					if isVoiced(s) {
						length = 6
					}
					c.filters = FLT_UPPER_H | FLT_UPPER_K
				} else if c2 >= 0x8a && c2 <= 0x8e { // ﾊ - ﾎ
					if isVoiced(s) || isSemiVoiced(s) { // voiced or semi voiced
						length = 6
					}
					c.filters = FLT_UPPER_H | FLT_UPPER_K
				} else if c2 >= 0x85 && c2 <= 0x9f { // ﾅ - ﾝﾞﾟ
					c.filters = FLT_UPPER_H | FLT_UPPER_K
				}
			}
		} else if c0 == 0xe3 {
			if c1 == 0x80 {
				if c2 == 0x80 { // Space
					c.filters = FLT_LOWER_S
				} else if c2 >= 0x81 && c2 <= 0x82 { // 、。
					c.filters = FLT_LOWER_H | FLT_LOWER_K
				}
			} else if c1 == 0x81 {
				if c2 >= 0x81 && c2 <= 0xbf { // ぁ - み
					c.filters = FLT_UPPER_C | FLT_LOWER_H
				}
			} else if c1 == 0x82 {
				if c2 >= 0x80 && c2 <= 0x93 { // む - ん
					c.filters = FLT_UPPER_C | FLT_LOWER_H
				} else if c2 >= 0x9b && c2 <= 0x9c { // ゛゜
					c.filters = FLT_LOWER_H | FLT_LOWER_K
				} else if c2 >= 0x9d && c2 <= 0x9e { // ゝゞ
					c.filters = FLT_UPPER_C
				} else if c2 >= 0xa1 && c2 <= 0xbf { // ァ - タ
					c.filters = FLT_LOWER_C | FLT_LOWER_K
				}
			} else if c1 == 0x83 {
				if c2 >= 0x80 && c2 <= 0xb3 { // チ - ン
					c.filters = FLT_LOWER_C | FLT_LOWER_K
				} else if c2 == 0xb4 { // ヴ
					c.filters = FLT_LOWER_K
				} else if c2 >= 0xbb && c2 <= 0xbc { // ・ー
					c.filters = FLT_LOWER_H | FLT_LOWER_K
				} else if c2 >= 0xbd && c2 <= 0xbe { // ヽ ヾ
					c.filters = FLT_LOWER_C
				}
			}
		}
	} else if is2Bytes(s) {
		length = 2
	} else if is4Bytes(s) {
		length = 4
	} else {
		length = 1
	}
	c.val = s[:length]
}

func conv(c *character, filters []filter) []byte {
	for _, f := range filters {
		f(c)
	}
	if len(c.cval) == 0 {
		asis(c)
	}
	return c.cval
}

func createFilters(mode string) []filter {
	filters := []filter{}
	exists := map[byte]bool{}
	for _, m := range []byte(mode) {
		if _, ok := exists[m]; ok {
			continue
		}
		switch m {
		case 'r':
			filters = append(filters, lowerR)
			exists[m] = true
		case 'R':
			filters = append(filters, upperR)
			exists[m] = true
		case 'n':
			filters = append(filters, lowerN)
			exists[m] = true
		case 'N':
			filters = append(filters, upperN)
			exists[m] = true
		case 'a':
			filters = append(filters, lowerA)
			exists[m] = true
		case 'A':
			filters = append(filters, upperA)
			exists[m] = true
		case 's':
			filters = append(filters, lowerS)
			exists[m] = true
		case 'S':
			filters = append(filters, upperS)
			exists[m] = true
		case 'k':
			filters = append(filters, lowerK)
			exists[m] = true
		case 'K':
			filters = append(filters, upperK)
			exists[m] = true
		case 'h':
			filters = append(filters, lowerH)
			exists[m] = true
		case 'H':
			filters = append(filters, upperH)
			exists[m] = true
		case 'c':
			filters = append(filters, lowerC)
			exists[m] = true
		case 'C':
			filters = append(filters, upperC)
			exists[m] = true
		}
	}
	return filters
}

package kanaco

import (
	"bufio"
	"bytes"
	"io"
)

const (
	hankaku  = 1
	zenkaku  = 2
	space    = 4
	numeric  = 8
	alphabet = 16
	kana     = 32
	voiced   = 64
	devoiced = 128
	extra    = 4096
	binary   = 8192
)

type Reader struct {
	r    *bufio.Reader
	mode string
}

type Writer struct {
	w    *bufio.Writer
	mode string
}

type word struct {
	val      []byte
	charType int
	len      int
}

var tbl map[string][]byte = map[string][]byte{
	"ぁ": {0xe3, 0x81, 0x81}, "あ": {0xe3, 0x81, 0x82},
	"ぃ": {0xe3, 0x81, 0x83}, "い": {0xe3, 0x81, 0x84},
	"ぅ": {0xe3, 0x81, 0x85}, "う": {0xe3, 0x81, 0x86},
	"ぇ": {0xe3, 0x81, 0x87}, "え": {0xe3, 0x81, 0x88},
	"ぉ": {0xe3, 0x81, 0x89}, "お": {0xe3, 0x81, 0x8a},
	"か": {0xe3, 0x81, 0x8b}, "が": {0xe3, 0x81, 0x8c},
	"き": {0xe3, 0x81, 0x8d}, "ぎ": {0xe3, 0x81, 0x8e},
	"く": {0xe3, 0x81, 0x8f}, "ぐ": {0xe3, 0x81, 0x90},
	"け": {0xe3, 0x81, 0x91}, "げ": {0xe3, 0x81, 0x92},
	"こ": {0xe3, 0x81, 0x93}, "ご": {0xe3, 0x81, 0x94},
	"さ": {0xe3, 0x81, 0x95}, "ざ": {0xe3, 0x81, 0x96},
	"し": {0xe3, 0x81, 0x97}, "じ": {0xe3, 0x81, 0x98},
	"す": {0xe3, 0x81, 0x99}, "ず": {0xe3, 0x81, 0x9a},
	"せ": {0xe3, 0x81, 0x9b}, "ぜ": {0xe3, 0x81, 0x9c},
	"そ": {0xe3, 0x81, 0x9d}, "ぞ": {0xe3, 0x81, 0x9e},
	"た": {0xe3, 0x81, 0x9f}, "だ": {0xe3, 0x81, 0xa0},
	"ち": {0xe3, 0x81, 0xa1}, "ぢ": {0xe3, 0x81, 0xa2},
	"っ": {0xe3, 0x81, 0xa3}, "つ": {0xe3, 0x81, 0xa4}, "づ": {0xe3, 0x81, 0xa5},
	"て": {0xe3, 0x81, 0xa6}, "で": {0xe3, 0x81, 0xa7},
	"と": {0xe3, 0x81, 0xa8}, "ど": {0xe3, 0x81, 0xa9},
	"な": {0xe3, 0x81, 0xaa},
	"に": {0xe3, 0x81, 0xab},
	"ぬ": {0xe3, 0x81, 0xac},
	"ね": {0xe3, 0x81, 0xad},
	"の": {0xe3, 0x81, 0xae},
	"は": {0xe3, 0x81, 0xaf}, "ば": {0xe3, 0x81, 0xb0}, "ぱ": {0xe3, 0x81, 0xb1},
	"ひ": {0xe3, 0x81, 0xb2}, "び": {0xe3, 0x81, 0xb3}, "ぴ": {0xe3, 0x81, 0xb4},
	"ふ": {0xe3, 0x81, 0xb5}, "ぶ": {0xe3, 0x81, 0xb6}, "ぷ": {0xe3, 0x81, 0xb7},
	"へ": {0xe3, 0x81, 0xb8}, "べ": {0xe3, 0x81, 0xb9}, "ぺ": {0xe3, 0x81, 0xba},
	"ほ": {0xe3, 0x81, 0xbb}, "ぼ": {0xe3, 0x81, 0xbc}, "ぽ": {0xe3, 0x81, 0xbd},
	"ま": {0xe3, 0x81, 0xbe},
	"み": {0xe3, 0x81, 0xbf},
	"む": {0xe3, 0x82, 128},
	"め": {0xe3, 0x82, 0x81},
	"も": {0xe3, 0x82, 0x82},
	"ゃ": {0xe3, 0x82, 0x83}, "や": {0xe3, 0x82, 0x84},
	"ゅ": {0xe3, 0x82, 0x85}, "ゆ": {0xe3, 0x82, 0x86},
	"ょ": {0xe3, 0x82, 0x87}, "よ": {0xe3, 0x82, 0x88},
	"ら": {0xe3, 0x82, 0x89},
	"り": {0xe3, 0x82, 0x8a},
	"る": {0xe3, 0x82, 0x8b},
	"れ": {0xe3, 0x82, 0x8c},
	"ろ": {0xe3, 0x82, 0x8d},
	"ゎ": {0xe3, 0x82, 0x8e},
	"わ": {0xe3, 0x82, 0x8f},
	"ゐ": {0xe3, 0x82, 0x90},
	"ゑ": {0xe3, 0x82, 0x91},
	"を": {0xe3, 0x82, 0x92},
	"ん": {0xe3, 0x82, 0x93},
	"ァ": {0xe3, 0x82, 0xa1}, "ア": {0xe3, 0x82, 0xa2},
	"ィ": {0xe3, 0x82, 0xa3}, "イ": {0xe3, 0x82, 0xa4},
	"ゥ": {0xe3, 0x82, 0xa5}, "ウ": {0xe3, 0x82, 0xa6},
	"ェ": {0xe3, 0x82, 0xa7}, "エ": {0xe3, 0x82, 0xa8},
	"ォ": {0xe3, 0x82, 0xa9}, "オ": {0xe3, 0x82, 0xaa},
	"カ": {0xe3, 0x82, 0xab}, "ガ": {0xe3, 0x82, 0xac},
	"キ": {0xe3, 0x82, 0xad}, "ギ": {0xe3, 0x82, 0xae},
	"ク": {0xe3, 0x82, 0xaf}, "グ": {0xe3, 0x82, 0xb0},
	"ケ": {0xe3, 0x82, 0xb1}, "ゲ": {0xe3, 0x82, 0xb2},
	"コ": {0xe3, 0x82, 0xb3}, "ゴ": {0xe3, 0x82, 0xb4},
	"サ": {0xe3, 0x82, 0xb5}, "ザ": {0xe3, 0x82, 0xb6},
	"シ": {0xe3, 0x82, 0xb7}, "ジ": {0xe3, 0x82, 0xb8},
	"ス": {0xe3, 0x82, 0xb9}, "ズ": {0xe3, 0x82, 0xba},
	"セ": {0xe3, 0x82, 0xbb}, "ゼ": {0xe3, 0x82, 0xbc},
	"ソ": {0xe3, 0x82, 0xbd}, "ゾ": {0xe3, 0x82, 0xbe},
	"タ": {0xe3, 0x82, 0xbf}, "ダ": {0xe3, 0x83, 128},
	"チ": {0xe3, 0x83, 0x81}, "ヂ": {0xe3, 0x83, 0x82},
	"ッ": {0xe3, 0x83, 0x83}, "ツ": {0xe3, 0x83, 0x84}, "ヅ": {0xe3, 0x83, 0x85},
	"テ": {0xe3, 0x83, 0x86}, "デ": {0xe3, 0x83, 0x87},
	"ト": {0xe3, 0x83, 0x88}, "ド": {0xe3, 0x83, 0x89},
	"ナ": {0xe3, 0x83, 0x8a},
	"ニ": {0xe3, 0x83, 0x8b},
	"ヌ": {0xe3, 0x83, 0x8c},
	"ネ": {0xe3, 0x83, 0x8d},
	"ノ": {0xe3, 0x83, 0x8e},
	"ハ": {0xe3, 0x83, 0x8f}, "バ": {0xe3, 0x83, 0x90}, "パ": {0xe3, 0x83, 0x91},
	"ヒ": {0xe3, 0x83, 0x92}, "ビ": {0xe3, 0x83, 0x93}, "ピ": {0xe3, 0x83, 0x94},
	"フ": {0xe3, 0x83, 0x95}, "ブ": {0xe3, 0x83, 0x96}, "プ": {0xe3, 0x83, 0x97},
	"ヘ": {0xe3, 0x83, 0x98}, "ベ": {0xe3, 0x83, 0x99}, "ペ": {0xe3, 0x83, 0x9a},
	"ホ": {0xe3, 0x83, 0x9b}, "ボ": {0xe3, 0x83, 0x9c}, "ポ": {0xe3, 0x83, 0x9d},
	"マ": {0xe3, 0x83, 0x9e},
	"ミ": {0xe3, 0x83, 0x9f},
	"ム": {0xe3, 0x83, 0xa0},
	"メ": {0xe3, 0x83, 0xa1},
	"モ": {0xe3, 0x83, 0xa2},
	"ャ": {0xe3, 0x83, 0xa3}, "ヤ": {0xe3, 0x83, 0xa4},
	"ュ": {0xe3, 0x83, 0xa5}, "ユ": {0xe3, 0x83, 0xa6},
	"ョ": {0xe3, 0x83, 0xa7}, "ヨ": {0xe3, 0x83, 0xa8},
	"ラ": {0xe3, 0x83, 0xa9},
	"リ": {0xe3, 0x83, 0xaa},
	"ル": {0xe3, 0x83, 0xab},
	"レ": {0xe3, 0x83, 0xac},
	"ロ": {0xe3, 0x83, 0xad},
	"ヮ": {0xe3, 0x83, 0xae}, "ワ": {0xe3, 0x83, 0xaf},
	"ヰ": {0xe3, 0x83, 0xb0},
	"ヱ": {0xe3, 0x83, 0xb1},
	"ヲ": {0xe3, 0x83, 0xb2},
	"ン": {0xe3, 0x83, 0xb3},
	"ー": {0xe3, 0x83, 0xbc},
	"ｦ": {0xef, 0xbd, 0xa6},
	"ｧ": {0xef, 0xbd, 0xa7}, "ｨ": {0xef, 0xbd, 0xa8}, "ｩ": {0xef, 0xbd, 0xa9}, "ｪ": {0xef, 0xbd, 0xaa}, "ｫ": {0xef, 0xbd, 0xab},
	"ｬ": {0xef, 0xbd, 0xac}, "ｭ": {0xef, 0xbd, 0xad}, "ｮ": {0xef, 0xbd, 0xae},
	"ｯ": {0xef, 0xbd, 0xaf}, "ｰ": {0xef, 0xbd, 0xb0},
	"ｱ": {0xef, 0xbd, 0xb1}, "ｲ": {0xef, 0xbd, 0xb2}, "ｳ": {0xef, 0xbd, 0xb3}, "ｴ": {0xef, 0xbd, 0xb4}, "ｵ": {0xef, 0xbd, 0xb5},
	"ｶ": {0xef, 0xbd, 0xb6}, "ｷ": {0xef, 0xbd, 0xb7}, "ｸ": {0xef, 0xbd, 0xb8}, "ｹ": {0xef, 0xbd, 0xb9}, "ｺ": {0xef, 0xbd, 0xba},
	"ｶﾞ": {0xef, 0xbd, 0xb6, 0xef, 0xbe, 0x9e}, "ｷﾞ": {0xef, 0xbd, 0xb7, 0xef, 0xbe, 0x9e}, "ｸﾞ": {0xef, 0xbd, 0xb8, 0xef, 0xbe, 0x9e},
	"ｹﾞ": {0xef, 0xbd, 0xb9, 0xef, 0xbe, 0x9e}, "ｺﾞ": {0xef, 0xbd, 0xba, 0xef, 0xbe, 0x9e},
	"ｻ": {0xef, 0xbd, 0xbb}, "ｼ": {0xef, 0xbd, 0xbc}, "ｽ": {0xef, 0xbd, 0xbd}, "ｾ": {0xef, 0xbd, 0xbe}, "ｿ": {0xef, 0xbd, 0xbf},
	"ｻﾞ": {0xef, 0xbd, 0xbb, 0xef, 0xbe, 0x9e}, "ｼﾞ": {0xef, 0xbd, 0xbc, 0xef, 0xbe, 0x9e}, "ｽﾞ": {0xef, 0xbd, 0xbd, 0xef, 0xbe, 0x9e},
	"ｾﾞ": {0xef, 0xbd, 0xbe, 0xef, 0xbe, 0x9e}, "ｿﾞ": {0xef, 0xbd, 0xbf, 0xef, 0xbe, 0x9e},
	"ﾀ": {0xef, 0xbe, 128}, "ﾁ": {0xef, 0xbe, 0x81}, "ﾂ": {0xef, 0xbe, 0x82}, "ﾃ": {0xef, 0xbe, 0x83}, "ﾄ": {0xef, 0xbe, 0x84},
	"ﾀﾞ": {0xef, 0xbe, 128, 0xef, 0xbe, 0x9e}, "ﾁﾞ": {0xef, 0xbe, 0x81, 0xef, 0xbe, 0x9e}, "ﾂﾞ": {0xef, 0xbe, 0x82, 0xef, 0xbe, 0x9e},
	"ﾃﾞ": {0xef, 0xbe, 0x83, 0xef, 0xbe, 0x9e}, "ﾄﾞ": {0xef, 0xbe, 0x84, 0xef, 0xbe, 0x9e},
	"ﾅ": {0xef, 0xbe, 0x85}, "ﾆ": {0xef, 0xbe, 0x86}, "ﾇ": {0xef, 0xbe, 0x87}, "ﾈ": {0xef, 0xbe, 0x88}, "ﾉ": {0xef, 0xbe, 0x89},
	"ﾊ": {0xef, 0xbe, 0x8a}, "ﾋ": {0xef, 0xbe, 0x8b}, "ﾌ": {0xef, 0xbe, 0x8c}, "ﾍ": {0xef, 0xbe, 0x8d}, "ﾎ": {0xef, 0xbe, 0x8e},
	"ﾊﾞ": {0xef, 0xbe, 0x8a, 0xef, 0xbe, 0x9e}, "ﾋﾞ": {0xef, 0xbe, 0x8b, 0xef, 0xbe, 0x9e}, "ﾌﾞ": {0xef, 0xbe, 0x8c, 0xef, 0xbe, 0x9e},
	"ﾍﾞ": {0xef, 0xbe, 0x8d, 0xef, 0xbe, 0x9e}, "ﾎﾞ": {0xef, 0xbe, 0x8e, 0xef, 0xbe, 0x9e},
	"ﾊﾟ": {0xef, 0xbe, 0x8a, 0xef, 0xbe, 0x9f}, "ﾋﾟ": {0xef, 0xbe, 0x8b, 0xef, 0xbe, 0x9f}, "ﾌﾟ": {0xef, 0xbe, 0x8c, 0xef, 0xbe, 0x9f},
	"ﾍﾟ": {0xef, 0xbe, 0x8d, 0xef, 0xbe, 0x9f}, "ﾎﾟ": {0xef, 0xbe, 0x8e, 0xef, 0xbe, 0x9f},
	"ﾏ": {0xef, 0xbe, 0x8f}, "ﾐ": {0xef, 0xbe, 0x90}, "ﾑ": {0xef, 0xbe, 0x91}, "ﾒ": {0xef, 0xbe, 0x92}, "ﾓ": {0xef, 0xbe, 0x93},
	"ﾔ": {0xef, 0xbe, 0x94}, "ﾕ": {0xef, 0xbe, 0x95}, "ﾖ": {0xef, 0xbe, 0x96},
	"ﾗ": {0xef, 0xbe, 0x97}, "ﾘ": {0xef, 0xbe, 0x98}, "ﾙ": {0xef, 0xbe, 0x99}, "ﾚ": {0xef, 0xbe, 0x9a}, "ﾛ": {0xef, 0xbe, 0x9b},
	"ﾜ": {0xef, 0xbe, 0x9c}, "ﾝ": {0xef, 0xbe, 0x9d},
	"ﾞ": {0xef, 0xbe, 0x9e}, "ﾟ": {0xef, 0xbe, 0x9f},
}

func Byte(b []byte, mode string) []byte {
	orders1 := []byte{}
	orders3 := []byte{}
	filters1 := map[byte]func(w *word) []byte{}
	filters3 := map[byte]func(w *word) []byte{}
	for _, m := range []byte(mode) {
		switch m {
		case 'r':
			if _, ok := filters3[m]; !ok {
				orders3 = append(orders3, m)
				filters3[m] = convAsSmallR
			}
		case 'R':
			if _, ok := filters1[m]; !ok {
				orders1 = append(orders1, m)
				filters1[m] = convAsLargeR
			}
		case 'n':
			if _, ok := filters3[m]; !ok {
				orders3 = append(orders3, m)
				filters3[m] = convAsSmallN
			}
		case 'N':
			if _, ok := filters1[m]; !ok {
				orders1 = append(orders1, m)
				filters1[m] = convAsLargeN
			}
		case 'a':
			if _, ok := filters3[m]; !ok {
				orders3 = append(orders3, m)
				filters3[m] = convAsSmallA
			}
		case 'A':
			if _, ok := filters1[m]; !ok {
				orders1 = append(orders1, m)
				filters1[m] = convAsLargeA
			}
		case 's':
			if _, ok := filters3[m]; !ok {
				orders3 = append(orders3, m)
				filters3[m] = convAsSmallS
			}
		case 'S':
			if _, ok := filters1[m]; !ok {
				orders1 = append(orders1, m)
				filters1[m] = convAsLargeS
			}
		case 'k':
			if _, ok := filters3[m]; !ok {
				orders3 = append(orders3, m)
				filters3[m] = convAsSmallK
			}
		case 'K':
			if _, ok := filters1[m]; !ok {
				orders1 = append(orders1, m)
				filters1[m] = convAsLargeK
			}
		case 'h':
			if _, ok := filters3[m]; !ok {
				orders3 = append(orders3, m)
				filters1[m] = convAsSmallH
			}
		case 'H':
			if _, ok := filters1[m]; !ok {
				orders1 = append(orders1, m)
				filters1[m] = convAsLargeH
			}
		case 'c':
			if _, ok := filters3[m]; !ok {
				orders3 = append(orders3, m)
				filters1[m] = convAsSmallC
			}
		case 'C':
			if _, ok := filters3[m]; !ok {
				orders3 = append(orders3, m)
				filters1[m] = convAsLargeC
			}
		}
	}
	byteCount := uint64(len(b))
	buf := make([]byte, 0, byteCount)
	for i := uint64(0); i < byteCount; i++ {
		word := extract(b[i:])
		if word.len == 1 {
			for _, ordr := range orders1 {
				f := filters1[ordr]
				buf = append(buf, f(word)...)
			}
		} else if word.len == 3 || (word.len == 6 && !word.one) {
			for _, ordr := range orders3 {
				f := filters3[ordr]
				buf = append(buf, f(word)...)
			}

		} else if word.len >= 4 && word.len <= 6 {
			buf = append(buf, word.val...)
		} else {
			i++
			continue
		}
		i += uint64(word.len) - 1
	}
	return buf
}

func String(s string, mode string) string {
	return string(Byte([]byte(s), mode))
}

func NewReader(r io.Reader, mode string) *Reader {
	reader := new(Reader)
	reader.r = bufio.NewReader(r)
	reader.mode = mode
	return reader
}

func (r *Reader) Read(p []byte) (int, error) {
	line, err := r.r.ReadBytes('\n')
	buf := Byte(line, r.mode)
	if len(p) < len(buf) {
		p = append(p[0:], buf...)
		return len(p), bufio.ErrBufferFull
	}
	copy(p, buf)
	n := len(buf)
	copy(p[n:], bytes.Repeat([]byte{0}, len(p)-len(buf)+1))
	return len(p), err
}

//func NewWriter (w io.Writer, mode string) *Writer {
//	writer := new(Writer)
//	writer.w = bufio.NewWriter(w)
//	writer.mode = mode
//	return writer
//}

//func (w *Writer) Write (p []byte) (int, error) {
//	return 0, nil
//}

//func (w *Writer) WriteString (s string) (int, error){
//	return 0, nil
//}

//func (w *Writer) Flush () error {
//
//}
func is1Byte(b []byte) bool {
	if b[0] < 0x80 {
		return true
	} else {
		return false
	}
}

func is2Bytes(b []byte) bool {
	if b[0] > 0xc1 && b[0] < 0xe0 && len(b) > 1 && b[1] > 0x7f && b[1] < 0xc0 {
		return true
	} else {
		return false
	}
}

func is3Bytes(b []byte) bool {
	if b[0] < 0xf0 && len(b) > 2 && b[1] > 0x7f && b[1] < 0xc0 && b[2] > 0x7f && b[2] < 0xc0 {
		return true
	} else {
		return false
	}
}

func is4Bytes(b []byte) bool {
	if b[0] > 0xef && b[0] < 0xf5 && len(b) > 3 && b[1] > 0x7f && b[1] < 0xc0 && b[2] > 0x7f && b[2] < 0xc0 && b[3] > 0x7f && b[3] < 0xc0 {
		return true
	} else {
		return false
	}
}

func extract(b []byte) *word {
	if is1Byte(b) {
		if b[0] >= 0x30 && b[0] <= 0x39 {
			return &word{b[0:1], hankaku + numeric, 1}
		} else if b[0] >= 33 && b[0] <= 125 && b[0] != 34 && b[0] != 39 && b[0] != 92 {
			return &word{b[0:1], hankaku + alphabet, 1}
		} else if b[0] < 128 {
			return &word{b[0:1], hankaku + extra, 1}
		}
	} else if is2Bytes(b) {
		return &word{b[0:2], zenkaku + extra, 2}
	} else if is3Bytes(b) {
		if b[0] == 0xef && b[1] == 0xbd {
			if b[2] >= 0xb6 && b[2] <= 0xbf { // カ・サ行
				if b[3] == 0xef && b[4] == 0xbe && b[5] == 0x9e {
					return &word{b[0:6], hankaku + voiced, 6}
				}
			}
		} else if b[0] == 0xef && b[1] == 0xbe {
			if b[2] >= 128 && b[2] <= 0x84 { // タ行
				if b[3] == 0xef && b[4] == 0xbe && b[5] == 0x9e {
					return &word{b[0:6], hankaku + voiced, 6}
				}
			} else if b[2] >= 0x8a && b[2] <= 0x8e { // ハ行
				if b[3] == 0xef && b[4] == 0xbe && b[5] == 0x9e {
					return &word{b[0:6], hankaku + voiced, 6}
				} else if b[3] == 0xef && b[4] == 0xbe && b[5] == 0x9f {
					return &word{b[0:6], hankaku + devoiced, 6}
				}
			}
		}
		return &word{b[0:3], zenkaku + extra, 3}
	} else if is4Bytes(b) {
		return &word{b[0:4], zenkaku + extra, 4}
	}
	return &word{b[0:1], binary, 1}
}

func isVoiced(w *word) bool {
	if w.len == 6 && w.val[5] == 0x9e && w.val[4] == 0xbe && w.val[3] == 0xef {
		return true
	} else {
		return false
	}
}

func isDevoiced(w *word) bool {
	if w.len == 6 && w.val[5] == 0x9f && w.val[4] == 0xbe && w.val[3] == 0xef {
		return true
	} else {
		return false
	}
}

/**
 * Hankaku Space -> Zenkaku Space
 */
func convAsLargeS(w *word) []byte {
	if w.len == 1 && w.val[0] == 32 {
		return []byte{0xe3, 128, 128}
	}
	return w.val
}

/**
 * Zenkaku Space -> Hankaku Space
 */
func convAsSmallS(w *word) []byte {
	if w.len == 3 && w.val[0] == 0xe3 && w.val[1] == 128 && w.val[2] == 128 {
		return []byte{32}
	}
	return w.val
}

/**
 * Hankaku Numeric -> Zenkaku Numeric
 */
func convAsLargeN(w *word) []byte {
	if w.len == 1 && w.val[0] >= 48 && w.val[0] <= 57 {
		return []byte{0xef, 0xbc, 96 + w.val[0]}
	}
	return w.val
}

/**
 * Zenkaku Numeric -> Hankaku Numeric
 */
func convAsSmallN(w *word) []byte {
	if w.len == 3 && w.val[0] == 0xef && w.val[1] == 0xbc && (w.val[2] >= 0x90 && w.val[2] <= 0x99) {
		return []byte{w.val[2] - 96}
	}
	return w.val
}

/**
 * Hankaku Alphabet -> Zenkaku Alphabet
 */
func convAsLargeR(w *word) []byte {
	if w.len != 1 {
		return w.val
	}
	// A-Z -> Ａ-Ｚ
	if w.val[0] >= 65 && w.val[0] <= 90 {
		return []byte{0xef, 0xbc, 96 + w.val[0]}
	}
	// a-z -> ａ-ｚ
	if w.val[0] >= 97 && w.val[0] <= 122 {
		return []byte{0xef, 0xbd, 32 + w.val[0]}
	}
	return w.val
}

/**
 * Zenkaku Alphabet -> Hankaku Alphabet
 */
func convAsSmallR(w *word) []byte {
	if w.len != 3 {
		return w.val
	}
	// Ａ-Ｚ -> A-Z
	if w.val[0] == 0xef && w.val[1] == 0xbc && (w.val[2] >= 0xa1 && w.val[2] <= 0xba) {
		return []byte{w.val[2] - 96}
	}
	// ａ-ｚ -> a-z
	if w.val[0] == 0xef && w.val[1] == 0xbd && (w.val[2] >= 0x81 && w.val[2] <= 0x9a) {
		return []byte{w.val[2] - 32}
	}
	return w.val
}

/**
 * Hankaku AlphaNumeric -> Zenkaku AlphaNumeric
 * !-}(Excluding ",',\)
 */
func convAsLargeA(w *word) []byte {
	if w.len != 1 {
		return w.val
	}
	if w.val[0] > 32 && w.val[0] < 96 {
		if w.val[0] == 34 || w.val[0] == 39 || w.val[0] == 92 {
			return w.val
		}
		return []byte{0xef, 0xbc, 96 + w.val[0]}
	} else if w.val[0] > 95 && w.val[0] < 126 {
		return []byte{0xef, 0xbd, 32 + w.val[0]}
	}
	return w.val
}

/**
 * Zenkaku AlphaNumeric -> Hankaku AlphaNumeric
 * !-}(Excluding ",',\)
 */
func convAsSmallA(w *word) []byte {
	if w.len != 3 {
		return w.val
	}
	// ！-＿ -> !-_
	if w.val[0] == 0xef && w.val[1] == 0xbc && (w.val[2] >= 0x81 && w.val[2] <= 0xbf) {
		return []byte{w.val[2] - 96}
	}
	// ｀-｝ -> `-}
	if w.val[0] == 0xef && w.val[1] == 0xbd && (w.val[2] >= 128 && w.val[2] <= 0x9d) {
		return []byte{w.val[2] - 32}
	}
	return w.val
}

/**
 * Zenkaku Katakana -> Hankaku Katakana
 */
func convAsSmallK(w *word) []byte {
	if w.len != 3 {
		return w.val
	}
	if w.val[1] == 0x82 && w.val[0] == 0xe3 {
		switch w.val[2] {
		case 0xa1: // ァ
			return tbl["ｧ"]
		case 0xa2: // ア
			return tbl["ｱ"]
		case 0xa3: // ィ
			return tbl["ｨ"]
		case 0xa4: // イ
			return tbl["ｲ"]
		case 0xa5: // ゥ
			return tbl["ｩ"]
		case 0xa6: // ウ
			return tbl["ｳ"]
		case 0xa7: // ェ
			return tbl["ｪ"]
		case 0xa8: // エ
			return tbl["ｴ"]
		case 0xa9: // ォ
			return tbl["ｫ"]
		case 0xaa: // オ
			return tbl["ｵ"]
		case 0xab: // カ
			return tbl["ｶ"]
		case 0xac: // ガ
			return tbl["ｶﾞ"]
		case 0xad:
			return tbl["ｷ"] // キ
		case 0xae:
			return tbl["ｷﾞ"] // ギ
		case 0xaf:
			return tbl["ｸ"] // ク
		case 0xb0:
			return tbl["ｸﾞ"] // グ
		case 0xb1:
			return tbl["ｹ"] // ケ
		case 0xb2:
			return tbl["ｹﾞ"] // ゲ
		case 0xb3:
			return tbl["ｺ"] // コ
		case 0xb4:
			return tbl["ｺﾞ"] // ゴ
		case 0xb5:
			return tbl["ｻ"] // サ
		case 0xb6:
			return tbl["ｻﾞ"] // ザ
		case 0xb7:
			return tbl["ｼ"] // シ
		case 0xb8:
			return tbl["ｼﾞ"] // ジ
		case 0xb9:
			return tbl["ｽ"] // ス
		case 0xba:
			return tbl["ｽﾞ"] // ズ
		case 0xbb:
			return tbl["ｾ"] // セ
		case 0xbc:
			return tbl["ｾﾞ"] // ゼ
		case 0xbd:
			return tbl["ｿ"] // ソ
		case 0xbe:
			return tbl["ｿﾞ"] // ゾ
		case 0xbf:
			return tbl["ﾀ"] // タ
		}
	} else if w.val[1] == 0x83 && w.val[0] == 0xe3 {
		switch w.val[2] {
		case 128:
			return tbl["ﾀﾞ"] // ダ
		case 0x81:
			return tbl["ﾁ"] // チ
		case 0x82:
			return tbl["ﾁﾞ"] // ヂ
		case 0x83:
			return tbl["ｯ"] // ッ
		case 0x84:
			return tbl["ﾂ"] // ツ
		case 0x85:
			return tbl["ﾂﾞ"] // ヅ
		case 0x86:
			return tbl["ﾃ"] // テ
		case 0x87:
			return tbl["ﾃﾞ"] // デ
		case 0x88:
			return tbl["ﾄ"] // ト
		case 0x89:
			return tbl["ﾄﾞ"] // ド
		case 0x8a:
			return tbl["ﾅ"] // ナ
		case 0x8b:
			return tbl["ﾆ"] // ニ
		case 0x8c:
			return tbl["ﾇ"] // ヌ
		case 0x8d:
			return tbl["ﾈ"] // ネ
		case 0x8e:
			return tbl["ﾉ"] // ノ
		case 0x8f:
			return tbl["ﾊ"] // ハ
		case 0x90:
			return tbl["ﾊﾞ"] // バ
		case 0x91:
			return tbl["ﾊﾟ"] // パ
		case 0x92:
			return tbl["ﾋ"] // ヒ
		case 0x93:
			return tbl["ﾋﾞ"] // ビ
		case 0x94:
			return tbl["ﾋﾟ"] // ピ
		case 0x95:
			return tbl["ﾌ"] // フ
		case 0x96:
			return tbl["ﾌﾞ"] // ブ
		case 0x97:
			return tbl["ﾌﾟ"] // プ
		case 0x98:
			return tbl["ﾍ"] // ヘ
		case 0x99:
			return tbl["ﾍﾞ"] // ベ
		case 0x9a:
			return tbl["ﾍﾟ"] // ペ
		case 0x9b:
			return tbl["ﾎ"] // ホ
		case 0x9c:
			return tbl["ﾎﾞ"] // ボ
		case 0x9d:
			return tbl["ﾎﾟ"] // ポ
		case 0x9e:
			return tbl["ﾏ"] // マ
		case 0x9f:
			return tbl["ﾐ"] // ミ
		case 0xa0:
			return tbl["ﾑ"] // ム
		case 0xa1:
			return tbl["ﾒ"] // メ
		case 0xa2:
			return tbl["ﾓ"] // モ
		case 0xa3:
			return tbl["ｬ"] // ャ
		case 0xa4:
			return tbl["ﾔ"] // ヤ
		case 0xa5:
			return tbl["ｭ"] // ュ
		case 0xa6:
			return tbl["ﾕ"] // ユ
		case 0xa7:
			return tbl["ｮ"] // ョ
		case 0xa8:
			return tbl["ﾖ"] // ヨ
		case 0xa9:
			return tbl["ﾗ"] // ラ
		case 0xaa:
			return tbl["ﾘ"] // リ
		case 0xab:
			return tbl["ﾙ"] // ル
		case 0xac:
			return tbl["ﾚ"] // レ
		case 0xad:
			return tbl["ﾛ"] // ロ
		case 0xaf:
			return tbl["ﾜ"] // ワ
		case 0xb0:
			return tbl["ｲ"] // ヰ
		case 0xb1:
			return tbl["ｴ"] // ヱ
		case 0xb2:
			return tbl["ｦ"] // ヲ
		case 0xb3:
			return tbl["ﾝ"] // ン
		case 0xbc:
			return tbl["ｰ"] // ー
		}
	}
	return w.val
}

/**
 * Hankaku Katakana -> Zenkaku Katakana
 */
func convAsLargeK(w *word) []byte {
	if w.len != 3 && w.len != 6 {
		return w.val
	}
	if w.val[1] == 0xbd && w.val[0] == 0xef {
		switch w.val[2] {
		case 0xa6: // ｦ
			return tbl["ヲ"]
		case 0xa7: // ｧ
			return tbl["ァ"]
		case 0xa8: // ｨ
			return tbl["ィ"]
		case 0xa9: // ｩ
			return tbl["ゥ"]
		case 0xaa: // ｪ
			return tbl["ェ"]
		case 0xab: // ｫ
			return tbl["ォ"]
		case 0xac: // ｬ
			return tbl["ャ"]
		case 0xad: // ｭ
			return tbl["ュ"]
		case 0xae: // ｮ
			return tbl["ョ"]
		case 0xaf: // ｯ
			return tbl["ッ"]
		case 0xb0: // ｰ
			return tbl["ー"]
		case 0xb1: // ｱ
			return tbl["ア"]
		case 0xb2: // ｲ
			return tbl["イ"]
		case 0xb3: // ｳ
			return tbl["ウ"]
		case 0xb4: // ｴ
			return tbl["エ"]
		case 0xb5: // ｵ
			return tbl["オ"]
		case 0xb6: // ｶ
			if isVoiced(w) {
				return tbl["ガ"]
			} else {
				return tbl["カ"]
			}
		case 0xb7: // ｷ
			if isVoiced(w) {
				return tbl["ギ"]
			} else {
				return tbl["キ"]
			}
		case 0xb8: // ｸ
			if isVoiced(w) {
				return tbl["グ"]
			} else {
				return tbl["ク"]
			}
		case 0xb9: // ｹ
			if isVoiced(w) {
				return tbl["ゲ"]
			} else {
				return tbl["ケ"]
			}
		case 0xba: // ｺ
			if isVoiced(w) {
				return tbl["ゴ"]
			} else {
				return tbl["コ"]
			}
		case 0xbb: // ｻ
			if isVoiced(w) {
				return tbl["ザ"]
			} else {
				return tbl["サ"]
			}
		case 0xbc: // ｼ
			if isVoiced(w) {
				return tbl["ジ"]
			} else {
				return tbl["シ"]
			}
		case 0xbd: // ｽ
			if isVoiced(w) {
				return tbl["ズ"]
			} else {
				return tbl["ス"]
			}
		case 0xbe: // ｾ
			if isVoiced(w) {
				return tbl["ゼ"]
			} else {
				return tbl["セ"]
			}
		case 0xbf: // ｿ
			if isVoiced(w) {
				return tbl["ゾ"]
			} else {
				return tbl["ソ"]
			}
		}
	} else if w.val[1] == 0xbe && w.val[0] == 0xef {
		switch w.val[2] {
		case 128: // ﾀ
			if isVoiced(w) {
				return tbl["ダ"]
			} else {
				return tbl["タ"]
			}
		case 0x81: // ﾁ
			if isVoiced(w) {
				return tbl["ヂ"]
			} else {
				return tbl["チ"]
			}
		case 0x82: // ﾂ
			if isVoiced(w) {
				return tbl["ヅ"]
			} else {
				return tbl["ツ"]
			}
		case 0x83: // ﾃ
			if isVoiced(w) {
				return tbl["デ"]
			} else {
				return tbl["テ"]
			}
		case 0x84: // ﾄ
			if isVoiced(w) {
				return tbl["ド"]
			} else {
				return tbl["ト"]
			}
		case 0x85: // ナ
			return tbl["ナ"]
		case 0x86: // ニ
			return tbl["ニ"]
		case 0x87: // ヌ
			return tbl["ヌ"]
		case 0x88: // ネ
			return tbl["ネ"]
		case 0x89: // ノ
			return tbl["ノ"]
		case 0x8a: // ハ
			if isDevoiced(w) {
				return tbl["パ"]
			} else if isVoiced(w) {
				return tbl["バ"]
			} else {
				return tbl["ハ"]
			}
		case 0x8b: // ヒ
			if isDevoiced(w) {
				return tbl["ピ"]
			} else if isVoiced(w) {
				return tbl["ビ"]
			} else {
				return tbl["ヒ"]
			}
		case 0x8c: // フ
			if isDevoiced(w) {
				return tbl["プ"]
			} else if isVoiced(w) {
				return tbl["ブ"]
			} else {
				return tbl["フ"]
			}
		case 0x8d: // ヘ
			if isDevoiced(w) {
				return tbl["ペ"]
			} else if isVoiced(w) {
				return tbl["ベ"]
			} else {
				return tbl["ヘ"]
			}
		case 0x8e: // ホ
			if isDevoiced(w) {
				return tbl["ポ"]
			} else if isVoiced(w) {
				return tbl["ボ"]
			} else {
				return tbl["ホ"]
			}
		case 0x8f: // マ
			return tbl["マ"]
		case 0x90: // ミ
			return tbl["ミ"]
		case 0x91: // ム
			return tbl["ム"]
		case 0x92: // メ
			return tbl["メ"]
		case 0x93: // モ
			return tbl["モ"]
		case 0x94: // ヤ
			return tbl["ヤ"]
		case 0x95: // ユ
			return tbl["ユ"]
		case 0x96: // ヨ
			return tbl["ヨ"]
		case 0x97: // ラ
			return tbl["ラ"]
		case 0x98: // リ
			return tbl["リ"]
		case 0x99: // ル
			return tbl["ル"]
		case 0x9a: // レ
			return tbl["レ"]
		case 0x9b: // ロ
			return tbl["ロ"]
		case 0x9c: // ワ
			return tbl["ワ"]
		case 0x9d: // ン
			return tbl["ン"]
		}
	}
	return w.val
}

/**
 * Zenkaku Hiragana -> Hankaku Katakana
 */
func convAsSmallH(w *word) []byte {
	if w.len != 3 {
		return w.val
	}
	if w.val[1] == 0x81 && w.val[0] == 0xe3 {
		switch w.val[2] {
		case 0x81: // ぁ
			return tbl["ｧ"]
		case 0x82: // あ
			return tbl["ｱ"]
		case 0x83: // ぃ
			return tbl["ｨ"]
		case 0x84: // い
			return tbl["ｲ"]
		case 0x85: // ぅ
			return tbl["ｩ"]
		case 0x86: // う
			return tbl["ｳ"]
		case 0x87: // ぇ
			return tbl["ｪ"]
		case 0x88: // え
			return tbl["ｴ"]
		case 0x89: // ぉ
			return tbl["ｫ"]
		case 0x8a: // お
			return tbl["ｵ"]
		case 0x8b: // か
			return tbl["ｶ"]
		case 0x8c: // が
			return tbl["ｶﾞ"]
		case 0x8d: // き
			return tbl["ｷ"]
		case 0x8e: // ぎ
			return tbl["ｷﾞ"]
		case 0x8f: // く
			return tbl["ｸ"]
		case 0x90: // ぐ
			return tbl["ｸﾞ"]
		case 0x91: // け
			return tbl["ｹ"]
		case 0x92: // げ
			return tbl["ｹﾞ"]
		case 0x93: // こ
			return tbl["ｺ"]
		case 0x94: // ご
			return tbl["ｺﾞ"]
		case 0x95: // さ
			return tbl["ｻ"]
		case 0x96: // ざ
			return tbl["ｻﾞ"]
		case 0x97: // し
			return tbl["ｼ"]
		case 0x98: // じ
			return tbl["ｼﾞ"]
		case 0x99: // す
			return tbl["ｽ"]
		case 0x9a: // ず
			return tbl["ｽﾞ"]
		case 0x9b: // せ
			return tbl["ｾ"]
		case 0x9c: // ぜ
			return tbl["ｾﾞ"]
		case 0x9d: // そ
			return tbl["ｿ"]
		case 0x9e: // ぞ
			return tbl["ｿﾞ"]
		case 0x9f: // た
			return tbl["ﾀ"]
		case 0xa0: // だ
			return tbl["ﾀﾞ"]
		case 0xa1: // ち
			return tbl["ﾁ"]
		case 0xa2: // ぢ
			return tbl["ﾁﾞ"]
		case 0xa3: // っ
			return tbl["ｯ"]
		case 0xa4: // つ
			return tbl["ﾂ"]
		case 0xa5: // づ
			return tbl["ﾂﾞ"]
		case 0xa6: // て
			return tbl["ﾃ"]
		case 0xa7: // で
			return tbl["ﾃﾞ"]
		case 0xa8: // と
			return tbl["ﾄ"]
		case 0xa9: // ど
			return tbl["ﾄﾞ"]
		case 0xaa: // な
			return tbl["ﾅ"]
		case 0xab: // に
			return tbl["ﾆ"]
		case 0xac: // ぬ
			return tbl["ﾇ"]
		case 0xad: // ね
			return tbl["ﾈ"]
		case 0xae: // の
			return tbl["ﾉ"]
		case 0xaf: // は
			return tbl["ﾊ"]
		case 0xb0: // ば
			return tbl["ﾊﾞ"]
		case 0xb1: // ぱ
			return tbl["ﾊﾟ"]
		case 0xb2: // ひ
			return tbl["ﾋ"]
		case 0xb3: // び
			return tbl["ﾋﾞ"]
		case 0xb4: // ぴ
			return tbl["ﾋﾟ"]
		case 0xb5: // ふ
			return tbl["ﾌ"]
		case 0xb6: // ぶ
			return tbl["ﾌﾞ"]
		case 0xb7: // ぷ
			return tbl["ﾌﾟ"]
		case 0xb8: // へ
			return tbl["ﾍ"]
		case 0xb9: // べ
			return tbl["ﾍﾞ"]
		case 0xba: // ぺ
			return tbl["ﾍﾟ"]
		case 0xbb: // ほ
			return tbl["ﾎ"]
		case 0xbc: // ぼ
			return tbl["ﾎﾞ"]
		case 0xbd: // ぽ
			return tbl["ﾎﾟ"]
		case 0xbe: // ま
			return tbl["ﾏ"]
		case 0xbf: // み
			return tbl["ﾐ"]
		}
	} else if w.val[1] == 0x82 && w.val[0] == 0xe3 {
		switch w.val[2] {
		case 128: // む
			return tbl["ﾑ"]
		case 0x81: // め
			return tbl["ﾒ"]
		case 0x82: // も
			return tbl["ﾓ"]
		case 0x83: // ゃ
			return tbl["ｬ"]
		case 0x84: // や
			return tbl["ﾔ"]
		case 0x85: // ゅ
			return tbl["ｭ"]
		case 0x86: // ゆ
			return tbl["ﾕ"]
		case 0x87: // ょ
			return tbl["ｮ"]
		case 0x88: // よ
			return tbl["ﾖ"]
		case 0x89: // ら
			return tbl["ﾗ"]
		case 0x8a: // り
			return tbl["ﾘ"]
		case 0x8b: // る
			return tbl["ﾙ"]
		case 0x8c: // れ
			return tbl["ﾚ"]
		case 0x8d: // ろ
			return tbl["ﾛ"]
		case 0x8f: // わ
			return tbl["ﾜ"]
		case 0x90: // ゐ
			return tbl["ｲ"]
		case 0x91: // ゑ
			return tbl["ｴ"]
		case 0x92: // を
			return tbl["ｦ"]
		case 0x93: // ん
			return tbl["ﾝ"]
		}
	} else if w.val[2] == 0xbc && w.val[1] == 0x83 && w.val[0] == 0xe3 {
		return tbl["ｰ"]
	}
	return w.val
}

/**
 * Hankaku Katakana -> Zenkaku Hiragana
 */
func convAsLargeH(w *word) []byte {
	if w.len != 3 && w.len != 6 {
		return w.val
	}
	if w.val[1] == 0xbd && w.val[0] == 0xef {
		switch w.val[2] {
		case 0xa6: // ｦ
			return tbl["を"]
		case 0xa7: // ｧ
			return tbl["ぁ"]
		case 0xa8: // ｨ
			return tbl["ぃ"]
		case 0xa9: // ｩ
			return tbl["ぅ"]
		case 0xaa: // ｪ
			return tbl["ぇ"]
		case 0xab: // ｫ
			return tbl["ぉ"]
		case 0xac: // ｬ
			return tbl["ゃ"]
		case 0xad: // ｭ
			return tbl["ゅ"]
		case 0xae: // ｮ
			return tbl["ょ"]
		case 0xaf: // ｯ
			return tbl["っ"]
		case 0xb0: // ｰ
			return tbl["ー"]
		case 0xb1: // ｱ
			return tbl["あ"]
		case 0xb2: // ｲ
			return tbl["い"]
		case 0xb3: // ｳ
			return tbl["う"]
		case 0xb4: // ｴ
			return tbl["え"]
		case 0xb5: // ｵ
			return tbl["お"]
		case 0xb6: // ｶ
			if isVoiced(w) {
				return tbl["が"]
			} else {
				return tbl["か"]
			}
		case 0xb7: // ｷ
			if isVoiced(w) {
				return tbl["ぎ"]
			} else {
				return tbl["き"]
			}
		case 0xb8: // ｸ
			if isVoiced(w) {
				return tbl["ぐ"]
			} else {
				return tbl["く"]
			}
		case 0xb9: // ｹ
			if isVoiced(w) {
				return tbl["げ"]
			} else {
				return tbl["け"]
			}
		case 0xba: // ｺ
			if isVoiced(w) {
				return tbl["ご"]
			} else {
				return tbl["こ"]
			}
		case 0xbb: // ｻ
			if isVoiced(w) {
				return tbl["ざ"]
			} else {
				return tbl["さ"]
			}
		case 0xbc: // ｼ
			if isVoiced(w) {
				return tbl["じ"]
			} else {
				return tbl["し"]
			}
		case 0xbd: // ｽ
			if isVoiced(w) {
				return tbl["ず"]
			} else {
				return tbl["す"]
			}
		case 0xbe: // ｾ
			if isVoiced(w) {
				return tbl["ぜ"]
			} else {
				return tbl["せ"]
			}
		case 0xbf: // ｿ
			if isVoiced(w) {
				return tbl["ぞ"]
			} else {
				return tbl["そ"]
			}
		}
	} else if w.val[1] == 0xbe && w.val[0] == 0xef {
		switch w.val[2] {
		case 128: // ﾀ
			if isVoiced(w) {
				return tbl["だ"]
			} else {
				return tbl["た"]
			}
		case 0x81: // ﾁ
			if isVoiced(w) {
				return tbl["ぢ"]
			} else {
				return tbl["ち"]
			}
		case 0x82: // ﾂ
			if isVoiced(w) {
				return tbl["づ"]
			} else {
				return tbl["つ"]
			}
		case 0x83: // ﾃ
			if isVoiced(w) {
				return tbl["で"]
			} else {
				return tbl["て"]
			}
		case 0x84: // ﾄ
			if isVoiced(w) {
				return tbl["ど"]
			} else {
				return tbl["と"]
			}
		case 0x85: // ﾅ
			return tbl["な"]
		case 0x86: // ﾆ
			return tbl["に"]
		case 0x87: // ﾇ
			return tbl["ぬ"]
		case 0x88: // ﾈ
			return tbl["ね"]
		case 0x89: // ﾉ
			return tbl["の"]
		case 0x8a: // ﾊ
			if isDevoiced(w) {
				return tbl["ぱ"]
			} else if isVoiced(w) {
				return tbl["ば"]
			} else {
				return tbl["は"]
			}
		case 0x8b: // ﾋ
			if isDevoiced(w) {
				return tbl["ぴ"]
			} else if isVoiced(w) {
				return tbl["び"]
			} else {
				return tbl["ひ"]
			}
		case 0x8c: // ﾌ
			if isDevoiced(w) {
				return tbl["ぷ"]
			} else if isVoiced(w) {
				return tbl["ぶ"]
			} else {
				return tbl["ふ"]
			}
		case 0x8d: // ﾍ
			if isDevoiced(w) {
				return tbl["ぺ"]
			} else if isVoiced(w) {
				return tbl["べ"]
			} else {
				return tbl["へ"]
			}
		case 0x8e: // ﾎ
			if isDevoiced(w) {
				return tbl["ぽ"]
			} else if isVoiced(w) {
				return tbl["ぼ"]
			} else {
				return tbl["ほ"]
			}
		case 0x8f: // ﾏ
			return tbl["ま"]
		case 0x90: // ﾐ
			return tbl["み"]
		case 0x91: // ﾑ
			return tbl["む"]
		case 0x92: // ﾒ
			return tbl["め"]
		case 0x93: // ﾓ
			return tbl["も"]
		case 0x94: // ﾔ
			return tbl["や"]
		case 0x95: // ﾕ
			return tbl["ゆ"]
		case 0x96: // ﾖ
			return tbl["よ"]
		case 0x97: // ﾗ
			return tbl["ら"]
		case 0x98: // ﾘ
			return tbl["り"]
		case 0x99: // ﾙ
			return tbl["る"]
		case 0x9a: // ﾚ
			return tbl["れ"]
		case 0x9b: // ﾛ
			return tbl["ろ"]
		case 0x9c: // ﾜ
			return tbl["わ"]
		case 0x9d: // ﾝ
			return tbl["ん"]
		}
	}
	return w.val
}

/**
 * Zenkaku Katakana -> Zenkaku Hiragana
 */
func convAsSmallC(w *word) []byte {
	if w.len != 3 {
		return w.val
	}
	if w.val[0] == 0xe3 {
		if w.val[1] == 0x82 { // ァ-タ
			if w.val[2] >= 0xa1 && w.val[2] <= 0xbf {
				return []byte{0xe3, 0x81, w.val[2] - 32}
			}
		} else if w.val[1] == 0x83 { // ダ-ン
			if w.val[2] >= 128 && w.val[2] <= 0x9f { // ダ-ミ
				return []byte{0xe3, 0x81, w.val[2] + 32}
			} else if w.val[2] >= 0xa0 && w.val[2] <= 0xb3 { // ム-ン
				return []byte{0xe3, 0x82, w.val[2] - 32}
			}
		}
	}
	return w.val
}

/**
 * Zenkaku Hiragana -> Zenkaku Katakana
 */
func convAsLargeC(w *word) []byte {
	if w.len != 3 {
		return w.val
	}
	if w.val[0] == 0xe3 {
		if w.val[1] == 0x81 { // ぁ-み
			if w.val[2] >= 0x81 && w.val[2] <= 0x9f { // ぁ-た
				return []byte{0xe3, 0x82, w.val[2] + 32}
			} else if w.val[2] >= 0xa0 && w.val[2] <= 0xbf { // だ-み
				return []byte{0xe3, 0x83, w.val[2] - 32}
			}
		} else if w.val[1] == 0x82 { // む-ん
			if w.val[2] >= 128 && w.val[2] <= 0x93 {
				return []byte{0xe3, 0x83, w.val[2] + 32}
			}
		}
	}
	return w.val
}

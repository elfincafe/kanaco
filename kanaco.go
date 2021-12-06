package kanaco

import (
	"bufio"
	"bytes"
	"io"
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
	val []byte
	one bool // 1byte character or not. (ex) ｶﾞ is false because of two characters.
	len int
}

var tbl map[string][]byte = map[string][]byte{
	"ぁ": {227, 129, 129}, "あ": {227, 129, 130},
	"ぃ": {227, 129, 131}, "い": {227, 129, 132},
	"ぅ": {227, 129, 133}, "う": {227, 129, 134},
	"ぇ": {227, 129, 135}, "え": {227, 129, 136},
	"ぉ": {227, 129, 137}, "お": {227, 129, 138},
	"か": {227, 129, 139}, "が": {227, 129, 140},
	"き": {227, 129, 141}, "ぎ": {227, 129, 142},
	"く": {227, 129, 143}, "ぐ": {227, 129, 144},
	"け": {227, 129, 145}, "げ": {227, 129, 146},
	"こ": {227, 129, 147}, "ご": {227, 129, 148},
	"さ": {227, 129, 149}, "ざ": {227, 129, 150},
	"し": {227, 129, 151}, "じ": {227, 129, 152},
	"す": {227, 129, 153}, "ず": {227, 129, 154},
	"せ": {227, 129, 155}, "ぜ": {227, 129, 156},
	"そ": {227, 129, 157}, "ぞ": {227, 129, 158},
	"た": {227, 129, 159}, "だ": {227, 129, 160},
	"ち": {227, 129, 161}, "ぢ": {227, 129, 162},
	"っ": {227, 129, 163}, "つ": {227, 129, 164}, "づ": {227, 129, 165},
	"て": {227, 129, 166}, "で": {227, 129, 167},
	"と": {227, 129, 168}, "ど": {227, 129, 169},
	"な": {227, 129, 170},
	"に": {227, 129, 171},
	"ぬ": {227, 129, 172},
	"ね": {227, 129, 173},
	"の": {227, 129, 174},
	"は": {227, 129, 175}, "ば": {227, 129, 176}, "ぱ": {227, 129, 177},
	"ひ": {227, 129, 178}, "び": {227, 129, 179}, "ぴ": {227, 129, 180},
	"ふ": {227, 129, 181}, "ぶ": {227, 129, 182}, "ぷ": {227, 129, 183},
	"へ": {227, 129, 184}, "べ": {227, 129, 185}, "ぺ": {227, 129, 186},
	"ほ": {227, 129, 187}, "ぼ": {227, 129, 188}, "ぽ": {227, 129, 189},
	"ま": {227, 129, 190},
	"み": {227, 129, 191},
	"む": {227, 130, 128},
	"め": {227, 130, 129},
	"も": {227, 130, 130},
	"ゃ": {227, 130, 131}, "や": {227, 130, 132},
	"ゅ": {227, 130, 133}, "ゆ": {227, 130, 134},
	"ょ": {227, 130, 135}, "よ": {227, 130, 136},
	"ら": {227, 130, 137},
	"り": {227, 130, 138},
	"る": {227, 130, 139},
	"れ": {227, 130, 140},
	"ろ": {227, 130, 141},
	"ゎ": {227, 130, 142},
	"わ": {227, 130, 143},
	"ゐ": {227, 130, 144},
	"ゑ": {227, 130, 145},
	"を": {227, 130, 146},
	"ん": {227, 130, 147},
	"ァ": {227, 130, 161}, "ア": {227, 130, 162},
	"ィ": {227, 130, 163}, "イ": {227, 130, 164},
	"ゥ": {227, 130, 165}, "ウ": {227, 130, 166},
	"ェ": {227, 130, 167}, "エ": {227, 130, 168},
	"ォ": {227, 130, 169}, "オ": {227, 130, 170},
	"カ": {227, 130, 171}, "ガ": {227, 130, 172},
	"キ": {227, 130, 173}, "ギ": {227, 130, 174},
	"ク": {227, 130, 175}, "グ": {227, 130, 176},
	"ケ": {227, 130, 177}, "ゲ": {227, 130, 178},
	"コ": {227, 130, 179}, "ゴ": {227, 130, 180},
	"サ": {227, 130, 181}, "ザ": {227, 130, 182},
	"シ": {227, 130, 183}, "ジ": {227, 130, 184},
	"ス": {227, 130, 185}, "ズ": {227, 130, 186},
	"セ": {227, 130, 187}, "ゼ": {227, 130, 188},
	"ソ": {227, 130, 189}, "ゾ": {227, 130, 190},
	"タ": {227, 130, 191}, "ダ": {227, 131, 128},
	"チ": {227, 131, 129}, "ヂ": {227, 131, 130},
	"ッ": {227, 131, 131}, "ツ": {227, 131, 132}, "ヅ": {227, 131, 133},
	"テ": {227, 131, 134}, "デ": {227, 131, 135},
	"ト": {227, 131, 136}, "ド": {227, 131, 137},
	"ナ": {227, 131, 138},
	"ニ": {227, 131, 139},
	"ヌ": {227, 131, 140},
	"ネ": {227, 131, 141},
	"ノ": {227, 131, 142},
	"ハ": {227, 131, 143}, "バ": {227, 131, 144}, "パ": {227, 131, 145},
	"ヒ": {227, 131, 146}, "ビ": {227, 131, 147}, "ピ": {227, 131, 148},
	"フ": {227, 131, 149}, "ブ": {227, 131, 150}, "プ": {227, 131, 151},
	"ヘ": {227, 131, 152}, "ベ": {227, 131, 153}, "ペ": {227, 131, 154},
	"ホ": {227, 131, 155}, "ボ": {227, 131, 156}, "ポ": {227, 131, 157},
	"マ": {227, 131, 158},
	"ミ": {227, 131, 159},
	"ム": {227, 131, 160},
	"メ": {227, 131, 161},
	"モ": {227, 131, 162},
	"ャ": {227, 131, 163}, "ヤ": {227, 131, 164},
	"ュ": {227, 131, 165}, "ユ": {227, 131, 166},
	"ョ": {227, 131, 167}, "ヨ": {227, 131, 168},
	"ラ": {227, 131, 169},
	"リ": {227, 131, 170},
	"ル": {227, 131, 171},
	"レ": {227, 131, 172},
	"ロ": {227, 131, 173},
	"ヮ": {227, 131, 174}, "ワ": {227, 131, 175},
	"ヰ": {227, 131, 176},
	"ヱ": {227, 131, 177},
	"ヲ": {227, 131, 178},
	"ン": {227, 131, 179},
	"ー": {227, 131, 188},
	"ｦ": {239, 189, 166},
	"ｧ": {239, 189, 167}, "ｨ": {239, 189, 168}, "ｩ": {239, 189, 169}, "ｪ": {239, 189, 170}, "ｫ": {239, 189, 171},
	"ｬ": {239, 189, 172}, "ｭ": {239, 189, 173}, "ｮ": {239, 189, 174},
	"ｯ": {239, 189, 175}, "ｰ": {239, 189, 176},
	"ｱ": {239, 189, 177}, "ｲ": {239, 189, 178}, "ｳ": {239, 189, 179}, "ｴ": {239, 189, 180}, "ｵ": {239, 189, 181},
	"ｶ": {239, 189, 182}, "ｷ": {239, 189, 183}, "ｸ": {239, 189, 184}, "ｹ": {239, 189, 185}, "ｺ": {239, 189, 186},
	"ｶﾞ": {239, 189, 182, 239, 190, 158}, "ｷﾞ": {239, 189, 183, 239, 190, 158}, "ｸﾞ": {239, 189, 184, 239, 190, 158},
	"ｹﾞ": {239, 189, 185, 239, 190, 158}, "ｺﾞ": {239, 189, 186, 239, 190, 158},
	"ｻ": {239, 189, 187}, "ｼ": {239, 189, 188}, "ｽ": {239, 189, 189}, "ｾ": {239, 189, 190}, "ｿ": {239, 189, 191},
	"ｻﾞ": {239, 189, 187, 239, 190, 158}, "ｼﾞ": {239, 189, 188, 239, 190, 158}, "ｽﾞ": {239, 189, 189, 239, 190, 158},
	"ｾﾞ": {239, 189, 190, 239, 190, 158}, "ｿﾞ": {239, 189, 191, 239, 190, 158},
	"ﾀ": {239, 190, 128}, "ﾁ": {239, 190, 129}, "ﾂ": {239, 190, 130}, "ﾃ": {239, 190, 131}, "ﾄ": {239, 190, 132},
	"ﾀﾞ": {239, 190, 128, 239, 190, 158}, "ﾁﾞ": {239, 190, 129, 239, 190, 158}, "ﾂﾞ": {239, 190, 130, 239, 190, 158},
	"ﾃﾞ": {239, 190, 131, 239, 190, 158}, "ﾄﾞ": {239, 190, 132, 239, 190, 158},
	"ﾅ": {239, 190, 133}, "ﾆ": {239, 190, 134}, "ﾇ": {239, 190, 135}, "ﾈ": {239, 190, 136}, "ﾉ": {239, 190, 137},
	"ﾊ": {239, 190, 138}, "ﾋ": {239, 190, 139}, "ﾌ": {239, 190, 140}, "ﾍ": {239, 190, 141}, "ﾎ": {239, 190, 142},
	"ﾊﾞ": {239, 190, 138, 239, 190, 158}, "ﾋﾞ": {239, 190, 139, 239, 190, 158}, "ﾌﾞ": {239, 190, 140, 239, 190, 158},
	"ﾍﾞ": {239, 190, 141, 239, 190, 158}, "ﾎﾞ": {239, 190, 142, 239, 190, 158},
	"ﾊﾟ": {239, 190, 138, 239, 190, 159}, "ﾋﾟ": {239, 190, 139, 239, 190, 159}, "ﾌﾟ": {239, 190, 140, 239, 190, 159},
	"ﾍﾟ": {239, 190, 141, 239, 190, 159}, "ﾎﾟ": {239, 190, 142, 239, 190, 159},
	"ﾏ": {239, 190, 143}, "ﾐ": {239, 190, 144}, "ﾑ": {239, 190, 145}, "ﾒ": {239, 190, 146}, "ﾓ": {239, 190, 147},
	"ﾔ": {239, 190, 148}, "ﾕ": {239, 190, 149}, "ﾖ": {239, 190, 150},
	"ﾗ": {239, 190, 151}, "ﾘ": {239, 190, 152}, "ﾙ": {239, 190, 153}, "ﾚ": {239, 190, 154}, "ﾛ": {239, 190, 155},
	"ﾜ": {239, 190, 156}, "ﾝ": {239, 190, 157},
	"ﾞ": {239, 190, 158}, "ﾟ": {239, 190, 159},
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
			for _, f := range filters1 {
				buf = append(buf, f(word)...)
			}
		} else if word.len == 3 || (word.len == 6 && !word.one) {
			for _, f := range filters3 {
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

func extract(b []byte) *word {
	if b[0] < 128 { // 1byte
		return &word{b[0:1], true, 1}
	} else if b[0] < 192 {
		return &word{b[0:0], true, 0}
	} else if b[0] < 224 { // 2byte
		return &word{b[0:2], true, 2}
	} else if b[0] < 240 { // 3byte
		if len(b) >= 6 {
			if b[0] == 239 && b[1] == 189 {
				if b[2] >= 182 && b[2] <= 191 { // カ・サ行
					if b[3] == 239 && b[4] == 190 && b[5] == 158 {
						return &word{b[0:6], false, 6}
					}
				}
			} else if b[0] == 239 && b[1] == 190 {
				if b[2] >= 128 && b[2] <= 132 { // タ行
					if b[3] == 239 && b[4] == 190 && b[5] == 158 {
						return &word{b[0:6], false, 6}
					}
				} else if b[2] >= 138 && b[2] <= 142 { // ハ行
					if b[3] == 239 && b[4] == 190 && (b[5] == 158 || b[5] == 159) {
						return &word{b[0:6], false, 6}
					}
				}
			}
		}
		return &word{b[0:3], true, 3}
	} else if b[0] < 248 { // 4byte
		return &word{b[0:4], true, 4}
	}
	return &word{b[0:0], true, 0}
}

func isVoiced(w *word) bool {
	if w.len == 6 && w.val[5] == 158 && w.val[4] == 190 && w.val[3] == 239 {
		return true
	} else {
		return false
	}
}

func isDevoiced(w *word) bool {
	if w.len == 6 && w.val[5] == 159 && w.val[4] == 190 && w.val[3] == 239 {
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
		return []byte{227, 128, 128}
	}
	return w.val
}

/**
 * Zenkaku Space -> Hankaku Space
 */
func convAsSmallS(w *word) []byte {
	if w.len == 3 && w.val[0] == 227 && w.val[1] == 128 && w.val[2] == 128 {
		return []byte{32}
	}
	return w.val
}

/**
 * Hankaku Numeric -> Zenkaku Numeric
 */
func convAsLargeN(w *word) []byte {
	if w.len == 1 && w.val[0] >= 48 && w.val[0] <= 57 {
		return []byte{239, 188, 96 + w.val[0]}
	}
	return w.val
}

/**
 * Zenkaku Numeric -> Hankaku Numeric
 */
func convAsSmallN(w *word) []byte {
	if w.len == 3 && w.val[0] == 239 && w.val[1] == 188 && (w.val[2] >= 144 && w.val[2] <= 153) {
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
		return []byte{239, 188, 96 + w.val[0]}
	}
	// a-z -> ａ-ｚ
	if w.val[0] >= 97 && w.val[0] <= 122 {
		return []byte{239, 189, 32 + w.val[0]}
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
	if w.val[0] == 239 && w.val[1] == 188 && (w.val[2] >= 161 && w.val[2] <= 186) {
		return []byte{w.val[2] - 96}
	}
	// ａ-ｚ -> a-z
	if w.val[0] == 239 && w.val[1] == 189 && (w.val[2] >= 129 && w.val[2] <= 154) {
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
		return []byte{239, 188, 96 + w.val[0]}
	} else if w.val[0] > 95 && w.val[0] < 126 {
		return []byte{239, 189, 32 + w.val[0]}
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
	if w.val[0] == 239 && w.val[1] == 188 && (w.val[2] >= 129 && w.val[2] <= 191) {
		return []byte{w.val[2] - 96}
	}
	// ｀-｝ -> `-}
	if w.val[0] == 239 && w.val[1] == 189 && (w.val[2] >= 128 && w.val[2] <= 157) {
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
	if w.val[1] == 130 && w.val[0] == 227 {
		switch w.val[2] {
		case 161: // ァ
			return tbl["ｧ"]
		case 162: // ア
			return tbl["ｱ"]
		case 163: // ィ
			return tbl["ｨ"]
		case 164: // イ
			return tbl["ｲ"]
		case 165: // ゥ
			return tbl["ｩ"]
		case 166: // ウ
			return tbl["ｳ"]
		case 167: // ェ
			return tbl["ｪ"]
		case 168: // エ
			return tbl["ｴ"]
		case 169: // ォ
			return tbl["ｫ"]
		case 170: // オ
			return tbl["ｵ"]
		case 171: // カ
			return tbl["ｶ"]
		case 172: // ガ
			return tbl["ｶﾞ"]
		case 173:
			return tbl["ｷ"] // キ
		case 174:
			return tbl["ｷﾞ"] // ギ
		case 175:
			return tbl["ｸ"] // ク
		case 176:
			return tbl["ｸﾞ"] // グ
		case 177:
			return tbl["ｹ"] // ケ
		case 178:
			return tbl["ｹﾞ"] // ゲ
		case 179:
			return tbl["ｺ"] // コ
		case 180:
			return tbl["ｺﾞ"] // ゴ
		case 181:
			return tbl["ｻ"] // サ
		case 182:
			return tbl["ｻﾞ"] // ザ
		case 183:
			return tbl["ｼ"] // シ
		case 184:
			return tbl["ｼﾞ"] // ジ
		case 185:
			return tbl["ｽ"] // ス
		case 186:
			return tbl["ｽﾞ"] // ズ
		case 187:
			return tbl["ｾ"] // セ
		case 188:
			return tbl["ｾﾞ"] // ゼ
		case 189:
			return tbl["ｿ"] // ソ
		case 190:
			return tbl["ｿﾞ"] // ゾ
		case 191:
			return tbl["ﾀ"] // タ
		}
	} else if w.val[1] == 131 && w.val[0] == 227 {
		switch w.val[2] {
		case 128:
			return tbl["ﾀﾞ"] // ダ
		case 129:
			return tbl["ﾁ"] // チ
		case 130:
			return tbl["ﾁﾞ"] // ヂ
		case 131:
			return tbl["ｯ"] // ッ
		case 132:
			return tbl["ﾂ"] // ツ
		case 133:
			return tbl["ﾂﾞ"] // ヅ
		case 134:
			return tbl["ﾃ"] // テ
		case 135:
			return tbl["ﾃﾞ"] // デ
		case 136:
			return tbl["ﾄ"] // ト
		case 137:
			return tbl["ﾄﾞ"] // ド
		case 138:
			return tbl["ﾅ"] // ナ
		case 139:
			return tbl["ﾆ"] // ニ
		case 140:
			return tbl["ﾇ"] // ヌ
		case 141:
			return tbl["ﾈ"] // ネ
		case 142:
			return tbl["ﾉ"] // ノ
		case 143:
			return tbl["ﾊ"] // ハ
		case 144:
			return tbl["ﾊﾞ"] // バ
		case 145:
			return tbl["ﾊﾟ"] // パ
		case 146:
			return tbl["ﾋ"] // ヒ
		case 147:
			return tbl["ﾋﾞ"] // ビ
		case 148:
			return tbl["ﾋﾟ"] // ピ
		case 149:
			return tbl["ﾌ"] // フ
		case 150:
			return tbl["ﾌﾞ"] // ブ
		case 151:
			return tbl["ﾌﾟ"] // プ
		case 152:
			return tbl["ﾍ"] // ヘ
		case 153:
			return tbl["ﾍﾞ"] // ベ
		case 154:
			return tbl["ﾍﾟ"] // ペ
		case 155:
			return tbl["ﾎ"] // ホ
		case 156:
			return tbl["ﾎﾞ"] // ボ
		case 157:
			return tbl["ﾎﾟ"] // ポ
		case 158:
			return tbl["ﾏ"] // マ
		case 159:
			return tbl["ﾐ"] // ミ
		case 160:
			return tbl["ﾑ"] // ム
		case 161:
			return tbl["ﾒ"] // メ
		case 162:
			return tbl["ﾓ"] // モ
		case 163:
			return tbl["ｬ"] // ャ
		case 164:
			return tbl["ﾔ"] // ヤ
		case 165:
			return tbl["ｭ"] // ュ
		case 166:
			return tbl["ﾕ"] // ユ
		case 167:
			return tbl["ｮ"] // ョ
		case 168:
			return tbl["ﾖ"] // ヨ
		case 169:
			return tbl["ﾗ"] // ラ
		case 170:
			return tbl["ﾘ"] // リ
		case 171:
			return tbl["ﾙ"] // ル
		case 172:
			return tbl["ﾚ"] // レ
		case 173:
			return tbl["ﾛ"] // ロ
		case 175:
			return tbl["ﾜ"] // ワ
		case 176:
			return tbl["ｲ"] // ヰ
		case 177:
			return tbl["ｴ"] // ヱ
		case 178:
			return tbl["ｦ"] // ヲ
		case 179:
			return tbl["ﾝ"] // ン
		case 188:
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
	if w.val[1] == 189 && w.val[0] == 239 {
		switch w.val[2] {
		case 166: // ｦ
			return tbl["ヲ"]
		case 167: // ｧ
			return tbl["ァ"]
		case 168: // ｨ
			return tbl["ィ"]
		case 169: // ｩ
			return tbl["ゥ"]
		case 170: // ｪ
			return tbl["ェ"]
		case 171: // ｫ
			return tbl["ォ"]
		case 172: // ｬ
			return tbl["ャ"]
		case 173: // ｭ
			return tbl["ュ"]
		case 174: // ｮ
			return tbl["ョ"]
		case 175: // ｯ
			return tbl["ッ"]
		case 176: // ｰ
			return tbl["ー"]
		case 177: // ｱ
			return tbl["ア"]
		case 178: // ｲ
			return tbl["イ"]
		case 179: // ｳ
			return tbl["ウ"]
		case 180: // ｴ
			return tbl["エ"]
		case 181: // ｵ
			return tbl["オ"]
		case 182: // ｶ
			if isVoiced(w) {
				return tbl["ガ"]
			} else {
				return tbl["カ"]
			}
		case 183: // ｷ
			if isVoiced(w) {
				return tbl["ギ"]
			} else {
				return tbl["キ"]
			}
		case 184: // ｸ
			if isVoiced(w) {
				return tbl["グ"]
			} else {
				return tbl["ク"]
			}
		case 185: // ｹ
			if isVoiced(w) {
				return tbl["ゲ"]
			} else {
				return tbl["ケ"]
			}
		case 186: // ｺ
			if isVoiced(w) {
				return tbl["ゴ"]
			} else {
				return tbl["コ"]
			}
		case 187: // ｻ
			if isVoiced(w) {
				return tbl["ザ"]
			} else {
				return tbl["サ"]
			}
		case 188: // ｼ
			if isVoiced(w) {
				return tbl["ジ"]
			} else {
				return tbl["シ"]
			}
		case 189: // ｽ
			if isVoiced(w) {
				return tbl["ズ"]
			} else {
				return tbl["ス"]
			}
		case 190: // ｾ
			if isVoiced(w) {
				return tbl["ゼ"]
			} else {
				return tbl["セ"]
			}
		case 191: // ｿ
			if isVoiced(w) {
				return tbl["ゾ"]
			} else {
				return tbl["ソ"]
			}
		}
	} else if w.val[1] == 190 && w.val[0] == 239 {
		switch w.val[2] {
		case 128: // ﾀ
			if isVoiced(w) {
				return tbl["ダ"]
			} else {
				return tbl["タ"]
			}
		case 129: // ﾁ
			if isVoiced(w) {
				return tbl["ヂ"]
			} else {
				return tbl["チ"]
			}
		case 130: // ﾂ
			if isVoiced(w) {
				return tbl["ヅ"]
			} else {
				return tbl["ツ"]
			}
		case 131: // ﾃ
			if isVoiced(w) {
				return tbl["デ"]
			} else {
				return tbl["テ"]
			}
		case 132: // ﾄ
			if isVoiced(w) {
				return tbl["ド"]
			} else {
				return tbl["ト"]
			}
		case 133: // ナ
			return tbl["ナ"]
		case 134: // ニ
			return tbl["ニ"]
		case 135: // ヌ
			return tbl["ヌ"]
		case 136: // ネ
			return tbl["ネ"]
		case 137: // ノ
			return tbl["ノ"]
		case 138: // ハ
			if isDevoiced(w) {
				return tbl["パ"]
			} else if isVoiced(w) {
				return tbl["バ"]
			} else {
				return tbl["ハ"]
			}
		case 139: // ヒ
			if isDevoiced(w) {
				return tbl["ピ"]
			} else if isVoiced(w) {
				return tbl["ビ"]
			} else {
				return tbl["ヒ"]
			}
		case 140: // フ
			if isDevoiced(w) {
				return tbl["プ"]
			} else if isVoiced(w) {
				return tbl["ブ"]
			} else {
				return tbl["フ"]
			}
		case 141: // ヘ
			if isDevoiced(w) {
				return tbl["ペ"]
			} else if isVoiced(w) {
				return tbl["ベ"]
			} else {
				return tbl["ヘ"]
			}
		case 142: // ホ
			if isDevoiced(w) {
				return tbl["ポ"]
			} else if isVoiced(w) {
				return tbl["ボ"]
			} else {
				return tbl["ホ"]
			}
		case 143: // マ
			return tbl["マ"]
		case 144: // ミ
			return tbl["ミ"]
		case 145: // ム
			return tbl["ム"]
		case 146: // メ
			return tbl["メ"]
		case 147: // モ
			return tbl["モ"]
		case 148: // ヤ
			return tbl["ヤ"]
		case 149: // ユ
			return tbl["ユ"]
		case 150: // ヨ
			return tbl["ヨ"]
		case 151: // ラ
			return tbl["ラ"]
		case 152: // リ
			return tbl["リ"]
		case 153: // ル
			return tbl["ル"]
		case 154: // レ
			return tbl["レ"]
		case 155: // ロ
			return tbl["ロ"]
		case 156: // ワ
			return tbl["ワ"]
		case 157: // ン
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
	if w.val[1] == 129 && w.val[0] == 227 {
		switch w.val[2] {
		case 129: // ぁ
			return tbl["ｧ"]
		case 130: // あ
			return tbl["ｱ"]
		case 131: // ぃ
			return tbl["ｨ"]
		case 132: // い
			return tbl["ｲ"]
		case 133: // ぅ
			return tbl["ｩ"]
		case 134: // う
			return tbl["ｳ"]
		case 135: // ぇ
			return tbl["ｪ"]
		case 136: // え
			return tbl["ｴ"]
		case 137: // ぉ
			return tbl["ｫ"]
		case 138: // お
			return tbl["ｵ"]
		case 139: // か
			return tbl["ｶ"]
		case 140: // が
			return tbl["ｶﾞ"]
		case 141: // き
			return tbl["ｷ"]
		case 142: // ぎ
			return tbl["ｷﾞ"]
		case 143: // く
			return tbl["ｸ"]
		case 144: // ぐ
			return tbl["ｸﾞ"]
		case 145: // け
			return tbl["ｹ"]
		case 146: // げ
			return tbl["ｹﾞ"]
		case 147: // こ
			return tbl["ｺ"]
		case 148: // ご
			return tbl["ｺﾞ"]
		case 149: // さ
			return tbl["ｻ"]
		case 150: // ざ
			return tbl["ｻﾞ"]
		case 151: // し
			return tbl["ｼ"]
		case 152: // じ
			return tbl["ｼﾞ"]
		case 153: // す
			return tbl["ｽ"]
		case 154: // ず
			return tbl["ｽﾞ"]
		case 155: // せ
			return tbl["ｾ"]
		case 156: // ぜ
			return tbl["ｾﾞ"]
		case 157: // そ
			return tbl["ｿ"]
		case 158: // ぞ
			return tbl["ｿﾞ"]
		case 159: // た
			return tbl["ﾀ"]
		case 160: // だ
			return tbl["ﾀﾞ"]
		case 161: // ち
			return tbl["ﾁ"]
		case 162: // ぢ
			return tbl["ﾁﾞ"]
		case 163: // っ
			return tbl["ｯ"]
		case 164: // つ
			return tbl["ﾂ"]
		case 165: // づ
			return tbl["ﾂﾞ"]
		case 166: // て
			return tbl["ﾃ"]
		case 167: // で
			return tbl["ﾃﾞ"]
		case 168: // と
			return tbl["ﾄ"]
		case 169: // ど
			return tbl["ﾄﾞ"]
		case 170: // な
			return tbl["ﾅ"]
		case 171: // に
			return tbl["ﾆ"]
		case 172: // ぬ
			return tbl["ﾇ"]
		case 173: // ね
			return tbl["ﾈ"]
		case 174: // の
			return tbl["ﾉ"]
		case 175: // は
			return tbl["ﾊ"]
		case 176: // ば
			return tbl["ﾊﾞ"]
		case 177: // ぱ
			return tbl["ﾊﾟ"]
		case 178: // ひ
			return tbl["ﾋ"]
		case 179: // び
			return tbl["ﾋﾞ"]
		case 180: // ぴ
			return tbl["ﾋﾟ"]
		case 181: // ふ
			return tbl["ﾌ"]
		case 182: // ぶ
			return tbl["ﾌﾞ"]
		case 183: // ぷ
			return tbl["ﾌﾟ"]
		case 184: // へ
			return tbl["ﾍ"]
		case 185: // べ
			return tbl["ﾍﾞ"]
		case 186: // ぺ
			return tbl["ﾍﾟ"]
		case 187: // ほ
			return tbl["ﾎ"]
		case 188: // ぼ
			return tbl["ﾎﾞ"]
		case 189: // ぽ
			return tbl["ﾎﾟ"]
		case 190: // ま
			return tbl["ﾏ"]
		case 191: // み
			return tbl["ﾐ"]
		}
	} else if w.val[1] == 130 && w.val[0] == 227 {
		switch w.val[2] {
		case 128: // む
			return tbl["ﾑ"]
		case 129: // め
			return tbl["ﾒ"]
		case 130: // も
			return tbl["ﾓ"]
		case 131: // ゃ
			return tbl["ｬ"]
		case 132: // や
			return tbl["ﾔ"]
		case 133: // ゅ
			return tbl["ｭ"]
		case 134: // ゆ
			return tbl["ﾕ"]
		case 135: // ょ
			return tbl["ｮ"]
		case 136: // よ
			return tbl["ﾖ"]
		case 137: // ら
			return tbl["ﾗ"]
		case 138: // り
			return tbl["ﾘ"]
		case 139: // る
			return tbl["ﾙ"]
		case 140: // れ
			return tbl["ﾚ"]
		case 141: // ろ
			return tbl["ﾛ"]
		case 143: // わ
			return tbl["ﾜ"]
		case 144: // ゐ
			return tbl["ｲ"]
		case 145: // ゑ
			return tbl["ｴ"]
		case 146: // を
			return tbl["ｦ"]
		case 147: // ん
			return tbl["ﾝ"]
		}
	} else if w.val[2] == 188 && w.val[1] == 131 && w.val[0] == 227 {
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
	if w.val[1] == 189 && w.val[0] == 239 {
		switch w.val[2] {
		case 166: // ｦ
			return tbl["を"]
		case 167: // ｧ
			return tbl["ぁ"]
		case 168: // ｨ
			return tbl["ぃ"]
		case 169: // ｩ
			return tbl["ぅ"]
		case 170: // ｪ
			return tbl["ぇ"]
		case 171: // ｫ
			return tbl["ぉ"]
		case 172: // ｬ
			return tbl["ゃ"]
		case 173: // ｭ
			return tbl["ゅ"]
		case 174: // ｮ
			return tbl["ょ"]
		case 175: // ｯ
			return tbl["っ"]
		case 176: // ｰ
			return tbl["ー"]
		case 177: // ｱ
			return tbl["あ"]
		case 178: // ｲ
			return tbl["い"]
		case 179: // ｳ
			return tbl["う"]
		case 180: // ｴ
			return tbl["え"]
		case 181: // ｵ
			return tbl["お"]
		case 182: // ｶ
			if isVoiced(w) {
				return tbl["が"]
			} else {
				return tbl["か"]
			}
		case 183: // ｷ
			if isVoiced(w) {
				return tbl["ぎ"]
			} else {
				return tbl["き"]
			}
		case 184: // ｸ
			if isVoiced(w) {
				return tbl["ぐ"]
			} else {
				return tbl["く"]
			}
		case 185: // ｹ
			if isVoiced(w) {
				return tbl["げ"]
			} else {
				return tbl["け"]
			}
		case 186: // ｺ
			if isVoiced(w) {
				return tbl["ご"]
			} else {
				return tbl["こ"]
			}
		case 187: // ｻ
			if isVoiced(w) {
				return tbl["ざ"]
			} else {
				return tbl["さ"]
			}
		case 188: // ｼ
			if isVoiced(w) {
				return tbl["じ"]
			} else {
				return tbl["し"]
			}
		case 189: // ｽ
			if isVoiced(w) {
				return tbl["ず"]
			} else {
				return tbl["す"]
			}
		case 190: // ｾ
			if isVoiced(w) {
				return tbl["ぜ"]
			} else {
				return tbl["せ"]
			}
		case 191: // ｿ
			if isVoiced(w) {
				return tbl["ぞ"]
			} else {
				return tbl["そ"]
			}
		}
	} else if w.val[1] == 190 && w.val[0] == 239 {
		switch w.val[2] {
		case 128: // ﾀ
			if isVoiced(w) {
				return tbl["だ"]
			} else {
				return tbl["た"]
			}
		case 129: // ﾁ
			if isVoiced(w) {
				return tbl["ぢ"]
			} else {
				return tbl["ち"]
			}
		case 130: // ﾂ
			if isVoiced(w) {
				return tbl["づ"]
			} else {
				return tbl["つ"]
			}
		case 131: // ﾃ
			if isVoiced(w) {
				return tbl["で"]
			} else {
				return tbl["て"]
			}
		case 132: // ﾄ
			if isVoiced(w) {
				return tbl["ど"]
			} else {
				return tbl["と"]
			}
		case 133: // ﾅ
			return tbl["な"]
		case 134: // ﾆ
			return tbl["に"]
		case 135: // ﾇ
			return tbl["ぬ"]
		case 136: // ﾈ
			return tbl["ね"]
		case 137: // ﾉ
			return tbl["の"]
		case 138: // ﾊ
			if isDevoiced(w) {
				return tbl["ぱ"]
			} else if isVoiced(w) {
				return tbl["ば"]
			} else {
				return tbl["は"]
			}
		case 139: // ﾋ
			if isDevoiced(w) {
				return tbl["ぴ"]
			} else if isVoiced(w) {
				return tbl["び"]
			} else {
				return tbl["ひ"]
			}
		case 140: // ﾌ
			if isDevoiced(w) {
				return tbl["ぷ"]
			} else if isVoiced(w) {
				return tbl["ぶ"]
			} else {
				return tbl["ふ"]
			}
		case 141: // ﾍ
			if isDevoiced(w) {
				return tbl["ぺ"]
			} else if isVoiced(w) {
				return tbl["べ"]
			} else {
				return tbl["へ"]
			}
		case 142: // ﾎ
			if isDevoiced(w) {
				return tbl["ぽ"]
			} else if isVoiced(w) {
				return tbl["ぼ"]
			} else {
				return tbl["ほ"]
			}
		case 143: // ﾏ
			return tbl["ま"]
		case 144: // ﾐ
			return tbl["み"]
		case 145: // ﾑ
			return tbl["む"]
		case 146: // ﾒ
			return tbl["め"]
		case 147: // ﾓ
			return tbl["も"]
		case 148: // ﾔ
			return tbl["や"]
		case 149: // ﾕ
			return tbl["ゆ"]
		case 150: // ﾖ
			return tbl["よ"]
		case 151: // ﾗ
			return tbl["ら"]
		case 152: // ﾘ
			return tbl["り"]
		case 153: // ﾙ
			return tbl["る"]
		case 154: // ﾚ
			return tbl["れ"]
		case 155: // ﾛ
			return tbl["ろ"]
		case 156: // ﾜ
			return tbl["わ"]
		case 157: // ﾝ
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
	if w.val[0] == 227 {
		if w.val[1] == 130 { // ァ-タ
			if w.val[2] >= 161 && w.val[2] <= 191 {
				return []byte{227, 129, w.val[2] - 32}
			}
		} else if w.val[1] == 131 { // ダ-ン
			if w.val[2] >= 128 && w.val[2] <= 159 { // ダ-ミ
				return []byte{227, 129, w.val[2] + 32}
			} else if w.val[2] >= 160 && w.val[2] <= 179 { // ム-ン
				return []byte{227, 130, w.val[2] - 32}
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
	if w.val[0] == 227 {
		if w.val[1] == 129 { // ぁ-み
			if w.val[2] >= 129 && w.val[2] <= 159 { // ぁ-た
				return []byte{227, 130, w.val[2] + 32}
			} else if w.val[2] >= 160 && w.val[2] <= 191 { // だ-み
				return []byte{227, 131, w.val[2] - 32}
			}
		} else if w.val[1] == 130 { // む-ん
			if w.val[2] >= 128 && w.val[2] <= 147 {
				return []byte{227, 131, w.val[2] + 32}
			}
		}
	}
	return w.val
}

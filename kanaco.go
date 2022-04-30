package kanaco

import (
	"bufio"
	"bytes"
	"io"
)

const (
	hankaku      = 1
	zenkaku      = 2
	space        = 4
	numeric      = 8
	alphabet     = 16
	alphanumeric = 32
	hiragana     = 64
	katakana     = 128
	voiced       = 256
	devoiced     = 512
	uppercase    = 1024
	lowercase    = 2048
	asIs         = 8192
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

/**
  map[string][0][]byte -> Byte order of myself
  map[string][1][]byte -> Byte order of myself
  map[string][2][]byte -> Byte order of myself
*/
var tbl map[string][2][]byte = map[string][2][]byte{
	"ぁ":  {[]byte("ァ"), []byte("ｧ")},
	"あ":  {[]byte("ア"), []byte("ｱ")},
	"ぃ":  {[]byte("ィ"), []byte("ｨ")},
	"い":  {[]byte("イ"), []byte("ｲ")},
	"ぅ":  {[]byte("ゥ"), []byte("ｩ")},
	"う":  {[]byte("ウ"), []byte("ｳ")},
	"ぇ":  {[]byte("ェ"), []byte("ｪ")},
	"え":  {[]byte("エ"), []byte("ｴ")},
	"ぉ":  {[]byte("ォ"), []byte("ｫ")},
	"お":  {[]byte("オ"), []byte("ｵ")},
	"か":  {[]byte("カ"), []byte("ｶ")},
	"が":  {[]byte("ガ"), []byte("ｶﾞ")},
	"き":  {[]byte("キ"), []byte("ｷ")},
	"ぎ":  {[]byte("ギ"), []byte("ｷﾞ")},
	"く":  {[]byte("ク"), []byte("ｸ")},
	"ぐ":  {[]byte("グ"), []byte("ｸﾞ")},
	"け":  {[]byte("ケ"), []byte("ｹ")},
	"げ":  {[]byte("ゲ"), []byte("ｹﾞ")},
	"こ":  {[]byte("コ"), []byte("ｺ")},
	"ご":  {[]byte("ゴ"), []byte("ｺﾞ")},
	"さ":  {[]byte("サ"), []byte("ｻ")},
	"ざ":  {[]byte("ザ"), []byte("ｻﾞ")},
	"し":  {[]byte("シ"), []byte("ｼ")},
	"じ":  {[]byte("ジ"), []byte("ｼﾞ")},
	"す":  {[]byte("ス"), []byte("ｽ")},
	"ず":  {[]byte("ズ"), []byte("ｽﾞ")},
	"せ":  {[]byte("セ"), []byte("ｾ")},
	"ぜ":  {[]byte("ゼ"), []byte("ｾﾞ")},
	"そ":  {[]byte("ソ"), []byte("ｿ")},
	"ぞ":  {[]byte("ゾ"), []byte("ｿﾞ")},
	"た":  {[]byte("タ"), []byte("ﾀ")},
	"だ":  {[]byte("ダ"), []byte("ﾀﾞ")},
	"ち":  {[]byte("チ"), []byte("ﾁ")},
	"ぢ":  {[]byte("ヂ"), []byte("ﾁﾞ")},
	"っ":  {[]byte("ッ"), []byte("ｯ")},
	"つ":  {[]byte("ツ"), []byte("ﾂ")},
	"づ":  {[]byte("ヅ"), []byte("ﾂﾞ")},
	"て":  {[]byte("テ"), []byte("ﾃ")},
	"で":  {[]byte("デ"), []byte("ﾃﾞ")},
	"と":  {[]byte("ト"), []byte("ﾄ")},
	"ど":  {[]byte("ド"), []byte("ﾄﾞ")},
	"な":  {[]byte("ナ"), []byte("ﾅ")},
	"に":  {[]byte("ニ"), []byte("ﾆ")},
	"ぬ":  {[]byte("ヌ"), []byte("ﾇ")},
	"ね":  {[]byte("ネ"), []byte("ﾈ")},
	"の":  {[]byte("ノ"), []byte("ﾉ")},
	"は":  {[]byte("ハ"), []byte("ﾊ")},
	"ば":  {[]byte("バ"), []byte("ﾊﾞ")},
	"ぱ":  {[]byte("パ"), []byte("ﾊﾟ")},
	"ひ":  {[]byte("ヒ"), []byte("ﾋ")},
	"び":  {[]byte("ビ"), []byte("ﾋﾞ")},
	"ぴ":  {[]byte("ピ"), []byte("ﾋﾟ")},
	"ふ":  {[]byte("フ"), []byte("ﾌ")},
	"ぶ":  {[]byte("ブ"), []byte("ﾌﾞ")},
	"ぷ":  {[]byte("プ"), []byte("ﾌﾟ")},
	"へ":  {[]byte("ヘ"), []byte("ﾍ")},
	"べ":  {[]byte("ベ"), []byte("ﾍﾞ")},
	"ぺ":  {[]byte("ペ"), []byte("ﾍﾟ")},
	"ほ":  {[]byte("ホ"), []byte("ﾎ")},
	"ぼ":  {[]byte("ボ"), []byte("ﾎﾞ")},
	"ぽ":  {[]byte("ポ"), []byte("ﾎﾟ")},
	"ま":  {[]byte("マ"), []byte("ﾏ")},
	"み":  {[]byte("ミ"), []byte("ﾐ")},
	"む":  {[]byte("ム"), []byte("ﾑ")},
	"め":  {[]byte("メ"), []byte("ﾒ")},
	"も":  {[]byte("モ"), []byte("ﾓ")},
	"ゃ":  {[]byte("ャ"), []byte("ｬ")},
	"や":  {[]byte("ヤ"), []byte("ﾔ")},
	"ゅ":  {[]byte("ュ"), []byte("ｭ")},
	"ゆ":  {[]byte("ユ"), []byte("ﾕ")},
	"ょ":  {[]byte("ョ"), []byte("ｮ")},
	"よ":  {[]byte("ヨ"), []byte("ﾖ")},
	"ら":  {[]byte("ラ"), []byte("ﾗ")},
	"り":  {[]byte("リ"), []byte("ﾘ")},
	"る":  {[]byte("ル"), []byte("ﾙ")},
	"れ":  {[]byte("レ"), []byte("ﾚ")},
	"ろ":  {[]byte("ロ"), []byte("ﾛ")},
	"ゎ":  {[]byte("ヮ"), []byte("ゎ")},
	"わ":  {[]byte("ワ"), []byte("ﾜ")},
	"ゐ":  {[]byte("ヰ"), []byte("ゐ")},
	"ゑ":  {[]byte("ヱ"), []byte("ゑ")},
	"を":  {[]byte("ヲ"), []byte("ｦ")},
	"ん":  {[]byte("ン"), []byte("ﾝ")},
	"ァ":  {[]byte("ぁ"), []byte("ｧ")},
	"ア":  {[]byte("あ"), []byte("ｱ")},
	"ィ":  {[]byte("ぃ"), []byte("ｨ")},
	"イ":  {[]byte("い"), []byte("ｲ")},
	"ゥ":  {[]byte("ぅ"), []byte("ｩ")},
	"ウ":  {[]byte("う"), []byte("ｳ")},
	"ェ":  {[]byte("ぇ"), []byte("ｪ")},
	"エ":  {[]byte("え"), []byte("ｴ")},
	"ォ":  {[]byte("ぉ"), []byte("ｫ")},
	"オ":  {[]byte("お"), []byte("ｵ")},
	"カ":  {[]byte("か"), []byte("ｶ")},
	"ガ":  {[]byte("が"), []byte("ｶﾞ")},
	"キ":  {[]byte("き"), []byte("ｷ")},
	"ギ":  {[]byte("ぎ"), []byte("ｷﾞ")},
	"ク":  {[]byte("く"), []byte("ｸ")},
	"グ":  {[]byte("ぐ"), []byte("ｸﾞ")},
	"ケ":  {[]byte("け"), []byte("ｹ")},
	"ゲ":  {[]byte("げ"), []byte("ｹﾞ")},
	"コ":  {[]byte("こ"), []byte("ｺ")},
	"ゴ":  {[]byte("ご"), []byte("ｺﾞ")},
	"サ":  {[]byte("さ"), []byte("ｻ")},
	"ザ":  {[]byte("ざ"), []byte("ｻﾞ")},
	"シ":  {[]byte("し"), []byte("ｼ")},
	"ジ":  {[]byte("じ"), []byte("ｼﾞ")},
	"ス":  {[]byte("す"), []byte("ｽ")},
	"ズ":  {[]byte("ず"), []byte("ｽﾞ")},
	"セ":  {[]byte("せ"), []byte("ｾ")},
	"ゼ":  {[]byte("ぜ"), []byte("ｾﾞ")},
	"ソ":  {[]byte("そ"), []byte("ｿ")},
	"ゾ":  {[]byte("ぞ"), []byte("ｿﾞ")},
	"タ":  {[]byte("た"), []byte("ﾀ")},
	"ダ":  {[]byte("だ"), []byte("ﾀﾞ")},
	"チ":  {[]byte("ち"), []byte("ﾁ")},
	"ヂ":  {[]byte("ぢ"), []byte("ﾁﾞ")},
	"ッ":  {[]byte("っ"), []byte("ｯ")},
	"ツ":  {[]byte("つ"), []byte("ﾂ")},
	"ヅ":  {[]byte("づ"), []byte("ﾂﾞ")},
	"テ":  {[]byte("て"), []byte("ﾃ")},
	"デ":  {[]byte("で"), []byte("ﾃﾞ")},
	"ト":  {[]byte("と"), []byte("ﾄ")},
	"ド":  {[]byte("ど"), []byte("ﾄﾞ")},
	"ナ":  {[]byte("な"), []byte("ﾅ")},
	"ニ":  {[]byte("に"), []byte("ﾆ")},
	"ヌ":  {[]byte("ぬ"), []byte("ﾇ")},
	"ネ":  {[]byte("ね"), []byte("ﾈ")},
	"ノ":  {[]byte("の"), []byte("ﾉ")},
	"ハ":  {[]byte("は"), []byte("ﾊ")},
	"バ":  {[]byte("ば"), []byte("ﾊﾞ")},
	"パ":  {[]byte("ぱ"), []byte("ﾊﾟ")},
	"ヒ":  {[]byte("ひ"), []byte("ﾋ")},
	"ビ":  {[]byte("び"), []byte("ﾋﾞ")},
	"ピ":  {[]byte("ぴ"), []byte("ﾋﾟ")},
	"フ":  {[]byte("ふ"), []byte("ﾌ")},
	"ブ":  {[]byte("ぶ"), []byte("ﾌﾞ")},
	"プ":  {[]byte("ぷ"), []byte("ﾌﾟ")},
	"ヘ":  {[]byte("へ"), []byte("ﾍ")},
	"ベ":  {[]byte("べ"), []byte("ﾍﾞ")},
	"ペ":  {[]byte("ぺ"), []byte("ﾍﾟ")},
	"ホ":  {[]byte("ほ"), []byte("ﾎ")},
	"ボ":  {[]byte("ぼ"), []byte("ﾎﾞ")},
	"ポ":  {[]byte("ぽ"), []byte("ﾎﾟ")},
	"マ":  {[]byte("ま"), []byte("ﾏ")},
	"ミ":  {[]byte("み"), []byte("ﾐ")},
	"ム":  {[]byte("む"), []byte("ﾑ")},
	"メ":  {[]byte("め"), []byte("ﾒ")},
	"モ":  {[]byte("も"), []byte("ﾓ")},
	"ャ":  {[]byte("ゃ"), []byte("ｬ")},
	"ヤ":  {[]byte("や"), []byte("ﾔ")},
	"ュ":  {[]byte("ゅ"), []byte("ｭ")},
	"ユ":  {[]byte("ゆ"), []byte("ﾕ")},
	"ョ":  {[]byte("ょ"), []byte("ｮ")},
	"ヨ":  {[]byte("よ"), []byte("ﾖ")},
	"ラ":  {[]byte("ら"), []byte("ﾗ")},
	"リ":  {[]byte("り"), []byte("ﾘ")},
	"ル":  {[]byte("る"), []byte("ﾙ")},
	"レ":  {[]byte("れ"), []byte("ﾚ")},
	"ロ":  {[]byte("ろ"), []byte("ﾛ")},
	"ヮ":  {[]byte("ゎ"), []byte("ヮ")},
	"ワ":  {[]byte("わ"), []byte("ﾜ")},
	"ヰ":  {[]byte("ゐ"), []byte("ヰ")},
	"ヱ":  {[]byte("ゑ"), []byte("ヱ")},
	"ヲ":  {[]byte("を"), []byte("ｦ")},
	"ン":  {[]byte("ん"), []byte("ﾝ")},
	"ｧ":  {[]byte("ぁ"), []byte("ァ")},
	"ｱ":  {[]byte("あ"), []byte("ア")},
	"ｨ":  {[]byte("ぃ"), []byte("ィ")},
	"ｲ":  {[]byte("い"), []byte("イ")},
	"ｩ":  {[]byte("ぅ"), []byte("ゥ")},
	"ｳ":  {[]byte("う"), []byte("ウ")},
	"ｪ":  {[]byte("ぇ"), []byte("ェ")},
	"ｴ":  {[]byte("え"), []byte("エ")},
	"ｫ":  {[]byte("ぉ"), []byte("ォ")},
	"ｵ":  {[]byte("お"), []byte("オ")},
	"ｶ":  {[]byte("か"), []byte("カ")},
	"ｶﾞ": {[]byte("が"), []byte("ガ")},
	"ｷ":  {[]byte("き"), []byte("キ")},
	"ｷﾞ": {[]byte("ぎ"), []byte("ギ")},
	"ｸ":  {[]byte("く"), []byte("ク")},
	"ｸﾞ": {[]byte("ぐ"), []byte("グ")},
	"ｹ":  {[]byte("け"), []byte("ケ")},
	"ｹﾞ": {[]byte("げ"), []byte("ゲ")},
	"ｺ":  {[]byte("こ"), []byte("コ")},
	"ｺﾞ": {[]byte("ご"), []byte("ゴ")},
	"ｻ":  {[]byte("さ"), []byte("サ")},
	"ｻﾞ": {[]byte("ざ"), []byte("ザ")},
	"ｼ":  {[]byte("し"), []byte("シ")},
	"ｼﾞ": {[]byte("じ"), []byte("ジ")},
	"ｽ":  {[]byte("す"), []byte("ス")},
	"ｽﾞ": {[]byte("ず"), []byte("ズ")},
	"ｾ":  {[]byte("せ"), []byte("セ")},
	"ｾﾞ": {[]byte("ぜ"), []byte("ゼ")},
	"ｿ":  {[]byte("そ"), []byte("ソ")},
	"ｿﾞ": {[]byte("ぞ"), []byte("ゾ")},
	"ﾀ":  {[]byte("た"), []byte("タ")},
	"ﾀﾞ": {[]byte("だ"), []byte("ダ")},
	"ﾁ":  {[]byte("ち"), []byte("チ")},
	"ﾁﾞ": {[]byte("ぢ"), []byte("ヂ")},
	"ｯ":  {[]byte("っ"), []byte("ッ")},
	"ﾂ":  {[]byte("つ"), []byte("ツ")},
	"ﾂﾞ": {[]byte("づ"), []byte("ヅ")},
	"ﾃ":  {[]byte("て"), []byte("テ")},
	"ﾃﾞ": {[]byte("で"), []byte("デ")},
	"ﾄ":  {[]byte("と"), []byte("ト")},
	"ﾄﾞ": {[]byte("ど"), []byte("ド")},
	"ﾅ":  {[]byte("な"), []byte("ナ")},
	"ﾆ":  {[]byte("に"), []byte("ニ")},
	"ﾇ":  {[]byte("ぬ"), []byte("ヌ")},
	"ﾈ":  {[]byte("ね"), []byte("ネ")},
	"ﾉ":  {[]byte("の"), []byte("ノ")},
	"ﾊ":  {[]byte("は"), []byte("ハ")},
	"ﾊﾞ": {[]byte("ば"), []byte("バ")},
	"ﾊﾟ": {[]byte("ぱ"), []byte("パ")},
	"ﾋ":  {[]byte("ひ"), []byte("ヒ")},
	"ﾋﾞ": {[]byte("び"), []byte("ビ")},
	"ﾋﾟ": {[]byte("ぴ"), []byte("ピ")},
	"ﾌ":  {[]byte("ふ"), []byte("フ")},
	"ﾌﾞ": {[]byte("ぶ"), []byte("ブ")},
	"ﾌﾟ": {[]byte("ぷ"), []byte("プ")},
	"ﾍ":  {[]byte("へ"), []byte("ヘ")},
	"ﾍﾞ": {[]byte("べ"), []byte("ベ")},
	"ﾍﾟ": {[]byte("ぺ"), []byte("ペ")},
	"ﾎ":  {[]byte("ほ"), []byte("ホ")},
	"ﾎﾞ": {[]byte("ぼ"), []byte("ボ")},
	"ﾎﾟ": {[]byte("ぽ"), []byte("ポ")},
	"ﾏ":  {[]byte("ま"), []byte("マ")},
	"ﾐ":  {[]byte("み"), []byte("ミ")},
	"ﾑ":  {[]byte("む"), []byte("ム")},
	"ﾒ":  {[]byte("め"), []byte("メ")},
	"ﾓ":  {[]byte("も"), []byte("モ")},
	"ｬ":  {[]byte("ゃ"), []byte("ャ")},
	"ﾔ":  {[]byte("や"), []byte("ヤ")},
	"ｭ":  {[]byte("ゅ"), []byte("ュ")},
	"ﾕ":  {[]byte("ゆ"), []byte("ユ")},
	"ｮ":  {[]byte("ょ"), []byte("ョ")},
	"ﾖ":  {[]byte("よ"), []byte("ヨ")},
	"ﾗ":  {[]byte("ら"), []byte("ラ")},
	"ﾘ":  {[]byte("り"), []byte("リ")},
	"ﾙ":  {[]byte("る"), []byte("ル")},
	"ﾚ":  {[]byte("れ"), []byte("レ")},
	"ﾛ":  {[]byte("ろ"), []byte("ロ")},
	"ﾜ":  {[]byte("わ"), []byte("ワ")},
	"ｦ":  {[]byte("を"), []byte("ヲ")},
	"ﾝ":  {[]byte("ん"), []byte("ン")},
}

var filters map[byte]func(w *word) []byte = map[byte]func(w *word) []byte{
	'r': convAsSmallR,
	'R': convAsLargeR,
	'n': convAsSmallN,
	'N': convAsLargeN,
	'a': convAsSmallA,
	'A': convAsLargeA,
	's': convAsSmallS,
	'S': convAsLargeS,
	'k': convAsSmallK,
	'K': convAsLargeK,
	'h': convAsSmallH,
	'H': convAsLargeH,
	'c': convAsSmallC,
	'C': convAsLargeC,
}

func Byte(b []byte, mode string) []byte {
	order := []byte{}
	for _, m := range []byte(mode) {
		flg := true
		for _, v := range order {
			if m == v {
				flg = false
				break
			}
		}
		if flg {
			order = append(order, m)
		}
	}
	byteCount := uint64(len(b))
	buf := make([]byte, 0, byteCount)
	i := uint64(0)
	for i < byteCount {
		w := extract(b[i:])
		val := w.val
		for _, o := range order {
			if _, ok := filters[o]; ok {
				val = filters[o](w)
				if bytes.Equal(w.val, val) {
					break
				} else {
					continue
				}
			}
		}
		buf = append(buf, val...)
		i += uint64(w.len)
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

func isVoiced(w *word) bool {
	if w.charType&voiced == voiced {
		return true
	} else {
		return false
	}
}

func isDevoiced(w *word) bool {
	if w.charType&devoiced == devoiced {
		return true
	} else {
		return false
	}
}

func isSpace(w *word) bool {
	if w.charType&space == space {
		return true
	} else {
		return false
	}
}

func isNum(w *word) bool {
	if w.charType&numeric == numeric {
		return true
	} else {
		return false
	}
}

func isAlpha(w *word) bool {
	if w.charType&alphabet == alphabet {
		return true
	} else {
		return false
	}
}

func isAlphaNum(w *word) bool {
	if w.charType&alphanumeric == alphanumeric {
		return true
	} else {
		return false
	}
}

func isHan(w *word) bool {
	if w.charType&hankaku == hankaku {
		return true
	} else {
		return false
	}
}

func isZen(w *word) bool {
	if w.charType&zenkaku == zenkaku {
		return true
	} else {
		return false
	}
}

func isHira(w *word) bool {
	if w.charType&hiragana == hiragana {
		return true
	} else {
		return false
	}
}

func isKata(w *word) bool {
	if w.charType&katakana == katakana {
		return true
	} else {
		return false
	}
}

func isUpper(w *word) bool {
	if w.charType&uppercase == uppercase {
		return true
	} else {
		return false
	}
}

func isLower(w *word) bool {
	if w.charType&lowercase == lowercase {
		return true
	} else {
		return false
	}
}

func extract(b []byte) *word {
	if is1Byte(b) {
		if b[0] == 0x20 {
			return &word{b[0:1], hankaku + space, 1}
		} else if b[0] >= 0x30 && b[0] <= 0x39 {
			return &word{b[0:1], hankaku + numeric + alphanumeric, 1}
		} else if b[0] >= 0x41 && b[0] <= 0x5a {
			return &word{b[0:1], hankaku + alphabet + alphanumeric + uppercase, 1}
		} else if b[0] >= 0x61 && b[0] <= 0x7a {
			return &word{b[0:1], hankaku + alphabet + alphanumeric + lowercase, 1}
		} else if b[0] >= 0x21 && b[0] <= 0x7d && b[0] != 0x22 && b[0] != 0x27 && b[0] != 0x5c {
			return &word{b[0:1], hankaku + alphabet, 1}
		} else if b[0] < 0x80 {
			return &word{b[0:1], asIs, 1}
		}
	} else if is2Bytes(b) {
		return &word{b[0:2], asIs, 2}
	} else if is3Bytes(b) {
		if b[0] == 0xef {
			switch b[1] {
			case 0xbc:
				if b[2] >= 0x90 && b[2] <= 0x99 { // ０ ～ ９
					return &word{b[0:3], zenkaku + numeric, 3}
				} else if b[2] >= 0xa1 && b[2] <= 0xba { // Ａ ～ Ｚ
					return &word{b[0:3], zenkaku + alphabet + alphanumeric + uppercase, 3}
				} else if b[2] >= 0x81 && b[2] == 0xbf && b[2] != 0x82 && b[2] != 0x87 && b[2] != 0xbc {
					return &word{b[0:3], zenkaku + alphanumeric, 3}
				} else {
					return &word{b[0:3], asIs, 3}
				}
			case 0xbd:
				if b[2] >= 0x81 && b[2] <= 0x9a {
					return &word{b[0:3], zenkaku + alphabet + alphanumeric + lowercase, 3}
				} else if b[2] >= 0x80 && b[2] <= 0x9d { // ｀ ～ ｝
					return &word{b[0:3], zenkaku + alphanumeric, 3}
				} else if b[2] >= 0xa6 && b[2] <= 0xbf {
					if len(b) >= 6 {
						if b[2] >= 0xb6 && b[2] <= 0xbf { // ｶ ～ ｿ
							if b[3] == 0xef && b[4] == 0xbe && b[5] == 0x9e { // 濁点
								return &word{b[0:6], hankaku + katakana + voiced, 6}
							}
						}
					}
					return &word{b[0:3], hankaku + katakana, 3}
				}
			case 0xbe:
				if b[2] >= 0x80 && b[2] <= 0x84 { // ﾀ ～ ﾄ
					if len(b) >= 6 {
						if b[3] == 0xef && b[4] == 0xbe && b[5] == 0x9e {
							return &word{b[0:6], hankaku + katakana + voiced, 6}
						}
					}
					return &word{b[0:3], hankaku + katakana, 3}
				} else if b[2] >= 0x8a && b[2] <= 0x8d { // ﾊ ～ ﾎ
					if len(b) >= 6 {
						if b[3] == 0xef && b[4] == 0xbe && b[5] == 0x9e {
							return &word{b[0:6], hankaku + katakana + voiced, 6}
						} else if b[3] == 0xef && b[4] == 0xbe && b[5] == 0x9f {
							return &word{b[0:6], hankaku + katakana + devoiced, 6}
						}
					}
					return &word{b[0:3], hankaku + katakana, 3}
				} else if b[2] >= 0x85 && b[2] <= 0x9d { // ﾅ ～ ﾝ
					return &word{b[0:3], hankaku + katakana, 3}
				}
			}
			return &word{b[0:3], asIs, 3}
		} else if b[0] == 0xe3 {
			// 全ひら・全カタ・全スペース
			switch b[1] {
			case 0x80:
				if b[2] == 0x80 { // スペース
					return &word{b[0:3], zenkaku + space, 3}
				}
			case 0x81:
				if b[2] >= 0x81 && b[2] <= 0xbf { // ぁ-み
					return &word{b[0:3], zenkaku + hiragana, 3}
				}
			case 0x82:
				if b[2] >= 0x80 && b[2] <= 0x93 { // む-ん
					return &word{b[0:3], zenkaku + hiragana, 3}
				} else if b[2] >= 0xa1 && b[2] <= 0xbf { // ァ-タ
					return &word{b[0:3], zenkaku + katakana, 3}
				}
			case 0x83:
				if b[2] >= 0x80 && b[2] <= 0xb3 { // チ-ン
					return &word{b[0:3], zenkaku + katakana, 3}
				}
			}
			return &word{b[0:3], asIs, 3}
		}
		return &word{b[0:3], asIs, 3}
	} else if is4Bytes(b) {
		return &word{b[0:4], asIs, 4}
	}
	return &word{b[0:1], asIs, 1}
}

/**
 * Hankaku Space -> Zenkaku Space
 */
func convAsLargeS(w *word) []byte {
	if isHan(w) && isSpace(w) {
		return []byte{0xe3, 0x80, 0x80}
	}
	return w.val
}

/**
 * Zenkaku Space -> Hankaku Space
 */
func convAsSmallS(w *word) []byte {
	if isZen(w) && isSpace(w) {
		return []byte{0x20}
	}
	return w.val
}

/**
 * Hankaku Numeric -> Zenkaku Numeric
 */
func convAsLargeN(w *word) []byte {
	if isHan(w) && isNum(w) {
		return []byte{0xef, 0xbc, 0x60 + w.val[0]}
	}
	return w.val
}

/**
 * Zenkaku Numeric -> Hankaku Numeric
 */
func convAsSmallN(w *word) []byte {
	if isZen(w) && isNum(w) {
		return []byte{w.val[2] - 0x60}
	}
	return w.val
}

/**
 * Hankaku Alphabet -> Zenkaku Alphabet
 */
func convAsLargeR(w *word) []byte {
	if isHan(w) && isAlpha(w) {
		// A-Z -> Ａ-Ｚ
		if isUpper(w) {
			return []byte{0xef, 0xbc, 0x60 + w.val[0]}
		}
		// a-z -> ａ-ｚ
		if isLower(w) {
			return []byte{0xef, 0xbd, 0x20 + w.val[0]}
		}
	}
	return w.val
}

/**
 * Zenkaku Alphabet -> Hankaku Alphabet
 */
func convAsSmallR(w *word) []byte {
	if isZen(w) && isAlpha(w) {
		// Ａ-Ｚ -> A-Z
		if isUpper(w) {
			return []byte{w.val[2] - 0x60}
		}
		// ａ-ｚ -> a-z
		if isLower(w) {
			return []byte{w.val[2] - 0x20}
		}
	}
	return w.val
}

/**
 * Hankaku AlphaNumeric -> Zenkaku AlphaNumeric
 * !-}(Excluding ",',\)
 */
func convAsLargeA(w *word) []byte {
	if isHan(w) && isAlphaNum(w) {
		if w.val[0] >= 0x21 && w.val[0] <= 0x5f {
			return []byte{0xef, 0xbc, 0x60 + w.val[0]}
		} else if w.val[0] >= 0x60 && w.val[0] <= 0x7d {
			return []byte{0xef, 0xbd, 0x20 + w.val[0]}
		}
	}
	// fmt.Println(string(w.val), w.charType)
	return w.val
}

/**
 * Zenkaku AlphaNumeric -> Hankaku AlphaNumeric
 * !-}(Excluding ",',\)
 */
func convAsSmallA(w *word) []byte {
	if isZen(w) && isAlphaNum(w) {
		if w.val[1] == 0xbc && w.val[2] >= 0x81 && w.val[2] <= 0xbf {
			return []byte{w.val[2] - 0x60}
		} else if w.val[1] == 0xbd && w.val[2] >= 0x80 && w.val[2] <= 0x9d {
			return []byte{w.val[2] - 0x20}
		}
	}
	return w.val
}

/**
 * Zenkaku Katakana -> Hankaku Katakana
 */
func convAsSmallK(w *word) []byte {
	if isZen(w) && isKata(w) {
		s := string(w.val)
		if _, ok := tbl[s]; ok {
			return tbl[s][1]
		}
	}
	return w.val
}

/**
 * Hankaku Katakana -> Zenkaku Katakana
 */
func convAsLargeK(w *word) []byte {
	if isHan(w) && isKata(w) {
		s := string(w.val)
		if _, ok := tbl[s]; ok {
			return tbl[s][1]
		}
	}
	return w.val
}

/**
 * Zenkaku Hiragana -> Hankaku Katakana
 */
func convAsSmallH(w *word) []byte {
	if isZen(w) && isHira(w) {
		s := string(w.val)
		if _, ok := tbl[s]; ok {
			return tbl[s][1]
		}
	}
	return w.val
}

/**
 * Hankaku Katakana -> Zenkaku Hiragana
 */
func convAsLargeH(w *word) []byte {
	if isHan(w) && isKata(w) {
		s := string(w.val)
		if _, ok := tbl[s]; ok {
			return tbl[s][0]
		}
	}
	return w.val
}

/**
 * Zenkaku Katakana -> Zenkaku Hiragana
 */
func convAsSmallC(w *word) []byte {
	if isZen(w) && isKata(w) {
		if w.val[1] == 0x82 { // ァ-タ
			if w.val[2] >= 0xa1 && w.val[2] <= 0xbf {
				return []byte{0xe3, 0x81, w.val[2] - 0x20}
			}
		} else if w.val[1] == 0x83 { // ダ-ン
			if w.val[2] >= 0x80 && w.val[2] <= 0x9f { // ダ-ミ
				return []byte{0xe3, 0x81, w.val[2] + 0x20}
			} else if w.val[2] >= 0xa0 && w.val[2] <= 0xb3 { // ム-ン
				return []byte{0xe3, 0x82, w.val[2] - 0x20}
			}
		}
	}
	return w.val
}

/**
 * Zenkaku Hiragana -> Zenkaku Katakana
 */
func convAsLargeC(w *word) []byte {
	if isZen(w) && isHira(w) {
		if w.val[1] == 0x81 { // ぁ-み
			if w.val[2] >= 0x81 && w.val[2] <= 0x9f { // ぁ-た
				return []byte{0xe3, 0x82, w.val[2] + 0x20}
			} else if w.val[2] >= 0xa0 && w.val[2] <= 0xbf { // だ-み
				return []byte{0xe3, 0x83, w.val[2] - 0x20}
			}
		} else if w.val[1] == 0x82 { // む-ん
			if w.val[2] >= 0x80 && w.val[2] <= 0x93 {
				return []byte{0xe3, 0x83, w.val[2] + 0x20}
			}
		}
	}
	return w.val
}

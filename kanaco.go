package kanaco

import (
	"bufio"
	// "bytes"
	"fmt"
	"io"
)

const (
	notEligible = errors.New("Not Eligible")
)

type Reader struct {
	r    *bufio.Reader
	mode string
}

type Writer struct {
	w    *bufio.Writer
	mode string
}

type char struct {
	val  []byte
	mode []byte
	len  int
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
	"ゎ":  {[]byte("ヮ"), []byte("ﾜ")},
	"わ":  {[]byte("ワ"), []byte("ﾜ")},
	"ゐ":  {[]byte("ヰ"), []byte("ｲ")},
	"ゑ":  {[]byte("ヱ"), []byte("ｴ")},
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
	"ヮ":  {[]byte("ゎ"), []byte("ﾜ")},
	"ワ":  {[]byte("わ"), []byte("ﾜ")},
	"ヰ":  {[]byte("ゐ"), []byte("ｲ")},
	"ヱ":  {[]byte("ゑ"), []byte("ｴ")},
	"ヲ":  {[]byte("を"), []byte("ｦ")},
	"ン":  {[]byte("ん"), []byte("ﾝ")},
	"ヴ":  {[]byte("ヴ"), []byte("ｳﾞ")},
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
	"ｳﾞ": {[]byte("う゛"), []byte("ヴ")},
	"ﾞ":  {[]byte("゛"), []byte("゛")},
	"ﾟ":  {[]byte("゜"), []byte("゜")},
	"､":  {[]byte("、"), []byte("、")},
	"｡":  {[]byte("。"), []byte("。")},
	"｢":  {[]byte("「"), []byte("「")},
	"｣":  {[]byte("」"), []byte("」")},
	"ヽ":  {[]byte("ゝ"), []byte("ゝ")},
	"ヾ":  {[]byte("ゞ"), []byte("ゞ")},
	"ゝ":  {[]byte("ヽ"), []byte("ヽ")},
	"ゞ":  {[]byte("ヾ"), []byte("ヾ")},
	"･":  {[]byte("・"), []byte("・")},
	"・":  {[]byte("･"), []byte("･")},
	"ｰ":  {[]byte("ー"), []byte("ー")},
	"ー":  {[]byte("ｰ"), []byte("ｰ")},
}

var filters map[byte]func(c *char) []byte = map[byte]func(c *char) []byte{
	'r': smallR,
	'R': largeR,
	'n': smallN,
	'N': largeN,
	'a': smallA,
	'A': largeA,
	's': smallS,
	'S': largeS,
	'k': smallK,
	'K': largeK,
	'h': smallH,
	'H': largeH,
	'c': smallC,
	'C': largeC,
}

func Byte(b []byte, mode string) []byte {
	orders := []byte{}
	for _, m := range []byte(mode) {
		flg := true
		for _, v := range orders {
			if m == v {
				flg = false
				break
			}
		}
		if flg {
			orders = append(orders, m)
		}
	}
fmt.Println(string(orders[0]),string(orders[1]),string(orders[2]))
	byteCount := uint64(len(b))
	buf := make([]byte, 0, byteCount)
	i := uint64(0)
	for i < byteCount {
		c := extract(b[i:])
		val := c.val
		for _, o := range orders {
			val = filters[o](c)
// fmt.Printf("%s <-> %s\n", c.val, val)
			if len(val)>0 {
				continue
			}
		}
		buf = append(buf, val...)
		i += uint64(c.len)
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
	if err == io.EOF {
		return 0, err
	}
	if err != nil {
		return 0, err
	}
	line = Byte(line, r.mode)
	if len(p) < len(line) {
		return 0, fmt.Errorf("Buffer size is not enough")
	}
	n := copy(p, line)
	return n, nil
}

func NewWriter(w io.Writer, mode string) *Writer {
	writer := new(Writer)
	writer.w = bufio.NewWriter(w)
	writer.mode = mode
	return writer
}

func (w *Writer) Write(p []byte) (int, error) {
	buf := Byte(p, w.mode)
	return w.w.Write(buf)
}

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

func extract(b []byte) *char {
	if is1Byte(b) {
		if b[0] == 0x20 { // Space
			return &char{b[0:1], []byte{'S'}, 1}
		} else if b[0] >= 0x30 && b[0] <= 0x39 { // 0 - 9
			return &char{b[0:1], []byte{'N', 'A'}, 1}
		} else if b[0] >= 0x41 && b[0] <= 0x5a { // A - Z
			return &char{b[0:1], []byte{'R', 'A'}, 1}
		} else if b[0] >= 0x61 && b[0] <= 0x7a { // a - z
			return &char{b[0:1], []byte{'R', 'A'}, 1}
		} else if b[0] >= 0x21 && b[0] <= 0x7d && b[0] != 0x22 && b[0] != 0x27 && b[0] != 0x5c {
			return &char{b[0:1], []byte{'A'}, 1}
		}
		return &char{b[:1], []byte{}, 1}
	} else if is2Bytes(b) {
		return &char{b[:2], []byte{}, 2}
	} else if is3Bytes(b) {
		if b[0] == 0xef {
			switch b[1] {
			case 0xbc:
				if b[2] >= 0x90 && b[2] <= 0x99 { // ０ ～ ９
					return &char{b[:3], []byte{'n', 'a'}, 3}
				} else if b[2] >= 0xa1 && b[2] <= 0xba { // Ａ ～ Ｚ
					return &char{b[:3], []byte{'r', 'a'}, 3}
				} else if b[2] != 0x82 && b[2] != 0x87 && b[2] != 0xbc { // ＂ ＇ ＼ 以外
					return &char{b[:3], []byte{'a'}, 3}
				}
				return &char{b[0:3], []byte{}, 3}
			case 0xbd:
				if b[2] >= 0x81 && b[2] <= 0x9a { // ａ ～ ｚ
					return &char{b[:3], []byte{'r', 'a'}, 3}
				} else if b[2] >= 0x80 && b[2] <= 0x9d { // ｀ ～ ｝
					return &char{b[:3], []byte{'a'}, 3}
				} else if b[2] >= 0xa1 && b[2] <= 0xbf { // ｡ ｢ ｣ ､ ･ ｦ ～ ｿ
					if len(b) >= 6 {
						if b[2] == 0xb3 || (b[2] >= 0xb6 && b[2] <= 0xbf) { // ｳ, ｶ ～ ｿ
							if b[3] == 0xef && b[4] == 0xbe && b[5] == 0x9e { // 濁点
								return &char{b[:6], []byte{'K', 'H'}, 6}
							}
						}
					}
					return &char{b[:3], []byte{'K', 'H'}, 3}
				}
				return &char{b[:3], []byte{}, 3}
			case 0xbe:
				if b[2] >= 0x80 && b[2] <= 0x84 { // ﾀ ～ ﾄ
					if len(b) >= 6 {
						if b[3] == 0xef && b[4] == 0xbe && b[5] == 0x9e {
							return &char{b[:6], []byte{'K', 'H'}, 6}
						}
					}
					return &char{b[:3], []byte{'K', 'H'}, 3}
				} else if b[2] >= 0x8a && b[2] <= 0x8e { // ﾊ ～ ﾎ
					if len(b) >= 6 {
						if b[3] == 0xef && b[4] == 0xbe && b[5] == 0x9e { // 濁点
							return &char{b[:6], []byte{'K', 'H'}, 6}
						} else if b[3] == 0xef && b[4] == 0xbe && b[5] == 0x9f { // 半濁点
							return &char{b[:6], []byte{'K', 'H'}, 6}
						}
					}
					return &char{b[:3], []byte{'K', 'H'}, 3}
				} else if b[2] >= 0x85 && b[2] <= 0x9f { // ﾅ ～ ﾝﾞﾟ
					return &char{b[:3], []byte{'K', 'H'}, 3}
				}
			}
			return &char{b[0:3], []byte{}, 3}
		} else if b[0] == 0xe3 {
			// 全ひら・全カタ・全スペース
			switch b[1] {
			case 0x80:
				if b[2] == 0x80 { // スペース
					return &char{b[:3], []byte{'s'}, 3}
				}
			case 0x81:
				if b[2] >= 0x81 && b[2] <= 0xbf { // ぁ-み
					return &char{b[:3], []byte{'h', 'C'}, 3}
				}
			case 0x82:
				if b[2] >= 0x80 && b[2] <= 0x93 { // む-ん
					return &char{b[:3], []byte{'h', 'C'}, 3}
				} else if b[2] >= 0x9d && b[2] <= 0x9e { // ゝゞ
					return &char{b[:3], []byte{'c'}, 3}
				} else if b[2] >= 0xa1 && b[2] <= 0xbf { // ァ-タ
					return &char{b[:3], []byte{'k', 'c'}, 3}
				}
			case 0x83:
				if b[2] >= 0x80 && b[2] <= 0xb4 { // チ-ヴ
					return &char{b[:3], []byte{'k', 'c'}, 3}
				} else if b[2] >= 0xbb && b[2] <= 0xbc { // ・ー
					return &char{b[:3], []byte{'h', 'k'}, 3}
				} else if b[2] >= 0xbd && b[2] <= 0xbe { // ヽ ヾ
					return &char{b[:3], []byte{'c'}, 3}
				}
			}
			return &char{b[:3], []byte{}, 3}
		}
		return &char{b[:3], []byte{}, 3}
	} else if is4Bytes(b) {
		return &char{b[:4], []byte{}, 4}
	}
	return &char{b[:1], []byte{}, 1}
}

/**
 * Hankaku Space -> Zenkaku Space
 */
func largeS(c *char) ([]byte, error) {
	for _, m := range c.mode {
		if m != 'S' {
			continue
		}
		return []byte{0xe3, 0x80, 0x80}, nil
	}
	return []byte{}, notEligible
}

/**
 * Zenkaku Space -> Hankaku Space
 */
func smallS(c *char) ([]byte, error) {
	for _, m := range c.mode {
		if m != 's' {
			continue
		}
		return []byte{0x20}
	}
	return []byte{}, notEligible
}

/**
 * Hankaku Numeric -> Zenkaku Numeric
 */
func largeN(c *char) ([]byte, error) {
	for _, m := range c.mode {
		if m != 'N' {
			continue
		}
		return []byte{0xef, 0xbc, 0x60 + c.val[0]}
	}
	return []byte{}, notEligible
}

/**
 * Zenkaku Numeric -> Hankaku Numeric
 */
func smallN(c *char) ([]byte, error) {
	for _, m := range c.mode {
		if m != 'n' {
			continue
		}
		return []byte{c.val[2] - 0x60}
	}
	return []byte{}, notEligible
}

/**
 * Hankaku Alphabet -> Zenkaku Alphabet
 */
func largeR(c *char) ([]byte, error) {
	for _, m := range c.mode {
		if m != 'R' {
			continue
		}
		// A-Z -> Ａ-Ｚ
		if c.val[0] >= 0x41 && c.val[0] <= 0x5a {
			return []byte{0xef, 0xbc, 0x60 + c.val[0]}
		}
		// a-z -> ａ-ｚ
		if c.val[0] >= 0x61 && c.val[0] <= 0x7a {
			return []byte{0xef, 0xbd, 0x20 + c.val[0]}
		}
		break
	}
	return []byte{}, notEligible
}

/**
 * Zenkaku Alphabet -> Hankaku Alphabet
 */
func smallR(c *char) ([]byte, error) {
	for _, m := range c.mode {
		if m != 'r' {
			continue
		}
		// Ａ-Ｚ -> A-Z
		if c.val[2] >= 0xa1 && c.val[2] <= 0xba {
			return []byte{c.val[2] - 0x60}
		}
		// ａ-ｚ -> a-z
		if c.val[2] >= 0x81 && c.val[2] <= 0x9a {
			return []byte{c.val[2] - 0x20}
		}
		break
	}
	return []byte{}, notEligible
}

/**
 * Hankaku AlphaNumeric -> Zenkaku AlphaNumeric
 * !-}(Excluding ",',\)
 */
func largeA(c *char) ([]byte, error) {
	for _, m := range c.mode {
		if m != 'A' {
			continue
		}
		if c.val[0] >= 0x21 && c.val[0] <= 0x5f {
			return []byte{0xef, 0xbc, 0x60 + c.val[0]}
		} else if c.val[0] >= 0x60 && c.val[0] <= 0x7d {
			return []byte{0xef, 0xbd, 0x20 + c.val[0]}
		}
		break
	}
	return []byte{}, notEligible
}

/**
 * Zenkaku AlphaNumeric -> Hankaku AlphaNumeric
 * !-}(Excluding ",',\)
 */
func smallA(c *char) ([]byte, error) {
	for _, m := range c.mode {
		if m != 'a' {
			continue
		}
		if c.val[1] == 0xbc && c.val[2] >= 0x81 && c.val[2] <= 0xbf {
			return []byte{c.val[2] - 0x60}
		} else if c.val[1] == 0xbd && c.val[2] >= 0x80 && c.val[2] <= 0x9d {
			return []byte{c.val[2] - 0x20}
		}
		break
	}
	return []byte{}, notEligible
}

/**
 * Hankaku Katakana -> Zenkaku Katakana
 */
func largeK(c *char) ([]byte, error) {
	for _, m := range c.mode {
		if m != 'K' {
			continue
		}
		s := string(c.val)
		if _, ok := tbl[s]; ok {
			return tbl[s][1]
		}
		break
	}
	return []byte{}, notEligible
}

/**
 * Zenkaku Katakana -> Hankaku Katakana
 */
func smallK(c *char) ([]byte, error) {
	for _, m := range c.mode {
		if m != 'k' {
			continue
		}
		s := string(c.val)
		if _, ok := tbl[s]; ok {
			return tbl[s][1]
		}
		break
	}
	return []byte{}, notEligible
}

/**
 * Hankaku Katakana -> Zenkaku Hiragana
 */
func largeH(c *char) ([]byte, error) {
	for _, m := range c.mode {
		if m != 'H' {
			continue
		}
		s := string(c.val)
		if _, ok := tbl[s]; ok {
			return tbl[s][0]
		}
		break
	}
	return []byte{}, notEligible
}

/**
 * Zenkaku Hiragana -> Hankaku Katakana
 */
func smallH(c *char) ([]byte, error) {
	for _, m := range c.mode {
		if m != 'h' {
			continue
		}
		s := string(c.val)
		if _, ok := tbl[s]; ok {
			return tbl[s][1]
		}
		break
	}
	return []byte{}, notEligible
}

/**
 * Zenkaku Hiragana -> Zenkaku Katakana
 */
func largeC(c *char) ([]byte, error) {
	for _, m := range c.mode {
		if m != 'C' {
			continue
		}
		switch c.val[1] {
		case 0x81: // ぁ-み
			if c.val[2] >= 0x81 && c.val[2] <= 0x9f { // ぁ-た
				return []byte{0xe3, 0x82, c.val[2] + 0x20}
			} else if c.val[2] >= 0xa0 && c.val[2] <= 0xbf { // だ-み
				return []byte{0xe3, 0x83, c.val[2] - 0x20}
			}
		case 0x82: // む-ん
			if c.val[2] >= 0x80 && c.val[2] <= 0x93 {
				return []byte{0xe3, 0x83, c.val[2] + 0x20}
			}
		}
		break
	}
	return []byte{}, notEligible
}

/**
 * Zenkaku Katakana -> Zenkaku Hiragana
 */
func smallC(c *char) ([]byte, error) {
	for _, m := range c.mode {
		if m != 'c' {
			continue
		}
		switch c.val[1] {
		case 0x82: // ァ-タ
			if c.val[2] >= 0xa1 && c.val[2] <= 0xbf {
				return []byte{0xe3, 0x81, c.val[2] - 0x20}
			}
		case 0x83: // ダ-ン
			if c.val[2] >= 0x80 && c.val[2] <= 0x9f { // ダ-ミ
				return []byte{0xe3, 0x81, c.val[2] + 0x20}
			} else if c.val[2] >= 0xa0 && c.val[2] <= 0xb3 { // ム-ン
				return []byte{0xe3, 0x82, c.val[2] - 0x20}
			} else if c.val[2] >= 0xbd && c.val[2] <= 0xbe { // ヽヾ
				return []byte{0xe3, 0x82, c.val[2] - 0x20}
			}
		}
		break
	}
	return []byte{}, notEligible
}

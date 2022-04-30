package kanaco

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

const (
	bufSize = 4096
)

type testCase struct {
	arg      *word
	expected string
}

var ks [][3]string = [][3]string{
	{"ｧ", "ァ", "ぁ"}, {"ｨ", "ィ", "ぃ"}, {"ｩ", "ゥ", "ぅ"}, {"ｪ", "ェ", "ぇ"}, {"ｫ", "ォ", "ぉ"},
	{"ｱ", "ア", "あ"}, {"ｲ", "イ", "い"}, {"ｳ", "ウ", "う"}, {"ｴ", "エ", "え"}, {"ｵ", "オ", "お"},
	{"ｶ", "カ", "か"}, {"ｷ", "キ", "き"}, {"ｸ", "ク", "く"}, {"ｹ", "ケ", "け"}, {"ｺ", "コ", "こ"},
	{"ｶﾞ", "ガ", "が"}, {"ｷﾞ", "ギ", "ぎ"}, {"ｸﾞ", "グ", "ぐ"}, {"ｹﾞ", "ゲ", "げ"}, {"ｺﾞ", "ゴ", "ご"},
	{"ｻ", "サ", "さ"}, {"ｼ", "シ", "し"}, {"ｽ", "ス", "す"}, {"ｾ", "セ", "せ"}, {"ｿ", "ソ", "そ"},
	{"ｻﾞ", "ザ", "ざ"}, {"ｼﾞ", "ジ", "じ"}, {"ｽﾞ", "ズ", "ず"}, {"ｾﾞ", "ゼ", "ぜ"}, {"ｿﾞ", "ゾ", "ぞ"},
	{"ﾀ", "タ", "た"}, {"ﾁ", "チ", "ち"}, {"ﾂ", "ツ", "つ"}, {"ﾃ", "テ", "て"}, {"ﾄ", "ト", "と"},
	{"ﾀﾞ", "ダ", "だ"}, {"ﾁﾞ", "ヂ", "ぢ"}, {"ﾂﾞ", "ヅ", "づ"}, {"ﾃﾞ", "デ", "で"}, {"ﾄﾞ", "ド", "ど"},
	{"ﾅ", "ナ", "な"}, {"ﾆ", "ニ", "に"}, {"ﾇ", "ヌ", "ぬ"}, {"ﾈ", "ネ", "ね"}, {"ﾉ", "ノ", "の"},
	{"ﾊ", "ハ", "は"}, {"ﾋ", "ヒ", "ひ"}, {"ﾌ", "フ", "ふ"}, {"ﾍ", "ヘ", "へ"}, {"ﾎ", "ホ", "ほ"},
	{"ﾊﾞ", "バ", "ば"}, {"ﾋﾞ", "ビ", "び"}, {"ﾌﾞ", "ブ", "ぶ"}, {"ﾍﾞ", "ベ", "べ"}, {"ﾎﾞ", "ボ", "ぼ"},
	{"ﾊﾟ", "パ", "ぱ"}, {"ﾋﾟ", "ピ", "ぴ"}, {"ﾌﾟ", "プ", "ぷ"}, {"ﾍﾟ", "ペ", "ぺ"}, {"ﾎﾟ", "ポ", "ぽ"},
	{"ﾏ", "マ", "ま"}, {"ﾐ", "ミ", "み"}, {"ﾑ", "ム", "む"}, {"ﾒ", "メ", "め"}, {"ﾓ", "モ", "も"},
	{"ﾗ", "ラ", "ら"}, {"ﾘ", "リ", "り"}, {"ﾙ", "ル", "る"}, {"ﾚ", "レ", "れ"}, {"ﾛ", "ロ", "ろ"},
	{"ﾜ", "ワ", "わ"}, {"ｦ", "ヲ", "を"}, {"ﾝ", "ン", "ん"},
}

var ss [][2]string = [][2]string{
	{" ", "　"},
}

var ns [][2]string = [][2]string{
	{"0", "０"}, {"1", "１"}, {"2", "２"}, {"3", "３"}, {"4", "４"},
	{"5", "５"}, {"6", "６"}, {"7", "７"}, {"8", "８"}, {"9", "９"},
}

var al [][2]string = [][2]string{
	{"a", "ａ"}, {"b", "ｂ"}, {"c", "ｃ"}, {"d", "ｄ"}, {"e", "ｅ"},
	{"f", "ｆ"}, {"g", "ｇ"}, {"h", "ｈ"}, {"i", "ｉ"}, {"j", "ｊ"},
	{"k", "ｋ"}, {"l", "ｌ"}, {"m", "ｍ"}, {"n", "ｎ"}, {"o", "ｏ"},
	{"p", "ｐ"}, {"q", "ｑ"}, {"r", "ｒ"}, {"s", "ｓ"}, {"t", "ｔ"},
	{"u", "ｕ"}, {"v", "ｖ"}, {"w", "ｗ"}, {"x", "ｘ"}, {"y", "ｙ"},
	{"z", "ｚ"},
}

var au [][2]string = [][2]string{
	{"A", "Ａ"}, {"B", "Ｂ"}, {"C", "Ｃ"}, {"D", "Ｄ"}, {"E", "Ｅ"},
	{"F", "Ｆ"}, {"G", "Ｇ"}, {"H", "Ｈ"}, {"I", "Ｉ"}, {"J", "Ｊ"},
	{"K", "Ｋ"}, {"L", "Ｌ"}, {"M", "Ｍ"}, {"N", "Ｎ"}, {"O", "Ｏ"},
	{"P", "Ｐ"}, {"Q", "Ｑ"}, {"R", "Ｒ"}, {"S", "Ｓ"}, {"T", "Ｔ"},
	{"U", "Ｕ"}, {"V", "Ｖ"}, {"W", "Ｗ"}, {"X", "Ｘ"}, {"Y", "Ｙ"},
	{"Z", "Ｚ"},
}

var ms [][2]string = [][2]string{
	{"!", "！"}, {"#", "＃"}, {"$", "＄"}, {"%", "％"}, {"&", "＆"},
	{"(", "（"}, {")", "）"}, {"*", "＊"}, {"+", "＋"}, {",", "，"},
	{"-", "－"}, {".", "．"}, {"/", "／"}, {":", "："}, {";", "；"},
	{"<", "＜"}, {"=", "＝"}, {">", "＞"}, {"?", "？"}, {"@", "＠"},
	{"[", "［"}, {"]", "］"}, {"^", "＾"}, {"_", "＿"}, {"`", "｀"},
	{"{", "｛"}, {"|", "｜"}, {"}", "｝"},
}

var mode []string = []string{
	"r", "R", "n", "N", "a", "A", "s", "S", "k", "K", "h", "H", "c", "C",
}

func createInput() []string {
	f, err := os.Open("./data/input.txt")
	if err != nil {
		return []string{}
	}
	tmp, err := ioutil.ReadAll(f)
	if err != nil {
		return []string{}
	}
	return toStringSlice(tmp)
}

func createExpected(str []string, mode string) []byte {
	path := fmt.Sprintf("./data/output.%s.txt", mode)
	f, err := os.Open(path)
	if err != nil {
		return []byte{}
	}
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return []byte{}
	}
	return buf
}

func execConvTest(t *testing.T, f func(w *word) []byte, cases []testCase) {
	for _, cs := range cases {
		r := f(cs.arg)
		if !bytes.Equal(r, []byte(cs.expected)) {
			t.Errorf("Fail to convert. %s != %s\n", string(r), cs.expected)
		}
	}
}

func toStringSlice(b []byte) []string {
	buf := make([]string, 0, len(b))
	for _, v := range strings.Split(string(b), "") {
		k := len(buf)
		if k > 0 && (v == "ﾞ" || v == "ﾟ") {
			buf[k-1] += v
			continue
		}
		buf = append(buf, v)
	}
	return buf
}

func TestByte(t *testing.T) {
	str := createInput()
	for _, m := range mode {
		ret := Byte([]byte(strings.Join(str, "")), m)
		inputs := toStringSlice(ret)
		outputs := toStringSlice(createExpected(str, m))
		for k, v := range inputs {
			if !bytes.Equal([]byte(v), []byte(outputs[k])) {
				t.Errorf("[Byte] Fail to convert. (%s mode)\nExpected: %s\nReturned: %s\n", m, outputs[k], v)
			}
		}
		break
	}
}

func TestString(t *testing.T) {
	str := createInput()
	for _, m := range mode {
		ret := String(strings.Join(str, ""), m)
		inputs := toStringSlice([]byte(ret))
		outputs := toStringSlice(createExpected(str, m))
		for k, v := range inputs {
			if !bytes.Equal([]byte(v), []byte(outputs[k])) {
				t.Errorf("[String] Fail to convert. (%s mode)\nExpected: %s\nReturned: %s\n", m, outputs[k], v)
			}
		}
		break
	}
}

func TestReaderRead(t *testing.T) {
	// for _, m := range mode {
	// 	f, err := os.Open("./data/input.txt")
	// 	if err != nil {
	// 		t.Errorf("[NewReader] Fail to open test file.")
	// 	}
	// 	r := NewReader(f, m)

	// 	buf := make([]byte, 4096)
	// 	for {
	// 		line := make([]byte, 4096)
	// 		_, err := r.Read(line)
	// 		if err == io.EOF {
	// 			break
	// 		}
	// 		if err != nil {
	// 			break
	// 		}
	// 		buf = append(buf, line...)
	// 	}
	// 	f.Close()
	// 	expected := createExpected(strings.Split(string(buf), ""), m)
	// 	for k, s := range strings.Split(string(buf), "") {
	// 		if !bytes.Equal([]byte(expected[k]), []byte(s)) {
	// 			t.Errorf("[NewReader] Expected:%s Returned:%s\n", expected[k], s)
	// 		}
	// 	}
	// 	break
	// }
}

func TestNewWriter(t *testing.T) {

}

/**
 * Zenkaku Space -> Hankaku Space
 */
func TestConvAsSmallS(t *testing.T) {
	w := &word{[]byte(ss[0][1]), zenkaku + space, len(ss[0][1])}
	e := ss[0][0]
	c := testCase{w, e}
	cases := []testCase{c}
	execConvTest(t, convAsSmallS, cases)
}

/**
 * Hankaku Space -> Zenkaku Space
 */
func TestConvAsLargeS(t *testing.T) {
	w := &word{[]byte(ss[0][0]), hankaku + space, len(ss[0][0])}
	e := ss[0][1]
	c := testCase{w, e}
	cases := []testCase{c}
	execConvTest(t, convAsLargeS, cases)
}

/**
 * Zenkaku Numeric -> Hankaku Numeric
 */
func TestConvAsSmallN(t *testing.T) {
	cases := []testCase{}
	for _, n := range ns {
		w := &word{[]byte(n[1]), zenkaku + numeric, len(n[1])}
		e := n[0]
		c := testCase{w, e}
		cases = append(cases, c)
	}
	execConvTest(t, convAsSmallN, cases)
}

/**
 * Hankaku Numeric -> Zenkaku Numeric
 */
func TestConvAsLargeN(t *testing.T) {
	cases := []testCase{}
	for _, n := range ns {
		w := &word{[]byte(n[0]), hankaku + numeric, len(n[0])}
		e := n[1]
		c := testCase{w, e}
		cases = append(cases, c)
	}
	execConvTest(t, convAsLargeN, cases)
}

/**
 * Zenkaku Alphabet -> Hankaku Alphabet
 */
func TestConvAsSmallR(t *testing.T) {
	cases := []testCase{}
	for _, n := range al {
		w := &word{[]byte(n[1]), zenkaku + alphabet + lowercase, len(n[1])}
		e := n[0]
		c := testCase{w, e}
		cases = append(cases, c)
	}
	for _, n := range au {
		w := &word{[]byte(n[1]), zenkaku + alphabet + uppercase, len(n[1])}
		e := n[0]
		c := testCase{w, e}
		cases = append(cases, c)
	}
	execConvTest(t, convAsSmallR, cases)
}

/**
 * Hankaku Alphabet -> Zenkaku Alphabet
 */
func TestConvAsLargeR(t *testing.T) {
	cases := []testCase{}
	for _, n := range al {
		w := &word{[]byte(n[0]), hankaku + alphabet + lowercase, len(n[0])}
		e := n[1]
		c := testCase{w, e}
		cases = append(cases, c)
	}
	for _, n := range au {
		w := &word{[]byte(n[0]), hankaku + alphabet + uppercase, len(n[0])}
		e := n[1]
		c := testCase{w, e}
		cases = append(cases, c)
	}
	execConvTest(t, convAsLargeR, cases)
}

/**
 * Zenkaku AlphaNumeric -> Hankaku AlphaNumeric
 * !-}(Excluding ",',\)
 */
func TestConvAsSmallA(t *testing.T) {
	cases := []testCase{}
	for _, n := range ns {
		w := &word{[]byte(n[1]), zenkaku + alphanumeric + numeric, len(n[1])}
		e := n[0]
		c := testCase{w, e}
		cases = append(cases, c)
	}
	for _, n := range al {
		w := &word{[]byte(n[1]), zenkaku + alphanumeric + alphabet + lowercase, len(n[1])}
		e := n[0]
		c := testCase{w, e}
		cases = append(cases, c)
	}
	for _, n := range au {
		w := &word{[]byte(n[1]), zenkaku + alphanumeric + alphabet + uppercase, len(n[1])}
		e := n[0]
		c := testCase{w, e}
		cases = append(cases, c)
	}
	for _, n := range ms {
		w := &word{[]byte(n[1]), zenkaku + alphanumeric, len(n[1])}
		e := n[0]
		c := testCase{w, e}
		cases = append(cases, c)
	}
	execConvTest(t, convAsSmallA, cases)
}

/**
 * Hankaku AlphaNumeric -> Zenkaku AlphaNumeric
 * !-}(Excluding ",',\)
 */
func TestConvAsLargeA(t *testing.T) {
	cases := []testCase{}
	for _, n := range ns {
		w := &word{[]byte(n[0]), hankaku + alphanumeric + numeric, len(n[0])}
		e := n[1]
		c := testCase{w, e}
		cases = append(cases, c)
	}
	for _, n := range al {
		w := &word{[]byte(n[0]), hankaku + alphanumeric + alphabet + lowercase, len(n[0])}
		e := n[1]
		c := testCase{w, e}
		cases = append(cases, c)
	}
	for _, n := range au {
		w := &word{[]byte(n[0]), hankaku + alphanumeric + alphabet + uppercase, len(n[0])}
		e := n[1]
		c := testCase{w, e}
		cases = append(cases, c)
	}
	for _, n := range ms {
		w := &word{[]byte(n[0]), hankaku + alphanumeric, len(n[0])}
		e := n[1]
		c := testCase{w, e}
		cases = append(cases, c)
	}
	execConvTest(t, convAsLargeA, cases)
}

/**
 * Zenkaku Katakana -> Hankaku Katakana
 */
func TestConvAsSmallK(t *testing.T) {
	cases := []testCase{}
	for _, n := range ks {
		w := &word{[]byte(n[1]), zenkaku + katakana, len(n[1])}
		e := n[0]
		c := testCase{w, e}
		cases = append(cases, c)
	}
	execConvTest(t, convAsSmallK, cases)
}

/**
 * Hankaku Katakana -> Zenkaku Katakana
 */
func TestConvAsLargeK(t *testing.T) {
	cases := []testCase{}
	for _, n := range ks {
		w := &word{[]byte(n[0]), hankaku + katakana, len(n[0])}
		e := n[1]
		c := testCase{w, e}
		cases = append(cases, c)
	}
	execConvTest(t, convAsLargeK, cases)
}

/**
 * Zenkaku Hiragana -> Hankaku Katakana
 */
func TestConvAsSmallH(t *testing.T) {
	cases := []testCase{}
	for _, n := range ks {
		w := &word{[]byte(n[2]), zenkaku + hiragana, len(n[2])}
		e := n[0]
		c := testCase{w, e}
		cases = append(cases, c)
	}
	execConvTest(t, convAsSmallH, cases)
}

/**
 * Hankaku Katakana -> Zenkaku Hiragana
 */
func TestConvAsLargeH(t *testing.T) {
	cases := []testCase{}
	for _, n := range ks {
		w := &word{[]byte(n[0]), hankaku + katakana, len(n[0])}
		e := n[2]
		c := testCase{w, e}
		cases = append(cases, c)
	}
	execConvTest(t, convAsLargeH, cases)
}

/**
 * Zenkaku Katakana -> Zenkaku Hiragana
 */
func TestConvAsSmallC(t *testing.T) {
	cases := []testCase{}
	for _, n := range ks {
		w := &word{[]byte(n[1]), zenkaku + katakana, len(n[1])}
		e := n[2]
		c := testCase{w, e}
		cases = append(cases, c)
	}
	execConvTest(t, convAsSmallC, cases)
}

/**
 * Zenkaku Hiragana -> Zenkaku Katakana
 */
func TestConvAsLargeC(t *testing.T) {
	cases := []testCase{}
	for _, n := range ks {
		w := &word{[]byte(n[2]), zenkaku + hiragana, len(n[2])}
		e := n[1]
		c := testCase{w, e}
		cases = append(cases, c)
	}
	execConvTest(t, convAsLargeC, cases)
}

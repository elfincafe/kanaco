package kanaco

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

type testCase struct {
	arg *word
	ret string
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
	{"ﾜ", "ワ", "わ"}, {"ｦ", "ヲ", "を"}, {"ﾝ", "ン", "ん"}, {"ｰ", "ー", "ー"},
}

var ss [][2]string = [][2]string{
	{" ", "　"},
}

var ns [][2]string = [][2]string{
	{"0", "０"}, {"1", "１"}, {"2", "２"}, {"3", "３"}, {"4", "４"},
	{"5", "５"}, {"6", "６"}, {"7", "７"}, {"8", "８"}, {"9", "９"},
}

var as [][2]string = [][2]string{
	{"a", "ａ"}, {"b", "ｂ"}, {"c", "ｃ"}, {"d", "ｄ"}, {"e", "ｅ"},
	{"f", "ｆ"}, {"g", "ｇ"}, {"h", "ｈ"}, {"i", "ｉ"}, {"j", "ｊ"},
	{"k", "ｋ"}, {"l", "ｌ"}, {"m", "ｍ"}, {"n", "ｎ"}, {"o", "ｏ"},
	{"p", "ｐ"}, {"q", "ｑ"}, {"r", "ｒ"}, {"s", "ｓ"}, {"t", "ｔ"},
	{"u", "ｕ"}, {"v", "ｖ"}, {"w", "ｗ"}, {"x", "ｘ"}, {"y", "ｙ"},
	{"z", "ｚ"},
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
	cases := []testCase{}
	for _, v := range ss {
		w := new(word)
		w.val = []byte(v[1])
		w.len = len(v[1])
		w.charType = zenkaku + space
		tc := testCase{arg: w, ret: v[0]}
		cases = append(cases, tc)
	}
	for _, c := range cases {
		r := convAsSmallS(c.arg)
		if !bytes.Equal(r, []byte(c.ret)) {
			t.Errorf("[convAsSmallS] Fail to convert. %s != %s\n", c.arg.val, string(r))
		}
	}
}

/**
 * Hankaku Space -> Zenkaku Space
 */
func TestConvAsLargeS(t *testing.T) {
	cases := []testCase{}
	for _, v := range ss {
		w := new(word)
		w.val = []byte(v[0])
		w.len = len(v[0])
		w.charType = zenkaku + space
		tc := testCase{arg: w, ret: v[1]}
		cases = append(cases, tc)
	}
	for _, c := range cases {
		r := convAsLargeS(c.arg)
		if !bytes.Equal(r, []byte(c.ret)) {
			t.Errorf("[convAsLargeS] Fail to convert. %s != %s\n", c.arg.val, string(r))
		}
	}
}

/**
 * Zenkaku Numeric -> Hankaku Numeric
 */
func TestConvAsSmallN(t *testing.T) {
	cases := []testCase{}
	for _, v := range ns {
		w := new(word)
		w.val = []byte(v[1])
		w.len = len(v[1])
		w.charType = zenkaku + numeric
		tc := testCase{arg: w, ret: v[0]}
		cases = append(cases, tc)
	}
	for _, c := range cases {
		r := convAsSmallN(c.arg)
		if !bytes.Equal(r, []byte(c.ret)) {
			t.Errorf("[convAsSmallN] Fail to convert. %s != %s\n", c.arg.val, string(r))
		}
	}
}

/**
 * Hankaku Numeric -> Zenkaku Numeric
 */
func TestConvAsLargeN(t *testing.T) {
	cases := []testCase{}
	for _, v := range ns {
		w := new(word)
		w.val = []byte(v[0])
		w.len = len(v[0])
		w.charType = hankaku + numeric
		tc := testCase{arg: w, ret: v[1]}
		cases = append(cases, tc)
	}
	for _, c := range cases {
		r := convAsLargeN(c.arg)
		if !bytes.Equal(r, []byte(c.ret)) {
			t.Errorf("[convAsLargeN] Fail to convert. %s != %s\n", c.arg.val, string(r))
		}
	}
}

/**
 * Zenkaku Alphabet -> Hankaku Alphabet
 */
func TestConvAsSmallR(t *testing.T) {
	cases := []testCase{}
	for _, v := range as {
		w := new(word)
		w.val = []byte(v[1])
		w.len = len(v[1])
		w.charType = zenkaku + alphabet
		tc := testCase{arg: w, ret: v[0]}
		cases = append(cases, tc)
	}
	for _, c := range cases {
		r := convAsSmallR(c.arg)
		if !bytes.Equal(r, []byte(c.ret)) {
			t.Errorf("[convAsSmallR] Fail to convert. %s != %s\n", c.arg.val, string(r))
		}
	}
}

/**
 * Hankaku Alphabet -> Zenkaku Alphabet
 */
func TestConvAsLargeR(t *testing.T) {
	cases := []testCase{}
	for _, v := range as {
		w := new(word)
		w.val = []byte(v[0])
		w.len = len(v[0])
		w.charType = hankaku + alphabet
		tc := testCase{arg: w, ret: v[1]}
		cases = append(cases, tc)
	}
	for _, c := range cases {
		r := convAsLargeR(c.arg)
		if !bytes.Equal(r, []byte(c.ret)) {
			t.Errorf("[convAsLargeR] Fail to convert. %s != %s\n", c.arg.val, string(r))
		}
	}
}

/**
 * Zenkaku AlphaNumeric -> Hankaku AlphaNumeric
 * !-}(Excluding ",',\)
 */
func TestConvAsSmallA(t *testing.T) {
	cases := []testCase{}
	as = append(as, ns...)
	as = append(as, ms...)
	for _, v := range as {
		w := new(word)
		w.val = []byte(v[1])
		w.len = len(v[1])
		w.charType = zenkaku + alphabet + numeric
		tc := testCase{arg: w, ret: v[0]}
		cases = append(cases, tc)
	}
	cases = append(cases, testCase{arg: &word{[]byte("”"), asIs, 3}, ret: "”"})
	cases = append(cases, testCase{arg: &word{[]byte("’"), asIs, 3}, ret: "’"})
	cases = append(cases, testCase{arg: &word{[]byte("＼"), asIs, 3}, ret: "＼"})
	cases = append(cases, testCase{arg: &word{[]byte("～"), asIs, 3}, ret: "～"})
	for _, c := range cases {
		r := convAsSmallA(c.arg)
		if !bytes.Equal(r, []byte(c.ret)) {
			t.Errorf("[convAsSmallA] Fail to convert. %s != %s %v\n", c.arg.val, string(r), r)
		}
	}
}

/**
 * Hankaku AlphaNumeric -> Zenkaku AlphaNumeric
 * !-}(Excluding ",',\)
 */
func TestConvAsLargeA(t *testing.T) {
	cases := []testCase{}
	as = append(as, ns...)
	as = append(as, ms...)
	cases = append(cases, testCase{arg: &word{[]byte("\""), asIs, 1}, ret: "\""})
	cases = append(cases, testCase{arg: &word{[]byte("'"), asIs, 1}, ret: "'"})
	cases = append(cases, testCase{arg: &word{[]byte("\\"), asIs, 1}, ret: "\\"})
	cases = append(cases, testCase{arg: &word{[]byte("~"), asIs, 1}, ret: "~"})
	for _, v := range as {
		w := new(word)
		w.val = []byte(v[0])
		w.len = len(v[0])
		w.charType = hankaku + alphabet + numeric
		tc := testCase{arg: w, ret: v[1]}
		cases = append(cases, tc)
	}
	for _, c := range cases {
		r := convAsLargeA(c.arg)
		if !bytes.Equal(r, []byte(c.ret)) {
			t.Errorf("[convAsLargeA] Fail to convert. %s != %s\n", c.arg.val, string(r))
		}
	}
}

/**
 * Zenkaku Katakana -> Hankaku Katakana
 */
func TestConvAsSmallK(t *testing.T) {
	cases := []testCase{}
	for _, v := range ks {
		w := new(word)
		w.val = []byte(v[1])
		w.len = len(v[1])
		w.charType = zenkaku + katakana
		tc := testCase{arg: w, ret: v[0]}
		cases = append(cases, tc)
	}
	for _, c := range cases {
		r := convAsSmallK(c.arg)
		if !bytes.Equal(r, []byte(c.ret)) {
			t.Errorf("[convAsSmallK] Fail to convert. %s != %s\n", c.arg.val, string(r))
		}
	}
}

/**
 * Hankaku Katakana -> Zenkaku Katakana
 */
func TestConvAsLargeK(t *testing.T) {
	cases := []testCase{}
	for _, v := range ks {
		w := new(word)
		w.val = []byte(v[0])
		w.len = len(v[0])
		if w.len == 3 {
			w.charType = hankaku + alphabet
		} else if w.len == 6 {
			if w.val[5] == 0x9e {
				w.charType = hankaku + katakana + voiced
			} else if w.val[5] == 0x9f {
				w.charType = hankaku + katakana + devoiced
			}
		}
		tc := testCase{arg: w, ret: v[1]}
		cases = append(cases, tc)
	}
	for _, c := range cases {
		r := convAsLargeK(c.arg)
		if !bytes.Equal(r, []byte(c.ret)) {
			t.Errorf("[convAsLargeK] Fail to convert. %s != %s\n", c.arg.val, string(r))
		}
	}
}

/**
 * Zenkaku Hiragana -> Hankaku Katakana
 */
func TestConvAsSmallH(t *testing.T) {
	cases := []testCase{}
	for _, v := range ks {
		w := new(word)
		w.val = []byte(v[1])
		w.len = len(v[1])
		w.charType = zenkaku + hiragana
		tc := testCase{arg: w, ret: v[0]}
		cases = append(cases, tc)
	}
	for _, c := range cases {
		r := convAsSmallH(c.arg)
		if !bytes.Equal(r, []byte(c.ret)) {
			t.Errorf("[convAsSmallH] Fail to convert. %s != %s\n", c.ret, string(r))
		}
	}
}

/**
 * Hankaku Katakana -> Zenkaku Hiragana
 */
func TestConvAsLargeH(t *testing.T) {
	cases := []testCase{}
	for _, v := range ks {
		w := new(word)
		w.val = []byte(v[0])
		w.len = len(v[0])
		if w.len == 3 {
			w.charType = hankaku + alphabet
		} else if w.len == 6 {
			if w.val[5] == 0x9e {
				w.charType = hankaku + katakana + voiced
			} else if w.val[5] == 0x9f {
				w.charType = hankaku + katakana + devoiced
			}
		}
		tc := testCase{arg: w, ret: v[2]}
		cases = append(cases, tc)
	}
	for _, c := range cases {
		r := convAsLargeH(c.arg)
		if !bytes.Equal(r, []byte(c.ret)) {
			t.Errorf("[convAsLargeH] Fail to convert. %s != %s\n", c.ret, string(r))
		}
	}
}

/**
 * Zenkaku Katakana -> Zenkaku Hiragana
 */
func TestConvAsSmallC(t *testing.T) {
	cases := []testCase{}
	for _, v := range ks {
		w := new(word)
		w.val = []byte(v[1])
		w.len = len(v[1])
		w.charType = zenkaku + katakana
		tc := testCase{arg: w, ret: v[2]}
		cases = append(cases, tc)
	}
	for _, c := range cases {
		r := convAsSmallC(c.arg)
		if !bytes.Equal(r, []byte(c.ret)) {
			t.Errorf("[convAsSmallC] Fail to convert. %s != %s\n", c.ret, string(r))
		}
	}
}

/**
 * Zenkaku Hiragana -> Zenkaku Katakana
 */
func TestConvAsLargeC(t *testing.T) {
	cases := []testCase{}
	for _, v := range ks {
		w := new(word)
		w.val = []byte(v[2])
		w.len = len(v[2])
		w.charType = zenkaku + hiragana
		tc := testCase{arg: w, ret: v[1]}
		cases = append(cases, tc)
	}
	for _, c := range cases {
		r := convAsLargeC(c.arg)
		if !bytes.Equal(r, []byte(c.ret)) {
			t.Errorf("[convAsLargeC] Fail to convert. %s != %s\n", c.ret, string(r))
		}
	}
}

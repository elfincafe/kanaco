package kanaco

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

type testCase struct {
	arg string
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
//	"r",
}

func createInput() string {
	s := strings.Builder{}
	for _, v := range ks {
		s.WriteString(v[0])
		s.WriteString(v[1])
		s.WriteString(v[2])
	}
	for _, v1 := range [][][2]string{ss, ns, as, ms} {
		for _, v2 := range v1 {
			s.WriteString(v2[0])
			s.WriteString(v2[1])
		}
	}
	return s.String()
}

func createExpected(s string, mode string) string {
	for _, m := range mode {
		switch m {
		case 'r':
			for _, v := range as {
				s = strings.ReplaceAll(s, v[1], v[0])
			}
		case 'R':
			for _, v := range as {
				s = strings.ReplaceAll(s, v[0], v[1])
			}
		case 'n':
			for _, v := range ns {
				s = strings.ReplaceAll(s, v[1], v[0])
			}
		case 'N':
			for _, v := range ns {
				s = strings.ReplaceAll(s, v[0], v[1])
			}
		case 'a':
			for _, v1 := range [][][2]string{ms,ns,as} {
				for _,v2 := range v1 {
					s = strings.ReplaceAll(s, v2[1], v2[0])
				}
			}
		case 'A':
			for _, v1 := range [][][2]string{ms,ns,as} {
				for _,v2 := range v1 {
					s = strings.ReplaceAll(s, v2[0], v2[1])
				}
			}
		case 's':
			for _, v := range ss {
				s = strings.ReplaceAll(s, v[1], v[0])
			}
		case 'S':
			for _, v := range ss {
				s = strings.ReplaceAll(s, v[0], v[1])
			}
		case 'k':
			cnt := len(ks) - 1
			for i:=cnt; i>=0; i-- {
				v := ks[i]
				s = strings.ReplaceAll(s, v[1], v[0])
			}
		case 'K':
			cnt := len(ks) - 1
			for i:=cnt; i>=0; i-- {
				v := ks[i]
				s = strings.ReplaceAll(s, v[0], v[1])
			}
		case 'h':
			cnt := len(ks) - 1
			for i:=cnt; i>=0; i-- {
				v := ks[i]
				s = strings.ReplaceAll(s, v[2], v[0])
			}
		case 'H':
			cnt := len(ks) - 1
			for i:=cnt; i>=0; i-- {
				v := ks[i]
				s = strings.ReplaceAll(s, v[0], v[2])
			}
		case 'c':
			cnt := len(ks) - 1
			for i:=cnt; i>=0; i-- {
				v := ks[i]
				s = strings.ReplaceAll(s, v[1], v[2])
			}
		case 'C':
			cnt := len(ks) - 1
			for i:=cnt; i>=0; i-- {
				v := ks[i]
				s = strings.ReplaceAll(s, v[2], v[1])
			}
		}
	}
	return s
}

func TestByte(t *testing.T) {
	s := createInput()
	for _, m := range mode {
		ret := Byte([]byte(s), m)
		expected := []rune(createExpected(s, m))
		for k, v := range []rune(string(ret)) {
			if expected[k]!=v {
				t.Errorf("[Byte] Fail to convert. (%s mode)\nExpected: %c\nReturned: %c\n", m, expected[k], v)
			}
		}
	}
}

func TestString(t *testing.T) {
	s := createInput()
	for _, m := range mode {
		ret := String(s, m)
		expected := []rune(createExpected(s, m))
		for k, v := range []rune(string(ret)) {
			if expected[k]!=v {
				t.Errorf("[String] Fail to convert. (%s mode)\nExpected: %c\nReturned: %c\n", m, expected[k], v)
			}
		}
	}
}

func TestReaderRead (t *testing.T) {
	dir := os.TempDir()
	f, err := os.CreateTemp(dir, "TestNewReader.*.txt")
	if err != nil {
		t.Errorf("[NewReader] Fail to create test file.")
	}
	s := createInput()
	s = strings.ReplaceAll(s, "ァ", "\nァ")
	s = strings.ReplaceAll(s, "あ", "\nあ")
	s = strings.ReplaceAll(s, "z", "\nz")
	s = strings.ReplaceAll(s, "5", "\n5")
	f.WriteString(s)
	path := f.Name()
	f.Close()
	defer os.Remove(path)

	for _, m := range mode {
		f, err = os.Open(path)
		if err != nil {
			t.Errorf("[NewReader] Fail to open test file.")
		}
		r := NewReader(f, m)
		expected := createExpected(s, m)
		expected = strings.ReplaceAll(expected, "ァ", "\nァ")
		expected = strings.ReplaceAll(expected, "あ", "\nあ")
		expected = strings.ReplaceAll(expected, "z", "\nz")
		expected = strings.ReplaceAll(expected, "5", "\n5")
		buf := make([]byte, 0, len(expected)+1024)
		for {
			b := make([]byte, 4096)
			_, err := r.Read(b)
			if err == io.EOF {
				buf = append(buf, b...)
				break
			}
			if err != nil {
				break
			}
			buf = append(buf, b...)
		}
		f.Close()
		if bytes.Equal([]byte(expected), buf) {
			runeE := []rune(expected)
			runeB := []rune(string(buf))
			for k, v := range runeE {
				v++
				if v!= runeB[k] {
					t.Errorf("[NewReader] Expected:%c Returned:%c\n", v, runeB[k])
				}
			}
		}
	}
}

func TestNewWriter (t *testing.T) {

}

/**
 * Zenkaku Space -> Hankaku Space
 */
func TestConvAsSmallS(t *testing.T) {
	cases := []testCase{}
	for _, v := range ss {
		tc := testCase{arg: v[1], ret: v[0]}
		cases = append(cases, tc)
	}
	for _, c := range cases {
		r := convAsSmallS([]byte(c.arg))
		if !bytes.Equal(r, []byte(c.ret)) {
			t.Errorf("[convAsSmallS] Fail to convert. %s != %s\n", c.arg, string(r))
		}
	}
}

/**
 * Hankaku Space -> Zenkaku Space
 */
func TestConvAsLargeS(t *testing.T) {
	cases := []testCase{}
	for _, v := range ss {
		tc := testCase{arg: v[0], ret: v[1]}
		cases = append(cases, tc)
	}
	for _, c := range cases {
		r := convAsLargeS([]byte(c.arg))
		if !bytes.Equal(r, []byte(c.ret)) {
			t.Errorf("[convAsLargeS] Fail to convert. %s != %s\n", c.arg, string(r))
		}
	}
}

/**
 * Zenkaku Numeric -> Hankaku Numeric
 */
func TestConvAsSmallN(t *testing.T) {
	cases := []testCase{}
	for _, v := range ns {
		tc := testCase{arg: v[1], ret: v[0]}
		cases = append(cases, tc)
	}
	for _, c := range cases {
		r := convAsSmallN([]byte(c.arg))
		if !bytes.Equal(r, []byte(c.ret)) {
			t.Errorf("[convAsSmallN] Fail to convert. %s != %s\n", c.arg, string(r))
		}
	}
}

/**
 * Hankaku Numeric -> Zenkaku Numeric
 */
func TestConvAsLargeN(t *testing.T) {
	cases := []testCase{}
	for _, v := range ns {
		tc := testCase{arg: v[0], ret: v[1]}
		cases = append(cases, tc)
	}
	for _, c := range cases {
		r := convAsLargeN([]byte(c.arg))
		if !bytes.Equal(r, []byte(c.ret)) {
			t.Errorf("[convAsLargeN] Fail to convert. %s != %s\n", c.arg, string(r))
		}
	}
}

/**
 * Zenkaku Alphabet -> Hankaku Alphabet
 */
func TestConvAsSmallR(t *testing.T) {
	cases := []testCase{}
	for _, v := range as {
		tc := testCase{arg: v[1], ret: v[0]}
		cases = append(cases, tc)
	}
	for _, c := range cases {
		r := convAsSmallR([]byte(c.arg))
		if !bytes.Equal(r, []byte(c.ret)) {
			t.Errorf("[convAsSmallR] Fail to convert. %s != %s\n", c.arg, string(r))
		}
	}
}

/**
 * Hankaku Alphabet -> Zenkaku Alphabet
 */
func TestConvAsLargeR(t *testing.T) {
	cases := []testCase{}
	for _, v := range as {
		tc := testCase{arg: v[0], ret: v[1]}
		cases = append(cases, tc)
	}
	for _, c := range cases {
		r := convAsLargeR([]byte(c.arg))
		if !bytes.Equal(r, []byte(c.ret)) {
			t.Errorf("[convAsLargeR] Fail to convert. %s != %s\n", c.arg, string(r))
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
		tc := testCase{arg: v[1], ret: v[0]}
		cases = append(cases, tc)
	}
	cases = append(cases, testCase{arg: "”", ret: "”"})
	cases = append(cases, testCase{arg: "’", ret: "’"})
	cases = append(cases, testCase{arg: "～", ret: "～"})
	for _, c := range cases {
		r := convAsSmallA([]byte(c.arg))
		if !bytes.Equal(r, []byte(c.ret)) {
			t.Errorf("[convAsSmallA] Fail to convert. %s != %s %v\n", c.arg, string(r), r)
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
	cases = append(cases, testCase{arg: "\"", ret: "\""})
	cases = append(cases, testCase{arg: "'", ret: "'"})
	cases = append(cases, testCase{arg: "~", ret: "~"})
	for _, v := range as {
		tc := testCase{arg: v[0], ret: v[1]}
		cases = append(cases, tc)
	}
	for _, c := range cases {
		r := convAsLargeA([]byte(c.arg))
		if !bytes.Equal(r, []byte(c.ret)) {
			t.Errorf("[convAsLargeA] Fail to convert. %s != %s\n", c.arg, string(r))
		}
	}
}

/**
 * Zenkaku Katakana -> Hankaku Katakana
 */
func TestConvAsSmallK(t *testing.T) {
	cases := []testCase{}
	for _, v := range ks {
		tc := testCase{arg: v[1], ret: v[0]}
		cases = append(cases, tc)
	}
	for _, c := range cases {
		r := convAsSmallK([]byte(c.arg))
		if !bytes.Equal(r, []byte(c.ret)) {
			t.Errorf("[convAsSmallK] Fail to convert. %s != %s\n", c.arg, string(r))
		}
	}
}

/**
 * Hankaku Katakana -> Zenkaku Katakana
 */
func TestConvAsLargeK(t *testing.T) {
	cases := []testCase{}
	for _, v := range ks {
		tc := testCase{arg: v[0], ret: v[1]}
		cases = append(cases, tc)
	}
	for _, c := range cases {
		r := convAsLargeK([]byte(c.arg))
		if !bytes.Equal(r, []byte(c.ret)) {
			t.Errorf("[convAsLargeK] Fail to convert. %s != %s\n", c.arg, string(r))
		}
	}
}

/**
 * Zenkaku Hiragana -> Hankaku Katakana
 */
func TestConvAsSmallH(t *testing.T) {
	cases := []testCase{}
	for _, v := range ks {
		tc := testCase{arg: v[2], ret: v[0]}
		cases = append(cases, tc)
	}
	for _, c := range cases {
		r := convAsSmallH([]byte(c.arg))
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
		tc := testCase{arg: v[0], ret: v[2]}
		cases = append(cases, tc)
	}
	for _, c := range cases {
		r := convAsLargeH([]byte(c.arg))
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
		tc := testCase{arg: v[1], ret: v[2]}
		cases = append(cases, tc)
	}
	for _, c := range cases {
		r := convAsSmallC([]byte(c.arg))
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
		tc := testCase{arg: v[2], ret: v[1]}
		cases = append(cases, tc)
	}
	for _, c := range cases {
		r := convAsLargeC([]byte(c.arg))
		if !bytes.Equal(r, []byte(c.ret)) {
			t.Errorf("[convAsLargeC] Fail to convert. %s != %s\n", c.ret, string(r))
		}
	}
}

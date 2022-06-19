package kanaco

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

type Reader struct {
	r    *bufio.Reader
	mode string
}

type char struct {
	val  []byte
	mode []byte
	len  int
}

var (
	tbl map[string]map[byte][]byte = map[string]map[byte][]byte{
		"ぁ": {'h': []byte("ｧ")},
		"あ": {'h': []byte("ｱ")},
		"ぃ": {'h': []byte("ｨ")},
		"い": {'h': []byte("ｲ")},
		"ぅ": {'h': []byte("ｩ")},
		"う": {'h': []byte("ｳ")},
		"ぇ": {'h': []byte("ｪ")},
		"え": {'h': []byte("ｴ")},
		"ぉ": {'h': []byte("ｫ")},
		"お": {'h': []byte("ｵ")},
		"か": {'h': []byte("ｶ")},
		"が": {'h': []byte("ｶﾞ")},
		"き": {'h': []byte("ｷ")},
		"ぎ": {'h': []byte("ｷﾞ")},
		"く": {'h': []byte("ｸ")},
		"ぐ": {'h': []byte("ｸﾞ")},
		"け": {'h': []byte("ｹ")},
		"げ": {'h': []byte("ｹﾞ")},
		"こ": {'h': []byte("ｺ")},
		"ご": {'h': []byte("ｺﾞ")},
		"さ": {'h': []byte("ｻ")},
		"ざ": {'h': []byte("ｻﾞ")},
		"し": {'h': []byte("ｼ")},
		"じ": {'h': []byte("ｼﾞ")},
		"す": {'h': []byte("ｽ")},
		"ず": {'h': []byte("ｽﾞ")},
		"せ": {'h': []byte("ｾ")},
		"ぜ": {'h': []byte("ｾﾞ")},
		"そ": {'h': []byte("ｿ")},
		"ぞ": {'h': []byte("ｿﾞ")},
		"た": {'h': []byte("ﾀ")},
		"だ": {'h': []byte("ﾀﾞ")},
		"ち": {'h': []byte("ﾁ")},
		"ぢ": {'h': []byte("ﾁﾞ")},
		"っ": {'h': []byte("ｯ")},
		"つ": {'h': []byte("ﾂ")},
		"づ": {'h': []byte("ﾂﾞ")},
		"て": {'h': []byte("ﾃ")},
		"で": {'h': []byte("ﾃﾞ")},
		"と": {'h': []byte("ﾄ")},
		"ど": {'h': []byte("ﾄﾞ")},
		"な": {'h': []byte("ﾅ")},
		"に": {'h': []byte("ﾆ")},
		"ぬ": {'h': []byte("ﾇ")},
		"ね": {'h': []byte("ﾈ")},
		"の": {'h': []byte("ﾉ")},
		"は": {'h': []byte("ﾊ")},
		"ば": {'h': []byte("ﾊﾞ")},
		"ぱ": {'h': []byte("ﾊﾟ")},
		"ひ": {'h': []byte("ﾋ")},
		"び": {'h': []byte("ﾋﾞ")},
		"ぴ": {'h': []byte("ﾋﾟ")},
		"ふ": {'h': []byte("ﾌ")},
		"ぶ": {'h': []byte("ﾌﾞ")},
		"ぷ": {'h': []byte("ﾌﾟ")},
		"へ": {'h': []byte("ﾍ")},
		"べ": {'h': []byte("ﾍﾞ")},
		"ぺ": {'h': []byte("ﾍﾟ")},
		"ほ": {'h': []byte("ﾎ")},
		"ぼ": {'h': []byte("ﾎﾞ")},
		"ぽ": {'h': []byte("ﾎﾟ")},
		"ま": {'h': []byte("ﾏ")},
		"み": {'h': []byte("ﾐ")},
		"む": {'h': []byte("ﾑ")},
		"め": {'h': []byte("ﾒ")},
		"も": {'h': []byte("ﾓ")},
		"ゃ": {'h': []byte("ｬ")},
		"や": {'h': []byte("ﾔ")},
		"ゅ": {'h': []byte("ｭ")},
		"ゆ": {'h': []byte("ﾕ")},
		"ょ": {'h': []byte("ｮ")},
		"よ": {'h': []byte("ﾖ")},
		"ら": {'h': []byte("ﾗ")},
		"り": {'h': []byte("ﾘ")},
		"る": {'h': []byte("ﾙ")},
		"れ": {'h': []byte("ﾚ")},
		"ろ": {'h': []byte("ﾛ")},
		"ゎ": {'h': []byte("ﾜ")},
		"わ": {'h': []byte("ﾜ")},
		"ゐ": {'h': []byte("ｲ")},
		"ゑ": {'h': []byte("ｴ")},
		"を": {'h': []byte("ｦ")},
		"ん": {'h': []byte("ﾝ")},

		"ァ": {'k': []byte("ｧ")},
		"ア": {'k': []byte("ｱ")},
		"ィ": {'k': []byte("ｨ")},
		"イ": {'k': []byte("ｲ")},
		"ゥ": {'k': []byte("ｩ")},
		"ウ": {'k': []byte("ｳ")},
		"ェ": {'k': []byte("ｪ")},
		"エ": {'k': []byte("ｴ")},
		"ォ": {'k': []byte("ｫ")},
		"オ": {'k': []byte("ｵ")},
		"カ": {'k': []byte("ｶ")},
		"ガ": {'k': []byte("ｶﾞ")},
		"キ": {'k': []byte("ｷ")},
		"ギ": {'k': []byte("ｷﾞ")},
		"ク": {'k': []byte("ｸ")},
		"グ": {'k': []byte("ｸﾞ")},
		"ケ": {'k': []byte("ｹ")},
		"ゲ": {'k': []byte("ｹﾞ")},
		"コ": {'k': []byte("ｺ")},
		"ゴ": {'k': []byte("ｺﾞ")},
		"サ": {'k': []byte("ｻ")},
		"ザ": {'k': []byte("ｻﾞ")},
		"シ": {'k': []byte("ｼ")},
		"ジ": {'k': []byte("ｼﾞ")},
		"ス": {'k': []byte("ｽ")},
		"ズ": {'k': []byte("ｽﾞ")},
		"セ": {'k': []byte("ｾ")},
		"ゼ": {'k': []byte("ｾﾞ")},
		"ソ": {'k': []byte("ｿ")},
		"ゾ": {'k': []byte("ｿﾞ")},
		"タ": {'k': []byte("ﾀ")},
		"ダ": {'k': []byte("ﾀﾞ")},
		"チ": {'k': []byte("ﾁ")},
		"ヂ": {'k': []byte("ﾁﾞ")},
		"ッ": {'k': []byte("ｯ")},
		"ツ": {'k': []byte("ﾂ")},
		"ヅ": {'k': []byte("ﾂﾞ")},
		"テ": {'k': []byte("ﾃ")},
		"デ": {'k': []byte("ﾃﾞ")},
		"ト": {'k': []byte("ﾄ")},
		"ド": {'k': []byte("ﾄﾞ")},
		"ナ": {'k': []byte("ﾅ")},
		"ニ": {'k': []byte("ﾆ")},
		"ヌ": {'k': []byte("ﾇ")},
		"ネ": {'k': []byte("ﾈ")},
		"ノ": {'k': []byte("ﾉ")},
		"ハ": {'k': []byte("ﾊ")},
		"バ": {'k': []byte("ﾊﾞ")},
		"パ": {'k': []byte("ﾊﾟ")},
		"ヒ": {'k': []byte("ﾋ")},
		"ビ": {'k': []byte("ﾋﾞ")},
		"ピ": {'k': []byte("ﾋﾟ")},
		"フ": {'k': []byte("ﾌ")},
		"ブ": {'k': []byte("ﾌﾞ")},
		"プ": {'k': []byte("ﾌﾟ")},
		"ヘ": {'k': []byte("ﾍ")},
		"ベ": {'k': []byte("ﾍﾞ")},
		"ペ": {'k': []byte("ﾍﾟ")},
		"ホ": {'k': []byte("ﾎ")},
		"ボ": {'k': []byte("ﾎﾞ")},
		"ポ": {'k': []byte("ﾎﾟ")},
		"マ": {'k': []byte("ﾏ")},
		"ミ": {'k': []byte("ﾐ")},
		"ム": {'k': []byte("ﾑ")},
		"メ": {'k': []byte("ﾒ")},
		"モ": {'k': []byte("ﾓ")},
		"ャ": {'k': []byte("ｬ")},
		"ヤ": {'k': []byte("ﾔ")},
		"ュ": {'k': []byte("ｭ")},
		"ユ": {'k': []byte("ﾕ")},
		"ョ": {'k': []byte("ｮ")},
		"ヨ": {'k': []byte("ﾖ")},
		"ラ": {'k': []byte("ﾗ")},
		"リ": {'k': []byte("ﾘ")},
		"ル": {'k': []byte("ﾙ")},
		"レ": {'k': []byte("ﾚ")},
		"ロ": {'k': []byte("ﾛ")},
		"ヮ": {'k': []byte("ﾜ")},
		"ワ": {'k': []byte("ﾜ")},
		"ヰ": {'k': []byte("ｲ")},
		"ヱ": {'k': []byte("ｴ")},
		"ヲ": {'k': []byte("ｦ")},
		"ン": {'k': []byte("ﾝ")},
		"ヴ": {'k': []byte("ｳﾞ")},

		"・": {'h': []byte("･"), 'k': []byte("･")},
		"ー": {'h': []byte("ｰ"), 'k': []byte("ｰ")},
		"、": {'h': []byte("､"), 'k': []byte("､")},
		"。": {'h': []byte("｡"), 'k': []byte("｡")},
		"゛": {'h': []byte("ﾞ"), 'k': []byte("ﾞ")},
		"゜": {'h': []byte("ﾟ"), 'k': []byte("ﾟ")},

		"ｧ":  {'H': []byte("ぁ"), 'K': []byte("ァ")},
		"ｱ":  {'H': []byte("あ"), 'K': []byte("ア")},
		"ｨ":  {'H': []byte("ぃ"), 'K': []byte("ィ")},
		"ｲ":  {'H': []byte("い"), 'K': []byte("イ")},
		"ｩ":  {'H': []byte("ぅ"), 'K': []byte("ゥ")},
		"ｳ":  {'H': []byte("う"), 'K': []byte("ウ")},
		"ｪ":  {'H': []byte("ぇ"), 'K': []byte("ェ")},
		"ｴ":  {'H': []byte("え"), 'K': []byte("エ")},
		"ｫ":  {'H': []byte("ぉ"), 'K': []byte("ォ")},
		"ｵ":  {'H': []byte("お"), 'K': []byte("オ")},
		"ｶ":  {'H': []byte("か"), 'K': []byte("カ")},
		"ｶﾞ": {'H': []byte("が"), 'K': []byte("ガ")},
		"ｷ":  {'H': []byte("き"), 'K': []byte("キ")},
		"ｷﾞ": {'H': []byte("ぎ"), 'K': []byte("ギ")},
		"ｸ":  {'H': []byte("く"), 'K': []byte("ク")},
		"ｸﾞ": {'H': []byte("ぐ"), 'K': []byte("グ")},
		"ｹ":  {'H': []byte("け"), 'K': []byte("ケ")},
		"ｹﾞ": {'H': []byte("げ"), 'K': []byte("ゲ")},
		"ｺ":  {'H': []byte("こ"), 'K': []byte("コ")},
		"ｺﾞ": {'H': []byte("ご"), 'K': []byte("ゴ")},
		"ｻ":  {'H': []byte("さ"), 'K': []byte("サ")},
		"ｻﾞ": {'H': []byte("ざ"), 'K': []byte("ザ")},
		"ｼ":  {'H': []byte("し"), 'K': []byte("シ")},
		"ｼﾞ": {'H': []byte("じ"), 'K': []byte("ジ")},
		"ｽ":  {'H': []byte("す"), 'K': []byte("ス")},
		"ｽﾞ": {'H': []byte("ず"), 'K': []byte("ズ")},
		"ｾ":  {'H': []byte("せ"), 'K': []byte("セ")},
		"ｾﾞ": {'H': []byte("ぜ"), 'K': []byte("ゼ")},
		"ｿ":  {'H': []byte("そ"), 'K': []byte("ソ")},
		"ｿﾞ": {'H': []byte("ぞ"), 'K': []byte("ゾ")},
		"ﾀ":  {'H': []byte("た"), 'K': []byte("タ")},
		"ﾀﾞ": {'H': []byte("だ"), 'K': []byte("ダ")},
		"ﾁ":  {'H': []byte("ち"), 'K': []byte("チ")},
		"ﾁﾞ": {'H': []byte("ぢ"), 'K': []byte("ヂ")},
		"ｯ":  {'H': []byte("っ"), 'K': []byte("ッ")},
		"ﾂ":  {'H': []byte("つ"), 'K': []byte("ツ")},
		"ﾂﾞ": {'H': []byte("づ"), 'K': []byte("ヅ")},
		"ﾃ":  {'H': []byte("て"), 'K': []byte("テ")},
		"ﾃﾞ": {'H': []byte("で"), 'K': []byte("デ")},
		"ﾄ":  {'H': []byte("と"), 'K': []byte("ト")},
		"ﾄﾞ": {'H': []byte("ど"), 'K': []byte("ド")},
		"ﾅ":  {'H': []byte("な"), 'K': []byte("ナ")},
		"ﾆ":  {'H': []byte("に"), 'K': []byte("ニ")},
		"ﾇ":  {'H': []byte("ぬ"), 'K': []byte("ヌ")},
		"ﾈ":  {'H': []byte("ね"), 'K': []byte("ネ")},
		"ﾉ":  {'H': []byte("の"), 'K': []byte("ノ")},
		"ﾊ":  {'H': []byte("は"), 'K': []byte("ハ")},
		"ﾊﾞ": {'H': []byte("ば"), 'K': []byte("バ")},
		"ﾊﾟ": {'H': []byte("ぱ"), 'K': []byte("パ")},
		"ﾋ":  {'H': []byte("ひ"), 'K': []byte("ヒ")},
		"ﾋﾞ": {'H': []byte("び"), 'K': []byte("ビ")},
		"ﾋﾟ": {'H': []byte("ぴ"), 'K': []byte("ピ")},
		"ﾌ":  {'H': []byte("ふ"), 'K': []byte("フ")},
		"ﾌﾞ": {'H': []byte("ぶ"), 'K': []byte("ブ")},
		"ﾌﾟ": {'H': []byte("ぷ"), 'K': []byte("プ")},
		"ﾍ":  {'H': []byte("へ"), 'K': []byte("ヘ")},
		"ﾍﾞ": {'H': []byte("べ"), 'K': []byte("ベ")},
		"ﾍﾟ": {'H': []byte("ぺ"), 'K': []byte("ペ")},
		"ﾎ":  {'H': []byte("ほ"), 'K': []byte("ホ")},
		"ﾎﾞ": {'H': []byte("ぼ"), 'K': []byte("ボ")},
		"ﾎﾟ": {'H': []byte("ぽ"), 'K': []byte("ポ")},
		"ﾏ":  {'H': []byte("ま"), 'K': []byte("マ")},
		"ﾐ":  {'H': []byte("み"), 'K': []byte("ミ")},
		"ﾑ":  {'H': []byte("む"), 'K': []byte("ム")},
		"ﾒ":  {'H': []byte("め"), 'K': []byte("メ")},
		"ﾓ":  {'H': []byte("も"), 'K': []byte("モ")},
		"ｬ":  {'H': []byte("ゃ"), 'K': []byte("ャ")},
		"ﾔ":  {'H': []byte("や"), 'K': []byte("ヤ")},
		"ｭ":  {'H': []byte("ゅ"), 'K': []byte("ュ")},
		"ﾕ":  {'H': []byte("ゆ"), 'K': []byte("ユ")},
		"ｮ":  {'H': []byte("ょ"), 'K': []byte("ョ")},
		"ﾖ":  {'H': []byte("よ"), 'K': []byte("ヨ")},
		"ﾗ":  {'H': []byte("ら"), 'K': []byte("ラ")},
		"ﾘ":  {'H': []byte("り"), 'K': []byte("リ")},
		"ﾙ":  {'H': []byte("る"), 'K': []byte("ル")},
		"ﾚ":  {'H': []byte("れ"), 'K': []byte("レ")},
		"ﾛ":  {'H': []byte("ろ"), 'K': []byte("ロ")},
		"ﾜ":  {'H': []byte("わ"), 'K': []byte("ワ")},
		"ｦ":  {'H': []byte("を"), 'K': []byte("ヲ")},
		"ﾝ":  {'H': []byte("ん"), 'K': []byte("ン")},
		"ｳﾞ": {'H': []byte("う゛"), 'K': []byte("ヴ")},
		"ﾞ":  {'H': []byte("゛"), 'K': []byte("゛")},
		"ﾟ":  {'H': []byte("゜"), 'K': []byte("゜")},

		"､": {'H': []byte("、"), 'K': []byte("、")},
		"｡": {'H': []byte("。"), 'K': []byte("。")},
		"｢": {'H': []byte("「"), 'K': []byte("「")},
		"｣": {'H': []byte("」"), 'K': []byte("」")},
		"･": {'H': []byte("・"), 'K': []byte("・")},
		"ｰ": {'H': []byte("ー"), 'K': []byte("ー")},
	}
)

var filters map[byte]func(c *char) []byte = map[byte]func(c *char) []byte{
	'r': lowerR,
	'R': upperR,
	'n': lowerN,
	'N': upperN,
	'a': lowerA,
	'A': upperA,
	's': lowerS,
	'S': upperS,
	'k': lowerK,
	'K': upperK,
	'h': lowerH,
	'H': upperH,
	'c': lowerC,
	'C': upperC,
}

func createMode(m string) []byte {
	mode := []byte{}
	for _, m := range []byte(m) {
		flg := true
		for _, v := range mode {
			if m == v {
				flg = false
				break
			}
		}
		if flg {
			mode = append(mode, m)
		}
	}
	return mode
}

func Byte(b []byte, mode string) []byte {
	modes := createMode(mode)
	byteCount := uint64(len(b))
	buf := make([]byte, 0, byteCount)
	i := uint64(0)
	for i < byteCount {
		c := extract(b[i:])
		var v []byte
		for _, m := range modes {
			v = filters[m](c)
			if !bytes.Equal(c.val, v) {
				break
			}
		}
		buf = append(buf, v...)
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
		return 0, fmt.Errorf("buffer size is not enough")
	}
	n := copy(p, line)
	return n, nil
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

func isVoiced(b []byte) bool {
	if len(b) >= 6 && b[3] == 0xef && b[4] == 0xbe && b[5] == 0x9e {
		return true
	} else {
		return false
	}
}

func isSemiVoiced(b []byte) bool {
	if len(b) >= 6 && b[3] == 0xef && b[4] == 0xbe && b[5] == 0x9f {
		return true
	} else {
		return false
	}
}

func extract(b []byte) *char {
	if is1Byte(b) {
		if b[0] == 0x20 { // Space
			return &char{b[:1], []byte{'S'}, 1}
		} else if b[0] >= 0x30 && b[0] <= 0x39 { // 0 - 9
			return &char{b[:1], []byte{'N', 'A'}, 1}
		} else if b[0] >= 0x41 && b[0] <= 0x5a { // A - Z
			return &char{b[:1], []byte{'R', 'A'}, 1}
		} else if b[0] >= 0x61 && b[0] <= 0x7a { // a - z
			return &char{b[:1], []byte{'R', 'A'}, 1}
		} else if b[0] >= 0x21 && b[0] <= 0x7d && b[0] != 0x22 && b[0] != 0x27 && b[0] != 0x5c {
			return &char{b[:1], []byte{'A'}, 1}
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
				return &char{b[:3], []byte{}, 3}
			case 0xbd:
				if b[2] >= 0x81 && b[2] <= 0x9a { // ａ ～ ｚ
					return &char{b[:3], []byte{'r', 'a'}, 3}
				} else if b[2] >= 0x80 && b[2] <= 0x9d { // ｀ ～ ｝
					return &char{b[:3], []byte{'a'}, 3}
				} else if b[2] >= 0xa1 && b[2] <= 0xbf { // ｡ ｢ ｣ ､ ･ ｦ ～ ｿ
					if len(b) >= 6 {
						if b[2] == 0xb3 || (b[2] >= 0xb6 && b[2] <= 0xbf) { // ｳ, ｶ ～ ｿ
							if isVoiced(b) { // 濁点
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
					if isVoiced(b) { // 濁点
						return &char{b[:6], []byte{'K', 'H'}, 6}
					} else if isSemiVoiced(b) { // 半濁点
						return &char{b[:6], []byte{'K', 'H'}, 6}
					}
					return &char{b[:3], []byte{'K', 'H'}, 3}
				} else if b[2] >= 0x85 && b[2] <= 0x9f { // ﾅ ～ ﾝﾞﾟ
					return &char{b[:3], []byte{'K', 'H'}, 3}
				}
			}
			return &char{b[:3], []byte{}, 3}
		} else if b[0] == 0xe3 {
			// 全ひら・全カタ・全スペース
			switch b[1] {
			case 0x80:
				if b[2] == 0x80 { // スペース
					return &char{b[:3], []byte{'s'}, 3}
				} else if b[2] >= 0x81 && b[2] <= 0x82 { // 、。
					return &char{b[:3], []byte{'h', 'k', 'c'}, 3}
				}
			case 0x81:
				if b[2] >= 0x81 && b[2] <= 0xbf { // ぁ-み
					return &char{b[:3], []byte{'h', 'C'}, 3}
				}
			case 0x82:
				if b[2] >= 0x80 && b[2] <= 0x93 { // む-ん
					return &char{b[:3], []byte{'h', 'C'}, 3}
				} else if b[2] >= 0x9b && b[2] <= 0x9c { // ゛゜
					return &char{b[:3], []byte{'h', 'k'}, 3}
				} else if b[2] >= 0x9d && b[2] <= 0x9e { // ゝゞ
					return &char{b[:3], []byte{'C'}, 3}
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
func upperS(c *char) []byte {
	for _, m := range c.mode {
		if m != 'S' {
			continue
		}
		return []byte{0xe3, 0x80, 0x80}
	}
	return c.val
}

/**
 * Zenkaku Space -> Hankaku Space
 */
func lowerS(c *char) []byte {
	for _, m := range c.mode {
		if m != 's' {
			continue
		}
		return []byte{0x20}
	}
	return c.val
}

/**
 * Hankaku Numeric -> Zenkaku Numeric
 */
func upperN(c *char) []byte {
	for _, m := range c.mode {
		if m != 'N' {
			continue
		}
		return []byte{0xef, 0xbc, 0x60 + c.val[0]}
	}
	return c.val
}

/**
 * Zenkaku Numeric -> Hankaku Numeric
 */
func lowerN(c *char) []byte {
	for _, m := range c.mode {
		if m != 'n' {
			continue
		}
		return []byte{c.val[2] - 0x60}
	}
	return c.val
}

/**
 * Hankaku Alphabet -> Zenkaku Alphabet
 */
func upperR(c *char) []byte {
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
	return c.val
}

/**
 * Zenkaku Alphabet -> Hankaku Alphabet
 */
func lowerR(c *char) []byte {
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
	return c.val
}

/**
 * Hankaku AlphaNumeric -> Zenkaku AlphaNumeric
 * !-}(Excluding ",',\)
 */
func upperA(c *char) []byte {
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
	return c.val
}

/**
 * Zenkaku AlphaNumeric -> Hankaku AlphaNumeric
 * !-}(Excluding ",',\)
 */
func lowerA(c *char) []byte {
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
	return c.val
}

/**
 * Hankaku Katakana -> Zenkaku Katakana
 */
func upperK(c *char) []byte {
	key := byte('K')
	for _, m := range c.mode {
		if m != key {
			continue
		}
		s := string(c.val)
		// fmt.Printf("K: %s / %s\n", s, tbl[s][key])
		if _, ok := tbl[s][key]; ok {
			return tbl[s][key]
		}
		break
	}
	return c.val
}

/**
 * Zenkaku Katakana -> Hankaku Katakana
 */
func lowerK(c *char) []byte {
	key := byte('k')
	for _, m := range c.mode {
		if m != key {
			continue
		}
		s := string(c.val)
		if _, ok := tbl[s][key]; ok {
			return tbl[s][key]
		}
		break
	}
	return c.val
}

/**
 * Hankaku Katakana -> Zenkaku Hiragana
 */
func upperH(c *char) []byte {
	key := byte('H')
	for _, m := range c.mode {
		if m != key {
			continue
		}
		s := string(c.val)
		if _, ok := tbl[s][key]; ok {
			return tbl[s][key]
		}
		break
	}
	return c.val
}

/**
 * Zenkaku Hiragana -> Hankaku Katakana
 */
func lowerH(c *char) []byte {
	key := byte('h')
	for _, m := range c.mode {
		if m != key {
			continue
		}
		s := string(c.val)
		if _, ok := tbl[s][key]; ok {
			return tbl[s][key]
		}
		break
	}
	return c.val
}

/**
 * Zenkaku Hiragana -> Zenkaku Katakana
 */
func upperC(c *char) []byte {
	for _, m := range c.mode {
		if m != 'C' {
			continue
		}
		switch c.val[1] {
		case 0x81:
			if c.val[2] >= 0x81 && c.val[2] <= 0x9f { // ぁ-た
				return []byte{0xe3, 0x82, c.val[2] + 0x20}
			} else if c.val[2] >= 0xa0 && c.val[2] <= 0xbf { // だ-み
				return []byte{0xe3, 0x83, c.val[2] - 0x20}
			}
		case 0x82:
			if c.val[2] >= 0x80 && c.val[2] <= 0x93 { // む-ん
				return []byte{0xe3, 0x83, c.val[2] + 0x20}
			} else if c.val[2] >= 0x9d && c.val[2] <= 0x9e { // ゝゞ
				return []byte{0xe3, 0x83, c.val[2] + 0x20}
			}
		}
		break
	}
	return c.val
}

/**
 * Zenkaku Katakana -> Zenkaku Hiragana
 */
func lowerC(c *char) []byte {
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
	return c.val
}

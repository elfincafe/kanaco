package kanaco

import (
	"bufio"
	"io"
)

type Reader struct {
	r *bufio.Reader
	mode string
}

var tbl map[string][]byte = map[string][]byte{
	"ぁ": []byte{227, 129, 129}, "あ": []byte{227, 129, 130},
	"ぃ": []byte{227, 129, 131}, "い": []byte{227, 129, 132},
	"ぅ": []byte{227, 129, 133}, "う": []byte{227, 129, 134},
	"ぇ": []byte{227, 129, 135}, "え": []byte{227, 129, 136},
	"ぉ": []byte{227, 129, 137}, "お": []byte{227, 129, 138},
	"か": []byte{227, 129, 139}, "が": []byte{227, 129, 140},
	"き": []byte{227, 129, 141}, "ぎ": []byte{227, 129, 142},
	"く": []byte{227, 129, 143}, "ぐ": []byte{227, 129, 144},
	"け": []byte{227, 129, 145}, "げ": []byte{227, 129, 146},
	"こ": []byte{227, 129, 147}, "ご": []byte{227, 129, 148},
	"さ": []byte{227, 129, 149}, "ざ": []byte{227, 129, 150},
	"し": []byte{227, 129, 151}, "じ": []byte{227, 129, 152},
	"す": []byte{227, 129, 153}, "ず": []byte{227, 129, 154},
	"せ": []byte{227, 129, 155}, "ぜ": []byte{227, 129, 156},
	"そ": []byte{227, 129, 157}, "ぞ": []byte{227, 129, 158},
	"た": []byte{227, 129, 159}, "だ": []byte{227, 129, 160},
	"ち": []byte{227, 129, 161}, "ぢ": []byte{227, 129, 162},
	"っ": []byte{227, 129, 163}, "つ": []byte{227, 129, 164}, "づ": []byte{227, 129, 165},
	"て": []byte{227, 129, 166}, "で": []byte{227, 129, 167},
	"と": []byte{227, 129, 168}, "ど": []byte{227, 129, 169},
	"な": []byte{227, 129, 170},
	"に": []byte{227, 129, 171},
	"ぬ": []byte{227, 129, 172},
	"ね": []byte{227, 129, 173},
	"の": []byte{227, 129, 174},
	"は": []byte{227, 129, 175}, "ば": []byte{227, 129, 176}, "ぱ": []byte{227, 129, 177},
	"ひ": []byte{227, 129, 178}, "び": []byte{227, 129, 179}, "ぴ": []byte{227, 129, 180},
	"ふ": []byte{227, 129, 181}, "ぶ": []byte{227, 129, 182}, "ぷ": []byte{227, 129, 183},
	"へ": []byte{227, 129, 184}, "べ": []byte{227, 129, 185}, "ぺ": []byte{227, 129, 186},
	"ほ": []byte{227, 129, 187}, "ぼ": []byte{227, 129, 188}, "ぽ": []byte{227, 129, 189},
	"ま": []byte{227, 129, 190},
	"み": []byte{227, 129, 191},
	"む": []byte{227, 130, 128},
	"め": []byte{227, 130, 129},
	"も": []byte{227, 130, 130},
	"ゃ": []byte{227, 130, 131}, "や": []byte{227, 130, 132},
	"ゅ": []byte{227, 130, 133}, "ゆ": []byte{227, 130, 134},
	"ょ": []byte{227, 130, 135}, "よ": []byte{227, 130, 136},
	"ら": []byte{227, 130, 137},
	"り": []byte{227, 130, 138},
	"る": []byte{227, 130, 139},
	"れ": []byte{227, 130, 140},
	"ろ": []byte{227, 130, 141},
	"ゎ": []byte{227, 130, 142},
	"わ": []byte{227, 130, 143},
	"ゐ": []byte{227, 130, 144},
	"ゑ": []byte{227, 130, 145},
	"を": []byte{227, 130, 146},
	"ん": []byte{227, 130, 147},
	"ァ": []byte{227, 130, 161}, "ア": []byte{227, 130, 162},
	"ィ": []byte{227, 130, 163}, "イ": []byte{227, 130, 164},
	"ゥ": []byte{227, 130, 165}, "ウ": []byte{227, 130, 166},
	"ェ": []byte{227, 130, 167}, "エ": []byte{227, 130, 168},
	"ォ": []byte{227, 130, 169}, "オ": []byte{227, 130, 170},
	"カ": []byte{227, 130, 171}, "ガ": []byte{227, 130, 172},
	"キ": []byte{227, 130, 173}, "ギ": []byte{227, 130, 174},
	"ク": []byte{227, 130, 175}, "グ": []byte{227, 130, 176},
	"ケ": []byte{227, 130, 177}, "ゲ": []byte{227, 130, 178},
	"コ": []byte{227, 130, 179}, "ゴ": []byte{227, 130, 180},
	"サ": []byte{227, 130, 181}, "ザ": []byte{227, 130, 182},
	"シ": []byte{227, 130, 183}, "ジ": []byte{227, 130, 184},
	"ス": []byte{227, 130, 185}, "ズ": []byte{227, 130, 186},
	"セ": []byte{227, 130, 187}, "ゼ": []byte{227, 130, 188},
	"ソ": []byte{227, 130, 189}, "ゾ": []byte{227, 130, 190},
	"タ": []byte{227, 130, 191}, "ダ": []byte{227, 131, 128},
	"チ": []byte{227, 131, 129}, "ヂ": []byte{227, 131, 130},
	"ッ": []byte{227, 131, 131}, "ツ": []byte{227, 131, 132}, "ヅ": []byte{227, 131, 133},
	"テ": []byte{227, 131, 134}, "デ": []byte{227, 131, 135},
	"ト": []byte{227, 131, 136}, "ド": []byte{227, 131, 137},
	"ナ": []byte{227, 131, 138},
	"ニ": []byte{227, 131, 139},
	"ヌ": []byte{227, 131, 140},
	"ネ": []byte{227, 131, 141},
	"ノ": []byte{227, 131, 142},
	"ハ": []byte{227, 131, 143}, "バ": []byte{227, 131, 144}, "パ": []byte{227, 131, 145},
	"ヒ": []byte{227, 131, 146}, "ビ": []byte{227, 131, 147}, "ピ": []byte{227, 131, 148},
	"フ": []byte{227, 131, 149}, "ブ": []byte{227, 131, 150}, "プ": []byte{227, 131, 151},
	"ヘ": []byte{227, 131, 152}, "ベ": []byte{227, 131, 153}, "ペ": []byte{227, 131, 154},
	"ホ": []byte{227, 131, 155}, "ボ": []byte{227, 131, 156}, "ポ": []byte{227, 131, 157},
	"マ": []byte{227, 131, 158},
	"ミ": []byte{227, 131, 159},
	"ム": []byte{227, 131, 160},
	"メ": []byte{227, 131, 161},
	"モ": []byte{227, 131, 162},
	"ャ": []byte{227, 131, 163}, "ヤ": []byte{227, 131, 164},
	"ュ": []byte{227, 131, 165}, "ユ": []byte{227, 131, 166},
	"ョ": []byte{227, 131, 167}, "ヨ": []byte{227, 131, 168},
	"ラ": []byte{227, 131, 169},
	"リ": []byte{227, 131, 170},
	"ル": []byte{227, 131, 171},
	"レ": []byte{227, 131, 172},
	"ロ": []byte{227, 131, 173},
	"ヮ": []byte{227, 131, 174}, "ワ": []byte{227, 131, 175},
	"ヰ": []byte{227, 131, 176},
	"ヱ": []byte{227, 131, 177},
	"ヲ": []byte{227, 131, 178},
	"ン": []byte{227, 131, 179},
	"ー": []byte{227, 131, 188},
	"ｦ": []byte{239, 189, 166},
	"ｧ": []byte{239, 189, 167}, "ｨ": []byte{239, 189, 168}, "ｩ": []byte{239, 189, 169}, "ｪ": []byte{239, 189, 170}, "ｫ": []byte{239, 189, 171},
	"ｬ": []byte{239, 189, 172}, "ｭ": []byte{239, 189, 173}, "ｮ": []byte{239, 189, 174},
	"ｯ": []byte{239, 189, 175}, "ｰ": []byte{239, 189, 176},
	"ｱ": []byte{239, 189, 177}, "ｲ": []byte{239, 189, 178}, "ｳ": []byte{239, 189, 179}, "ｴ": []byte{239, 189, 180}, "ｵ": []byte{239, 189, 181},
	"ｶ": []byte{239, 189, 182}, "ｷ": []byte{239, 189, 183}, "ｸ": []byte{239, 189, 184}, "ｹ": []byte{239, 189, 185}, "ｺ": []byte{239, 189, 186},
	"ｶﾞ": []byte{239, 189, 182,239, 190, 158}, "ｷﾞ": []byte{239, 189, 183,239, 190, 158}, "ｸﾞ": []byte{239, 189, 184,239, 190, 158},
	"ｹﾞ": []byte{239, 189, 185,239, 190, 158}, "ｺﾞ": []byte{239, 189, 186,239, 190, 158},
	"ｻ": []byte{239, 189, 187}, "ｼ": []byte{239, 189, 188}, "ｽ": []byte{239, 189, 189}, "ｾ": []byte{239, 189, 190}, "ｿ": []byte{239, 189, 191},
	"ｻﾞ": []byte{239, 189, 187,239, 190, 158}, "ｼﾞ": []byte{239, 189, 188,239, 190, 158}, "ｽﾞ": []byte{239, 189, 189,239, 190, 158},
	"ｾﾞ": []byte{239, 189, 190,239, 190, 158}, "ｿﾞ": []byte{239, 189, 191,239, 190, 158},
	"ﾀ": []byte{239, 190, 128}, "ﾁ": []byte{239, 190, 129}, "ﾂ": []byte{239, 190, 130}, "ﾃ": []byte{239, 190, 131}, "ﾄ": []byte{239, 190, 132},
	"ﾀﾞ": []byte{239, 190, 128,239, 190, 158}, "ﾁﾞ": []byte{239, 190, 129,239, 190, 158}, "ﾂﾞ": []byte{239, 190, 130,239, 190, 158},
	"ﾃﾞ": []byte{239, 190, 131,239, 190, 158}, "ﾄﾞ": []byte{239, 190, 132,239, 190, 158},
	"ﾅ": []byte{239, 190, 133}, "ﾆ": []byte{239, 190, 134}, "ﾇ": []byte{239, 190, 135}, "ﾈ": []byte{239, 190, 136}, "ﾉ": []byte{239, 190, 137},
	"ﾊ": []byte{239, 190, 138}, "ﾋ": []byte{239, 190, 139}, "ﾌ": []byte{239, 190, 140}, "ﾍ": []byte{239, 190, 141}, "ﾎ": []byte{239, 190, 142},
	"ﾊﾞ": []byte{239, 190, 138,239, 190, 158}, "ﾋﾞ": []byte{239, 190, 139,239, 190, 158}, "ﾌﾞ": []byte{239, 190, 140,239, 190, 158},
	"ﾍﾞ": []byte{239, 190, 141,239, 190, 158}, "ﾎﾞ": []byte{239, 190, 142,239, 190, 158},
	"ﾊﾟ": []byte{239, 190, 138,239, 190, 159}, "ﾋﾟ": []byte{239, 190, 139,239, 190, 159}, "ﾌﾟ": []byte{239, 190, 140,239, 190, 159},
	"ﾍﾟ": []byte{239, 190, 141,239, 190, 159}, "ﾎﾟ": []byte{239, 190, 142,239, 190, 159},
	"ﾏ": []byte{239, 190, 143}, "ﾐ": []byte{239, 190, 144}, "ﾑ": []byte{239, 190, 145}, "ﾒ": []byte{239, 190, 146}, "ﾓ": []byte{239, 190, 147},
	"ﾔ": []byte{239, 190, 148}, "ﾕ": []byte{239, 190, 149}, "ﾖ": []byte{239, 190, 150},
	"ﾗ": []byte{239, 190, 151}, "ﾘ": []byte{239, 190, 152}, "ﾙ": []byte{239, 190, 153}, "ﾚ": []byte{239, 190, 154}, "ﾛ": []byte{239, 190, 155},
	"ﾜ": []byte{239, 190, 156}, "ﾝ": []byte{239, 190, 157},
	"ﾞ": []byte{239, 190, 158}, "ﾟ": []byte{239, 190, 159},
}


func Byte (b []byte, mode string) []byte {
	filters := map[byte]func([]byte) []byte{}
	for _, m := range mode {
		switch m {
		case 'r':
			filters['r'] = convAsSmallR
		case 'R':
			filters['R'] = convAsLargeR
		case 'n':
			filters['n'] = convAsSmallN
		case 'N':
			filters['N'] = convAsLargeN
		case 'a':
			filters['a'] = convAsSmallA
		case 'A':
			filters['A'] = convAsLargeA
		case 's':
			filters['s'] = convAsSmallS
		case 'S':
			filters['S'] = convAsLargeS
		case 'k':
			filters['k'] = convAsSmallK
		case 'K':
			filters['K'] = convAsLargeK
		case 'h':
			filters['h'] = convAsSmallH
		case 'H':
			filters['H'] = convAsLargeH
		case 'c':
			filters['c'] = convAsSmallC
		case 'C':
			filters['C'] = convAsLargeC
		}
	}
	byteCount := uint64(len(b))
	buf := make([]byte, 0, byteCount)
	for i := uint64(0); i < byteCount; i++ {
		c := make([]byte, 0, 6)
		count := count(b[i:])
		e := i + count
		c = b[i:e]
		for _, f := range filters {
			c = f(c)
		}
		buf = append(buf, c...)
		i += count - 1
	}
	return buf
}

func String (s string, mode string) string {
	return string(Byte([]byte(s), mode))
}

func NewReader (r io.Reader, mode string) *Reader {
	reader := new(Reader)
	reader.r = bufio.NewReader(r)
	reader.mode = mode
	return reader
}

func (r *Reader) Read(p []byte) (int, error) {
	line, err := r.r.ReadBytes('\n')
	line = Byte(line, r.mode)
	n := copy(p, line)
	return n, err
}

func count(b []byte) uint64 {
	if b[0] < 128 { // 1byte
		return uint64(1)
	} else if b[0] < 192 {
		return uint64(0)
	} else if b[0] < 224 { // 2byte
		return uint64(2)
	} else if b[0] < 240 { // 3byte
		if len(b) >= 6 {
			if b[0] == 239 && b[1] == 189 {
				if b[2] >= 182 && b[2] <= 191 { // カ・サ行
					if b[3] == 239 && b[4] == 190 && b[5] == 158 {
						return uint64(6)
					}
				}
			} else if b[0] == 239 && b[1] == 190 {
				if b[2] >= 128 && b[2] <= 132 { // タ行
					if b[3] == 239 && b[4] == 190 && b[5] == 158 {
						return uint64(6)
					}
				} else if b[2] >= 138 && b[2] <= 142 { // ハ行
					if b[3] == 239 && b[4] == 190 && (b[5] == 158 || b[5] == 159) {
						return uint64(6)
					}
				}
			}
		}
		return uint64(3)
	} else if b[0] < 248 { // 4byte
		return uint64(4)
	}
	return uint64(0)
}

/**
 * Hankaku Space -> Zenkaku Space
 */
func convAsLargeS(b []byte) []byte {
	if len(b) == 1 && b[0] == 32 {
		return []byte{227, 128, 128}
	}
	return b
}

/**
 * Zenkaku Space -> Hankaku Space
 */
func convAsSmallS(b []byte) []byte {
	if len(b) == 3 && b[0] == 227 && b[1] == 128 && b[2] == 128 {
		return []byte{32}
	}
	return b
}

/**
 * Hankaku Numeric -> Zenkaku Numeric
 */
func convAsLargeN(b []byte) []byte {
	if len(b) == 1 && b[0] >= 48 && b[0] <= 57 {
		return []byte{239, 188, 96 + b[0]}
	}
	return b
}

/**
 * Zenkaku Numeric -> Hankaku Numeric
 */
func convAsSmallN(b []byte) []byte {
	if len(b) == 3 && b[0] == 239 && b[1] == 188 && (b[2] >= 144 && b[2] <= 153) {
		return []byte{b[2] - 96}
	}
	return b
}

/**
 * Hankaku Alphabet -> Zenkaku Alphabet
 */
func convAsLargeR(b []byte) []byte {
	if len(b) != 1 {
		return b
	}
	// A-Z -> Ａ-Ｚ
	if b[0] >= 65 && b[0] <= 90 {
		return []byte{239, 188, 96 + b[0]}
	}
	// a-z -> ａ-ｚ
	if b[0] >= 97 && b[0] <= 122 {
		return []byte{239, 189, 32 + b[0]}
	}
	return b
}

/**
 * Zenkaku Alphabet -> Hankaku Alphabet
 */
func convAsSmallR(b []byte) []byte {
	if len(b) != 3 {
		return b
	}
	// Ａ-Ｚ -> A-Z
	if b[0] == 239 && b[1] == 188 && (b[2] >= 161 && b[2] <= 186) {
		return []byte{b[2] - 96}
	}
	// ａ-ｚ -> a-z
	if b[0] == 239 && b[1] == 189 && (b[2] >= 129 && b[2] <= 154) {
		return []byte{b[2] - 32}
	}
	return b
}

/**
 * Hankaku AlphaNumeric -> Zenkaku AlphaNumeric
 * !-}(Excluding ",',\)
 */
func convAsLargeA(b []byte) []byte {
	if len(b) != 1 {
		return b
	}
	if b[0] > 32 && b[0] < 96 {
		if b[0] == 34 || b[0] == 39 || b[0] == 92 {
			return b
		}
		return []byte{239, 188, 96 + b[0]}
	} else if b[0] > 95 && b[0] < 126 {
		return []byte{239, 189, 32 + b[0]}
	}
	return b
}

/**
 * Zenkaku AlphaNumeric -> Hankaku AlphaNumeric
 * !-}(Excluding ",',\)
 */
func convAsSmallA(b []byte) []byte {
	if len(b) != 3 {
		return b
	}
	// ！-＿ -> !-_
	if b[0] == 239 && b[1] == 188 && (b[2] >= 129 && b[2] <= 191) {
		return []byte{b[2] - 96}
	}
	// ｀-｝ -> `-}
	if b[0] == 239 && b[1] == 189 && (b[2] >= 128 && b[2] <= 157) {
		return []byte{b[2] - 32}
	}
	return b
}

/**
 * Zenkaku Katakana -> Hankaku Katakana
 */
func convAsSmallK(b []byte) []byte {
	if len(b) != 3 {
		return b
	}
	if b[1] == 130 && b[0]==227{
		switch b[2] {
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
	} else if b[1] == 131 && b[0]==227 {
		switch b[2] {
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
	return b
}

/**
 * Hankaku Katakana -> Zenkaku Katakana
 */
func convAsLargeK(b []byte) []byte {
	len := len(b)
	if len != 3 && len != 6 {
		return b
	}
	if b[1] == 189 && b[0] == 239 {
		switch b[2] {
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
			if len == 6 && b[5] == 158 && b[4] == 190 && b[3] == 239 {
				return tbl["ガ"]
			}
			return tbl["カ"]
		case 183: // ｷ
			if len == 6 && b[5] == 158 && b[4] == 190 && b[3] == 239 {
				return tbl["ギ"]
			}
			return tbl["キ"]
		case 184: // ｸ
			if len == 6 && b[5] == 158 && b[4] == 190 && b[3] == 239 {
				return tbl["グ"]
			}
			return tbl["ク"]
		case 185: // ｹ
			if len == 6 && b[5] == 158 && b[4] == 190 && b[3] == 239 {
				return tbl["ゲ"]
			}
			return tbl["ケ"]
		case 186: // ｺ
			if len == 6 && b[5] == 158 && b[4] == 190 && b[3] == 239 {
				return tbl["ゴ"]
			}
			return tbl["コ"]
		case 187: // ｻ
			if len == 6 && b[5] == 158 && b[4] == 190 && b[3] == 239 {
				return tbl["ザ"]
			}
			return tbl["サ"]
		case 188: // ｼ
			if len == 6 && b[5] == 158 && b[4] == 190 && b[3] == 239 {
				return tbl["ジ"]
			}
			return tbl["シ"]
		case 189: // ｽ
			if len == 6 && b[5] == 158 && b[4] == 190 && b[3] == 239 {
				return tbl["ズ"]
			}
			return tbl["ス"]
		case 190: // ｾ
			if len == 6 && b[5] == 158 && b[4] == 190 && b[3] == 239 {
				return tbl["ゼ"]
			}
			return tbl["セ"]
		case 191: // ｿ
			if len == 6 && b[5] == 158 && b[4] == 190 && b[3] == 239 {
				return tbl["ゾ"]
			}
			return tbl["ソ"]
		}
	} else if b[1] == 190 && b[0] == 239 {
		switch b[2] {
		case 128: // ﾀ
			if len == 6 && b[5] == 158 && b[4] == 190 && b[3] == 239 {
				return tbl["ダ"]
			}
			return tbl["タ"]
		case 129: // ﾁ
			if len == 6 && b[5] == 158 && b[4] == 190 && b[3] == 239 {
				return tbl["ヂ"]
			}
			return tbl["チ"]
		case 130: // ﾂ
			if len == 6 && b[5] == 158 && b[4] == 190 && b[3] == 239 {
				return tbl["ヅ"]
			}
			return tbl["ツ"]
		case 131: // ﾃ
			if len == 6 && b[5] == 158 && b[4] == 190 && b[3] == 239 {
				return tbl["デ"]
			}
			return tbl["テ"]
		case 132: // ﾄ
			if len == 6 && b[5] == 158 && b[4] == 190 && b[3] == 239 {
				return tbl["ド"]
			}
			return tbl["ト"]
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
			if len == 6 && b[5] == 159 && b[4] == 190 && b[3] == 239 {
				return tbl["パ"]
			} else if len == 6 && b[5] == 158 && b[4] == 190 && b[3] == 239 {
				return tbl["バ"]
			}
			return tbl["ハ"]
		case 139: // ヒ
			if len == 6 && b[5] == 159 && b[4] == 190 && b[3] == 239 {
				return tbl["ピ"]
			} else if len == 6 && b[5] == 158 && b[4] == 190 && b[3] == 239 {
				return tbl["ビ"]
			}
			return tbl["ヒ"]
		case 140: // フ
			if len == 6 && b[5] == 159 && b[4] == 190 && b[3] == 239 {
				return tbl["プ"]
			} else if len == 6 && b[5] == 158 && b[4] == 190 && b[3] == 239 {
				return tbl["ブ"]
			}
			return tbl["フ"]
		case 141: // ヘ
			if len == 6 && b[5] == 159 && b[4] == 190 && b[3] == 239 {
				return tbl["ペ"]
			} else if len == 6 && b[5] == 158 && b[4] == 190 && b[3] == 239 {
				return tbl["ベ"]
			}
			return tbl["ヘ"]
		case 142: // ホ
			if len == 6 && b[5] == 159 && b[4] == 190 && b[3] == 239 {
				return tbl["ポ"]
			} else if len == 6 && b[5] == 158 && b[4] == 190 && b[3] == 239 {
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
	return b
}

/**
 * Zenkaku Hiragana -> Hankaku Katakana
 */
func convAsSmallH(b []byte) []byte {
	len := len(b)
	if len != 3 {
		return b
	}
	if b[1] == 129 && b[0] == 227 {
		switch b[2] {
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
	} else if b[1] == 130 && b[0] == 227 {
		switch b[2] {
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
	} else if b[2]==188 && b[1] == 131 && b[0]==227 {
		return tbl["ｰ"]
	}
	return b
}

/**
 * Hankaku Katakana -> Zenkaku Hiragana
 */
func convAsLargeH(b []byte) []byte {
	len := len(b)
	if len != 3 && len != 6 {
		return b
	}
	if b[1] == 189 && b[0] == 239 {
		switch b[2] {
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
			if len == 6 && b[5] == 158 && b[4] == 190 && b[3] == 239 {
				return tbl["が"]
			}
			return tbl["か"]
		case 183: // ｷ
			if len == 6 && b[5] == 158 && b[4] == 190 && b[3] == 239 {
				return tbl["ぎ"]
			}
			return tbl["き"]
		case 184: // ｸ
			if len == 6 && b[5] == 158 && b[4] == 190 && b[3] == 239 {
				return tbl["ぐ"]
			}
			return tbl["く"]
		case 185: // ｹ
			if len == 6 && b[5] == 158 && b[4] == 190 && b[3] == 239 {
				return tbl["げ"]
			}
			return tbl["け"]
		case 186: // ｺ
			if len == 6 && b[5] == 158 && b[4] == 190 && b[3] == 239 {
				return tbl["ご"]
			}
			return tbl["こ"]
		case 187: // ｻ
			if len == 6 && b[5] == 158 && b[4] == 190 && b[3] == 239 {
				return tbl["ざ"]
			}
			return tbl["さ"]
		case 188: // ｼ
			if len == 6 && b[5] == 158 && b[4] == 190 && b[3] == 239 {
				return tbl["じ"]
			}
			return tbl["し"]
		case 189: // ｽ
			if len == 6 && b[5] == 158 && b[4] == 190 && b[3] == 239 {
				return tbl["ず"]
			}
			return tbl["す"]
		case 190: // ｾ
			if len == 6 && b[5] == 158 && b[4] == 190 && b[3] == 239 {
				return tbl["ぜ"]
			}
			return tbl["せ"]
		case 191: // ｿ
			if len == 6 && b[5] == 158 && b[4] == 190 && b[3] == 239 {
				return tbl["ぞ"]
			}
			return tbl["そ"]
		}
	} else if b[1] == 190 && b[0] == 239 {
		switch b[2] {
		case 128: // ﾀ
			if len == 6 && b[5] == 158 && b[4] == 190 && b[3] == 239 {
				return tbl["だ"]
			}
			return tbl["た"]
		case 129: // ﾁ
			if len == 6 && b[5] == 158 && b[4] == 190 && b[3] == 239 {
				return tbl["ぢ"]
			}
			return tbl["ち"]
		case 130: // ﾂ
			if len == 6 && b[5] == 158 && b[4] == 190 && b[3] == 239 {
				return tbl["づ"]
			}
			return tbl["つ"]
		case 131: // ﾃ
			if len == 6 && b[5] == 158 && b[4] == 190 && b[3] == 239 {
				return tbl["で"]
			}
			return tbl["て"]
		case 132: // ﾄ
			if len == 6 && b[5] == 158 && b[4] == 190 && b[3] == 239 {
				return tbl["ど"]
			}
			return tbl["と"]
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
			if len == 6 && b[5] == 159 && b[4] == 190 && b[3] == 239 {
				return tbl["ぱ"]
			} else if len == 6 && b[5] == 158 && b[4] == 190 && b[3] == 239 {
				return tbl["ば"]
			}
			return tbl["は"]
		case 139: // ﾋ
			if len == 6 && b[5] == 159 && b[4] == 190 && b[3] == 239 {
				return tbl["ぴ"]
			} else if len == 6 && b[5] == 158 && b[4] == 190 && b[3] == 239 {
				return tbl["び"]
			}
			return tbl["ひ"]
		case 140: // ﾌ
			if len == 6 && b[5] == 159 && b[4] == 190 && b[3] == 239 {
				return tbl["ぷ"]
			} else if len == 6 && b[5] == 158 && b[4] == 190 && b[3] == 239 {
				return tbl["ぶ"]
			}
			return tbl["ふ"]
		case 141: // ﾍ
			if len == 6 && b[5] == 159 && b[4] == 190 && b[3] == 239 {
				return tbl["ぺ"]
			} else if len == 6 && b[5] == 158 && b[4] == 190 && b[3] == 239 {
				return tbl["べ"]
			}
			return tbl["へ"]
		case 142: // ﾎ
			if len == 6 && b[5] == 159 && b[4] == 190 && b[3] == 239 {
				return tbl["ぽ"]
			} else if len == 6 && b[5] == 158 && b[4] == 190 && b[3] == 239 {
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
	return b
}

/**
 * Zenkaku Katakana -> Zenkaku Hiragana
 */
func convAsSmallC(b []byte) []byte {
	if len(b) != 3 {
		return b
	}
	if b[0] != 227 {
		return b
	}
	if b[1] == 130 { // ァ-タ
		if b[2] >= 161 && b[2] <= 191 {
			return []byte{227, 129, b[2] - 32}
		}
	} else if b[1] == 131 { // ダ-ン
		if b[2] >= 128 && b[2] <= 159 { // ダ-ミ
			return []byte{227, 129, b[2] + 32}
		} else if b[2] >= 160 && b[2] <= 179 { // ム-ン
			return []byte{227, 130, b[2] - 32}
		}
	}
	return b
}

/**
 * Zenkaku Hiragana -> Zenkaku Katakana
 */
func convAsLargeC(b []byte) []byte {
	if len(b) != 3 {
		return b
	}
	if b[0] != 227 {
		return b
	}
	if b[1] == 129 { // ぁ-み
		if b[2] >= 129 && b[2] <= 159 { // ぁ-た
			return []byte{227, 130, b[2] + 32}
		} else if b[2] >= 160 && b[2] <= 191 { // だ-み
			return []byte{227, 131, b[2] - 32}
		}
	} else if b[1] == 130 { // む-ん
		if b[2] >= 128 && b[2] <= 147 {
			return []byte{227, 131, b[2] + 32}
		}
	}
	return b
}

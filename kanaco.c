#include "kanaco.h"

#include <stdbool.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#define MODE_BUF 16
#define RET_BUF 4096

bool is_1byte(const char *s, int len) {
#ifdef DEBUG
  printf("Func: %s\n", "is_1byte");
#endif

  uint8_t c0 = *s & 0xff;
  if (len > 0 && c0 < 0x80) {
    return true;
  }
  return false;
}

bool is_2bytes(const char *s, int len) {
#ifdef DEBUG
  printf("Func: %s\n", "is_2bytes");
#endif

  uint8_t c0 = *s & 0xff, c1 = *(s + 1) & 0xff;
  if (len > 1 && c0 > 0xc1 && c0 < 0xe0 && c1 > 0x7f && c1 < 0xc0) {
    return true;
  }
  return false;
}

bool is_3bytes(const char *s, int len) {
#ifdef DEBUG
  printf("Func: %s\n", "is_3bytes");
#endif
  uint8_t c0 = *s & 0xff, c1 = *(s + 1) & 0xff, c2 = *(s + 2) & 0xff;
  if (len > 2 && c0 > 0xdf && c0 < 0xf0 && c1 > 0x7f && c1 < 0xc0 &&
      c2 > 0x7f && c2 < 0xc0) {
    return true;
  }
  return false;
}

bool is_4bytes(const char *s, int len) {
#ifdef DEBUG
  printf("Func: %s\n", "is_4bytes");
#endif
  uint8_t c0 = *s & 0xff, c1 = *(s + 1) & 0xff, c2 = *(s + 2) & 0xff,
          c3 = *(s + 3) & 0xff;
  if (len > 3 && c0 > 0xef && c0 < 0xf5 && c1 > 0x7f && c1 < 0xc0 &&
      c2 > 0x7f && c2 < 0xc0 && c3 > 0x7f && c3 < 0xc0) {
    return true;
  }
  return false;
}

bool is_voiced(const char *s, int len) {
#ifdef DEBUG
  printf("Func: %s\n", "is_voiced");
#endif
  if (len > 5) {
    uint8_t c3 = *(s + 3) & 0xff, c4 = *(s + 4) & 0xff, c5 = *(s + 5) & 0xff;
    if (c3 == 0xef && c4 == 0xbe && c5 == 0x9e) {
      uint8_t c0 = *s & 0xff, c1 = *(s + 1) & 0xff, c2 = *(s + 2) & 0xff;
      if (c0 == 0xef && c1 == 0xbd && c2 > 0xb5 && c2 < 0xc0) {  // ｶ - ｿ
        return true;
      } else if (c0 == 0xef && c1 == 0xbe && c2 > 0x79 && c2 < 0x85) {  // ﾀ - ﾄ
        return true;
      } else if (c0 == 0xef && c1 == 0xbe && c2 > 0x89 && c2 < 0x8f) {  // ﾊ - ﾎ
        return true;
      } else if (c0 == 0xef && c1 == 0xbd && c2 == 0xb3) {  // ｳ
        return true;
      }
    }
  }
  return false;
}

bool is_semi_voiced(const char *s, int len) {
#ifdef DEBUG
  printf("Func: %s\n", "is_semi_voiced");
#endif
  if (len > 5) {
    uint8_t c3 = *(s + 3) & 0xff, c4 = *(s + 4) & 0xff, c5 = *(s + 5) & 0xff;
    if (c3 == 0xef && c4 == 0xbe && c5 == 0x9f) {
      uint8_t c0 = *s & 0xff, c1 = *(s + 1) & 0xff, c2 = *(s + 2) & 0xff;
      if (c0 == 0xef && c1 == 0xbe && c2 > 0x89 && c2 < 0x8f) {  // ﾊ - ﾎ
        return true;
      }
    }
  }
  return false;
}

void lower_r(character *c) {
#ifdef DEBUG
  printf("Func: %s [%s]\n", "lower_r", c->val);
#endif

  if (!(c->conv & CNV_LOWER_R)) {
    return;
  }
  uint8_t c2 = (uint8_t)(*(c->val + 2) & 0xff);
  // Ａ-Ｚ -> A-Z
  if (c2 >= 0xa1 && c2 <= 0xba) {
    *(c->cval + 0) = (char)((c2 - 0x60) & 0xff);
    *(c->cval + 1) = 0x00;
    c->clen = 1;
    return;
  }
  // ａ-ｚ -> a-z
  if (c2 >= 0x81 && c2 <= 0x9a) {
    *(c->cval + 0) = (char)((c2 - 0x20) & 0xff);
    *(c->cval + 1) = 0x00;
    c->clen = 1;
    return;
  }
}

void upper_r(character *c) {
#ifdef DEBUG
  printf("Func: %s [%s]\n", "upper_r", c->val);
#endif

  if (!(c->conv & CNV_UPPER_R)) {
    return;
  }
  uint8_t c0 = (uint8_t)(*(c->val + 0));
  if (c0 >= 0x41 && c0 <= 0x5a) {  // A-Z -> Ａ-Ｚ
    c->clen = 3;
    *(c->cval + 0) = (char)(0xef & 0xff);
    *(c->cval + 1) = (char)(0xbc & 0xff);
    *(c->cval + 2) = (char)((c0 + 0x60) & 0xff);
    *(c->cval + c->clen) = 0x00;
  } else if (c0 >= 0x61 && c0 <= 0x7a) {  // a-z -> ａ-ｚ
    c->clen = 3;
    *(c->cval + 0) = (char)(0xef & 0xff);
    *(c->cval + 1) = (char)(0xbd & 0xff);
    *(c->cval + 2) = (char)((c0 + 0x20) & 0xff);
    *(c->cval + c->clen) = 0x00;
  }
}

void lower_n(character *c) {
#ifdef DEBUG
  printf("Func: %s [%s]\n", "lower_n", c->val);
#endif

  if (!(c->conv & CNV_LOWER_N)) {
    return;
  }
  uint8_t c2 = (uint8_t)(*(c->val + 2));
  c->clen = 1;
  *(c->cval + 0) = (char)((c2 - 0x60) & 0xff);
  *(c->cval + c->clen) = 0x00;
}

void upper_n(character *c) {
#ifdef DEBUG
  printf("Func: %s [%s]\n", "upper_n", c->val);
#endif
  if (!(c->conv & CNV_UPPER_N)) {
    return;
  }
  uint8_t c0 = (uint8_t)(*(c->val + 0) & 0xff);
  c->clen = 3;
  *(c->cval + 0) = (char)(0xef & 0xff);
  *(c->cval + 1) = (char)(0xbc & 0xff);
  *(c->cval + 2) = (char)((c0 + 0x60) & 0xff);
  *(c->cval + c->clen) = 0x00;
}

void lower_a(character *c) {
#ifdef DEBUG
  printf("Func: %s [%s]\n", "lower_a", c->val);
#endif

  if (!(c->conv & CNV_LOWER_A)) {
    return;
  }
  uint8_t c1 = (uint8_t)(*(c->val + 1)), c2 = (uint8_t)(*(c->val + 2));
  if (c1 == 0xbc && c2 >= 0x81 && c2 <= 0xbf) {
    c->clen = 1;
    *(c->cval + 0) = (char)((c2 - 0x60) & 0xff);
    *(c->cval + c->clen) = 0x00;
  } else if (c1 == 0xbd && c2 >= 0x80 && c2 <= 0x9d) {
    c->clen = 1;
    *(c->cval + 0) = (char)((c2 - 0x20) & 0xff);
    *(c->cval + c->clen) = 0x00;
  }
}

void upper_a(character *c) {
#ifdef DEBUG
  printf("Func: %s [%s]\n", "upper_a", c->val);
#endif

  if (!(c->conv & CNV_UPPER_A)) {
    return;
  }
  uint8_t c0 = (uint8_t)(*(c->val + 0));
  if (c0 >= 0x21 && c0 <= 0x5f) {
    c->clen = 3;
    *(c->cval + 0) = (char)(0xef & 0xff);
    *(c->cval + 1) = (char)(0xbc & 0xff);
    *(c->cval + 2) = (char)((c0 + 0x60) & 0xff);
    *(c->cval + c->clen) = 0x00;
  } else if (c0 >= 0x60 && c0 <= 0x7d) {
    c->clen = 3;
    *(c->cval + 0) = (char)(0xef & 0xff);
    *(c->cval + 1) = (char)(0xbd & 0xff);
    *(c->cval + 2) = (char)((c0 + 0x20) & 0xff);
    *(c->cval + c->clen) = 0x00;
  }
}

void lower_s(character *c) {
#ifdef DEBUG
  printf("Func: %s [%s]\n", "lower_s", c->val);
#endif

  if (!(c->conv & CNV_LOWER_S)) {
    return;
  }
  c->clen = 1;
  *(c->cval + 0) = (char)(0x20 & 0xff);
  *(c->cval + c->clen) = 0x00;
}

void upper_s(character *c) {
#ifdef DEBUG
  printf("Func: %s\n", "upper_s");
#endif

  if (!(c->conv & CNV_UPPER_S)) {
    return;
  }
  c->clen = 3;
  *(c->cval + 0) = (char)(0xe3 & 0xff);
  *(c->cval + 1) = (char)(0x80 & 0xff);
  *(c->cval + 2) = (char)(0x80 & 0xff);
  *(c->cval + c->clen) = 0x00;
}

void lower_k(character *c) {
#ifdef DEBUG
  printf("Func: %s [%s]\n", "lower_k", c->val);
#endif

  if (!(c->conv & CNV_LOWER_K)) {
    return;
  }
  uint8_t c1 = (uint8_t)c->val[1], c2 = (uint8_t)c->val[2];
  c->clen = 3;
  if (c1 == 0x80) {
    switch (c2) {
      case 0x81:
        strncpy(c->cval, "､", c->clen);
        break; /* 、 */
      case 0x82:
        strncpy(c->cval, "｡", c->clen);
        break; /* 。 */
    }
  } else if (c1 == 0x82) {
    switch (c2) {
      case 0x9b:
        strncpy(c->cval, "ﾞ", c->clen);
        break; /* ゛ */
      case 0x9c:
        strncpy(c->cval, "ﾟ", c->clen);
        break; /* ゜ */
      case 0xa1:
        strncpy(c->cval, "ｧ", c->clen);
        break;
      case 0xa2:
        strncpy(c->cval, "ｱ", c->clen);
        break;
      case 0xa3:
        strncpy(c->cval, "ｨ", c->clen);
        break;
      case 0xa4:
        strncpy(c->cval, "ｲ", c->clen);
        break;
      case 0xa5:
        strncpy(c->cval, "ｩ", c->clen);
        break;
      case 0xa6:
        strncpy(c->cval, "ｳ", c->clen);
        break;
      case 0xa7:
        strncpy(c->cval, "ｪ", c->clen);
        break;
      case 0xa8:
        strncpy(c->cval, "ｴ", c->clen);
        break;
      case 0xa9:
        strncpy(c->cval, "ｫ", c->clen);
        break;
      case 0xaa:
        strncpy(c->cval, "ｵ", c->clen);
        break;
      case 0xab:
        strncpy(c->cval, "ｶ", c->clen);
        break;
      case 0xac:
        c->clen = 6;
        strncpy(c->cval, "ｶﾞ", c->clen);
        break;
      case 0xad:
        strncpy(c->cval, "ｷ", c->clen);
        break;
      case 0xae:
        c->clen = 6;
        strncpy(c->cval, "ｷﾞ", c->clen);
        break;
      case 0xaf:
        strncpy(c->cval, "ｸ", c->clen);
        break;
      case 0xb0:
        c->clen = 6;
        strncpy(c->cval, "ｸﾞ", c->clen);
        break;
      case 0xb1:
        strncpy(c->cval, "ｹ", c->clen);
        break;
      case 0xb2:
        c->clen = 6;
        strncpy(c->cval, "ｹﾞ", c->clen);
        break;
      case 0xb3:
        strncpy(c->cval, "ｺ", c->clen);
        break;
      case 0xb4:
        c->clen = 6;
        strncpy(c->cval, "ｺﾞ", c->clen);
        break;
      case 0xb5:
        strncpy(c->cval, "ｻ", c->clen);
        break;
      case 0xb6:
        c->clen = 6;
        strncpy(c->cval, "ｻﾞ", c->clen);
        break;
      case 0xb7:
        strncpy(c->cval, "ｼ", c->clen);
        break;
      case 0xb8:
        c->clen = 6;
        strncpy(c->cval, "ｼﾞ", c->clen);
        break;
      case 0xb9:
        strncpy(c->cval, "ｽ", c->clen);
        break;
      case 0xba:
        c->clen = 6;
        strncpy(c->cval, "ｽﾞ", c->clen);
        break;
      case 0xbb:
        strncpy(c->cval, "ｾ", c->clen);
        break;
      case 0xbc:
        c->clen = 6;
        strncpy(c->cval, "ｾﾞ", c->clen);
        break;
      case 0xbd:
        strncpy(c->cval, "ｿ", c->clen);
        break;
      case 0xbe:
        c->clen = 6;
        strncpy(c->cval, "ｿﾞ", c->clen);
        break;
      case 0xbf:
        strncpy(c->cval, "ﾀ", c->clen);
        break;
    }
  } else if (c1 == 0x83) {
    switch (c2) {
      case 0x80:
        c->clen = 6;
        strncpy(c->cval, "ﾀﾞ", c->clen);
        break;
      case 0x81:
        strncpy(c->cval, "ﾁ", c->clen);
        break;
      case 0x82:
        c->clen = 6;
        strncpy(c->cval, "ﾁﾞ", c->clen);
        break;
      case 0x83:
        strncpy(c->cval, "ｯ", c->clen);
        break;
      case 0x84:
        strncpy(c->cval, "ﾂ", c->clen);
        break;
      case 0x85:
        c->clen = 6;
        strncpy(c->cval, "ﾂﾞ", c->clen);
        break;
      case 0x86:
        strncpy(c->cval, "ﾃ", c->clen);
        break;
      case 0x87:
        c->clen = 6;
        strncpy(c->cval, "ﾃﾞ", c->clen);
        break;
      case 0x88:
        strncpy(c->cval, "ﾄ", c->clen);
        break;
      case 0x89:
        c->clen = 6;
        strncpy(c->cval, "ﾄﾞ", c->clen);
        break;
      case 0x8a:
        strncpy(c->cval, "ﾅ", c->clen);
        break;
      case 0x8b:
        strncpy(c->cval, "ﾆ", c->clen);
        break;
      case 0x8c:
        strncpy(c->cval, "ﾇ", c->clen);
        break;
      case 0x8d:
        strncpy(c->cval, "ﾈ", c->clen);
        break;
      case 0x8e:
        strncpy(c->cval, "ﾉ", c->clen);
        break;
      case 0x8f:
        strncpy(c->cval, "ﾊ", c->clen);
        break;
      case 0x90:
        c->clen = 6;
        strncpy(c->cval, "ﾊﾞ", c->clen);
        break;
      case 0x91:
        c->clen = 6;
        strncpy(c->cval, "ﾊﾟ", c->clen);
        break;
      case 0x92:
        strncpy(c->cval, "ﾋ", c->clen);
        break;
      case 0x93:
        c->clen = 6;
        strncpy(c->cval, "ﾋﾞ", c->clen);
        break;
      case 0x94:
        c->clen = 6;
        strncpy(c->cval, "ﾋﾟ", c->clen);
        break;
      case 0x95:
        strncpy(c->cval, "ﾌ", c->clen);
        break;
      case 0x96:
        c->clen = 6;
        strncpy(c->cval, "ﾌﾞ", c->clen);
        break;
      case 0x97:
        c->clen = 6;
        strncpy(c->cval, "ﾌﾟ", c->clen);
        break;
      case 0x98:
        strncpy(c->cval, "ﾍ", c->clen);
        break;
      case 0x99:
        c->clen = 6;
        strncpy(c->cval, "ﾍﾞ", c->clen);
        break;
      case 0x9a:
        c->clen = 6;
        strncpy(c->cval, "ﾍﾟ", c->clen);
        break;
      case 0x9b:
        strncpy(c->cval, "ﾎ", c->clen);
        break;
      case 0x9c:
        c->clen = 6;
        strncpy(c->cval, "ﾎﾞ", c->clen);
        break;
      case 0x9d:
        c->clen = 6;
        strncpy(c->cval, "ﾎﾟ", c->clen);
        break;
      case 0x9e:
        strncpy(c->cval, "ﾏ", c->clen);
        break;
      case 0x9f:
        strncpy(c->cval, "ﾐ", c->clen);
        break;
      case 0xa0:
        strncpy(c->cval, "ﾑ", c->clen);
        break;
      case 0xa1:
        strncpy(c->cval, "ﾒ", c->clen);
        break;
      case 0xa2:
        strncpy(c->cval, "ﾓ", c->clen);
        break;
      case 0xa3:
        strncpy(c->cval, "ｬ", c->clen);
        break;
      case 0xa4:
        strncpy(c->cval, "ﾔ", c->clen);
        break;
      case 0xa5:
        strncpy(c->cval, "ｭ", c->clen);
        break;
      case 0xa6:
        strncpy(c->cval, "ﾕ", c->clen);
        break;
      case 0xa7:
        strncpy(c->cval, "ｮ", c->clen);
        break;
      case 0xa8:
        strncpy(c->cval, "ﾖ", c->clen);
        break;
      case 0xa9:
        strncpy(c->cval, "ﾗ", c->clen);
        break;
      case 0xaa:
        strncpy(c->cval, "ﾘ", c->clen);
        break;
      case 0xab:
        strncpy(c->cval, "ﾙ", c->clen);
        break;
      case 0xac:
        strncpy(c->cval, "ﾚ", c->clen);
        break;
      case 0xad:
        strncpy(c->cval, "ﾛ", c->clen);
        break;
      case 0xae:
        strncpy(c->cval, "ﾜ", c->clen);
        break;
      case 0xaf:
        strncpy(c->cval, "ﾜ", c->clen);
        break;
      case 0xb0:
        strncpy(c->cval, "ｲ", c->clen);
        break;
      case 0xb1:
        strncpy(c->cval, "ｴ", c->clen);
        break;
      case 0xb2:
        strncpy(c->cval, "ｦ", c->clen);
        break;
      case 0xb3:
        strncpy(c->cval, "ﾝ", c->clen);
        break;
      case 0xb4:
        c->clen = 6;
        strncpy(c->cval, "ｳﾞ", c->clen);
        break;
      case 0xbb:
        strncpy(c->cval, "･", c->clen);
        break; /* ・ */
      case 0xbc:
        strncpy(c->cval, "ｰ", c->clen);
        break; /* ー */
    }
  }
  *(c->cval + c->clen) = 0x00;
}

void upper_k(character *c) {
#ifdef DEBUG
  printf("Func: %s [%s]\n", "upper_k", c->val);
#endif

  if (!(c->conv & CNV_UPPER_K)) {
    return;
  }
  uint8_t c1 = *(c->val + 1) & 0xff, c2 = *(c->val + 2) & 0xff, c5 = 0x00;
  if (c->len > 5) {
    c5 = *(c->val + 5) & 0xff;
  }
  c->clen = 3;
  if (c1 == 0xbd) {
    switch (c2) {
      case 0xa1: /* ｡ */
        strncpy(c->cval, "。", c->clen);
        break;
      case 0xa2: /* ｢ */
        strncpy(c->cval, "「", c->clen);
        break;
      case 0xa3: /* ｣ */
        strncpy(c->cval, "」", c->clen);
        break;
      case 0xa4: /* ､ */
        strncpy(c->cval, "、", c->clen);
        break;
      case 0xa5: /* ･ */
        strncpy(c->cval, "・", c->clen);
        break;
      case 0xa6:
        strncpy(c->cval, "ヲ", c->clen);
        break;
      case 0xa7:
        strncpy(c->cval, "ァ", c->clen);
        break;
      case 0xa8:
        strncpy(c->cval, "ィ", c->clen);
        break;
      case 0xa9:
        strncpy(c->cval, "ゥ", c->clen);
        break;
      case 0xaa:
        strncpy(c->cval, "ェ", c->clen);
        break;
      case 0xab:
        strncpy(c->cval, "ォ", c->clen);
        break;
      case 0xac:
        strncpy(c->cval, "ャ", c->clen);
        break;
      case 0xad:
        strncpy(c->cval, "ュ", c->clen);
        break;
      case 0xae:
        strncpy(c->cval, "ョ", c->clen);
        break;
      case 0xaf:
        strncpy(c->cval, "ッ", c->clen);
        break;
      case 0xb0: /* ｰ */
        strncpy(c->cval, "ー", c->clen);
        break;
      case 0xb1:
        strncpy(c->cval, "ア", c->clen);
        break;
      case 0xb2:
        strncpy(c->cval, "イ", c->clen);
        break;
      case 0xb3:
        strncpy(c->cval, (c5 == 0x9e ? "ヴ" : "ウ"), c->clen);
        break;
      case 0xb4:
        strncpy(c->cval, "エ", c->clen);
        break;
      case 0xb5:
        strncpy(c->cval, "オ", c->clen);
        break;
      case 0xb6:
        strncpy(c->cval, (c5 == 0x9e ? "ガ" : "カ"), c->clen);
        break;
      case 0xb7:
        strncpy(c->cval, (c5 == 0x9e ? "ギ" : "キ"), c->clen);
        break;
      case 0xb8:
        strncpy(c->cval, (c5 == 0x9e ? "グ" : "ク"), c->clen);
        break;
      case 0xb9:
        strncpy(c->cval, (c5 == 0x9e ? "ゲ" : "ケ"), c->clen);
        break;
      case 0xba:
        strncpy(c->cval, (c5 == 0x9e ? "ゴ" : "コ"), c->clen);
        break;
      case 0xbb:
        strncpy(c->cval, (c5 == 0x9e ? "ザ" : "サ"), c->clen);
        break;
      case 0xbc:
        strncpy(c->cval, (c5 == 0x9e ? "ジ" : "シ"), c->clen);
        break;
      case 0xbd:
        strncpy(c->cval, (c5 == 0x9e ? "ズ" : "ス"), c->clen);
        break;
      case 0xbe:
        strncpy(c->cval, (c5 == 0x9e ? "ゼ" : "セ"), c->clen);
        break;
      case 0xbf:
        strncpy(c->cval, (c5 == 0x9e ? "ゾ" : "ソ"), c->clen);
        break;
    }
  } else if (c1 == 0xbe) {
    switch (c2) {
      case 0x80:
        strncpy(c->cval, (c5 == 0x9e ? "ダ" : "タ"), c->clen);
        break;
      case 0x81:
        strncpy(c->cval, (c5 == 0x9e ? "ヂ" : "チ"), c->clen);
        break;
      case 0x82:
        strncpy(c->cval, (c5 == 0x9e ? "ヅ" : "ツ"), c->clen);
        break;
      case 0x83:
        strncpy(c->cval, (c5 == 0x9e ? "デ" : "テ"), c->clen);
        break;
      case 0x84:
        strncpy(c->cval, (c5 == 0x9e ? "ド" : "ト"), c->clen);
        break;
      case 0x85:
        strncpy(c->cval, "ナ", c->clen);
        break;
      case 0x86:
        strncpy(c->cval, "ニ", c->clen);
        break;
      case 0x87:
        strncpy(c->cval, "ヌ", c->clen);
        break;
      case 0x88:
        strncpy(c->cval, "ネ", c->clen);
        break;
      case 0x89:
        strncpy(c->cval, "ノ", c->clen);
        break;
      case 0x8a:
        if (c5 == 0x9e) {
          strncpy(c->cval, "バ", c->clen);
        } else if (c5 == 0x9f) {
          strncpy(c->cval, "パ", c->clen);
        } else {
          strncpy(c->cval, "ハ", c->clen);
        }
        break;
      case 0x8b:
        if (c5 == 0x9e) {
          strncpy(c->cval, "ビ", c->clen);
        } else if (c5 == 0x9f) {
          strncpy(c->cval, "ピ", c->clen);
        } else {
          strncpy(c->cval, "ヒ", c->clen);
        }
        break;
      case 0x8c:
        if (c5 == 0x9e) {
          strncpy(c->cval, "ブ", c->clen);
        } else if (c5 == 0x9f) {
          strncpy(c->cval, "プ", c->clen);
        } else {
          strncpy(c->cval, "フ", c->clen);
        }
        break;
      case 0x8d:
        if (c5 == 0x9e) {
          strncpy(c->cval, "ベ", c->clen);
        } else if (c5 == 0x9f) {
          strncpy(c->cval, "ペ", c->clen);
        } else {
          strncpy(c->cval, "ヘ", c->clen);
        }
        break;
      case 0x8e:
        if (c5 == 0x9e) {
          strncpy(c->cval, "ボ", c->clen);
        } else if (c5 == 0x9f) {
          strncpy(c->cval, "ポ", c->clen);
        } else {
          strncpy(c->cval, "ホ", c->clen);
        }
        break;
      case 0x8f:
        strncpy(c->cval, "マ", c->clen);
        break;
      case 0x90:
        strncpy(c->cval, "ミ", c->clen);
        break;
      case 0x91:
        strncpy(c->cval, "ム", c->clen);
        break;
      case 0x92:
        strncpy(c->cval, "メ", c->clen);
        break;
      case 0x93:
        strncpy(c->cval, "モ", c->clen);
        break;
      case 0x94:
        strncpy(c->cval, "ヤ", c->clen);
        break;
      case 0x95:
        strncpy(c->cval, "ユ", c->clen);
        break;
      case 0x96:
        strncpy(c->cval, "ヨ", c->clen);
        break;
      case 0x97:
        strncpy(c->cval, "ラ", c->clen);
        break;
      case 0x98:
        strncpy(c->cval, "リ", c->clen);
        break;
      case 0x99:
        strncpy(c->cval, "ル", c->clen);
        break;
      case 0x9a:
        strncpy(c->cval, "レ", c->clen);
        break;
      case 0x9b:
        strncpy(c->cval, "ロ", c->clen);
        break;
      case 0x9c:
        strncpy(c->cval, "ワ", c->clen);
        break;
      case 0x9d:
        strncpy(c->cval, "ン", c->clen);
        break;
      case 0x9e: /* ﾞ */
        strncpy(c->cval, "゛", c->clen);
        break;
      case 0x9f: /* ﾟ */
        strncpy(c->cval, "゜", c->clen);
        break;
    }
  }
  *(c->cval + c->clen) = 0x00;
}

void lower_h(character *c) {
#ifdef DEBUG
  printf("Func: %s [%s]\n", "lower_h", c->val);
#endif

  if (!(c->conv & CNV_LOWER_H)) {
    return;
  }
  uint8_t c1 = (uint8_t)(*(c->val + 1) & 0xff),
          c2 = (uint8_t)(*(c->val + 2) & 0xff);
  c->clen = 3;
  if (c1 == 0x80) {
    switch (c2) {
      case 0x81: /* 、 */
        strncpy(c->cval, "､", c->clen);
        break;
      case 0x82: /* 。 */
        strncpy(c->cval, "｡", c->clen);
        break;
    }
  } else if (c1 == 0x81) {
    switch (c2) {
      case 0x81:
        strncpy(c->cval, "ｧ", c->clen);
        break;
      case 0x82:
        strncpy(c->cval, "ｱ", c->clen);
        break;
      case 0x83:
        strncpy(c->cval, "ｨ", c->clen);
        break;
      case 0x84:
        strncpy(c->cval, "ｲ", c->clen);
        break;
      case 0x85:
        strncpy(c->cval, "ｩ", c->clen);
        break;
      case 0x86:
        strncpy(c->cval, "ｳ", c->clen);
        break;
      case 0x87:
        strncpy(c->cval, "ｪ", c->clen);
        break;
      case 0x88:
        strncpy(c->cval, "ｴ", c->clen);
        break;
      case 0x89:
        strncpy(c->cval, "ｫ", c->clen);
        break;
      case 0x8a:
        strncpy(c->cval, "ｵ", c->clen);
        break;
      case 0x8b:
        strncpy(c->cval, "ｶ", c->clen);
        break;
      case 0x8c:
        c->clen = 6;
        strncpy(c->cval, "ｶﾞ", c->clen);
        break;
      case 0x8d:
        strncpy(c->cval, "ｷ", c->clen);
        break;
      case 0x8e:
        c->clen = 6;
        strncpy(c->cval, "ｷﾞ", c->clen);
        break;
      case 0x8f:
        strncpy(c->cval, "ｸ", c->clen);
        break;
      case 0x90:
        c->clen = 6;
        strncpy(c->cval, "ｸﾞ", c->clen);
        break;
      case 0x91:
        strncpy(c->cval, "ｹ", c->clen);
        break;
      case 0x92:
        c->clen = 6;
        strncpy(c->cval, "ｹﾞ", c->clen);
        break;
      case 0x93:
        strncpy(c->cval, "ｺ", c->clen);
        break;
      case 0x94:
        c->clen = 6;
        strncpy(c->cval, "ｺﾞ", c->clen);
        break;
      case 0x95:
        strncpy(c->cval, "ｻ", c->clen);
        break;
      case 0x96:
        c->clen = 6;
        strncpy(c->cval, "ｻﾞ", c->clen);
        break;
      case 0x97:
        strncpy(c->cval, "ｼ", c->clen);
        break;
      case 0x98:
        c->clen = 6;
        strncpy(c->cval, "ｼﾞ", c->clen);
        break;
      case 0x99:
        strncpy(c->cval, "ｽ", c->clen);
        break;
      case 0x9a:
        c->clen = 6;
        strncpy(c->cval, "ｽﾞ", c->clen);
        break;
      case 0x9b:
        strncpy(c->cval, "ｾ", c->clen);
        break;
      case 0x9c:
        c->clen = 6;
        strncpy(c->cval, "ｾﾞ", c->clen);
        break;
      case 0x9d:
        strncpy(c->cval, "ｿ", c->clen);
        break;
      case 0x9e:
        c->clen = 6;
        strncpy(c->cval, "ｿﾞ", c->clen);
        break;
      case 0x9f:
        strncpy(c->cval, "ﾀ", c->clen);
        break;
      case 0xa0:
        c->clen = 6;
        strncpy(c->cval, "ﾀﾞ", c->clen);
        break;
      case 0xa1:
        strncpy(c->cval, "ﾁ", c->clen);
        break;
      case 0xa2:
        c->clen = 6;
        strncpy(c->cval, "ﾁﾞ", c->clen);
        break;
      case 0xa3:
        strncpy(c->cval, "ｯ", c->clen);
        break;
      case 0xa4:
        strncpy(c->cval, "ﾂ", c->clen);
        break;
      case 0xa5:
        c->clen = 6;
        strncpy(c->cval, "ﾂﾞ", c->clen);
        break;
      case 0xa6:
        strncpy(c->cval, "ﾃ", c->clen);
        break;
      case 0xa7:
        c->clen = 6;
        strncpy(c->cval, "ﾃﾞ", c->clen);
        break;
      case 0xa8:
        strncpy(c->cval, "ﾄ", c->clen);
        break;
      case 0xa9:
        c->clen = 6;
        strncpy(c->cval, "ﾄﾞ", c->clen);
        break;
      case 0xaa:
        strncpy(c->cval, "ﾅ", c->clen);
        break;
      case 0xab:
        strncpy(c->cval, "ﾆ", c->clen);
        break;
      case 0xac:
        strncpy(c->cval, "ﾇ", c->clen);
        break;
      case 0xad:
        strncpy(c->cval, "ﾈ", c->clen);
        break;
      case 0xae:
        strncpy(c->cval, "ﾉ", c->clen);
        break;
      case 0xaf:
        strncpy(c->cval, "ﾊ", c->clen);
        break;
      case 0xb0:
        c->clen = 6;
        strncpy(c->cval, "ﾊﾞ", c->clen);
        break;
      case 0xb1:
        c->clen = 6;
        strncpy(c->cval, "ﾊﾟ", c->clen);
        break;
      case 0xb2:
        strncpy(c->cval, "ﾋ", c->clen);
        break;
      case 0xb3:
        c->clen = 6;
        strncpy(c->cval, "ﾋﾞ", c->clen);
        break;
      case 0xb4:
        c->clen = 6;
        strncpy(c->cval, "ﾋﾟ", c->clen);
        break;
      case 0xb5:
        strncpy(c->cval, "ﾌ", c->clen);
        break;
      case 0xb6:
        c->clen = 6;
        strncpy(c->cval, "ﾌﾞ", c->clen);
        break;
      case 0xb7:
        c->clen = 6;
        strncpy(c->cval, "ﾌﾟ", c->clen);
        break;
      case 0xb8:
        strncpy(c->cval, "ﾍ", c->clen);
        break;
      case 0xb9:
        c->clen = 6;
        strncpy(c->cval, "ﾍﾞ", c->clen);
        break;
      case 0xba:
        c->clen = 6;
        strncpy(c->cval, "ﾍﾟ", c->clen);
        break;
      case 0xbb:
        strncpy(c->cval, "ﾎ", c->clen);
        break;
      case 0xbc:
        c->clen = 6;
        strncpy(c->cval, "ﾎﾞ", c->clen);
        break;
      case 0xbd:
        c->clen = 6;
        strncpy(c->cval, "ﾎﾟ", c->clen);
        break;
      case 0xbe:
        strncpy(c->cval, "ﾏ", c->clen);
        break;
      case 0xbf:
        strncpy(c->cval, "ﾐ", c->clen);
        break;
    }
  } else if (c1 == 0x82) {
    switch (c2) {
      case 0x80:
        strncpy(c->cval, "ﾑ", c->clen);
        break;
      case 0x81:
        strncpy(c->cval, "ﾒ", c->clen);
        break;
      case 0x82:
        strncpy(c->cval, "ﾓ", c->clen);
        break;
      case 0x83:
        strncpy(c->cval, "ｬ", c->clen);
        break;
      case 0x84:
        strncpy(c->cval, "ﾔ", c->clen);
        break;
      case 0x85:
        strncpy(c->cval, "ｭ", c->clen);
        break;
      case 0x86:
        strncpy(c->cval, "ﾕ", c->clen);
        break;
      case 0x87:
        strncpy(c->cval, "ｮ", c->clen);
        break;
      case 0x88:
        strncpy(c->cval, "ﾖ", c->clen);
        break;
      case 0x89:
        strncpy(c->cval, "ﾗ", c->clen);
        break;
      case 0x8a:
        strncpy(c->cval, "ﾘ", c->clen);
        break;
      case 0x8b:
        strncpy(c->cval, "ﾙ", c->clen);
        break;
      case 0x8c:
        strncpy(c->cval, "ﾚ", c->clen);
        break;
      case 0x8d:
        strncpy(c->cval, "ﾛ", c->clen);
        break;
      case 0x8e:
        strncpy(c->cval, "ﾜ", c->clen);
        break;
      case 0x8f:
        strncpy(c->cval, "ﾜ", c->clen);
        break;
      case 0x90:
        strncpy(c->cval, "ｲ", c->clen);
        break;
      case 0x91:
        strncpy(c->cval, "ｴ", c->clen);
        break;
      case 0x92:
        strncpy(c->cval, "ｦ", c->clen);
        break;
      case 0x93:
        strncpy(c->cval, "ﾝ", c->clen);
        break;
      case 0x9b: /* ゛ */
        strncpy(c->cval, "ﾞ", c->clen);
        break;
      case 0x9c: /* ゜ */
        strncpy(c->cval, "ﾟ", c->clen);
        break;
    }
  } else if (c1 == 0x83) {
    switch (c2) {
      case 0xbb: /* ・ */
        strncpy(c->cval, "･", c->clen);
        break;
      case 0xbc: /* ー */
        strncpy(c->cval, "ｰ", c->clen);
        break;
    }
  }
  *(c->cval + c->clen) = 0x00;
}

void upper_h(character *c) {
#ifdef DEBUG
  printf("Func: %s [%s]\n", "upper_h", c->val);
#endif

  if (!(c->conv & CNV_UPPER_H)) {
    return;
  }
  uint8_t c1 = *(c->val + 1) & 0xff, c2 = *(c->val + 2) & 0xff, c5 = 0x00;
  if (c->len > 5) {
    c5 = *(c->val + 5) & 0xff;
  }

  c->clen = 3;
  if (c1 == 0xbd) {
    switch (c2) {
      case 0xa1: /* ｡ */
        strncpy(c->cval, "。", c->clen);
        break;
      case 0xa2: /* ｢ */
        strncpy(c->cval, "「", c->clen);
        break;
      case 0xa3: /* ｣ */
        strncpy(c->cval, "」", c->clen);
        break;
      case 0xa4: /* ､ */
        strncpy(c->cval, "、", c->clen);
        break;
      case 0xa5: /* ･ */
        strncpy(c->cval, "・", c->clen);
        break;
      case 0xa6:
        strncpy(c->cval, "を", c->clen);
        break;
      case 0xa7:
        strncpy(c->cval, "ぁ", c->clen);
        break;
      case 0xa8:
        strncpy(c->cval, "ぃ", c->clen);
        break;
      case 0xa9:
        strncpy(c->cval, "ぅ", c->clen);
        break;
      case 0xaa:
        strncpy(c->cval, "ぇ", c->clen);
        break;
      case 0xab:
        strncpy(c->cval, "ぉ", c->clen);
        break;
      case 0xac:
        strncpy(c->cval, "ゃ", c->clen);
        break;
      case 0xad:
        strncpy(c->cval, "ゅ", c->clen);
        break;
      case 0xae:
        strncpy(c->cval, "ょ", c->clen);
        break;
      case 0xaf:
        strncpy(c->cval, "っ", c->clen);
        break;
      case 0xb0: /* ｰ */
        strncpy(c->cval, "ー", c->clen);
        break;
      case 0xb1:
        strncpy(c->cval, "あ", c->clen);
        break;
      case 0xb2:
        strncpy(c->cval, "い", c->clen);
        break;
      case 0xb3:
        if (is_voiced(c->cval, c->clen)) {
          c->clen = 6;
          strncpy(c->cval, "う゛", c->clen);
        } else {
          strncpy(c->cval, "う", c->clen);
        }
        break;
      case 0xb4:
        strncpy(c->cval, "え", c->clen);
        break;
      case 0xb5:
        strncpy(c->cval, "お", c->clen);
        break;
      case 0xb6:
        strncpy(c->cval, (c5 == 0x9e ? "が" : "か"), c->clen);
        break;
      case 0xb7:
        strncpy(c->cval, (c5 == 0x9e ? "ぎ" : "き"), c->clen);
        break;
      case 0xb8:
        strncpy(c->cval, (c5 == 0x9e ? "ぐ" : "く"), c->clen);
        break;
      case 0xb9:
        strncpy(c->cval, (c5 == 0x9e ? "げ" : "け"), c->clen);
        break;
      case 0xba:
        strncpy(c->cval, (c5 == 0x9e ? "ご" : "こ"), c->clen);
        break;
      case 0xbb:
        strncpy(c->cval, (c5 == 0x9e ? "ざ" : "さ"), c->clen);
        break;
      case 0xbc:
        strncpy(c->cval, (c5 == 0x9e ? "じ" : "し"), c->clen);
        break;
      case 0xbd:
        strncpy(c->cval, (c5 == 0x9e ? "ず" : "す"), c->clen);
        break;
      case 0xbe:
        strncpy(c->cval, (c5 == 0x9e ? "ぜ" : "せ"), c->clen);
        break;
      case 0xbf:
        strncpy(c->cval, (c5 == 0x9e ? "ぞ" : "そ"), c->clen);
        break;
    }
  } else if (c1 == 0xbe) {
    switch (c2) {
      case 0x80:
        strncpy(c->cval, (c5 == 0x9e ? "だ" : "た"), c->clen);
        break;
      case 0x81:
        strncpy(c->cval, (c5 == 0x9e ? "ぢ" : "ち"), c->clen);
        break;
      case 0x82:
        strncpy(c->cval, (c5 == 0x9e ? "づ" : "つ"), c->clen);
        break;
      case 0x83:
        strncpy(c->cval, (c5 == 0x9e ? "で" : "て"), c->clen);
        break;
      case 0x84:
        strncpy(c->cval, (c5 == 0x9e ? "ど" : "と"), c->clen);
        break;
      case 0x85:
        strncpy(c->cval, "な", c->clen);
        break;
      case 0x86:
        strncpy(c->cval, "に", c->clen);
        break;
      case 0x87:
        strncpy(c->cval, "ぬ", c->clen);
        break;
      case 0x88:
        strncpy(c->cval, "ね", c->clen);
        break;
      case 0x89:
        strncpy(c->cval, "の", c->clen);
        break;
      case 0x8a:
        if (c5 == 0x9e) {
          strncpy(c->cval, "ば", c->clen);
        } else if (c5 == 0x9f) {
          strncpy(c->cval, "ぱ", c->clen);
        } else {
          strncpy(c->cval, "は", c->clen);
        }
        break;
      case 0x8b:
        if (c5 == 0x9e) {
          strncpy(c->cval, "び", c->clen);
        } else if (c5 == 0x9f) {
          strncpy(c->cval, "ぴ", c->clen);
        } else {
          strncpy(c->cval, "ひ", c->clen);
        }
        break;
      case 0x8c:
        if (c5 == 0x9e) {
          strncpy(c->cval, "ぶ", c->clen);
        } else if (c5 == 0x9f) {
          strncpy(c->cval, "ぷ", c->clen);
        } else {
          strncpy(c->cval, "ふ", c->clen);
        }
        break;
      case 0x8d:
        if (c5 == 0x9e) {
          strncpy(c->cval, "べ", c->clen);
        } else if (c5 == 0x9f) {
          strncpy(c->cval, "ぺ", c->clen);
        } else {
          strncpy(c->cval, "へ", c->clen);
        }
        break;
      case 0x8e:
        if (c5 == 0x9e) {
          strncpy(c->cval, "ぼ", c->clen);
        } else if (c5 == 0x9f) {
          strncpy(c->cval, "ぽ", c->clen);
        } else {
          strncpy(c->cval, "ほ", c->clen);
        }
        break;
      case 0x8f:
        strncpy(c->cval, "ま", c->clen);
        break;
      case 0x90:
        strncpy(c->cval, "み", c->clen);
        break;
      case 0x91:
        strncpy(c->cval, "む", c->clen);
        break;
      case 0x92:
        strncpy(c->cval, "め", c->clen);
        break;
      case 0x93:
        strncpy(c->cval, "も", c->clen);
        break;
      case 0x94:
        strncpy(c->cval, "や", c->clen);
        break;
      case 0x95:
        strncpy(c->cval, "ゆ", c->clen);
        break;
      case 0x96:
        strncpy(c->cval, "よ", c->clen);
        break;
      case 0x97:
        strncpy(c->cval, "ら", c->clen);
        break;
      case 0x98:
        strncpy(c->cval, "り", c->clen);
        break;
      case 0x99:
        strncpy(c->cval, "る", c->clen);
        break;
      case 0x9a:
        strncpy(c->cval, "れ", c->clen);
        break;
      case 0x9b:
        strncpy(c->cval, "ろ", c->clen);
        break;
      case 0x9c:
        strncpy(c->cval, "わ", c->clen);
        break;
      case 0x9d:
        strncpy(c->cval, "ん", c->clen);
        break;
      case 0x9e: /* ﾞ */
        strncpy(c->cval, "゛", c->clen);
        break;
      case 0x9f: /* ﾟ */
        strncpy(c->cval, "゜", c->clen);
        break;
    }
  }
  *(c->cval + c->clen) = 0x00;
}

void lower_c(character *c) {
#ifdef DEBUG
  printf("Func: %s [%s]\n", "lower_c", c->val);
#endif

  if (!(c->conv & CNV_LOWER_C)) {
    return;
  }

  uint8_t c0 = (uint8_t)c->val[0], c1 = (uint8_t)c->val[1],
          c2 = (uint8_t)c->val[2];
  c->clen = 3;
  switch (c1) {
    case 0x82:  // ァ - タ
      if (c2 >= 0xa1 && c2 <= 0xbf) {
        *(c->cval + 0) = (char)(0xe3 & 0xff);
        *(c->cval + 1) = (char)(0x81 & 0xff);
        *(c->cval + 2) = (char)((c2 - 0x20) & 0xff);
        *(c->cval + c->clen) = 0x00;
      }
      break;
    case 0x83:
      if (c2 >= 0x80 && c2 <= 0x9f) {  // ダ - ミ
        *(c->cval + 0) = (char)(0xe3 & 0xff);
        *(c->cval + 1) = (char)(0x81 & 0xff);
        *(c->cval + 2) = (char)((c2 + 0x20) & 0xff);
        *(c->cval + c->clen) = 0x00;
      } else if (c2 >= 0xa0 && c2 <= 0xb3) {  // ム - ン
        *(c->cval + 0) = (char)(0xe3 & 0xff);
        *(c->cval + 1) = (char)(0x82 & 0xff);
        *(c->cval + 2) = (char)((c2 - 0x20) & 0xff);
        *(c->cval + c->clen) = 0x00;
      } else if (c2 >= 0xbd && c2 <= 0xbe) {  // ヽヾ
        *(c->cval + 0) = (char)(0xe3 & 0xff);
        *(c->cval + 1) = (char)(0x82 & 0xff);
        *(c->cval + 2) = (char)((c2 - 0x20) & 0xff);
        *(c->cval + c->clen) = 0x00;
      }
      break;
  }
}

void upper_c(character *c) {
#ifdef DEBUG
  printf("Func: %s [%s]\n", "upper_c", c->val);
#endif

  if (!(c->conv & CNV_UPPER_C)) {
    return;
  }
  uint8_t c0 = (uint8_t)c->val[0], c1 = (uint8_t)c->val[1],
          c2 = (uint8_t)c->val[2];
  c->clen = 3;
  switch (c1) {
    case 0x81:
      if (c2 >= 0x81 && c2 <= 0x9f) {  // ぁ - た
        *(c->cval + 0) = (char)(0xe3 & 0xff);
        *(c->cval + 1) = (char)(0x82 & 0xff);
        *(c->cval + 2) = (char)((c2 + 0x20) & 0xff);
        *(c->cval + c->clen) = 0x00;
      } else if (c2 >= 0xa0 && c2 <= 0xbf) {  // だ - み
        *(c->cval + 0) = (char)(0xe3 & 0xff);
        *(c->cval + 1) = (char)(0x83 & 0xff);
        *(c->cval + 2) = (char)((c2 - 0x20) & 0xff);
        *(c->cval + c->clen) = 0x00;
      }
      break;
    case 0x82:
      if (c2 >= 0x80 && c2 <= 0x93) {  // む - ん
        *(c->cval + 0) = (char)(0xe3 & 0xff);
        *(c->cval + 1) = (char)(0x83 & 0xff);
        *(c->cval + 2) = (char)((c2 + 0x20) & 0xff);
        *(c->cval + c->clen) = 0x00;
      } else if (c2 >= 0x9d && c2 <= 0x9e) {  // ゝゞ
        *(c->cval + 0) = (char)(0xe3 & 0xff);
        *(c->cval + 1) = (char)(0x83 & 0xff);
        *(c->cval + 2) = (char)((c2 + 0x20) & 0xff);
        *(c->cval + c->clen) = 0x00;
      }
      break;
  }
}

void asis(character *c) {
#ifdef DEBUG
  printf("Func: %s [%s]\n", "asis", c->val);
#endif

  for (uint8_t i = 0; i < c->len; i++) {
    *(c->cval + i) = (char)(*(c->val + i) & 0xff);
  }
  c->clen = c->len;
  *(c->cval + c->clen) = 0x00;
}

void extract(character *c, const char *s, int len) {
#ifdef DEBUG
  printf("Func: %s\n", "extract");
#endif
  bool flg = false;
  c->len = 1;
  if (is_1byte(s, len)) {
    uint8_t c0 = *s & 0xff;
    if (c0 == 0x20) {  // Space
      c->conv = CNV_UPPER_S;
    } else if (c0 >= 0x30 && c0 <= 0x39) {  // 0 - 9
      c->conv = CNV_UPPER_A | CNV_UPPER_N;
    } else if (c0 >= 0x41 && c0 <= 0x5a) {  // A - Z
      c->conv = CNV_UPPER_A | CNV_UPPER_R;
    } else if (c0 >= 0x61 && c0 <= 0x7a) {  // a - z
      c->conv = CNV_UPPER_A | CNV_UPPER_R;
    } else if (c0 >= 0x21 && c0 <= 0x7d && c0 != 0x22 && c0 != 0x27 &&
               c0 != 0x5c) {
      c->conv = CNV_UPPER_A;
    }
  } else if (is_3bytes(s, len)) {
    uint8_t c0 = *s & 0xff, c1 = *(s + 1) & 0xff, c2 = *(s + 2) & 0xff;
    c->len = 3;
    if (c0 == 0xef) {
      if (c1 == 0xbc) {
        if (c2 >= 0x90 && c2 <= 0x99) {  // ０ - ９
          flg = true;
          c->conv = CNV_LOWER_A | CNV_LOWER_N;
        } else if (c2 >= 0xa1 && c2 <= 0xba) {  // Ａ - Ｚ
          c->conv = CNV_LOWER_A | CNV_LOWER_R;
        } else if (c2 != 0x82 && c2 != 0x87 && c2 != 0xbc) {  // except ＂ ＇ ＼
          c->conv = CNV_LOWER_A;
        }
      } else if (c1 == 0xbd) {
        if (c2 >= 0x81 && c2 <= 0x9a) {  // ａ - ｚ
          c->conv = CNV_LOWER_A | CNV_LOWER_R;
        } else if (c2 >= 0x80 && c2 <= 0x9d) {  // ｀ - ｝
          c->conv = CNV_LOWER_A;
        } else if (c2 >= 0xa1 && c2 <= 0xbf) {  // ｡ ｢ ｣ ､ ･ ｦ - ｿ
          if (is_voiced(s, len)) {              // voiced
            c->len = 6;
            if (c2 == 0xb3) {
              c->conv = CNV_UPPER_H | CNV_UPPER_K;
            } else {
              c->conv = CNV_UPPER_H | CNV_UPPER_K;
            }
          } else {
            c->conv = CNV_UPPER_H | CNV_UPPER_K;
          }
        }
      } else if (c1 == 0xbe) {
        if (c2 >= 0x80 && c2 <= 0x84) {  // ﾀ - ﾄ
          if (is_voiced(s, len)) {
            c->len = 6;
          }
          c->conv = CNV_UPPER_H | CNV_UPPER_K;
        } else if (c2 >= 0x8a && c2 <= 0x8e) {  // ﾊ - ﾎ
          if (is_voiced(s, len) ||
              is_semi_voiced(s, len)) {  // voiced or semi voiced
            c->len = 6;
          }
          c->conv = CNV_UPPER_H | CNV_UPPER_K;
        } else if (c2 >= 0x85 && c2 <= 0x9f) {  // ﾅ - ﾝﾞﾟ
          c->conv = CNV_UPPER_H | CNV_UPPER_K;
        }
      }
    } else if (c0 == 0xe3) {
      if (c1 == 0x80) {
        if (c2 == 0x80) {  // Space
          c->conv = CNV_LOWER_S;
        } else if (c2 >= 0x81 && c2 <= 0x82) {  // 、。
          c->conv = CNV_LOWER_H | CNV_LOWER_K;
        }
      } else if (c1 == 0x81) {
        if (c2 >= 0x81 && c2 <= 0xbf) {  // ぁ - み
          c->conv = CNV_UPPER_C | CNV_LOWER_H;
        }
      } else if (c1 == 0x82) {
        if (c2 >= 0x80 && c2 <= 0x93) {  // む - ん
          c->conv = CNV_UPPER_C | CNV_LOWER_H;
        } else if (c2 >= 0x9b && c2 <= 0x9c) {  // ゛゜
          c->conv = CNV_LOWER_H | CNV_LOWER_K;
        } else if (c2 >= 0x9d && c2 <= 0x9e) {  // ゝゞ
          c->conv = CNV_UPPER_C;
        } else if (c2 >= 0xa1 && c2 <= 0xbf) {  // ァ - タ
          c->conv = CNV_LOWER_C | CNV_LOWER_K;
        }
      } else if (c1 == 0x83) {
        if (c2 >= 0x80 && c2 <= 0xb3) {  // チ - ン
          c->conv = CNV_LOWER_C | CNV_LOWER_K;
        } else if (c2 == 0xb4) {  // ヴ
          c->conv = CNV_LOWER_K;
        } else if (c2 >= 0xbb && c2 <= 0xbc) {  // ・ー
          c->conv = CNV_LOWER_H | CNV_LOWER_K;
        } else if (c2 >= 0xbd && c2 <= 0xbe) {  // ヽ ヾ
          c->conv = CNV_LOWER_C;
        }
      }
    }
  } else if (is_2bytes(s, len)) {
    c->len = 2;
  } else if (is_4bytes(s, len)) {
    c->len = 4;
  } else {
    c->len = 1;
  }

  for (int i = 0; i < c->len; i++) {
    *(c->val + i) = *(s + i) & 0xff;
  }
  *(c->val + c->len) = 0x00;
}

void conv(character *c, filter *filters) {
#ifdef DEBUG
  printf("Func: %s\n", "conv");
#endif

  int i = 0;
  while (i < MODE_BUF && *(filters + i) != 0x00) {
    (*(filters + i))(c);
    if (c->clen > 0) {
      break;
    }
    i++;
  }
  if (c->clen == 0) {
    asis(c);
  }
}

filter *create_filters(const char *mode_str, int mode_len) {
#ifdef DEBUG
  printf("Func: %s\n", "create_filters");
#endif

  filter *filters = (filter *)calloc(MODE_BUF, sizeof(filter));
  if (filters == NULL) {
    return NULL;
  }

  char m;
  char mode[MODE_BUF] = {};
  bool flg;
  int k = 0;
  for (int i = 0; i < mode_len; i++) {
    m = *(mode_str + i);
    flg = false;
    for (int j = 0; j < MODE_BUF; j++) {
      if (mode[j] == m) {
        flg = false;
        break;
      }
      flg = true;
    }

    if (flg) {
      mode[k] = m;
      switch (m) {
        case 'r':
          *(filters + k) = lower_r;
          k++;
          break;
        case 'R':
          *(filters + k) = upper_r;
          k++;
          break;
        case 'n':
          *(filters + k) = lower_n;
          k++;
          break;
        case 'N':
          *(filters + k) = upper_n;
          k++;
          break;
        case 'a':
          *(filters + k) = lower_a;
          k++;
          break;
        case 'A':
          *(filters + k) = upper_a;
          k++;
          break;
        case 's':
          *(filters + k) = lower_s;
          k++;
          break;
        case 'S':
          *(filters + k) = upper_s;
          k++;
          break;
        case 'k':
          *(filters + k) = lower_k;
          k++;
          break;
        case 'K':
          *(filters + k) = upper_k;
          k++;
          break;
        case 'h':
          *(filters + k) = lower_h;
          k++;
          break;
        case 'H':
          *(filters + k) = upper_h;
          k++;
          break;
        case 'c':
          *(filters + k) = lower_c;
          k++;
          break;
        case 'C':
          *(filters + k) = upper_c;
          k++;
          break;
      }
    }
  }
  return filters;
}

void init_character(character *c) {
#ifdef DEBUG
  printf("Func: %s\n", "init_character");
#endif

  c->val[0] = 0x00;
  c->len = 0;
  c->conv = CNV_ASIS;
  c->cval[0] = 0x00;
  c->clen = 0;
}

char *convert(const char *str, int str_len, const char *mode_str,
              int mode_len) {
#ifdef DEBUG
  printf("Func: %s\n", "convert");
#endif

  if (mode_len == 0) {
    return NULL;
  }

  filter *filters = create_filters(mode_str, mode_len);

  // Allocate for return value
  int buf_len = RET_BUF;
  char *buf = (char *)malloc(buf_len);
  if (buf == NULL) {
    return NULL;
  }

  character c;
  int offset = 0, offset_buf = 0, total_len = 0;
  while (str_len > 0) {
    init_character(&c);
    extract(&c, str + offset, str_len);
    conv(&c, filters);
    total_len += c.clen;
    if (buf_len - total_len < 32) {
      char *tmp = (char *)realloc(buf, buf_len + RET_BUF);
      if (tmp == NULL) {
        return NULL;
      }
      buf = tmp;
      buf_len += RET_BUF;
    }
    strncpy(buf + offset_buf, c.cval, c.clen);
    offset += c.len;
    str_len -= c.len;
    offset_buf += c.clen;
  }
  *(buf + offset_buf) = 0x00;
  if (filters != NULL) {
    free(filters);
  }

  return buf;
}

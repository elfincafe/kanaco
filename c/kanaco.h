#ifndef _INCLUDE_KANACO_H
#define _INCLUDE_KANACO_H

#include <stdbool.h>
#include <stdint.h>

#define CNV_UNKNOWN 0
#define CNV_LOWER_R 1
#define CNV_UPPER_R 2
#define CNV_LOWER_N 4
#define CNV_UPPER_N 8
#define CNV_LOWER_A 16
#define CNV_UPPER_A 32
#define CNV_LOWER_S 64
#define CNV_UPPER_S 128
#define CNV_LOWER_K 256
#define CNV_UPPER_K 512
#define CNV_LOWER_H 1028
#define CNV_UPPER_H 2048
#define CNV_LOWER_C 4096
#define CNV_UPPER_C 8192

typedef struct _character {
    uint8_t val[8];
    uint8_t len;
    uint16_t conv;  // CNV_LOWER_* or CNV_UPPER_*
    char cval[8];   // converted value
    uint8_t clen;   // converted value length;
} character;

typedef void (*filter)(character*);

bool is_1byte(char*, int);
bool is_2byte(char*, int);
bool is_3byte(char*, int);
bool is_4byte(char*, int);
bool is_voiced(char*, int);
bool is_semi_voiced(char*, int);

void lower_r(character*);
void upper_r(character*);
void lower_n(character*);
void upper_n(character*);
void lower_a(character*);
void upper_a(character*);
void lower_s(character*);
void upper_s(character*);
void lower_k(character*);
void upper_k(character*);
void lower_h(character*);
void upper_h(character*);
void lower_c(character*);
void upper_c(character*);
void unknown(character*);

char* create_mode(char*, int);
void init_character(character*);
void conv(character*, char*);
void extract(character* c, char* s, int len);

extern char* convert(char*, int, char*, int);

#endif

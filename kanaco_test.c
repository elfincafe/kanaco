#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <assert.h>
#include "kanaco.h"

void test_is_1byte()
{
    char* s1 = "abc";
    assert(is_1byte(s1, strlen(s1)));
    char* s2 = "あいう";
    assert(is_1byte(s2, strlen(s2)));
}

void test_is_2bytes()
{
}

void test_is_3bytes()
{
}

void test_is_4bytes()
{
}

void test_is_voiced()
{

}

void test_is_semi_voiced()
{

}

int main(int argc, char** argv)
{
    test_is_1byte();
    test_is_2bytes();
    test_is_3bytes();
    test_is_4bytes();
    test_is_voiced();
    test_is_semi_voiced();
}

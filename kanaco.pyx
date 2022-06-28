
cdef extern from "kanaco.h":
    char *convert(char *s, int length, char *mode, int mode_len)
    void freeMemory(void *m)
cdef extern from "stdlib.h":
    void free(void*)

def conv(s, str m):
    cdef str string
    # cdef bytes tmp
    cdef str ret
    string = str(s)
    tmp = convert(<bytes>string.encode(), len(string), <bytes>m.encode(), len(m))
    ret = <str>tmp.decode()
    free(tmp)
    return ret

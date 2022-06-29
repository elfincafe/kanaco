from libc.stdlib cimport free

cdef extern from "kanaco.h":
    char *convert(char *s, int length, char *mode, int mode_len)

def conv(s, str m):
    cdef char *tmp = NULL
    cdef bytes ret
    try:
        t = type(s)
        if t is str:
            tmp = convert(<bytes>s.encode(), len(s), <bytes>m.encode(), len(m))
        elif t is int:
            tmp = convert(<bytes>str(s).encode(), len(str(s)), <bytes>m.encode(), len(m))
        elif t is bytes:
            tmp = convert(s, len(s), <bytes>m.encode(), len(m))
        else:
            raise TypeError("Invalid Data Type.")
        ret = <bytes>tmp
    finally:
        free(tmp)
    return <str>ret

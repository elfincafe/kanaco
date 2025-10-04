include ../common.mk

OBJS=kanaco.c

%.c: %.pyx
	$(CYTHON) -3 -f -o $@ $^

all: $(OBJS)
	$(CC) $(CFLAGS) $(INCLUDES) -I../c -fPIC -c $(INCLUDES) -o kanaco.o $^

clean:
	$(RM) -rf *.c *.o __pycache__ .pytest_cache

test:
	$(PYTEST) -s -vv test_kanaco.py

.PHONY: all clean test

include ../common.mk

%.static.o: %.c
	$(CC) $(CFLAGS) $(INCLUDES) -c -o $@ $^

%.shared.o: %.c
	$(CC) $(CFLAGS) $(INCLUDES) -fPIC -c -o $@ $^

all: kanaco.static.o kanaco.shared.o
	$(AR) rcsv libkanaco.a kanaco.static.o
	$(CC) $(CFLAGS) $(INCLUDES) -o libkanaco.so kanaco.shared.o -shared

clean:
	$(RM) -f *.o *.a *.so

.PHONY: all clean

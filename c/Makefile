CC=gcc
AR=ar
RM=rm

all: clean build

build:
	$(CC) -c kanaco.c -o kanaco.o
	$(AR) rsv libkanaco.a kanaco.o
	$(CC) -o kanaco main.c -I. -L. -static -lkanaco

clean:
	$(RM) -f *.a *.o core kanaco

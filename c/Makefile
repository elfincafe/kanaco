CC=gcc
AR=ar
RM=rm

all: clean build

build:
	$(CC) -c kanaco.c -o kanaco.o
	$(AR) rsv libkanaco.a kanaco.o

clean:
	$(RM) -f *.a *.o

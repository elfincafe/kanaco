CC=gcc
AR=ar
GO=go
GOBUILD=go build
GOTEST=go test
CYTHON=cython
PYCONFIG=python3-config
PYCFLAGS=python3-config --cflags
PYLDFLAGS=python3-config --ldflags
MAKE=make
MKDIR=mkdir
RM=rm
INCLUDE=-I./c
WORKDIR=.build


all: clean cbuild pybuild

cbuild:
	$(MKDIR) -p $(WORKDIR)
	$(CC) -c kanaco.c -o $(WORKDIR)/kanaco.c.o -fPIC -I.
	$(AR) rsv $(WORKDIR)/libkanaco.a $(WORKDIR)/kanaco.c.o
	# $(CC) -c kanaco.c -o $(WORKDIR)/kanaco.c.o -fPIC -I.
	# $(CC) $(WORKDIR)/kanaco.c.o -o $(WORKDIR)/libkanaco.so -shared

gotest: cbuild
	$(GOTEST) -v

pybuild: cbuild
	$(CYTHON) -3 -f -o $(WORKDIR)/kanaco.py.c kanaco.pyx
	$(CC) -o $(WORKDIR)/kanaco.py.o -c $(WORKDIR)/kanaco.py.c -fPIC -I. $$($(PYCONFIG) --cflags)
	$(CC) $(WORKDIR)/kanaco.py.o -o $(WORKDIR)/kanaco.so -shared $$($(PYCONFIG) --ldflags) -L$(WORKDIR) -lkanaco

clean:
	$(RM) -rf .build core


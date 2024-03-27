CC=clang
AR=ar
PKGCONFIG=pkg-config
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
INCLUDE="-I. -I/usr/include"
WORKDIR=/tmp/.kanaco
STATICDIR=$(WORKDIR)/static
SHAREDDIR=$(WORKDIR)/shared
UTFLAGS=$$($(PKGCONFIG) --cflags --libs cunit)


all: clean cbuild pybuild

cbuild:
	$(MKDIR) -p $(STATICDIR) $(SHAREDDIR)
	$(CC) -c kanaco.c -o $(STATICDIR)/kanaco.c.o -Wall
	$(AR) rsv $(STATICDIR)/libkanaco.a $(STATICDIR)/kanaco.c.o
	$(CC) -c kanaco.c -o $(SHAREDDIR)/libkanaco.so -fPIC -Wall

ctest: cbuild
	$(CC) -o $(STATICDIR)/kanaco_test -I. -L$(STATICDIR) kanaco_test.c -lkanaco $$($(PKGCONFIG) --cflags --libs cunit) && $(STATICDIR)/kanaco_test
	$(CC) -o $(SHAREDDIR)/kanaco_test -I. -L$(SHAREDDIR) kanaco_test.c -lkanaco $$($(PKGCONFIG) --cflags --libs cunit) && $(SHAREDDIR)/kanaco_test

gotest: cbuild
	$(GOTEST) -v

gobuild: cbuild

pybuild: cbuild
	$(RM) -f $(WORKDIR)/libkanaco.so
	$(CYTHON) -3 -f -o $(WORKDIR)/kanaco.py.c kanaco.pyx
	$(CC) -o $(WORKDIR)/kanaco.py.o -c $(WORKDIR)/kanaco.py.c -fPIC -I. $$($(PYCONFIG) --cflags)
	$(CC) $(WORKDIR)/kanaco.py.o -o $(WORKDIR)/kanaco.so -shared $$($(PYCONFIG) --ldflags) -L$(WORKDIR) -lkanaco

pytest:


clean:
	$(RM) -rf $(WORKDIR)

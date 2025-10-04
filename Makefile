include common.mk

UTFLAGS=$$($(PKGCONFIG) --cflags --libs cunit)

# ctest: cbuild
# 	$(CC) -o $(STATICDIR)/kanaco_test -I. -L$(STATICDIR) kanaco_test.c -lkanaco $$($(PKGCONFIG) --cflags --libs cunit) && $(STATICDIR)/kanaco_test
# 	$(CC) -o $(SHAREDDIR)/kanaco_test -I. -L$(SHAREDDIR) kanaco_test.c -lkanaco $$($(PKGCONFIG) --cflags --libs cunit) && $(SHAREDDIR)/kanaco_test

all:
	$(MAKE) -C c
	$(MAKE) -C cython
	$(CC) $(CFLAGS) $(INCLUDES) -o cython/kanaco.so cython/kanaco.o c/kanaco.shared.o -shared

clean:
	$(MAKE) -C c clean
	$(MAKE) -C cython clean
	$(RM) -rf *.so .pytest_cache

test:
	$(MAKE) -C cython test

c:
	$(MAKE) -C c

.PHONY: c all clean test

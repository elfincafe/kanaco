CC=gcc
AR=ar
CYTHON=cython
MAKE=make
RM=rm
PYCONFIG=python3-config
PYTEST=pytest

INCLUDES=-I. -I/usr/include $$($(PYCONFIG) --includes)
CFLAGS=-Wall -O2 $$($(PYCONFIG) --cflags --libs)

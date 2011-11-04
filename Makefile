include $(GOROOT)/src/Make.inc

TARG     = github.com/mattn/go-v8/v8
CGOFILES = v8.go

CGO_CFLAGS  = 
CGO_LDFLAGS  = -L. -lv8wrap -lstdc++
ifeq ($(GOOS),windows)
CGO_DEPS = v8wrap.dll
else
CGO_DEPS = libv8wrap.so
endif

include $(GOROOT)/src/Make.pkg

ifeq ($(GOOS),windows)
v8wrap.dll : v8wrap.cc
	g++ -shared -o v8wrap.dll -Ic:/mingw/include/v8 v8wrap.cc -lv8 -lstdc++ -lws2_32 -lwinmm
	dlltool -d v8wrap.def -l libv8wrap.a
else
libv8wrap.so : v8wrap.cc
	g++ -shared -o libv8wrap.so v8wrap.cc -lv8
endif

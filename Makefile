ifndef GOOS
GOOS := `go env GOOS`
endif

ifeq ($(GOOS),windows)
v8wrap.dll : v8wrap.cc v8.go
	g++ -shared -o v8wrap.dll -I. -Ic:/mingw/include/v8 v8wrap.cc -lv8 -lstdc++ -lws2_32 -lwinmm
	dlltool -d v8wrap.def -l libv8wrap.a
	go build -x .

clean:
	rm -f *.dll
else
.PHONY: go-install

all: libv8wrap.a go-install

libv8wrap.a : v8wrap.cc v8.go
	g++ `go env GOGCCFLAGS` -I. -c v8wrap.cc -lv8
	ar rvs libv8wrap.a v8wrap.o
	rm v8wrap.o

go-install:
	go install

clean:
	rm -f *.a
endif


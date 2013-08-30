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

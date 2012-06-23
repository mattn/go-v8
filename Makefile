ifeq ($(GOOS),windows)
v8wrap.dll : v8wrap.cc v8.go
	g++ -shared -o v8wrap.dll -I. -Ic:/mingw/include/v8 v8wrap.cc -lv8 -lstdc++ -lws2_32 -lwinmm
	dlltool -d v8wrap.def -l libv8wrap.a
	go build -x .

clean:
	rm -f *.dll
else
libv8wrap.so : v8wrap.cc v8.go
	g++ -fPIC -shared -o libv8wrap.so -I. v8wrap.cc -lv8
	go build -x .

clean:
	rm -f *.so
endif


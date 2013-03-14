go-v8
=====

WHATS:
------

  Go bindings for V8

INSTALL:
--------

WIN32

	# To Build v8 go package:
	# make v8wrap.dll
	# go install

	# To run go-v8 tests:
	# go test

	# To run example go exec:
	# cd example
	# copy v8wrap.dll
	# go build example.go
	# ./example

LINUX

	# To Build v8 go package:
	# make libv8wrap.so
	# go install

	# To run go-v8 tests:
	# LD_LIBRARY_PATH=. go test

	# To run example go exec:
	# cd example
	# go build example.go
	# LD_LIBRARY_PATH=.. ./example

MAC OS X

	# To Build v8 go package:
	# make libv8wrap.so
	# go install

	# To run go-v8 tests:
	# DYLD_LIBRARY_PATH=. go test

	# To run example go exec:
	# cd example
	# go build example.go
	# DYLD_LIBRARY_PATH=.. ./example

LICENSE:
--------

  under the MIT License: http://mattn.mit-license.org/2013

AUTHOR:
-------

  * Yasuhiro Matsumoto

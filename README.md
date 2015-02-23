go-v8
=====

WHATS:
------

  Go bindings for V8
  V8 Version 3.25.30

INSTALL:
--------

WIN32

	# To Build v8 go package:
	# CMake .
	# make
	# make install
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
	# CMake .
	# make
	# make install
	# go install

	# To run go-v8 tests:
	# LD_LIBRARY_PATH=. go test

	# To run example go exec:
	# cd example
	# go build example.go
	# ./example

MAC OS X

	# To Build v8 go package:

	# CMake .
	# make
	# make install
	# go install

	# To run go-v8 tests:
	# go test

	# To run example go exec:
	# cd example
	# go build example.go
	# ./example

LICENSE:
--------

  under the MIT License: http://mattn.mit-license.org/2013

AUTHOR:
-------

  * Yasuhiro Matsumoto

go-v8
=====

WHATS:
------

  Go bindings for V8

UPDATE for V8 3.21
------------------

Only got this going against v8 for OSX.  Build V8 for OSX as described by https://code.google.com/p/v8/wiki/BuildingWithGYP
Copy $(V8SRCROOT)/include/ files to /usr/local/include
Copy $(V8SRCROOT)/native/out/libv8\* files to /usr/local/lib

'''
go test
# testmain
github.com/mattn/go-v8(__DATA/__datacoal_nt): unexpected reloc for dynamic symbol _ZTVN10__cxxabiv117__class_type_infoE
github.com/mattn/go-v8(__DATA/__datacoal_nt): unhandled relocation for _ZTVN10__cxxabiv117__class_type_infoE (type 28 rtype 120)
FAIL    github.com/mattn/go-v8 [build failed]
'''

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

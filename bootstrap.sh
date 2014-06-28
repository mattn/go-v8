#!/bin/bash

## script from github.com/idada/v8-go

# find download tool
download=''
if hash curl 2>/dev/null; then
	download='curl -o'
elif hash wget 2>/dev/null; then
	download='wget -O'
else
	echo >&2 "You need to install 'curl' or 'wget'."
	exit 1
fi

v8_version="3.28.4.1"
v8_path="v8-$v8_version"

# check v8 installation
need_v8='false'
if [ ! -d $v8_path ] || [ ! -d $v8_path/out/native/ ]; then
	need_v8='true'
else
	libv8_base="`find $v8_path/out/native/ -name 'libv8_base.*.a' | head -1`"
	libv8_snapshot="`find $v8_path/out/native/ -name 'libv8_snapshot.a' | head -1`"

	if [ libv8_base == '' ] || [ libv8_snapshot == '' ]; then
		need_v8='true'
	fi
fi

# download and build v8
if [ $need_v8 == 'true' ]; then
	# download
	if [ ! -f $v8_path.tar.gz ]; then
		$download $v8_path.tar.gz https://codeload.github.com/v8/v8/tar.gz/$v8_version
	fi
	tar -xzvf $v8_path.tar.gz

	# begin
	cd $v8_path

	# we don't need ICU library
	svn checkout --force http://gyp.googlecode.com/svn/trunk build/gyp --revision 1685

	# build
	make i18nsupport=off native component=shared_library

  #end
  cd ..
fi

# for Linux
librt=''
if [ `go env | grep GOHOSTOS` == 'GOHOSTOS="linux"' ]; then
	librt='-lrt'
fi

# for Mac
libstdcpp=''
if  [ `go env | grep GOHOSTOS` == 'GOHOSTOS="darwin"' ]; then
	libstdcpp='-stdlib=libstdc++'
fi

echo "Name: v8
Description: v8 javascript engine
Version: $v8_version
Cflags: $libstdcpp -I`pwd` -I`pwd`/$v8_path/include
Libs: $libstdcpp `pwd`/$v8_path/out/native/obj.target/tools/gyp/libv8_base.a `pwd`/$v8_path/out/native/obj.target/tools/gyp/libv8_libbase.a `pwd`/$v8_path/out/native/obj.target/tools/gyp/libv8_snapshot.a $librt" > v8.pc

go install

go test

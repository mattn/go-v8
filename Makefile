CC=g++
CFLAGS= -I/usr/include -lv8 -dynamiclib -o $(TARGET) -DDEBUG -g -O2
SOURCES=v8wrap.cc
OBJECTS=$(SOURCES:.cc=.o) $(V8_DYLIB)
TARGET=libv8wrap.dylib

all: $(TARGET)

$(TARGET): $(OBJECTS)
	$(CC) $(CFLAGS) $< -o $@

clean:
	rm $(TARGET) $(OBJECTS)

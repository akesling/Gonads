CFLAGS=$(shell pkg-config --cflags xcb)
LFLAGS=$(shell pkg-config --libs xcb)

all: xcb_shapes

xcb_shapes:
	gcc $(CFLAGS) -o xcb_shapes xcb_shapes.c $(LFLAGS)

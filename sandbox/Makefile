CFLAGS=$(shell pkg-config --cflags xcb)
LFLAGS=$(shell pkg-config --libs xcb)

all: xcb_shapes xcb_shape_go

xcb_shapes:
	gcc $(CFLAGS) -o xcb_shapes xcb_shapes.c $(LFLAGS)

xcb_shapes_go:
	go build xcb_shapes_go.go

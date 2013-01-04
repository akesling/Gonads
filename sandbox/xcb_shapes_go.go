package main

/* Shameless cgo port of xcb_shapes.c
 * Michael Pratt <michael@pratt.im> */

import "unsafe"

/*
#cgo pkg-config: xcb
#include <stdlib.h>
#include <stdio.h>
#include <xcb/xcb.h>
*/
import "C"

// cgo doesn't support preproccessor defines
// see /usr/include/xcb/xproto.h
const (
	XCB_GC_FOREGROUND             = 4
	XCB_GC_GRAPHICS_EXPOSURES     = 65536
	XCB_CW_BACK_PIXEL             = 2
	XCB_CW_EVENT_MASK             = 2048
	XCB_EVENT_MASK_EXPOSURE       = 32768
	XCB_COPY_FROM_PARENT          = 0
	XCB_WINDOW_CLASS_INPUT_OUTPUT = 1
	XCB_EXPOSE                    = 12
	XCB_COORD_MODE_ORIGIN         = 0
	XCB_COORD_MODE_PREVIOUS       = 1
)

func main() {
	// Open connection
	c := C.xcb_connect(nil, nil)

	// Get first screen
	root_iterator := C.xcb_setup_roots_iterator(C.xcb_get_setup(c))
	screen := root_iterator.data

	// Create black (foreground) graphic context
	win := C.xcb_drawable_t(screen.root)

	foreground := C.xcb_gcontext_t(C.xcb_generate_id(c))
	mask := C.uint32_t(XCB_GC_FOREGROUND) | C.uint32_t(XCB_GC_GRAPHICS_EXPOSURES)
	values := []C.uint32_t{screen.black_pixel, 0}

	C.xcb_create_gc(c, foreground, win, mask, &values[0])

	// Ask for our window's id
	win = C.xcb_drawable_t(C.xcb_generate_id(c))

	// Create window
	mask = C.uint32_t(XCB_CW_BACK_PIXEL) | C.uint32_t(XCB_CW_EVENT_MASK)
	values[0] = screen.white_pixel
	values[1] = XCB_EVENT_MASK_EXPOSURE
	C.xcb_create_window(c, // Connection
		XCB_COPY_FROM_PARENT, // depth
		C.xcb_window_t(win),  // window id
		screen.root,          // parent
		0, 0,                 // x, y
		150, 150, // width, height
		10, // border_width
		XCB_WINDOW_CLASS_INPUT_OUTPUT, // class
		screen.root_visual,            // visual
		mask, &values[0])

	// Map window to screen
	C.xcb_map_window(c, C.xcb_window_t(win))

	// Flush request
	C.xcb_flush(c)

	points := []C.xcb_point_t{{10, 10}, {10, 20}, {20, 10}, {20, 20}}
	polyline := []C.xcb_point_t{{50, 10}, {5, 20}, {25, -20}, {10, 10}}
	segments := []C.xcb_segment_t{{100, 10, 140, 30}, {110, 25, 130, 60}}
	rectangles := []C.xcb_rectangle_t{{10, 50, 40, 20}, {80, 50, 10, 40}}
	arcs := []C.xcb_arc_t{{10, 100, 60, 40, 0, 90 << 6}, {90, 100, 55, 40, 0, 270 << 6}}

	// There must be a better way
	for e := C.xcb_wait_for_event(c); e != nil; e = C.xcb_wait_for_event(c) {
		switch int(e.response_type) & ^0x80 {
		case XCB_EXPOSE:
			// Draw points
			C.xcb_poly_point(c, XCB_COORD_MODE_ORIGIN, win, foreground, 4, &points[0])

			// Draw line
			C.xcb_poly_line(c, XCB_COORD_MODE_PREVIOUS, win, foreground, 4, &polyline[0])

			// Draw segments
			C.xcb_poly_segment(c, win, foreground, 2, &segments[0])

			// Draw rectangles
			C.xcb_poly_rectangle(c, win, foreground, 2, &rectangles[0])

			// Draw arcs
			C.xcb_poly_arc(c, win, foreground, 2, &arcs[0])

			C.xcb_flush(c)

			break
		default:
			// Unknown
			break
		}

		C.free(unsafe.Pointer(e))
	}
}

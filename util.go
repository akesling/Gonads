package main

import (
    "fmt"
    "log"
    "github.com/BurntSushi/xgb"
    "github.com/BurntSushi/xgb/xproto"
)


// Open the connection to the X server
func xConnect() *xgb.Conn {
    X, err := xgb.NewConn()
    if err != nil {
        log.Fatal(err)
    }
    return X
}

// Replace existing window manager
func usurpWM(X *xgb.Conn, screen *xproto.ScreenInfo) {
    wmName := fmt.Sprintf("WM_S%d", X.DefaultScreen)
    managerAtom, err := xproto.InternAtom(X, true, uint16(len(wmName)), wmName).Reply()
    if err != nil {
        log.Fatal(err)
    }

    fakeWindow, _ := xproto.NewWindowId(X)
    xproto.CreateWindow(X,                  // Connection
            screen.RootDepth,               // Depth
            fakeWindow,                     // Window Id
            screen.Root,                    // Parent Window
            -1000, -1000,                   // x, y
            1, 1,                           // width, height
            0,                              // border_width
            xproto.WindowClassInputOutput,  // class
            screen.RootVisual,              // visual
            xproto.CwEventMask|xproto.CwOverrideRedirect,
            []uint32{1, xproto.EventMaskPropertyChange})      // masks
    xproto.MapWindow(X, fakeWindow)
    err = xproto.SetSelectionOwnerChecked(X, fakeWindow, managerAtom.Atom, xproto.TimeCurrentTime).Check()
    if err != nil {
        log.Fatal(err)
    }
}

package main

import (
    "fmt"
    "log"
    "github.com/BurntSushi/xgb"
    "github.com/BurntSushi/xgb/xproto"
)

func main() {
    // Open the connection to the X server
    X, err := xgb.NewConn()
    if err != nil {
        log.Fatal(err)
    }
    defer X.Close()

    setup := xproto.Setup(X)
    // Get the first screen
    screen := setup.DefaultScreen(X)

    // Replace existing window manager
    wmName := fmt.Sprintf("WM_S%d", X.DefaultScreen)
    managerAtom, err := xproto.InternAtom(X, true, uint16(len(wmName)), wmName).Reply()
    if err != nil {
        log.Fatal(err)
    }

    fakeWindow, _ := xproto.NewWindowId(X)
    xproto.CreateWindow(X,                  // Connection
            screen.RootDepth,               // Depth
            blankWindow,                    // Window Id
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
        fmt.Println("foo")
        log.Fatal(err)
    }

    arcs := []xproto.Arc{
        {10, 100, 60, 40, 0, 90 << 6},
        {90, 100, 55, 40, 0, 270 << 6}};

    // Create black (foreground) graphic context
    foreground, _ := xproto.NewGcontextId(X)
    mask := uint32(xproto.GcForeground | xproto.GcGraphicsExposures)
    values := []uint32{screen.BlackPixel, 0}
    xproto.CreateGC(X, foreground, xproto.Drawable(screen.Root), mask, values)

    // Ask for our window's Id
    win, _ := xproto.NewWindowId(X)
    winDrawable := xproto.Drawable(win)

    // Create the window
    mask = uint32(xproto.CwBackPixel | xproto.CwEventMask)
    values = []uint32{screen.WhitePixel, xproto.EventMaskExposure}
    xproto.CreateWindow(X,                  // Connection
            screen.RootDepth,               // Depth
            win,                            // Window Id
            screen.Root,                    // Parent Window
            0, 0,                           // x, y
            150, 150,                       // width, height
            10,                             // border_width
            xproto.WindowClassInputOutput,  // class
            screen.RootVisual,              // visual
            mask, values)                   // masks

    // Map the window on the screen
    xproto.MapWindow(X, win)

    // Obey the window-delete protocol
    tp := "WM_PROTOCOLS"
    prp := "WM_DELETE_WINDOW"
    typeAtom, _ := xproto.InternAtom(X, true, uint16(len(tp)), tp).Reply()
    propertyAtom, _ := xproto.InternAtom(X, true, uint16(len(prp)), prp).Reply()

    data := make([]byte, 4)
    xgb.Put32(data, uint32(propertyAtom.Atom))
    xproto.ChangeProperty(X, xproto.PropModeReplace, win, typeAtom.Atom, xproto.AtomAtom, 32, 1, data)

    // Main loop
    for {
        evt, err := X.WaitForEvent()
        fmt.Printf("An event of type %T occured.\n", evt)

        if evt == nil && err == nil {
            fmt.Println("Exiting....")
            return
        } else if err != nil {
            log.Fatal(err)
        }

        switch event := evt.(type) {
            case xproto.ExposeEvent:
                /* We draw the arcs */
                xproto.PolyArc(X, winDrawable, foreground, arcs)
            case xproto.ClientMessageEvent:
                if len(event.Data.Data32) > 0 {
                    data := xproto.Atom(event.Data.Data32[0])
                    if data == propertyAtom.Atom {
                        return
                    } else {
                        atomName, _ := xproto.GetAtomName(X, data).Reply()
                        fmt.Println(atomName.Name)
                    }
                } else {
                    atomName, _ := xproto.GetAtomName(X, event.Type).Reply()
                    fmt.Println(atomName.Name)
                }
            default:
                /* Unknown event type, ignore it */
        }
    }
    return
}

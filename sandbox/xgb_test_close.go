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

    // geometric objects
    points := []xproto.Point{
        {10, 10},
        {10, 20},
        {20, 10},
        {20, 20}};

    polyline := []xproto.Point{
        {50, 10},
        { 5, 20},     // rest of points are relative
        {25,-20},
        {10, 10}};

    segments := []xproto.Segment{
        {100, 10, 140, 30},
        {110, 25, 130, 60}};

    rectangles := []xproto.Rectangle{
        { 10, 50, 40, 20},
        { 80, 50, 10, 40}};

    arcs := []xproto.Arc{
        {10, 100, 60, 40, 0, 90 << 6},
        {90, 100, 55, 40, 0, 270 << 6}};

    setup := xproto.Setup(X)
    // Get the first screen
    screen := setup.DefaultScreen(X)

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

    // It turns out that we need the window ID as a byte-stream... WTF!!
    // xprop.ChangeProp(xu, win, 8, "WM_NAME", "STRING", ([]byte)(name))
    // ChangeProp(xu *xgbutil.XUtil, win xproto.Window, format byte, prop string, typ string, data []byte)
    data := make([]byte, 4)
    xgb.Put32(data, uint32(propertyAtom.Atom))
    xproto.ChangeProperty(X, xproto.PropModeReplace, win, typeAtom.Atom, xproto.AtomAtom, 32, 1, data)

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
                /* We draw the points */
                xproto.PolyPoint(X, xproto.CoordModeOrigin, winDrawable, foreground, points)

                /* We draw the polygonal line */
                xproto.PolyLine(X, xproto.CoordModePrevious, winDrawable, foreground, polyline)

                /* We draw the segments */
                xproto.PolySegment(X, winDrawable, foreground, segments)

                /* We draw the rectangles */
                xproto.PolyRectangle(X, winDrawable, foreground, rectangles)

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

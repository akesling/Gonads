Gonads
======

Gonads is a X-centric window manager: Go + XMonad - Haskell = Gonads.

Gonads, "the baller window manager," is inspired by a combination of XMonad and Rob Pike's NewSqueak concurrent windowing system (http://swtch.com/~rsc/thread/cws.pdf).  It is first and foremost a tiling window manager and secondly many other things.

Gonads Axioms (these may change over time):

    Axiom) *Design with a specific customer in mind*:
        Response) We build for programmers.
    A) *Minimize mouse activity*:
        R) By focusing on a tiling paradigm we try to optimize keyboard interaction with the windowing system.
    A) *Don't make the user wait*:
        R) We regard speed very highly.  The user _ever_ waiting due to the window system is a failing on our part.
    A) *Less is infinitely more*:
        R) This does one thing well: manage windows.  One should forget the system exists and be able to get their work done.
        R) Fewer lines and core concepts generally translates to fewer bugs and a more elegant experience.
    A) *Follow the principle of least astonishment*:
        R) Do what users expect and little (preferably nothing) else. We do our best not to second-guess the user.
        R) Be consistent in general and use consistent mechanisms for related actions.

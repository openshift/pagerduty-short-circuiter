package ui

import (
	"log"
	"os/exec"
	"sync"

	tcellterm "git.sr.ht/~rockorager/tcell-term"
	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
	"github.com/rivo/tview"
)

type Terminal struct {
	*tview.Box

	term    *tcellterm.VT
	running bool
	cmd     *exec.Cmd
	tui     *TUI
	w       int
	h       int
	sync.RWMutex
}

func NewTerminal(cmd *exec.Cmd, index int, tui *TUI) *Terminal {
	t := &Terminal{
		Box:  tview.NewBox(),
		term: tcellterm.New(),
		cmd:  cmd,
		tui:  tui,
	}
	return t
}

func (t *Terminal) Draw(s tcell.Screen) {
	t.Box.DrawForSubclass(s, t)

	x, y, w, h := t.GetInnerRect()
	view := views.NewViewPort(s, x, y, w, h)
	t.term.SetSurface(view)
	if w != t.w || h != t.h {
		t.w = w
		t.h = h
		t.term.Resize(w, h)
	}

	if !t.running {
		err := t.term.Start(t.cmd)
		if err != nil {
			log.Print(err)
			return
			// panic(err)
		}
		t.term.Attach(t.HandleEvent)
		t.running = true
	}
	if t.HasFocus() {
		cy, cx, style, vis := t.term.Cursor()
		if vis {
			s.ShowCursor(cx+x, cy+y)
			s.SetCursorStyle(style)
		} else {
			s.HideCursor()
		}
	}
	t.term.Draw()
}

// the tterm tview wrapper swallows the eventclosed
// that is emitted here so we add the event type to switch
// and handle accordingly?
func (t *Terminal) HandleEvent(ev tcell.Event) {
	switch event := ev.(type) {
	case *tcellterm.EventClosed:
		go func() {
			t.tui.App.QueueUpdateDraw(func() {
				// run closed in main event  loop thread
				t.Closed(event)
			})
		}()
	case *tcellterm.EventRedraw:
		go func() {
			t.tui.App.QueueUpdateDraw(func() {})
		}()
	}
}

func (t *Terminal) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return t.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		t.term.HandleEvent(event)
	})
}

func (t *Terminal) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
	return t.WrapMouseHandler(func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
		return t.term.HandleEvent(event), nil
	})
}

// Function for Closing a Terminal
func (t *Terminal) Closed(ev *tcellterm.EventClosed) {
	ExitSlide(CurrentActivePage, t.tui)
}

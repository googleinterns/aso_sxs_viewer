package createwindow

import (
	"fmt"

	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
)

//CreateInputWindow creates window to capture keycodes
func CreateInputWindow(layout Layout, X *xgb.Conn, screenInfo *xproto.ScreenInfo, a *QuitStruct) (*xgb.Conn, xproto.Window, *QuitStruct, error) {
	wid, _ := xproto.NewWindowId(X)
	cookie := xproto.CreateWindowChecked(X, screenInfo.RootDepth, wid, screenInfo.Root,
		0, 0, uint16(layout.w), uint16(layout.h), 0,
		xproto.WindowClassInputOutput, screenInfo.RootVisual,
		xproto.CwBackPixel|xproto.CwEventMask,
		[]uint32{
			0xffffffff,
			xproto.EventMaskEnterWindow |
				xproto.EventMaskLeaveWindow |
				xproto.EventMaskKeyPress |
				xproto.EventMaskKeyRelease |
				xproto.EventMaskStructureNotify})

	if err := cookie.Check(); err != nil {
		return nil, 0, a, err
	}

	xproto.MapWindow(X, wid)
	xproto.ConfigureWindow(X, wid,
		xproto.ConfigWindowX|xproto.ConfigWindowY,
		[]uint32{
			uint32(layout.y), uint32(layout.x),
		})

	a.Quitters = append(a.Quitters, InputWindow{wid, X, true})

	return X, wid, a, nil
}

func keyPresshandler(X *xgb.Conn, wid xproto.Window, a *QuitStruct, quitfunc func(*QuitStruct), b xproto.KeyPressEvent) {
	fmt.Println(b.Detail)
}

func keyReleasehandler(X *xgb.Conn, wid xproto.Window, a *QuitStruct, quitfunc func(*QuitStruct), b xproto.KeyReleaseEvent) {
	fmt.Println(b.Detail)
}

func enterNotifyhandler(X *xgb.Conn, wid xproto.Window, a *QuitStruct, quitfunc func(*QuitStruct), b xproto.EnterNotifyEvent) error {
	cookie := xproto.GrabKeyboard(X, true, wid, xproto.TimeCurrentTime, xproto.GrabModeAsync, xproto.GrabModeAsync)
	if _, err := cookie.Reply(); err != nil {
		return err
	}
	// keybinding.Focus = false
	return nil
}

func leaveNotifyhandler(X *xgb.Conn, wid xproto.Window, a *QuitStruct, quitfunc func(*QuitStruct), b xproto.LeaveNotifyEvent) error {
	cookie := xproto.UngrabKeyboardChecked(X, xproto.TimeCurrentTime)
	if err := cookie.Check(); err != nil {
		return err
	}
	return nil
}

func unmapNotifyhandler(X *xgb.Conn, wid xproto.Window, a *QuitStruct, quitfunc func(*QuitStruct), b xproto.UnmapNotifyEvent) {
	fmt.Println("unmap notify event")
	fmt.Println("connection interrupted")
	a.Quitters[len(a.Quitters)-1].SetToClose(false)
	quitfunc(a)
}

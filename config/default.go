package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jezek/xgb"
	"github.com/jezek/xgb/randr"
	"github.com/jezek/xgb/xproto"
)

var (
	AsoSxSViewerDir    = filepath.Join(os.Getenv("HOME"), ".aso_sxs_viewer")
	AsoSxSViewerConfig = filepath.Join(AsoSxSViewerDir, "config.textproto")
	DefaultURL         = "https://mail.google.com/"
	DefaultCSSSelector = struct {
		Selector string
		Position int32
	}{
		Selector: "input",
		Position: 7,
	}
	DefaultBrowserCount      = int32(2)
	DefaultUseCookies        = false
	DefaultUserDataDirPrefix = filepath.Join(AsoSxSViewerDir, "profiles")
	DefaultRootWindowWidth   = int32(1600)
	DefaultRootWindowHeight  = int32(900)

	BrowserWindowExample = fmt.Sprintf(`# You may use the template below to add window_overrides	
#	window_overrides: {
#		url: "%s"
#		css_selector: {
#			selector: "%s"
#			position: %d
#		}
#	}`, DefaultURL, DefaultCSSSelector.Selector, DefaultCSSSelector.Position)
)

func DefaultRootWindowSize() (width, height uint16) {
	width, height = uint16(DefaultRootWindowWidth), uint16(DefaultRootWindowHeight)

	X, err := xgb.NewConn()
	if err != nil {
		return
	}
	if err := randr.Init(X); err != nil {
		return
	}

	screens, err := randr.GetScreenResourcesCurrent(X, xproto.Setup(X).DefaultScreen(X).Root).Reply()
	if err != nil {
		return
	}

	crtc, err := randr.GetCrtcInfo(X, screens.Crtcs[0], xproto.TimeCurrentTime).Reply()
	if err != nil {
		return
	}
	return uint16(0.8 * float64(crtc.Width)), uint16(0.8 * float64(crtc.Height))
}

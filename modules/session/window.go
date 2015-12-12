package session

import (
    "fmt"
    "io"
    "strconv"
    "github.com/BurntSushi/xgb/xproto"
    "github.com/BurntSushi/xgbutil/ewmh"
    "github.com/BurntSushi/xgbutil/xgraphics"
    "github.com/BurntSushi/xgbutil/xwindow"
)

type SessionProcess struct {
    PID     uint                        `json:"pid"`
    Command string                      `json:"command"`
    User    string                      `json:"user"`
    UID     uint                        `json:"uid"`
}

type SessionWindowGeometry struct {
    X      int                          `json:"x"`
    Y      int                          `json:"y"`
    Width  uint                         `json:"width"`
    Height uint                         `json:"height"`
}

type SessionWindow struct {
    ID            uint32                `json:"id"`
    Title         string                `json:"name"`
    Roles         []string              `json:"roles,omitempty"`
    Flags         map[string]bool       `json:"flags"`
    IconUri       string                `json:"icon,omitempty"`
    Workspace     uint                  `json:"workspace"`
    AllWorkspaces bool                  `json:"all_workspaces,omitempty"`
    Screen        uint                  `json:"screen,omitempty"`
    Dimensions    SessionWindowGeometry `json:"dimensions"`
    Process       SessionProcess        `json:"process"`
    Active        bool                  `json:"active,omitempty"`
}

func (self *SessionModule) GetWindow(window_id string) (SessionWindow, error) {
    window := SessionWindow{}

    if id, err := self.toX11WindowId(window_id); err == nil {
        xgbWindow           := xwindow.New(self.X, id)
        geom, _             := xgbWindow.DecorGeometry()
        process             := SessionProcess{}
        window.ID            = uint32(id)
        window.Title, _      = ewmh.WmNameGet(self.X, id)
        //window.IconUri           = r.Path() + "/" + id
        windowWorkspace, _  := ewmh.WmDesktopGet(self.X, id)
        activeWinId, _      := ewmh.ActiveWindowGet(self.X)

        if windowWorkspace == 0xFFFFFFFF {
            window.AllWorkspaces = true
        }else{
            window.Workspace = windowWorkspace
        }

        if id == activeWinId {
            window.Active = true
        }

    //  calculate window dimensions from desktop and window frame boundaries
        window.Dimensions.Width  = uint(geom.Width())
        window.Dimensions.Height = uint(geom.Height())
        window.Dimensions.X      = geom.X()
        window.Dimensions.Y      = geom.Y()

    //  fill in process details
        process.PID, _           = ewmh.WmPidGet(self.X, id)
        window.Process           = process

    //  get window state flags
        window.Flags = make(map[string]bool)

    //  minimized
        if self.x11HasWmState(id, `_NET_WM_STATE_HIDDEN`) {
            window.Flags["minimized"] = true
        }

    //  shaded
        if self.x11HasWmState(id, `_NET_WM_STATE_SHADED`) {
            window.Flags[`shaded`] = true
        }

    //  maximized
        if self.x11HasWmState(id, `_NET_WM_STATE_MAXIMIZED_VERT`) && self.x11HasWmState(id, `_NET_WM_STATE_MAXIMIZED_HORZ`) {
            window.Flags[`maximized`] = true
        }

    //  above
        if self.x11HasWmState(id, `_NET_WM_STATE_ABOVE`) {
            window.Flags[`above`] = true
        }

    //  below
        if self.x11HasWmState(id, `_NET_WM_STATE_BELOW`) {
            window.Flags[`below`] = true
        }

    //  urgent
        if self.x11HasWmState(id, `_NET_WM_STATE_DEMANDS_ATTENTION`) {
            window.Flags[`urgent`] = true
        }

    //  skip_taskbar
        if self.x11HasWmState(id, `_NET_WM_STATE_SKIP_TASKBAR`) {
            window.Flags[`skip_taskbar`] = true
        }

    //  skip_pager
        if self.x11HasWmState(id, `_NET_WM_STATE_SKIP_PAGER`) {
            window.Flags[`skip_pager`] = true
        }

    //  sticky
        if self.x11HasWmState(id, `_NET_WM_STATE_STICKY`) {
            window.Flags[`sticky`] = true
        }

    //  fullscreen
        if self.x11HasWmState(id, `_NET_WM_STATE_FULLSCREEN`) {
            window.Flags[`fullscreen`] = true
        }

    //  modal
        if self.x11HasWmState(id, `_NET_WM_STATE_MODAL`) {
            window.Flags[`modal`] = true
        }

        return window, nil
    }else{
        return window, err
    }
}

func (self *SessionModule) WriteWindowIcon(window_id string, width uint, height uint, w io.Writer) error {
    if id, err := self.toX11WindowId(window_id); err == nil {
        if icon, err := xgraphics.FindIcon(self.X, id, int(width), int(height)); err == nil {
            return icon.WritePng(w)
        }else{
            return err
        }
    }else{
        return err
    }
}

func (self *SessionModule) WriteWindowImage(window_id string, w io.Writer) error {
    if id, err := self.toX11WindowId(window_id); err == nil {
        drawable   := xproto.Drawable(id)

        if image, err := xgraphics.NewDrawable(self.X, drawable); err == nil {
            image.XSurfaceSet(id)
            image.XDraw()

            return image.WritePng(w)
        }else{
            return err
        }
    }else{
        return err
    }
}

func (self *SessionModule) GetAllWindows() ([]SessionWindow, error) {
    clients, _ := ewmh.ClientListGet(self.X)

//  allocate window objects
    windows := make([]SessionWindow, len(clients))

//  for each window...
    for i, id := range clients {
        windows[i], _ = self.GetWindow(strconv.Itoa(int(id)))
    }

    return windows, nil
}

func (self *SessionModule) RaiseWindow(window_id string) error {
    if id, err := self.toX11WindowId(window_id); err == nil {
    //  unhide the window
        self.removeState(id, `_NET_WM_STATE_HIDDEN`)

    //  move the window to the stop of the stacking order
        ewmh.RestackWindow(self.X, id)

    //  activate the window
        ewmh.ActiveWindowReq(self.X, id)
    }else{
        return err
    }

    return nil
}

func (self *SessionModule) ActionWindow(window_id string, atom string) error {
    if id, err := self.toX11WindowId(window_id); err == nil {
        self.addState(id, atom)
    }else{
        return err
    }

    return nil
}

func (self *SessionModule) MaximizeWindowHorizontal(window_id string) error {
    return self.ActionWindow(window_id, "_NET_WM_STATE_MAXIMIZED_HORZ")
}

func (self *SessionModule) MaximizeWindowVertical(window_id string) error {
    return self.ActionWindow(window_id, "_NET_WM_STATE_MAXIMIZED_VERT")
}

func (self *SessionModule) MaximizeWindow(window_id string) error {
    if err := self.MaximizeWindowHorizontal(window_id); err == nil {
        return self.MaximizeWindowVertical(window_id)
    }else{
        return err
    }
}

func (self *SessionModule) RestoreWindow(window_id string) error {
    if id, err := self.toX11WindowId(window_id); err == nil {
        self.removeState(id, `_NET_WM_STATE_MAXIMIZED_VERT`)
        self.removeState(id, `_NET_WM_STATE_MAXIMIZED_HORZ`)
    }else{
        return err
    }

    return nil
}

func (self *SessionModule) MinimizeWindow(window_id string) error {
    return self.ActionWindow(window_id, `_NET_WM_STATE_HIDDEN`)
}

func (self *SessionModule) HideWindow(window_id string) error {
    return self.ActionWindow(window_id, `_NET_WM_STATE_HIDDEN`)
}

func (self *SessionModule) ShowWindow(window_id string) error {
    if id, err := self.toX11WindowId(window_id); err == nil {
    //  unhide the window
        self.removeState(id, `_NET_WM_STATE_HIDDEN`)
    }else{
        return err
    }

    return nil
}


func (self *SessionModule) MoveWindow(window_id string, x int, y int) error {
    if id, err := self.toX11WindowId(window_id); err == nil {
        return ewmh.MoveWindow(self.X, id, x, y)
    }else{
        return err
    }
}


func (self *SessionModule) ResizeWindow(window_id string, width uint, height uint) error {
    return fmt.Errorf("Not implemented")
}

func (self *SessionModule) toX11WindowId(window_id string) (xproto.Window, error) {
    if id_number, err := strconv.Atoi(window_id); err == nil {
        return xproto.Window(uint32(id_number)), nil
    }else{
        return xproto.Window(0), err
    }
}

func (self *SessionModule) x11HasWmState(id xproto.Window, state string) bool {
    states, _ := ewmh.WmStateGet(self.X, id)

    for _, s := range states {
        if s == state {
            return true
        }
    }

    return false
}

func (self *SessionModule) addState(id xproto.Window, state string) error {
    return ewmh.WmStateReq(self.X, id, 0, state)
}

func (self *SessionModule) removeState(id xproto.Window, state string) error {
    states, _ := ewmh.WmStateGet(self.X, id)

    for i, s := range states {
        if s == state {
            states = append(states[:i], states[i+1:]...)

            if err := ewmh.WmStateSet(self.X, id, states); err != nil {
                return err
            }
        }
    }

    return nil
}

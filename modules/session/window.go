package main

import (
    "github.com/BurntSushi/xgb/xproto"
    "github.com/BurntSushi/xgbutil/ewmh"
    "github.com/BurntSushi/xgbutil/xgraphics"
    "github.com/BurntSushi/xgbutil/xwindow"
    "io"
    "strconv"
)

type SessionProcess struct {
    PID     uint   `json:"pid"`
    Command string `json:"command"`
    User    string `json:"user"`
    UID     uint   `json:"uid"`
}

type SessionWindowGeometry struct {
    X      int  `json:"x"`
    Y      int  `json:"y"`
    Width  uint `json:"width"`
    Height uint `json:"height"`
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

func (self *SessionPlugin) GetWindow(window_id string) (window SessionWindow, err error) {
    id, _      := self.Atox(window_id)
    xgb_window := xwindow.New(self.X, id)
    window = SessionWindow{}
    process := SessionProcess{}

    window.ID = uint32(id)
    window.Title, _ = ewmh.WmNameGet(self.X, id)
    //window.IconUri           = r.Path() + "/" + id
    window_workspace, _ := ewmh.WmDesktopGet(self.X, id)
    active_window, _ := ewmh.ActiveWindowGet(self.X)

    if window_workspace == 0xFFFFFFFF {
        window.AllWorkspaces = true
    } else {
        window.Workspace = window_workspace
    }

    if id == active_window {
        window.Active = true
    }

    window_geometry, _ := xgb_window.DecorGeometry()

    //  calculate window dimensions from desktop and window frame boundaries
    window.Dimensions.Width = uint(window_geometry.Width())
    window.Dimensions.Height = uint(window_geometry.Height())
    window.Dimensions.X = window_geometry.X()
    window.Dimensions.Y = window_geometry.Y()

    process.PID, _ = ewmh.WmPidGet(self.X, id)

    window.Process = process

    //  get window state flags
    window.Flags = make(map[string]bool)

    //  minimized
    if self.hasState(id, "_NET_WM_STATE_HIDDEN") {
        window.Flags["minimized"] = true
    }

    //  shaded
    if self.hasState(id, "_NET_WM_STATE_SHADED") {
        window.Flags["shaded"] = true
    }

    //  maximized
    if self.hasState(id, "_NET_WM_STATE_MAXIMIZED_VERT") && self.hasState(id, "_NET_WM_STATE_MAXIMIZED_HORZ") {
        window.Flags["maximized"] = true
    }

    //  above
    if self.hasState(id, "_NET_WM_STATE_ABOVE") {
        window.Flags["above"] = true
    }

    //  below
    if self.hasState(id, "_NET_WM_STATE_BELOW") {
        window.Flags["below"] = true
    }

    //  urgent
    if self.hasState(id, "_NET_WM_STATE_DEMANDS_ATTENTION") {
        window.Flags["urgent"] = true
    }

    //  skip_taskbar
    if self.hasState(id, "_NET_WM_STATE_SKIP_TASKBAR") {
        window.Flags["skip_taskbar"] = true
    }

    //  skip_pager
    if self.hasState(id, "_NET_WM_STATE_SKIP_PAGER") {
        window.Flags["skip_pager"] = true
    }

    //  sticky
    if self.hasState(id, "_NET_WM_STATE_STICKY") {
        window.Flags["sticky"] = true
    }

    //  fullscreen
    if self.hasState(id, "_NET_WM_STATE_FULLSCREEN") {
        window.Flags["fullscreen"] = true
    }

    //  modal
    if self.hasState(id, "_NET_WM_STATE_MODAL") {
        window.Flags["modal"] = true
    }

    return
}

func (self *SessionPlugin) WriteWindowIcon(window_id string, width uint, height uint, w io.Writer) (err error) {
    id, err     := self.Atox(window_id)
    if err != nil {
       return
    }

    icon, err := xgraphics.FindIcon(self.X, id, int(width), int(height))

    if err != nil {
        return
    }

    err = icon.WritePng(w)
    return
}

func (self *SessionPlugin) WriteWindowImage(window_id string, w io.Writer) (err error) {
    id, err      := self.Atox(window_id)
    if err != nil {
       return
    }

    drawable   := xproto.Drawable(id)
    image, err := xgraphics.NewDrawable(self.X, drawable)

    if err != nil {
        return
    }

    image.XSurfaceSet(id)
    image.XDraw()

    err = image.WritePng(w)
    return
}

func (self *SessionPlugin) GetAllWindows() ([]SessionWindow, error) {
    clients, _ := ewmh.ClientListGet(self.X)

    //  allocate window objects
    windows := make([]SessionWindow, len(clients))

    //  for each window...
    for i, id := range clients {
        windows[i], _ = self.GetWindow(strconv.Itoa(int(id)))
    }

    return windows, nil
}

func (self *SessionPlugin) RaiseWindow(window_id string) (err error) {
    id, err := self.Atox(window_id)

    if err != nil {
       return
    }

//  unhide the window
    self.removeState(id, "_NET_WM_STATE_HIDDEN")

//  move the window to the stop of the stacking order
    ewmh.RestackWindow(self.X, id)

//  activate the window
    ewmh.ActiveWindowReq(self.X, id)

    return
}

func (self *SessionPlugin) ActionWindow(window_id string, atom string) (err error) {
    id, err := self.Atox(window_id)

    if err != nil {
       return
    }

    self.addState(id, atom)
    return
}

func (self *SessionPlugin) MaximizeWindowHorizontal(window_id string) (err error) {
    return self.ActionWindow(window_id, "_NET_WM_STATE_MAXIMIZED_HORZ")
}

func (self *SessionPlugin) MaximizeWindowVertical(window_id string) (err error) {
    return self.ActionWindow(window_id, "_NET_WM_STATE_MAXIMIZED_VERT")
}

func (self *SessionPlugin) MaximizeWindow(window_id string) (err error) {
    self.MaximizeWindowHorizontal(window_id)
    self.MaximizeWindowVertical(window_id)
    return
}

func (self *SessionPlugin) RestoreWindow(window_id string) (err error) {
    id, err := self.Atox(window_id)

    if err != nil {
       return
    }

    self.removeState(id, "_NET_WM_STATE_MAXIMIZED_VERT")
    self.removeState(id, "_NET_WM_STATE_MAXIMIZED_HORZ")
    return
}

func (self *SessionPlugin) MinimizeWindow(window_id string) (err error) {
    return self.ActionWindow(window_id, "_NET_WM_STATE_HIDDEN")
}

func (self *SessionPlugin) HideWindow(window_id string) (err error) {
    return self.ActionWindow(window_id, "_NET_WM_STATE_HIDDEN")
}

func (self *SessionPlugin) ShowWindow(window_id string) (err error) {
    id, err := self.Atox(window_id)

    if err != nil {
       return
    }

//  unhide the window
    self.removeState(id, "_NET_WM_STATE_HIDDEN")
    return
}


func (self *SessionPlugin) MoveWindow(window_id string, x int, y int) (err error) {
    id, err := self.Atox(window_id)

    if err != nil {
       return
    }

    ewmh.MoveWindow(self.X, id, x, y)
    return
}


func (self *SessionPlugin) ResizeWindow(window_id string, width uint, height uint) (err error) {

    return
}




func (self *SessionPlugin) Atox(window_id string) (id xproto.Window, err error) {
    id_number, err := strconv.Atoi(window_id)

    if err != nil {
        return
    }

    id = xproto.Window(uint32(id_number))
    return
}

func (self *SessionPlugin) hasState(id xproto.Window, state string) bool {
    states, _ := ewmh.WmStateGet(self.X, id)

    if indexOf(states, state) >= 0 {
        return true
    }

    return false
}

func (self *SessionPlugin) addState(id xproto.Window, state string) (err error) {
    err = ewmh.WmStateReq(self.X, id, 0, state)
    return
}

func (self *SessionPlugin) removeState(id xproto.Window, state string) (err error) {
    states, _ := ewmh.WmStateGet(self.X, id)

    if i := indexOf(states, state); i >= 0 {
        states = append(states[:i], states[i+1:]...)
        err = ewmh.WmStateSet(self.X, id, states)
    }

    return
}
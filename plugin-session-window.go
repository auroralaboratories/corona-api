package main

import (
    "strconv"
    "io"
    "github.com/BurntSushi/xgb/xproto"
    "github.com/BurntSushi/xgbutil/xwindow"
    "github.com/BurntSushi/xgbutil/xgraphics"
    "github.com/BurntSushi/xgbutil/ewmh"
)

type SessionProcess struct {
    PID      uint                           `json:"pid"`
    Command  string                         `json:"command"`
    User     string                         `json:"user"`
    UID      uint                           `json:"uid"`
}

type SessionWindowGeometry struct {
    X        int                            `json:"x"`
    Y        int                            `json:"y"`
    Width    uint                           `json:"width"`
    Height   uint                           `json:"height"`
}

type SessionWindow struct {
    ID            uint32                      `json:"id"`
    Title         string                      `json:"name"`
    Roles         []string                    `json:"roles,omitempty"`
    Flags         []string                    `json:"flags,omitempty"`
    IconUri       string                      `json:"icon,omitempty"`
    Workspace     uint                        `json:"workspace"`
    AllWorkspaces bool                        `json:"all_workspaces,omitempty"`
    Screen        uint                        `json:"screen,omitempty"`
    Dimensions    SessionWindowGeometry       `json:"dimensions"`
    Process       SessionProcess              `json:"process"`
    Active        bool                        `json:"active,omitempty"`
}


func (self *SessionPlugin) GetWindow(window_id string) (window SessionWindow, err error) {
    id_number, _            := strconv.Atoi(window_id)
    id                      := xproto.Window(uint32(id_number))
    xgb_window              := xwindow.New(self.X, id)
    window                   = SessionWindow{}
    process                 := SessionProcess{}


    window.ID                = uint32(id)
    window.Title, _          = ewmh.WmNameGet(self.X, id)
    //window.IconUri           = r.Path() + "/" + id
    window_workspace, _     := ewmh.WmDesktopGet(self.X, id)
    active_window, _        := ewmh.ActiveWindowGet(self.X)


    if window_workspace == 0xFFFFFFFF {
        window.AllWorkspaces = true
    }else{
        window.Workspace     = window_workspace
    }

    if id == active_window {
        window.Active = true
    }

    window_geometry, _      := xgb_window.DecorGeometry()

//  calculate window dimensions from desktop and window frame boundaries
    window.Dimensions.Width  = uint(window_geometry.Width())
    window.Dimensions.Height = uint(window_geometry.Height())
    window.Dimensions.X      = window_geometry.X()
    window.Dimensions.Y      = window_geometry.Y()

    process.PID, _           = ewmh.WmPidGet(self.X, id)

    window.Process = process

    return
}


func (self *SessionPlugin) WriteWindowIcon(window_id string, width uint, height uint, w io.Writer) (err error) {
    id_number, _            := strconv.Atoi(window_id)
    id                      := xproto.Window(uint32(id_number))
    icon, err               := xgraphics.FindIcon(self.X, id, int(width), int(height))

    if err != nil {
        return
    }

    err = icon.WritePng(w)
    return
}


func (self *SessionPlugin) WriteWindowImage(window_id string, w io.Writer) (err error) {
    id_number, _            := strconv.Atoi(window_id)
    id                      := xproto.Window(uint32(id_number))
    drawable                := xproto.Drawable(id)
    image, err              := xgraphics.NewDrawable(self.X, drawable)

    if err != nil {
        return
    }

    image.XSurfaceSet(id)
    image.XDraw()


    err = image.WritePng(w)
    return
}


func (self *SessionPlugin) GetAllWindows() ([]SessionWindow, error) {
    clients, _          := ewmh.ClientListGet(self.X)

//  allocate window objects
    windows             := make([]SessionWindow, len(clients))

//  for each window...
    for i, id := range clients {
        windows[i], _ = self.GetWindow(strconv.Itoa(int(id)))
    }

    return windows, nil
}


func (self *SessionPlugin) RaiseWindow(window_id string) (err error) {
    id_number, _            := strconv.Atoi(window_id)
    id                      := xproto.Window(uint32(id_number))
    err                      = ewmh.ActiveWindowSet(self.X, id)

    if err != nil {
        return
    }

    err                      = ewmh.RestackWindow(self.X, id)

    return
}
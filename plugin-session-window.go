package main

import (
    "strconv"
    "github.com/BurntSushi/xgb/xproto"
    "github.com/BurntSushi/xgbutil/xwindow"
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
    Workspace     uint                        `json:"workspace,omitempty"`
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

    if err != nil {
        return window, err
    }


    window.ID                = uint32(id)
    window.Title, _          = ewmh.WmNameGet(self.X, id)
    //window.IconUri           = r.Path() + "/" + id
    window_workspace, _     := ewmh.WmDesktopGet(self.X, id)
    active_window, _        := ewmh.ActiveWindowGet(self.X)

    if window_workspace == 0xFFFFFFFF {
        window.AllWorkspaces = true
    }else{
        window.Workspace = window_workspace
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
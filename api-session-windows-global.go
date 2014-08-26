package main

import (
    "github.com/ant0ine/go-json-rest/rest"
    "github.com/BurntSushi/xgbutil/xwindow"
    "github.com/BurntSushi/xgbutil/ewmh"
)

type ApiSessionProcess struct {
    PID      uint                           `json:"pid"`
    Command  string                         `json:"command"`
    User     string                         `json:"user"`
    UID      uint                           `json:"uid"`
}

type ApiSessionWindowGeometry struct {
    X        int                            `json:"x"`
    Y        int                            `json:"y"`
    Width    uint                           `json:"width"`
    Height   uint                           `json:"height"`
}

type ApiSessionWindow struct {
    ID            uint32                      `json:"id"`
    Title         string                      `json:"name"`
    Roles         []string                    `json:"roles,omitempty"`
    Flags         []string                    `json:"flags,omitempty"`
    IconUri       string                      `json:"icon,omitempty"`
    Workspace     uint                        `json:"workspace,omitempty"`
    AllWorkspaces bool                        `json:"all_workspaces,omitempty"`
    Screen        uint                        `json:"screen,omitempty"`
    Dimensions    ApiSessionWindowGeometry    `json:"dimensions"`
    Process       ApiSessionProcess           `json:"process"`
    Active        bool                        `json:"active,omitempty"`
}

func (self *SprinklesAPI) GetWindows(w rest.ResponseWriter, r *rest.Request) {
    clients,   _        := ewmh.ClientListGet(self.X)
    active_window, _    := ewmh.ActiveWindowGet(self.X)

//  allocate window objects
    windows             := make([]ApiSessionWindow, len(clients))

//  for each window...
    for i, id := range clients {
        xgb_window              := xwindow.New(self.X, id)
        window                  := ApiSessionWindow{}
        process                 := ApiSessionProcess{}

        window.ID                = uint32(id)
        window.Title, _          = ewmh.WmNameGet(self.X, id)
        //window.IconUri           = r.Path() + "/" + id
        window_workspace, _     := ewmh.WmDesktopGet(self.X, id)

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
        windows[i] = window
    }

//  output
    w.WriteJson(&windows)
}
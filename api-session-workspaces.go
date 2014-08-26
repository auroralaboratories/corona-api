package main

import (
    "fmt"
    "strconv"
    "github.com/ant0ine/go-json-rest/rest"
    "github.com/BurntSushi/xgbutil/ewmh"
)

type ApiSessionWorkspace struct {
    Number      uint                        `json:"number"`
    Name        string                      `json:"name"`
    Screen      uint                        `json:"screen,omitempty"`
    IsCurrent   bool                        `json:"current,omitempty"`
    WindowCount uint                        `json:"window_count,omitempty"`
}

type ApiSessionWorkspaces struct {
    CurrentWorkspace uint                   `json:"current"`
    Workspaces       []ApiSessionWorkspace  `json:"workspaces"`
}



func (self *SprinklesAPI) GetWorkspaces(w rest.ResponseWriter, r *rest.Request) {
    workspace_count,   _  := ewmh.NumberOfDesktopsGet(self.X)
    current_workspace, _  := ewmh.CurrentDesktopGet(self.X)
    workspace_names, _    := ewmh.DesktopNamesGet(self.X)

//  allocate workspace objects
    workspaces            := make([]ApiSessionWorkspace, workspace_count)

//  for each workspace...
    for i := uint(0); i < workspace_count; i++ {
        workspace           := ApiSessionWorkspace{}
        workspace.Number     = i
        workspace.Name       = workspace_names[i]

    //  flag current if this is the current workspace
        if i == current_workspace {
          workspace.IsCurrent = true
        }

        workspaces[i] = workspace
    }

//  output
    w.WriteJson(&ApiSessionWorkspaces{
        CurrentWorkspace: current_workspace,
        Workspaces:       workspaces,
    })
}



func (self *SprinklesAPI) GetCurrentWorkspace(w rest.ResponseWriter, r *rest.Request) {
    current_workspace, _  := ewmh.CurrentDesktopGet(self.X)
    workspace_names, _    := ewmh.DesktopNamesGet(self.X)

    workspace             := ApiSessionWorkspace{}
    workspace.Number       = current_workspace
    workspace.Name         = workspace_names[current_workspace]

//  output
    w.WriteJson(&workspace)
}



func (self *SprinklesAPI) GetWorkspace(w rest.ResponseWriter, r *rest.Request) {
    workspace_number, _   := strconv.Atoi(r.PathParam("number"))
    workspace_count,   _  := ewmh.NumberOfDesktopsGet(self.X)
    workspace_names, _    := ewmh.DesktopNamesGet(self.X)

    if uint(workspace_number) >= workspace_count {
        w.WriteHeader(404)
        return
    }

    workspace             := ApiSessionWorkspace{}
    workspace.Number       = uint(workspace_number)
    workspace.Name         = workspace_names[workspace_number]

//  output
    w.WriteJson(&workspace)
}



func (self *SprinklesAPI) SetWorkspace(w rest.ResponseWriter, r *rest.Request) {
    workspace_number, _   := strconv.Atoi(r.PathParam("number"))
    workspace_count,   _  := ewmh.NumberOfDesktopsGet(self.X)

    if uint(workspace_number) >= workspace_count {
        w.WriteHeader(400)
        w.WriteJson(&SprinklesAPIError{
            Code:    400,
            Message: fmt.Sprintf("Cannot change to non-existent workspace %d", workspace_number),
        })

        return
    }

//  set the current workspace
    ewmh.CurrentDesktopSet(self.X, uint(workspace_number))

//  show the workspace we just switched to
    self.GetWorkspace(w, r)
}



package main

import (
    "fmt"
    "errors"
    "github.com/BurntSushi/xgbutil/ewmh"
)

type SessionWorkspace struct {
    Number      uint                        `json:"number"`
    Name        string                      `json:"name"`
    Screen      uint                        `json:"screen,omitempty"`
    IsCurrent   bool                        `json:"current,omitempty"`
    WindowCount uint                        `json:"window_count,omitempty"`
}


func (self *SessionPlugin) GetWorkspace(workspace_number uint) (workspace SessionWorkspace, err error) {
    workspace              = SessionWorkspace{}
    workspace_count,   _  := ewmh.NumberOfDesktopsGet(self.X)
    workspace_names, _    := ewmh.DesktopNamesGet(self.X)
    current_workspace, _  := ewmh.CurrentDesktopGet(self.X)

    if uint(workspace_number) >= workspace_count {
        err = errors.New(fmt.Sprintf("Cannot get workspace index %d, only %d workspaces exist", workspace_number, workspace_count))
        return
    }

    workspace.Number       = workspace_number
    
    if workspace_number < len(workspace_names) {
        workspace.Name         = workspace_names[workspace_number]
    }

    if workspace_number == current_workspace {
        workspace.IsCurrent = true
    }

    return
}


func (self *SessionPlugin) GetAllWorkspaces() (workspaces []SessionWorkspace, err error) {
    workspace_count,   _  := ewmh.NumberOfDesktopsGet(self.X)

//  allocate workspace objects
    workspaces             = make([]SessionWorkspace, workspace_count)

//  for each workspace...
    for i := uint(0); i < workspace_count; i++ {
        workspace, _ := self.GetWorkspace(i)
        workspaces[i] = workspace
    }

    return
}
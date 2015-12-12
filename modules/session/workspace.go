package session

import (
    "fmt"
    "github.com/BurntSushi/xgbutil/ewmh"
)

type SessionWorkspace struct {
    Number      uint                        `json:"number"`
    Name        string                      `json:"name"`
    Screen      uint                        `json:"screen,omitempty"`
    IsCurrent   bool                        `json:"current,omitempty"`
    WindowCount uint                        `json:"window_count,omitempty"`
}


func (self *SessionModule) GetWorkspace(workspaceNumber uint) (SessionWorkspace, error) {
    workspace := SessionWorkspace{
        Number: workspaceNumber,
    }

    count,   _           := ewmh.NumberOfDesktopsGet(self.X)
    names, _             := ewmh.DesktopNamesGet(self.X)
    currentWorkspace, _  := ewmh.CurrentDesktopGet(self.X)

    if uint(workspaceNumber) >= count {
        return workspace, fmt.Errorf("Cannot get workspace index %d, only %d workspaces exist", workspaceNumber, count)
    }

//  set workspace name if one is given
//  TODO: what does this array look like if only a few workspaces have names, but not all?
//        might get a nil assign error
    if int(workspaceNumber) < len(names) {
        workspace.Name  = names[workspaceNumber]
    }

//  flag which workspace is the active one
    if workspaceNumber == currentWorkspace {
        workspace.IsCurrent = true
    }

    return workspace, nil
}


func (self *SessionModule) GetAllWorkspaces() ([]SessionWorkspace, error) {
    if count, err  := ewmh.NumberOfDesktopsGet(self.X); err == nil {
    //  allocate workspace objects
        workspaces := make([]SessionWorkspace, count)

    //  for each workspace...
        for i := uint(0); i < count; i++ {
            if workspace, err := self.GetWorkspace(i); err == nil {
                workspaces[i] = workspace
            }else{
                return nil, err
            }
        }

        return workspaces, nil
    }else{
        return nil, err
    }
}
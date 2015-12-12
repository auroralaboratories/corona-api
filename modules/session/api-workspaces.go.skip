package main

import (
    "strconv"
    "github.com/ant0ine/go-json-rest/rest"
)



func (self *CoronaAPI) GetWorkspaces(w rest.ResponseWriter, r *rest.Request) {
    workspaces, _ := self.Plugin("Session").(*SessionPlugin).GetAllWorkspaces()
    w.WriteJson(&workspaces)
}

func (self *CoronaAPI) GetCurrentWorkspace(w rest.ResponseWriter, r *rest.Request) {
    workspaces, _ := self.Plugin("Session").(*SessionPlugin).GetAllWorkspaces()

    for _, ws := range workspaces {
        if ws.IsCurrent {
            w.WriteJson(&ws)
            return
        }
    }

    w.WriteHeader(404)
}

func (self *CoronaAPI) GetWorkspace(w rest.ResponseWriter, r *rest.Request) {
    workspace_number, _ := strconv.Atoi(r.PathParam("number"))
    workspace, err := self.Plugin("Session").(*SessionPlugin).GetWorkspace(uint(workspace_number))

    if err != nil {
        rest.Error(w, err.Error(), 400)
    }else{
        w.WriteJson(&workspace)
    }
}



// func (self *CoronaAPI) SetWorkspace(w rest.ResponseWriter, r *rest.Request) {
//     workspace_number, _   := strconv.Atoi(r.PathParam("number"))
//     workspace_count,   _  := ewmh.NumberOfDesktopsGet(self.X)

//     if uint(workspace_number) >= workspace_count {
//         w.WriteHeader(400)
//         w.WriteJson(&CoronaAPIError{
//             Code:    400,
//             Message: fmt.Sprintf("Cannot change to non-existent workspace %d", workspace_number),
//         })

//         return
//     }

// //  set the current workspace
//     ewmh.CurrentDesktopSet(self.X, uint(workspace_number))

// //  show the workspace we just switched to
//     self.GetWorkspace(w, r)
// }



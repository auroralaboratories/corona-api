package main

import (
    "fmt"
    "github.com/ant0ine/go-json-rest/rest"
)

func (self *SprinklesAPI) GetWindows(w rest.ResponseWriter, r *rest.Request) {
    windows, _ := self.Plugin("Session").(*SessionPlugin).GetAllWindows()

    for i, _ := range windows {
      windows[i].IconUri = fmt.Sprintf("%s/v1/session/windows/%d/icon", r.BaseUrl(), windows[i].ID)
    }

//  output
    w.WriteJson(&windows)
}

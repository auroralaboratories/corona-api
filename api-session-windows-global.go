package main

import (
    "github.com/ant0ine/go-json-rest/rest"
)

func (self *SprinklesAPI) GetWindows(w rest.ResponseWriter, r *rest.Request) {
    windows, _ := self.Plugin("Session").(*SessionPlugin).GetAllWindows()

//  output
    w.WriteJson(&windows)
}


func (self *SprinklesAPI) GetWindow(w rest.ResponseWriter, r *rest.Request) {
    window, _ := self.Plugin("Session").(*SessionPlugin).GetWindow(r.PathParam("id"))

//  output
    w.WriteJson(&window)
}
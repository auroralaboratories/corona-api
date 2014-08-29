package main

import (
    "bytes"
    "net/http"
    "github.com/ant0ine/go-json-rest/rest"
)


func (self *SprinklesAPI) GetWindow(w rest.ResponseWriter, r *rest.Request) {
    window, _ := self.Plugin("Session").(*SessionPlugin).GetWindow(r.PathParam("id"))

//  output
    w.WriteJson(&window)
}

func (self *SprinklesAPI) GetWindowIcon(w rest.ResponseWriter, r *rest.Request) {
    var buffer bytes.Buffer

    err := self.Plugin("Session").(*SessionPlugin).WriteWindowIcon(r.PathParam("id"), 0, 0, &buffer)

    if err != nil {
      rest.Error(w, err.Error(), 400)
    }else{
      w.Header().Set("Content-Type", "image/png")
      w.(http.ResponseWriter).Write(buffer.Bytes())
    }
}
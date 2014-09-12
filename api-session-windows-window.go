package main

import (
    "fmt"
    "bytes"
    "strconv"
    "net/http"
    "github.com/ant0ine/go-json-rest/rest"
)


func (self *CoronaAPI) GetWindow(w rest.ResponseWriter, r *rest.Request) {
    window, _ := self.Plugin("Session").(*SessionPlugin).GetWindow(r.PathParam("id"))

    window.IconUri = fmt.Sprintf("%s/v1/session/windows/%s/icon", r.BaseUrl(), r.PathParam("id"))

//  output
    w.WriteJson(&window)
}

func (self *CoronaAPI) GetWindowIcon(w rest.ResponseWriter, r *rest.Request) {
    var buffer bytes.Buffer
    var width  uint
    var height uint

    if r.PathParam("w") != "" {
      w, _ := strconv.Atoi("w")
      width = uint(w)
    }

    if r.PathParam("h") != "" {
      h, _ := strconv.Atoi("h")
      height = uint(h)
    }

    err := self.Plugin("Session").(*SessionPlugin).WriteWindowIcon(r.PathParam("id"), width, height, &buffer)

    if err != nil {
      rest.Error(w, err.Error(), 400)
    }else{
      w.Header().Set("Content-Type", "image/png")
      w.(http.ResponseWriter).Write(buffer.Bytes())
    }
}


func (self *CoronaAPI) GetWindowImage(w rest.ResponseWriter, r *rest.Request) {
    var buffer bytes.Buffer

    err := self.Plugin("Session").(*SessionPlugin).WriteWindowImage(r.PathParam("id"), &buffer)

    if err != nil {
      rest.Error(w, err.Error(), 400)
    }else{
      w.Header().Set("Content-Type", "image/png")
      w.(http.ResponseWriter).Write(buffer.Bytes())
    }
}


func (self *CoronaAPI) RaiseWindow(w rest.ResponseWriter, r *rest.Request) {
  err := self.Plugin("Session").(*SessionPlugin).RaiseWindow(r.PathParam("id"))

  if err != nil {
    rest.Error(w, err.Error(), 500)
  }
}
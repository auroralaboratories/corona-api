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

func (self *CoronaAPI) ActionWindow(w rest.ResponseWriter, r *rest.Request) {
  plugin := self.Plugin("Session").(*SessionPlugin)
  id     := r.PathParam("id")

  var err error

  switch r.PathParam("action") {
  case "maximize":
    plugin.MaximizeWindow(id)

  case "max-x":
    plugin.MaximizeWindowHorizontal(id)

  case "max-y":
    plugin.MaximizeWindowVertical(id)

  case "minimize":
    plugin.MinimizeWindow(id)

  case "restore":
    plugin.RestoreWindow(id)

  case "hide":
    plugin.HideWindow(id)

  case "show":
    plugin.ShowWindow(id)

  case "raise":
    plugin.RaiseWindow(id)
  }


  if err != nil {
    rest.Error(w, err.Error(), 500)
  }
}


func (self *CoronaAPI) MoveWindow(w rest.ResponseWriter, r *rest.Request) {
  plugin := self.Plugin("Session").(*SessionPlugin)
  id     := r.PathParam("id")

  x, err := strconv.Atoi(r.PathParam("x"))
  if err != nil {
    rest.Error(w, err.Error(), 500)
    return
  }

  y, err := strconv.Atoi(r.PathParam("y"))
  if err != nil {
    rest.Error(w, err.Error(), 500)
    return
  }

  plugin.MoveWindow(id, x, y)
}

func (self *CoronaAPI) ResizeWindow(w rest.ResponseWriter, r *rest.Request) {
  plugin := self.Plugin("Session").(*SessionPlugin)
  id     := r.PathParam("id")

  width, err := strconv.Atoi(r.PathParam("width"))
  if err != nil {
    rest.Error(w, err.Error(), 500)
    return
  }

  height, err := strconv.Atoi(r.PathParam("height"))
  if err != nil {
    rest.Error(w, err.Error(), 500)
    return
  }

  plugin.ResizeWindow(id, uint(width), uint(height))
}
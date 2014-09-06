package main

import (
    "github.com/ant0ine/go-json-rest/rest"
)

type ApiStatus struct {
      Ok bool `json:"ok"`
}

func (self *SprinklesAPI) GetApiStatus(w rest.ResponseWriter, r *rest.Request) {
      w.WriteJson(&ApiStatus{
          Ok: true,
      })
}

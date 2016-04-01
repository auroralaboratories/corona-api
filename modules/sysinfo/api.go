package main

import (
	"github.com/ant0ine/go-json-rest/rest"
)

func (self *CoronaAPI) GetSystemStats(w rest.ResponseWriter, r *rest.Request) {
	stats, _ := self.Plugin("System").(*SystemPlugin).GetAllStats()
	w.WriteJson(&stats)
}

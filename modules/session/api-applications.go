package main


import (
    "strings"
    "github.com/ant0ine/go-json-rest/rest"
)

func (self *CoronaAPI) GetApplications(w rest.ResponseWriter, r *rest.Request) {
    applications := self.Plugin("Session").(*SessionPlugin).GetAppList()
    w.WriteJson(&applications)
}

func (self *CoronaAPI) GetAppByName(w rest.ResponseWriter, r *rest.Request) {
    applications := self.Plugin("Session").(*SessionPlugin).GetAppByName(strings.Replace(r.PathParam("name") , "%20", " ", -1))
    w.WriteJson(&applications)
}

func (self *CoronaAPI) LaunchAppByName(w rest.ResponseWriter, r *rest.Request) {
    self.Plugin("Session").(*SessionPlugin).LaunchAppByName(strings.Replace(r.PathParam("name") , "%20", " ", -1))
    w.WriteHeader(204)
}

func (self *CoronaAPI) SearchAppByName(w rest.ResponseWriter, r *rest.Request) {
    applications := self.Plugin("Session").(*SessionPlugin).SearchAppByName(strings.Replace(r.PathParam("pattern") , "%20", " ", -1))
    w.WriteJson(&applications)
}

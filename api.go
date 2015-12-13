package main

import (
    "fmt"
    "net/http"

    "github.com/julienschmidt/httprouter"
    "github.com/codegangsta/negroni"
    "github.com/auroralaboratories/corona-api/util"
    "github.com/auroralaboratories/corona-api/modules/session"
)

type API struct {
    Address           string
    Port              int
    UiDirectory       string

    router            *httprouter.Router
    server            *negroni.Negroni
}

func NewApi() *API {
    return &API{
        Address:     `localhost`,
        Port:        25672,
        UiDirectory: `embedded`,
    }
}

func (self *API) Serve() error {
    self.router    = httprouter.New()

    if err := self.loadRoutes(); err != nil {
        return err
    }

    self.server = negroni.New()
    self.server.Use(negroni.NewRecovery())
    self.server.UseHandler(self.router)

    self.server.Run(fmt.Sprintf("%s:%d", self.Address, self.Port))
    return nil
}

func (self *API) loadRoutes() error {
    self.router.GET(`/api/status`, func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
        util.Respond(w, http.StatusOK, map[string]interface{}{
            `started_at`: util.StartedAt,
        }, nil)
    })

    if err := session.LoadRoutes(self.router); err != nil {
        return err
    }

    return nil
}


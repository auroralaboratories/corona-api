package main

import (
    "fmt"
    "github.com/ant0ine/go-json-rest/rest"
    "net/http"
)

type SprinklesAPI struct {
    RestHandler rest.ResourceHandler
    HttpHandler http.Handler
    Port        uint
    Interface   string
    Plugins     map[string]IPlugin
}

type SprinklesAPIError struct {
    Code    int
    Message string
}

func (self *SprinklesAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if _, ok := r.Header["Upgrade"]; ok {
        self.HttpHandler.ServeHTTP(w, r)
    } else {
        self.RestHandler.ServeHTTP(w, r)
    }
}

func (self *SprinklesAPI) Throw(err error, code int, w rest.ResponseWriter) {
    rest.Error(w, err.Error(), code)
}

func (self *SprinklesAPI) Plugin(name string) IPlugin {
    return self.Plugins[name]
}

func (self *SprinklesAPI) Init() (err error) {
    if self.Port == 0 {
        self.Port = 9521
    }

    self.Plugins = make(map[string]IPlugin)

    self.Plugins["Session"] = &SessionPlugin{}
    self.Plugins["Config"]  = &ConfigPlugin{}
    self.Plugins["Bus"]     = &BusPlugin{}

    for name, plugin := range self.Plugins {
        logger.Infof("Initializing plugin: %s", name)
        err := plugin.Init()
        PanicIfErr(err)
    }

//  initialize non-rest HTTP handlers
    mux := http.NewServeMux()
    mux.HandleFunc("/v1/bus", self.WebsocketClientConnect)
    self.HttpHandler = mux


//  setup CORS middleware
    self.RestHandler = rest.ResourceHandler{
        PreRoutingMiddlewares: []rest.Middleware{
            &rest.CorsMiddleware{
                RejectNonCorsRequests: false,
                OriginValidator: func(origin string, request *rest.Request) bool {
                    return true
                },
                AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"},
                AllowedHeaders: []string{
                    "Accept", "Content-Type", "X-Custom-Header", "Origin"},
                AccessControlAllowCredentials: true,
                AccessControlMaxAge:           3600,
            },
        },
    }

    //  setup routes
    err = self.RestHandler.SetRoutes(
        &rest.Route{"GET", "/",
            self.GetApiStatus },

        &rest.Route{"GET", "/v1",
            self.GetApiStatus },

        &rest.Route{"GET", "/v1/session/workspaces",
            self.GetWorkspaces },

        &rest.Route{"GET", "/v1/session/workspaces/current",
            self.GetCurrentWorkspace },

        &rest.Route{"GET", "/v1/session/workspaces/:number",
            self.GetWorkspace },

        // &rest.Route{"PUT", "/v1/session/workspaces/:number",
        //     self.SetWorkspace },

        &rest.Route{"GET", "/v1/session/windows",
            self.GetWindows },

        &rest.Route{"GET", "/v1/session/windows/:id",
            self.GetWindow },

        &rest.Route{"GET", "/v1/session/windows/:id/icon",
            self.GetWindowIcon },

        &rest.Route{"GET", "/v1/session/windows/:id/image",
            self.GetWindowImage },

        &rest.Route{"PUT", "/v1/session/windows/:id/raise",
            self.RaiseWindow },
    )

    if err != nil {
        logger.Errorf("Error occurred registering routes: %s", err)
    }

    return
}

func (self *SprinklesAPI) Serve() error {
    err := http.ListenAndServe(fmt.Sprintf("%s:%d", self.Interface, self.Port), self)
    return err
}

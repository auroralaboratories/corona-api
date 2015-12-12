// package main

// import (
//     "fmt"
//     "os"
//     "io"
//     "io/ioutil"
//     "path"
//     "strings"
//     "github.com/ant0ine/go-json-rest/rest"
//     "github.com/shurcooL/go/github_flavored_markdown"
//     "net/http"
// )

// type CoronaOptions struct {
//     StaticRoot          string
//     NoMarkdownAutoparse bool
// }

// type CoronaAPI struct {
//     RestHandler rest.ResourceHandler
//     HttpHandler http.Handler
//     FileHandler http.Handler
//     Port        uint
//     Interface   string
//     Plugins     map[string]IPlugin
//     Options     CoronaOptions
// }

// type CoronaAPIError struct {
//     Code    int
//     Message string
// }

// func (self *CoronaAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// //  handle websocket HTTP connection upgrade
//     if _, ok := r.Header["Upgrade"]; ok {
//         self.HttpHandler.ServeHTTP(w, r)
//     } else {
//     //  paths that start with '/static' route to the file server
//         if strings.HasPrefix(r.URL.Path, "/static/") {
//         //  allow user to get raw source with the ?raw=true query string
//             raw := (r.URL.Query().Get("raw") == "true")

//             if !self.Options.NoMarkdownAutoparse && !raw {
//             //  paths ending in '.md' are markdown and returned as HTML
//                 if strings.HasSuffix(r.URL.Path, ".md") {
//                 //  open target file
//                     mdpath := path.Join(self.Options.StaticRoot, strings.TrimPrefix(r.URL.Path, "/static"))
//                     buffer, err := ioutil.ReadFile(mdpath)

//                     if err != nil {
//                         w.WriteHeader(500)
//                         return
//                     }

//                     parsed := github_flavored_markdown.Markdown(buffer)
//                     w.Header().Set("Content-Type", "text/html")
//                     io.WriteString(w, `<html><head><meta charset="utf-8"><style>code, div.highlight { tab-size: 4; }</style><link href="https://github.com/assets/github.css" media="all" rel="stylesheet" type="text/css" /></head><body><article class="markdown-body entry-content" style="padding: 30px;">`)
//                     w.Write(parsed)
//                     io.WriteString(w, `</article></body></html>`)
//                     return
//                 }
//             }

//             self.FileHandler.ServeHTTP(w, r)

//         }else if strings.HasPrefix(r.URL.Path, "/static") {
//             http.Redirect(w, r, "/static/", 301)
//             return

//     //  all other paths are treated as REST calls
//         }else{
//             self.RestHandler.ServeHTTP(w, r)
//         }
//     }
// }

// func (self *CoronaAPI) Throw(err error, code int, w rest.ResponseWriter) {
//     rest.Error(w, err.Error(), code)
// }

// func (self *CoronaAPI) Plugin(name string) IPlugin {
//     return self.Plugins[name]
// }

// func (self *CoronaAPI) Init() (err error) {
//     if self.Port == 0 {
//         self.Port = 9521
//     }

//     self.Plugins            = make(map[string]IPlugin)
//     self.Plugins["Session"] = &SessionPlugin{}
//     self.Plugins["Config"]  = &ConfigPlugin{}
//     self.Plugins["Bus"]     = &BusPlugin{}
//     self.Plugins["System"]  = &SystemPlugin{}

//     for name, plugin := range self.Plugins {
//         logger.Infof("Initializing plugin: %s", name)
//         err := plugin.Init()

//         if err != nil {
//             logger.Fatalf("Unable to initialize plugin %s: %s", name, err.Error())
//         }
//     }

// //  initialize non-rest HTTP handlers
//     mux := http.NewServeMux()
//     mux.HandleFunc("/v1/bus", self.WebsocketClientConnect)
//     self.HttpHandler = mux


// //  initialize static file server
//     if entry, err := os.Stat(self.Options.StaticRoot); err == nil {
//         if entry.IsDir() {
//             logger.Infof("Initializing static assets root at %s", self.Options.StaticRoot)
//             self.FileHandler = http.StripPrefix("/static/", http.FileServer(http.Dir(self.Options.StaticRoot)))
//         }
//     }else{
//         logger.Errorf("Unable to read static assets root at %s", self.Options.StaticRoot)
//     }



// //  setup CORS middleware
//     self.RestHandler = rest.ResourceHandler{
//         PreRoutingMiddlewares: []rest.Middleware{
//             &rest.CorsMiddleware{
//                 RejectNonCorsRequests: false,
//                 OriginValidator: func(origin string, request *rest.Request) bool {
//                     return true
//                 },
//                 AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"},
//                 AllowedHeaders: []string{
//                     "Accept", "Content-Type", "X-Custom-Header", "Origin"},
//                 AccessControlAllowCredentials: true,
//                 AccessControlMaxAge:           3600,
//             },
//         },
//     }

//     //  setup routes
//     err = self.RestHandler.SetRoutes(
//         &rest.Route{"GET", "/",
//             self.GetApiStatus },

//         &rest.Route{"GET", "/v1",
//             self.GetApiStatus },

//         &rest.Route{"GET", "/v1/session/workspaces",
//             self.GetWorkspaces },

//         &rest.Route{"GET", "/v1/session/workspaces/current",
//             self.GetCurrentWorkspace },

//         &rest.Route{"GET", "/v1/session/workspaces/:number",
//             self.GetWorkspace },

//         // &rest.Route{"PUT", "/v1/session/workspaces/:number",
//         //     self.SetWorkspace },

//         &rest.Route{"GET", "/v1/session/windows",
//             self.GetWindows },

//         &rest.Route{"GET", "/v1/session/windows/:id",
//             self.GetWindow },

//         &rest.Route{"GET", "/v1/session/windows/:id/icon",
//             self.GetWindowIcon },

//         &rest.Route{"GET", "/v1/session/windows/:id/image",
//             self.GetWindowImage },

//         &rest.Route{"PUT", "/v1/session/windows/:id/:action",
//             self.ActionWindow },

//         &rest.Route{"PUT", "/v1/session/windows/:id/move/:x/:y",
//             self.MoveWindow },

//         &rest.Route{"PUT", "/v1/session/windows/:id/resize/:width/:height",
//             self.ResizeWindow },

//         &rest.Route{"GET", "/v1/session/applications",
//             self.GetApplications },

//         &rest.Route{"GET", "/v1/session/applications/:name",
//             self.GetAppByName },

//         &rest.Route{"GET", "/v1/session/applications/find/:pattern",
//             self.SearchAppByName },

//         &rest.Route{"GET", "/v1/session/applications/:name/launch",
//             self.LaunchAppByName },

//         &rest.Route{"GET", "/v1/system/stats",
//             self.GetSystemStats },
//     )

//     if err != nil {
//         logger.Errorf("Error occurred registering routes: %s", err)
//     }

//     return
// }

// func (self *CoronaAPI) Serve() error {
//     err := http.ListenAndServe(fmt.Sprintf("%s:%d", self.Interface, self.Port), self)
//     return err
// }

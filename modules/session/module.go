package session

import (
    "bytes"
    "fmt"
    "net/http"
    "strings"
    "github.com/BurntSushi/xgbutil"
    "github.com/codegangsta/cli"
    "github.com/julienschmidt/httprouter"
    "github.com/shutterstock/go-stockutil/stringutil"
    "github.com/auroralaboratories/corona-api/util"
    "github.com/auroralaboratories/corona-api/modules"
)

type SessionModule struct {
    modules.BaseModule

    X *xgbutil.XUtil
}

func RegisterSubcommands() []cli.Command {
    return []cli.Command{}
}

func LoadRoutes(router *httprouter.Router) error {
    session := SessionModule{}

    if err := session.Init(); err == nil {
        router.GET(`/api/session/workspaces`, func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
            if workspaces, err := session.GetAllWorkspaces(); err == nil {
                util.Respond(w, http.StatusOK, workspaces, nil)
            }else{
                util.Respond(w, http.StatusInternalServerError, nil, err)
            }
        })

        router.GET(`/api/session/workspaces/current`, func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
            if workspaces, err := session.GetAllWorkspaces(); err == nil {
                for _, workspace := range workspaces {
                    if workspace.IsCurrent {
                        util.Respond(w, http.StatusOK, workspace, nil)
                        return
                    }
                }

                util.Respond(w, http.StatusNotFound, nil, fmt.Errorf("Current workspace not found"))
            }else{
                util.Respond(w, http.StatusInternalServerError, nil, err)
            }
        })

        router.GET(`/api/session/windows`, func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
            if windows, err := session.GetAllWindows(); err == nil {
                for i, _ := range windows {
                    windows[i].IconUri = fmt.Sprintf("/api/session/windows/%d/icon", windows[i].ID)
                }

                util.Respond(w, http.StatusOK, windows, nil)
            }else{
                util.Respond(w, http.StatusInternalServerError, nil, err)
            }
        })

        router.GET(`/api/session/windows/:id`, func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
            if window, err := session.GetWindow(params.ByName(`id`)); err == nil {
                window.IconUri = fmt.Sprintf("/api/session/windows/%s/icon", params.ByName(`id`))

                util.Respond(w, http.StatusOK, window, nil)
            }else{
                util.Respond(w, http.StatusInternalServerError, nil, err)
            }
        })

        router.GET(`/api/session/windows/:id/icon`, func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
            var buffer bytes.Buffer
            width  := uint(16)
            height := uint(16)

            if w := req.URL.Query().Get(`w`); w != `` {
                if value, err := stringutil.ConvertToInteger(w); err == nil {
                    width = uint(value)
                }
            }

            if h := req.URL.Query().Get(`h`); h != `` {
                if value, err := stringutil.ConvertToInteger(h); err == nil {
                    height = uint(value)
                }
            }

            if height != width {
                height = width
            }

            if err := session.WriteWindowIcon(params.ByName(`id`), width, height, &buffer); err == nil {
                w.Header().Set(`Content-Type`, `image/png`)
                w.Write(buffer.Bytes())
                return
            }else{
                util.Respond(w, http.StatusInternalServerError, nil, err)
            }
        })

        router.GET(`/api/session/windows/:id/image`, func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
            var buffer bytes.Buffer

            if err := session.WriteWindowImage(params.ByName(`id`), &buffer); err == nil {
                w.Header().Set(`Content-Type`, `image/png`)
                w.Write(buffer.Bytes())
                return
            }else{
                util.Respond(w, http.StatusInternalServerError, nil, err)
            }
        })

        router.PUT(`/api/session/windows/:id/do/:action`, func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
            var err error
            id := params.ByName(`id`)

            switch params.ByName(`action`) {
            case `maximize`:
                err = session.MaximizeWindow(id)

            case `max-x`:
                err = session.MaximizeWindowHorizontal(id)

            case `max-y`:
                err = session.MaximizeWindowVertical(id)

            case `minimize`:
                err = session.MinimizeWindow(id)

            case `restore`:
                err = session.RestoreWindow(id)

            case `hide`:
                err = session.HideWindow(id)

            case `show`:
                err = session.ShowWindow(id)

            case `raise`:
                err = session.RaiseWindow(id)

            default:
                util.Respond(w, http.StatusBadRequest, nil, fmt.Errorf("Unknown action '%s'", params.ByName(`action`)))
                return
            }

            if err == nil {
                util.Respond(w, http.StatusAccepted, nil, nil)
            }else{
                util.Respond(w, http.StatusInternalServerError, nil, err)
            }
        })

        router.PUT(`/api/session/windows/:id/move/:x/:y`, func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
            id := params.ByName(`id`)
            var x, y int

            if value, err := stringutil.ConvertToInteger(params.ByName(`x`)); err == nil {
                x = int(value)
            }else{
                util.Respond(w, http.StatusBadRequest, nil, err)
                return
            }

            if value, err := stringutil.ConvertToInteger(params.ByName(`y`)); err == nil {
                y = int(value)
            }else{
                util.Respond(w, http.StatusBadRequest, nil, err)
                return
            }

            if err := session.MoveWindow(id, x, y); err == nil {
                util.Respond(w, http.StatusAccepted, nil, nil)
            }else{
                util.Respond(w, http.StatusInternalServerError, nil, err)
            }
        })

        router.PUT(`/api/session/windows/:id/resize/:x/:y`, func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
            id := params.ByName(`id`)
            var width, height uint

            if value, err := stringutil.ConvertToInteger(params.ByName(`width`)); err == nil {
                width = uint(value)
            }else{
                util.Respond(w, http.StatusBadRequest, nil, err)
                return
            }

            if value, err := stringutil.ConvertToInteger(params.ByName(`height`)); err == nil {
                height = uint(value)
            }else{
                util.Respond(w, http.StatusBadRequest, nil, err)
                return
            }

            if err := session.ResizeWindow(id, width, height); err == nil {
                util.Respond(w, http.StatusAccepted, nil, nil)
            }else{
                util.Respond(w, http.StatusInternalServerError, nil, err)
            }
        })

        router.GET(`/api/session/applications`, func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
            util.Respond(w, http.StatusOK, session.GetAppList(), nil)
        })

        router.GET(`/api/session/applications/:name`, func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
            if app, err := session.GetAppByName(params.ByName(`name`)); err == nil {
                util.Respond(w, http.StatusOK, app, nil)
            }else{
                util.Respond(w, http.StatusNotFound, nil, err)
            }
        })

        // router.GET(`/api/session/applications/find/:pattern`, func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
        // })

        router.PUT(`/api/session/applications/:name/launch`, func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
            appName := strings.Replace(params.ByName(`name`) , `%20`, " ", -1)

            if err := session.LaunchAppByName(appName); err == nil {
                util.Respond(w, http.StatusAccepted, nil, nil)
            }else{
                util.Respond(w, http.StatusNotFound, nil, err)
            }
        })
    }else{
        return err
    }

    return nil
}

func (self *SessionModule) Init() (err error) {
    self.X, err = xgbutil.NewConn()
    return
}

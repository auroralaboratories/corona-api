package session

import (
    "bytes"
    "fmt"
    "io/ioutil"
    "net/http"
    "sort"
    "strings"
    "github.com/BurntSushi/xgbutil"
    "github.com/codegangsta/cli"
    "github.com/julienschmidt/httprouter"
    "github.com/shutterstock/go-stockutil/sliceutil"
    "github.com/shutterstock/go-stockutil/stringutil"
    "github.com/auroralaboratories/corona-api/util"
    "github.com/auroralaboratories/corona-api/modules"
    "github.com/auroralaboratories/freedesktop/desktop"
    "github.com/auroralaboratories/freedesktop/icons"
)

type SessionModule struct {
    modules.BaseModule

    X            *xgbutil.XUtil
    Applications *desktop.EntrySet
    Themeset     *icons.Themeset
}

func RegisterSubcommands() []cli.Command {
    return []cli.Command{}
}

func (self *SessionModule) LoadRoutes(router *httprouter.Router) error {
    router.GET(`/api/session/workspaces`, func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
        if workspaces, err := self.GetAllWorkspaces(); err == nil {
            util.Respond(w, http.StatusOK, workspaces, nil)
        }else{
            util.Respond(w, http.StatusInternalServerError, nil, err)
        }
    })

    router.GET(`/api/session/workspaces/current`, func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
        if workspaces, err := self.GetAllWorkspaces(); err == nil {
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
        if windows, err := self.GetAllWindows(); err == nil {
            for i, _ := range windows {
                windows[i].IconUri = fmt.Sprintf("/api/session/windows/%d/icon", windows[i].ID)
            }

            util.Respond(w, http.StatusOK, windows, nil)
        }else{
            util.Respond(w, http.StatusInternalServerError, nil, err)
        }
    })

    router.GET(`/api/session/windows/:id`, func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
        if window, err := self.GetWindow(params.ByName(`id`)); err == nil {
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

        if err := self.WriteWindowIcon(params.ByName(`id`), width, height, &buffer); err == nil {
            w.Header().Set(`Content-Type`, `image/png`)
            w.Write(buffer.Bytes())
            return
        }else{
            util.Respond(w, http.StatusInternalServerError, nil, err)
        }
    })

    router.GET(`/api/session/windows/:id/image`, func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
        var buffer bytes.Buffer

        if err := self.WriteWindowImage(params.ByName(`id`), &buffer); err == nil {
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
            err = self.MaximizeWindow(id)

        case `max-x`:
            err = self.MaximizeWindowHorizontal(id)

        case `max-y`:
            err = self.MaximizeWindowVertical(id)

        case `minimize`:
            err = self.MinimizeWindow(id)

        case `restore`:
            err = self.RestoreWindow(id)

        case `hide`:
            err = self.HideWindow(id)

        case `show`:
            err = self.ShowWindow(id)

        case `raise`:
            err = self.RaiseWindow(id)

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

        if err := self.MoveWindow(id, x, y); err == nil {
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

        if err := self.ResizeWindow(id, width, height); err == nil {
            util.Respond(w, http.StatusAccepted, nil, nil)
        }else{
            util.Respond(w, http.StatusInternalServerError, nil, err)
        }
    })

    router.GET(`/api/session/applications`, func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
        keys := make([]string, 0)

        for key, _ := range self.Applications.Entries {
            keys = append(keys, key)
        }

        sort.Strings(keys)

        util.Respond(w, http.StatusOK, keys, nil)
    })

    router.GET(`/api/session/applications/:name`, func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
        key := params.ByName(`name`)

        if app, ok := self.Applications.Entries[key]; ok {
            util.Respond(w, http.StatusOK, app, nil)
        }else{
            util.Respond(w, http.StatusNotFound, nil, fmt.Errorf("Could not locate application '%s'", key))
        }
    })


    router.GET(`/api/session/icons/list/:type`, func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
        var filterMinSize, filterMaxSize int

        rv := make([]string, 0)
        listType := params.ByName(`type`)


    //  filters
        filterThemes           := strings.Split(req.URL.Query().Get(`themes`), `,`)
        filterIconContextTypes := strings.Split(req.URL.Query().Get(`contexts`), `,`)
        filterIconFileTypes    := strings.Split(req.URL.Query().Get(`filetypes`), `,`)
        filterMinSizeS         := req.URL.Query().Get(`minsize`)
        filterMaxSizeS         := req.URL.Query().Get(`maxsize`)
        filterIsScalable       := req.URL.Query().Get(`scalable`)

        if filterMinSizeS != `` {
            if v, err := stringutil.ConvertToInteger(filterMinSizeS); err == nil {
                filterMinSize = int(v)
            }else{
                util.Respond(w, http.StatusBadRequest, nil, err)
                return
            }
        }

        if filterMaxSizeS != `` {
            if v, err := stringutil.ConvertToInteger(filterMaxSizeS); err == nil {
                filterMaxSize = int(v)
            }else{
                util.Respond(w, http.StatusBadRequest, nil, err)
                return
            }
        }

        for _, theme := range self.Themeset.Themes {
            if len(filterThemes) > 0 && filterThemes[0] != `` {
                if !sliceutil.ContainsString(filterThemes, strings.ToLower(theme.InternalName)) {
                    inInherited := false

                    for _, inheritedThemeName := range theme.Inherits {
                        if sliceutil.ContainsString(filterThemes, strings.ToLower(inheritedThemeName)) {
                            inInherited = true
                            break
                        }
                    }

                    if !inInherited {
                        continue
                    }
                }
            }

            switch listType {
            case `themes`:
                if !sliceutil.ContainsString(rv, theme.Name) {
                    rv = append(rv, theme.InternalName)
                }
            default:
                for _, icon := range theme.Icons {
                //  filter context types
                    if len(filterIconContextTypes) > 0 && filterIconContextTypes[0] != `` {
                        if !sliceutil.ContainsString(filterIconContextTypes, strings.ToLower(icon.Context.Name)) {
                            continue
                        }
                    }

                //  filter icon filetypes
                    if len(filterIconFileTypes) > 0 && filterIconFileTypes[0] != `` {
                        if !sliceutil.ContainsString(filterIconFileTypes, icon.Type) {
                            continue
                        }
                    }

                //  filter icon size contraints
                    if filterMinSize > 0 && icon.Context.MinSize < filterMinSize {
                        continue
                    }

                    if filterMaxSize > 0 && icon.Context.MaxSize > filterMaxSize {
                        continue
                    }

                //  filter for scalable/non-scalable icons
                    if filterIsScalable == `true` && icon.Context.Type != icons.IconContextScalable {
                        continue
                    }else if filterIsScalable == `false` && icon.Context.Type == icons.IconContextScalable {
                        continue
                    }

                    var value string

                    switch listType {
                    case `names`:
                        value = icon.Name
                    case `contexts`:
                        value = strings.ToLower(icon.Context.Name)
                    case `display-names`:
                        value = icon.DisplayName
                    default:
                        util.Respond(w, http.StatusBadRequest, nil, fmt.Errorf("Unrecognized list type '%s'", listType))
                        return
                    }

                    if value != `` {
                        if !sliceutil.ContainsString(rv, value) {
                            rv = append(rv, value)
                        }
                    }
                }
            }
        }

        sort.Strings(rv)

        util.Respond(w, http.StatusOK, rv, nil)
    })

    router.GET(`/api/session/icons/view/:name/size/:size`, func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
        var iconSize int

        iconNames := strings.Split(params.ByName(`name`), `,`)
        iconSizeS := params.ByName(`size`)

        if v, err := stringutil.ConvertToInteger(iconSizeS); err == nil {
            iconSize = int(v)

            if icon, ok := self.Themeset.FindIconViaTheme(`Faenza-Dark`, iconNames, iconSize); ok {
                var contentType string

                switch icon.Type {
                case `png`:
                    contentType = `image/png`
                case `svg`:
                    contentType = `image/svg+xml`
                default:
                    util.Respond(w, http.StatusBadRequest, nil, fmt.Errorf("Unsupported icon type '%s'", icon.Type))
                    return
                }

                defer icon.Close()

                if data, err := ioutil.ReadAll(icon); err == nil {
                    w.Header().Set(`Content-Type`, contentType)
                    w.Write(data)
                }else{
                    util.Respond(w, http.StatusBadRequest, nil, err)
                }
            }else{
                util.Respond(w, http.StatusNotFound, nil, fmt.Errorf("Could not locate icon"))
            }
        }else{
            util.Respond(w, http.StatusBadRequest, nil, err)
        }
    })

    // router.GET(`/api/session/applications/find/:pattern`, func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
    // })

    router.PUT(`/api/session/applications/:name/launch`, func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
        if err := self.Applications.LaunchEntry(params.ByName(`name`)); err == nil {
            util.Respond(w, http.StatusAccepted, nil, nil)
        }else{
            util.Respond(w, http.StatusNotFound, nil, err)
        }

    })

    return nil
}

func (self *SessionModule) Init() (err error) {
    self.X, err = xgbutil.NewConn()
    self.Applications = desktop.NewEntrySet()
    self.Themeset = icons.NewThemeset()

    if err := self.Applications.Refresh(); err != nil {
        return err
    }


    if err := self.Themeset.Load(); err != nil {
        return err
    }

    return nil
}

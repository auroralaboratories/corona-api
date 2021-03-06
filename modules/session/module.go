package session

import (
	"bytes"
	"fmt"
	"github.com/BurntSushi/xgbutil"
	"github.com/auroralaboratories/corona-api/modules"
	"github.com/auroralaboratories/corona-api/util"
	"github.com/auroralaboratories/freedesktop/desktop"
	"github.com/auroralaboratories/freedesktop/icons"
	"github.com/ghetzel/cli"
	"github.com/husobee/vestigo"
	"github.com/shutterstock/go-stockutil/sliceutil"
	"github.com/shutterstock/go-stockutil/stringutil"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
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

func (self *SessionModule) LoadRoutes(router *vestigo.Router) error {
	router.Get(`/api/session/workspaces`, func(w http.ResponseWriter, req *http.Request) {
		if workspaces, err := self.GetAllWorkspaces(); err == nil {
			util.Respond(w, http.StatusOK, workspaces, nil)
		} else {
			util.Respond(w, http.StatusInternalServerError, nil, err)
		}
	})

	router.Get(`/api/session/workspaces/current`, func(w http.ResponseWriter, req *http.Request) {
		if workspaces, err := self.GetAllWorkspaces(); err == nil {
			for _, workspace := range workspaces {
				if workspace.IsCurrent {
					util.Respond(w, http.StatusOK, workspace, nil)
					return
				}
			}

			util.Respond(w, http.StatusNotFound, nil, fmt.Errorf("Current workspace not found"))
		} else {
			util.Respond(w, http.StatusInternalServerError, nil, err)
		}
	})

	router.Get(`/api/session/windows`, func(w http.ResponseWriter, req *http.Request) {
		if windows, err := self.GetAllWindows(); err == nil {
			for i := range windows {
				windows[i].IconUri = fmt.Sprintf("/api/session/windows/%d/icon", windows[i].ID)
			}

			util.Respond(w, http.StatusOK, windows, nil)
		} else {
			util.Respond(w, http.StatusInternalServerError, nil, err)
		}
	})

	router.Get(`/api/session/windows/:id`, func(w http.ResponseWriter, req *http.Request) {
		id := vestigo.Param(req, `id`)

		if window, err := self.GetWindow(id); err == nil {
			window.IconUri = fmt.Sprintf("/api/session/windows/%s/icon", id)

			util.Respond(w, http.StatusOK, window, nil)
		} else {
			util.Respond(w, http.StatusInternalServerError, nil, err)
		}
	})

	router.Get(`/api/session/windows/:id/icon`, func(w http.ResponseWriter, req *http.Request) {
		var buffer bytes.Buffer
		width := uint(16)
		height := uint(16)
		id := vestigo.Param(req, `id`)

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

		if err := self.WriteWindowIcon(id, width, height, &buffer); err == nil {
			w.Header().Set(`Content-Type`, `image/png`)
			w.Write(buffer.Bytes())
			return
		} else {
			util.Respond(w, http.StatusInternalServerError, nil, err)
		}
	})

	router.Get(`/api/session/windows/:id/image`, func(w http.ResponseWriter, req *http.Request) {
		var buffer bytes.Buffer
		id := vestigo.Param(req, `id`)

		if err := self.WriteWindowImage(id, &buffer); err == nil {
			w.Header().Set(`Content-Type`, `image/png`)
			w.Write(buffer.Bytes())
			return
		} else {
			util.Respond(w, http.StatusInternalServerError, nil, err)
		}
	})

	router.Put(`/api/session/windows/:id/do/:action`, func(w http.ResponseWriter, req *http.Request) {
		var err error
		id := vestigo.Param(req, `id`)
		action := vestigo.Param(req, `action`)

		switch action {
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
			util.Respond(w, http.StatusBadRequest, nil, fmt.Errorf("Unknown action '%s'", action))
			return
		}

		if err == nil {
			util.Respond(w, http.StatusAccepted, nil, nil)
		} else {
			util.Respond(w, http.StatusInternalServerError, nil, err)
		}
	})

	router.Put(`/api/session/windows/:id/move/:x/:y`, func(w http.ResponseWriter, req *http.Request) {
		id := vestigo.Param(req, `id`)

		var x, y int

		if value, err := stringutil.ConvertToInteger(vestigo.Param(req, `x`)); err == nil {
			x = int(value)
		} else {
			util.Respond(w, http.StatusBadRequest, nil, err)
			return
		}

		if value, err := stringutil.ConvertToInteger(vestigo.Param(req, `y`)); err == nil {
			y = int(value)
		} else {
			util.Respond(w, http.StatusBadRequest, nil, err)
			return
		}

		if err := self.MoveWindow(id, x, y); err == nil {
			util.Respond(w, http.StatusAccepted, nil, nil)
		} else {
			util.Respond(w, http.StatusInternalServerError, nil, err)
		}
	})

	router.Put(`/api/session/windows/:id/resize/:x/:y`, func(w http.ResponseWriter, req *http.Request) {
		id := vestigo.Param(req, `id`)
		var width, height uint

		if value, err := stringutil.ConvertToInteger(vestigo.Param(req, `width`)); err == nil {
			width = uint(value)
		} else {
			util.Respond(w, http.StatusBadRequest, nil, err)
			return
		}

		if value, err := stringutil.ConvertToInteger(vestigo.Param(req, `height`)); err == nil {
			height = uint(value)
		} else {
			util.Respond(w, http.StatusBadRequest, nil, err)
			return
		}

		if err := self.ResizeWindow(id, width, height); err == nil {
			util.Respond(w, http.StatusAccepted, nil, nil)
		} else {
			util.Respond(w, http.StatusInternalServerError, nil, err)
		}
	})

	router.Get(`/api/session/applications`, func(w http.ResponseWriter, req *http.Request) {
		keys := make([]string, 0)

		for key := range self.Applications.Entries {
			keys = append(keys, key)
		}

		sort.Strings(keys)

		util.Respond(w, http.StatusOK, keys, nil)
	})

	router.Get(`/api/session/applications/:name`, func(w http.ResponseWriter, req *http.Request) {
		key := vestigo.Param(req, `key`)

		if app, ok := self.Applications.Entries[key]; ok {
			util.Respond(w, http.StatusOK, app, nil)
		} else {
			util.Respond(w, http.StatusNotFound, nil, fmt.Errorf("Could not locate application '%s'", key))
		}
	})

	router.Get(`/api/session/icons/list/:type`, func(w http.ResponseWriter, req *http.Request) {
		var filterMinSize, filterMaxSize int

		rv := make([]string, 0)
		listType := vestigo.Param(req, `type`)

		//  filters
		filterThemes := strings.Split(req.URL.Query().Get(`themes`), `,`)
		filterIconContextTypes := strings.Split(req.URL.Query().Get(`contexts`), `,`)
		filterIconFileTypes := strings.Split(req.URL.Query().Get(`filetypes`), `,`)
		filterMinSizeS := req.URL.Query().Get(`minsize`)
		filterMaxSizeS := req.URL.Query().Get(`maxsize`)
		filterIsScalable := req.URL.Query().Get(`scalable`)

		if filterMinSizeS != `` {
			if v, err := stringutil.ConvertToInteger(filterMinSizeS); err == nil {
				filterMinSize = int(v)
			} else {
				util.Respond(w, http.StatusBadRequest, nil, err)
				return
			}
		}

		if filterMaxSizeS != `` {
			if v, err := stringutil.ConvertToInteger(filterMaxSizeS); err == nil {
				filterMaxSize = int(v)
			} else {
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
					} else if filterIsScalable == `false` && icon.Context.Type == icons.IconContextScalable {
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
					case `sizes`:
						if v, err := stringutil.ToString(icon.Context.Size); err == nil {
							value = v
						}
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

	router.Get(`/api/session/icons/view/:name/size/:size`, func(w http.ResponseWriter, req *http.Request) {
		var iconSize int

		name := vestigo.Param(req, `name`)
		iconNames := strings.Split(name, `,`)
		iconSizeS := vestigo.Param(req, `size`)
		themeName := req.URL.Query().Get(`theme`)

		if themeName == `` {
			themeName = self.Themeset.DefaultTheme
		}

		if v, err := stringutil.ConvertToInteger(iconSizeS); err == nil {
			iconSize = int(v)

			var icon *icons.Icon

			switch req.URL.Query().Get(`mode`) {
			case `hicolor-first`:
				if hiColorIcon, ok := self.Themeset.FindIconViaTheme(`hicolor`, iconNames, iconSize); ok {
					icon = hiColorIcon
				} else if themeIcon, ok := self.Themeset.FindIconViaTheme(themeName, iconNames, iconSize); ok {
					icon = themeIcon
				}
			default:
				if themeIcon, ok := self.Themeset.FindIconViaTheme(themeName, iconNames, iconSize); ok {
					icon = themeIcon
				}
			}

			if icon != nil {
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
				} else {
					util.Respond(w, http.StatusBadRequest, nil, err)
				}
			} else {
				util.Respond(w, http.StatusNotFound, nil, fmt.Errorf("Could not locate icon"))
			}
		} else {
			util.Respond(w, http.StatusBadRequest, nil, err)
		}
	})

	// router.Get(`/api/session/applications/find/:pattern`, func(w http.ResponseWriter, req *http.Request) {
	// })

	router.Put(`/api/session/applications/:name/launch`, func(w http.ResponseWriter, req *http.Request) {
		name := vestigo.Param(req, `name`)

		if err := self.Applications.LaunchEntry(name); err == nil {
			util.Respond(w, http.StatusAccepted, nil, nil)
		} else {
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

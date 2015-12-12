package session

import (
    "fmt"
    "regexp"
    "strings"
    "os/user"
    "os/exec"
    "io/ioutil"
    "github.com/vaughan0/go-ini"
    "github.com/shutterstock/go-stockutil/stringutil"
)

type Application struct {
    Type            string    `json:"type"`
    Version         string    `json:"version,omitempty"`
    Name            string    `json:"name"`
    GenericName     string    `json:"generic_name,omitempty"`
    Comment         string    `json:"comment,omitempty"`
    Icon            string    `json:"icon,omitempty"`
    TryExec         string    `json:"tryexec,omitempty"`
    Exec            string    `json:"exec,omitempty"`
    Path            string    `json:"path,omitempty"`
    StartupWMClass  string    `json:"startup_wm_class,omitempty"`
    URL             string    `json:"url,omitempty"`
    DesktopfilePath string    `json:"desktop_file_path,omitempty"`
    NoDisplay       bool      `json:"no_display,omitempty"`
    StartupNotify   bool      `json:"startup_notify,omitempty"`
    Hidden          bool      `json:"hidden,omitempty"`
    DBusActivatable bool      `json:"dbus_activateable,omitempty"`
    Terminal        bool      `json:"terminal,omitempty"`
    OnlyShowIn      []string  `json:"only_show_in,omitempty"`
    NotShowIn       []string  `json:"not_shown_in,omitempty"`
    Actions         []string  `json:"actions,omitempty"`
    MimeType        []string  `json:"mimetypes,omitempty"`
    Categories      []string  `json:"categories,omitempty"`
    Implements      []string  `json:"implements,omitempty"`
    Keywords        []string  `json:"keywords,omitempty"`
}

func (self *SessionModule) SearchAppByName(pattern string) []Application {
    results := make([]Application, 0)

    for _, app := range self.GetAppList() {
        if match, err := regexp.MatchString(pattern, strings.ToLower(app.Name)); err == nil && match {
            results = append(results, app)
        }
    }

    return results
}

func (self *SessionModule) SearchAppByExec(pattern string) []Application {
    results := make([]Application, 0)

    for _, app := range self.GetAppList(){
        if match, err := regexp.MatchString(pattern, strings.ToLower(app.Exec)); err == nil && match {
            results = append(results, app)
        }
    }

    return results
}

func (self *SessionModule) LaunchAppByName(name string) error {
    if app, err := self.GetAppByName(name); err == nil {
        cmd := strings.Split(app.Exec, " ")

        return self.Launch(cmd)
    }else{
        return err
    }
}

func (self *SessionModule) Launch(command []string) error {
    exe := command[0]
    args := command[1:len(command)]

    cmd := exec.Command(exe, args...)
    return cmd.Start()
}

func (self *SessionModule) GetAppByName(name string) (Application, error) {
    for _, app := range self.GetAppList(){
        if strings.ToLower(app.Name) == strings.ToLower(name){
            return app, nil
        }
    }

    return Application{}, fmt.Errorf("Cannot find application '%s'", name)
}

func (self *SessionModule) GetAppList() []Application {
    applist := make([]Application, 0)

    if files, err := self.getFileList(self.getPaths()); err == nil {
        for _, i := range files {
            if file, err := ini.LoadFile(i); err == nil {
                app                 := Application{}
                app.DesktopfilePath  = i
                app.Type             = file["Desktop Entry"]["Type"]
                app.Version          = file["Desktop Entry"]["Version"]
                app.Name             = file["Desktop Entry"]["Name"]
                app.GenericName      = file["Desktop Entry"]["GenericName"]
                app.Comment          = file["Desktop Entry"]["Comment"]
                app.Icon             = file["Desktop Entry"]["Icon"]
                app.TryExec          = file["Desktop Entry"]["TryExec"]
                app.Exec             = file["Desktop Entry"]["Exec"]
                app.Path             = file["Desktop Entry"]["Path"]
                app.StartupWMClass   = file["Desktop Entry"]["StartupWMClass"]
                app.URL              = file["Desktop Entry"]["URL"]
                app.DesktopfilePath  = file["Desktop Entry"]["DesktopfilePath"]

                if v, err := stringutil.ConvertToBool(file["Desktop Entry"]["NoDisplay"]); err == nil {
                    app.NoDisplay = v
                }

                if v, err := stringutil.ConvertToBool(file["Desktop Entry"]["StartupNotify"]); err == nil {
                    app.StartupNotify = v
                }

                if v, err := stringutil.ConvertToBool(file["Desktop Entry"]["Hidden"]); err == nil {
                    app.Hidden = v
                }

                if v, err := stringutil.ConvertToBool(file["Desktop Entry"]["DBusActivatable"]); err == nil {
                    app.DBusActivatable = v
                }

                if v, err := stringutil.ConvertToBool(file["Desktop Entry"]["Terminal"]); err == nil {
                    app.Terminal = v
                }

                app.OnlyShowIn      = self.getDesktopStringList(file["Desktop Entry"]["OnlyShowIn"])
                app.NotShowIn       = self.getDesktopStringList(file["Desktop Entry"]["NotShowIn"])
                app.Actions         = self.getDesktopStringList(file["Desktop Entry"]["Actions"])
                app.MimeType        = self.getDesktopStringList(file["Desktop Entry"]["MimeType"])
                app.Categories      = self.getDesktopStringList(file["Desktop Entry"]["Categories"])
                app.Implements      = self.getDesktopStringList(file["Desktop Entry"]["Implements"])
                app.Keywords        = self.getDesktopStringList(file["Desktop Entry"]["Keywords"])

                applist = append(applist, app)
            }
        }
    }

    return applist
}

func (self *SessionModule) getPaths() []string {
    if usr, err := user.Current(); err == nil {
        homedir := usr.HomeDir
        return []string{ "/usr/share/applications/", homedir + "/.local/share/applications/" }
    }else{
        return nil
    }
}

func (self *SessionModule) getFileList(paths []string) ([]string, error) {
    var files = make([]string, 0)

    for path := range paths {
        if filelist, err := ioutil.ReadDir(paths[path]); err == nil {
            for _, f := range filelist {
                files = append(files, paths[path]+f.Name())
            }
        }else{
            return files, err
        }
    }

    return files, nil
}

func (self *SessionModule) getDesktopStringList(in string) []string {
    rv := make([]string, 0)

    for _, value := range strings.Split(in, `;`) {
        value = strings.TrimSpace(value)

        if value != `` {
            rv = append(rv, value)
        }
    }

    return rv
}

package main

import (
  "log"
  "regexp"
  "strings"
  "os/user"
  "os/exec"
  "io/ioutil"
  "github.com/vaughan0/go-ini"
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
  OnlyShowIn      []string    `json:"only_show_in,omitempty"`
  NotShowIn       []string    `json:"not_shown_in,omitempty"`
  Actions         []string    `json:"actions,omitempty"`
  MimeType        []string    `json:"mimetypes,omitempty"`
  Categories      []string    `json:"categories,omitempty"`
  Implements      []string    `json:"implements,omitempty"`
  Keywords        []string    `json:"keywords,omitempty"`
}

func (self *SessionPlugin) SearchAppByName(pattern string) []Application{
  results := []Application{}
  for _, app := range self.GetAppList(){
    match, _ := regexp.MatchString(pattern, strings.ToLower(app.Name))
    if match {
      results = append(results, app)
    }
  }
  return results
}

func (self *SessionPlugin) SearchAppByExec(pattern string) []Application{
  results := []Application{}
  for _, app := range self.GetAppList(){
    match, _ := regexp.MatchString(pattern, strings.ToLower(app.Exec))
    if match {
      results = append(results, app)
    }
  }
  return results
}

func (self *SessionPlugin) LaunchAppByName(name string){
  cmd := strings.Split(self.GetAppByName(name).Exec, " ")[0]
  self.Launch(cmd)
}

func (self *SessionPlugin) Launch(command string){
  cmd := exec.Command(command)
  // cmd.SysProcAttr.Setpgid = true
  err := cmd.Start()
  if err != nil {
    log.Fatal(err)
  }
}

func (self *SessionPlugin) GetAppByName(name string) Application{
  for _, app := range self.GetAppList(){
    if app.Name == name{
      return app
    }
  }
  return Application{}
}

func (self *SessionPlugin) GetAppList() []Application{
  files   := self.getFileList(self.getPaths())
  applist := []Application{}

  for _, i := range files {
    file, err := ini.LoadFile(i)
    if err != nil {
      log.Printf("%s - %s", i, err)
      continue
    }
    app := Application{}
    app.DesktopfilePath = i
    app.Type            = file["Desktop Entry"]["Type"]
    app.Version         = file["Desktop Entry"]["Version"]
    app.Name            = file["Desktop Entry"]["Name"]
    app.GenericName     = file["Desktop Entry"]["GenericName"]
    app.Comment         = file["Desktop Entry"]["Comment"]
    app.Icon            = file["Desktop Entry"]["Icon"]
    app.TryExec         = file["Desktop Entry"]["TryExec"]
    app.Exec            = file["Desktop Entry"]["Exec"]
    app.Path            = file["Desktop Entry"]["Path"]
    app.StartupWMClass  = file["Desktop Entry"]["StartupWMClass"]
    app.URL             = file["Desktop Entry"]["URL"]
    app.DesktopfilePath = file["Desktop Entry"]["DesktopfilePath"]
    app.NoDisplay       = Stob(file["Desktop Entry"]["NoDisplay"])
    app.StartupNotify   = Stob(file["Desktop Entry"]["StartupNotify"])
    app.Hidden          = Stob(file["Desktop Entry"]["Hidden"])
    app.DBusActivatable = Stob(file["Desktop Entry"]["DBusActivatable"])
    app.Terminal        = Stob(file["Desktop Entry"]["Terminal"])
    app.OnlyShowIn      = Stosl(file["Desktop Entry"]["OnlyShowIn"])
    app.NotShowIn       = Stosl(file["Desktop Entry"]["NotShowIn"])
    app.Actions         = Stosl(file["Desktop Entry"]["Actions"])
    app.MimeType        = Stosl(file["Desktop Entry"]["MimeType"])
    app.Categories      = Stosl(file["Desktop Entry"]["Categories"])
    app.Implements      = Stosl(file["Desktop Entry"]["Implements"])
    app.Keywords        = Stosl(file["Desktop Entry"]["Keywords"])

    applist = append(applist, app)
  }
  return applist
}

func (self *SessionPlugin) getPaths() []string{
  usr, _ := user.Current()
  homedir := usr.HomeDir
  return []string{"/usr/share/applications/", homedir + "/.local/share/applications/"}
}

func (self *SessionPlugin) getFileList(paths []string) []string{
  var files = make([]string, 0)
  for path := range paths {
    filelist, err := ioutil.ReadDir(paths[path])
    if err != nil {
      log.Fatal(err)
    }
    for _, f := range filelist {
      files = append(files, paths[path]+f.Name())
    }
  }
  return files
}

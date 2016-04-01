package desktop

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/shutterstock/go-stockutil/stringutil"
	"github.com/vaughan0/go-ini"
)

type EntrySet struct {
	Entries map[string]Entry
}

func NewEntrySet() *EntrySet {
	rv := &EntrySet{
		Entries: make(map[string]Entry),
	}

	return rv
}

func (self *EntrySet) Refresh() error {
	if files, err := self.getFileList(self.getPaths()); err == nil {
		for _, filename := range files {
			if file, err := ini.LoadFile(filename); err == nil {
				entryKey := strings.Replace(strings.ToLower(strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))), `_`, `-`, -1)

				entry := Entry{}
				entry.Key = entryKey
				entry.DesktopfilePath = filename
				entry.Type = file[`Desktop Entry`][`Type`]
				entry.Version = file[`Desktop Entry`][`Version`]
				entry.Name = file[`Desktop Entry`][`Name`]
				entry.GenericName = file[`Desktop Entry`][`GenericName`]
				entry.Comment = file[`Desktop Entry`][`Comment`]
				entry.Icon = file[`Desktop Entry`][`Icon`]
				entry.TryExec = file[`Desktop Entry`][`TryExec`]
				entry.Exec = file[`Desktop Entry`][`Exec`]
				entry.Path = file[`Desktop Entry`][`Path`]
				entry.StartupWMClass = file[`Desktop Entry`][`StartupWMClass`]
				entry.URL = file[`Desktop Entry`][`URL`]
				entry.DesktopfilePath = file[`Desktop Entry`][`DesktopfilePath`]

				if v, err := stringutil.ConvertToBool(file[`Desktop Entry`][`NoDisplay`]); err == nil {
					entry.NoDisplay = v
				}

				if v, err := stringutil.ConvertToBool(file[`Desktop Entry`][`StartupNotify`]); err == nil {
					entry.StartupNotify = v
				}

				if v, err := stringutil.ConvertToBool(file[`Desktop Entry`][`Hidden`]); err == nil {
					entry.Hidden = v
				}

				if v, err := stringutil.ConvertToBool(file[`Desktop Entry`][`DBusActivatable`]); err == nil {
					entry.DBusActivatable = v
				}

				if v, err := stringutil.ConvertToBool(file[`Desktop Entry`][`Terminal`]); err == nil {
					entry.Terminal = v
				}

				entry.OnlyShowIn = self.getDesktopStringList(file[`Desktop Entry`][`OnlyShowIn`])
				entry.NotShowIn = self.getDesktopStringList(file[`Desktop Entry`][`NotShowIn`])
				entry.Actions = self.getDesktopStringList(file[`Desktop Entry`][`Actions`])
				entry.MimeType = self.getDesktopStringList(file[`Desktop Entry`][`MimeType`])
				entry.Categories = self.getDesktopStringList(file[`Desktop Entry`][`Categories`])
				entry.Implements = self.getDesktopStringList(file[`Desktop Entry`][`Implements`])
				entry.Keywords = self.getDesktopStringList(file[`Desktop Entry`][`Keywords`])

				self.Entries[entry.Key] = entry
			} else {
				return err
			}
		}
	} else {
		return err
	}

	return nil
}

func (self *EntrySet) SearchByName(pattern string) []Entry {
	results := make([]Entry, 0)

	for _, entry := range self.Entries {
		if match, err := regexp.MatchString(pattern, strings.ToLower(entry.Name)); err == nil && match {
			results = append(results, entry)
		}
	}

	return results
}

func (self *EntrySet) SearchByExec(pattern string) []Entry {
	results := make([]Entry, 0)

	for _, entry := range self.Entries {
		if match, err := regexp.MatchString(pattern, strings.ToLower(entry.Exec)); err == nil && match {
			results = append(results, entry)
		}
	}

	return results
}

func (self *EntrySet) LaunchEntry(key string) error {
	if entry, ok := self.Entries[key]; ok {
		return entry.Launch()
	} else {
		return fmt.Errorf("Could not locate application '%s'", key)
	}
}

func (self *EntrySet) getPaths() []string {
	rv := []string{`/usr/share/applications`}

	if usr, err := user.Current(); err == nil {
		rv = append(rv, fmt.Sprintf("%s/.local/share/applications", usr.HomeDir))
	}

	return rv
}

func (self *EntrySet) getFileList(paths []string) ([]string, error) {
	var files = make([]string, 0)

	for _, path := range paths {
		filepath.Walk(path, func(filename string, info os.FileInfo, err error) error {
			if err == nil {
				if info.Mode().IsRegular() {
					if filepath.Ext(filename) == `.desktop` {
						if absPath, err := filepath.Abs(filename); err == nil {
							files = append(files, absPath)
						} else {
							return err
						}
					}
				} else if filename != path && info.IsDir() {
					return filepath.SkipDir
				}
			}

			return nil
		})
	}

	return files, nil
}

func (self *EntrySet) getDesktopStringList(in string) []string {
	rv := make([]string, 0)

	for _, value := range strings.Split(in, `;`) {
		value = strings.TrimSpace(value)

		if value != `` {
			rv = append(rv, value)
		}
	}

	return rv
}

package icons

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	// "log"

	"github.com/shutterstock/go-stockutil/stringutil"
	"github.com/vaughan0/go-ini"
)

const (
	DEFAULT_ICONTHEME_INDEX_FILE = `index.theme`
	DEFAULT_ICONTHEME_INHERIT    = `hicolor`
)

type Theme struct {
	Comment      string
	Contexts     map[string]IconContext
	Directories  []string
	Example      string
	Hidden       bool
	Icons        []*Icon
	IndexFile    string
	Inherits     []string
	InternalName string
	Name         string
	ThemeDirs    []string

	loadedDef bool
}

func NewTheme(internalName string) *Theme {
	return &Theme{
		Contexts:     make(map[string]IconContext),
		Icons:        make([]*Icon, 0),
		IndexFile:    DEFAULT_ICONTHEME_INDEX_FILE,
		Inherits:     make([]string, 0),
		InternalName: internalName,
		Name:         internalName,
		ThemeDirs:    GetIconThemePaths(),
	}
}

func (self *Theme) Refresh() error {
	if self.InternalName == `` {
		return fmt.Errorf("Cannot refresh theme without an internal name")
	}

	if err := self.refreshThemeDefinition(); err != nil {
		return err
	}

	if err := self.refreshIcons(); err != nil {
		return err
	}

	return nil
}

func (self *Theme) refreshThemeDefinition() error {
	for _, themeDir := range self.ThemeDirs {
		themeIndexFilename := path.Join(themeDir, self.InternalName, self.IndexFile)

		if themeIndex, err := ini.LoadFile(themeIndexFilename); err == nil {
			if config, ok := themeIndex[`Icon Theme`]; ok {
				if v, ok := config[`Name`]; ok {
					self.Name = v
				}

				self.Inherits = strings.Split(config[`Inherits`], `,`)

				if len(self.Inherits) == 0 {
					self.Inherits = []string{DEFAULT_ICONTHEME_INHERIT}
				}

				self.Comment = config[`Comment`]
				self.Example = config[`Example`]

				if v, ok := config[`Hidden`]; ok {
					self.Hidden = (v == `true`)
				}

				if v, ok := config[`Directories`]; ok {
					self.Directories = strings.Split(v, `,`)

					for _, directory := range self.Directories {
						if contextConfig, ok := themeIndex[directory]; ok {
							context := IconContext{
								Subdirectory: directory,
							}

							if v, err := stringutil.ConvertToInteger(contextConfig[`Size`]); err == nil {
								context.Size = int(v)
							}

							if v, ok := contextConfig[`Context`]; ok {
								context.Name = v
							}

							context.MinSize = context.Size
							context.MaxSize = context.Size
							context.Threshold = DEFAULT_ICON_CONTEXT_THRESHOLD

							switch strings.ToLower(contextConfig[`Type`]) {
							case `fixed`:
								context.Type = IconContextFixed

								if minSize, ok := contextConfig[`MinSize`]; ok {
									if v, err := stringutil.ConvertToInteger(minSize); err == nil {
										context.MinSize = int(v)
									}
								}

								if maxSize, ok := contextConfig[`MaxSize`]; ok {
									if v, err := stringutil.ConvertToInteger(maxSize); err == nil {
										context.MaxSize = int(v)
									}
								}

							case `scalable`:
								context.Type = IconContextScalable

							default:
								context.Type = IconContextThreshold

								if threshold, ok := contextConfig[`Threshold`]; ok {
									if v, err := stringutil.ConvertToInteger(threshold); err == nil {
										context.Threshold = int(v)
									}
								}
							}

							self.Contexts[directory] = context
						}
					}
				}

			} else {
				return fmt.Errorf("Cannot load theme at %s: missing [Icon Theme] section", themeIndexFilename)
			}

			self.loadedDef = true
			break
		}
	}

	if !self.loadedDef {
		if err := self.generateThemeDefinition(); err != nil {
			return fmt.Errorf("Unable to find a theme definition file for '%s' in any directories", self.InternalName)
		}
	}

	return nil
}

func (self *Theme) generateThemeDefinition() error {
	for _, themeDir := range self.ThemeDirs {
		//  glob for all sub-subdirs in a candidate directory: <base>/<name>/*/*
		//      e.g.:  /usr/share/icons/hicolor/apps/16
		//
		candidateDir := path.Join(themeDir, self.InternalName)

		if files, err := filepath.Glob(path.Join(candidateDir, `*/*`)); err == nil {
			for _, filename := range files {
				//  verify we have a directory
				if stat, err := os.Stat(filename); err == nil && stat.IsDir() {
					themeSubdir := strings.TrimPrefix(strings.TrimPrefix(filename, candidateDir), `/`)
					self.Directories = append(self.Directories, themeSubdir)

					//  try to figure out the context of this subdirectory
					if parts := strings.SplitN(themeSubdir, `/`, 2); len(parts) == 2 {
						sizeParts := strings.Split(parts[1], `x`)
						context := IconContext{
							Subdirectory: themeSubdir,
							Threshold:    DEFAULT_ICON_CONTEXT_THRESHOLD,
						}

						if sizeParts[0] == `scalable` {
							context.Type = IconContextScalable
							context.Size = DEFAULT_ICON_CONTEXT_SCALABLE_SIZE
							context.MinSize = DEFAULT_ICON_CONTEXT_SCALABLE_MIN
							context.MaxSize = DEFAULT_ICON_CONTEXT_SCALABLE_MAX

						} else if v, err := stringutil.ConvertToInteger(sizeParts[0]); err == nil {
							context.Type = IconContextFixed
							context.Size = int(v)
						}

						self.Contexts[themeSubdir] = context
					}
				}
			}
		}
	}

	//  if we got to this point and still don't have any directories in our list,
	//  give up. we really, really tried....
	if len(self.Directories) == 0 {
		return fmt.Errorf("Unable to generate theme definition")
	}

	//  we're here! set spec defaults
	self.Inherits = []string{DEFAULT_ICONTHEME_INHERIT}

	return nil
}

func (self *Theme) refreshIcons() error {
	//  populate icons
	for _, themeDir := range self.ThemeDirs {
		for _, directory := range self.Directories {
			iconBaseDir := path.Join(themeDir, self.InternalName, directory)

			if stat, err := os.Stat(iconBaseDir); err == nil && stat.IsDir() {
				if files, err := filepath.Glob(path.Join(iconBaseDir, `*.*`)); err == nil {
					for _, iconFilename := range files {
						switch filepath.Ext(iconFilename) {
						case `.png`, `.svg`, `.xpm`:
							icon := NewIcon(iconFilename, self)

							if context, ok := self.Contexts[directory]; ok {
								icon.Context = context
							}

							if err := icon.Refresh(); err == nil {
								self.Icons = append(self.Icons, icon)
							} else {
								return err
							}
						}
					}
				}
			}
		}
	}

	return nil
}

//  Lookup an icon by name and desired size
//
//  The icon lookup mechanism has two global settings, the list of base directories and the internal name of the
//  current theme. Given these we need to specify how to look up an icon file from the icon name and the nominal size.
//
//  The lookup is done first in the current theme, and then recursively in each of the current theme's parents, and
//  finally in the default theme called "hicolor" (implementations may add more default themes before "hicolor", but
//  "hicolor" must be last). As soon as there is an icon of any size that matches in a theme, the search is stopped.
//  Even if there may be an icon with a size closer to the correct one in an inherited theme, we don't want to use it.
//  Doing so may generate an inconsistant change in an icon when you change icon sizes (e.g. zoom in).
//
//  The lookup inside a theme is done in three phases. First all the directories are scanned for an exact match, e.g.
//  one where the allowed size of the icon files match what was looked up. Then all the directories are scanned for any
//  icon that matches the name. If that fails we finally fall back on unthemed icons. If we fail to find any icon at
//  all it is up to the application to pick a good fallback, as the correct choice depends on the context.
//
func (self *Theme) FindIcon(names []string, size int) (*Icon, bool) {
	//  search through a list of given names
	for _, name := range names {
		//  search for a matching icon (size given)
		for _, icon := range self.Icons {
			if icon.IsMatch(name, size) {
				return icon, true
			}
		}

		//  min size starts at MAXINT
		minimalSize := int(^uint(0) >> 1)

		var closestIcon *Icon

		//  search for a matching icon (name only)
		for _, icon := range self.Icons {
			if icon.Name == name {
				if distance := icon.DistanceFromSize(size); distance < minimalSize {
					minimalSize = distance
					closestIcon = icon
				}
			}
		}

		if closestIcon != nil {
			return closestIcon, true
		}
	}

	return nil, false
}

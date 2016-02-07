package icons

import (
    "os"
    "path"
    "path/filepath"

    "github.com/shutterstock/go-stockutil/sliceutil"
)

type Themeset struct {
    Themes     []*Theme
    ThemeDirs  []string

    themeIndex map[string]*Theme
}

func NewThemeset() *Themeset {
    return &Themeset{
        Themes:     make([]*Theme, 0),
        ThemeDirs:  GetIconThemePaths(),

        themeIndex: make(map[string]*Theme),
    }
}

func (self *Themeset) Load() error {
    themeNames := make([]string, 0)

//  build a unique list of theme names from all ThemeDirs
    for _, iconPath := range self.ThemeDirs {
        if themeDirs, err := filepath.Glob(path.Join(iconPath, `*`)); err == nil {
            for _, themeDir := range themeDirs {
                if stat, err := os.Stat(themeDir); err == nil && stat.IsDir() {
                    name := path.Base(themeDir)

                    if !sliceutil.ContainsString(themeNames, name) {
                        themeNames = append(themeNames, name)
                    }
                }
            }
        }
    }

//  assemble and load all themes
    for _, themeName := range themeNames {
        theme := NewTheme(themeName)
        theme.ThemeDirs = self.ThemeDirs

        if err := theme.Refresh(); err == nil {
            self.Themes = append(self.Themes, theme)

            if _, ok := self.themeIndex[theme.InternalName]; !ok {
                self.themeIndex[theme.InternalName] = theme
            }
        }else{
            continue
        }
    }

    return nil
}

func (self *Themeset) GetTheme(name string) (*Theme, bool) {
    if rv, ok := self.themeIndex[name]; ok {
        return rv, true
    }else{
        return nil, false
    }
}

// Locate an icon by name and preferred size, starting with the named theme and resursively
// searching
func (self *Themeset) FindIconViaTheme(themeName string, names []string, size int) (*Icon, bool) {
    if theme, ok := self.themeIndex[themeName]; ok {
        if icon, ok := theme.FindIcon(names, size); ok {
            return icon, true
        }else{
            for _, inheritedThemeName := range theme.Inherits {
                if inheritedThemeName != `` {
                    if icon, ok := self.FindIconViaTheme(inheritedThemeName, names, size); ok {
                        return icon, true
                    }
                }
            }
        }
    }

    return nil, false
}
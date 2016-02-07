package icons

import (
    "os"
    "path"
    "path/filepath"

    "github.com/shutterstock/go-stockutil/sliceutil"
)

const (
    DEFAULT_THEMESET_THEME_INTERNAL_NAME = `hicolor`
    DEFAULT_AUTO_HICOLOR_PENALTY         = 0.0
)

type Themeset struct {
    Themes             []*Theme
    ThemeDirs          []string
    DefaultTheme       string
    AutoHiColorPenalty float64

    themeIndex map[string]*Theme
}

func NewThemeset() *Themeset {
    return &Themeset{
        Themes:             make([]*Theme, 0),
        ThemeDirs:          GetIconThemePaths(),
        DefaultTheme:       DEFAULT_THEMESET_THEME_INTERNAL_NAME,
        AutoHiColorPenalty: DEFAULT_AUTO_HICOLOR_PENALTY,

        themeIndex:         make(map[string]*Theme),
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

// Retrieve a theme by its internal name.
//
func (self *Themeset) GetTheme(name string) (*Theme, bool) {
    if rv, ok := self.themeIndex[name]; ok {
        return rv, true
    }else{
        return nil, false
    }
}


// Locate an icon by name and preferred size using the DefaultTheme as a starting point
// for the search.
//
func (self *Themeset) FindIcon(names []string, size int) (*Icon, bool) {
    return self.FindIconViaTheme(self.DefaultTheme, names, size)
}


// Locate an icon by name and preferred size, starting with the named theme and resursively
// searching.
//
func (self *Themeset) FindIconViaTheme(themeName string, names []string, size int) (*Icon, bool) {
    var hiColorIcon *Icon

//  If desired, a hicolor icon of the given name will be searched for. If the found hicolor
//  icon's (size * penalty) is closer to the preferred size than found theme icon size, use it.
//
//  This is to address the fact that most applications tend to install their branded icons in
//  the hicolor theme because it is guaranteed to be present on all compliant systems.
//
    if themeName != `hicolor` && self.AutoHiColorPenalty > 0.0 {
    //  attempt to locate the icon from the hicolor theme
        if v, ok := self.FindIconViaTheme(`hicolor`, names, size); ok {
            hiColorIcon = v
        }
    }

    if theme, ok := self.themeIndex[themeName]; ok {
        if icon, ok := theme.FindIcon(names, size); ok {
            if hiColorIcon != nil && (self.AutoHiColorPenalty * float64(hiColorIcon.DistanceFromSize(size))) < float64(icon.DistanceFromSize(size)) {
                return hiColorIcon, true
            }else{
                return icon, true
            }
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
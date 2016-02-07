// The icons package implements helper utilities for locating named XDG icons and themes.
//
// An icon theme is a set of icons that share a common look and feel. The user can then select the icon theme that they
// want to use, and all applications will use icons from the theme.
//
// See also: http://standards.freedesktop.org/icon-theme-spec/icon-theme-spec-latest.html
//
package icons

import (
    "os"
    "path"

    "github.com/auroralaboratories/freedesktop"
)

var IconThemesHome   string = os.ExpandEnv("${HOME}/.icons")
var IconThemesLegacy string = `/usr/share/pixmaps`

func GetIconThemePaths() []string {
    rv := make([]string, 0)

    if stat, err := os.Stat(IconThemesHome); err == nil && stat.IsDir() {
        rv = append(rv, IconThemesHome)
    }

    for _, dir := range freedesktop.GetXdgDataPaths() {
        iconsDir := path.Join(dir, `icons`)

        if stat, err := os.Stat(iconsDir); err == nil && stat.IsDir() {
            rv = append(rv, iconsDir)
        }
    }

    if stat, err := os.Stat(IconThemesLegacy); err == nil && stat.IsDir() {
        rv = append(rv, IconThemesLegacy)
    }

    return rv
}

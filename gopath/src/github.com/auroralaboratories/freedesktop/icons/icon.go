package icons

type Icon struct {

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
func FindIcon(name string, size int) (Icon, bool) {
    return Icon{}, false
}


//  Lookup an icon by name and desired size from a specific theme
//
func FindIconInTheme(name string, size int, themeName string) (Icon, bool) {
    // Psuedocode from http://standards.freedesktop.org/icon-theme-spec/icon-theme-spec-latest.html
    //
    // for each subdir in $(theme subdir list) {
    // for each directory in $(basename list) {
    //   for extension in ("png", "svg", "xpm") {
    //     if DirectoryMatchesSize(subdir, size) {
    //       filename = directory/$(themename)/subdir/iconname.extension
    //       if exist filename
    //     return filename
    //     }
    //   }
    // }
    // }
    // minimal_size = MAXINT
    // for each subdir in $(theme subdir list) {
    // for each directory in $(basename list) {
    //   for extension in ("png", "svg", "xpm") {
    //     filename = directory/$(themename)/subdir/iconname.extension
    //     if exist filename and DirectorySizeDistance(subdir, size) < minimal_size {
    //    closest_filename = filename
    //    minimal_size = DirectorySizeDistance(subdir, size)
    //     }
    //   }
    // }
    // }
    // if closest_filename set
    //  return closest_filename
    // return none
    return Icon{}, false
}
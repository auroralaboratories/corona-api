package icons

import (
	"strings"
	"testing"
)

func TestGetIconThemePaths(t *testing.T) {
	have := GetIconThemePaths()

	if len(have) > 0 {
		for i, p := range have {
			t.Logf("Path %d: %s", i, p)
		}
	} else {
		t.Errorf("Did not discover any icon theme paths")
	}
}

func TestLoadThemes(t *testing.T) {
	themeset := NewThemeset()

	if err := themeset.Load(); err == nil {
		if len(themeset.Themes) > 0 {
			for i, p := range themeset.Themes {
				t.Logf("Theme %02d: %s", i, p.InternalName)
				t.Logf("          Name:     %s", p.Name)
				t.Logf("          Inherits: %s", strings.Join(p.Inherits, `, `))
				t.Logf("          Icons:    %d", len(p.Icons))
				t.Logf("\n")
			}
		} else {
			t.Errorf("Did not discover any icon themes")
		}
	} else {
		t.Errorf("Error loading icon themeset: %v", err)
	}
}

// func TestGetCurrentTheme(t *testing.T) {
// }

func TestFindIconFromTheme(t *testing.T) {
	themeset := NewThemeset()

	if err := themeset.Load(); err == nil {
		var baseTheme *Theme

		for _, theme := range themeset.Themes {
			if theme.InternalName == `hicolor` {
				baseTheme = theme
				break
			}
		}

		if baseTheme != nil {
			if icon, ok := baseTheme.FindIcon([]string{`dropbox`}, 40); ok {
				t.Logf("Got icon: %+v", icon)
			} else {
				t.Errorf("Could not find 16x16 blank icon")
			}

		} else {
			t.Errorf("Error loading hicolor icon theme")
		}
	} else {
		t.Errorf("Error loading icon themeset: %v", err)
	}
}

func TestFindIconFromAllThemes(t *testing.T) {
	themeset := NewThemeset()

	if err := themeset.Load(); err == nil {
		if icon, ok := themeset.FindIconViaTheme(`Faenza-Dark`, []string{`playonlinux`}, 41); ok {
			t.Logf("Got icon: %+v", icon)
		} else {
			t.Errorf("Could not find icon 'playonlinux'")
		}
	} else {
		t.Errorf("Error loading icon themeset: %v", err)
	}
}

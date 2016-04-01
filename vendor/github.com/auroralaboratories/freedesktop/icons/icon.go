package icons

import (
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/vaughan0/go-ini"
)

type Icon struct {
	io.Reader

	Context               IconContext
	Filename              string
	Name                  string
	Type                  string
	Theme                 *Theme
	DisplayName           string
	EmbeddedTextRectangle *Rectangle
	AttachPoints          []Point

	dataFilename string
	fileHandle   *os.File
}

// Allocate a new icon from a given filename and theme instance
//
func NewIcon(filename string, theme *Theme) *Icon {
	rv := &Icon{
		Filename:     filename,
		Name:         strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename)),
		Theme:        theme,
		AttachPoints: make([]Point, 0),
	}

	rv.DisplayName = rv.Name
	rv.Type = strings.ToLower(strings.TrimPrefix(filepath.Ext(rv.Filename), `.`))

	return rv
}

// Return whether the icon has an associated Icon Data file
//
func (self *Icon) HasDataFile() bool {
	return (self.dataFilename != ``)
}

// Reload the icon's data from disk (if present)
//
func (self *Icon) Refresh() error {
	if !self.Context.IsValid() {
		return fmt.Errorf("Cannot process icon '%s' in theme '%s' without a valid theme context", self.Name, self.Theme.Name)
	}

	dataFilename := strings.TrimSuffix(self.Filename, filepath.Ext(self.Filename)) + `.icon`

	if iconData, err := ini.LoadFile(dataFilename); err == nil {
		self.dataFilename = dataFilename

		if d, ok := iconData[`Icon Data`]; ok {
			if v, ok := d[`DisplayName`]; ok {
				self.DisplayName = v
			}

			if v, ok := d[`EmbeddedTextRectangle`]; ok {
				if rect, err := CreateRectangleFromString(v); err == nil {
					self.EmbeddedTextRectangle = rect
				}
			}

			if v, ok := d[`AttachPoints`]; ok {
				self.AttachPoints = CreatePointsFromString(v)
			}
		}
	}

	return nil
}

// Return whether the given name/size matches this icon or not, accounting
// for icon type (fixed, scalable, or threshold).
//
func (self *Icon) IsMatch(name string, size int) bool {
	//  non-matching name immediately fails
	if self.Name != name {
		return false
	}

	//  test size details
	switch self.Context.Type {
	//  fixed context: exact match only
	case IconContextFixed:
		return (self.Context.Size == size)

		//  scalable context: size between (min, max) inclusive
	case IconContextScalable:
		return (size >= self.Context.MinSize && size <= self.Context.MaxSize)

		//  threshold context: size between context Size +/- threshold
	case IconContextThreshold:
		return (size >= (self.Context.Size-self.Context.Threshold) && size <= (self.Context.Size+self.Context.Threshold))
	}

	//  default false
	return false
}

// Return the absolute difference between this icon and the given size
//
func (self *Icon) DistanceFromSize(size int) int {
	//  test size details
	switch self.Context.Type {
	//  fixed context: exact match only
	case IconContextFixed:
		return int(math.Abs(float64(self.Context.Size - size)))

		//  scalable context: size between min and max (inclusive)
	case IconContextScalable:
		if size < self.Context.MinSize {
			return (self.Context.MinSize - size)
		} else if size > self.Context.MaxSize {
			return (size - self.Context.MaxSize)
		}

		//  threshold context: size between context Size +/- threshold
	case IconContextThreshold:
		if size < (self.Context.Size - self.Context.Threshold) {
			return (self.Context.MinSize - size)
		} else if size > (self.Context.Size + self.Context.Threshold) {
			return (size - self.Context.MaxSize)
		}
	}

	//  fallback to zero distance
	return 0
}

func (self *Icon) Open() (*os.File, error) {
	if file, err := os.Open(self.Filename); err == nil {
		self.fileHandle = file
		return self.fileHandle, nil
	} else {
		return nil, err
	}
}

func (self *Icon) Read(buffer []byte) (int, error) {
	if self.fileHandle == nil {
		if _, err := self.Open(); err != nil {
			return 0, err
		}
	}

	return self.fileHandle.Read(buffer)
}

func (self *Icon) Close() error {
	if self.fileHandle == nil {
		return fmt.Errorf("Cannot close unopened file")
	}

	err := self.fileHandle.Close()
	self.fileHandle = nil
	return err
}

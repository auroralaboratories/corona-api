package icons

import (
	"fmt"
	"strings"

	"github.com/shutterstock/go-stockutil/stringutil"
)

type Point struct {
	X int
	Y int
}

func CreatePointsFromString(spec string) []Point {
	points := make([]Point, 0)

	pairs := strings.Split(spec, `|`)

	for _, pair := range pairs {
		p := strings.SplitN(pair, `,`, 2)

		if len(p) == 2 {
			if x, err := stringutil.ConvertToInteger(p[0]); err == nil {
				if y, err := stringutil.ConvertToInteger(p[1]); err == nil {
					points = append(points, Point{
						X: int(x),
						Y: int(y),
					})
				}
			}
		}
	}

	return points
}

type Rectangle struct {
	TopLeft     Point
	BottomRight Point
}

func CreateRectangleFromString(spec string) (*Rectangle, error) {
	coords := strings.SplitN(spec, `,`, 4)

	if len(coords) == 4 {
		rect := &Rectangle{}

		if x0, err := stringutil.ConvertToInteger(coords[0]); err == nil {
			rect.TopLeft.X = int(x0)
		} else {
			return nil, fmt.Errorf("Invalid x0 coordinate: %v", err)
		}

		if y0, err := stringutil.ConvertToInteger(coords[1]); err == nil {
			rect.TopLeft.Y = int(y0)
		} else {
			return nil, fmt.Errorf("Invalid y0 coordinate: %v", err)
		}

		if x1, err := stringutil.ConvertToInteger(coords[2]); err == nil {
			rect.BottomRight.X = int(x1)
		} else {
			return nil, fmt.Errorf("Invalid x1 coordinate: %v", err)
		}

		if y1, err := stringutil.ConvertToInteger(coords[3]); err == nil {
			rect.BottomRight.Y = int(y1)
		} else {
			return nil, fmt.Errorf("Invalid y1 coordinate: %v", err)
		}

		return rect, nil
	}

	return nil, fmt.Errorf("Expected 4 values, got %d", len(coords))
}

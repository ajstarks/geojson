package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type LGeoJSON struct {
	Type     string      `json:"type"`
	Features []LFeatures `json:"features"`
}

type LFeatures struct {
	Type     string    `json:"type"`
	Geometry LGeometry `json:"geometry"`
}

type LGeometry struct {
	Type        string      `json:"type"`
	Coordinates [][]float64 `json:"coordinates"`
}

// coordinates as polygon
type GeoJSON struct {
	Type     string     `json:"type"`
	Features []Features `json:"features"`
}

type Features struct {
	Type     string   `json:"type"`
	Geometry Geometry `json:"geometry"`
}

type Geometry struct {
	Type        string        `json:"type"`
	Coordinates [][][]float64 `json:"coordinates"`
}

// coordinates as multipolygon
type MGeoJSON struct {
	Type     string      `json:"type"`
	Features []MFeatures `json:"features"`
}

type MFeatures struct {
	Type     string    `json:"type"`
	Geometry MGeometry `json:"geometry"`
}

type MGeometry struct {
	Type        string          `json:"type"`
	Coordinates [][][][]float64 `json:"coordinates"`
}

// Parse parses GeoJSON geometry
func Parse(r io.Reader) (GeoJSON, error) {
	var data GeoJSON
	err := json.NewDecoder(r).Decode(&data)
	return data, err
}

// LineParse parses GeoJSON lines
func LineParse(r io.Reader) (LGeoJSON, error) {
	var data LGeoJSON
	err := json.NewDecoder(r).Decode(&data)
	return data, err
}

// MultiParse parses muliple GeoJSON structure
func MultiParse(r io.Reader) (MGeoJSON, error) {
	var data MGeoJSON
	err := json.NewDecoder(r).Decode(&data)
	return data, err
}

// LineCoordinates returns a set of line lat/long coordinates
func LineCoordinates(data LGeoJSON) ([]float64, []float64) {
	var xp, yp []float64
	for _, f := range data.Features {
		if f.Geometry.Type == "LineString" {
			coords := f.Geometry.Coordinates
			for _, c := range coords {
				xp = append(xp, c[1])
				yp = append(yp, c[0])
			}
		}
	}
	return xp, yp
}

// Coordinates returns polygon lat long coordinates
func Coordinates(data GeoJSON) ([]float64, []float64) {
	var xp, yp []float64
	for _, f := range data.Features {
		if f.Geometry.Type == "Polygon" {
			coords := f.Geometry.Coordinates
			for _, c := range coords[0] {
				xp = append(xp, c[1])
				yp = append(yp, c[0])
			}
		}
	}
	return xp, yp
}

// MultiCoordinates returns multipolygon lat/long coordinates
func MultiCoordinates(data MGeoJSON) ([]float64, []float64) {
	var xp, yp []float64
	for _, f := range data.Features {
		if f.Geometry.Type == "MultiPolygon" {
			coords := f.Geometry.Coordinates
			for _, c := range coords[0][0] {
				xp = append(xp, c[1])
				yp = append(yp, c[0])
			}
		}
	}
	return xp, yp
}

// corderr processes coordinate errors
func coorderr(s string, err error) {
	if len(s) > 0 {
		fmt.Fprintf(os.Stderr,
			"coordinates encoded as %q (use -type %s)\n", s, strings.ToLower(s[0:1]))
	}
	fmt.Fprintln(os.Stderr, err)
}

// coords writes lat/long coordinates to an io.Writer read from an io.Reader
func coords(w io.Writer, r io.Reader, ptype string) {
	var x, y []float64
	switch ptype {
	case "linestring", "l", "ls":
		data, err := LineParse(r)
		if err != nil {
			coorderr(data.Features[0].Geometry.Type, err)
			return
		}
		x, y = LineCoordinates(data)

	case "polygon", "p", "poly":
		data, err := Parse(r)
		if err != nil {
			coorderr(data.Features[0].Geometry.Type, err)
			return
		}
		x, y = Coordinates(data)

	case "multipolygon", "m", "mp":
		data, err := MultiParse(r)
		if err != nil {
			coorderr(data.Features[0].Geometry.Type, err)
			return
		}
		x, y = MultiCoordinates(data)
	}

	for i := 0; i < len(x); i++ {
		fmt.Fprintf(w, "%v %v\n", x[i], y[i])
	}
}

func main() {

	var ptype string
	flag.StringVar(&ptype, "type", "polygon", "type of coordinate ([l]linestring, [p]olygon, [m]ultipolygon)")
	flag.Parse()

	// if no args read/write from stdin/stdout
	if len(flag.Args()) == 0 {
		coords(os.Stdout, os.Stdin, ptype)
		return
	}

	// for every file read and write coordinates
	for _, filename := range flag.Args() {
		r, err := os.Open(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			continue
		}
		coords(os.Stdout, r, ptype)
		r.Close()
	}
}

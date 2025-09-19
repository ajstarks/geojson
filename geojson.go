package geojson

import (
	"encoding/json"
	"io"
)

// coordinates as linestring
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

func LineParse(r io.Reader) (LGeoJSON, error) {
	var data LGeoJSON
	err := json.NewDecoder(r).Decode(&data)
	return data, err
}

func Parse(r io.Reader) (GeoJSON, error) {
	var data GeoJSON
	err := json.NewDecoder(r).Decode(&data)
	return data, err
}

func MultiParse(r io.Reader) (MGeoJSON, error) {
	var data MGeoJSON
	err := json.NewDecoder(r).Decode(&data)
	return data, err
}

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

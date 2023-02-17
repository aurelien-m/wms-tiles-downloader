/*
Code below ports some functions from Python's https://github.com/mapbox/mercantile package.
*/

package mercantile

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"

	geoshapes "github.com/ctessum/geom"
	geojson "github.com/paulmach/go.geojson"
)

// Bbox represents Web Mercator Bounding Box.
type Bbox struct {
	Left   float64
	Bottom float64
	Right  float64
	Top    float64
}

// TileID represents ID of the tile in X, Y, Z format.
type TileID struct {
	X int
	Y int
	Z int
}

// LngLat represents point in space.
type LngLat struct {
	Lng float64
	Lat float64
}

// Xy retrieves Spherical Mercator (x, y) in meters.
func Xy(lngLat LngLat) (x, y float64) {
	lng := lngLat.Lng * (math.Pi / 180.0)
	lat := lngLat.Lat * (math.Pi / 180.0)
	x = 6378137.0 * lng
	y = 6378137.0 * math.Log(math.Tan((math.Pi*0.25)+(0.5*lat)))
	return x, y
}

// Ul retrieves upper left coordinates (lon, lat) of a tile.
func Ul(tile TileID) LngLat {
	n := math.Pow(2.0, float64(tile.Z))
	lonDeg := float64(tile.X)/n*360.0 - 180.0
	latRad := math.Atan(math.Sinh(math.Pi * (1 - 2*float64(tile.Y)/n)))
	latDeg := (180.0 / math.Pi) * latRad
	return LngLat{lonDeg, latDeg}
}

// Br retrieves bottom right coordinates (lon, lat) of a tile.
func Br(tile TileID) LngLat {
	nextTile := TileID{tile.X + 1, tile.Y + 1, tile.Z}
	return Ul(nextTile)
}

// XyBounds retrieves Spherical Mercator Bounding Box of a tile.
func XyBounds(tile TileID) Bbox {
	left, top := Xy(Ul(tile))
	nextTile := TileID{tile.X + 1, tile.Y + 1, tile.Z}
	right, bottom := Xy(Ul(nextTile))
	return Bbox{left, bottom, right, top}
}

// Tile retrieves tile containing given longitude and latitude.
func Tile(lng float64, lat float64, zoom int) TileID {
	lat = lat * (math.Pi / 180.0)
	n := math.Pow(2.0, float64(zoom))
	tileX := int(math.Floor((lng + 180.0) / 360.0 * n))
	tileY := int(math.Floor((1.0 - math.Log(math.Tan(lat)+(1.0/math.Cos(lat)))/math.Pi) / 2.0 * n))
	return TileID{tileX, tileY, zoom}
}

type PolygonOption func(tile TileID) float64

func WithPolygon(polygon geoshapes.Polygon) PolygonOption {
	return func(tile TileID) float64 {
		ul := Ul(tile)
		br := Br(tile)
		tile_latlon := geoshapes.Polygon{{
			geoshapes.Point{X: ul.Lng, Y: ul.Lat},
			geoshapes.Point{X: br.Lng, Y: ul.Lat},
			geoshapes.Point{X: br.Lng, Y: br.Lat},
			geoshapes.Point{X: ul.Lng, Y: br.Lat},
		}}

		inter := polygon.Intersection(tile_latlon)
		return inter.Area()
	}
}

func Tiles(bboxes [][]float64, zooms []int, options ...PolygonOption) []TileID {
	var tiles []TileID
	for _, bbox := range bboxes {
		w := math.Max(-180.0, bbox[0])
		s := math.Max(-85.051129, bbox[1])
		e := math.Min(180.0, bbox[2])
		n := math.Min(85.051129, bbox[3])

		for _, z := range zooms {
			ll := Tile(w, s, z)
			ur := Tile(e, n, z)

			var llx int
			var ury int

			if ll.X < 0 {
				llx = 0
			} else {
				llx = ll.X
			}
			if ur.Y < 0 {
				ury = 0
			} else {
				ury = ur.Y
			}

			for i := llx; i < int(math.Min(float64(ur.X)+1.0, math.Pow(2.0, float64(z)))); i++ {
				for j := ury; j < int(math.Min(float64(ll.Y)+1.0, math.Pow(2.0, float64(z)))); j++ {
					tile := TileID{i, j, z}

					if len(options) == 1 {
						area := options[0](tile)
						fmt.Printf("area: %f\n", area)
						if area == 0 {
							continue
						}
					}

					tiles = append(tiles, tile)
				}
			}
		}
	}

	return tiles
}

// TilesFromBbox retrieves tiles intersecting a geographic bounding box.
func TilesFromBbox(west, south, east, north float64, zooms []int) []TileID {
	var bboxes [][]float64
	if west > east {
		bboxWest := []float64{-180.0, south, east, north}
		bboxEast := []float64{west, south, 180.0, north}
		bboxes = [][]float64{bboxWest, bboxEast}
	} else {
		bboxes = [][]float64{{west, south, east, north}}
	}
	return Tiles(bboxes, zooms)
}

// TilesFromGeoJSON retrieves tiles intersecting a GeoJSON polygon.
func TilesFromGeoJSON(geojson_path string, zooms []int) []TileID {
	jsonFile, err := os.Open(geojson_path)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	fc, err := geojson.UnmarshalFeatureCollection(byteValue)
	if err != nil {
		fmt.Printf("error: %v", err)
	}

	west := 0.0
	south := 0.0
	east := 0.0
	north := 0.0

	points := []geoshapes.Point{}
	for _, point := range fc.Features[0].Geometry.MultiPolygon[0][0] {
		points = append(points, geoshapes.Point{X: point[0], Y: point[1]})

		if point[0] < west || west == 0.0 {
			west = point[0]
		}
		if point[0] > east || east == 0.0 {
			east = point[0]
		}
		if point[1] < south || south == 0.0 {
			south = point[1]
		}
		if point[1] > north || north == 0.0 {
			north = point[1]
		}
	}
	polygon := geoshapes.Polygon{points}

	var bboxes [][]float64
	if west > east {
		bboxWest := []float64{-180.0, south, east, north}
		bboxEast := []float64{west, south, 180.0, north}
		bboxes = [][]float64{bboxWest, bboxEast}
	} else {
		bboxes = [][]float64{{west, south, east, north}}
	}

	return Tiles(bboxes, zooms, WithPolygon(polygon))
}

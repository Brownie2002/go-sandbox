package main

import (
	"fmt"

	geojson "github.com/kpawlik/geojson"
)

func main() {
	fc := geojson.NewFeatureCollection([]*geojson.Feature{})

	// feature
	p := geojson.NewPoint(geojson.Coordinate{12, 3.123})
	f1 := geojson.NewFeature(p, nil, nil)
	fc.AddFeatures(f1)

	// feature with propertises
	props := map[string]interface{}{"name": "location", "code": 107}
	f2 := geojson.NewFeature(p, props, nil)
	fc.AddFeatures(f2)

	// feature with propertises and id
	f3 := geojson.NewFeature(p, props, 11101)
	fc.AddFeatures(f3)

	ls := geojson.NewLineString(geojson.Coordinates{{1, 1}, {2.001, 3}, {4001, 1223}})
	f4 := geojson.NewFeature(ls, nil, nil)
	fc.AddFeatures(f4)

	// Polygon feature

	var arrayCoordinates []geojson.Coordinates

	arrayCoordinates = append(arrayCoordinates, geojson.Coordinates{{2, 50}, {2, 51}, {3, 51}, {3, 50}, {2, 50}})

	poly := geojson.NewPolygon(arrayCoordinates)
	f5 := geojson.NewFeature(poly, nil, nil)

	fc.AddFeatures(f5)

	if gjstr, err := geojson.Marshal(fc); err != nil {
		panic(err)
	} else {
		fmt.Println(gjstr)
	}
}

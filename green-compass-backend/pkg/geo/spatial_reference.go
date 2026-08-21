package geo

import "fmt"

const SRID4326 = 4326

func (p Point) WKT() string {
	return fmt.Sprintf("POINT(%v %v)", p.Lon, p.Lat)
}

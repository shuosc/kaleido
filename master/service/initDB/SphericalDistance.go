package initDB

import "math"

const EarthRadius = 6378

func rad(degree float64) float64 {
	return degree * math.Pi / 180.0
}

func distance(lat1 float64, lng1 float64, lat2 float64, lng2 float64) float64 {
	radLat1 := rad(lat1)
	radLat2 := rad(lat2)
	a := math.Pow(math.Sin((radLat1-radLat2)/2), 2)
	b := math.Cos(radLat1) * math.Cos(radLat2) * math.Pow(math.Sin((rad(lng1)-rad(lng2))/2), 2)
	return EarthRadius * 2 * math.Asin(math.Sqrt(a+b))
}

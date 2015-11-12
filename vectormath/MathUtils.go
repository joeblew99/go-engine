package vectormath

import "math"

func Round(val float64, roundOn float64, places int) float64 {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	return round / pow
}

func RoundHalfUp(val float64) (newVal int) {
	return (int)(Round(val, .5, 0))
}

func ApproxEqual(value1, value2, epsilon float64) bool {
	return math.Abs(value1-value2) <= epsilon
}

//PointToPlaneDist distance from plane (a,b,c) to point
func PointToPlaneDist(a, b, c, point Vector3) float64 {
	ab := b.Subtract(a)
	ac := c.Subtract(a)
	ap := point.Subtract(a)
	normal := ab.Cross(ac).Normalize()
	return math.Abs(ap.Dot(normal))
}

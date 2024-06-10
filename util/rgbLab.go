package util

import (
	"math"

	"strconv"

	"strings"
)

// This is a very (not) nice function to convert RGB values to CIELAB Colors
// This assumes an observer angle of 2 degrees, and a white point of D65
// Frankly i'm not going to document the math since its fairly widely available online
func RgbToLab(rScale uint8, gScale uint8, bScale uint8) (float64, float64, float64) {

	r := float64(rScale) / 255.0

	g := float64(gScale) / 255.0

	b := float64(bScale) / 255.0

	var x, y, z float64

	r = Ternary((r > 0.04045), math.Pow((r+0.055)/1.055, 2.4), r/12.92).(float64)
	g = Ternary((g > 0.04045), math.Pow((g+0.055)/1.055, 2.4), g/12.92).(float64)
	b = Ternary((b > 0.04045), math.Pow((b+0.055)/1.055, 2.4), b/12.92).(float64)

	x = (r*0.4124 + g*0.3576 + b*0.1805) / 0.95047
	y = (r*0.2126 + g*0.7152 + b*0.0722) / 1.00000
	z = (r*0.0193 + g*0.1192 + b*0.9505) / 1.08883

	x = Ternary((x > 0.008856), math.Pow(x, 1.0/3.0), (7.787*x)+16.0/116.0).(float64)
	y = Ternary((y > 0.008856), math.Pow(y, 1.0/3.0), (7.787*y)+16.0/116.0).(float64)
	z = Ternary((z > 0.008856), math.Pow(z, 1.0/3.0), (7.787*z)+16.0/116.0).(float64)

	return (116 * y) - 16, 500 * (x - y), 200 * (y - z)

}

// Cauculate the DeltaE distance metric between two lab colors under the CIE74 standard
func DeltaE(l1 float64, a1 float64, b1 float64, l2 float64, a2 float64, b2 float64) float64 {

	//Literally just euclidean distance
	return math.Sqrt(math.Pow(l2-l1, 2.0) + math.Pow(a2-a1, 2.0) + math.Pow(b2-b1, 2.0))

}

// Convert HEX to RGB
func HexToRGB(hex string) (uint8, uint8, uint8) {

	//Remove the hashtag if present and convert into a uint64
	values, _ := strconv.ParseUint(strings.Replace(string(hex), "#", "", 1), 16, 32)

	//Do the standard bitshifting
	return uint8(values >> 16), uint8((values >> 8) & 0xFF), uint8(values & 0xFF)

}

// Convert Hex to LAB, just calls the HexToRGB and then RgbToLab
func HexToLAB(hex string) (float64, float64, float64) {
	r, g, b := HexToRGB(hex)

	return RgbToLab(r, g, b)
}

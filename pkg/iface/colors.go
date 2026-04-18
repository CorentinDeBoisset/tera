package iface

import (
	"fmt"
	"image/color"
	"math"
)

type ColorHsl struct {
	H, S, L float64
}

func (c ColorHsl) Clamped() ColorHsl {
	return ColorHsl{
		H: min(max(c.H, 0), 360),
		S: min(max(c.S, 0), 1),
		L: min(max(c.L, 0), 1),
	}
}

func (c ColorHsl) Hex() string {
	return hslToHex(c)
}

// The color calulation methods are inspired from: https://gist.github.com/ciembor/1494530
// Kudos to them :)

func colorToHsl(bgColor color.Color) ColorHsl {
	r, g, b, _ := bgColor.RGBA()

	// Normalize r, g and b to [0, 1]
	normR := max(min(float64(r)/65535, 1), 0)
	normG := max(min(float64(g)/65535, 1), 0)
	normB := max(min(float64(b)/65535, 1), 0)

	cMax := max(normR, normG, normB)
	cMin := min(normR, normG, normB)

	chroma := cMax - cMin
	light := (cMax + cMin) / 2

	if chroma == 0 {
		return ColorHsl{0, 0, light}
	}

	var hue, sat float64

	// Calculate the saturation
	// The denominator is never 0, because if so chroma=0 and this case is handled above
	sat = chroma / (1 - math.Abs(2*light-1))

	// Calculate the hue
	switch cMax {
	case normR:
		// The double module is to ensure we have a positive result
		hue = math.Mod(math.Mod((normG-normB)/chroma, 6)+6, 6)
	case normG:
		hue = (normB-normR)/chroma + 2
	default:
		hue = (normR-normG)/chroma + 4
	}
	hue *= 60

	return ColorHsl{hue, sat, light}
}

// Convert a hue to r, g, or b
// The float is in [0, 1]
func hueToRgb(p, q, t float64) float64 {
	if t < 0 {
		t += 1
	}
	if t > 1 {
		t -= 1
	}
	if t < 1./6 {
		return p + (q-p)*6*t
	}
	if t < 1./2 {
		return q
	}
	if t < 2./3 {
		return p + (q-p)*(2./3-t)*6
	}

	return p
}

// Convert an HSL color to its hex RGB representation ("#A1C3EE" for instance)
// It assumes c.H is in [0, 360], and c.S and c.L are within [0, 1]
func hslToHex(c ColorHsl) string {
	var r, g, b uint8
	if c.S == 0 {
		// Add 0.5 to avoid rounding errors
		r, g, b = uint8(c.L*255.+0.5), uint8(c.L*255.+0.5), uint8(c.L*255.+0.5)
	} else {
		var q float64
		if c.L < 0.5 {
			q = c.L * (1 + c.S)
		} else {
			q = c.L + c.S - c.L*c.S
		}
		p := 2*c.L - q

		r = uint8(hueToRgb(p, q, c.H/360+1./3)*255. + 0.5)
		g = uint8(hueToRgb(p, q, c.H/360)*255. + 0.5)
		b = uint8(hueToRgb(p, q, c.H/360-1./3)*255. + 0.5)
	}

	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
}

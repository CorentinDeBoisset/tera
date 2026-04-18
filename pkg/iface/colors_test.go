package iface

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestColorToHsl(t *testing.T) {
	t.Parallel()

	// black
	assert.Equal(t, ColorHsl{0, 0, 0}, colorToHsl(color.RGBA{0, 0, 0, 255}))

	// middle grey
	assert.Equal(t, ColorHsl{0, 0, 0.5019607843137255}, colorToHsl(color.RGBA{128, 128, 128, 255}))

	// purple
	assert.Equal(t, ColorHsl{300, 0.4980392156862745, 0.5}, colorToHsl(color.RGBA{191, 64, 191, 255}))

	// transparent and opaque green - the opacity should have no effect
	assert.Equal(t, ColorHsl{61.875, 0.6399999999999999, 0.3921568627450981}, colorToHsl(color.RGBA{160, 164, 36, 0}))
	assert.Equal(t, ColorHsl{61.875, 0.6399999999999999, 0.3921568627450981}, colorToHsl(color.RGBA{160, 164, 36, 255}))

	// light-blue
	assert.Equal(t, ColorHsl{195.0, 0.8045977011494253, 0.6588235294117647}, colorToHsl(color.RGBA{98, 203, 238, 255}))
}

func TestHslToHex(t *testing.T) {
	t.Parallel()

	// black
	assert.Equal(t, "#000000", hslToHex(ColorHsl{0, 0, 0}))

	// middle grey
	assert.Equal(t, "#808080", hslToHex(ColorHsl{0, 0, 0.5}))

	// purple
	assert.Equal(t, "#bf40bf", hslToHex(ColorHsl{300, 0.5, 0.5}))

	// green
	assert.Equal(t, "#a0a424", hslToHex(ColorHsl{61.87, 0.64, 0.392}))

	// light-blue
	assert.Equal(t, "#62cbee", hslToHex(ColorHsl{195, 0.805, 0.659}))

	// pink
	assert.Equal(t, "#ba20d5", hslToHex(ColorHsl{291, 0.74, 0.48}))
}

func TestHexToColor(t *testing.T) {
	t.Parallel()

	assert.Equal(t, color.RGBA{0, 0, 0, 255}, hexToColor("#000000"))

	assert.Equal(t, color.RGBA{128, 128, 128, 255}, hexToColor("#808080"))
	assert.Equal(t, color.RGBA{128, 128, 128, 255}, hexToColor("808080"))

	assert.Equal(t, color.RGBA{160, 164, 36, 255}, hexToColor("#a0a424"))

	assert.Equal(t, color.RGBA{98, 203, 238, 255}, hexToColor("#62cbee"))

	assert.Equal(t, color.RGBA{186, 32, 213, 255}, hexToColor("#ba20d5"))

	// Test short hex colors
	assert.Equal(t, color.RGBA{34, 34, 34, 255}, hexToColor("#222"))
	assert.Equal(t, color.RGBA{255, 170, 68, 255}, hexToColor("#FA4"))
	assert.Equal(t, color.RGBA{255, 170, 68, 255}, hexToColor("FA4"))

	// Test error inputs
	assert.Equal(t, color.RGBA{0, 0, 0, 0}, hexToColor(""))                // too short string
	assert.Equal(t, color.RGBA{0, 0, 0, 0}, hexToColor("very long input")) // too long string
	assert.Equal(t, color.RGBA{0, 0, 0, 0}, hexToColor("somthng"))         // 7 letter input without '#'
	assert.Equal(t, color.RGBA{0, 0, 0, 0}, hexToColor("somt"))            // 4 letter input without '#'
	assert.Equal(t, color.RGBA{0, 0, 0, 0}, hexToColor("supere"))          // 6 letter input but non-parseable
	assert.Equal(t, color.RGBA{0, 0, 0, 0}, hexToColor("sup"))             // 3 letter input but non-parseable
}

func TestColorHsl(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "#a0a424", ColorHsl{61.87, 0.64, 0.392}.Hex())

	assert.Equal(t, ColorHsl{360., 1., 0.}, ColorHsl{362, 1.2, -0.1}.Clamped())
}

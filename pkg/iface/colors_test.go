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

func TestColorHsl(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "#a0a424", ColorHsl{61.87, 0.64, 0.392}.Hex())

	assert.Equal(t, ColorHsl{360., 1., 0.}, ColorHsl{362, 1.2, -0.1}.Clamped())
}

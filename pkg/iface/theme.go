package iface

import (
	"image/color"
	"log"
	"os"
	"sync"

	"charm.land/lipgloss/v2"
)

// This is the actual color of the background of the terminal
var BackgroundColor = sync.OnceValue(func() color.Color {
	rawColor, err := lipgloss.BackgroundColor(os.Stdin, os.Stderr)
	if err != nil || rawColor == nil {
		log.Printf("the request of the background color failed: %v", err)

		return color.RGBA{0, 0, 0, 255}
	}

	return rawColor
})

// Generic theme used in the TUI applications
type Theme struct {
	NoticeableSurface         lipgloss.Style
	FocusedHighlightSurface   lipgloss.Style
	UnfocusedHighlightSurface lipgloss.Style

	FocusedHighlightText   lipgloss.Style
	UnfocusedHighlightText lipgloss.Style

	Separator lipgloss.Style

	FocusedOutputBorder   lipgloss.Style
	UnfocusedOutputBorder lipgloss.Style

	SuccessTextColor color.Color
	WarningTextColor color.Color
	ErrorTextColor   color.Color
}

// Theme used to display the usage and help insructions
type HelpTheme struct {
	BaseTitle  lipgloss.Style
	ErrorTitle lipgloss.Style

	Command    lipgloss.Style
	SubCommand lipgloss.Style
	DimmedArg  lipgloss.Style
	Flag       lipgloss.Style

	Codeblock           lipgloss.Style
	CodeblockCommand    lipgloss.Style
	CodeblockSubCommand lipgloss.Style
	CodeblockDimmedArg  lipgloss.Style
}

func LoadHelpTheme(background color.Color) HelpTheme {
	bgColor := colorToHsl(background)
	lightDark := lipgloss.LightDark(bgColor.L < 0.5)

	var codeblockSurface, codeblockForeground ColorHsl
	var baseTitleColor, commandColor, subCommandColor, dimmedArgColor, flagColor ColorHsl
	var codeblockCommandColor, codeblockSubCommandColor, codeblockDimmedArgColor ColorHsl

	body_hue := 43.
	green_hue := 82.5
	red_hue := 8.
	yellow_hue := 53.

	if bgColor.L <= 0.2 {
		// Lighter surfaces, light text
		lScale := bgColor.L / 0.2

		baseTitleColor = ColorHsl{green_hue, 0.70, 0.35 + 0.15*lScale}.Clamped()
		commandColor = ColorHsl{red_hue, 0.80, 0.6 + 0.1*lScale}.Clamped()
		subCommandColor = ColorHsl{yellow_hue, 0.85, 0.6 + 0.1*lScale}.Clamped()
		dimmedArgColor = ColorHsl{bgColor.H, bgColor.S, 0.4 + 0.1*lScale}.Clamped()
		flagColor = ColorHsl{yellow_hue, 0.7, 0.45 + 0.15*lScale}.Clamped()
		codeblockSurface = ColorHsl{bgColor.H, bgColor.S, 0.15 + 0.08*lScale}.Clamped()
		codeblockForeground = ColorHsl{body_hue, 0.58, 0.4 + 0.4*lScale}.Clamped()
		codeblockCommandColor = ColorHsl{red_hue, 0.80, 0.55 + 0.1*lScale}.Clamped()
		codeblockSubCommandColor = ColorHsl{green_hue, 0.70, 0.4 + 0.3*lScale}.Clamped()
		codeblockDimmedArgColor = ColorHsl{bgColor.H, bgColor.S, 0.4 + 0.2*lScale}.Clamped()
	} else if bgColor.L < 0.5 {
		// Darker surfaces, light text
		lScale := (bgColor.L - 0.2) / 0.3

		baseTitleColor = ColorHsl{green_hue, 0.80, 0.6 + 0.15*lScale}.Clamped()
		commandColor = ColorHsl{red_hue, 0.70, 0.6 + 0.25*lScale}.Clamped()
		subCommandColor = ColorHsl{yellow_hue, 0.85, 0.7 + 0.05*lScale}.Clamped()
		dimmedArgColor = ColorHsl{bgColor.H, bgColor.S, 0.5 + 0.3*lScale}.Clamped()
		flagColor = ColorHsl{yellow_hue, 0.7, 0.6 + 0.25*lScale}.Clamped()
		codeblockSurface = ColorHsl{bgColor.H, bgColor.S, 0.17 + 0.2*lScale}.Clamped()
		codeblockForeground = ColorHsl{body_hue, 0.58, 0.8 + 0.15*lScale}.Clamped()
		codeblockCommandColor = ColorHsl{red_hue, 0.80, 0.6 + 0.2*lScale}.Clamped()
		codeblockSubCommandColor = ColorHsl{green_hue, 0.70, 0.6 + 0.22*lScale}.Clamped()
		codeblockDimmedArgColor = ColorHsl{bgColor.H, bgColor.S, 0.6 + 0.2*lScale}.Clamped()
	} else if bgColor.L < 0.75 {
		// Lighter surfaces, dark text
		lScale := (bgColor.L - 0.5) / 0.25

		baseTitleColor = ColorHsl{green_hue, 0.80, 0.2 + 0.05*lScale}.Clamped()
		commandColor = ColorHsl{red_hue, 1.0, 0.35 + 0.05*lScale}.Clamped()
		subCommandColor = ColorHsl{yellow_hue, 0.85, 0.15 + 0.05*lScale}.Clamped()
		dimmedArgColor = ColorHsl{bgColor.H, bgColor.S, 0.23 + 0.05*lScale}.Clamped()
		flagColor = ColorHsl{yellow_hue, 0.7, 0.15 + 0.1*lScale}.Clamped()
		codeblockSurface = ColorHsl{bgColor.H, bgColor.S, 0.6 + 0.2*lScale}.Clamped()
		codeblockForeground = ColorHsl{body_hue, 0.58, 0.07 + 0.08*lScale}.Clamped()
		codeblockCommandColor = ColorHsl{red_hue, 1.0, 0.35 + 0.04*lScale}.Clamped()
		codeblockSubCommandColor = ColorHsl{green_hue, 0.8, 0.17 + 0.12*lScale}.Clamped()
		codeblockDimmedArgColor = ColorHsl{bgColor.H, bgColor.S, 0.23 + 0.1*lScale}.Clamped()
	} else {
		// Darker surfaces, dark text
		lScale := (bgColor.L - 0.75) / 0.25

		baseTitleColor = ColorHsl{green_hue, 0.80, 0.25 + 0.05*lScale}.Clamped()
		commandColor = ColorHsl{red_hue, 1.0, 0.4 + 0.05*lScale}.Clamped()
		subCommandColor = ColorHsl{yellow_hue, 0.85, 0.2 + 0.05*lScale}.Clamped()
		dimmedArgColor = ColorHsl{bgColor.H, bgColor.S, 0.25 + 0.1*lScale}.Clamped()
		flagColor = ColorHsl{yellow_hue, 0.7, 0.25 + 0.1*lScale}.Clamped()
		codeblockSurface = ColorHsl{bgColor.H, 0.05 + 0.8*bgColor.S, 0.65 + 0.3*lScale}.Clamped()
		codeblockForeground = ColorHsl{body_hue, 0.58, 0.07 + 0.08*lScale}.Clamped()
		codeblockCommandColor = ColorHsl{red_hue, 1.0, 0.4 + 0.04*lScale}.Clamped()
		codeblockSubCommandColor = ColorHsl{green_hue, 0.8, 0.2 + 0.07*lScale}.Clamped()
		codeblockDimmedArgColor = ColorHsl{bgColor.H, bgColor.S, 0.35 + 0.15*lScale}.Clamped()
	}

	baseCodeBlock := lipgloss.NewStyle().Background(lipgloss.Color(hslToHex(codeblockSurface)))

	return HelpTheme{
		BaseTitle:  lipgloss.NewStyle().Foreground(lipgloss.Color(baseTitleColor.Hex())),
		ErrorTitle: lipgloss.NewStyle().Background(lipgloss.Red).Foreground(lightDark(lipgloss.BrightWhite, lipgloss.Black)),

		Command:    lipgloss.NewStyle().Foreground(lipgloss.Color(commandColor.Hex())),
		SubCommand: lipgloss.NewStyle().Foreground(lipgloss.Color(subCommandColor.Hex())),
		DimmedArg:  lipgloss.NewStyle().Foreground(lipgloss.Color(dimmedArgColor.Hex())),
		Flag:       lipgloss.NewStyle().Foreground(lipgloss.Color(flagColor.Hex())),

		Codeblock:           baseCodeBlock.Padding(1, 2).Margin(0, 1).Foreground(lipgloss.Color(codeblockForeground.Hex())),
		CodeblockCommand:    baseCodeBlock.Foreground(lipgloss.Color(codeblockCommandColor.Hex())),
		CodeblockSubCommand: baseCodeBlock.Foreground(lipgloss.Color(codeblockSubCommandColor.Hex())),
		CodeblockDimmedArg:  baseCodeBlock.Foreground(lipgloss.Color(codeblockDimmedArgColor.Hex())),
	}
}

func LoadTheme(background color.Color) Theme {
	bgColor := colorToHsl(background)

	var noticeableSurfaceColor, unfocusedHighlightSurfaceColor, highlightSurfaceColor ColorHsl
	var invertedUnfocusedHighlightSurfaceColor, invertedHighlightSurfaceColor ColorHsl
	var bodyColorOnNoticeable, bodyColorOnUnfocusedHighlight, bodyColorOnHighlight ColorHsl
	var bodyColorOnInvertedUnfocusedHighlight, bodyColorOnInvertedHighlight ColorHsl
	var separatorCol, focusedOutputBorderColor ColorHsl
	var successTextCol, warningTextCol, errorTextCol ColorHsl

	/*
		The color theme is heavily dependent on the lightness of the background.

		The calculations are not the same for all ranges of lightness, to always have a coherent color palette.
		Thoses ranges are:
			* L in [0 ; 0.2]    -> the surfaces will be lighter, and the text will be light
			* L in ]0.2 ; 0.5]  -> the surfaces will be darker, and the text will be light
			* L in ]0.5 ; 0.75] -> the surfaces will be lighter, and the text will be dark
			* L in ]0.75 ; 1]   -> the surfaces will be darker, and the text will be light
	*/

	brickHue := 92.
	if bgColor.H >= 65 && bgColor.H < 190 {
		// For all green background, we use blue bricks rather than green-on-green
		brickHue = 211
	}

	borderHue := 26.
	if bgColor.H >= 5 && bgColor.H < 52 {
		// For orange backgrounds, the border hue is a blue instead
		borderHue = 198
	}

	body_hue := 43.

	if bgColor.L <= 0.2 {
		// Lighter surfaces, light text
		lScale := bgColor.L / 0.2

		noticeableSurfaceColor = ColorHsl{bgColor.H, bgColor.S, 0.15 + 0.1*lScale}.Clamped()
		bodyColorOnNoticeable = ColorHsl{body_hue, 0.58, 0.8 + 0.15*lScale}.Clamped()

		unfocusedHighlightSurfaceColor = ColorHsl{brickHue, 0.35, 0.10 + 0.08*lScale}.Clamped()
		bodyColorOnUnfocusedHighlight = ColorHsl{body_hue, 0.58, 0.8 + 0.15*lScale}.Clamped()

		highlightSurfaceColor = ColorHsl{brickHue, 0.42, 0.18 + 0.04*lScale}.Clamped()
		bodyColorOnHighlight = ColorHsl{body_hue, 1.0, 0.95}

		invertedUnfocusedHighlightSurfaceColor = ColorHsl{26.7, 0.65, 0.55}
		bodyColorOnInvertedUnfocusedHighlight = ColorHsl{0, 0, 0.07}

		invertedHighlightSurfaceColor = ColorHsl{26.7, 0.95, 0.95}
		bodyColorOnInvertedHighlight = ColorHsl{0, 0, 0.1}

		separatorCol = ColorHsl{0, 0, 0.4}
		focusedOutputBorderColor = ColorHsl{borderHue, 0.87, 0.55}

		successTextCol = ColorHsl{66.7, 61, 0.5 + 0.1*lScale}.Clamped()
		warningTextCol = ColorHsl{47, 87, 0.5 + 0.1*lScale}.Clamped()
		errorTextCol = ColorHsl{11, 92, 0.65 + 0.06*lScale}.Clamped()
	} else if bgColor.L < 0.5 {
		// Darker surfaces, light text
		lScale := (bgColor.L - 0.2) / 0.3

		noticeableSurfaceColor = ColorHsl{bgColor.H, bgColor.S, 0.17 + 0.2*lScale}.Clamped()
		bodyColorOnNoticeable = ColorHsl{body_hue, 0.58, 0.8 + 0.15*lScale}.Clamped()

		unfocusedHighlightSurfaceColor = ColorHsl{brickHue, 0.28, 0.14 + 0.04*lScale}.Clamped()
		bodyColorOnUnfocusedHighlight = ColorHsl{body_hue, 0.58, 0.8 + 0.15*lScale}.Clamped()

		highlightSurfaceColor = ColorHsl{brickHue, 0.42, 0.18 + 0.04*lScale}.Clamped()
		bodyColorOnHighlight = ColorHsl{body_hue, 1.0, 0.95}

		invertedUnfocusedHighlightSurfaceColor = ColorHsl{26.7, 0.65, 0.55}
		bodyColorOnInvertedUnfocusedHighlight = ColorHsl{0, 0, 0.07}

		invertedHighlightSurfaceColor = ColorHsl{26.7, 0.95, 0.95}
		bodyColorOnInvertedHighlight = ColorHsl{0, 0, 0.1}

		separatorCol = ColorHsl{0, 0, 0.4}
		focusedOutputBorderColor = ColorHsl{borderHue, 0.87, 0.72}

		successTextCol = ColorHsl{66.7, 73, 0.6 + 0.1*lScale}.Clamped()
		warningTextCol = ColorHsl{47, 87, 0.6 + 0.1*lScale}.Clamped()
		errorTextCol = ColorHsl{11, 92, 0.71 + 0.06*lScale}.Clamped()
	} else if bgColor.L < 0.75 {
		// Lighter surfaces, dark text
		lScale := (bgColor.L - 0.5) / 0.25

		noticeableSurfaceColor = ColorHsl{bgColor.H, bgColor.S, 0.6 + 0.2*lScale}.Clamped()
		bodyColorOnNoticeable = ColorHsl{body_hue, 0.58, 0.07 + 0.08*lScale}.Clamped()

		unfocusedHighlightSurfaceColor = ColorHsl{brickHue, 0.25, 0.5 + 0.25*bgColor.L}.Clamped()
		bodyColorOnUnfocusedHighlight = ColorHsl{body_hue, 0.58, 0.07 + 0.08*lScale}.Clamped()

		highlightSurfaceColor = ColorHsl{brickHue, 0.48, 0.55 + 0.25*lScale}.Clamped()
		bodyColorOnHighlight = ColorHsl{body_hue, 0.58, 0.1}

		invertedUnfocusedHighlightSurfaceColor = ColorHsl{26, 0.87, 0.6 + 0.1*lScale}.Clamped()
		bodyColorOnInvertedUnfocusedHighlight = ColorHsl{0, 0, 0.1}

		invertedHighlightSurfaceColor = ColorHsl{26, 0.9, 0.17 + 0.05*lScale}.Clamped()
		bodyColorOnInvertedHighlight = ColorHsl{0, 0, 0.95}

		separatorCol = ColorHsl{0, 0, 0.15}
		focusedOutputBorderColor = ColorHsl{borderHue, 0.87, 0.30}

		successTextCol = ColorHsl{66.7, 85, 0.1 + 0.02*lScale}.Clamped()
		warningTextCol = ColorHsl{47, 87, 0.15 + 0.1*lScale}.Clamped()
		errorTextCol = ColorHsl{11, 92, 0.27 + 0.05*lScale}.Clamped()
	} else {
		// Darker surfaces, dark text
		lScale := (bgColor.L - 0.75) / 0.25

		noticeableSurfaceColor = ColorHsl{bgColor.H, 0.05 + 0.8*bgColor.S, 0.65 + 0.3*lScale}.Clamped()
		bodyColorOnNoticeable = ColorHsl{body_hue, 0.58, 0.07 + 0.08*lScale}.Clamped()

		unfocusedHighlightSurfaceColor = ColorHsl{brickHue, 0.2, 0.75 + 0.1*lScale}.Clamped()
		bodyColorOnUnfocusedHighlight = ColorHsl{body_hue, 0.58, 0.07 + 0.08*lScale}.Clamped()

		highlightSurfaceColor = ColorHsl{brickHue, 0.27, 0.65 + 0.1*lScale}.Clamped()
		bodyColorOnHighlight = ColorHsl{body_hue, 0.58, 0.1}

		invertedUnfocusedHighlightSurfaceColor = ColorHsl{26, 0.87, 0.6 + 0.1*lScale}.Clamped()
		bodyColorOnInvertedUnfocusedHighlight = ColorHsl{0, 0, 0.1}

		invertedHighlightSurfaceColor = ColorHsl{26, 0.9, 0.17 + 0.05*lScale}.Clamped()
		bodyColorOnInvertedHighlight = ColorHsl{0, 0, 0.95}

		separatorCol = ColorHsl{0, 0, 0.15}
		focusedOutputBorderColor = ColorHsl{borderHue, 0.87, 0.45}

		successTextCol = ColorHsl{66.7, 85, 0.12 + 0.05*lScale}.Clamped()
		warningTextCol = ColorHsl{47, 87, 0.25 + 0.05*lScale}.Clamped()
		errorTextCol = ColorHsl{11, 92, 0.32 + 0.05*lScale}.Clamped()
	}

	// Convert them all to lipgloss colors
	lgNoticeableColor := lipgloss.Color(noticeableSurfaceColor.Hex())
	lgBodyColorOnNoticeable := lipgloss.Color(bodyColorOnNoticeable.Hex())

	lgUnfocusedHighlightSurfaceColor := lipgloss.Color(unfocusedHighlightSurfaceColor.Hex())
	lgBodyColorOnUnfocusedHighlight := lipgloss.Color(bodyColorOnUnfocusedHighlight.Hex())

	lgHighlightSurfaceColor := lipgloss.Color(highlightSurfaceColor.Hex())
	lgBodyColorOnHighlight := lipgloss.Color(bodyColorOnHighlight.Hex())

	lgInvertedHighlightSurfaceColor := lipgloss.Color(invertedHighlightSurfaceColor.Hex())
	lgBodyColorOnInvertedHighlight := lipgloss.Color(bodyColorOnInvertedHighlight.Hex())

	lgInvertedUnfocusedHighlightSurfaceColor := lipgloss.Color(invertedUnfocusedHighlightSurfaceColor.Hex())
	lgBodyColorOnInvertedUnfocusedHighlight := lipgloss.Color(bodyColorOnInvertedUnfocusedHighlight.Hex())

	return Theme{
		NoticeableSurface: lipgloss.NewStyle().
			Background(lgNoticeableColor).
			BorderBackground(lgNoticeableColor).
			Foreground(lgBodyColorOnNoticeable),

		FocusedHighlightSurface: lipgloss.NewStyle().
			Background(lgHighlightSurfaceColor).
			BorderBackground(lgHighlightSurfaceColor).
			Foreground(lgBodyColorOnHighlight),

		UnfocusedHighlightSurface: lipgloss.NewStyle().
			Background(lgUnfocusedHighlightSurfaceColor).
			BorderBackground(lgUnfocusedHighlightSurfaceColor).
			Foreground(lgBodyColorOnUnfocusedHighlight),

		FocusedHighlightText: lipgloss.NewStyle().
			Background(lgInvertedHighlightSurfaceColor).
			BorderBackground(lgInvertedHighlightSurfaceColor).
			Foreground(lgBodyColorOnInvertedHighlight),

		UnfocusedHighlightText: lipgloss.NewStyle().
			Background(lgInvertedUnfocusedHighlightSurfaceColor).
			BorderBackground(lgInvertedUnfocusedHighlightSurfaceColor).
			Foreground(lgBodyColorOnInvertedUnfocusedHighlight),

		Separator:             lipgloss.NewStyle().BorderForeground(lipgloss.Color(separatorCol.Hex())),
		FocusedOutputBorder:   lipgloss.NewStyle().BorderForeground(lipgloss.Color(focusedOutputBorderColor.Hex())),
		UnfocusedOutputBorder: lipgloss.NewStyle().BorderForeground(lipgloss.Color("#808080")),

		SuccessTextColor: lipgloss.Color(successTextCol.Hex()),
		WarningTextColor: lipgloss.Color(warningTextCol.Hex()),
		ErrorTextColor:   lipgloss.Color(errorTextCol.Hex()),
	}
}

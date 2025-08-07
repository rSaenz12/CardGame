package Ui

import "image/color"

var (
	materialRed  = color.NRGBA{R: 0xF4, G: 0x43, B: 0x36, A: 0xFF} // Hex: #F44336 (Similar to your existing red)
	materialBlue = color.NRGBA{R: 0x21, G: 0x96, B: 0xF3, A: 0xFF} // Hex: #2196F3 (Similar to your existing blue)
	black        = color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	white        = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
)

var (
	backgroundImagePath = "goFishGame/tableTop.png"
	faceDownCardPath    = "1B"
)

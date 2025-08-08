package uiFunctions

import (
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"image"
	"image/png"
	"log"
	"os"
)

// LoadImage loads the images using the path
func LoadImage(path string) image.Image {
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("ERROR: Failed to open image %s: %v\n", path, err)
	}
	img, err := png.Decode(f)
	if err != nil {
		log.Fatalf("ERROR: Failed to decode image %s: %v\n", path, err)
	}
	// Check if closing the file results in an error
	if err := f.Close(); err != nil {
		log.Fatalf("ERROR: Failed to close image file %s: %v\n", path, err)
	}
	return img
}

// GetCardImage grabs the path,loads images, adds to an array of images
func GetCardImage(currentHand []string) []image.Image {
	var images []image.Image

	//loops through the current hand calling each card, adding them as slices of images to the array
	for _, card := range currentHand {
		path := "deckImages/" + card + ".png"

		f, err := os.Open(path)
		if err != nil {
			log.Fatalf("ERROR: Failed to open image %s: %v\n", path, err)
		}

		img, err := png.Decode(f)
		if err != nil {
			log.Fatalf("ERROR: Failed to decode image %s: %v\n", path, err)
		}

		// Check if closing the file results in an error
		if err := f.Close(); err != nil {
			log.Fatalf("ERROR: Failed to close image file %s: %v\n", path, err)
		}

		images = append(images, img)
	}
	return images
}

func PrintCards(cards []image.Image, yOffset int, scale float32, cardHeight int) layout.FlexChild {
	return layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		// apply a vertical offset for the row of cards
		xOffset := (gtx.Constraints.Min.X / 2) - (15 * len(cards))
		offset := op.Offset(image.Pt(xOffset, yOffset)).Push(gtx.Ops)
		defer offset.Pop()

		// Loop through cards and draw each with spacing
		cardWidth := 25 // approx width after scaling
		spacing := 0
		startingOffset := int(float32(gtx.Constraints.Max.X) / 2)

		for i, card := range cards {
			cardOffset := op.Offset(image.Pt(startingOffset+(i*(cardWidth+spacing)), cardHeight)).Push(gtx.Ops)
			scaleOp := op.Affine(f32.Affine2D{}.Scale(f32.Pt(0, 0), f32.Pt(scale, scale))).Push(gtx.Ops)

			//prints card
			paint.NewImageOp(card).Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)

			scaleOp.Pop()
			cardOffset.Pop()
		}
		return layout.Dimensions{Size: image.Pt(len(cards)*(cardWidth+spacing), cardWidth)}
	})
}

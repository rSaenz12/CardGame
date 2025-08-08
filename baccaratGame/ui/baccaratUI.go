//UI and buttons for actual card bet playing. 
package ui

import (
	"CombinedCardgames/baccaratGame/game"
	ui2 "CombinedCardgames/uiFunctions"
	//"fmt"
	"image/color"
	"strconv"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

//This functio returns errors trigger by button clicks to play the game. 
func RunBaccaratUI(window *app.Window, g *game.Game) error {
	var ops op.Ops
	th := material.NewTheme()

	var playerBetBtn, bankerBetBtn, tieBetBtn, dealBtn, menuBtn, playAgainBtn widget.Clickable

	background := ui2.LoadImage(backgroundPath)
	bgOp := paint.NewImageOp(background)
	var imageWidget = widget.Image{Src: bgOp, Fit: widget.Cover}

	g.NewRound()

	for {
		e := window.Event()
		switch e := e.(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			if playerBetBtn.Clicked(gtx) {
				g.PlaceBet("player")
			}
			if bankerBetBtn.Clicked(gtx) {
				g.PlaceBet("banker")
			}
			if tieBetBtn.Clicked(gtx) {
				g.PlaceBet("tie")
			}
			if dealBtn.Clicked(gtx) {
				g.DealHand()
			}
			if menuBtn.Clicked(gtx) {
				return nil // Return to the BaccaratMain loop
			}
			if playAgainBtn.Clicked(gtx) {
				g.NewRound()
			}

			layout.Stack{}.Layout(gtx,
				layout.Expanded(imageWidget.Layout),
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					// Main flex container for the whole screen
					return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
						// Top Section: Points and Result
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return layout.Inset{Top: unit.Dp(20)}.Layout(gtx,
								func(gtx layout.Context) layout.Dimensions {
									return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceAround}.Layout(gtx,
										layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
											pointsLabel := material.H5(th, "Points: "+strconv.Itoa(g.UserPoints))
											pointsLabel.Color = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
											return pointsLabel.Layout(gtx)
										}),
										layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
											resultLabel := material.H5(th, g.LastResult)
											resultLabel.Color = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
											return resultLabel.Layout(gtx)
										}),
									)
								})
						}),

						// Middle Section: Cards (expands to fill space and centers content)
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
									layout.Rigid(renderHand(th, "Banker", g.BankerHand, g.CardsRevealed)),
									layout.Rigid(layout.Spacer{Height: unit.Dp(50)}.Layout),
									layout.Rigid(renderHand(th, "Player", g.PlayerHand, g.CardsRevealed)),
								)
							})
						}),

						// Button Section: Controls
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return layout.Inset{Bottom: unit.Dp(30)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								if g.CurrentPhase == "betting" {
									return renderBettingControls(th, &playerBetBtn, &bankerBetBtn, &tieBetBtn, &dealBtn, &menuBtn, g.UserBet)(gtx)
								}
								return renderResultControls(th, &playAgainBtn)(gtx)
							})
						}),
					)
				}),
			)
			e.Frame(gtx.Ops)
		}
	}
}
//deals and displays your hand and dealer hands. 
func renderHand(th *material.Theme, title string, hand []game.Card, revealed bool) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				titleLabel := material.H6(th, title)
				titleLabel.Color = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
				return titleLabel.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				var handToRender []string
				if revealed {
					for _, c := range hand {
						handToRender = append(handToRender, c.Rank+c.Suit)
					}
				} else {
					for range hand {
						handToRender = append(handToRender, "1B") // Card back image
					}
				}
				//selects and draws resources for cards used in other parts of application. 
				cardImages := ui2.GetCardImage(handToRender)
				children := make([]layout.FlexChild, len(cardImages))
				for i, img := range cardImages {
					cardImg := widget.Image{Src: paint.NewImageOp(img), Fit: widget.Contain}
					children[i] = layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						scale := f32.Pt(0.5, 0.5)
						macro := op.Record(gtx.Ops)
						dims := cardImg.Layout(gtx)
						call := macro.Stop()
						op.Affine(f32.Affine2D{}.Scale(f32.Point{}, scale)).Add(gtx.Ops)
						call.Add(gtx.Ops)
						return dims
					})
				}
				return layout.Flex{}.Layout(gtx, children...)
			}),
		)
	}
}
//game controls tied to buttons. 
func renderBettingControls(th *material.Theme, pBtn, bBtn, tBtn, dBtn, mBtn *widget.Clickable, currentBet string) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Alignment: layout.Middle, Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				betLabel := material.Body1(th, "Place Your Bet")
				betLabel.Color = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
				return betLabel.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Spacing: layout.SpaceEvenly}.Layout(gtx,
					layout.Rigid(button(th, pBtn, "Player", currentBet == "player")),
					layout.Rigid(button(th, bBtn, "Banker", currentBet == "banker")),
					layout.Rigid(button(th, tBtn, "Tie", currentBet == "tie")),
				)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Spacing: layout.SpaceEvenly}.Layout(gtx,
					layout.Rigid(button(th, dBtn, "Deal", false)),
					layout.Rigid(button(th, mBtn, "Menu", false)),
				)
			}),
		)
	}
}
//button options after win/loss
func renderResultControls(th *material.Theme, paBtn *widget.Clickable) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Center.Layout(gtx, button(th, paBtn, "Play Again", false))
	}
}
//button layout definition
func button(th *material.Theme, click *widget.Clickable, text string, selected bool) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		btn := material.Button(th, click, text)
		if selected {
			btn.Background = color.NRGBA{G: 0x80, A: 0xFF}
		}
		return btn.Layout(gtx)
	}
}

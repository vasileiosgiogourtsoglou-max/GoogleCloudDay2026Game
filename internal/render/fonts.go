package render

import (
	"bytes"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/gofont/goregular"
)

type FontManager struct {
	RegularSource *text.GoTextFaceSource
	BoldSource    *text.GoTextFaceSource

	TitleLarge  *text.GoTextFace
	TitleMedium *text.GoTextFace
	TitleSmall  *text.GoTextFace
	CardTitle   *text.GoTextFace
	CardTag     *text.GoTextFace
	CardBody    *text.GoTextFace
	HUD         *text.GoTextFace
}

var Fonts *FontManager

func InitFonts() {
	if Fonts != nil {
		return
	}
	regSrc, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		log.Fatalf("render: load regular font: %v", err)
	}
	boldSrc, err := text.NewGoTextFaceSource(bytes.NewReader(gobold.TTF))
	if err != nil {
		log.Fatalf("render: load bold font: %v", err)
	}

	Fonts = &FontManager{
		RegularSource: regSrc,
		BoldSource:    boldSrc,
		TitleLarge:    &text.GoTextFace{Source: boldSrc, Size: 22},
		TitleMedium:   &text.GoTextFace{Source: boldSrc, Size: 14},
		TitleSmall:    &text.GoTextFace{Source: boldSrc, Size: 11},
		CardTitle:     &text.GoTextFace{Source: boldSrc, Size: 11},
		CardTag:       &text.GoTextFace{Source: boldSrc, Size: 9.5},
		CardBody:      &text.GoTextFace{Source: regSrc, Size: 9.5},
		HUD:           &text.GoTextFace{Source: boldSrc, Size: 11},
	}
}

func DrawCentered(dst *ebiten.Image, msg string, face *text.GoTextFace, cx, cy float64, clr color.Color) {
	w, h := text.Measure(msg, face, 0)
	op := &text.DrawOptions{}
	op.GeoM.Translate(mathRound(cx-w/2.0), mathRound(cy-h/2.0))
	op.ColorScale.ScaleWithColor(clr)
	text.Draw(dst, msg, face, op)
}

func mathRound(val float64) float64 {
	return float64(int(val + 0.5))
}

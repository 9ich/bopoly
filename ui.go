package main

import (
	"fmt"
	"image/color"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/ebitenutil"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	eb "github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
)

var (
	normFont   font.Face
	bigFont    font.Face
	whiteImage *eb.Image
	richBlack  = color.RGBA{0x01, 0x0b, 0x13, 0xff}
	paper      = color.RGBA{0xfd, 0xf5, 0xe8, 0xff}
)

type Game struct{}

func (g *Game) Update(screen *eb.Image) error {
	if eb.IsDrawingSkipped() {
		return nil
	}

	sq := mkDeedCard(&board[int(time.Now().UnixNano()/int64(time.Second/2))%len(board)])
	if err := screen.DrawImage(sq, &eb.DrawImageOptions{}); err != nil {
		log.Fatal(err)
	}
	return nil
}

func (g *Game) Layout(outW, outH int) (int, int) {
	return outW, outH
}

func testAlign(img *eb.Image) {
	y := 100
	ebitenutil.DrawLine(img, 0, float64(y), 600, float64(y), color.RGBA{255, 0, 255, 255})
	drawText(img, "top-_xTgl", bigFont, 0, y, left, top, color.White)
	drawText(img, "center-_xTgl", bigFont, 200, y, left, center, color.White)
	drawText(img, "bottom-_xTgl", bigFont, 400, y, left, bottom, color.White)
}

func ui() {
	uiinit()
	eb.SetWindowSize(640, 480)
	eb.SetWindowTitle("bopoly")
	g := &Game{}
	if err := eb.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func uiinit() {
	whiteImage, _ = eb.NewImage(2, 2, ebiten.FilterNearest)
	whiteImage.Fill(color.White)

	ttf, err := ioutil.ReadFile("Alata-Regular.ttf")
	if err != nil {
		log.Fatal(err)
	}

	tt, err := truetype.Parse(ttf)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	normFont = truetype.NewFace(tt, &truetype.Options{
		Size:    18,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	bigFont = truetype.NewFace(tt, &truetype.Options{
		Size:    32,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
}

func mkDeedCard(sq *Square) *eb.Image {
	const W, H, margin = 318, 500, 5
	bg := paper
	fg := richBlack
	img, _ := eb.NewImage(W, H, eb.FilterDefault)
	img.Fill(bg)
	drawRect(img, 0+margin, 0+margin, W-margin*2, H-margin*2, fg)
	{
		const (
			h, margin = 40, margin * 2
			x, w      = margin, W - 2*margin
		)
		clr := groupColor(sq.group)
		ebitenutil.DrawRect(img, x, x, w, h, clr)
		drawRect(img, x, x, w, h, fg)
		drawText(img, strings.ToUpper(sq.name), normFont, W/2, h/2+margin, center, center, juxta(clr, bg, fg))
	}
	for i := range sq.rent {
		const y, yadv, pad = 70, 30, 60
		houses := ""
		switch i {
		default:
			houses = fmt.Sprintf("%d houses", i)
		case 0:
			houses = "no houses"
		case 1:
			houses = "1 house"
		case 5:
			houses = "a hotel"
		}
		s := fmt.Sprintf("With %s", houses)
		drawText(img, s, normFont, 0+pad, y+i*yadv, left, top, fg)
		r := fmt.Sprintf("%c%d", bux, sq.rent[i])
		drawText(img, r, normFont, W-pad, y+i*yadv, right, top, fg)
	}
	return img
}

func groupColor(g int) color.Color {
	switch g {
	case SBlue:
		return color.RGBA{121, 159, 235, 0xff}
	case SGreen:
		return color.RGBA{50, 152, 42, 0xff}
	case SNavy:
		return color.RGBA{0, 47, 104, 0xff}
	case SOrange:
		return color.RGBA{220, 78, 45, 0xff}
	case SYellow:
		return color.RGBA{246, 244, 84, 0xff}
	case SPink:
		return color.RGBA{230, 68, 223, 0xff}
	case SBrown:
		return color.RGBA{128, 14, 48, 0xff}
	default:
		return richBlack
	}
}

func juxta(clr, light, dark color.Color) color.Color {
	r, g, b, _ := clr.RGBA()
	if (r+g+b)/3 > 0x7fff {
		return dark
	}
	return light
}

func drawRect(img *eb.Image, x, y, w, h int, clr color.Color) {
	fx, fy, fw, fh := float64(x), float64(y), float64(w), float64(h)
	ebitenutil.DrawLine(img, fx, fy, fx+fw, fy, clr)
	ebitenutil.DrawLine(img, fx+fw, fy, fx+fw, fy+fh, clr)
	ebitenutil.DrawLine(img, fx+fw, fy+fh, fx, fy+fh, clr)
	ebitenutil.DrawLine(img, fx, fy+fh, fx, fy, clr)
}

func drawText(dst *eb.Image, s string, face font.Face, x, y int, halign, valign int, clr color.Color) {
	xofs, yofs := alignText(s, face, halign, valign)
	text.Draw(dst, s, face, x+xofs, y+yofs, clr)
}

const (
	left = iota
	center
	right
	top    = left
	bottom = right
)

func alignText(s string, face font.Face, horz, vert int) (xofs int, yofs int) {
	b, _ := font.BoundString(face, s)
	w := (b.Max.X - b.Min.X).Floor()
	gb, _, _ := face.GlyphBounds('S')
	h := (gb.Max.Y - gb.Min.Y).Round()
	return align(w, h, horz, vert)
}

func align(w, h int, horz, vert int) (xofs int, yofs int) {
	switch horz {
	case center:
		xofs = -w / 2
	case right:
		xofs = -w
	}
	switch vert {
	case center:
		yofs = h / 2
	case top:
		yofs = h
	}
	return xofs, yofs
}

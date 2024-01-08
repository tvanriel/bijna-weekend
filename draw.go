package bijnaweekend

import (
	"image"
	"image/color"
	"io"
	"math/rand"

	"github.com/fogleman/gg"
	"github.com/muesli/gamut"
	"golang.org/x/image/font"
)

type Mascotte struct {
	X, Y  float64
	Image image.Image
}

type Tagline struct {
	X, Y    float64
	Text    string
	Font    font.Face
	Palette TextColor
}

type TextColor struct {
	Foreground color.Color
	Shadow     color.Color
}

type GradientColors struct {
	From color.Color
	To   color.Color
}

type Config struct {
	Height, Width int
	Mascotte      Mascotte
	Tagline       Tagline
	Palette       GradientColors
	Writer        io.Writer
}

func PastelTextColor() (TextColor, error) {
	p, err := Pastel()
	if err != nil {
		return TextColor{}, err
	}
	return TextColor{
		Foreground: p[0],
		Shadow:     gamut.Darker(p[0], 0.3),
	}, nil
}

func PastelGradient() (GradientColors, error) {
	p, err := Pastel()
	if err != nil {
		return GradientColors{}, nil
	}

	return GradientColors{
		From: p[0],
		To:   p[1],
	}, nil
}

func NewTagline(x, y float64, filename string, points float64, text string, palette TextColor) (Tagline, error) {

	font, err := gg.LoadFontFace(filename, points)

	return Tagline{
		X:       x,
		Y:       y,
		Text:    text,
		Font:    font,
		Palette: palette,
	}, err
}

func NewMascotte(x, y float64, filename string) (Mascotte, error) {
	im, err := gg.LoadImage(filename)

	return Mascotte{
		X:     x,
		Y:     y,
		Image: im,
	}, err
}

func BijnaWeekend(c *Config) error {
	ctx := gg.NewContext(c.Width, c.Height)
	drawBackground(ctx, c.Palette)
	drawMascotte(ctx, c.Mascotte)
	drawTagline(ctx, c.Tagline)

	return ctx.EncodePNG(c.Writer)
}

func Pastel() (*[2]color.Color, error) {
	p, err := gamut.Generate(2, gamut.PastelGenerator{})
	if err != nil {
		return nil, err
	}
	return &[2]color.Color{
		p[0], p[1],
	}, nil
}

func drawBackground(ctx *gg.Context, palette GradientColors) {
	bgPick := rand.Intn(3)

	var bg gg.Gradient
	switch bgPick {
	case 0:
		bg = gg.NewLinearGradient(0, 0, 0, float64(ctx.Width()))
	case 1:
		bg = gg.NewLinearGradient(0, 0, float64(ctx.Height()), float64(ctx.Width()))
	case 2:
		bg = gg.NewLinearGradient(0, float64(ctx.Width()), 0, 0)
	}
	bg.AddColorStop(0, palette.From)
	bg.AddColorStop(1, palette.To)
	ctx.SetFillStyle(bg)
	ctx.DrawRectangle(0, 0, float64(ctx.Width()), float64(ctx.Height()))
	ctx.Fill()
}

func drawMascotte(ctx *gg.Context, mascotte Mascotte) {
	ctx.MoveTo(
		mascotte.X,
		mascotte.Y,
	)
	ctx.DrawImage(mascotte.Image, 0, 0)
}

func drawTagline(ctx *gg.Context, t Tagline) {
	ctx.SetFillStyle(gg.NewSolidPattern(t.Palette.Shadow))

	distance := 2.0
	ctx.DrawString(
		t.Text,
		t.X+distance,
		t.Y+distance,
	)
	ctx.SetFillStyle(gg.NewSolidPattern(t.Palette.Foreground))
	ctx.DrawString(
		t.Text,
		t.X,
		t.Y,
	)

}

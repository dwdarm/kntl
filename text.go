package kntl

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type Font interface {
	LoadFromFile(path string, size int) error
	GetTTFFont() *ttf.Font
	Destory()
}

type FontImp struct {
	font *ttf.Font
	size int
}

func NewFont() Font {
	return &FontImp{}
}

func (f *FontImp) LoadFromFile(path string, size int) error {
	font, err := ttf.OpenFont(path, size)
	if err != nil {
		return err
	}

	f.font = font
	f.size = size

	return nil
}

func (f *FontImp) GetTTFFont() *ttf.Font {
	return f.font
}

func (f *FontImp) Destory() {
	if f.font != nil {
		f.font.Close()
	}
}

type Text interface {
	SetText(text string)
	SetColor(color *Color)
	SetOutlineSize(size int)
	SetOutlineColor(color *Color)
	SetPosition(pos *Vector2)
	GetSize() *Vector2
	BuildSurface()
	Draw()
	Destroy()
}

type TextImp struct {
	game         Game
	font         Font
	text         string
	color        Color
	outlineSize  int
	outlineColor Color
	position     Vector2
	surface      *sdl.Surface
}

func NewText(game Game, font Font, text string, color *Color) Text {
	return &TextImp{
		game:  game,
		font:  font,
		text:  text,
		color: *color,
	}
}

func (t *TextImp) SetText(text string) {
	t.text = text
}

func (t *TextImp) SetColor(color *Color) {
	t.color = *color
}

func (t *TextImp) SetOutlineSize(size int) {
	t.outlineSize = size
}

func (t *TextImp) SetOutlineColor(color *Color) {
	t.outlineColor = *color
}

func (t *TextImp) SetPosition(pos *Vector2) {
	t.position = *pos
}

func (t *TextImp) GetSize() *Vector2 {
	size := &Vector2{}

	if t.surface != nil {
		size.X = float32(t.surface.W)
		size.Y = float32(t.surface.H)
	}

	return size
}

func (t *TextImp) BuildSurface() {
	if t.surface != nil {
		t.surface.Free()
	}

	font := t.font.GetTTFFont()
	font.SetOutline(0)

	if t.outlineSize <= 0 {
		surface, err := font.RenderUTF8Blended(t.text, *t.color.ToSDLColor())
		if err != nil {
			panic(err)
		}

		t.surface = surface
	} else {
		fg, err := font.RenderUTF8Blended(t.text, *t.color.ToSDLColor())
		if err != nil {
			panic(err)
		}
		defer fg.Free()

		font.SetOutline(t.outlineSize)
		bg, err := font.RenderUTF8Blended(t.text, *t.outlineColor.ToSDLColor())
		if err != nil {
			panic(err)
		}

		fg.SetBlendMode(sdl.BLENDMODE_BLEND)
		fg.Blit(nil, bg, &sdl.Rect{
			X: int32(t.outlineSize),
			Y: int32(t.outlineSize),
			W: fg.W,
			H: fg.H,
		})

		t.surface = bg
	}
}

func (t *TextImp) Draw() {
	renderer := t.game.GetRenderer()
	surface := t.surface

	texture, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		panic(err)
	}
	defer texture.Destroy()

	renderer.Copy(texture, nil, &sdl.Rect{
		X: int32(t.position.X),
		Y: int32(t.position.Y),
		W: surface.W,
		H: surface.H,
	})
}

func (t *TextImp) Destroy() {
	if t.surface != nil {
		t.surface.Free()
	}
}

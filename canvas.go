// Pix
package pix3d

import (
	"image"
	"image/color"
	"math"

	"github.com/disintegration/imaging"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// 2D.
// Operation2D — функция, которая рисует на *image.NRGBA.
type Operation2D func(*image.NRGBA)

// Canvas — холст с набором 2D и 3D-операций.
type Canvas struct {
	img              *image.NRGBA
	width            int
	height           int
	op               []Operation2D
	opCount          int
	Scale            int // Размер модели(экранный)
	CenterX, CenterY int
	ModelColor       color.Color
}

// 2D.
// NewCanvas создаёт новый холст заданного размера, заполненный чёрным.
func NewCanvas(width, height int) *Canvas {
	return &Canvas{
		img:     imaging.New(width, height, color.Black),
		width:   width,
		height:  height,
		Scale:   800,
		CenterX: width / 2,
		CenterY: height / 2,
	}
}

// 2D.
// Fill заливает холст указанным цветом.
func (c *Canvas) Fill(clr color.Color) {
	for i := range len(c.img.Pix) >> 2 {
		r, g, b, a := clr.RGBA()
		c.img.Pix[i<<2] = uint8(r >> 8)
		c.img.Pix[i<<2+1] = uint8(g >> 8)
		c.img.Pix[i<<2+2] = uint8(b >> 8)
		c.img.Pix[i<<2+3] = uint8(a >> 8)
	}
}

// Save сохраняет холст в файл изобравжения(PNG, JPEG).
func (c *Canvas) Save(filename string) {
	imaging.Save(c.img, filename)
}

// Image возвращает текущее изображение холста.
func (c *Canvas) Image() image.Image {
	return c.img
}

// 2D.
// Add запускае.
func (c *Canvas) Op(op Operation2D) {
	op(c.img)
}

// 2D.
// FillSquare возвращает операцию рисования залитого прямоугольника.
func FillSquare(x1, y1, x2, y2 int, col color.Color) Operation2D {
	return func(img *image.NRGBA) {
		for y := y1; y <= y2; y++ {
			for x := x1; x <= x2; x++ {
				img.Set(x, y, col)
			}
		}
	}
}

// 2D.
// FillCircle возвращает операцию рисования залитого круга.
func FillCircle(x, y, radius int, col color.Color) Operation2D {
	return func(img *image.NRGBA) {
		r2 := radius * radius
		for dy := -radius; dy <= radius; dy++ {
			dx := int(math.Sqrt(float64(r2 - dy*dy)))
			x1 := x - dx
			x2 := x + dx
			yCoord := y + dy
			if yCoord < 0 || yCoord >= img.Bounds().Max.Y {
				continue
			}
			for xi := x1; xi <= x2; xi++ {
				if xi >= 0 && xi < img.Bounds().Max.X {
					img.Set(xi, yCoord, col)
				}
			}
		}
	}
}

// 2D.
// DrawLine возвращает операцию рисования линии (алгоритм Брезенхема).
func DrawLine(x1, y1, x2, y2 int, col color.Color) Operation2D {
	return func(img *image.NRGBA) {
		dx := abs(x2 - x1)
		dy := abs(y2 - y1)
		sx := 1
		if x1 > x2 {
			sx = -1
		}
		sy := 1
		if y1 > y2 {
			sy = -1
		}
		err := dx - dy

		for {
			img.Set(x1, y1, col)
			if x1 == x2 && y1 == y2 {
				break
			}
			e2 := err * 2
			if e2 > -dy {
				err -= dy
				x1 += sx
			}
			if e2 < dx {
				err += dx
				y1 += sy
			}
		}
	}
}

// 2D.
// DrawText возвращает операцию рисования текста (использует встроенный шрифт 7x13).
func DrawText(clr color.Color, label string, x, y int) Operation2D {
	return func(img *image.NRGBA) {
		point := fixed.Point26_6{
			X: fixed.Int26_6(x * 64),
			Y: fixed.Int26_6(y * 64),
		}
		d := &font.Drawer{
			Dst:  img,
			Src:  image.NewUniform(clr),
			Face: basicfont.Face7x13,
			Dot:  point,
		}
		d.DrawString(label)
	}
}

// 2D.
// DrawEllipse рисует контур эллипса.
func DrawEllipse(cx, cy, rx, ry int, col color.Color) Operation2D {
	return func(img *image.NRGBA) {
		steps := rx + ry
		if steps < 360 {
			steps = 360
		}
		for i := 0; i <= steps; i++ {
			angle := 2 * math.Pi * float64(i) / float64(steps)
			x := cx + int(float64(rx)*math.Cos(angle))
			y := cy + int(float64(ry)*math.Sin(angle))
			if x >= 0 && x < img.Bounds().Max.X && y >= 0 && y < img.Bounds().Max.Y {
				img.Set(x, y, col)
			}
		}
	}
}

// 2D.
// DrawScaledCircle рисует контур эллипса, заданного множителями.
func DrawScaledCircle(cx, cy, radius int, mulX, mulY float64, col color.Color) Operation2D {
	rx := int(math.Abs(float64(radius) * mulX))
	ry := int(math.Abs(float64(radius) * mulY))
	if rx == 0 && ry == 0 {
		rx, ry = 1, 1
	}
	return DrawEllipse(cx, cy, rx, ry, col)
}

// 2D.
// Fill заливает холст цветом (используйте как Operation2D).
func Fill(clr color.Color) Operation2D {
	return func(img *image.NRGBA) {
		for i := range len(img.Pix) >> 2 {
			r, g, b, a := clr.RGBA()
			img.Pix[i<<2] = uint8(r >> 8)
			img.Pix[i<<2+1] = uint8(g >> 8)
			img.Pix[i<<2+2] = uint8(b >> 8)
			img.Pix[i<<2+3] = uint8(a >> 8)
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

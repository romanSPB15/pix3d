# Pix3D
Легковесный 3D-движок, написанный на чистом Go. Может быть совмещён с различными GUI библиотеками. Вычисления проводятся на CPU.
<br><br>
<img src="demo.png" alt="demo.png"><br><br>
# Особенности
* Полностью поддерживает Wavefront OBJ.
* Источников света может быть несколько.

## Установка
```bash
go get -u github.com/romanSPB15/pix3d
```

## Быстрый старт

```
package main

import (
	"image/color"
	"log"
	"math"

	"github.com/romanSPB15/pix3d"
)

func main() {
	cnv := pix3d.NewCanvas(1000, 1000)
	cnv.Fill(color.RGBA{30, 30, 30, 255})

	tris, err := pix3d.ParseOBJ("stanford-bunny.obj")
	if err != nil {
		log.Fatal(err)
	}

	tris = pix3d.CenterAndScaleModel(tris, 1.0)

	cnv.Scale = 1800

	cnv.DrawModel(tris, pix3d.Yellow)
	cnv.Save("render.png")
}

```

## Производительность
* У меня на крупных моделях 4.6 млн. тругольников в секунду
* На мелких моделях 3~4 млн.
* С GUI(например, Fyne) FPS падает в несколько раз, так что эти цифры - производительность самого движка.
* В ветке `feature/indexed-geometry` используется индексированная геометрия, так что потребление памяти упало в 3 раза.
* Потребление памяти `feature/indexed-geometry` от 128 байт на треугольник.
> Все тесты проводились на разрешении 1000x1000.

Графика основана на [https://github.com/romanSPB15/pixgl](PixGL)

## Лицензия
`MIT`

package pix3d

import (
	"bufio"
	"image/color"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

// Стандартные цвета
var (
	Red    = color.RGBA{255, 0, 0, 255}
	Green  = color.RGBA{0, 255, 0, 255}
	Blue   = color.RGBA{0, 0, 255, 255}
	Yellow = color.RGBA{255, 255, 0, 255}
	Violet = color.RGBA{0, 255, 255, 255}
)

type Vec3 struct{ X, Y, Z float64 }

type Triangle struct {
	Verts  [3]Vec3
	Normal Vec3
}

type Light struct {
	Position  Vec3        // Направление
	Color     color.Color // Цветвета
	Intensity float64     // Яркость (0..1)
	Ambient   float64     // Фоновое освещение (0..1)
	Diffuse   float64     // Коэффициент рассеянного света
}

type Mesh struct {
	Vertices []Vec3
	Normals  []Vec3
	Indices  []int
}

// RotateX поворачивает все вершины меша вокруг оси X на угол angle (радианы).
func (m *Mesh) RotateX(angle float64) {
	cos := math.Cos(angle)
	sin := math.Sin(angle)
	for i := range m.Vertices {
		y := m.Vertices[i].Y
		z := m.Vertices[i].Z
		m.Vertices[i].Y = y*cos - z*sin
		m.Vertices[i].Z = y*sin + z*cos
	}

	for i := range m.Normals {
		y := m.Normals[i].Y
		z := m.Normals[i].Z
		m.Normals[i].Y = y*cos - z*sin
		m.Normals[i].Z = y*sin + z*cos
	}
}

// RotateY поворачивает все вершины меша вокруг оси Y на угол angle (радианы).
func (m *Mesh) RotateY(angle float64) {
	cos := math.Cos(angle)
	sin := math.Sin(angle)
	for i := range m.Vertices {
		x := m.Vertices[i].X
		z := m.Vertices[i].Z
		m.Vertices[i].X = x*cos + z*sin
		m.Vertices[i].Z = -x*sin + z*cos
	}
	// Если есть вершинные нормали – поворачиваем и их
	for i := range m.Normals {
		x := m.Normals[i].X
		z := m.Normals[i].Z
		m.Normals[i].X = x*cos + z*sin
		m.Normals[i].Z = -x*sin + z*cos
		// Нормализация не требуется, если они были единичными
	}
}

// RotateZ поворачивает все вершины меша вокруг оси Z на угол angle (радианы).
func (m *Mesh) RotateZ(angle float64) {
	cos := math.Cos(angle)
	sin := math.Sin(angle)
	for i := range m.Vertices {
		x := m.Vertices[i].X
		y := m.Vertices[i].Y
		m.Vertices[i].X = x*cos - y*sin
		m.Vertices[i].Y = x*sin + y*cos
	}
	for i := range m.Normals {
		x := m.Normals[i].X
		y := m.Normals[i].Y
		m.Normals[i].X = x*cos - y*sin
		m.Normals[i].Y = x*sin + y*cos
	}
}

func triangleNormal(v0, v1, v2 Vec3) Vec3 {
	e1 := Vec3{v1.X - v0.X, v1.Y - v0.Y, v1.Z - v0.Z}
	e2 := Vec3{v2.X - v0.X, v2.Y - v0.Y, v2.Z - v0.Z}
	nx := e1.Y*e2.Z - e1.Z*e2.Y
	ny := e1.Z*e2.X - e1.X*e2.Z
	nz := e1.X*e2.Y - e1.Y*e2.X
	len := math.Sqrt(nx*nx + ny*ny + nz*nz)
	if len == 0 {
		return Vec3{0, 0, 0}
	}
	return Vec3{nx / len, ny / len, nz / len}
}

func LoadMesh(filename string) (*Mesh, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var vertices []Vec3
	var indices []int

	type faceVert struct{ v int }
	var currentFace []faceVert

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		switch fields[0] {
		case "v":
			if len(fields) < 4 {
				continue
			}
			x, _ := strconv.ParseFloat(fields[1], 64)
			y, _ := strconv.ParseFloat(fields[2], 64)
			z, _ := strconv.ParseFloat(fields[3], 64)
			vertices = append(vertices, Vec3{x, y, z})
		case "f":
			currentFace = currentFace[:0]
			for _, part := range fields[1:] {
				parts := strings.Split(part, "/")
				vIdx, _ := strconv.Atoi(parts[0])
				if vIdx < 0 {
					vIdx = len(vertices) + vIdx + 1
				}
				currentFace = append(currentFace, faceVert{v: vIdx})
			}
			for i := 1; i+1 < len(currentFace); i++ {
				indices = append(indices,
					currentFace[0].v-1,
					currentFace[i].v-1,
					currentFace[i+1].v-1)
			}
		}
	}

	mesh := &Mesh{
		Vertices: vertices,
		Indices:  indices,
		Normals:  nil,
	}
	return mesh, scanner.Err()
}

func (m *Mesh) CenterAndScale(targetSize float64) {
	if len(m.Vertices) == 0 {
		return
	}
	minX, maxX := math.MaxFloat64, -math.MaxFloat64
	minY, maxY := math.MaxFloat64, -math.MaxFloat64
	minZ, maxZ := math.MaxFloat64, -math.MaxFloat64
	for _, v := range m.Vertices {
		if v.X < minX {
			minX = v.X
		}
		if v.X > maxX {
			maxX = v.X
		}
		if v.Y < minY {
			minY = v.Y
		}
		if v.Y > maxY {
			maxY = v.Y
		}
		if v.Z < minZ {
			minZ = v.Z
		}
		if v.Z > maxZ {
			maxZ = v.Z
		}
	}
	centerX := (minX + maxX) / 2
	centerY := (minY + maxY) / 2
	centerZ := (minZ + maxZ) / 2
	sizeX := maxX - minX
	sizeY := maxY - minY
	sizeZ := maxZ - minZ
	maxSize := math.Max(sizeX, math.Max(sizeY, sizeZ))
	if maxSize == 0 {
		maxSize = 1
	}
	scale := targetSize / maxSize
	for i := range m.Vertices {
		m.Vertices[i].X = (m.Vertices[i].X - centerX) * scale
		m.Vertices[i].Y = (m.Vertices[i].Y - centerY) * scale
		m.Vertices[i].Z = (m.Vertices[i].Z - centerZ) * scale
	}
}

func (c *Canvas) DrawTriangle(tri Triangle, clr color.Color) {
	cameraZ := 3.0

	project := func(v Vec3) (int, int) {
		dz := cameraZ - v.Z
		if dz <= 0 {
			return -1, -1
		}
		factor := float64(c.Scale) / dz
		sx := int(float64(c.CenterX) + v.X*factor)
		sy := int(float64(c.CenterY) - v.Y*factor)
		return sx, sy
	}

	p0x, p0y := project(tri.Verts[0])
	p1x, p1y := project(tri.Verts[1])
	p2x, p2y := project(tri.Verts[2])

	w := c.img.Bounds().Max.X
	h := c.img.Bounds().Max.Y
	if (p0x < 0 && p1x < 0 && p2x < 0) ||
		(p0x >= w && p1x >= w && p2x >= w) ||
		(p0y < 0 && p1y < 0 && p2y < 0) ||
		(p0y >= h && p1y >= h && p2y >= h) {
		return
	}

	var totalR, totalG, totalB float64

	for _, light := range c.Lights {
		lightDir := light.Position
		lenL := math.Sqrt(lightDir.X*lightDir.X + lightDir.Y*lightDir.Y + lightDir.Z*lightDir.Z)
		if lenL > 0 {
			lightDir.X /= lenL
			lightDir.Y /= lenL
			lightDir.Z /= lenL
		}

		dot := tri.Normal.X*lightDir.X + tri.Normal.Y*lightDir.Y + tri.Normal.Z*lightDir.Z
		if dot < 0 {
			dot = 0
		}
		intensity := light.Ambient + light.Diffuse*dot
		intensity *= light.Intensity

		lr, lg, lb, _ := light.Color.RGBA()
		lr8, lg8, lb8 := float64(lr>>8), float64(lg>>8), float64(lb>>8)

		totalR += lr8 * intensity
		totalG += lg8 * intensity
		totalB += lb8 * intensity
	}

	if totalR > 255 {
		totalR = 255
	}
	if totalG > 255 {
		totalG = 255
	}
	if totalB > 255 {
		totalB = 255
	}

	rM, gM, bM, _ := clr.RGBA()
	r8, g8, b8 := float64(rM>>8), float64(gM>>8), float64(bM>>8)

	r := uint8(r8 * totalR / 255)
	g := uint8(g8 * totalG / 255)
	b := uint8(b8 * totalB / 255)

	x1, y1, x2, y2, x3, y3 := p0x, p0y, p1x, p1y, p2x, p2y
	if y1 > y2 {
		x1, x2 = x2, x1
		y1, y2 = y2, y1
	}
	if y1 > y3 {
		x1, x3 = x3, x1
		y1, y3 = y3, y1
	}
	if y2 > y3 {
		x2, x3 = x3, x2
		y2, y3 = y3, y2
	}

	interp := func(x1, y1, x2, y2, y int) int {
		if y1 == y2 {
			return x1
		}
		t := float64(y-y1) / float64(y2-y1)
		return x1 + int(float64(x2-x1)*t)
	}

	bounds := c.img.Bounds()
	for y := y1; y <= y3 && y < bounds.Max.Y; y++ {
		if y < 0 {
			continue
		}
		var xStart, xEnd int
		if y < y2 {
			xStart = interp(x1, y1, x2, y2, y)
			xEnd = interp(x1, y1, x3, y3, y)
		} else {
			xStart = interp(x2, y2, x3, y3, y)
			xEnd = interp(x1, y1, x3, y3, y)
		}
		if xStart > xEnd {
			xStart, xEnd = xEnd, xStart
		}
		if xStart < 0 {
			xStart = 0
		}
		if xEnd >= bounds.Max.X {
			xEnd = bounds.Max.X - 1
		}
		offset := c.img.PixOffset(xStart, y)
		for x := xStart; x <= xEnd; x++ {
			c.img.Pix[offset] = r
			c.img.Pix[offset+1] = g
			c.img.Pix[offset+2] = b
			c.img.Pix[offset+3] = 255
			offset += 4
		}
	}
}

func (c *Canvas) DrawMesh(mesh *Mesh, clr color.Color) {
	triCount := len(mesh.Indices) / 3
	type triDepth struct {
		idx   int
		depth float64
	}
	depths := make([]triDepth, triCount)
	for i := 0; i < triCount; i++ {
		i0 := mesh.Indices[3*i]
		i1 := mesh.Indices[3*i+1]
		i2 := mesh.Indices[3*i+2]
		v0 := mesh.Vertices[i0]
		v1 := mesh.Vertices[i1]
		v2 := mesh.Vertices[i2]
		avgZ := (v0.Z + v1.Z + v2.Z) / 3.0
		depths[i] = triDepth{idx: 3 * i, depth: avgZ}
	}
	sort.Slice(depths, func(i, j int) bool { return depths[i].depth < depths[j].depth })

	for _, td := range depths {
		i0 := mesh.Indices[td.idx]
		i1 := mesh.Indices[td.idx+1]
		i2 := mesh.Indices[td.idx+2]
		tri := Triangle{
			Verts: [3]Vec3{
				mesh.Vertices[i0],
				mesh.Vertices[i1],
				mesh.Vertices[i2],
			},
		}
		tri.Normal = triangleNormal(tri.Verts[0], tri.Verts[1], tri.Verts[2])
		c.DrawTriangle(tri, clr)
	}
}

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

func RotateX(tris []Triangle, angle float64) []Triangle {
	cos := math.Cos(angle)
	sin := math.Sin(angle)
	res := make([]Triangle, len(tris))
	for i, tri := range tris {
		var newVerts [3]Vec3
		for j, p := range tri.Verts {
			newVerts[j] = Vec3{
				X: p.X,
				Y: p.Y*cos - p.Z*sin,
				Z: p.Y*sin + p.Z*cos,
			}
		}
		n := tri.Normal
		newNormal := Vec3{
			X: n.X,
			Y: n.Y*cos - n.Z*sin,
			Z: n.Y*sin + n.Z*cos,
		}
		lenN := math.Sqrt(newNormal.X*newNormal.X + newNormal.Y*newNormal.Y + newNormal.Z*newNormal.Z)
		if lenN > 0 {
			newNormal.X /= lenN
			newNormal.Y /= lenN
			newNormal.Z /= lenN
		}
		res[i] = Triangle{Verts: newVerts, Normal: newNormal}
	}
	return res
}

func RotateY(tris []Triangle, angle float64) []Triangle {
	cos := math.Cos(angle)
	sin := math.Sin(angle)
	res := make([]Triangle, len(tris))
	for i, tri := range tris {
		var newVerts [3]Vec3
		for j, p := range tri.Verts {
			newVerts[j] = Vec3{
				X: p.X*cos + p.Z*sin,
				Y: p.Y,
				Z: -p.X*sin + p.Z*cos,
			}
		}
		n := tri.Normal
		newNormal := Vec3{
			X: n.X*cos + n.Z*sin,
			Y: n.Y,
			Z: -n.X*sin + n.Z*cos,
		}
		lenN := math.Sqrt(newNormal.X*newNormal.X + newNormal.Y*newNormal.Y + newNormal.Z*newNormal.Z)
		if lenN > 0 {
			newNormal.X /= lenN
			newNormal.Y /= lenN
			newNormal.Z /= lenN
		}
		res[i] = Triangle{Verts: newVerts, Normal: newNormal}
	}
	return res
}

func RotateZ(tris []Triangle, angle float64) []Triangle {
	cos := math.Cos(angle)
	sin := math.Sin(angle)
	res := make([]Triangle, len(tris))
	for i, tri := range tris {
		var newVerts [3]Vec3
		for j, p := range tri.Verts {
			newVerts[j] = Vec3{
				X: p.X*cos - p.Y*sin,
				Y: p.X*sin + p.Y*cos,
				Z: p.Z,
			}
		}
		n := tri.Normal
		newNormal := Vec3{
			X: n.X*cos - n.Y*sin,
			Y: n.X*sin + n.Y*cos,
			Z: n.Z,
		}
		lenN := math.Sqrt(newNormal.X*newNormal.X + newNormal.Y*newNormal.Y + newNormal.Z*newNormal.Z)
		if lenN > 0 {
			newNormal.X /= lenN
			newNormal.Y /= lenN
			newNormal.Z /= lenN
		}
		res[i] = Triangle{Verts: newVerts, Normal: newNormal}
	}
	return res
}

func TriangleNormal(v0, v1, v2 Vec3) Vec3 {
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

func ParseOBJ(filename string) ([]Triangle, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var vertices []Vec3
	var triangles []Triangle
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
			vertices = append(vertices, Vec3{X: x, Y: y, Z: z})
		case "f":
			if len(fields) < 4 {
				continue
			}
			indices := make([]int, 0, len(fields)-1)
			for _, part := range fields[1:] {
				idxStr := strings.Split(part, "/")[0]
				idx, _ := strconv.Atoi(idxStr)
				if idx > 0 {
					idx--
				} else if idx < 0 {
					idx = len(vertices) + idx
				} else {
					continue
				}
				if idx < 0 || idx >= len(vertices) {
					continue
				}
				indices = append(indices, idx)
			}
			if len(indices) < 3 {
				continue
			}
			for i := 1; i < len(indices)-1; i++ {
				tri := Triangle{
					Verts: [3]Vec3{
						vertices[indices[0]],
						vertices[indices[i]],
						vertices[indices[i+1]],
					},
				}
				// Вычисляем нормаль сразу (пока без центрирования)
				tri.Normal = TriangleNormal(tri.Verts[0], tri.Verts[1], tri.Verts[2])
				triangles = append(triangles, tri)
			}

		}
	}
	return triangles, scanner.Err()
}

func CenterAndScaleModel(tris []Triangle, targetSize float64) []Triangle {
	if len(tris) == 0 {
		return tris
	}
	minX, maxX := math.MaxFloat64, -math.MaxFloat64
	minY, maxY := math.MaxFloat64, -math.MaxFloat64
	minZ, maxZ := math.MaxFloat64, -math.MaxFloat64
	for _, tri := range tris {
		for _, v := range tri.Verts {
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
	result := make([]Triangle, len(tris))
	for i, tri := range tris {
		var newVerts [3]Vec3
		for j, v := range tri.Verts {
			newVerts[j] = Vec3{
				X: (v.X - centerX) * scale,
				Y: (v.Y - centerY) * scale,
				Z: (v.Z - centerZ) * scale,
			}
		}
		newNormal := TriangleNormal(newVerts[0], newVerts[1], newVerts[2])
		result[i] = Triangle{Verts: newVerts, Normal: newNormal}
	}
	return result
}

func (c *Canvas) DrawTriangle(tri Triangle, clr color.Color) {
	cameraZ := 3.0
	lightDir := Vec3{1, 1, 1}
	lenLight := math.Sqrt(lightDir.X*lightDir.X + lightDir.Y*lightDir.Y + lightDir.Z*lightDir.Z)
	lightDir.X /= lenLight
	lightDir.Y /= lenLight
	lightDir.Z /= lenLight

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
	if p0x < 0 || p1x < 0 || p2x < 0 {
		return
	}

	normal := tri.Normal
	intensity := normal.X*lightDir.X + normal.Y*lightDir.Y + normal.Z*lightDir.Z
	if intensity < 0 {
		intensity = 0
	}
	intensity = 0.2 + 0.8*intensity

	rM, gM, bM, _ := clr.RGBA()

	r := uint8(float64(rM) * intensity)
	g := uint8(float64(gM) * intensity)
	b := uint8(float64(bM) * intensity)

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

func (c *Canvas) DrawModel(tris []Triangle, clr color.Color) {
	type TriWithDepth struct {
		tri   Triangle
		depth float64
	}
	items := make([]TriWithDepth, len(tris))
	for i, t := range tris {
		avgZ := (t.Verts[0].Z + t.Verts[1].Z + t.Verts[2].Z) / 3.0
		items[i] = TriWithDepth{t, avgZ}
	}
	sort.Slice(items, func(i, j int) bool { return items[i].depth < items[j].depth })
	for _, it := range items {
		c.DrawTriangle(it.tri, clr)
	}
}

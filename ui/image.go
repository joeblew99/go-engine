package ui

import (
	"image"
	"image/color"

	"github.com/walesey/go-engine/renderer"
	vmath "github.com/walesey/go-engine/vectormath"
)

type ImageElement struct {
	width, height int
	rotation      float64
	node          *renderer.Node
	img           image.Image
}

func (ie *ImageElement) Render(offset vmath.Vector2) vmath.Vector2 {
	width, height := ie.width, ie.height
	if width == 0 && height == 0 {
		width, height = ie.img.Bounds().Size().X, ie.img.Bounds().Size().Y
	} else if width == 0 {
		width = ie.img.Bounds().Size().X * height / ie.img.Bounds().Size().Y
	} else if height == 0 {
		height = ie.img.Bounds().Size().Y * width / ie.img.Bounds().Size().X
	}
	size := vmath.Vector2{float64(width), float64(height)}
	ie.node.SetScale(size.ToVector3())
	ie.node.SetTranslation(offset.ToVector3())
	return size
}

func (ie *ImageElement) Spatial() renderer.Spatial {
	return ie.node
}

func (ie *ImageElement) SetSize(width, height int) {
	ie.width, ie.height = width, height
}

func (ie *ImageElement) SetRotation(rotation float64) {
	ie.rotation = rotation
}

func NewImageElement(img image.Image) *ImageElement {
	box := renderer.CreateBoxWithOffset(1, 1, 0, 0)
	box.SetColor(color.RGBA{255, 0, 0, 255})
	// mat := renderer.CreateMaterial()
	// mat.Diffuse = img
	// mat.LightingMode = renderer.MODE_UNLIT
	// box.Material = mat
	node := renderer.CreateNode()
	node.Add(box)
	return &ImageElement{
		width:    img.Bounds().Size().X,
		height:   img.Bounds().Size().Y,
		rotation: 0,
		node:     node,
		img:      img,
	}
}

package ui

import (
	"fmt"
	"image/color"

	"github.com/walesey/go-engine/renderer"
	vmath "github.com/walesey/go-engine/vectormath"
)

type Window struct {
	node, elementNode, background *renderer.Node
	element                       Element
	size                          vmath.Vector2
}

func (w *Window) Draw(renderer renderer.Renderer) {
	w.node.Draw(renderer)
}

func (w *Window) Centre() vmath.Vector3 {
	return w.node.Centre()
}

func (w *Window) Optimize(geometry *renderer.Geometry, transform renderer.Transform) {
	w.node.Optimize(geometry, transform)
}

func (w *Window) SetScale(scale vmath.Vector3) {
	w.background.SetScale(scale)
	w.size = vmath.Vector2{scale.X, scale.Y}
	w.Render()
}

func (w *Window) SetTranslation(translation vmath.Vector3) {
	w.node.SetTranslation(translation)
}

func (w *Window) SetOrientation(orientation vmath.Quaternion) {
	w.node.SetOrientation(orientation)
}

func (w *Window) SetElement(element Element) {
	if w.element != nil {
		w.elementNode.Remove(w.element.Spatial())
	}
	w.element = element
	w.elementNode.Add(element.Spatial())
	w.Render()
}

func (w *Window) Render() {
	size := w.element.Render(vmath.Vector2{0, 0})
	scale := vmath.Vector3{1, 1, 1}
	if size.X > w.size.X {
		scale.X = w.size.X / size.X
	}
	if size.Y > w.size.Y {
		scale.Y = w.size.Y / size.Y
	}
	fmt.Println(scale)
	w.elementNode.SetScale(scale)
}

func NewWindow() *Window {
	node := renderer.CreateNode()
	elementNode := renderer.CreateNode()
	background := renderer.CreateNode()
	box := renderer.CreateBoxWithOffset(1, 1, 0, 0)
	box.SetColor(color.RGBA{255, 255, 255, 255})
	mat := renderer.CreateMaterial()
	mat.LightingMode = renderer.MODE_UNLIT
	box.Material = mat
	background.Add(box)
	node.Add(background)
	node.Add(elementNode)
	return &Window{
		node:        node,
		background:  background,
		elementNode: elementNode,
		size:        vmath.Vector2{1, 1},
	}
}

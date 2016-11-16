package engine

import (
	"fmt"
	"image/color"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/walesey/go-engine/renderer"
	"github.com/walesey/go-engine/ui"
)

// Engine is a wrapper for all the go-engine boilerblate code.
// It sets up a basic render / Update loop and provides a nice interface for writing games.
type Engine interface {
	Start(Init func())
	AddOrtho(spatial renderer.Spatial)
	AddSpatial(spatial renderer.Spatial)
	AddSpatialTransparent(spatial renderer.Spatial)
	RemoveSpatial(spatial renderer.Spatial, destroy bool)
	RemoveOrtho(spatial renderer.Spatial, destroy bool)
	AddUpdatable(updatable Updatable)
	RemoveUpdatable(updatable Updatable)
	Camera() *renderer.Camera
	Renderer() renderer.Renderer
	SetFpsCap(FpsCap float64)
	FPS() float64
	InitFpsDial()
	Update()
}

type EngineImpl struct {
	fpsMeter       *renderer.FPSMeter
	renderer       renderer.Renderer
	sceneGraph     *renderer.SceneGraph
	orthoNode      *renderer.Node
	camera         *renderer.Camera
	updatableStore *UpdatableStore
}

func (engine *EngineImpl) Start(Init func()) {
	if engine.renderer != nil {
		engine.renderer.SetInit(Init)
		engine.renderer.SetUpdate(engine.Update)
		engine.renderer.SetRender(engine.Render)
		engine.renderer.SetCamera(engine.camera)
		engine.renderer.Start()
	} else {
		Init()
		for {
			engine.Update()
		}
	}
}

func (engine *EngineImpl) Update() {
	dt := engine.fpsMeter.UpdateFPSMeter()
	engine.updatableStore.UpdateAll(dt)
}

func (engine *EngineImpl) Render() {
	engine.camera.Ortho = false
	engine.sceneGraph.RenderScene(engine.renderer, engine.camera.Translation)
	engine.camera.Ortho = true
	engine.orthoNode.Draw(engine.renderer, mgl32.Ident4())
}

func (engine *EngineImpl) AddOrtho(spatial renderer.Spatial) {
	if engine.orthoNode != nil {
		engine.orthoNode.Add(spatial)
	}
}

func (engine *EngineImpl) RemoveOrtho(spatial renderer.Spatial, destroy bool) {
	if engine.orthoNode != nil {
		engine.orthoNode.Remove(spatial, destroy)
	}
}

func (engine *EngineImpl) AddSpatial(spatial renderer.Spatial) {
	if engine.sceneGraph != nil {
		engine.sceneGraph.Add(spatial)
	}
}

func (engine *EngineImpl) AddSpatialTransparent(spatial renderer.Spatial) {
	if engine.sceneGraph != nil {
		engine.sceneGraph.AddTransparent(spatial)
	}
}

func (engine *EngineImpl) RemoveSpatial(spatial renderer.Spatial, destroy bool) {
	if engine.sceneGraph != nil {
		engine.sceneGraph.Remove(spatial, destroy)
		engine.sceneGraph.RemoveTransparent(spatial, destroy)
	}
}

func (engine *EngineImpl) AddUpdatable(updatable Updatable) {
	engine.updatableStore.Add(updatable)
}

func (engine *EngineImpl) RemoveUpdatable(updatable Updatable) {
	engine.updatableStore.Remove(updatable)
}

func (engine *EngineImpl) Camera() *renderer.Camera {
	return engine.camera
}

func (engine *EngineImpl) Renderer() renderer.Renderer {
	return engine.renderer
}

func (engine *EngineImpl) SetFpsCap(FpsCap float64) {
	engine.fpsMeter.FpsCap = FpsCap
}

func (engine *EngineImpl) FPS() float64 {
	return engine.fpsMeter.Value()
}

func (engine *EngineImpl) InitFpsDial() {
	window := ui.NewWindow()
	window.SetTranslation(mgl32.Vec3{10, 10, 1})
	window.SetScale(mgl32.Vec3{400, 0, 1})
	window.SetBackgroundColor(0, 0, 0, 0)

	container := ui.NewContainer()
	container.SetBackgroundColor(0, 0, 0, 0)
	window.SetElement(container)

	text := ui.NewTextElement("0", color.RGBA{255, 0, 0, 255}, 18, nil)
	container.AddChildren(text)
	engine.AddUpdatable(UpdatableFunc(func(dt float64) {
		fps := engine.FPS()
		text.SetText(fmt.Sprintf("%v", int(fps)))
		switch {
		case fps < 20:
			text.SetTextColor(color.RGBA{255, 0, 0, 255})
		case fps < 30:
			text.SetTextColor(color.RGBA{255, 90, 0, 255})
		case fps < 50:
			text.SetTextColor(color.RGBA{255, 255, 0, 255})
		default:
			text.SetTextColor(color.RGBA{0, 255, 0, 255})
		}
		text.ReRender()
	}))

	window.Render()
	engine.AddOrtho(window)
}

func NewEngine(r renderer.Renderer) Engine {
	fpsMeter := renderer.CreateFPSMeter(1.0)
	fpsMeter.FpsCap = 144

	sceneGraph := renderer.CreateSceneGraph()
	orthoNode := renderer.NewNode()
	updatableStore := NewUpdatableStore()
	camera := renderer.CreateCamera()

	return &EngineImpl{
		fpsMeter:       fpsMeter,
		sceneGraph:     sceneGraph,
		orthoNode:      orthoNode,
		updatableStore: updatableStore,
		renderer:       r,
		camera:         camera,
	}
}

func NewHeadlessEngine() Engine {
	updatableStore := NewUpdatableStore()
	fpsMeter := renderer.CreateFPSMeter(1.0)
	fpsMeter.FpsCap = 144

	return &EngineImpl{
		fpsMeter:       fpsMeter,
		updatableStore: updatableStore,
	}
}

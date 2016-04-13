package main

import (
	"os"
	"runtime"

	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/walesey/go-engine/controller"
	"github.com/walesey/go-engine/engine"
	"github.com/walesey/go-engine/renderer"
	"github.com/walesey/go-engine/ui"
	vmath "github.com/walesey/go-engine/vectormath"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func main() {
	glRenderer := renderer.NewOpenglRenderer("GoEngine", 1900, 1000)
	gameEngine := engine.NewEngine(glRenderer)

	// assetLib, err := assets.LoadAssetLibrary("TestAssets/ui.asset")
	// if err != nil {
	// 	panic(err)
	// }
	// alienwareImg := assetLib.GetImage("alienware")

	audiowideFont, _ := ui.LoadFont("TestAssets/Audiowide-Regular.ttf")
	htmlAssets := ui.NewHtmlAssets()
	htmlAssets.AddFont("default", audiowideFont)

	htmlAssets.AddCallback("activateTextField", func(element ui.Element, args ...interface{}) {
		container, ok := element.(*ui.Container)
		if ok {
			child := container.GetChild(0)
			text, ok := child.(*ui.TextElement)
			if ok {
				text.Activate()
			}
		}
	})
	htmlAssets.AddCallback("hover", func(element ui.Element, args ...interface{}) {
		container, ok := element.(*ui.Container)
		if ok {
			container.SetMargin(ui.Margin{10, 0, 0, 10})
		}
	})
	htmlAssets.AddCallback("unhover", func(element ui.Element, args ...interface{}) {
		container, ok := element.(*ui.Container)
		if ok {
			container.SetMargin(ui.NewMargin(0))
		}
	})

	gameEngine.Start(func() {

		window := ui.NewWindow()
		window.SetTranslation(vmath.Vector3{200, 200, 1})
		window.SetScale(vmath.Vector3{1200, 1, 1})
		window.SetBackgroundColor(90, 0, 255, 255)
		gameEngine.AddOrtho(window)

		tab := ui.NewContainer()
		tab.SetBackgroundColor(120, 120, 120, 255)
		tab.SetHeight(40)

		mainContainer := ui.NewContainer()
		window.SetElement(mainContainer)
		container := ui.NewContainer()
		mainContainer.AddChildren(tab, container)

		css, err := os.Open("TestAssets/test.css")
		if err != nil {
			panic(err)
		}
		defer css.Close()

		html, err := os.Open("TestAssets/test.html")
		if err != nil {
			panic(err)
		}
		defer html.Close()

		ui.LoadPage(container, html, css, htmlAssets)

		gameEngine.AddUpdatable(engine.UpdatableFunc(func(dt float64) {
			window.Render()
		}))

		//input/controller manager
		controllerManager := controller.NewControllerManager(glRenderer.Window)

		uiController := ui.NewUiController(window)
		controllerManager.AddController(uiController)

		//custom controller
		customController := controller.NewActionMap()
		controllerManager.AddController(customController)
		ui.ClickAndDragWindow(window, tab.Hitbox, customController)

		customController.BindMouseAction(func() {
			deactivateAllTextElements(mainContainer)
		}, glfw.MouseButton1, glfw.Press)
	})
}

func deactivateAllTextElements(container *ui.Container) {
	for i := 0; i < container.GetNbChildren(); i++ {
		child := container.GetChild(i)
		childContainer, ok := child.(*ui.Container)
		if ok {
			deactivateAllTextElements(childContainer)
		}
		text, ok := child.(*ui.TextElement)
		if ok {
			text.Deactivate()
		}
	}
}
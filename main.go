// Package main provides various examples of Fyne API capabilities.
package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"github.com/sirupsen/logrus"
	"idmiss/mqtt/cli/page"
	"idmiss/mqtt/cli/resource/images"
	"idmiss/mqtt/cli/runner"
	"idmiss/mqtt/cli/util/logger"
)

var a fyne.App
var Loading fyne.Window

func init() {
	a = app.NewWithID("com.idmiss.mqtt")
	driver := fyne.CurrentApp().Driver()
	if drv, ok := driver.(desktop.Driver); ok {
		Loading = drv.CreateSplashWindow()
		Loading.SetContent(widget.NewLabelWithStyle("正在加载……",
			fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))
		Loading.Hide()
	}
	runner.Runner()
}

const preferenceCurrentTutorial = "currentTutorial"

var topWindow fyne.Window

func shortcutFocused(s fyne.Shortcut, w fyne.Window) {
	if focused, ok := w.Canvas().Focused().(fyne.Shortcutable); ok {
		focused.TypedShortcut(s)
	}
}

func main() {

	a.Settings().SetTheme(page.MyTheme{})
	a.SetIcon(images.ResourceIdmisstxPng)
	w := a.NewWindow("MQTT模拟设备客户端 - go-mqtt-cli ")
	w.SetFixedSize(true)
	topWindow = w
	tutorial := page.Pages["welcome"].SetView(w)

	databaseItem := fyne.NewMenuItem("连接配置", func() {
		logger.Log.WithFields(logrus.Fields{"data": ""}).Info("连接配置")
		tutorial = page.Pages["database"].SetView(w)
		w.SetContent(tutorial)
	})

	helpMenu := fyne.NewMenu("帮助",
		fyne.NewMenuItem("查看文档", func() {
			tutorial = page.Pages["help"].SetView(w)
			w.SetContent(tutorial)
		}))

	welcome := fyne.NewMenuItem("首页", func() {
		tutorial = page.Pages["welcome"].SetView(w)
		w.SetContent(tutorial)
	})
	autocodeMenu := fyne.NewMenuItem("状态消息", func() {
		tutorial = page.Pages["status"].SetView(w)
		w.SetContent(tutorial)
		logger.Log.WithFields(logrus.Fields{"data": ""}).Info("状态消息")
	})
	mainMenu := fyne.NewMainMenu(
		// a quit item will be appended to our first menu
		fyne.NewMenu("文件", welcome),
		fyne.NewMenu("设置", databaseItem),
		fyne.NewMenu("查看", autocodeMenu),
		helpMenu,
	)

	w.SetMainMenu(mainMenu)
	w.SetMaster()

	w.SetContent(tutorial)

	w.Resize(fyne.NewSize(1280, 920))
	w.ShowAndRun()

	//os.Unsetenv("FYNE_FONT")
}

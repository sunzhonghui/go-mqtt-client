package page

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"idmiss/mqtt/cli/util/conf"
)

type PageDetail struct {
	Title, Intro string
	View         func(win fyne.Window) fyne.CanvasObject
}

var (
	Pages = map[string]*PageDetail{
		"welcome":  {Title: "MQTT模拟设备客户端 go-mqtt-client " + conf.Version, View: welcomeScreen},
		"database": {Title: "设置MQTT客户端连接", View: DatabaseScreen},
		"status":   {Title: "客户端状态，遥测消息，控制消息", View: StatusScreen},
		"help":     {Title: "模拟设备客户端使用方法", View: HelpScreen},
	}
)

func (d *PageDetail) SetView(win fyne.Window) *fyne.Container {
	return container.NewBorder(container.NewVBox(widget.NewLabelWithStyle(d.Title, fyne.TextAlignCenter, fyne.TextStyle{Bold: true}), widget.NewSeparator()), nil, nil, nil, d.View(win))
}

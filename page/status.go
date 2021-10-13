package page

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"idmiss/mqtt/cli/util/conf"
)

func StatusScreen(win fyne.Window) fyne.CanvasObject {

	return container.NewBorder(
		container.NewHBox(widget.NewLabel("客户端状态："), conf.StatusLabel, widget.NewLabel("    接口1状态："), conf.Open1, widget.NewLabel("    接口2状态："), conf.Open2), nil, nil, nil, container.NewHSplit(
			container.NewBorder(widget.NewLabel("已发送遥测数据"), nil, nil, nil, conf.LeftTextList.T),
			container.NewBorder(widget.NewLabel("接收的控制指令"), nil, nil, nil, conf.RightTextList.T),
		),
	)
}

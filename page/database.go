package page

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"idmiss/mqtt/cli/util/conf"
)

func DatabaseScreen(win fyne.Window) fyne.CanvasObject {
	urlName := widget.NewEntry()
	urlName.SetPlaceHolder("测试连接")
	urlName.Validator = validation.NewRegexp(`^.{1,50}$`, "请输入连接名称")
	urlName.SetText(conf.Database.UrlName)

	clientId := widget.NewEntry()
	clientId.SetPlaceHolder("yy-0001")
	clientId.Validator = validation.NewRegexp(`^.{1,50}$`, "输入正确得客户端id")
	clientId.SetText(conf.Database.ClientId)

	serverUrl := widget.NewEntry()
	serverUrl.Validator = validation.NewRegexp(`^.{1,50}$`, "请输入连接地址")
	serverUrl.SetPlaceHolder("ws://mqtt.youyanghealth.com:8083/mqtt")
	serverUrl.SetText(conf.Database.ServerUrl)

	userName := widget.NewEntry()
	userName.SetPlaceHolder("")
	userName.SetText(conf.Database.UserName)

	password := widget.NewPasswordEntry()
	password.SetPlaceHolder("")
	password.SetText(conf.Database.Password)

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "连接名：", Widget: urlName, HintText: "输入连接名称"},
			{Text: "客户端ID：", Widget: clientId, HintText: "输入客户端ID"},
			{Text: "服务器地址：", Widget: serverUrl, HintText: "输入服务器地址"},
			{Text: "用户名：", Widget: userName, HintText: "输入用户名"},
			{Text: "密码：", Widget: password, HintText: "输入密码"},
		},
		OnCancel: func() {
			//mysqlUrl := userName.Text + ":" + password.Text + "@(" + ip.Text + ":" + port.Text + ")/" + databaseName.Text + "?charset=utf8&parseTime=True&loc=Local"
			//logger.Log.WithFields(logrus.Fields{"data": mysqlUrl}).Info("数据库连接地址")
			//db, err := gorm.Open(mysql.Open(mysqlUrl), &gorm.Config{})
			//if err != nil {
			//	dialog.ShowError(errors.New("连接失败"), win)
			//} else {
			//	conf.DB = db
			//	dialog.ShowInformation("提示", "连接成功", win)
			//}
			conf.Database.DisCon()
			dialog.ShowInformation("提示", "断开链接", win)
		},
		OnSubmit: func() {
			conf.Database.UrlName = urlName.Text
			conf.Database.ClientId = clientId.Text
			conf.Database.ServerUrl = serverUrl.Text
			conf.Database.UserName = userName.Text
			conf.Database.Password = password.Text
			conf.Database.Save()
			if err := conf.Database.GetDB(); err == nil {
				dialog.ShowInformation("提示", "连接成功", win)
				//tutorial := Pages["status"].SetView(win)
				//win.SetContent(tutorial)
			} else {
				dialog.ShowInformation("提示", "连接失败", win)
			}

		},
	}
	form.CancelText = "断开"
	form.SubmitText = "连接"

	return form
}

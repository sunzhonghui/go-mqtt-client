package font

import (
	"github.com/sirupsen/logrus"
	"idmiss/mqtt/cli/util/common"
	"idmiss/mqtt/cli/util/logger"
	"os"
)

const (
	zhFontUrl = "http://f.idmiss.com/file/font/WenQuanWeiMiHei.ttf"
)

func GetZhFont() {
	logger.Log.WithFields(logrus.Fields{"data": ""}).Info("检测字体是否存在...")
	if _, err := common.DownloadFileForUrl(zhFontUrl, "resource/font/", "WenQuanWeiMiHei.ttf"); err != nil {
		logger.Log.WithFields(logrus.Fields{"err": err}).Error("字体下载失败...")

	}
	logger.Log.Info("开始配置字体...")
	os.Setenv("FYNE_FONT", "resource/font/WenQuanWeiMiHei.ttf")
	logger.Log.Info("配置字体成功")
}

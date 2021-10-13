package runner

import (
	"idmiss/mqtt/cli/util/conf"
	"idmiss/mqtt/cli/util/logger"
)

func Runner() {
	//font.GetZhFont()
	//font.GetTemp()
	conf.Init()
	logger.InitLog()
}

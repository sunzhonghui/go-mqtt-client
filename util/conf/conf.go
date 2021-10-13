package conf

import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2/widget"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"idmiss/mqtt/cli/util/logger"
	mqtt2 "idmiss/mqtt/cli/util/mqtt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var Version = "v1.1.1"

type DatabaseConf struct {
	UrlName   string `json:"urlName"`  //连接名称
	ClientId  string `json:"clientId"` //client id
	ServerUrl string `json:"serverUrl"`
	UserName  string `json:"userName"`
	Password  string `json:"password"`
	//WillTopic string `json:"willTopic"` //遗言主题
	//WillData string `json:"willData"` //遗言内容
}

type MqttTempData struct {
	Bat   float64           `json:"Bat"`
	Cmd   int64             `json:"Cmd"`
	DTime string            `json:"D_time"`
	Imei  string            `json:"Imei"`
	Lat   string            `json:"Lat"`
	Lon   string            `json:"Lon"`
	Por   string            `json:"Por"`
	Ver   string            `json:"Ver"`
	Data  []MqttTempDataArr `json:"Data"`
}
type MqttOpenStatus struct {
	Solenoid1 int64 `json:"Solenoid1"`
	Solenoid2 int64 `json:"Solenoid2"`
}
type MqttTempDataArr struct {
	Humid       float64 `json:"Humid"`
	LiquidLevel int64   `json:"LiquidLevel"`
	O2          int64   `json:"O2"`
	O2Pressure  float64 `json:"O2_pressure"`
	RTDTemp     float64 `json:"RTD_temp"`
	RunStatus   int64   `json:"RunStatus"`
	MqttOpenStatus
	Temp   float64 `json:"Temp"`
	Weight int64   `json:"Weight"`
}

type MqttTempCall struct { // 指令接受
	Cmd   int64              `json:"Cmd"`
	DTime string             `json:"D_time"`
	Data  []MqttTempCallData `json:"Data"`
	Imei  string             `json:"Imei"`
}

type MqttTempCallData struct {
	Solenoid1 string `json:"Solenoid1"`
	Solenoid2 string `json:"Solenoid2"`
}

type MqttTempCallBack struct { // 指令回复
	Cmd   int64            `json:"Cmd"`
	DTime string           `json:"D_time"`
	Data  []MqttOpenStatus `json:"Data"`
	Imei  string           `json:"Imei"`
}

type TestTempList struct {
	L []struct {
		Msg  string `json:"msg"`
		Date string `json:"date"`
	} `json:"l"`
	T *widget.Entry `json:"t"`
}

var Database = &DatabaseConf{}
var MqttOpen = &MqttOpenStatus{2, 2} // 设备开启状态
var MqCli mqtt.Client
var cro = cron.New(cron.WithSeconds()) //精确到秒
var croId cron.EntryID
var croId2 cron.EntryID
var StatusLabel = widget.NewLabel("离线") //设备状态 label
var Open1 = widget.NewLabel("关闭")       //设备状态 label
var Open2 = widget.NewLabel("关闭")       //设备状态 label
func SetOpen1() {
	if MqttOpen.Solenoid1 == 2 {
		Open1.SetText("关闭")
	} else {
		Open1.SetText("打开")
	}
}
func SetOpen2() {
	if MqttOpen.Solenoid2 == 2 {
		Open2.SetText("关闭")
	} else {
		Open2.SetText("打开")
	}
}

//var LeftText = widget.NewMultiLineEntry()
//LeftText.Wrapping = fyne.TextWrapWord
//var RightText = widget.NewMultiLineEntry()
//LeftText.Wrapping = fyne.TextWrapWord
var LeftTextList = TestTempList{[]struct {
	Msg  string `json:"msg"`
	Date string `json:"date"`
}{}, widget.NewMultiLineEntry()}
var RightTextList = TestTempList{[]struct {
	Msg  string `json:"msg"`
	Date string `json:"date"`
}{}, widget.NewMultiLineEntry()}

func Init() {

	//viper.SetConfigType("yaml") // 如果配置文件的名称中没有扩展名，则需要配置此项
	viper.AddConfigPath("resource") // 设置读取路径：就是在此路径下搜索配置文件。
	//viper.AddConfigPath("$HOME/.appname")  // 多次调用以添加多个搜索路径
	viper.SetConfigName("conf") // 设置被读取文件的全名，包括扩展名。
	//viper.SetConfigName("server") // 设置被读取文件的名字： 这个方法 和 SetConfigFile实际上仅使用一个就够了
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig() // 读取配置文件： 这一步将配置文件变成了 Go语言的配置文件对象包含了 map，string 等对象。
	if err != nil {
		logger.Log.WithFields(logrus.Fields{"err": err}).Warn("未找到配置文件，创建配置文件")
		viper.SetDefault("database.uname", "测试连接")
		viper.SetDefault("database.cid", "yy-0001")
		viper.SetDefault("database.surl", "ws://mqtt.youyangjiankang.com:8083/mqtt")
		viper.SetDefault("database.username", "")
		viper.SetDefault("database.password", "")
		//viper.SetDefault("database.will.topic", "$sta/rep/{clientId}")
		//viper.SetDefault("database.will.data", "offline")

		//项目默认配置
		viper.WriteConfigAs("resource/conf.yaml")
	}

	ResetData()

	// 控制台输出： map[first:panda last:8z] 99 panda [Coding Movie Swimming]

}

func ResetData() {
	Database.UrlName = viper.GetString("database.uname")
	Database.ClientId = viper.GetString("database.cid")
	Database.ServerUrl = viper.GetString("database.surl")
	Database.UserName = viper.GetString("database.username")
	Database.Password = viper.GetString("database.password")

}

func (d *DatabaseConf) Save() {
	logger.Log.WithFields(logrus.Fields{"data": d}).Info("保存数据库配置到文件")

	viper.Set("database.uname", d.UrlName)
	viper.Set("database.cid", d.ClientId)
	viper.Set("database.username", d.UserName)
	viper.Set("database.password", d.Password)
	viper.Set("database.surl", d.ServerUrl)
	err := viper.WriteConfigAs("resource/conf.yaml")
	if err != nil {
		logger.Log.WithFields(logrus.Fields{"data": d, "err": err}).Error("保存数据库配置到文件失败")
	}
}

func (d *DatabaseConf) GetDB() error {

	if MqCli != nil && MqCli.IsConnected() {
		logger.Log.Info("已经连接了。")
		return nil
	}
	opts := mqtt.NewClientOptions().AddBroker(Database.ServerUrl).SetClientID(Database.ClientId)

	opts.SetKeepAlive(60 * time.Second)
	// 设置消息回调处理函数
	opts.SetDefaultPublishHandler(mqtt2.MqttBackFun)
	opts.SetPingTimeout(1 * time.Second)
	//opts.SetWill("$sta/rep/"+d.ClientId,"offline",2,true)
	opts.SetOnConnectHandler(func(client mqtt.Client) {
		logger.Log.WithFields(logrus.Fields{"data": d}).Info("连接Mqtt服务器成功")
		StatusLabel.SetText("在线")
	})
	opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		StatusLabel.SetText("离线")
		logger.Log.WithFields(logrus.Fields{"err": err, "clientId": d.ClientId}).Info("Mqtt 连接断开")
		if token := MqCli.Connect(); token.Wait() && token.Error() != nil {
			logger.Log.WithFields(logrus.Fields{"err": token.Error()}).Error("断开重连，连接失败")
		}
	})

	MqCli = mqtt.NewClient(opts)
	if token := MqCli.Connect(); token.Wait() && token.Error() != nil {
		logger.Log.WithFields(logrus.Fields{"err": token.Error()}).Error("连接mqtt服务器失败")
		return token.Error()
	}
	//else {
	//	token := MqCli.Publish("$sta/rep/"+d.ClientId,2,true,"online")
	//	if token.Wait()&& token.Error()!=nil{
	//		logger.Log.WithFields(logrus.Fields{"err":token.Error(),"topic":"$sta/rep/"+d.ClientId}).Error("发送状态主题失败")
	//		return token.Error()
	//	}
	//
	//}

	// 订阅主题
	if token := MqCli.Subscribe("$creq/"+d.ClientId+"/cmd", 2, func(client mqtt.Client, message mqtt.Message) {
		logger.Log.WithFields(logrus.Fields{"msg": string(message.Payload()), "topic": "$creq/" + d.ClientId + "/cmd"}).Info("控制主题收到消息")

		payload := message.Payload()[1 : len(message.Payload())-1]

		call := &MqttTempCall{}
		if err := json.Unmarshal(payload, call); err != nil {
			logger.Log.WithFields(logrus.Fields{"err": err, "topic": "$creq/" + d.ClientId + "/cmd"}).Error("控制主题消息格式错误", string(payload))
		}
		//RightTextList = append(LeftTextList,TestTemp{string(message.Payload()),time.Now().Format("15:04:05")})
		RightTextList.set(string(message.Payload()))
		if call.Cmd == 3 && len(call.Data) > 0 {
			if strings.Index(call.Data[0].Solenoid1, "1") == 0 {
				MqttOpen.Solenoid1 = 1
			}
			if strings.Index(call.Data[0].Solenoid2, "1") == 0 {
				MqttOpen.Solenoid2 = 1
			}
			if strings.Index(call.Data[0].Solenoid1, "2") == 0 {
				MqttOpen.Solenoid1 = 2
			}
			if strings.Index(call.Data[0].Solenoid2, "2") == 0 {
				MqttOpen.Solenoid2 = 2
			}
			SetOpen1()
			SetOpen2()
			d.sendCall3()
		}
		if call.Cmd == 4 && len(call.Data) > 0 {
			d.SendMs()
		}
		//d.SendMs()
	}); token.Wait() && token.Error() != nil {
		logger.Log.WithFields(logrus.Fields{"err": token.Error(), "topic": "$creq/" + d.ClientId + "/cmd"}).Error("订阅主题失败")
		return token.Error()
	}

	if croId > 0 {
		cro.Remove(croId)
	}
	spec1 := "*/15 * * * * ?" //cron表达式，每15s一次
	croId, _ = cro.AddFunc(spec1, func() {
		d.SendMs()
	})

	if croId2 > 0 {
		cro.Remove(croId2)
	}
	spec2 := "*/30 * * * * ?" //cron表达式，每30s一次
	croId, _ = cro.AddFunc(spec2, func() {
		d.sendHb()
	})
	cro.Start()
	return nil
}

//断开连接
func (d DatabaseConf) DisCon() {
	if cro != nil {
		cro.Stop()
	}
	if MqCli != nil && MqCli.IsConnected() {
		logger.Log.Info("手动断开连接")
		//MqCli.Publish("$sta/rep/"+d.ClientId,2,true,"offline")
		MqCli.Disconnect(100)
		StatusLabel.SetText("离线")
	}
}

func (d DatabaseConf) SendMs() {
	// 15s 上报一次
	if MqCli == nil {
		logger.Log.WithFields(logrus.Fields{"topic": "/$dp"}).Error("发送遥测主题失败,mqtt 对象空")
		return
	}
	data := MqttTempData{}
	data.Ver = "V1.0"
	data.Por = d.ClientId
	data.Imei = d.ClientId
	rand.Seed(time.Now().UnixNano())
	data.Lon = Round2Str(100*rand.Float64(), 4)
	rand.Seed(time.Now().UnixNano() + 10)
	data.Lat = Round2Str(100*rand.Float64(), 4)
	data.Cmd = 2
	rand.Seed(time.Now().UnixNano() + 20)
	data.Bat = Round2(100*rand.Float64(), 2)
	data.DTime = time.Now().Format("2006-01-02 15:04:05")

	arr := MqttTempDataArr{}
	arr.Solenoid1 = MqttOpen.Solenoid1
	arr.Solenoid2 = MqttOpen.Solenoid2

	rand.Seed(time.Now().UnixNano() + 30)
	arr.LiquidLevel = rand.Int63n(100)
	arr.O2Pressure = Round2(1*rand.Float64(), 1)
	arr.Temp = Round2(100*rand.Float64(), 2)
	arr.Humid = Round2(100*rand.Float64(), 2)
	arr.RTDTemp = Round2(100*rand.Float64(), 2)
	arr.Weight = rand.Int63n(100)
	arr.O2 = rand.Int63n(100)
	arr.RunStatus = 0
	data.Data = []MqttTempDataArr{arr}
	jdata, _ := json.Marshal(data)
	sendDataStr := "$" + string(jdata) + "#"
	token := MqCli.Publish("/$dp", 2, false, sendDataStr)
	if token.Wait() && token.Error() != nil {
		logger.Log.WithFields(logrus.Fields{"err": token.Error(), "topic": "/$dp", "data": sendDataStr}).Error("发送遥测数据失败")
		return
	}
	logger.Log.WithFields(logrus.Fields{"topic": "/$dp", "data": sendDataStr}).Info("发送遥测数据成功")
	//LeftTextList = append(LeftTextList,TestTemp{sendDataStr,time.Now().Format("15:04:05")})
	LeftTextList.set(sendDataStr)
}

type HeartBeat struct {
	Imei string `json:"Imei"`
	Cmd  int64  `json:"Cmd"`
}

func (d DatabaseConf) sendHb() {
	data := HeartBeat{Imei: d.ClientId, Cmd: 1}
	jdata, _ := json.Marshal(data)
	sendDataStr := "$" + string(jdata) + "#"
	token := MqCli.Publish("/$dp", 2, false, sendDataStr)
	if token.Wait() && token.Error() != nil {
		logger.Log.WithFields(logrus.Fields{"err": token.Error(), "topic": "/$dp", "data": sendDataStr}).Error("发送心跳失败")
		return
	}
	logger.Log.WithFields(logrus.Fields{"topic": "/$dp", "data": sendDataStr}).Info("发送心跳成功")
	//LeftTextList = append(LeftTextList,TestTemp{sendDataStr,time.Now().Format("15:04:05")})
	LeftTextList.set(sendDataStr)

}
func (d DatabaseConf) sendCall3() {
	data := MqttTempCallBack{Imei: d.ClientId, Cmd: 3, DTime: time.Now().Format("2006-01-02 15:04:05")}
	data.Data = []MqttOpenStatus{*MqttOpen}

	jdata, _ := json.Marshal(data)
	sendDataStr := "$" + string(jdata) + "#"
	token := MqCli.Publish("/$dp", 2, false, sendDataStr)
	if token.Wait() && token.Error() != nil {
		logger.Log.WithFields(logrus.Fields{"err": token.Error(), "topic": "/$dp", "data": sendDataStr}).Error("发送控制3回执失败")
		return
	}
	logger.Log.WithFields(logrus.Fields{"topic": "/$dp", "data": sendDataStr}).Info("发送控制3回执成功")
	//LeftTextList = append(LeftTextList,TestTemp{sendDataStr,time.Now().Format("15:04:05")})
	LeftTextList.set(sendDataStr)

}
func (t *TestTempList) set(d string) {

	t.L = append([]struct {
		Msg  string `json:"msg"`
		Date string `json:"date"`
	}{{Msg: d, Date: time.Now().Format("15:04:05")}}, t.L...)
	if len(t.L) > 15 {
		t.L = t.L[:15]
	}
	s := ""
	for _, v := range t.L {
		s += v.Date + "\r\n" + v.Msg + "\r\n\r\n"
	}
	t.T.SetText(s)
}
func Round2(f float64, n int) float64 {
	floatStr := fmt.Sprintf("%."+strconv.Itoa(n)+"f", f)
	inst, _ := strconv.ParseFloat(floatStr, 64)
	return inst
}
func Round2Str(f float64, n int) string {
	floatStr := fmt.Sprintf("%."+strconv.Itoa(n)+"f", f)
	return floatStr
}

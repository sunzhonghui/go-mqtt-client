package page

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func HelpScreen(win fyne.Window) fyne.CanvasObject {

	helpEnt := widget.NewMultiLineEntry()
	helpEnt.SetText(`
		点击设置 - 连接配置

		输入客户端ID(相当于设备编号) 跟 服务器地址即可连接mqtt服务器
		
		设备客户端消息体：

		{
			"Por":"xxxx",
			"Ver":"V1.0",
			"Imei":"xxxx",
			"Lon":"100.0000",
			"Lat":"90.0000",
			"Bat":11.03,
			"D_time":"2021-09-18  16:01:07",
			"Cmd":2,
			"Data":[
				{
					"LiquidLevel":80,
					"O2_pressure":0.1,
					"Temp":24.3,
					"Humid":23.4,
					"RTD_temp":22.63,
					"Weight":0,
					"O2":0,
					"Solenoid1":2,
					"Solenoid2":2,
					"RunStatus":2
				}
			]
		}

		设定操作：

		1.该设备会对 /$dp  的定时发布如下消息，会对$creq/{clientId}/cmd  订阅消息

		2.每15秒，该设备发布设备状态：
			
			${"Por":"guwei-test3","Ver":"V1.0","Imei":"guwei-test3","Lon":"100.0000","Lat":"90.0000","Bat":11.03,"D_time":"2021-09-18  16:01:07","Cmd":2,
			"Data":[{"LiquidLevel":80.00,"O2_pressure":0.1,"Temp":24.3,"Humid":23.4,"RTD_temp":22.63,"Weight":0,"O2":0.0,"Solenoid1":2,"Solenoid2":2,"RunStatus":2}]}#

		3.每30秒，该设备发布一次心跳，消息体示例如下：
			${"Imei":"guwei-test3","Cmd":1}#

		4.指令下发与回执
			注：  1开头说明要开启 2说明要关闭 0则不处理

			①系统下发接口1开指令如下：
				${"Imei":"guwei-test3","D_time":"2021-10-4 11:11:05","Cmd":"3","Data":[{"Solenoid1": "1,30","Solenoid2": "0,30"}]}#
			②设备接收接口1开，回执指令如下：
				${"Imei":"guwei-test3","D_time":"2021-10-04 11:11:17","Cmd":3,"Data":[{"Solenoid1":1,"Solenoid2":0}]}#
			③系统下发接口1关指令如下：
				${"Imei":"guwei-test3","D_time":"2021-10-4 11:12:59","Cmd":"3","Data":[{"Solenoid1": "2,30","Solenoid2": "0,30"}]}#
			④设备接收接口1关指令，回执如下：
				${"Imei":"guwei-test3","D_time":"2021-10-04 11:13:12","Cmd":3,"Data":[{"Solenoid1":2,"Solenoid2":0}]}#

		5.系统下发设备查询指令如下：
			${"imei":"guwei-test3","cmd":4}#
			
		  设备接收到查询指令，回执如下：
			${"Por":"guwei-test3","Ver":"V1.0","Imei":"guwei-test3","Lon":"100.0000","Lat":"90.0000","Bat":11.03,"D_time":"2021-09-18  16:01:07","Cmd":2,
			"Data":[{"LiquidLevel":80.00,"O2_pressure":0.1,"Temp":24.3,"Humid":23.4,"RTD_temp":22.63,"Weight":0,"O2":0.0,"Solenoid1":2,"Solenoid2":2,"RunStatus":2}]}#
		点击 - 查看 - 状态消息 

		可以查看当前设备客户端 的状态

		查看设备发送的遥测数据 跟 设备接收到的 控制指令

		字段说明
| **字段**                		| **字段类型**  	| **是否为空**  	| **字段说明**                                                 	| **注释**                                                     	|
| ---------------------------- 	| ------------ 	| ------------	| ------------------------------------------------------------	| ------------------------------------------------------------ 	|
| Por                          	| varchar(25)  	| 否           	| 设备名称                                                   	| COMB000001：项目主板名称，数字根据项目规模确定，即每个设备名称唯一 |
| Ver                          	| varchar(25)  	| 否           	| 版本号                                                       	| 协议版本号                                                   	|
| Imei                         	| varchar(25)  	| 否           	| 设备编号                                                   	|                                                              	|
| lon                          	| float        	| 否           	| 经度                                                      		|                                                              	|
| lat                          	| float        	| 否           	| 纬度                                                      		|                                                              	|
| bat                          	| float        	| 否           	| 电池电量                                                   	|                                                              	|
| d_time                       	| timestamp    	| 否           	| 数据采集时间                                                	| 设备所在时间                                                 	|
| cmd                          	| varchar(10)  	| 否           	| 上传命令   1：心跳   2：数据上传  3：下发控制电机  4：下发查询状态	| 数据上传时间15秒，心跳时间30秒                               	|
| data                         	| varchar(500) 	| 否           	| 数据                                                          	| 一组包含各传感器和设备状态的数据组                           		|
| **data****数组包含以下内容**	|              	| 否         	|                                                              	|                                                              	|
| liquidLevel                  	| float        	| 否           	| 液位                                                          	|                                                              	|
| O2_pressure                  	| float        	| 否           	| 气压                                                          	|                                                              	|
| temp                         	| float        	| 否           	| 环境温度                                                   	|                                                              	|
| humid                        	| float        	| 否           	| 环境湿度                                                   	|                                                              	|
| RTD_temp                     	| float        	| 否           	| 信号强度                                                   	|                                                              	|
| Weight                       	| float        	| 否           	| 重量                                                      		|                                                              	|
| O2                           	| float        	| 否           	| 氧气浓度                                                     	| 环境氧气浓度                                                 	|
| solenoid1                    	| bool         	| 否           	| 电磁阀1                                                   		| 打开状态：1   关闭状态：2                                    	|
| solenoid2                    	| bool         	| 否           	| 电磁阀2                                                   		| 打开状态：1   关闭状态：2                                  		|
| RunStatus                    	| int(11)      	| 否           	| 设备运行状态                                                	| 正常状态：0      关注状态：1   补能状态：2    故障状态：3    		|
`)
	helpEnt.Disable()
	return container.NewBorder(nil, nil, nil, nil, helpEnt)
}

module github.com/company/mcutie

go 1.15

replace github.com/company/mcutie/test v1.0.0 => /Users/dan/Nextcloud/programming/golang/mcutie/test

replace github.com/company/mcutie/getstats v1.0.0 => /Users/dan/Nextcloud/programming/golang/mcutie/getstats

require (
	github.com/company/mcutie/getstats v1.0.0
	github.com/deckarep/gosx-notifier v0.0.0-20180201035817-e127226297fb
	github.com/distatus/battery v0.10.0 // indirect
	github.com/eclipse/paho.mqtt.golang v1.3.0
	github.com/gen2brain/beeep v0.0.0-20200526185328-e9c15c258e28
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/glendc/go-external-ip v0.0.0-20200601212049-c872357d968e // indirect
	github.com/go-ini/ini v1.62.0 // indirect
	github.com/shirou/gopsutil v3.20.11+incompatible // indirect
	github.com/shirou/gopsutil/v3 v3.20.11 // indirect
	github.com/sirupsen/logrus v1.7.0
	github.com/smartystreets/goconvey v1.6.4 // indirect
	github.com/urfave/cli v1.22.5 // indirect
	github.com/zpatrick/go-config v0.0.0-20191118215128-80ba6b3e54f6
	gopkg.in/ini.v1 v1.62.0 // indirect
)

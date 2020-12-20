package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/company/mcutie/getstats"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/zpatrick/go-config"

	gosxnotifier "github.com/deckarep/gosx-notifier"
	"github.com/gen2brain/beeep"
)

// We use gosx-notifier instead of beep on OSX. This looks better than "beeep" which
// always shows Osascript in Subject and links to the applescript editor

// Command json recieved via MQTT
type Command struct {
	Prog string
	Arg1 string
	Arg2 string
	Arg3 string
}

// HostNameSafe - this machines hostname with any suffic removed
var HostNameSafe string

// This is where we deal with messages published to our subscribed topics
var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.Info("messagePubHandler: Received message")
	log.Info("messagePubHandler: From topic: ", msg.Topic())
	decodeme := msg.Payload()
	var jcmd Command
	json.Unmarshal([]byte(decodeme), &jcmd)
	switch os := runtime.GOOS; os {
	case "darwin":
		log.Info("messagePubHandler: System: OS X")
		if jcmd.Prog == "notify" {
			log.Info("messagePubHandler: notify")
			osxnotify(jcmd.Arg1, jcmd.Arg2)
		}
		if jcmd.Prog == "execute" {
			log.Info("messagePubHandler: Running blind system command")
			execute(jcmd.Arg1, jcmd.Arg2, jcmd.Arg3)
		}

	case "linux":
		log.Info("messagePubHandler: System: Linux")
	case "windows":
		log.Info("messagePubHandler: Windows")
		if jcmd.Prog == "notify" {
			windowsNotify(jcmd.Arg1, jcmd.Arg2)
		}
	default:
		log.Info("messagePubHandler: case default reached - I don't know what OS I'm on!")
	}
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	log.Info("Connected to MQTT server")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.Warn("Connect lost: ", err)
}

func windowsNotify(title string, subject string) {
	err := beeep.Notify(title, subject, "gopher.png")
	if err != nil {
		log.Error(err)
	}

}

func osxnotify(title string, subject string) {
	log.Info("Title: ", title, " Subject: ", subject)
	note := gosxnotifier.NewNotification(subject)
	note.Title = title

	note.Group = "com.unique.yourapp.identifier"

	note.Sender = "com.apple.Safari"

	err := note.Push()
	if err != nil {
		log.Error("osxnotify error")
	}
}

func execute(arg1 string, arg2 string, arg3 string) {

	log.Info("function: execute")

	out, err := exec.Command(arg1, arg2, arg3).Output()
	if err != nil {
		log.Error(err)
	}
	log.Info("Command Successfully Executed")
	output := string(out[:])
	log.Info(output)

}

func publishHomeAssistantAutoConfigData(client mqtt.Client, myGroup string, mySensor string, unitOfMasurement string, iconChoice string) {
	rootTopic := "homeassistant/sensor/mcutie/"
	// myGroup := "system"
	// mySensor := "cpu"
	// unitOfMasurement := "%"
	// iconChoice := "mdi:speedometer"
	log.Info("Fuction: publishHomeAssistantAutoConfigData hostNameSafe = ", HostNameSafe)
	myTopic := rootTopic + HostNameSafe + "_" + myGroup + "_" + mySensor + "/config"
	availabilityTopic := "mcutie/" + HostNameSafe + "/lwt"
	stateTopic := "mcutie/" + HostNameSafe + "/stats/" + myGroup + "/" + mySensor
	myName := HostNameSafe + " " + mySensor
	myUniqueID := HostNameSafe + "_" + mySensor

	jsonData :=
		"{ \"unit_of_measurement\": \"" + unitOfMasurement + "\"" + "," +
			"\"icon\": \"" + iconChoice + "\"" + "," +
			"\"availability_topic\": \"" + availabilityTopic + "\"" + "," +
			"\"state_topic\": \"" + stateTopic + "\"" + "," +
			"\"name\": \"" + myName + "\"" + "," +
			"\"unique_id\": \"" + myUniqueID + "\"" + "," +
			"\"payload_available\": \"" + "ON" + "\"" + "," +
			"\"payload_not_available\": \"" + "OFF" + "\"" + "," +
			"\"device\": {" +
			"\"identifiers\": [" +
			"\"" + HostNameSafe + "_" + mySensor + "\"" +
			"]," +
			"\"name\": \"" + HostNameSafe + " " + mySensor + "\"" + "," +
			"\"model\": \"" + "v1.0" + "\"" + "," +
			"\"manufacturer\": \"" + "MCutie" + "\"" + "}}"

	token := client.Publish(myTopic, 0, false, jsonData)

	token.Wait()
}

func publishStats(client mqtt.Client, updateInterval int) {
	for {

		myTopic := "mcutie/" + HostNameSafe + "/lwt"
		token := client.Publish(myTopic, 0, false, "ON")

		myTopic = "mcutie/" + HostNameSafe + "/stats/system/hostname"
		token = client.Publish(myTopic, 0, false, HostNameSafe)

		myTopic = "mcutie/" + HostNameSafe + "/stats/system/cpu"
		token = client.Publish(myTopic, 0, false, getstats.CPUUsage())

		myTopic = "mcutie/" + HostNameSafe + "/stats/system/user"
		token = client.Publish(myTopic, 0, false, getstats.CurrentUser())

		runtimeOS := runtime.GOOS
		myTopic = "mcutie/" + HostNameSafe + "/stats/system/os"
		token = client.Publish(myTopic, 0, false, runtimeOS)

		myTopic = "mcutie/" + HostNameSafe + "/stats/system/uptime"
		token = client.Publish(myTopic, 0, false, getstats.UpTime())

		myTopic = "mcutie/" + HostNameSafe + "/stats/power/battery"
		token = client.Publish(myTopic, 0, false, getstats.BatteryLevel())

		myTopic = "mcutie/" + HostNameSafe + "/stats/net/iplocal"
		token = client.Publish(myTopic, 0, false, getstats.LocalIP())

		myTopic = "mcutie/" + HostNameSafe + "/stats/net/ipwan"
		token = client.Publish(myTopic, 0, false, getstats.ExternalIP())

		myTopic = "mcutie/" + HostNameSafe + "/stats/memory/mem_total"
		token = client.Publish(myTopic, 0, false, getstats.MemTotal())

		myTopic = "mcutie/" + HostNameSafe + "/stats/memory/mem_used"
		token = client.Publish(myTopic, 0, false, getstats.MemUsed())

		myTopic = "mcutie/" + HostNameSafe + "/stats/memory/mem_free"
		token = client.Publish(myTopic, 0, false, getstats.MemFree())

		myTopic = "mcutie/" + HostNameSafe + "/stats/memory/mem_percent"
		token = client.Publish(myTopic, 0, false, getstats.MemUsedPercent())

		myTopic = "mcutie/" + HostNameSafe + "/stats/disk/disk_total"
		token = client.Publish(myTopic, 0, false, getstats.DiskTotal())

		myTopic = "mcutie/" + HostNameSafe + "/stats/disk/disk_used"
		token = client.Publish(myTopic, 0, false, getstats.DiskUsed())

		myTopic = "mcutie/" + HostNameSafe + "/stats/disk/disk_free"
		token = client.Publish(myTopic, 0, false, getstats.DiskFree())

		token.Wait()

		time.Sleep(time.Duration(updateInterval) * time.Second)
	}
}

func sub(client mqtt.Client) {
	commandTopic := "mcutie/" + HostNameSafe + "/command"
	token := client.Subscribe(commandTopic, 1, nil)
	token.Wait()
	log.Info("Subscribed to topic: ", commandTopic)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func initConfig() *config.Config {
	yamlFile := config.NewYAMLFile("config.yaml")
	return config.NewConfig([]config.Provider{yamlFile})
}

// *******************************************
// main
// *******************************************

func main() {

	// open a file
	f, err := os.OpenFile("mcutie.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}

	// don't forget to close it
	defer f.Close()

	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.TextFormatter{})

	// Output to stderr instead of stdout, could also be a file.
	log.SetOutput(f)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
	log.Info("***")
	log.Info("*** Start program execution ***")
	log.Info("***")
	log.Info("Cleaning up hostname to remove suffix (if present)")
	varTemp := getstats.HostName()
	varTemp2 := strings.Split(varTemp, ".")
	HostNameSafe = varTemp2[0]
	HostNameSafe = strings.ToUpper(HostNameSafe)

	log.Info("Safe hostname = ", HostNameSafe)

	log.Info("read data from config.yaml")
	conf := initConfig()

	connectURL, err := conf.String("url")
	if err != nil {
		log.Fatal(err)
	}
	connectUsername, err := conf.String("username")
	if err != nil {
		log.Fatal(err)
	}
	connectPassword, err := conf.String("password")
	if err != nil {
		log.Fatal(err)
	}
	updateInterval, err := conf.Int("updateinterval")
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Connect to MQQT server")
	opts := mqtt.NewClientOptions()
	opts.AddBroker(connectURL)
	opts.SetClientID(HostNameSafe) // Every device that connects to the broker needs a unique ID
	opts.SetUsername(connectUsername)
	opts.SetPassword(connectPassword)
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Info("Network error: Sleeping 5 minutes before allowing service to reconnect")
		time.Sleep(time.Duration(300) * time.Second)
		log.Fatal(token.Error())
	}

	log.Info("Subscribe to 'command' topic")
	sub(client)

	log.Info("Publish HA auto-config sensors")
	publishHomeAssistantAutoConfigData(client, "system", "cpu", "%", "mdi:speedometer")
	publishHomeAssistantAutoConfigData(client, "system", "hostname", "", "mdi:account")
	if runtime.GOOS == "darwin" {
		publishHomeAssistantAutoConfigData(client, "system", "os", "", "mdi:apple")
	}
	if runtime.GOOS == "windows" {
		publishHomeAssistantAutoConfigData(client, "system", "os", "", "mdi:windows")
	}
	if runtime.GOOS == "linux" {
		publishHomeAssistantAutoConfigData(client, "system", "os", "", "mdi:linux")
	}
	publishHomeAssistantAutoConfigData(client, "system", "user", "", "mdi:account")
	publishHomeAssistantAutoConfigData(client, "system", "uptime", "", "mdi:calendar-clock")
	publishHomeAssistantAutoConfigData(client, "power", "battery", "%", "mdi:battery")
	publishHomeAssistantAutoConfigData(client, "net", "iplocal", "", "mdi:lan")
	publishHomeAssistantAutoConfigData(client, "net", "ipwan", "", "mdi:wan")
	publishHomeAssistantAutoConfigData(client, "memory", "mem_total", "GB", "mdi:memory")
	publishHomeAssistantAutoConfigData(client, "memory", "mem_used", "GB", "mdi:speedometer")
	publishHomeAssistantAutoConfigData(client, "memory", "mem_free", "GB", "mdi:speedometer")
	publishHomeAssistantAutoConfigData(client, "memory", "mem_percent", "%", "mdi:speedometer")
	publishHomeAssistantAutoConfigData(client, "disk", "disk_total", "GB", "mdi:harddisk")
	publishHomeAssistantAutoConfigData(client, "disk", "disk_used", "GB", "mdi:speedometer")
	publishHomeAssistantAutoConfigData(client, "disk", "disk_free", "GB", "mdi:speedometer")

	log.Info("Publish device stats in a loop")
	publishStats(client, updateInterval)

	client.Disconnect(250)
}

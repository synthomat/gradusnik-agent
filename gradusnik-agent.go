package main

import (
	"encoding/json"
	"flag"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	host "github.com/shirou/gopsutil/host"
	"os"
	"strings"
	"time"
)

type Sensor struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
	Host  string  `json:"host"`
	Time  string  `json:"time"`
	Type string `json:"type"`
}

func main() {
	hostname, _ := os.Hostname()

	topic := flag.String("topic", "sensors", "The topic name to/from which to publish/subscribe")
	broker := flag.String("broker", "tcp://192.168.1.105:1883", "The broker URI. ex: tcp://10.10.1.1:1883")
	password := flag.String("password", "", "The password (optional)")
	user := flag.String("user", "", "The User (optional)")
	flag.Parse()

	opts := mqtt.NewClientOptions().AddBroker(*broker)
	opts.SetUsername(*user)
	opts.SetPassword(*password)

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	tempData, _ := host.SensorsTemperatures()

	for _, temp := range tempData {
		sensor := Sensor{
			Name:  strings.ToLower(temp.SensorKey),
			Value: temp.Temperature,
			Host:  hostname,
			Time:  time.Now().Format(time.RFC3339),
			Type: "temp",
		}
		data, _ := json.Marshal(sensor)
		token := c.Publish(*topic, 0, false, string(data))
		token.Wait()
	}
}

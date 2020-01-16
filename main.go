//
// author: armaan roshani
//
// derived partially from: github.com/go-ble/ble/examples/basic/scanner
//					  and: github.com/eclipse/paho.mqtt.golang/cmd/sample
//

package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"
	//"os"

	"github.com/go-ble/ble"
	"github.com/go-ble/ble/examples/lib/dev"
	"github.com/pkg/errors"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type advType struct {
	c    chan string
	name string
}

var (
	// BLE
	device = flag.String("device", "default", "implementation of ble")
	du     = flag.Duration("du", 5*time.Second, "scanning duration")
	dup    = flag.Bool("dup", true, "allow duplicate reported")

	// MQTT
	topic     = flag.String("topic", "ble/test", "The topic name to publish to. (default ble/test)")
	broker    = flag.String("broker", "tcp://127.0.0.1:1883", "The broker URI. ex: tcp://10.10.1.1:1883")
	password  = flag.String("password", "", "The password (optional)")
	user      = flag.String("user", "", "The User (optional)")
	id        = flag.String("id", "bleTest", "The ClientID (optional)")
	cleansess = flag.Bool("clean", false, "Set Clean Session (default false)")
	qos       = flag.Int("qos", 0, "The Quality of Service 0,1,2 (default 0)")
	store     = flag.String("store", ":memory:", "The Store Directory (default use memory store)")
)

// values to convert to JSON
type BleDev struct {
	Addr     string
	Rssi     int64
	Contable bool
	Name     string
	Services []string
	ManData  []byte
	//Scanned 	time.Time
}

func main() {
	flag.Parse()

	// BLE setup
	d, err := dev.NewDevice(*device)
	if err != nil {
		log.Fatalf("can't new device : %s", err)
	}
	ble.SetDefaultDevice(d)

	// MQTT setup
	opts := MQTT.NewClientOptions()
	opts.AddBroker(*broker)
	opts.SetClientID(*id)
	opts.SetUsername(*user)
	opts.SetPassword(*password)
	opts.SetCleanSession(*cleansess)
	if *store != ":memory:" {
		opts.SetStore(MQTT.NewFileStore(*store))
	}

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	myChan := make(chan string, 1000)
	// create new object with channel
	newAdv := &advType{
		c:    myChan,
		name: "main advertiser", // give this advertiser a name
	}

	go bleScan(*newAdv)
	mqttSend(client, *newAdv)
}

func bleScan(newAdv advType) {
	// BLE Scan for specified duration, or until interrupted by user.
	fmt.Printf("Scanning for %s...\n", *du)
	ctx := ble.WithSigHandler(context.WithTimeout(context.Background(), *du))
	chkErr(ble.Scan(ctx, *dup, newAdv.advHandler, nil)) // pass adv method to scan func
	defer close(newAdv.c)
}

// make advHandler method of adv type
func (adv *advType) advHandler(a ble.Advertisement) {

	bDev := BleDev{
		Addr:     a.Addr().String(),
		Rssi:     int64(a.RSSI()),
		Contable: false,
	}

	if a.Connectable() {
		bDev.Contable = true
	}
	if len(a.LocalName()) > 0 {
		bDev.Name = a.LocalName()
	}
	if len(a.Services()) > 0 {
		//bDev.Services = string(a.Services())
	}
	if len(a.ManufacturerData()) > 0 {
		bDev.ManData = a.ManufacturerData()
	}

	var jsonData []byte
	jsonData, err := json.Marshal(bDev)
	if err != nil {
		log.Println(err)
	}
	//fmt.Println(string(jsonData))
	adv.c <- string(jsonData)
}

func mqttSend(client MQTT.Client, adv advType) {
	var i int
	for m := range adv.c {
		client.Publish(*topic, byte(*qos), false, m)
		i++
		//fmt.Println("In mqttSend: ", m)
	}
	fmt.Println("Sent ", i, " MQTT messages")
}

func chkErr(err error) {
	switch errors.Cause(err) {
	case nil:
	case context.DeadlineExceeded:
		fmt.Printf("done\n")
	case context.Canceled:
		fmt.Printf("canceled\n")
	default:
		log.Fatalf(err.Error())
	}
}

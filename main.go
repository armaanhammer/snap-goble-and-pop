// 
// author: armaan roshani
//
// derived partially from: github.com/go-ble/ble/examples/basic/scanner
//					  and: github.com/eclipse/paho.mqtt.golang/cmd/sample
//

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"
	"encoding/json"
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
	du     = flag.Duration("du", 2*time.Second, "scanning duration")
	dup    = flag.Bool("dup", true, "allow duplicate reported")

	// MQTT
	topic  	  = flag.String("topic", "ble/test", "The topic name to publish to. (default ble/test)")
	broker 	  = flag.String("broker", "tcp://127.0.0.1:1883", "The broker URI. ex: tcp://10.10.1.1:1883")
	password  = flag.String("password", "", "The password (optional)")
	user 	  = flag.String("user", "", "The User (optional)")
	id 		  = flag.String("id", "bleTest", "The ClientID (optional)")
	cleansess = flag.Bool("clean", false, "Set Clean Session (default false)")
	qos 	  = flag.Int("qos", 0, "The Quality of Service 0,1,2 (default 0)")
	store 	  = flag.String("store", ":memory:", "The Store Directory (default use memory store)")
)

// values to convert to JSON
type BleDev struct {
    Addr		string
	Rssi		int64
	Contable	bool
    Name    	string
	Services	[]string
    ManData     []byte
    Scanned 	time.Time
}

func main() {
	flag.Parse()

	//messages := make(chan string, 100)

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

	//fmt.Println("Sample Publisher Started")
	//for i := 0; i < *num; i++ {
	//	fmt.Println("---- doing publish ----")
	//	token := client.Publish(*topic, byte(*qos), false, *payload)
	//	token.Wait()
	//}

	//client.Disconnect(250)
	//fmt.Println("Sample Publisher Disconnected")


	myChan := make(chan string, 1000)
	// create new object with channel
	newAdv := &advType{
		c:    myChan, 
		name: "main advertiser",         // give this advertiser a name
	}


	bleScan(*newAdv)
	mqttSend(client, *newAdv)


}


func bleScan(newAdv advType) {
	// BLE Scan for specified duration, or until interrupted by user.
	fmt.Printf("Scanning for %s...\n", *du)
	ctx := ble.WithSigHandler(context.WithTimeout(context.Background(), *du))
	chkErr(ble.Scan(ctx, *dup, newAdv.advHandler, nil))   // pass adv method to scan func
	defer close(newAdv.c)
}


//func advHandler(a ble.Advertisement) {

// make advHandler method of adv type
func (adv *advType) advHandler(a ble.Advertisement) {
	
	// print the adv's name
	//fmt.Printf("Hello from inside adv: %s", adv.name)
	// send something on the channel
	//adv.c <- "hello!"

	bDev := BleDev{
		Addr:		a.Addr().String(),
		Rssi:		int64(a.RSSI()),
		Contable:	false,
	}

	if a.Connectable() {
		bDev.Contable = true
		//fmt.Printf("[%s] C %3d:", a.Addr(), a.RSSI())
	} else {
		//fmt.Printf("[%s] N %3d:", a.Addr(), a.RSSI())
	}
	//comma := ""
	if len(a.LocalName()) > 0 {
		bDev.Name = a.LocalName()
		//fmt.Printf(" Name: %s", a.LocalName())
		//comma = ","
	}
	if len(a.Services()) > 0 {
		//bDev.Services = string(a.Services())
		//fmt.Printf("%s Svcs: %v", comma, a.Services())
		//comma = ","
	}
	if len(a.ManufacturerData()) > 0 {
		bDev.ManData = a.ManufacturerData()
		//fmt.Printf("%s MD: %X", comma, a.ManufacturerData())
	}
	//fmt.Printf("\n")

	var jsonData []byte
	jsonData, err := json.Marshal(bDev)
	if err != nil {
    	log.Println(err)
	}
	//testString := string(jsonData)
	fmt.Println(string(jsonData))

	adv.c <- string(jsonData)

	

	//var client MQTT.Client
	//client := MQTT.Client
}


func mqttSend(client MQTT.Client,adv advType) {
	for m := range adv.c {
		//test := <-m
		//token := client.Publish(*topic, byte(*qos), false, test)
		client.Publish(*topic, byte(*qos), false, m)
		fmt.Println("In mqttSend: ", m)
   		//fmt.Printf("I read string: %s from channel\n", m)
	}
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

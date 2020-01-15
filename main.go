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

var (
	// BLE
	device = flag.String("device", "default", "implementation of ble")
	du     = flag.Duration("du", 5*time.Second, "scanning duration")
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

	// BLE Scan for specified duration, or until interrupted by user.
	fmt.Printf("Scanning for %s...\n", *du)
	ctx := ble.WithSigHandler(context.WithTimeout(context.Background(), *du))
	chkErr(ble.Scan(ctx, *dup, advHandler, nil))

	client.Disconnect(250)
}


func advHandler(a ble.Advertisement) {
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
	testString := string(jsonData)
	fmt.Println(string(jsonData))

	client := MQTT.Client
	token := client.Publish(*topic, byte(*qos), false, testString)
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

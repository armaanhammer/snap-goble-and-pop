// 
// author: armaan roshani
//
// derived from: github.com/go-ble/ble/examples/basic/scanner
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
	//"github.com/eclipse/paho.mqtt.golang"
)

var (
	device = flag.String("device", "default", "implementation of ble")
	du     = flag.Duration("du", 5*time.Second, "scanning duration")
	dup    = flag.Bool("dup", true, "allow duplicate reported")
)

// Struct to hold values to convert to JSON
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

	d, err := dev.NewDevice(*device)
	if err != nil {
		log.Fatalf("can't new device : %s", err)
	}
	ble.SetDefaultDevice(d)



	// Scan for specified durantion, or until interrupted by user.
	fmt.Printf("Scanning for %s...\n", *du)
	ctx := ble.WithSigHandler(context.WithTimeout(context.Background(), *du))
	chkErr(ble.Scan(ctx, *dup, advHandler, nil))
}

func advHandler(a ble.Advertisement) {
	bDev := BleDev{
		Addr:		a.Addr().String(),
		Rssi:		int64(a.RSSI()),
		Contable:	false,
	}

	if a.Connectable() {
		bDev.Contable = true
		fmt.Printf("[%s] C %3d:", a.Addr(), a.RSSI())
	} else {
		fmt.Printf("[%s] N %3d:", a.Addr(), a.RSSI())
	}
	comma := ""
	if len(a.LocalName()) > 0 {
		bDev.Name = a.LocalName()
		fmt.Printf(" Name: %s", a.LocalName())
		comma = ","
	}
	if len(a.Services()) > 0 {
		//bDev.Services = string(a.Services())
		fmt.Printf("%s Svcs: %v", comma, a.Services())
		comma = ","
	}
	if len(a.ManufacturerData()) > 0 {
		bDev.ManData = a.ManufacturerData()
		fmt.Printf("%s MD: %X", comma, a.ManufacturerData())
	}
	fmt.Printf("\n")

	var jsonData []byte
	jsonData, err := json.Marshal(bDev)
	if err != nil {
    	log.Println(err)
	}
	fmt.Println(string(jsonData))
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

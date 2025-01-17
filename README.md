# Snap, GoBLE and Pop
"snap, gobble [and pop](https://en.wikipedia.org/wiki/Snap,_Crackle_and_Pop)"

An IoT project written in Go for BLE devices, to be distributed as a [Snap](https://snapcraft.io/). 

Uses:
- [Go-BLE](http://github.com/go-ble/ble)
- [Paho](http://github.com/eclipse/paho.mqtt.golang)

The Snap binary currently exists in this repo. It will be added to the [Snap Store](https://snapcraft.io/store) once it is past proof-of-concept stage.

## Current
Current iteration is proof-of-concept. When activated, it scans for all BLE devices, packages info found about devices as JSON, then sends that JSON out as MQTT.

#### To install snap:
1. Clone this git repo.
1. Navigate to /snap/ directory within repo.
1. Enter: `sudo snap install --devmode snap-goble-and-pop_0.2_amd64.snap `

#### To activate the snap:
1. Make sure /snap/bin/ exists in your $PATH
1. Make sure you have a BLE device attached to your system
1. Enter: `sudo snap-goble-and-pop`

#### By default:
- The snap connects to an unencrypted broker on localhost port 1883 running TCP (default "tcp://127.0.0.1:1883"). This can be changed by passing the `-broker` flag to the snap. 
- The snap scans for a duration of five seconds (default 5s). This can be changed by passing the `-du` flag.
- The snap publishes to MQTT topic ble/test (default "ble/test"). This can be changed by passing the `-topic` flag.

A complete list of flags can be generated by passing the --help flag.

#### To view the MQTT messages using the defaults:
1. Run an MQTT broker such as [Mosquitto](https://mosquitto.org/) on the local machine.
1. Point an MQTT client to the broker and subscribe to the topic. For example, if using mosquitto.sub, enter: `mosquitto.sub -h 127.0.0.1 -p 1883 -t ble/#`.

## Future
Plan is to develop a modular thermostat using BLE thermometers, a central controller, and relays, IR blasters, and device APIs to control window air conditioners and plug-in heaters.

#### Next steps:
- Add security configurations to snap so snap can connect to AWS for example.
- Target known BLE devices and pass only relevant information.

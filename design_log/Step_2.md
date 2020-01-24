# Pipe data over MQTT

## Dataflow

Basing device data processing on go-ble library scanner example. Example writes info about all ble devices to standard out. Instead, I want to send MQTT data.

### Problem

Having troule figuring out how to pass a channel or any other variable into and out of [`advHandler` function](https://github.com/go-ble/ble/blob/67ae735460239645b70e728470f4b725f505c459/examples/basic/scanner/main.go#L36). I have:
`client := MQTT.NewClient(opts)` in `func main()` so it only connects to the broker once, but I can't figure out how to pass `client` into `advHandler`. Ideally I'd like a sperate routine that handles MQTT, but I run into the same problem trying to pass a channel into it.

### Sought help

I had been talking to Oliver R. about Go recently, so reached out to him to see if he had any ideas. He suggested making a type whose method is `advHandler` and instantiating a new type with the channel inside it, like so:

```
// type that will have an advHandler method
type advType struct {
	c    chan string
	name string
}

func main() {
  ...
}

// make advHandler method of adv type
func (adv *advType) advHandler(a ble.Advertisement) {
	// print the adv's name
	fmt.Printf("Hello from inside adv: %s", adv.name)
	// send something on the channel
	adv.c <- "hello!"
```

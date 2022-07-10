package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	mux "github.com/gorilla/mux"
)

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

// defaultHandler is a http request handler for route / .
func defaultHandler(w http.ResponseWriter, r *http.Request) {
	currentTime := time.Now()
	helloMsg := "Hello IOx " + currentTime.Format("2006.01.02 15:04:05") + "\n"
	log.Printf(helloMsg)

	w.Write([]byte(helloMsg))
}

// pingHandler is a http request handler for route /ping .
func pingHandler(w http.ResponseWriter, r *http.Request) {
	currentTime := time.Now()
	mqttMsg := fmt.Sprintf("Hello from sebBeagle %s", currentTime.Format("2006.01.02 15:04:05"))

	doPublish(mqttMsg)

	pongMsg := "Published mqtt message from sebBeagle - " + mqttMsg + "\n"
	w.Write([]byte(pongMsg))
}

func doPublish(msg string) {

	publish(mqttClient, MqttTopic, msg)

	pongMsg := "Published mqtt message from sebBeagle - " + msg + "\n"
	log.Printf(pongMsg)
}

func doPublishLoop() {
	time.Sleep(5 * time.Second) // delay start
	for {
		currentTime := time.Now()
		mqttMsg := fmt.Sprintf("Hello from sebBeagle %s", currentTime.Format("2006.01.02 15:04:05"))

		publish(mqttClient, MqttTopic, mqttMsg)

		time.Sleep(10 * time.Second)
	}
}

const DefaultMqttQoS = 1
const MqttTopic = "devices/sebBeagle/messages/events/" // topic - devices/{device_id}/modules/{module_id}/messages/events/
var mqttClient mqtt.Client

// initializes mqtt client connection to broker
func init() {
	var broker = "<hubname>.azure-devices.net"
	var port = 8883
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tls://%s:%d", broker, port))
	opts.SetProtocolVersion(4)
	opts.SetClientID("sebBeagle")
	opts.SetUsername("<hubname>.azure-devices.net/sebBeagle/?api-version=2020-09-30")

	// TODO: Need to manually generate SAS for now and paste into code. Actual build script can generate and store in ENV. Code can then read from ENV.
	dummySAS := "SharedAccessSignature sr=<hubname>.azure-devices.net%2Fdevices%2F<device name>&sig=...%2F..."
	// az iot hub generate-sas-token -d sebBeagle -n seb-hub --du 86400
	// arm32 - GOOS=linux GOARCH=arm GOARM=5 go build -o gomqttpubarm32v7 main.go
	// Monitor event hub from <path to codein>/go/eventhub
	opts.SetPassword(dummySAS)

	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	mqttClient = mqtt.NewClient(opts)
}

func main() {
	// Setting up a simple HTTP REST /ping request - where tester can ping to send mqtt msg
	portPtr := flag.String("port", "8282", "port number")
	flag.Parse()

	httpPort := *portPtr
	httpURL := "0.0.0.0:" + httpPort
	log.Printf("HTTP %s up and listening...\n", httpURL)

	r := mux.NewRouter()
	// Routes consist of a path and a handler function.
	r.HandleFunc("/ping", pingHandler)
	r.HandleFunc("/", defaultHandler)

	go doPublishLoop() // A go-routine to send mqtt in a loop

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(httpURL, r))
}

func publish(client mqtt.Client, topic string, msg string) {

	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	defer mqttClient.Disconnect(250)

	log.Printf("Publishing to topic: %s\n", topic)
	log.Printf("Sending message: %s\n", msg)
	token := client.Publish(topic, DefaultMqttQoS, false, msg)
	token.Wait()
}

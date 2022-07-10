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
	helloMsg := "Hello gomqttpub " + currentTime.Format("2006.01.02 15:04:05") + "\n"
	log.Printf(helloMsg)

	w.Write([]byte(helloMsg))
}

// pingHandler is a http request handler for route /ping .
func pingHandler(w http.ResponseWriter, r *http.Request) {
	currentTime := time.Now()
	mqttMsg := fmt.Sprintf("Hello %s", currentTime.Format("2006.01.02 15:04:05"))
	pongMsg := "Publishing mqtt message - " + mqttMsg + "\n"
	log.Printf(pongMsg)

	publish(mqttClient, MqttTopic, mqttMsg)

	w.Write([]byte(pongMsg))
}

const DefaultMqttQoS = 1
const MqttTopic = "devices/sebEdgeDevice/messages/events/" // topic - devices/{device_id}/modules/{module_id}/messages/events/
var mqttClient mqtt.Client

// initializes mqtt client connection to broker
func init() {
	var broker = "seb-hub.azure-devices.net"
	var port = 8883
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tls://%s:%d", broker, port))
	opts.SetProtocolVersion(4)
	opts.SetClientID("sebEdgeDevice")
	opts.SetUsername("seb-hub.azure-devices.net/sebEdgeDevice/?api-version=2020-09-30")

	// TODO: Need to manually generate SAS
	dummySAS := "SharedAccessSignature sr=<hubname>.azure-devices.net%2Fdevices%2F<device name>&sig=...%2F..."

	// az iot hub generate-sas-token -d sebEdgeDevice -n seb-hub --du 28800
	// amd64 - CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o gomqttpub main.go
	opts.SetPassword(dummySAS)

	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	mqttClient = mqtt.NewClient(opts)
}

func main() {
	// Setting up a simple HTTP REST /ping request
	portPtr := flag.String("port", "8282", "port number")
	flag.Parse()

	httpPort := *portPtr
	httpURL := "0.0.0.0:" + httpPort
	log.Printf("HTTP %s up and listening...\n", httpURL)

	r := mux.NewRouter()
	// Routes consist of a path and a handler function.
	r.HandleFunc("/ping", pingHandler)
	r.HandleFunc("/", defaultHandler)

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
	/*
		num := 5
		for i := 0; i < num; i++ {
			text := fmt.Sprintf("Hello-Message-%d", i)
			log.Printf("Sending message: %s\n", text)
			token := client.Publish(topic, DefaultQoS, false, text)
			token.Wait()
			time.Sleep(time.Second)
		}
	*/
}

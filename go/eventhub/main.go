package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	eventhub "github.com/Azure/azure-event-hubs-go/v3"
)

func main() {
	hosturl := "sb://<hub>namespace.servicebus.windows.net/"
	connStr := "Endpoint=" + hosturl + ";SharedAccessKeyName=service;SharedAccessKey=...;EntityPath=iothub-ehub-..."
	hub, err := eventhub.NewHubFromConnectionString(connStr)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Connected to event hub - " + hosturl)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// send a single message into a random partition
	// requires Azure admin update Event Hub to allow SEND permission
	// https://stackoverflow.com/questions/49777839/how-to-send-a-message-to-an-event-hub-from-an-azure-function
	//err = hub.Send(ctx, eventhub.NewEventFromString("hello go eventhub!"))
	//if err != nil {
	//	fmt.Printf("Failed send event message:\n%v\n", err)
	//	return
	//}

	handler := func(c context.Context, event *eventhub.Event) error {
		fmt.Println(string(event.Data))
		return nil
	}

	// listen to each partition of the Event Hub
	runtimeInfo, err := hub.GetRuntimeInformation(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, partitionID := range runtimeInfo.PartitionIDs {
		// Start receiving messages
		//
		// Receive blocks while attempting to connect to hub, then runs until listenerHandle.Close() is called
		// <- listenerHandle.Done() signals listener has stopped
		// listenerHandle.Err() provides the last error the receiver encountered
		//listenerHandle, err := hub.Receive(ctx, partitionID, handler, eventhub.ReceiveWithLatestOffset())
		_, err := hub.Receive(ctx, partitionID, handler, eventhub.ReceiveWithLatestOffset())
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	// Wait for a signal to quit:
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	<-signalChan

	err = hub.Close(context.Background())
	if err != nil {
		fmt.Println(err)
	}
}

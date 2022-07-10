# azure-iot-poc
POC for Azure IoT Edge and Azure IoT Hub.  
Lots of provisioning and connectivity steps, but summarized as below for *"a picture paints a thousands words"*.  

Where available, individual READMEs are provided in individual packages.  

Codes are samples-only, developers should follow respective coding styles and quality guidelines for product developement.  

Feel free to reach out to the author for any clarifications.  

## Go

This section is purely Go.  
Pre-requisites for Go development environment are assumed.

### go/pub/mqtt
To build, please do a `go mod init <your path>` again.  

This is a first attempt on using MQTT in my Azure registered iot edge device sebEdgeDevice.  
This package demonstrates using pure MQTT to publish to Azure IoT Hub.  

### go/sub/amqp
To build, please do a `go mod init <your path>` again.  

The Azure IoT Hub is seb-hub.  
Microsoft has Azure AMQP SDK for Go, but not a subset for Azure IoT SDK for Go.  
This package demonstrates using underlying Go Azure AMQP SDK to receive MQTT messages.  
It is a rundown version of amenzhinsky's codes, with no abstraction from the plumbery of AMQP and Azure Event Hub.  
Notn for faint-hearted, you can ignore this by using Azure Event Hub below.

### go/eventhub
To build, please do a `go mod init <your path>` again.  

This package demonstrates receiving IoT Hub messages from Azure Event Hub.  


## Azure IoT SDK

This section contains Azure IoT SDK codes using NodeJs.  
Custom Go is also possible since we are merely deploying Go binary in Docker.  

### Device-to-Cloud azure-iot-sdk/sebEdgeGoMqttPub
This is similar to [go/pub/mqtt] above but just using VS Code Azure IoT extension for packaging.  

To build, please do a `go mod init <your path>` again.  
Do test it out in your local Docker before deployment to Edge device.  

This application is first deployed to Azure Repository and then pushed down to Edge Device.  

### Device-to-Cloud azure-iot-sdk/sebEdgeGoMqttPubModuleId
This is similar to [azure-iot-sdk/sebEdgeGoMqttPub] above but ModuleId is include as part of MQTT topic.  

### Cloud-to-Device azure-iot-sdk/sebNodeJS

Demonstrates Cloud-to-Device by the device exposing Azure IoT SDK **Direct Method**.  
Direct Methods can also be entry points as a *Facade* to internal Go apps via local IPC/GRPC/REST.  

### Cloud-to-Device azure-iot-sdk/sebNodeJsModuleId
This is similar to [azure-iot-sdk/sebNodeJS] above but ModuleId is added as part the device identity.  
This is so that messages are directed towards the application module in the device.

## NodeJS

### azure-iot-sdk/nodejs/iot-hub/back-end-application

This is a backend application demonstrating C2D messaging via **Direct Method** exposed by the Edge device.  
BackEndApplication.js calls the sebEdgeDevice to change the sending interval.  
CloudToBeagle.js calls BeagleBone to do the same.  

### azure-iot-sdk/nodejs/iot-hub/back-end-application-moduleid

This is similar to [azure-iot-sdk/nodejs/iot-hub/back-end-application] but **Direct Method** is directed towards a ModuleId in the Edge device.

### azure-iot-sdk/nodejs/iot-hub/read-d2c-messages

This application receives IoT Hub messages from other devices via Azure Event Hub.  
The concept is similar to [go/eventhub].  

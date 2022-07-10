# GoMqttPubModuleId

Demonstrates Go application publishing MQTT message to Azure IoT Hub Topic with **ModuleId**.  
```sh
devices/{device_id}/modules/{module_id}/messages/events/
```

# Build and Deploy Steps

## Edge Device
Ensure Edge Device is on-line and its Azure IoT Edge Runtime service is up and running.  

## Generate SAS Key

Manually generate SAS for now and paste into code.  
Actual build script can generate and store in ENV.  
Code can then read from ENV.  
E.g. SAS for app in device publishing to a ModuleId named *NodeJsModuleId*  
```sh
az iot hub generate-sas-token -d sebEdgeDevice -m NodeJsModuleId -n seb-hub --du 86400
```

## Build Go binary

For Linux amd64 like Ubuntu:  
```sh
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o gomqttpubmoduleid main.go
```

For Linux arm32 like Raspberry PI and BeagleBone:  
```sh
GOOS=linux GOARCH=arm GOARM=5 go build -o gomqttpubmoduleidarm32v7 main.go
```

## Build and Push IoT Edge Solution using VS Code

For this, your VS must have Azure IoT Edge add-on and your machine must have Docker:  

Right-click on **deployment.template.json**, select **Build and Push IoT Edge Solution**.  

Verify the image is deployed in Azure Registry.  

## Build using Az CLI command

If developer machine does not have Docker, you can queue your build to Azure cloud.

Example 1 - Build and Push NodeJs app for arm32v7 to Azure Registry(a DockerHub):  
```sh
az acr build -t sebregistry.azurecr.io/nodejsmodule_arm32v7:0.0.1-arm32v7 -r sebregistry . -f Dockerfile.arm32v7 --platform linux/arm/v7
```

Example 2 - Build and Push Go app for arm32v7 to Azure Registry(a DockerHub):  
The Go binary specified in Dockerfile is already be pre-built for arm, so just deploy it to Azure Registry  
```sh
az acr build -t sebregistry.azurecr.io/gomqttpubmodule_arm32v7:0.0.1-arm32v7 -r sebregistry . -f Dockerfile.arm32v7 --platform linux/arm/v7
```

Verify the image is deployed in Azure Registry.  

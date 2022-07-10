# Build and Deploy Steps

## Edge Device
Ensure Edge Device is on-line and its Azure IoT Edge Runtime service is up and running.  

## Generate SAS Key

Manually generate SAS for now and paste into code.  
Actual build script can generate and store in ENV.  
Code can then read from ENV.  
E.g.  
```sh
az iot hub generate-sas-token -d sebBeagle -n seb-hub --du 86400
```

## Build Go binary

For Linux amd64 like Ubuntu:  
```sh
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o gomqttpub main.go
```

For Linux arm32 like Raspberry PI and BeagleBone:  
```sh
GOOS=linux GOARCH=arm GOARM=5 go build -o gomqttpubarm32v7 main.go
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

## More about Docker and ARM-based Builds

https://docs.docker.com/docker-for-mac/multi-arch/  
https://www.danielstechblog.io/building-arm-based-container-images-with-vsts-for-your-azure-iot-edge-deployments/  


## Push the Deployed Image to Edge Device
This is the actual deployment of your Docker image from IoT Hub to the Edge device.  

From IoT Hub | IoT Edge | Devices | <the device id> | Set Modules  
- add the image to be pushed down to device.  
E.g. 
`sebregistry.azurecr.io/gomqttpubmodule:0.0.1-amd64`  
`sebregistry.azurecr.io/gomqttpubmodule_arm32v7:0.0.1-amd64`  

Verify the image is deployed and running.  

## Device-to-Cloud Monitor IoT Hub Events
In general, use CLI:  
```sh
az iot hub monitor-events -n seb-hub
```

## Device-to-Cloud AMQP Subscribing from Azure Event Hub
This simulates a Cloud application receiving messages from Device.  
go run the subscriber from:  
`<path-to>/code.in.spdigital.sg/sebastianmlk/azure-iot-poc/go/eventhub`  

## Test MQTT Publishing
Do this step if your Go module is not publishing MQTT message automatically.  
SSH to the device, get the running go module container id, get the running container IP address, then publish test message:  
```sh
curl http:<container IP>:8282/ping
```

## My Docker cheatsheet

Below are for my FYIs.  

### Refresher steps for docker build and run
My docker basics - https://github.com/sebmacisco/cisco-iox-go/tree/master/gosafeentry/gateway  
```sh
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o gomqttpub main.go
docker build -t gomqttpub-alpine -f Dockerfile.alpine8282 .
sudo docker run -p 8282:8282 --entrypoint=/bin/sh sebregistry.azurecr.io/gomqttpubmodule:0.0.1-amd64
sudo docker run -p 8282:8282 --rm -it --entrypoint=/bin/sh sebregistry.azurecr.io/gomqttpubmodule:0.0.1-amd64
docker run -d -p 8282:8282 --entrypoint=/bin/sh sebregistry.azurecr.io/gomqttpubmodule
```

To attach to a running Alpine container:  
```sh
docker exec -i -t <containerId/Name> /bin/sh
```

### Get the IP of Docker Container
```sh 
sudo docker inspect -f \
'{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' \
<container id>
172.18.0.4
sebma@vm-4smhrehgorf5i:~$ curl http://172.18.0.4:8282/ping
Publishing mqtt message - Hello 2021.04.07 09:24:26
```

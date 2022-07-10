# sebNodeJS

Demonstrates Cloud-to-Device using Azure IoT SDK **Direct Method**.  

## Start read-d2c-message (optional)
If you already have a Go app listening to Azure Event Hub you can skip this.  
Go to: 
`<path-to>/nodejs/iot-hub/Quickstarts/read-d2c-messages`  

## Start Edge Device
Ensure `iotedge` service on edge device is on-line and its iot edge runtime and this **NodeJsModule** is started.  

To view actual invocation of **Direct Method**, SSH to device to see real-time logs:  
```sh
sudo iotedge logs -f NodeJsModule
```

Go eventhub app or *read-d2c-message* should also be displaying messages it receives.  

## Send Direct Method

Run the app from `<path-to>/nodejs/iot-hub/Quickstarts/back-end-application`  

## Capturing underlying Build and Push IoT Edge Solution

This is for automating builds in your CI/CD process e.g. using Jenkins  

```sh
$ docker build  --rm -f "/home/ubuntu/go/src/github.com/sebmaspd/rnd/azure/dotnet/sebNodeJsModuleId/modules/NodeJsModuleId/Dockerfile.amd64" -t sebregistry.azurecr.io/nodejsmoduleid:0.0.1-amd64 "/home/ubuntu/go/src/github.com/sebmaspd/rnd/azure/dotnet/sebNodeJsModuleId/modules/NodeJsModuleId" && docker push sebregistry.azurecr.io/nodejsmoduleid:0.0.1-amd64
Sending build context to Docker daemon  18.43kB
Step 1/7 : FROM node:10-alpine
 ---> 863024ec4a19
Step 2/7 : WORKDIR /app/
 ---> Using cache
 ---> b2a0c9c3a13d
Step 3/7 : COPY package*.json ./
 ---> 2cfe2f049739
Step 4/7 : RUN npm install --production
 ---> Running in 8aa71184cb88
npm notice created a lockfile as package-lock.json. You should commit this file.
added 153 packages from 210 contributors and audited 153 packages in 11.644s

5 packages are looking for funding
  run `npm fund` for details

found 0 vulnerabilities

Removing intermediate container 8aa71184cb88
 ---> 077331fc6b92
Step 5/7 : COPY app.js ./
 ---> 8fbf0355417e
Step 6/7 : USER node
 ---> Running in aecb3edb307f
Removing intermediate container aecb3edb307f
 ---> 2e74b183b163
Step 7/7 : CMD ["node", "app.js"]
 ---> Running in 927eacfd795c
Removing intermediate container 927eacfd795c
 ---> 25beb7b9d032
Successfully built 25beb7b9d032
Successfully tagged sebregistry.azurecr.io/nodejsmoduleid:0.0.1-amd64
The push refers to repository [sebregistry.azurecr.io/nodejsmoduleid]
d1a5115f5a72: Pushed 
43227c173e5d: Pushed 
665337bdf9fe: Pushed 
1be332370db5: Mounted from nodejsmodule 
fd7d7fddbeff: Mounted from nodejsmodule 
6f154f1607cd: Mounted from nodejsmodule 
cbd8fabb2ce5: Mounted from nodejsmodule 
2b2bcc6e6724: Mounted from nodejsmodule 
0.0.1-amd64: digest: sha256:2431298712032c02ed96698085fd687774cf990e510f7a4e8ac75170672a3f1b size: 1991
ubuntu@ubuntu1804:~/go/src/github.com/sebmaspd/rnd/azure/dotnet/sebNodeJS$ 
```
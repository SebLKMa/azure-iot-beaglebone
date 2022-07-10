# sebNodeJsModuleId

Similar to **sebNodeJS**.  
Demonstrates Cloud-to-Device using Azure IoT SDK **Direct Method**.  
However, a device can be running more than 1 application module.  

**ModuleId** is added as part the device identity, so that messages are directed towards the application module in the device.  

## az command line to create module id

See:  
https://docs.microsoft.com/en-us/cli/azure/iot/hub/module-identity?view=azure-cli-latest#az_iot_hub_module_identity_create  
```sh
az iot hub module-identity create --device-id
                                  --module-id
                                  [--am {shared_private_key, x509_ca, x509_thumbprint}]
                                  [--hub-name]
                                  [--login]
                                  [--od]
                                  [--primary-thumbprint]
                                  [--resource-group]
                                  [--secondary-thumbprint]
                                  [--valid-days]
```

```sh
ubuntu@ubuntu1804:~$ az iot hub module-identity create --device-id sebEdgeDevice --module-id NodeJsModuleId --hub-name seb-hub
```

```sh
az iot hub module-identity show --device-id
                                --module-id
                                [--hub-name]
                                [--login]
                                [--resource-group]
```

```sh
az iot hub module-identity show --device-id sebEdgeDevice --module-id NodeJsModuleId --hub-name seb-hub
```

## az command line to show module id connection string

**This command has to be run from Azure Cloud bash shell.**

See:  
https://docs.microsoft.com/en-us/cli/azure/iot/hub/module-identity/connection-string?view=azure-cli-latest  

```sh
az iot hub module-identity connection-string show --device-id
                                                  --module-id
                                                  [--hub-name]
                                                  [--key-type {primary, secondary}]
                                                  [--login]
                                                  [--resource-group]
```

**This command has to be run from Azure Cloud bash shell.**  
```sh
az iot hub module-identity connection-string show --device-id sebEdgeDevice --module-id NodeJsModuleId --hub-name seb-hub

```

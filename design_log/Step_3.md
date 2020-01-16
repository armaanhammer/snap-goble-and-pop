# Snapcraft

## Checklist
https://snapcraft.io/docs/snapcraft-checklist

### 1. Language/Framework/Build system (mandatory)
The go plugin
https://snapcraft.io/docs/go-plugin

### 2. Toolkits and desktop support (optional)
Most likely not needed

### 3. System integration (optional)
1. BLE
2. MQTT

----------------
Uncertain what interfaces I will need. Perhaps because I am talking to Go apps that already use interfaces, I will not need them. Chances are I will though, so will see what connections each of the apps makes.

## Global metadata
https://snapcraft.io/docs/adding-global-metadata

## Parts
https://snapcraft.io/docs/adding-parts


---------------
#### Notes:
Using an Ubuntu 16.04 VM inside of Virtual Box. Unable to use Snapcraft with Multipass, presumably because inside of a VM. Using lxd instead.


Strange error encountered when attempting to build snap: 
```
[Errno 2] No such file or directory: '/root/parts/snap-goble-and-pop/go/bin'
```
created the directory inside the container, then re-ran. Build seemed to work.

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
created the directory inside the container, then re-ran. Build seemed to work. Did not have a service though, so trying to rebuild including a service.

#### Trouble with service:

running command `snapcraft --use-lxd --debug` gives:

```
Pulling snap-goble-and-pop 
Please consider setting `go-importpath` for the 'snap-goble-and-pop' part
go get -t -d ./project/...
go: missing Git command. See https://golang.org/s/gogetcmd
package github.com/eclipse/paho.mqtt.golang: exec: "git": executable file not found in $PATH
go: missing Git command. See https://golang.org/s/gogetcmd
package github.com/go-ble/ble: exec: "git": executable file not found in $PATH
package github.com/go-ble/ble/examples/lib/dev: cannot find package "github.com/go-ble/ble/examples/lib/dev" in any of:
	/snap/go/4901/src/github.com/go-ble/ble/examples/lib/dev (from $GOROOT)
	/root/parts/snap-goble-and-pop/go/src/github.com/go-ble/ble/examples/lib/dev (from $GOPATH)
go: missing Git command. See https://golang.org/s/gogetcmd
package github.com/pkg/errors: exec: "git": executable file not found in $PATH
Failed to run 'go get -t -d ./project/...' for 'snap-goble-and-pop': Exited with code 1.
Verify that the part is using the correct parameters and try again.
snapcraft-hello # exit
exit
osboxes@osboxes:~/go/src/snap-goble-and-pop$ echo $PATH
/home/osboxes/bin:/home/osboxes/.local/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/games:/usr/local/games:/snap/bin:/snap/bin:/var/lib/snapd/snap/bin:/snap/bin:/var/lib/snapd/snap/bin:/snap/bin:/var/lib/snapd/snap/bin
osboxes@osboxes:~/go/src/snap-goble-and-pop$ which git
/usr/bin/git
osboxes@osboxes:~/go/src/snap-goble-and-pop$ echo $GOROOT

osboxes@osboxes:~/go/src/snap-goble-and-pop$
```
Uncertain whether paths and directories mentioned are within the virtual build environment or within the main system.

### Some success
After hours and much gnashing of teeth: found two problems building from git resources:
1. Bug in github.com/go-ble/ble
2. References in github.com/rigado/ble that point back to github.com/go-ble/ble which has bug.

Seems only way to build project currently is to pull github.com/go-ble/ble locally, remove the bug, and build using the local directory.

### No Joy

Unfortunately, even after cloning the repo and programmatically removing every reference to `go-ble` from it, snapcraft still returns this:

```
snapcraft-snap-goble-and-pop # snapcraft clean
Cleaning up staging area
Cleaning up parts directory
snapcraft-snap-goble-and-pop # snapcraft
The 'snap' directory is meant specifically for snapcraft, but it contains the following non-snapcraft-related paths, which is unsupported and will cause unexpected behavior:
- .snapcraft.yaml.swp
- snap-goble-and-pop_0.1_amd64.snap
- snapcraft.yaml.even.newer
- snapcraft.yaml.old
- snapcraft.yaml.old.newer
- wtf.yaml

If you must store these files within the 'snap' directory, move them to 'snap/local', which is ignored by snapcraft.
Pulling snap-goble-and-pop 
Cloning into '/root/parts/snap-goble-and-pop/src'...
done.
go get -t -d ./github.com/armaanhammer/snap-goble-and-pop/...
Building snap-goble-and-pop 
go build -o /root/parts/snap-goble-and-pop/go/bin/snap-goble-and-pop github.com/armaanhammer/snap-goble-and-pop
# github.com/go-ble/ble/linux/hci
github.com/go-ble/ble/linux/hci/hci.go:316:5: undefined: log
Failed to run 'go build -o /root/parts/snap-goble-and-pop/go/bin/snap-goble-and-pop github.com/armaanhammer/snap-goble-and-pop' for 'snap-goble-and-pop': Exited with code 2.
Verify that the part is using the correct parameters and try again.
snapcraft-snap-goble-and-pop # 
```

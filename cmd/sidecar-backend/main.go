package main

import (
	"fmt"
	sidecar "github.com/Lmare/lightning-playground/sidecar"
)


const namespace = "lightning-playground"
const secretlnd = "lnd-credentials"
const mountedBackendStorage = "/app/nodes-storage"

func main() {
	fmt.Println("[Sidecar backend] Starting")
	getLndSecretAndSaveThem()
	watchLndSecretModifiedAndSaveThem()
	fmt.Println("[Sidecar backend] Complete -> go to sleep")
    select {}
}


// Get the lnd credentials and Store them
func getLndSecretAndSaveThem() {
	fmt.Println("[Sidecar backend] getLndSecretAndSaveThem...")
	callback := &sidecar.Callback {
	    Context : map[string]interface{} {
			"namespace" : namespace,
			"secretName" : secretlnd,
			"mountedVolume" : mountedBackendStorage,
		},
		Chain : []func(*sidecar.Callback) error {
			sidecar.ReadSecret,
			sidecar.MountedVolumeToBasePath,
			sidecar.SecretToPath,
		},
	}
	if err := callback.CallNext(); err != nil {
		fmt.Println("[Sidecar backend] getLndSecretAndSaveThem... Failed !")
		panic("error on getLndSecretAndSaveThem : %s" + err.Error())
	}
	fmt.Println("[Sidecar backend] getLndSecretAndSaveThem... Done !")
}



// Get the lnd credentials and Store them
func watchLndSecretModifiedAndSaveThem() {
	fmt.Println("[Sidecar backend] watchLndSecretModifiedAndSaveThem...")
	callback := &sidecar.Callback {
	    Context : map[string]interface{} {
			"namespace" : namespace,
			"secretName" : secretlnd,
			"mountedVolume" : mountedBackendStorage,
		},
		Chain : []func(*sidecar.Callback) error {
			sidecar.WatchSecretModified,
			sidecar.MountedVolumeToBasePath,
			sidecar.SecretToPath,
		},
	}
	if err := callback.CallNext(); err != nil {
		fmt.Println("[Sidecar backend] watchLndSecretModifiedAndSaveThem... Failed !")
		panic("error on watchLndSecretModifiedAndSaveThem : %s" + err.Error())
	}
	fmt.Println("[Sidecar backend] watchLndSecretModifiedAndSaveThem... Done !")
}

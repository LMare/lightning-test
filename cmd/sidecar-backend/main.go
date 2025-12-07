package main

import (
	sidecar "github.com/Lmare/lightning-playground/sidecar"
)


const namespace = "lightning-playground"
const secretlnd = "lnd-credentials"
const mountedBackendStorage = "/app/nodes-storage"

func main() {
	getLndSecretAndSaveThem()
	watchLndSecretModifiedAndSaveThem()
}


// Get the lnd credentials and Store them
func getLndSecretAndSaveThem() {
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
		panic("error on getLndSecretAndSaveThem : %s" + err.Error())
	}
}



// Get the lnd credentials and Store them
func watchLndSecretModifiedAndSaveThem() {
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
		panic("error on watchLndSecretModifiedAndSaveThem : %s" + err.Error())
	}
}

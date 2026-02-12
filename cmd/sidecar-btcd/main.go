package main

import (
	"fmt"
	sidecar "github.com/Lmare/lightning-playground/sidecar"
)



const namespace = "lightning-playground"
const secretBtcd = "btcd-credentials"
const mountedBtcdVol = "/root/.btcd/"
const certName = "rpc.cert"

func main() {
	// it's more a initContainers than a sidecar because there isn't rootine
	// to stay alive
	fmt.Println("[Sidecar btcd] Starting")
	readBtcdCertAsSecretAndPatch()
	fmt.Println("[Sidecar btcd] Complete")
}


// Read the btcd cert and share it
func readBtcdCertAsSecretAndPatch() {
	fmt.Println("[Sidecar btcd] readBtcdCertAsSecretAndPatch...")
	callback := &sidecar.Callback {
	    Context : map[string]interface{} {
			"namespace" : namespace,
			"secretName" : secretBtcd,
			"mountedVolume" : mountedBtcdVol,
			"filePathInVolume" : certName,
		},
		Chain : []func(*sidecar.Callback) error {
			sidecar.ReadFileAsSecret,
			sidecar.PatchSecret,
		},
	}
	if err := callback.CallNext(); err != nil {
		fmt.Println("[Sidecar btcd] readBtcdCertAsSecretAndPatch... Failed !")
		panic("error on readBtcdCertAsSecretAndPatch : %s" + err.Error())
	}
	fmt.Println("[Sidecar btcd] readBtcdCertAsSecretAndPatch... Done !")
}

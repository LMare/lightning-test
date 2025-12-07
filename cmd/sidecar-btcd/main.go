package main

import (
	sidecar "github.com/Lmare/lightning-playground/sidecar"
)



const namespace = "lightning-playground"
const secretBtcd = "btcd-credentials"
const mountedBtcdVol = "/root/.btcd/"
const certName = "tls.cert"

func main() {
	readBtcdCertAsSecretAndPatch()
}


// Read the btcd cert and share it
func readBtcdCertAsSecretAndPatch() {
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
		panic("error on readBtcdCertAsSecretAndPatch : %s" + err.Error())
	}
}

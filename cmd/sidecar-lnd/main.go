package main

import (
	"os"

	sidecar "github.com/Lmare/lightning-playground/sidecar"
)


const namespace = "lightning-playground"
const secretBtcd = "btcd-credentials"
const secretLnd = "lnd-credentials"
const mountedLndVol = "/root/.lnd/"
const mountedBtcdVol = "/root/.btcd/"
const macarronPattern = "*.macaroon"
const certName = "tls.cert"
var lndNodeName = os.Getenv("NODE_NAME")

func main() {
	pullBtcdSecret()
	readLndCertAsSecretAndPatch()
	watchLndMacaroonAsSecretAndPatch()
}


// Read the lnd cert and share it
func readLndCertAsSecretAndPatch() {
	callback := &sidecar.Callback {
	    Context : map[string]interface{} {
			"namespace" : namespace,
			"secretName" : secretLnd,
			"mountedVolume" : mountedLndVol,
			"filePathInVolume" : certName,
			"secretPrefix" : lndNodeName,
		},
		Chain : []func(*sidecar.Callback) error {
			sidecar.ReadFileAsSecret,
			sidecar.PatchSecret,
		},
	}
	if err := callback.CallNext(); err != nil {
		panic("error on readLndCertAsSecretAndPatch : %s" + err.Error())
	}
}


// Get the btcd cert and Store it
func pullBtcdSecret() {
	callback := &sidecar.Callback {
	    Context : map[string]interface{} {
			"namespace" : namespace,
			"secretName" : secretBtcd,
			"mountedVolume" : mountedBtcdVol,
		},
		Chain : []func(*sidecar.Callback) error {
			sidecar.ReadSecret,
			sidecar.MountedVolumeToBasePath,
			sidecar.SecretToPath,
		},
	}
	if err := callback.CallNext(); err != nil {
		panic("error on pullBtcdSecret : %s" + err.Error())
	}
}


// Check Macaroon File Update & send it as secret
func watchLndMacaroonAsSecretAndPatch() {
	callback := &sidecar.Callback {
	    Context : map[string]interface{} {
			"namespace" : namespace,
			"secretName" : secretLnd,
			"mountedVolume" : mountedLndVol,
			"pattern" : macarronPattern,
			"secretPrefix" : lndNodeName,
		},
		Chain : []func(*sidecar.Callback) error {
			sidecar.MountedVolumeToBasePath,
			sidecar.WatchFilePattern,
			sidecar.IfIsEventOfWriting,
			sidecar.FilePathToFilePathInVolume,
			sidecar.ReadFileAsSecret,
			sidecar.PatchSecret,
		},
	}
	if err := callback.CallNext(); err != nil {
		panic("error on watchLndMacaroonAsSecretAndPatch : " + err.Error())
	}
}

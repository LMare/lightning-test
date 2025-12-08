package main

import (
	"os"
	"fmt"
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
	fmt.Println("[Sidecar lnd] Starting on " + lndNodeName)
	pullBtcdSecret()
	readLndCertAsSecretAndPatch()
	watchLndMacaroonAsSecretAndPatch()
	fmt.Println("[Sidecar lnd] Complete -> go to sleep")
	select {}
}


// Read the lnd cert and share it
func readLndCertAsSecretAndPatch() {
	fmt.Println("[Sidecar lnd] readLndCertAsSecretAndPatch...")
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
		fmt.Println("[Sidecar lnd] readLndCertAsSecretAndPatch... Failed !")
		panic("error on readLndCertAsSecretAndPatch : %s" + err.Error())
	}
	fmt.Println("[Sidecar lnd] readLndCertAsSecretAndPatch... Done !")
}


// Get the btcd cert and Store it
func pullBtcdSecret() {
	fmt.Println("[Sidecar lnd] pullBtcdSecret...")
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
		fmt.Println("[Sidecar lnd] pullBtcdSecret... Failed !")
		panic("error on pullBtcdSecret : %s" + err.Error())
	}
	fmt.Println("[Sidecar lnd] pullBtcdSecret... Done !")
}


// Check Macaroon File Update & send it as secret
func watchLndMacaroonAsSecretAndPatch() {
	fmt.Println("[Sidecar lnd] watchLndMacaroonAsSecretAndPatch...")
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
		fmt.Println("[Sidecar lnd] watchLndMacaroonAsSecretAndPatch... Failed !")
		panic("error on watchLndMacaroonAsSecretAndPatch : " + err.Error())
	}
	fmt.Println("[Sidecar lnd] watchLndMacaroonAsSecretAndPatch... Done !")
}

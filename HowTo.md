#conf local
export PATH="$PATH:$HOME/Perso/Outils/go1.25.0.windows-amd64/go/bin"
export PATH="$PATH:$HOME/go/bin/"
export PATH="$PATH:$HOME/Perso/Outils/gprotoc-32.0-win64/bin"


#Exécuter localement
go run ./cmd/backend 
go run ./cmd/frontend



#Pour générer les Clients gRPC
1. installation de protoc : https://protobuf.dev/installation/ 
2. installation des plugins go pour protoc 
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
3. télécharger des fichier *.proto ex : https://github.com/lightningnetwork/lnd/tree/master/lnrpc
4. générer les clients go : protoc --go_out=. --go-grpc_out=. lightning.proto 


#Inclure une lib externe : 
go get github.com/joho/godotenv


#Get docker image  from btcsuite/btcd (no official image in dockerHub) :
git clone --depth 1 https://github.com/btcsuite/btcd.git
sudo docker build -t btcsuite/btcd:latest .

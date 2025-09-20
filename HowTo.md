#conf local
export PATH="$PATH:$HOME/Perso/Outils/go1.25.0.windows-amd64/go/bin"
export PATH="$PATH:$HOME/go/bin/"
export PATH="$PATH:$HOME/Perso/Outils/gprotoc-32.0-win64/bin"


#Exécuter localement
go run ./cmd/backend
go run ./cmd/frontend



#Generate gRPC  clients
1. installation de protoc : https://protobuf.dev/installation/
2. installation des plugins go pour protoc
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
3. télécharger des fichier *.proto ex : https://github.com/lightningnetwork/lnd/tree/master/lnrpc
4. générer les clients go : protoc --go_out=. --go-grpc_out=. lightning.proto


#Inclure external lib :
go get github.com/joho/godotenv


#Get docker image  from btcsuite/btcd (no official image in dockerHub) :
git clone --depth 1 https://github.com/btcsuite/btcd.git
docker build -t btcsuite/btcd:latest .


#start container docker
docker-compose up -d

#connexion in the container
docker exec -it lnd1 bash

#init the wallet ln to generate macaroons then unlock it
lncli --network=simnet create
lncli --network=simnet unlock

##Copy macaroons + certificate from the lnd container :
# Certificat TLS
docker cp lnd1:/root/.lnd/tls.cert .
# Macaroons
docker cp lnd1:/root/.lnd/data/chain/bitcoin/simnet/admin.macaroon .
docker cp lnd1:/root/.lnd/data/chain/bitcoin/simnet/invoice.macaroon .
docker cp lnd1:/root/.lnd/data/chain/bitcoin/simnet/readonly.macaroon .

#Activate taproot by gererating some blocks in simnet
btcctl --simnet generate 1500

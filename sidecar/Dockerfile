### Utiliser comme ceci pour générer l'image
### docker build -t lightning-playground-backend:latest -f backend/Dockerfile .

ARG GO_TAG=1.25.0

# 1. Compilation Go du sidecar en fonction type demandé
FROM golang:$GO_TAG AS compiler
ARG TYPE_SIDECAR=lnd
ENV TYPE_SIDECAR=$TYPE_SIDECAR
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY sidecar sidecar/
COPY cmd/sidecar-backend cmd/sidecar-backend/
COPY cmd/sidecar-btcd cmd/sidecar-btcd/
COPY cmd/sidecar-lnd cmd/sidecar-lnd/
RUN CGO_ENABLED=0 go build -o sidecar-$TYPE_SIDECAR ./cmd/sidecar-$TYPE_SIDECAR

# 2. Construction de l'image finale
FROM scratch
ARG TYPE_SIDECAR=lnd
ENV TYPE_SIDECAR=$TYPE_SIDECAR
WORKDIR /app
COPY --from=compiler /app/sidecar-$TYPE_SIDECAR ./sidecar
ENTRYPOINT ["/app/sidecar"]

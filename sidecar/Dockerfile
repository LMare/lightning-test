### Utiliser comme ceci pour générer l'image
### docker build -t lightning-playground-backend:latest -f backend/Dockerfile .

ARG GO_TAG=1.25.0

# 1. Compilation Go du sidecar en fonction type demandé
FROM golang:$GO_TAG AS compiler
WORKDIR /app
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod,id=gomod,sharing=locked \
	go mod download
COPY sidecar sidecar/
COPY cmd cmd
RUN --mount=type=cache,target=/root/.cache/go-build,id=gobuild-sidecar,sharing=locked \
    --mount=type=cache,target=/go/pkg/mod,id=gomod,sharing=locked \
    CGO_ENABLED=0 \
    go build -o sidecar-lnd ./cmd/sidecar-lnd && \
    go build -o sidecar-btcd ./cmd/sidecar-btcd && \
	go build -o sidecar-backend ./cmd/sidecar-backend

# 2. Construction de l'image finale
FROM scratch
ARG TYPE_SIDECAR=lnd
ENV TYPE_SIDECAR=$TYPE_SIDECAR
WORKDIR /app
COPY --from=compiler /app/sidecar-$TYPE_SIDECAR ./sidecar
ENTRYPOINT ["/app/sidecar"]

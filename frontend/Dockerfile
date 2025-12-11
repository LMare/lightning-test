### Utiliser comme ceci pour générer l'image
### docker build -t lightning-playground-frontend:latest -f frontend/Dockerfile .

ARG ALPINE_TAG=3.22
ARG GO_TAG=1.25.0

# 1. Build de l'application frontend
FROM golang:$GO_TAG AS compiler
WORKDIR /app
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod,id=gomod,sharing=locked \
	go mod download
COPY *.go ./
COPY cmd/frontend cmd/frontend/
RUN --mount=type=cache,target=/root/.cache/go-build,id=gobuild-frontend,sharing=locked \
    --mount=type=cache,target=/go/pkg/mod,id=gomod,sharing=locked \
	CGO_ENABLED=0 \
	go build -o app-frontend ./cmd/frontend; \

# 2. Construction du dossier /app avec le binaire + ressources
FROM scratch AS builder
WORKDIR /app
COPY --from=compiler /app/app-frontend ./
COPY .env ./
COPY frontend ./frontend/

# 3. Construction de l'image finale
FROM scratch AS frontend-scratch
WORKDIR /app
COPY --from=builder /app ./
ENTRYPOINT ["/app/app-frontend"]

FROM alpine:$ALPINE_TAG AS frontend-alpine
WORKDIR /app
COPY --from=builder /app ./
ENTRYPOINT ["/app/app-frontend"]

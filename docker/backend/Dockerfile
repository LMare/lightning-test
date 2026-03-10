### Utiliser comme ceci pour générer l'image
### docker build -t lightning-playground-backend:latest -f backend/Dockerfile .

ARG ALPINE_TAG=3.22
ARG GO_TAG=1.25.0

# 1. Source Go l'application backend
FROM golang:$GO_TAG AS compiler
WORKDIR /app
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod,id=gomod,sharing=locked \
	go mod download
COPY *.go ./
COPY cmd/backend cmd/backend/
COPY backend backend/
# compilation statique (dépendances des libs C)
RUN --mount=type=cache,target=/root/.cache/go-build,id=gobuild-backend,sharing=locked \
	--mount=type=cache,target=/go/pkg/mod,id=gomod,sharing=locked \
	CGO_ENABLED=0 \
	go build -o app-backend ./cmd/backend

# 2. Construction du dossier /app avec le binaire + ressources
FROM scratch AS builder
WORKDIR /app
COPY --from=compiler /app/app-backend ./
COPY .env ./
COPY backend/templates ./backend/templates/
COPY migrations migrations/

# 3. Construction de l'image finale (alpine ou scratch)

FROM scratch AS backend-scratch
WORKDIR /app
COPY --from=builder /app ./
ENTRYPOINT ["/app/app-backend"]

FROM alpine:$ALPINE_TAG AS backend-alpine
WORKDIR /app
COPY --from=builder /app ./
ENTRYPOINT ["/app/app-backend"]

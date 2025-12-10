### Utiliser comme ceci pour générer l'image
### docker build -t lightning-playground-frontend:latest -f frontend/Dockerfile .

ARG ALPINE_TAG=3.22
ARG GO_TAG=1.25.0

# 1. Build de l'application frontend
FROM golang:$GO_TAG AS compiler
ARG COMPILATION=static
ENV COMPILATION=$COMPILATION
WORKDIR /app
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
	go mod download
COPY *.go ./
COPY cmd/frontend cmd/frontend/
RUN --mount=type=cache,target=/root/.cache/go-build \
	if [ "$COMPILATION" = "static" ]; then \
		echo "CGO_ENABLED=0 go build -o app-frontend ./cmd/frontend"; \
		CGO_ENABLED=0 go build -o app-frontend ./cmd/frontend; \
	else \
		echo "go build -o app-frontend ./cmd/frontend"; \
		go build -o app-frontend ./cmd/frontend; \
	fi

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

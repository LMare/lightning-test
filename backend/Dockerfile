### Utiliser comme ceci pour générer l'image
### docker build -t lightning-playground-backend:latest -f backend/Dockerfile .

ARG ALPINE_TAG=3.22
ARG GO_TAG=1.25.0

# 1. Source Go l'application backend
FROM golang:$GO_TAG AS compiler
ARG COMPILATION=static
ENV COMPILATION=$COMPILATION
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
COPY cmd/backend cmd/backend/
COPY backend backend/
# compilation statique / dynamique (dépendances des libs C)
RUN if [ "$COMPILATION" = "static" ]; then \
		echo "CGO_ENABLED=0 go build -o app-backend ./cmd/backend"; \
		CGO_ENABLED=0 go build -o app-backend ./cmd/backend; \
	else \
		echo "go build -o app-backend ./cmd/backend"; \
		go build -o app-backend ./cmd/backend; \
	fi

# 2. Construction du dossier /app avec le binaire + ressources
FROM scratch AS builder
WORKDIR /app
COPY --from=compiler /app/app-backend ./
COPY .env nodes.yaml ./
COPY backend/templates ./backend/templates/

# 3. Construction de l'image finale (alpine ou scratch)

FROM scratch AS backend-scratch
WORKDIR /app
COPY --from=builder /app ./
ENTRYPOINT ["/app/app-backend"]

FROM alpine:$ALPINE_TAG AS backend-alpine
WORKDIR /app
COPY --from=builder /app ./
ENTRYPOINT ["/app/app-backend"]

variable "APP_VERSION" {
  default = "v0.4.0-0"
}
variable "BTCD_VERSION" {
  default = "v0.25.0"
}

variable "LND_TAG" {
  default = "v0.20.0-beta-custom"
}

variable "ALPINE_TAG" {
  default = "3.22.2"
}
variable "GO_TAG" {
  default = "1.25.4-alpine3.22"
}

target "common-app-args" {
	args = {
		ALPINE_TAG = "${ALPINE_TAG}"
		GO_TAG = "${GO_TAG}"
	}
}

group "compiler" {
	targets = ["sidecar-compiler", "backend-compiler", "frontend-compiler"]
}

group "app-scratch" {
	targets = ["backend-scratch", "frontend-scratch"]
}
group "app-alpine" {
	targets = ["backend-alpine", "frontend-alpine"]
}
group "backend" {
	targets = ["backend-alpine", "backend-scratch"]
}
group "frontend" {
	targets = ["frontend-alpine", "frontend-scratch"]
}

group "sidecars" {
	targets = ["sidecar-backend", "sidecar-lnd", "sidecar-btcd"]
}

group "data" {
  targets = ["duckdb-elt-service", "prefect-agent"]
}

group "default" {
	targets = ["frontend", "backend", "btcd", "lnd", "sidecars", "data"]
}


target "backend-compiler" {
  context = "."
  dockerfile = "./docker/backend/Dockerfile"
  target = "compiler"
  inherits = ["common-app-args"]
  tags = ["backend-compiler:latest"]
}


target "backend-scratch" {
  context = "."
  dockerfile = "./docker/backend/Dockerfile"
  target = "backend-scratch"
  inherits = ["common-app-args"]
  tags = ["LMare/lightning-playground-backend:${APP_VERSION}"]
}

target "backend-alpine" {
  context = "."
  dockerfile = "./docker/backend/Dockerfile"
  target = "backend-alpine"
  inherits = ["common-app-args"]
  tags = ["LMare/lightning-playground-backend:${APP_VERSION}-alpine-${ALPINE_TAG}"]
}

target "frontend-compiler" {
  context = "."
  dockerfile = "./docker/frontend/Dockerfile"
  target = "compiler"
  inherits = ["common-app-args"]
  tags = ["frontend-compiler:latest"]
}

target "frontend-scratch" {
  context = "."
  dockerfile = "./docker/frontend/Dockerfile"
  target = "frontend-scratch"
  inherits = ["common-app-args"]
  tags = ["LMare/lightning-playground-frontend:${APP_VERSION}"]
}

target "frontend-alpine" {
  context = "."
  dockerfile = "./docker/frontend/Dockerfile"
  target = "frontend-alpine"
  inherits = ["common-app-args"]
  tags = ["LMare/lightning-playground-frontend:${APP_VERSION}-alpine-${ALPINE_TAG}"]
}

target "btcd" {
  context = "https://github.com/btcsuite/btcd.git#${BTCD_VERSION}"
  dockerfile = "Dockerfile"
  tags = ["btcsuite/btcd:${BTCD_VERSION}"]
}

target "lnd" {
  context = "https://github.com/LMare/lnd.git#feature/gRPC-alias-color"
  dockerfile = "Dockerfile"
  args = {
	  checkout = "feature/gRPC-alias-color"
	  git_url = "https://github.com/LMare/lnd"
  }
  tags = ["LMare/lnd:${LND_TAG}"]
}

target "sidecar-backend" {
  context = "."
  dockerfile = "./docker/sidecar/Dockerfile"
  args = { TYPE_SIDECAR = "backend" }
  inherits = ["common-app-args"]
  tags = ["LMare/lightning-playground-sidecar-backend:${APP_VERSION}"]
}

target "sidecar-compiler" {
  context = "."
  dockerfile = "./docker/sidecar/Dockerfile"
  inherits = ["common-app-args"]
  target = "compiler"
  tags = ["sidecar-compiler:latest"]
}


target "sidecar-lnd" {
  context = "."
  dockerfile = "./docker/sidecar/Dockerfile"
  args = { TYPE_SIDECAR = "lnd" }
  inherits = ["common-app-args"]
  tags = ["LMare/lightning-playground-sidecar-lnd:${APP_VERSION}"]
}


target "sidecar-btcd" {
  context = "."
  dockerfile = "./docker/sidecar/Dockerfile"
  args = { TYPE_SIDECAR = "btcd" }
  inherits = ["common-app-args"]
  tags = ["LMare/lightning-playground-sidecar-btcd:${APP_VERSION}"]
}


target "prefect-agent" {
  context = "./docker/prefect-agent/prefect"
  dockerfile = "../Dockerfile"
  tags = ["LMare/prefect-agent:${APP_VERSION}"]
}

target "duckdb-elt-service" {
  context = "./docker/duckdb-elt-service"
  dockerfile = "Dockerfile"
  tags = ["LMare/duckdb-elt-service:${APP_VERSION}"]
}

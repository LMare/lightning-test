variable "APP_VERSION" {
  default = "v0.2.0-5"
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
group "default" {
	targets = ["frontend", "backend", "btcd", "lnd"]
}


target "backend-scratch" {
  context = "."
  dockerfile = "./backend/Dockerfile"
  target = "backend-scratch"
  args = { COMPILATION = "static" }
  inherits = ["common-app-args"]
  tags = ["LMare/lightning-playground-backend:${APP_VERSION}"]
}

target "backend-alpine" {
  context = "."
  dockerfile = "./backend/Dockerfile"
  target = "backend-alpine"
  args = { COMPILATION = "dynamic" }
  inherits = ["common-app-args"]
  tags = ["LMare/lightning-playground-backend:${APP_VERSION}-alpine-${ALPINE_TAG}"]
}

target "frontend-scratch" {
  context = "."
  dockerfile = "./frontend/Dockerfile"
  target = "frontend-scratch"
  args = { COMPILATION = "static" }
  inherits = ["common-app-args"]
  tags = ["LMare/lightning-playground-frontend:${APP_VERSION}"]
}

target "frontend-alpine" {
  context = "."
  dockerfile = "./frontend/Dockerfile"
  target = "frontend-alpine"
  args = { COMPILATION = "dynamic" }
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
  dockerfile = "./sidecar/Dockerfile"
  args = { TYPE_SIDECAR = "backend" }
  inherits = ["common-app-args"]
  tags = ["LMare/lightning-playground-sidecar-backend:${APP_VERSION}"]
}

target "sidecar-lnd" {
  context = "."
  dockerfile = "./sidecar/Dockerfile"
  args = { TYPE_SIDECAR = "lnd" }
  inherits = ["common-app-args"]
  tags = ["LMare/lightning-playground-sidecar-lnd:${APP_VERSION}"]
}

target "sidecar-btcd" {
  context = "."
  dockerfile = "./sidecar/Dockerfile"
  args = { TYPE_SIDECAR = "btcd" }
  inherits = ["common-app-args"]
  tags = ["LMare/lightning-playground-sidecar-btcd:${APP_VERSION}"]
}

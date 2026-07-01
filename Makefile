.PHONY: build test clean docker-build docker-run

FRONTEND_DIR := frontend
BACKEND_DIR := backend
ASSETS_DIR := $(BACKEND_DIR)/internal/assets/dist
IMAGE_TAG := calendar-booking:latest
PORT ?= 8080

build:
	@echo "==> Building frontend"
	npm ci --prefix $(FRONTEND_DIR) --legacy-peer-deps
	npm run build --prefix $(FRONTEND_DIR)

	@echo "==> Copying frontend assets into backend embed directory"
	rm -rf $(ASSETS_DIR)
	mkdir -p $(ASSETS_DIR)
	cp -R $(FRONTEND_DIR)/dist/. $(ASSETS_DIR)/
	mkdir -p $(ASSETS_DIR)/assets
	touch $(ASSETS_DIR)/assets/.gitkeep

	@echo "==> Building backend binary"
	make -C $(BACKEND_DIR) build

test:
	@echo "==> Running backend tests"
	make -C $(BACKEND_DIR) test

	@echo "==> Running frontend tests"
	npm test --prefix $(FRONTEND_DIR)

clean:
	rm -rf $(ASSETS_DIR)
	make -C $(BACKEND_DIR) clean

docker-build:
	docker build -t $(IMAGE_TAG) .

docker-run:
	docker run --rm -p $(PORT):$(PORT) -e PORT=$(PORT) $(IMAGE_TAG)

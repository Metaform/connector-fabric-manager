.PHONY: help build test clean docker-build docker-clean install-tools generate-mocks

#==============================================================================
# Multi-Service Orchestration - Delegates to Individual Makefiles
#==============================================================================

# Service directories
COMMON_DIR=common
PMANAGER_DIR=pmanager
TMANAGER_DIR=tmanager
EDCV_DIR=agent/edcv

E2E_DIR=e2e

# Docker settings
DOCKER_REGISTRY=
DOCKER_TAG=latest

#==============================================================================
# Help
#==============================================================================

help:
	@echo "Multi-Service Orchestration Commands (delegates to individual Makefiles):"
	@echo ""
	@echo "Build Commands:"
	@echo "  build          - Build all services"
	@echo "  build-pmanager - Build pmanager service"
	@echo "  build-tmanager - Build tmanager service"
	@echo "  build-all      - Build all services for all platforms"
	@echo ""
	@echo "Test Commands:"
	@echo "  test           - Run tests for all services"
	@echo "  test-pmanager  - Test pmanager service"
	@echo "  test-tmanager  - Test tmanager service"
	@echo ""
	@echo "Development Commands:"
	@echo "  dev-pmanager   - Run pmanager in development mode"
	@echo "  dev-tmanager   - Run tmanager in development mode"
	@echo "  clean          - Clean all build artifacts"
	@echo ""
	@echo "Docker Commands:"
	@echo "  docker-build   - Build all Docker images"
	@echo "  docker-clean   - Remove all Docker images"
	@echo ""
	@echo "Tool Commands:"
	@echo "  install-tools  - Install development tools for all services"
	@echo "  generate-mocks - Generate mocks for all services"

#==============================================================================
# Build Commands - Delegate to Service Makefiles
#==============================================================================

build:
	@echo "Building all modules..."
	$(MAKE) -C $(PMANAGER_DIR) build
	$(MAKE) -C $(TMANAGER_DIR) build
	$(MAKE) -C $(EDCV_DIR) build

build-pmanager:
	@echo "Building pmanager..."
	$(MAKE) -C $(PMANAGER_DIR) build

build-tmanager:
	@echo "Building tmanager..."
	$(MAKE) -C $(TMANAGER_DIR) build

build-edcv:
	@echo "Building tmanager..."
	$(MAKE) -C $(EDCV_DIR) build

build-all:
	@echo "Building all modules for all platforms..."
	$(MAKE) -C $(PMANAGER_DIR) build-all
	$(MAKE) -C $(TMANAGER_DIR) build-all
	$(MAKE) -C $(EDCV_DIR) build-all

#==============================================================================
# Test Commands - Delegate to Service Makefiles
#==============================================================================

test:
	@echo "Testing all services..."
	$(MAKE) -C $(COMMON_DIR) test
	$(MAKE) -C $(PMANAGER_DIR) test
	$(MAKE) -C $(TMANAGER_DIR) test
	$(MAKE) -C $(EDCV_DIR) test
	$(MAKE) -C $(E2E_DIR) test

test-common:
	@echo "Testing common..."
	$(MAKE) -C $(COMMON_DIR) test

test-pmanager:
	@echo "Testing pmanager..."
	$(MAKE) -C $(PMANAGER_DIR) test

test-tmanager:
	@echo "Testing tmanager..."
	$(MAKE) -C $(TMANAGER_DIR) test

test-edcv:
	@echo "Testing EDC-V agent..."
	$(MAKE) -C $(EDCV_DIR) test

#==============================================================================
# Development Commands - Delegate to Service Makefiles
#==============================================================================

dev-pmanager:
	@echo "Starting pmanager in development mode..."
	$(MAKE) -C $(PMANAGER_DIR) dev-server

dev-tmanager:
	@echo "Starting tmanager in development mode..."
	$(MAKE) -C $(TMANAGER_DIR) dev-server

clean:
	@echo "Cleaning all services..."
	$(MAKE) -C $(COMMON_DIR) clean
	$(MAKE) -C $(PMANAGER_DIR) clean
	$(MAKE) -C $(TMANAGER_DIR) clean
	$(MAKE) -C $(EDCV_DIR) clean

#==============================================================================
# Tool Commands - Delegate to Service Makefiles
#==============================================================================

install-tools:
	@echo "Installing tools for all services..."
	$(MAKE) -C $(PMANAGER_DIR) install-tools
	$(MAKE) -C $(TMANAGER_DIR) install-tools

generate-mocks:
	@echo "Generating mocks for all services..."
	$(MAKE) -C $(COMMON_DIR) generate-mocks
	$(MAKE) -C $(PMANAGER_DIR) generate-mocks

generate-docs:
	$(MAKE) -C $(TMANAGER_DIR) generate-docs
	$(MAKE) -C $(PMANAGER_DIR) generate-docs

#==============================================================================
# Docker Commands - Handled at Top Level
#==============================================================================

docker-build: docker-build-pmanager docker-build-tmanager docker-build-testagent

docker-build-pmanager:
	@echo "Building pmanager Docker image..."
	docker build -f Dockerfile.pmanager.dockerfile -t $(DOCKER_REGISTRY)pmanager:$(DOCKER_TAG) .

docker-build-tmanager:
	@echo "Building tmanager Docker image..."
	docker build -f Dockerfile.tmanager.dockerfile -t $(DOCKER_REGISTRY)tmanager:$(DOCKER_TAG) .

docker-build-testagent:
	@echo "Building test agent Docker image..."
	docker build -f Dockerfile.testagent.dockerfile -t $(DOCKER_REGISTRY)testagent:$(DOCKER_TAG) .

docker-clean: docker-clean-pmanager docker-clean-tmanager docker-clean-testagent

docker-clean-pmanager:
	docker rmi $(DOCKER_REGISTRY)pmanager:$(DOCKER_TAG) || true

docker-clean-tmanager:
	docker rmi $(DOCKER_REGISTRY)tmanager:$(DOCKER_TAG) || true

docker-clean-testagent:
	docker rmi $(DOCKER_REGISTRY)testagent:$(DOCKER_TAG) || true

#==============================================================================
# Combined Commands
#==============================================================================

all: build docker-build
	@echo "Built all services and Docker images"

deploy: build-all docker-build
	@echo "Built all services for all platforms and Docker images"

dev-setup: install-tools generate-mocks build
	@echo "Development environment ready for all services"

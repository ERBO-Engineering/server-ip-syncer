.PHONY: build clean install remote-install remote-uninstall

# Build configuration
BINARY_NAME=ip-leaker
BUILD_DIR=build
GOOS?=linux
GOARCH?=amd64

# SSH configuration
SSH_HOST?=erbo
SSH_USER?=bob
SSH_PORT?=22
REMOTE_DIR?=/home/$(SSH_USER)/ip-leaker

# Build the binary
build: clean
	mkdir -p $(BUILD_DIR)
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/ip-leaker

# Clean build directory
clean:
	rm -rf $(BUILD_DIR)

# Install locally (for testing)
install: build
	sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	sudo cp ip-leaker.service /etc/systemd/system/
	sudo sed -i "s/%i/$(shell whoami)/g" /etc/systemd/system/ip-leaker.service
	sudo systemctl daemon-reload
	sudo systemctl enable ip-leaker
	sudo systemctl start ip-leaker

# Uninstall locally (for testing)
uninstall:
	sudo systemctl stop ip-leaker
	sudo systemctl disable ip-leaker
	sudo rm /etc/systemd/system/ip-leaker.service
	sudo rm /usr/local/bin/$(BINARY_NAME)
	sudo systemctl daemon-reload

# Remote installation
remote-install: build
	@echo "Creating remote directory..."
	ssh -t -p $(SSH_PORT) $(SSH_USER)@$(SSH_HOST) "mkdir -p $(REMOTE_DIR)"
	
	@echo "Copying files..."
	scp -P $(SSH_PORT) $(BUILD_DIR)/$(BINARY_NAME) $(SSH_USER)@$(SSH_HOST):$(REMOTE_DIR)/
	scp -P $(SSH_PORT) ip-leaker.service $(SSH_USER)@$(SSH_HOST):$(REMOTE_DIR)/
	scp -P $(SSH_PORT) key.json $(SSH_USER)@$(SSH_HOST):$(REMOTE_DIR)/
	
	@echo "Installing service..."
	@echo "Enter your sudo password when prompted:"
	ssh -t -p $(SSH_PORT) $(SSH_USER)@$(SSH_HOST) "sudo cp $(REMOTE_DIR)/$(BINARY_NAME) /usr/local/bin/ && \
		sudo cp $(REMOTE_DIR)/ip-leaker.service /etc/systemd/system/ && \
		sudo sed -i 's/%i/$(SSH_USER)/g' /etc/systemd/system/ip-leaker.service && \
		sudo systemctl daemon-reload && \
		sudo systemctl enable ip-leaker && \
		sudo systemctl start ip-leaker"

# Remote uninstallation
remote-uninstall:
	@echo "Uninstalling service..."
	@echo "Enter your sudo password when prompted:"
	ssh -t -p $(SSH_PORT) $(SSH_USER)@$(SSH_HOST) "sudo systemctl stop ip-leaker && \
		sudo systemctl disable ip-leaker && \
		sudo rm /etc/systemd/system/ip-leaker.service && \
		sudo rm /usr/local/bin/$(BINARY_NAME) && \
		sudo systemctl daemon-reload && \
		rm -rf $(REMOTE_DIR)"

# Remote status check
remote-status:
	@echo "Checking service status..."
	@echo "Enter your sudo password when prompted:"
	ssh -t -p $(SSH_PORT) $(SSH_USER)@$(SSH_HOST) "sudo systemctl status ip-leaker"

# Remote logs
remote-logs:
	@echo "Showing service logs..."
	@echo "Enter your sudo password when prompted:"
	ssh -t -p $(SSH_PORT) $(SSH_USER)@$(SSH_HOST) "sudo journalctl -u ip-leaker -f" 
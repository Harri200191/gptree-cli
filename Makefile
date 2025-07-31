APP_NAME = gptree
VERSION := 1.0
BUILD_DIR := .build
ARCHS := amd64 arm arm64 aarch64 x86_64
DEB_DIR := debuild

all: $(BUILD_DIR)/$(APP_NAME) $(BUILD_DIR)/$(APP_NAME).exe deb-packages

.PHONY: clean deb-packages

$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

# Native build (Linux/WSL)
$(BUILD_DIR)/$(APP_NAME): | $(BUILD_DIR)
	go build -o $@ .

# Windows build
$(BUILD_DIR)/$(APP_NAME).exe: | $(BUILD_DIR)
	go build -o $@ .

# Build binaries for all architectures
$(BUILD_DIR)/$(APP_NAME)_%: | $(BUILD_DIR)
	go build -o $@ .

# Build .deb packages
deb-packages: $(foreach arch,$(ARCHS),$(BUILD_DIR)/$(APP_NAME)_$(arch).deb)

$(BUILD_DIR)/$(APP_NAME)_%.deb: $(BUILD_DIR)/$(APP_NAME)_%
	@echo "📦 Packaging $@"
	rm -rf $(DEB_DIR)
	mkdir -p $(DEB_DIR)/DEBIAN
	mkdir -p $(DEB_DIR)/usr/bin

	cp $< $(DEB_DIR)/usr/bin/$(APP_NAME)
	chmod 755 $(DEB_DIR)/usr/bin/$(APP_NAME)

	echo "Package: $(APP_NAME)" > $(DEB_DIR)/DEBIAN/control
	echo "Version: $(VERSION)" >> $(DEB_DIR)/DEBIAN/control
	echo "Section: utils" >> $(DEB_DIR)/DEBIAN/control
	echo "Priority: optional" >> $(DEB_DIR)/DEBIAN/control
	echo "Architecture: $*" >> $(DEB_DIR)/DEBIAN/control
	echo "Maintainer: Haris Rehman <harisrehmanchugtai@gmail.com>" >> $(DEB_DIR)/DEBIAN/control
	echo "Description: GPT-friendly directory summarizer and prompt generator." >> $(DEB_DIR)/DEBIAN/control

	chmod -R 755 $(DEB_DIR)/DEBIAN
	dpkg-deb --build $(DEB_DIR)
	mv $(DEB_DIR).deb $@

clean:
	@echo "🧹 Cleaning..."
	rm -rf $(BUILD_DIR) $(DEB_DIR)

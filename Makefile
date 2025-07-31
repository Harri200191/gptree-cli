APP_NAME := gptree
VERSION := 1.0
BUILD_DIR := .build
DEB_DIR := debuild
ARCHS := amd64 arm arm64

all: $(BUILD_DIR)/$(APP_NAME) $(BUILD_DIR)/$(APP_NAME).exe deb-packages

.PHONY: clean deb-packages

$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

# Native Linux build
$(BUILD_DIR)/$(APP_NAME): | $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -o $@ .

# Windows build
$(BUILD_DIR)/$(APP_NAME).exe: | $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 go build -o $@ .

# Define architecture-specific binary targets explicitly
BINARIES := $(foreach arch,$(ARCHS),$(BUILD_DIR)/$(APP_NAME)-bin_$(arch))

# Rule for building architecture-specific binaries
$(BINARIES): $(BUILD_DIR)/$(APP_NAME)-bin_%: | $(BUILD_DIR)
	GOOS=linux GOARCH=$* go build -o $@ .

# Define .deb package targets
DEB_PACKAGES := $(foreach arch,$(ARCHS),$(BUILD_DIR)/$(APP_NAME)_$(arch).deb)

# Target to build all .deb packages
deb-packages: $(DEB_PACKAGES)

# Rule to create .deb packages from binaries
$(BUILD_DIR)/$(APP_NAME)_%.deb: $(BUILD_DIR)/$(APP_NAME)-bin_%
	@echo "📦 Packaging $@"
	rm -rf $(DEB_DIR)
	mkdir -p $(DEB_DIR)/DEBIAN
	mkdir -p $(DEB_DIR)/usr/bin
	chmod 755 $(DEB_DIR)/DEBIAN
	cp $< $(DEB_DIR)/usr/bin/$(APP_NAME)
	chmod 755 $(DEB_DIR)/usr/bin/$(APP_NAME)
	echo "Package: $(APP_NAME)" > $(DEB_DIR)/DEBIAN/control
	echo "Version: $(VERSION)" >> $(DEB_DIR)/DEBIAN/control
	echo "Section: utils" >> $(DEB_DIR)/DEBIAN/control
	echo "Priority: optional" >> $(DEB_DIR)/DEBIAN/control
	echo "Architecture: $*" >> $(DEB_DIR)/DEBIAN/control
	echo "Maintainer: Haris Rehman <harisrehmanchugtai@gmail.com>" >> $(DEB_DIR)/DEBIAN/control
	echo "Description: GPT-friendly directory summarizer and prompt generator." >> $(DEB_DIR)/DEBIAN/control
	chmod 644 $(DEB_DIR)/DEBIAN/control
	dpkg-deb --build $(DEB_DIR)
	mv $(DEB_DIR).deb $@

clean:
	@echo "🧹 Cleaning..."
	rm -rf $(BUILD_DIR) $(DEB_DIR)

.RECIPEPREFIX := +
.PHONY: build install lint clean install-magick

export CGO_CFLAGS_ALLOW = -Xpreprocessor

build: install-magick
	+go build -o wallify .

install:
	+go install .

lint:
	+golangci-lint run

clean:
	+rm -f wallify

install-magick:
ifeq ($(OS), Windows_NT)
	+@echo "Checking for ImageMagick on Windows..."
	+@where magick >nul 2>nul || (
		echo "ImageMagick not found, installing with Chocolatey..."
		@"%SystemRoot%\System32\WindowsPowerShell\v1.0\powershell.exe" -NoProfile -InputFormat None -ExecutionPolicy Bypass -Command "[System.Net.ServicePointManager]::SecurityProtocol = 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))" && SET "PATH=%PATH%;%ALLUSERSPROFILE%\chocolatey\bin"
		choco install imagemagick --no-progress -y
	)
else
    ifeq ($(shell uname -s), Linux)
		+@echo "Checking for ImageMagick on Linux..."
		+@if ! command -v magick &> /dev/null; then \
			echo "ImageMagick not found, attempting to install with a known package manager..."; \
			if command -v pacman &> /dev/null; then \
				echo "Using pacman..."; \
				sudo pacman -S --noconfirm imagemagick; \
			elif command -v dnf &> /dev/null; then \
				echo "Using dnf..."; \
				sudo dnf install -y ImageMagick; \
			elif command -v yum &> /dev/null; then \
				echo "Using yum..."; \
				sudo yum install -y ImageMagick; \
			elif command -v zypper &> /dev/null; then \
				echo "Using zypper..."; \
				sudo zypper install -y ImageMagick; \
			elif command -v apt-get &> /dev/null; then \
				echo "Using apt-get..."; \
				sudo apt-get update && sudo apt-get install -y imagemagick; \
			else \
				echo "ERROR: Could not find a known package manager (apt-get, pacman, dnf, yum, zypper). Please install ImageMagick manually."; \
				exit 1; \
			fi \
		fi
    else ifeq ($(shell uname -s), Darwin)
		+@echo "Checking for ImageMagick on macOS..."
		+@if ! command -v magick &> /dev/null; then \
			echo "ImageMagick not found, installing with Homebrew..."; \
			brew install imagemagick; \
		fi
    endif
endif

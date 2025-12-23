#!/bin/bash
set -e

# Bootstrap script for Arch Linux GUI Installer
# This script prepares the live environment and launches the GUI.

echo ">>> Initializing Arch GUI Installer Bootstrap..."

# 1. Verify Root
if [[ $EUID -ne 0 ]]; then
   echo "This script must be run as root." 
   exit 1
fi

# 2. Check Network & Space
echo ">>> Checking network connectivity..."
if ! ping -c 1 google.com &> /dev/null; then
    echo "Error: No internet connection. Please verify your network."
    exit 1
fi

echo ">>> Expanding COW space..."
# Remount COW (RAM disk) to use more space. Required for Xorg install on standard ISOs.
# We attempt to resize to 2G or 75% of RAM, whichever is safe? 
# Simplest approach: Remount with size=4G (kernel handles limits if RAM is lower, usually OOMs, but better than instant fail)
mount -o remount,size=4G /run/archiso/cow || echo "Warning: Failed to remount COW. Disk space might be low."

# 3. Install Dependencies
echo ">>> Installing Xorg and Window Manager..."
# Update database first
pacman -Sy

# Install minimal X environment
# xorg-server: display server
# xorg-xinit: handles startx
# fluxbox: lightweight WM
# xterm: terminal for debug
# glibc: usually present, but ensure base deps
pacman -S --noconfirm --needed xorg-server xorg-xinit fluxbox xterm ttf-dejavu

# 4. Download Components
WORK_DIR="/opt/arch-installer" # Use /opt or /root
mkdir -p "$WORK_DIR/backend"
cd "$WORK_DIR"

# Base URL for downloading artifacts
# User's Repo: https://github.com/mathisen99/Archlinux-Gui-Installer
REPO_URL="https://raw.githubusercontent.com/mathisen99/Archlinux-Gui-Installer/main"
RELEASE_URL="https://github.com/mathisen99/Archlinux-Gui-Installer/releases/latest/download"

echo ">>> Downloading installer components..."

# Download Backend Script
if [[ ! -f "backend/arch-install.sh" ]]; then
    echo " -> Downloading backend/arch-install.sh..."
    # Using curl -f to fail on 404
    curl -fsSL -o backend/arch-install.sh "$REPO_URL/backend/arch-install.sh" || {
        echo "Error: Failed to download backend script. Check REPO_URL."
        exit 1
    }
    chmod +x backend/arch-install.sh
fi

# Download GUI Binary
INSTALLER_BIN="./archgui"
if [[ ! -f "$INSTALLER_BIN" ]]; then
    echo " -> Downloading archgui binary..."
    curl -fsSL -o archgui "$RELEASE_URL/archgui" || {
        echo "Error: Failed to download archgui binary. Check RELEASE_URL."
        # Fallback for testing if file exists in current dir (dev mode)
        if [[ -f "/tmp/archgui" ]]; then
             cp /tmp/archgui .
        else 
             echo "Critical: Installer binary not found."
             exit 1
        fi
    }
    chmod +x archgui
fi

# 5. Create xinitrc
echo ">>> Configuring X session..."
cat > /root/.xinitrc <<EOF
# Start Fluxbox
fluxbox &

# Start our Installer
# We must be in the WORK_DIR so it can find backend/arch-install.sh
cd "$WORK_DIR"
exec xterm -e ./archgui
EOF

# 6. Launch
echo ">>> Starting Graphical Interface..."
startx

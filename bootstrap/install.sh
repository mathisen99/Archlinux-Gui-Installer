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

echo ">>> Debug: Current Disk and RAM usage:"
df -h / /run/archiso/cow 2>/dev/null || df -h /
free -h

echo ">>> Expanding COW space..."
# Try to find the actual mount point of the COW partition
# Usually /run/archiso/cow or /run/archiso/cowspace
COW_PATH="/run/archiso/cow"
if ! mountpoint -q "$COW_PATH"; then
    COW_PATH="/run/archiso/cowspace"
fi

# Remount with 4G explicitly
mount -o remount,size=4G "$COW_PATH" || echo "Warning: Failed to remount $COW_PATH"

echo ">>> Post-resize Disk usage:"
df -h "$COW_PATH" || df -h /

# Check if we actually have enough space (need ~500MB+ for Xorg)
FREE_KB=$(df -k "$COW_PATH" | awk 'NR==2 {print $4}')
if [[ "$FREE_KB" -lt 1000000 ]]; then # Less than 1GB free
    echo "WARNING: Less than 1GB free in COW. Installation might fail."
    echo "Recommendation: Increase VM RAM to at least 2GB/4GB."
fi

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

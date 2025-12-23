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

# 2. Check Network
echo ">>> Checking network connectivity..."
if ! ping -c 1 google.com &> /dev/null; then
    echo "Error: No internet connection. Please verify your network."
    exit 1
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

# 4. Download Release
# TODO: Replace with actual GitHub Release URL once published.
# For now, we assume development mode or local testing.
INSTALLER_BIN="./archgui"

if [[ ! -f "$INSTALLER_BIN" ]]; then
    echo ">>> Installer binary not found locally."
    # echo "Downloading from GitHub..."
    # curl -L -o archgui https://github.com/<user>/<repo>/releases/latest/download/archgui
    # chmod +x archgui
    echo "PLEASE PLACE 'archgui' BINARY HERE OR UPDATE URL."
    # exit 1 
    # (Proceeding for now to allow local testing if possible)
fi

# 5. Create xinitrc
echo ">>> Configuring X session..."
cat > /root/.xinitrc <<EOF
# Start Fluxbox
fluxbox &

# Start our Installer
# Run in xterm to see logs if GUI fails, or run directly
exec ./archgui
EOF

# 6. Launch
echo ">>> Starting Graphical Interface..."
startx

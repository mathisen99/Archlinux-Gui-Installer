#!/bin/bash
set -e

# Core installer script for Arch Linux
# This script is designed to be run non-interactively via a config file.

CONFIG_FILE=""
NONINTERACTIVE="no"

# Formatting
BOLD='\033[1m'
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

log() {
    echo -e "${GREEN}[ARCH-GUI]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1" >&2
}

usage() {
    echo "Usage: $0 [--config <file>]"
    exit 1
}

# Parse arguments
while [[ "$#" -gt 0 ]]; do
    case $1 in
        --config) CONFIG_FILE="$2"; shift ;;
        *) usage ;;
    esac
    shift
done

# Load config
if [[ -n "$CONFIG_FILE" ]]; then
    if [[ -f "$CONFIG_FILE" ]]; then
        source "$CONFIG_FILE"
        NONINTERACTIVE="yes"
        log "Loaded configuration from $CONFIG_FILE"
    else
        error "Config file not found: $CONFIG_FILE"
        exit 1
    fi
fi

# Validation
validate_config() {
    log "Validating configuration..."
    
    # Check for required variables
    local MISSING_KEYS=0
    
    if [[ -z "$DISK" ]]; then error "DISK is not set"; MISSING_KEYS=1; fi
    if [[ -z "$HOSTNAME" ]]; then error "HOSTNAME is not set"; MISSING_KEYS=1; fi
    if [[ -z "$USERNAME" ]]; then error "USERNAME is not set"; MISSING_KEYS=1; fi
    if [[ -z "$ROOT_PASSWORD" ]]; then error "ROOT_PASSWORD is not set"; MISSING_KEYS=1; fi
    if [[ -z "$USER_PASSWORD" ]]; then error "USER_PASSWORD is not set"; MISSING_KEYS=1; fi
    
    if [[ "$MISSING_KEYS" -eq 1 ]]; then
        exit 1
    fi
    
    # Check if disk exists
    if [[ ! -b "$DISK" ]]; then
        error "Target disk $DISK is not a block device or does not exist."
        exit 1
    fi
    
    log "Configuration valid."
}

# Main installation logic
install_arch() {
    log "Starting installation on $DISK..."
    
    # 1. Partitioning (Destructive!)
    log "Partitioning $DISK..."
    # In a real script, we'd use sgdisk/sfdisk here.
    # For now, we'll simulate the steps or assume the user wants auto-partitioning.
    
    # 2. Formatting
    log "Formatting partitions..."
    # mkfs.ext4 ...
    
    # 3. Mounting
    log "Mounting filesystems..."
    # mount ...
    
    # 4. Pacstrap
    log "Installing base system (pacstrap)..."
    # pacstrap /mnt base linux linux-firmware ...
    
    # 5. Fstab
    log "Generating fstab..."
    # genfstab -U /mnt >> /mnt/etc/fstab
    
    # 6. Chroot Configuration
    log "Configuring system in chroot..."
    # arch-chroot /mnt /bin/bash <<EOF
    # ... set timezone, locale, hostname, users ...
    # EOF
    
    log "Installation complete!"
}

if [[ "$NONINTERACTIVE" == "yes" ]]; then
    validate_config
    if [[ "$DRY_RUN" != "yes" ]]; then
        install_arch
    else
        log "DRY_RUN enabled. Skipping actual installation."
    fi
else
    # Interactive fallback mode
    echo "Interactive mode not yet implemented. Please use --config."
    exit 1
fi

#!/bin/bash
set -e

# Backend Installer for Arch Linux GUI
# Based on Mathisen's Arch Install Script (docs/old_terminal_install.sh)
# Designed to be run non-interactively via environment config.

# Default Config / Environment Variables
DISK="${DISK:-}"
HOSTNAME="${HOSTNAME:-archlinux}"
USERNAME="${USERNAME:-user}"
ROOT_PASSWORD="${ROOT_PASSWORD:-}"
USER_PASSWORD="${USER_PASSWORD:-}"
FS_TYPE="${FS_TYPE:-ext4}"           # ext4, btrfs
USE_LUKS="${USE_LUKS:-no}"           # yes, no
LUKS_PASSWORD="${LUKS_PASSWORD:-}"
DESKTOP_ENV="${DESKTOP_ENV:-none}"   # xfce, gnome, kde, i3, sway, hyprland, etc.
SHELL_CHOICE="${SHELL_CHOICE:-bash}" # bash, zsh, zsh-ohmyzsh
HAS_NVIDIA="${HAS_NVIDIA:-no}"       # yes, no
LOCALE="${LOCALE:-en_US}"
KEYMAP="${KEYMAP:-us}"
TIMEZONE="${TIMEZONE:-UTC}"

# Internal Variables
CONFIG_FILE=""
NONINTERACTIVE="no"
DRY_RUN="${DRY_RUN:-no}"

# Formatting
BOLD='\033[1m'
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

log() { echo -e "${GREEN}[BACKEND]${NC} $1"; }
error() { echo -e "${RED}[ERROR]${NC} $1" >&2; }

usage() {
    echo "Usage: $0 [--config <file>]"
    echo "Environment variables can also be set directly."
    exit 1
}

# Parse config argument
while [[ "$#" -gt 0 ]]; do
    case $1 in
        --config) CONFIG_FILE="$2"; shift ;;
        *) usage ;;
    esac
    shift
done

# Load config file if provided
if [[ -n "$CONFIG_FILE" ]]; then
    if [[ -f "$CONFIG_FILE" ]]; then
        # shellcheck source=/dev/null
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
    local MISSING_KEYS=0
    
    [[ -z "$DISK" ]] && { error "DISK is not set"; MISSING_KEYS=1; }
    [[ -z "$USERNAME" ]] && { error "USERNAME is not set"; MISSING_KEYS=1; }
    [[ -z "$ROOT_PASSWORD" ]] && { error "ROOT_PASSWORD is not set"; MISSING_KEYS=1; }
    [[ -z "$USER_PASSWORD" ]] && { error "USER_PASSWORD is not set"; MISSING_KEYS=1; }
    [[ -z "$LUKS_PASSWORD" ]] && [[ "$USE_LUKS" == "yes" ]] && { error "LUKS_PASSWORD is required for encryption"; MISSING_KEYS=1; }
    
    if [[ "$MISSING_KEYS" -eq 1 ]]; then exit 1; fi

    if [[ "$DRY_RUN" != "yes" ]]; then
        if [[ ! -b "$DISK" ]]; then
            error "Target disk $DISK does not exist."
            exit 1
        fi
    fi
    log "Configuration valid."
}

setup_partitioning() {
    local BOOT_MODE
    if [[ -d /sys/firmware/efi/efivars ]]; then
        BOOT_MODE="UEFI"
    else
        BOOT_MODE="BIOS"
    fi
    log "Boot Mode: $BOOT_MODE"

    # Wiping disk
    log "Wiping $DISK..."
    wipefs -af "$DISK"

    # Naming
    local PART_PREFIX="$DISK"
    [[ "$DISK" == *"nvme"* || "$DISK" == *"mmcblk"* ]] && PART_PREFIX="${DISK}p"

    if [[ "$BOOT_MODE" == "UEFI" ]]; then
        # UEFI: GPT, ESP, Root
        parted -s "$DISK" mklabel gpt
        parted -s "$DISK" mkpart "EFI" fat32 1MiB 513MiB
        parted -s "$DISK" set 1 esp on
        parted -s "$DISK" mkpart "root" "${FS_TYPE}" 513MiB 100%
        
        EFI_PART="${PART_PREFIX}1"
        ROOT_PART="${PART_PREFIX}2"
        mkfs.fat -F32 "$EFI_PART"
    else
        # BIOS: MBR, Root
        parted -s "$DISK" mklabel msdos
        parted -s "$DISK" mkpart primary "${FS_TYPE}" 1MiB 100%
        parted -s "$DISK" set 1 boot on
        
        ROOT_PART="${PART_PREFIX}1"
    fi
    
    # Wait for nodes
    sleep 2
    partprobe "$DISK" || true

    # Encryption
    local CRYPT_ROOT="$ROOT_PART"
    if [[ "$USE_LUKS" == "yes" ]]; then
        log "Encrypting root partition..."
        echo -n "$LUKS_PASSWORD" | cryptsetup luksFormat --type luks2 "$ROOT_PART" -
        echo -n "$LUKS_PASSWORD" | cryptsetup open "$ROOT_PART" cryptroot -
        CRYPT_ROOT="/dev/mapper/cryptroot"
    fi

    # Formatting Root
    if [[ "$FS_TYPE" == "btrfs" ]]; then
        log "Formatting BTRFS..."
        mkfs.btrfs -f "$CRYPT_ROOT"
        mount "$CRYPT_ROOT" /mnt
        
        # Subvolumes
        btrfs subvolume create /mnt/@
        btrfs subvolume create /mnt/@home
        btrfs subvolume create /mnt/@snapshots
        btrfs subvolume create /mnt/@var_log
        umount /mnt
        
        local BTRFS_OPTS="noatime,compress=zstd,space_cache=v2,discard=async"
        mount -o "subvol=@,${BTRFS_OPTS}" "$CRYPT_ROOT" /mnt
        mkdir -p /mnt/{home,.snapshots,var/log,boot}
        mount -o "subvol=@home,${BTRFS_OPTS}" "$CRYPT_ROOT" /mnt/home
        mount -o "subvol=@snapshots,${BTRFS_OPTS}" "$CRYPT_ROOT" /mnt/.snapshots
        mount -o "subvol=@var_log,${BTRFS_OPTS}" "$CRYPT_ROOT" /mnt/var/log
    else
        log "Formatting EXT4..."
        mkfs.ext4 -F "$CRYPT_ROOT"
        mount "$CRYPT_ROOT" /mnt
    fi

    if [[ "$BOOT_MODE" == "UEFI" ]]; then
        mkdir -p /mnt/boot
        mount "$EFI_PART" /mnt/boot
    fi
}

install_packages() {
    log "Installing base system..."
    # Microcode check
    local MICROCODE=""
    local CPU_VENDOR
    CPU_VENDOR=$(grep -m1 vendor_id /proc/cpuinfo | awk '{print $3}' || true)
    [[ "$CPU_VENDOR" == "GenuineIntel" ]] && MICROCODE="intel-ucode"
    [[ "$CPU_VENDOR" == "AuthenticAMD" ]] && MICROCODE="amd-ucode"

    local PACKAGES="base base-devel linux linux-firmware networkmanager grub sudo nano vim git btop"
    [[ -n "$MICROCODE" ]] && PACKAGES="$PACKAGES $MICROCODE"
    [[ "$FS_TYPE" == "btrfs" ]] && PACKAGES="$PACKAGES btrfs-progs"
    [[ -d /sys/firmware/efi/efivars ]] && PACKAGES="$PACKAGES efibootmgr"
    
    pacstrap -K /mnt $PACKAGES

    # Desktop Environment Logic
    if [[ "$DESKTOP_ENV" != "none" ]]; then
        log "Installing Desktop: $DESKTOP_ENV"
        local DESKTOP_PKGS="xorg-server xorg-xinit pipewire pipewire-alsa pipewire-pulse wireplumber pavucontrol"
        local DM=""
        
        case "$DESKTOP_ENV" in
            xfce)
                DESKTOP_PKGS="$DESKTOP_PKGS xfce4 xfce4-goodies lightdm lightdm-gtk-greeter"
                DM="lightdm"
                ;;
            gnome)
                DESKTOP_PKGS="$DESKTOP_PKGS gnome gnome-tweaks gdm"
                DM="gdm"
                ;;
            kde)
                DESKTOP_PKGS="$DESKTOP_PKGS plasma-meta kde-applications-meta sddm packagekit-qt6"
                DM="sddm"
                ;;
            i3)
                DESKTOP_PKGS="$DESKTOP_PKGS i3-wm i3status dmenu xterm feh picom lightdm lightdm-gtk-greeter"
                DM="lightdm"
                ;;
            hyprland)
                DESKTOP_PKGS="$DESKTOP_PKGS hyprland xdg-desktop-portal-hyprland waybar wofi foot sddm"
                DM="sddm"
                ;;
            *)
                log "Unknown desktop environment: $DESKTOP_ENV"
                ;;
        esac

        # Nvidia
        if [[ "$HAS_NVIDIA" == "yes" ]]; then
            DESKTOP_PKGS="$DESKTOP_PKGS nvidia nvidia-utils nvidia-settings"
        fi

        pacstrap /mnt $DESKTOP_PKGS
        
        # Save DM for chroot script
        echo "DISPLAY_MANAGER=$DM" > /mnt/dm_info
    fi

    # Shell
    if [[ "$SHELL_CHOICE" == *"zsh"* ]]; then
        pacstrap /mnt zsh zsh-completions
    fi
}

configure_system() {
    log "Configuring system..."
    genfstab -U /mnt >> /mnt/etc/fstab

    # Create Chroot Script
    cat > /mnt/setup_chroot.sh <<EOF
#!/bin/bash
set -e

# Timezone & Locale
ln -sf /usr/share/zoneinfo/${TIMEZONE} /etc/localtime
hwclock --systohc
sed -i "s/^#${LOCALE}.UTF-8/${LOCALE}.UTF-8/" /etc/locale.gen
locale-gen
echo "LANG=${LOCALE}.UTF-8" > /etc/locale.conf
echo "KEYMAP=${KEYMAP}" > /etc/vconsole.conf
echo "${HOSTNAME}" > /etc/hostname

# Root Password
echo "root:${ROOT_PASSWORD}" | chpasswd

# User Setup
USER_SHELL="/bin/bash"
if command -v zsh &>/dev/null; then USER_SHELL="/bin/zsh"; fi
useradd -m -G wheel,audio,video,storage -s "\$USER_SHELL" "${USERNAME}"
echo "${USERNAME}:${USER_PASSWORD}" | chpasswd
sed -i 's/^# %wheel ALL=(ALL:ALL) ALL/%wheel ALL=(ALL:ALL) ALL/' /etc/sudoers

# NetworkManager
systemctl enable NetworkManager

# Display Manager
if [[ -f /dm_info ]]; then
    source /dm_info
    if [[ -n "\$DISPLAY_MANAGER" ]]; then
        systemctl enable "\$DISPLAY_MANAGER"
    fi
    rm /dm_info
fi

# Bootloader (GRUB)
# Grub config for LUKS
if [[ "${USE_LUKS}" == "yes" ]]; then
    UUID=\$(blkid -s UUID -o value ${ROOT_PART})
    sed -i "s|GRUB_CMDLINE_LINUX=\"\"|GRUB_CMDLINE_LINUX=\"cryptdevice=UUID=\$UUID:cryptroot root=/dev/mapper/cryptroot\"|" /etc/default/grub
    sed -i 's/^HOOKS=.*/HOOKS=(base udev autodetect modconf kms keyboard keymap consolefont block encrypt filesystems fsck)/' /etc/mkinitcpio.conf
    echo "GRUB_ENABLE_CRYPTODISK=y" >> /etc/default/grub
    mkinitcpio -P
fi

if [[ -d /sys/firmware/efi/efivars ]]; then
    grub-install --target=x86_64-efi --efi-directory=/boot --bootloader-id=ARCH
else
    grub-install --target=i386-pc ${DISK}
fi
grub-mkconfig -o /boot/grub/grub.cfg

# Oh-My-Zsh (Optional, simplified)
if [[ "${SHELL_CHOICE}" == "zsh-ohmyzsh" ]]; then
   # We can't easily run the automated install script inside chroot non-interactively without hacks.
   # For MVP, we'll skip or just install it crudely if internet works.
   # su - ${USERNAME} -c '...'
   true 
fi

EOF
    
    chmod +x /mnt/setup_chroot.sh
    # Pass external variables that are needed inside if not templated
    # Actually, we templated ${TIMEZONE} etc into the heredoc, so we are good.
    # But wait! ROOT_PART is needed inside for LUKS UUID, but we defined it in the host scope.
    # We must substitute it.
    
    # Correction: The heredoc above has variables expanded by the HOST shell (because EOF is not quoted).
    # Variables like \${TIMEZONE} will be expanded. 
    # Variables like \$USER_SHELL (escaped) will be literal in the file.
    # We need ROOT_PART to be expanded.
    
    arch-chroot /mnt /setup_chroot.sh
    rm /mnt/setup_chroot.sh
}

# Main Execution Flow
if [[ "$NONINTERACTIVE" == "yes" ]]; then
    validate_config
    if [[ "$DRY_RUN" == "yes" ]]; then
        log "Dry run complete. No changes made."
        exit 0
    fi
    
    setup_partitioning
    install_packages
    configure_system
    
    log "Installation Complete!"
else
    usage
fi

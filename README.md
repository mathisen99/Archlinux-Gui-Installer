# Arch Linux GUI Installer

A modern, minimal GUI installer for Arch Linux, built with Fyne and a robust bash backend.

## 🚀 Quick Start (Simulation)

To test the installer components without wiping your drive:

### 1. Prerequisites
- **Go** (for building the GUI)
- **Fyne Dependencies** (C compiler, GL headers)
  - *Arch Linux*: `pacman -S base-devel go xorg-server xorg-xinit fluxbox`

### 2. GUI Testing

Run the GUI in development mode:

```bash
cd gui/cmd/archgui
go run .
```

- Fill in the form (Disk, Hostname, etc.).
- Click **Install Arch Linux**.
- **Result**: The GUI generates a config file at `/tmp/install.env` and attempts to run the backend.
  - *Note*: It will likely fail or show "Permission denied" if not root, or "Disk not found" if testing safely. This is expected behavior for testing.

### 3. Backend Testing

Run the backend script directly with a mock configuration to verify validation logic:

1. Create a mock config (already provided as `mock_config.env`):
   ```bash
   DISK=/dev/null
   HOSTNAME=arch-test
   USERNAME=testuser
   ROOT_PASSWORD=root
   USER_PASSWORD=user
   # NONINTERACTIVE=yes is inferred if --config is passed
   ```

2. Run the script:
   ```bash
   ./backend/arch-install.sh --config mock_config.env
   ```
   **Output**: Should error out safely because `/dev/null` isn't a valid install target, or proceed if using a loopback device.

### 4. Bootstrap Testing

The `bootstrap/install.sh` script is the entry point for the live ISO.

```bash
# Validate syntax
bash -n bootstrap/install.sh
```

---

## 📂 Project Structure

- **`bootstrap/`**: Scripts for the ISO entry point (install Xorg, launch GUI).
- **`gui/`**: Fyne-based graphical interface source code.
- **`backend/`**: Bash scripts that handle the actual Arch installation logic.

## ⚠️ Warning

This software is a **system installer**. In production mode, it **WILL WIPE** the target disk specified in the configuration. Always use the `--config` flag with caution or tested in a VM.
# Goal

User boots the *stock* Arch Linux install ISO (archiso), gets network online, then runs one command:

```bash
curl -fsSL https://<your-domain-or-raw-github>/install.sh | bash
```

That script installs prerequisites, starts a minimal X session, and launches your **Fyne** GUI, which then drives your existing install backend script.

---

# Non-negotiable constraints on stock Arch ISO

* The stock ISO boots to a TTY; it does **not** include a desktop or window system by default. Your bootstrap must install and start an X session.
* Your installer backend must be **non-interactive** (config-driven). A GUI cannot reliably drive `read` prompts or `dialog` in-process.

---

# Deliverables (Repo structure)

Create a repo with these top-level items:

```
.
├─ bootstrap/
│  ├─ install.sh                # curl entrypoint; safe, idempotent
│  ├─ xinitrc                    # launches your GUI under X
│  └─ pacman-packages.txt        # list of packages to install on ISO
├─ gui/
│  ├─ cmd/archgui/main.go        # Fyne app entry
│  ├─ internal/...               # pages, state machine, validators
│  └─ assets/...                 # icons, branding
├─ backend/
│  ├─ arch-install.sh            # your existing script (refactored)
│  ├─ schema.env.example         # config keys the GUI writes
│  └─ lib/*.sh                   # optional split-out functions
└─ releases/
   └─ checksums.txt              # sha256 of released binaries
```

---

# Phase 1 — Refactor backend script for GUI control (do this first)

## Objective

Make `backend/arch-install.sh` runnable without any prompts:

```bash
bash arch-install.sh --config /tmp/install.env
```

## Required changes

1. Add `--config <file>` support to `source` env vars.
2. Replace every `read`/`dialog` prompt with:

   * if `NONINTERACTIVE=yes`: use env var
   * else: keep current prompt (so you can still test manually)
3. Add **validation** early (fail fast) with clear errors:

   * DISK exists and is a block device
   * passwords provided if encryption enabled
   * required keys set
4. Ensure logs are always printed to stdout/stderr (no hidden UI-only output).

## Output contract (recommendation)

Adopt an explicit config schema, e.g.:

* `DISK=/dev/nvme0n1`
* `FS_TYPE=ext4|btrfs`
* `USE_LUKS=yes|no`
* `LUKS_PASSWORD=...`
* `HOSTNAME=...`
* `TIMEZONE=...`
* `LOCALE=en_US.UTF-8`
* `KEYMAP=us`
* `USERNAME=...`
* `ROOT_PASSWORD=...`
* `USER_PASSWORD=...`
* `DESKTOP_ENV=none|xfce|gnome|kde|...`

Your GUI writes these keys; backend consumes them.

---

# Phase 2 — Build the Fyne GUI (minimal viable)

## Objective

A single binary (`archgui`) that:

* gathers install inputs
* writes `/tmp/install.env` (0600)
* executes `backend/arch-install.sh --config /tmp/install.env`
* streams stdout/stderr into a scrolling log view
* shows progress state (stepper) and final success/fail

## Suggested screens

1. Welcome + safety warning
2. Disk selection (parse `lsblk -J`)
3. Filesystem + encryption
4. System settings (hostname, timezone, locale, keymap)
5. User setup (username/passwords)
6. Desktop selection
7. Summary
8. Install (live log + progress)

## Technical notes

* Run backend in a goroutine; capture stdout/stderr pipes.
* UI updates must be marshalled onto the main thread.
* Disable “Install” until required fields validate.

---

# Phase 3 — Bootstrap script (the curl entrypoint)

## Objective

`bootstrap/install.sh` should:

1. verify running as root
2. verify network connectivity
3. `pacman -Sy --needed` install Xorg + minimal WM + dependencies
4. download your **released** GUI binary and backend script from GitHub Releases (or a stable URL)
5. verify checksum/signature
6. write `/root/.xinitrc` (or use `startx /path/to/xinitrc`)
7. start X and launch the GUI

## Recommended package set (minimal)

* `xorg-server`
* `xorg-xinit`
* `openbox` (or `fluxbox`)
* `xterm` (debug fallback)
* `mesa` (basic GL stack; helps on some systems)
* `wget` or `curl` (already present, but ensure)

## Why a window manager

It gives you reliable window focus/placement and a way to close/exit.

---

# Phase 4 — Distribution strategy (don’t compile on the ISO)

## Objective

Users should not build Go/Fyne on the live ISO.

* Build `archgui` for linux/amd64 (and optionally arm64) in CI.
* Upload binaries to GitHub Releases.
* Bootstrap downloads the correct binary and verifies `sha256`.

---

# Phase 5 — Testing plan

1. VM test matrix:

   * BIOS + UEFI
   * SATA + NVMe naming
   * ext4 and btrfs
   * encryption on/off
2. Real hardware smoke test:

   * Wi-Fi and Ethernet
   * NVIDIA/AMD/iGPU (Xorg start)
3. Failure tests:

   * invalid disk
   * missing passwords
   * pacstrap failure
   * bootloader failure

---

# Implementation sequence (do in this order)

1. Backend: add `--config`, remove hard interactive dependence, add validation.
2. Fyne MVP: one-page form + install button + log streaming.
3. Bootstrap: install Xorg + WM + download artifacts + `startx` launch.
4. Expand GUI to multi-step flow + nicer UX.
5. CI + releases + checksum verification.

---

# Security and safety requirements (practical)

* Always display a confirmation screen that clearly states the target disk will be wiped.
* Require the user to type the disk name (or confirm twice) before proceeding.
* Store passwords only in-memory where possible; if written to env file, use `chmod 600` and delete after install.
* Consider GPG-signing releases; at minimum, enforce SHA256 verification.

---

# “curl command” example you will publish

Once implemented, you publish a stable URL to `bootstrap/install.sh`:

```bash
curl -fsSL https://raw.githubusercontent.com/<you>/<repo>/main/bootstrap/install.sh | bash
```

That script does the rest.

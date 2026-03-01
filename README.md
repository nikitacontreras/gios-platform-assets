# GIOS Platform Assets 🛰️📱

This repository centralizes the necessary assets for [GIOS](https://github.com/nikitacontreras/gios) to work seamlessly across all iOS platforms.

## 📦 What's inside?
- **SDKs**: Header and library files for cross-compiling Go to iOS.
- **DDIs**: Developer Disk Images required for debugging and advanced features (screenshots, reboot, etc.) on physical devices.

## 🔗 Original Sources
- SDKs based on [theos/sdks](https://github.com/theos/sdks).
- Disk Images based on [haikieu/xcode-developer-disk-image-all-platforms](https://github.com/haikieu/xcode-developer-disk-image-all-platforms).

## 🚀 Usage in GIOS
GIOS automatically consults the `assets.json` manifest in this repository to manage installations. Use `gios sdk add` or `gios mount` to interactively install these components.

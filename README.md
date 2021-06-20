# modpack

A simple and light mod pack for Fabric users with the bare essentials.

Comes with a GUI installer to vastly simplify installation and enable updates.

## Download

Download the latest version for your platform from [here](https://github.com/retrixe/modpack/releases), open it and let it download and install mods. Currently supports Windows 10, Windows 7/8.1 with Edge Beta or newer and Linux. macOS is also supported but there are no pre-built binaries available.

## Updates

Whenever mods are updated, an entry is made in  a `mods.json` in the zip file for that version of the mod. Upon installing when mods from this pack are already installed, the installer will match the major Minecraft version selected for install with the version of the mods in the `mods` folder, determined by `modsversion.txt`. If both belong to the same major Minecraft version, the existing mods will be updated instead by looping through `mods.json`, finding old mods, and replacing them. If they do not/the folder does not exist, the folder will be renamed to `oldmods` and a fresh install will occur.

## Servers and Development

It is possible to repurpose this mod installer for other packs of your own. The modpack in its current state can be made to use an alternate server to get mod data using the `MODS_JSON_URL` environment variable, however, it will only be able to use the pre-defined Minecraft versions. You can modify the headers, FAQ and available Minecraft versions fairly easy to create your own spins of the installer with all of its functionality.

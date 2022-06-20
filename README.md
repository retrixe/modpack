# modpack

A simple and light mod pack for Fabric users with the bare essentials.

This is a GUI installer that makes it easy to install and update the mods.

![Welcome Screen](https://cdn.discordapp.com/attachments/588340346841464835/988554355923886160/unknown.png)

## Installation

Download the latest version for your platform from [here](https://github.com/retrixe/modpack/releases), open it and let it download and install mods. Currently supports Windows 10, macOS and Linux.

## Update

This installer supports downloading and updating the mods in your `mods` folder with the latest versions currently available in the modpack, without removing any mods you added or adding back any mods you removed. Download the latest installer, and you will receive this prompt to update when you press Continue.

![Update Functionality](https://cdn.discordapp.com/attachments/839933525045149716/988550158830682112/unknown.png)

## For Developers

The `ui/` folder contains the webview UI written in React. To build modpack, run `yarn build` in the `ui/` folder to build a JavaScript bundle (after running `yarn` to install JavaScript dependencies), then run `go build -ldflags="-s -w" -v .` in the top-level folder to package a Golang executable. On Windows, you may need `WebView2Loader.dll` and `webview.dll` (check the `.github/workflows/go.yml` file).

For convenience when developing, `yarn start` in the `ui/` folder will start the Golang app automatically for you.

### Updates

Whenever mods are updated, an entry is made in  a `mods.json` in the zip file for that version of the mod. Upon installing when mods from this pack are already installed, the installer will match the major Minecraft version selected for install with the version of the mods in the `mods` folder, determined by `modsversion.txt`. If both belong to the same major Minecraft version, the existing mods will be updated instead by looping through `mods.json`, finding old mods, and replacing them. If they do not/the folder does not exist, the folder will be renamed to `oldmods` and a fresh install will occur.

## Servers

It is possible to repurpose this mod installer for other packs of your own. The modpack in its current state can be made to use an alternate server to get mod data using the `MODS_JSON_URL` environment variable, however, it will only be able to use the pre-defined Minecraft versions. You can modify the headers, FAQ and available Minecraft versions fairly easy to create your own spins of the installer with all of its functionality.

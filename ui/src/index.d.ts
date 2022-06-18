interface Window {
  // Golang APIs.
  promptForFolder: () => void
  updateMinecraftFolder: (folder: string) => void
  changeVersion: (version: string) => void

  // JavaScript APIs.
  setMinecraftVersionState: (version: string) => void
  setUpdatableVersionState: (version: string) => void
  setMinecraftFolderState: (version: string) => void
}

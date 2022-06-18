interface Window {
  // Golang APIs.
  installMods: () => void
  promptForFolder: () => void
  toggleInstallFabric: () => void
  updateMinecraftFolder: (folder: string) => void
  changeVersion: (version: string) => void
  respondQuery: (answer: boolean) => void

  // JavaScript APIs.
  setMinecraftVersionState: (version: string) => void
  setUpdatableVersionState: (version: string) => void
  setMinecraftFolderState: (version: string) => void
  setInProgressState: (inProgress: boolean) => void
  setMessageState: (progress: string) => void
  setErrorState: (error: string) => void
  setQueryState: (query: string) => void
}

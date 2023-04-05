import { useState, useEffect, useRef } from 'react'
import { Grid, Paper, Typography, Divider, Button } from '@mui/material'
import WelcomeScreen from './screens/WelcomeScreen'
import VersionSelectionScreen from './screens/VersionSelectionScreen'
import ModSelectionScreen from './screens/ModSelectionScreen'
import InstallationScreen from './screens/InstallationScreen'
import Faq from './Faq'

const App = (): JSX.Element => {
  const [faqOpen, setFaqOpen] = useState(false)
  const [inProgress, setInProgress] = useState(false)
  const [currentStep, setCurrentStep] = useState(1)
  const [minecraftFolder, setMinecraftFolderState] = useState('')
  const [minecraftVersion, setMinecraftVersionState] = useState('')
  const [updatableVersions, setUpdatableVersions] = useState<string[]>([])
  const [installFabric, setInstallFabricState] = useState(true)
  console.log(updatableVersions)

  // Any changes to Minecraft folder/version should propagate to Go.
  // Go can also edit the UI via folder prompt, as well as selected and updatable Minecraft versions.
  window.setInProgressState = setInProgress
  window.setUpdatableVersionsState = setUpdatableVersions
  window.setMinecraftVersionState = setMinecraftVersionState
  window.setMinecraftFolderState = setMinecraftFolderState
  const setMinecraftVersion = (newState: string): void => {
    window.changeVersion(newState)
    setMinecraftVersionState(newState)
  }
  const setMinecraftFolder = (newState: string): void => {
    window.updateMinecraftFolder(newState)
    setMinecraftFolderState(newState)
  }
  const toggleInstallFabric = (): void => {
    window.toggleInstallFabric()
    setInstallFabricState(state => !state)
  }
  // Request UI update once to get info about the current updatable version.
  const calledUiUpdateOnce = useRef(false)
  useEffect(() => {
    if (!calledUiUpdateOnce.current) {
      calledUiUpdateOnce.current = true
      window.updateMinecraftFolder('')
    }
  }, [])

  const handleFaqOpen = (): void => setFaqOpen(true)
  if (faqOpen) return <Faq close={() => setFaqOpen(false)} />
  return (
    <div css={{ height: '100%', padding: '8px', boxSizing: 'border-box' }}>
      <Grid container spacing={2} height='100%'>
        <Grid item xs={4} md={2}>
          <Paper css={{ height: '100%', padding: '8px', display: 'flex', flexDirection: 'column' }}>
            <Typography color={currentStep === 1 ? undefined : '#666666'} variant={currentStep === 1 ? 'h6' : undefined}>
              <b>Step 1:</b> Welcome
            </Typography>
            <Divider light sx={{ marginTop: '8px', marginBottom: '8px' }} />
            <Typography color={currentStep === 2 ? undefined : '#666666'} variant={currentStep === 2 ? 'h6' : undefined}>
              <b>Step 2:</b> Select Minecraft version
            </Typography>
            <Divider light sx={{ marginTop: '8px', marginBottom: '8px' }} />
            <Typography color={currentStep === 3 ? undefined : '#666666'} variant={currentStep === 3 ? 'h6' : undefined}>
              <b>Step 3:</b> Optional: Select mods to install
            </Typography>
            <Divider light sx={{ marginTop: '8px', marginBottom: '8px' }} />
            <Typography color={currentStep === 4 ? undefined : '#666666'} variant={currentStep === 4 ? 'h6' : undefined}>
              <b>Step 4:</b> Installation
            </Typography>
            <div css={{ flex: 1 }} />
            <Button variant='outlined' color='info' onClick={handleFaqOpen} disabled={inProgress}>
              FAQ
            </Button>
          </Paper>
        </Grid>
        <Grid item xs={8} md={10} css={{ display: 'flex', flexDirection: 'column', padding: '8px' }}>
          {currentStep === 1 && (
            <WelcomeScreen
              setCurrentStep={setCurrentStep}
              minecraftFolder={minecraftFolder}
              setMinecraftFolder={setMinecraftFolder}
            />
          )}
          {currentStep === 2 && (
            <VersionSelectionScreen
              setCurrentStep={setCurrentStep}
              minecraftVersion={minecraftVersion}
              updatableVersions={updatableVersions}
              setMinecraftVersion={setMinecraftVersion}
            />
          )}
          {currentStep === 3 && (
            <ModSelectionScreen
              setCurrentStep={setCurrentStep}
              installFabric={installFabric}
              toggleInstallFabric={toggleInstallFabric}
              minecraftVersion={minecraftVersion}
            />
          )}
          {currentStep === 4 && (
            <InstallationScreen inProgress={inProgress} setCurrentStep={setCurrentStep} />
          )}
        </Grid>
      </Grid>
    </div>
  )
}

export default App

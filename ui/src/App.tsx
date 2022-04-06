import { useState } from 'react'
import { Grid, Paper, Typography, Divider, Button } from '@mui/material'
import WelcomeScreen from './screens/WelcomeScreen'
import VersionSelectionScreen from './screens/VersionSelectionScreen'

const App = (): JSX.Element => {
  const [currentStep, setCurrentStep] = useState(1)

  return (
    <div css={{ height: '100%', padding: '8px', boxSizing: 'border-box' }}>
      <Grid container spacing={2} height='100%'>
        <Grid item xs={4} md={2}>
          <Paper css={{ height: '100%', padding: '8px', display: 'flex', flexDirection: 'column' }}>
            <Typography variant={currentStep === 1 ? 'h6' : undefined}>
              <b>Step 1:</b> Welcome
            </Typography>
            <Divider light sx={{ marginTop: '8px', marginBottom: '8px' }} />
            <Typography variant={currentStep === 2 ? 'h6' : undefined}>
              <b>Step 2:</b> Select Minecraft version
            </Typography>
            <Divider light sx={{ marginTop: '8px', marginBottom: '8px' }} />
            <Typography variant={currentStep === 3 ? 'h6' : undefined}>
              <b>Step 3:</b> Optional: Select mods to install
            </Typography>
            <Divider light sx={{ marginTop: '8px', marginBottom: '8px' }} />
            <Typography variant={currentStep === 4 ? 'h6' : undefined}>
              <b>Step 4:</b> Confirm installation
            </Typography>
            <div css={{ flex: 1 }} />
            <Button variant='outlined' color='secondary'>FAQ</Button>
          </Paper>
        </Grid>
        <Grid item xs={8} md={10} css={{ display: 'flex', flexDirection: 'column', padding: '8px' }}>
          {currentStep === 1 && <WelcomeScreen setCurrentStep={setCurrentStep} />}
          {currentStep === 2 && <VersionSelectionScreen setCurrentStep={setCurrentStep} />}
        </Grid>
      </Grid>
    </div>
  )
}

export default App

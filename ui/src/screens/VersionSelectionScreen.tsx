import { Button, Typography, FormControl, FormControlLabel, RadioGroup, Radio } from '@mui/material'

const VersionSelectionScreen = ({ setCurrentStep }: {
  setCurrentStep: React.Dispatch<React.SetStateAction<number>>
}): JSX.Element => {
  // TODO: Support oldmodfolder prompt
  // TODO: Receive upgrade info from Go
  // TODO: Actually set Minecraft version somewhere
  return (
    <>
      <Typography variant='h5' textAlign='center' gutterBottom>ibu's mod installer</Typography>
      <Typography variant='h6' fontSize='1.1rem' textAlign='center' gutterBottom>
        Select Minecraft version to install mods for:
      </Typography>
      <FormControl>
        <RadioGroup>
          <FormControlLabel value='upgrade' control={<Radio />} label='Upgrade existing mods installed for Minecraft 1.19' />
          <FormControlLabel value='1.19' control={<Radio />} label='Minecraft 1.19' />
          <FormControlLabel value='1.18.2' control={<Radio />} label='Minecraft 1.18.2' />
          <FormControlLabel value='1.17.1' control={<Radio />} label='Minecraft 1.17.1 (bug-fixes only)' />
          <FormControlLabel value='1.16.5' control={<Radio />} label='Minecraft 1.16.5 (bug-fixes only)' />
        </RadioGroup>
      </FormControl>
      <div css={{ flex: 1 }} />
      <div css={{ display: 'flex', alignItems: 'stretch' }}>
        <Button variant='outlined' color='warning' onClick={() => setCurrentStep(1)}>Back</Button>
        <div css={{ padding: '4px' }} />
        <Button variant='outlined' color='secondary' onClick={() => setCurrentStep(3)} css={{ flex: 1 }}>
          Continue
        </Button>
      </div>
    </>
  )
}

export default VersionSelectionScreen

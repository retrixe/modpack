import { Button, Typography, FormControl, FormControlLabel, RadioGroup, Radio } from '@mui/material'

const VersionSelectionScreen = ({ setCurrentStep, minecraftVersion, setMinecraftVersion, updatableMcVersion }: {
  updatableMcVersion: string
  minecraftVersion: string
  setMinecraftVersion: (newState: string) => void
  setCurrentStep: React.Dispatch<React.SetStateAction<number>>
}): JSX.Element => {
  const l = (value: string, label: string): string => {
    return value === updatableMcVersion ? `Upgrade existing mods installed for ${label}` : label
  }
  return (
    <>
      <Typography variant='h5' textAlign='center' gutterBottom>ibu's mod installer</Typography>
      <Typography variant='h6' fontSize='1.1rem' textAlign='center' gutterBottom>
        Select Minecraft version to install mods for:
      </Typography>
      <FormControl>
        <RadioGroup value={minecraftVersion} onChange={event => setMinecraftVersion(event.target.value)}>
          <FormControlLabel value='1.19' control={<Radio />} label={l('1.19', 'Minecraft 1.19.2')} />
          <FormControlLabel value='1.18' control={<Radio />} label={l('1.18', 'Minecraft 1.18.2')} />
          <FormControlLabel value='1.17' control={<Radio />} label={l('1.17', 'Minecraft 1.17.1 (UNSUPPORTED)')} />
          <FormControlLabel value='1.16' control={<Radio />} label={l('1.16', 'Minecraft 1.16.5 (UNSUPPORTED)')} />
          <FormControlLabel value='1.15' control={<Radio />} label={l('1.15', 'Minecraft 1.15.2 (UNSUPPORTED)')} />
          <FormControlLabel value='1.14' control={<Radio />} label={l('1.14', 'Minecraft 1.14.5 (UNSUPPORTED)')} />
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

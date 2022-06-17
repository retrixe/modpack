import { Button, Typography, FormControl, FormControlLabel, RadioGroup, Radio } from '@mui/material'

const VersionSelectionScreen = ({ setCurrentStep }: {
  setCurrentStep: React.Dispatch<React.SetStateAction<number>>
}): JSX.Element => {
  // TODO: Upgrade support
  // TODO: Actually set this state somewhere
  return (
    <>
      <Typography variant='h5' textAlign='center' gutterBottom>ibu's mod installer</Typography>
      {/* <Divider flexItem light /> */}
      <Typography variant='h6' fontSize='1.1rem' textAlign='center' gutterBottom>
        Select Minecraft version to install mods for:
      </Typography>
      <FormControl>
        <RadioGroup>
          <FormControlLabel value='1.19' control={<Radio />} label='Minecraft 1.19' />
          <FormControlLabel value='1.18.2' control={<Radio />} label='Minecraft 1.18.2' />
          <FormControlLabel value='1.17.1' control={<Radio />} label='Minecraft 1.17.1 (bug-fixes only)' />
          <FormControlLabel value='1.16.5' control={<Radio />} label='Minecraft 1.16.5 (bug-fixes only)' />
        </RadioGroup>
      </FormControl>
      <div css={{ flex: 1 }} />
      <Button variant='outlined' color='secondary' onClick={() => setCurrentStep(1)}>
        Continue
      </Button>
    </>
  )
}

export default VersionSelectionScreen

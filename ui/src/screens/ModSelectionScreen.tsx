import { Button, Typography } from '@mui/material'

const ModSelectionScreen = ({ setCurrentStep }: {
  setCurrentStep: React.Dispatch<React.SetStateAction<number>>
}): JSX.Element => {
  // TODO: Make this work with modpack v2 API.
  return (
    <>
      <Typography variant='h5' textAlign='center' gutterBottom>ibu's mod installer</Typography>
      <Typography variant='h6' fontSize='1.1rem' textAlign='center' gutterBottom>
        Currently, this option is unavailable. Proceed to installation.
      </Typography>
      <div css={{ flex: 1 }} />
      <div css={{ display: 'flex', alignItems: 'stretch' }}>
        <Button variant='outlined' color='warning' onClick={() => setCurrentStep(2)}>Back</Button>
        <div css={{ padding: '4px' }} />
        <Button variant='outlined' color='secondary' onClick={() => setCurrentStep(4)} css={{ flex: 1 }}>
          Continue
        </Button>
      </div>
    </>
  )
}

export default ModSelectionScreen

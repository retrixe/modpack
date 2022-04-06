import { Button, Typography } from '@mui/material'

const VersionSelectionScreen = ({ setCurrentStep }: {
  setCurrentStep: React.Dispatch<React.SetStateAction<number>>
}): JSX.Element => {
  return (
    <>
      <Typography variant='h5' textAlign='center'>
        Installer for the MythicMC Fabric modpack
      </Typography>
      {/* TODO: Add information. */}
      <div css={{ flex: 1 }} />
      <Button variant='outlined' color='secondary' onClick={() => setCurrentStep(1)}>
        Continue
      </Button>
    </>
  )
}

export default VersionSelectionScreen

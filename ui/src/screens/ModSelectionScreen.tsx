import { Button, Typography, FormControlLabel, Checkbox } from '@mui/material'

const ModSelectionScreen = ({ setCurrentStep, minecraftVersion, installFabric, toggleInstallFabric }: {
  setCurrentStep: React.Dispatch<React.SetStateAction<number>>
  toggleInstallFabric: () => void
  installFabric: boolean
  minecraftVersion: string
}): JSX.Element => {
  const toggleFabric = (event: React.ChangeEvent<HTMLInputElement>, checked: boolean): void => {
    if ((checked && !installFabric) || (!checked && installFabric)) toggleInstallFabric()
  }
  // Yes, it's not pretty, but it'll be fixed eventually.
  const fabricVersions = ['1.18', '1.17', '1.16', '1.15', '1.14']
  fabricVersions.map(version => minecraftVersion.startsWith(version))
  const loaderName = fabricVersions.includes(minecraftVersion.substring(0, 4)) ? 'Fabric' : 'Quilt'
  return (
    <>
      <Typography variant='h5' textAlign='center' gutterBottom>ibu's mod installer</Typography>
      <FormControlLabel
        label={`Install ${loaderName} Loader (uncheck if you already have it and don't want to update it)`}
        control={<Checkbox checked={installFabric} onChange={toggleFabric} />}
        css={{ marginBottom: '8px' }}
      />
      <Typography variant='subtitle2' gutterBottom>
        Currently, there is no option to select what mods you would like to install.
        This is being worked on.
      </Typography>
      <div css={{ flex: 1 }} />
      <div css={{ display: 'flex', alignItems: 'stretch' }}>
        <Button variant='outlined' color='warning' onClick={() => setCurrentStep(2)}>Back</Button>
        <div css={{ padding: '4px' }} />
        <Button variant='outlined' color='secondary' onClick={() => setCurrentStep(4)} css={{ flex: 1 }}>
          Install
        </Button>
      </div>
    </>
  )
}

export default ModSelectionScreen

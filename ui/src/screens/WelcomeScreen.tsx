import { Button, Typography } from '@mui/material'

const WelcomeScreen = ({ setCurrentStep }: {
  setCurrentStep: React.Dispatch<React.SetStateAction<number>>
}): JSX.Element => {
  return (
    <>
      <Typography variant='h5' textAlign='center' gutterBottom>ibu's mod installer</Typography>
      <Typography gutterBottom>
        This installs fair-play Fabric mods useful in single-player and Factions/Survival-style servers.
      </Typography>
      <Typography gutterBottom>
        Mods include Sodium (a high-performance OptiFine replacement paired with support
        for zoom, shaders, capes, connected blocks and dynamic lights), Xaero's Minimap,
        MiniHUD, Command Macros, and many other features and gameplay improvements.
      </Typography>
      <Typography gutterBottom>
        1.18 and 1.19 are fully supported, while 1.16 and 1.17 only receive bug fixes.
      </Typography>
      <Typography gutterBottom>
        See the FAQ for more info about older Minecraft versions, OptiFine support and
        a complete list of mods. Check your server's rules before using this.
      </Typography>
      <div css={{ flex: 1 }} />
      <Button variant='outlined' color='secondary' onClick={() => setCurrentStep(2)}>
        Continue
      </Button>
    </>
  )
}

export default WelcomeScreen

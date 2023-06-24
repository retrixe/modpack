import { Button, Typography, TextField, SvgIcon, IconButton } from '@mui/material'

const WelcomeScreen = ({ setCurrentStep, minecraftFolder, setMinecraftFolder }: {
  minecraftFolder: string
  setMinecraftFolder: (newState: string) => void
  setCurrentStep: React.Dispatch<React.SetStateAction<number>>
}): JSX.Element => {
  return (
    <>
      <Typography variant='h5' textAlign='center' gutterBottom>ibu's mod installer</Typography>
      <Typography gutterBottom>
        This installs fair-play Fabric/Quilt mods useful in single-player and Factions/Survival-style servers.
      </Typography>
      <Typography gutterBottom>
        Mods include Sodium (a high-performance OptiFine replacement paired with mods
        for zoom, shaders, capes, connected blocks and dynamic lights), Xaero's Minimap,
        MiniHUD, Command Macros, and many other features and gameplay improvements.
      </Typography>
      <Typography gutterBottom>
        1.18, 1.19 and 1.20 are fully supported. 1.14 through 1.17 are available, but no longer updated.
      </Typography>
      <Typography gutterBottom>
        See the FAQ for more info about Quilt, older Minecraft versions, OptiFine support and
        a complete list of mods. Check your server's rules before using these.
      </Typography>
      <div css={{ flex: 1 }} />
      <div css={{ display: 'flex', alignItems: 'center', marginBottom: '8px' }}>
        <TextField
          value={minecraftFolder}
          onChange={e => setMinecraftFolder(e.target.value)}
          label='Advanced users: Path to game install folder'
          variant='outlined'
          css={{ flex: 1, marginRight: '4px' }}
        />
        <IconButton size='large' color='primary' onClick={() => window.promptForFolder()}>
          <SvgIcon fontSize='inherit'>
            <path d='M0 0h24v24H0z' fill='none' />
            {/* M10 4H4c-1.1 0-1.99.9-1.99 2L2 18c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V8c0-1.1-.9-2-2-2h-8l-2-2z */}
            <path d='M20 6h-8l-2-2H4c-1.1 0-1.99.9-1.99 2L2 18c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V8c0-1.1-.9-2-2-2zm0 12H4V8h16v10z' />
          </SvgIcon>
        </IconButton>
      </div>
      <Button variant='outlined' color='secondary' onClick={() => setCurrentStep(2)}>
        Continue
      </Button>
    </>
  )
}

export default WelcomeScreen

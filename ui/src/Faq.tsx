import { AppBar, IconButton, SvgIcon, Toolbar, Typography } from '@mui/material'

const Faq = ({ close }: { close: () => void }): JSX.Element => (
  <div css={{ height: '100%', boxSizing: 'border-box' }}>
    <AppBar position='fixed'>
      <Toolbar>
        <IconButton onClick={close} size='large' sx={{ mr: 2 }}>
          <SvgIcon fontSize='inherit'>
            <path d='M0 0h24v24H0z' fill='none' />
            <path d='M19 6.41L17.59 5 12 10.59 6.41 5 5 6.41 10.59 12 5 17.59 6.41 19 12 13.41 17.59 19 19 17.59 13.41 12z' />
          </SvgIcon>
        </IconButton>
        <Typography variant='h6'>FAQ - ibu's mod installer</Typography>
      </Toolbar>
    </AppBar>
    <Toolbar />
    <div css={{ padding: '8px' }}>
      <p><strong>Is this Fabric or Forge?</strong><br /> Fabric.</p>

      <p><strong>What Minecraft versions are supported?</strong><br /> 1.16, 1.17 and 1.18 are being actively updated.
        1.14 and 1.15 are unsupported, may have bugs (e.g. updating Fabric manually will break them) and should be avoided.
      </p>

      <p><strong>Does this modpack get updates?</strong><br /> Yes! You can re-run the installer to get the latest updates
        as well as any new mods I add :^) It also detects any mods you added/removed yourself, and will not readd or
        remove them. Just don't delete the <code>modsversion.txt</code> file for this to work correctly.
      </p>

      <p><strong>How to use these mods?</strong><br /> You can check the Controls and the Mods menu on how to use most of them.
        For FMap, MobCountMod, Watson and MiniHUD, you will need special keybinds (Y+C, P+C, L+C and H+C respectively).
      </p>

      <p><strong>OptiFine and 1.16+?</strong><br /> Due to OptiFine being slow and problematic, it has been replaced with Iris+Sodium
        and replacement mods for capes, zoom, show fps, better grass, connected textures (1.17+) and dynamic lights. A full list
        of replacement mods can be found at <a>https://lambdaurora.dev/optifine_alternatives/</a> for resource pack features. You
        can also use Canvas instead of Iris+Indium+Sodium if you want. For displaying your fps, use the H key (H+C for settings).
      </p>

      <p><strong>Can I still use OptiFine on 1.16+?</strong><br /> Sort of. You can download OptiFine and OptiFabric,
        however, it can cause conflicts with other mods, hence it's recommended to stick to Iris+Sodium or Canvas.
        <strong>There is no reason anymore to use OptiFine with this pack anymore apart from some resource pack features.</strong>
        Disable Capes, LambdaBetterGrass, Continuity, lambDynamicLights, Indium, Iris, Sodium,
        Lithium, Hydrogen and Phosphor/Starlight before using OptiFine with my modpack. Report
        any incompatibilities to me, but it is likely you will be told not to use OptiFine.
      </p>

      <p><strong>Full list for 1.16/1.17/1.18:</strong><br /> Architectury <i>(lib)</i>, ChunkBorders (1.16 only),
        Chat Utils, Cloth Config <i>(lib)</i>, Command Macros, Dark Loading Screen on 1.17+/Splash on 1.16,
        Dynamic FPS, EasierChests, Fabric API <i>(lib)</i>, Fabric Kotlin <i>(lib)</i>, FMapOverlayMod,
        VoxelMap (1.16/1.17), Hwyla/Wthit, Capes/Continuity (1.17+)/LambdaBetterGrass/lambDynamicLights/Ok
        Zoomer (1.16) or WI-Zoom (1.17+)/Hydrogen (1.16/1.17)/Phosphor or Starlight/Sodium/Indium <i>(lib)</i>/Iris, LightOverlay, Lithium, M-Tape,
        MaLiLib <i>(lib)</i>, Mc122477Fix, MiniHUD, MobCountMod, ModMenu, MyBrightness, ResolutionControl+, ScreenshotToClipboard,
        ShulkerBoxTooltip, ToroHealth, Watson, WorldEditCUI, Xaero's World Map/Minimap Fair-Play on 1.18+
      </p>

      <p><strong>What are the mods marked as <i>(lib)</i>?</strong><br /> These are mods required by
        certain mods in the pack, e.g. Architectury is needed by LightOverlay. Do not delete these
        from your mods folder, or Minecraft may fail to load.
      </p>

      <p><strong>What happens to my existing mods?</strong><br /> The <code>mods</code> folder will be renamed to <code>oldmods</code>.</p>

      <p><strong>Why does Minecraft slow down in the background?</strong><br /> This pack comes with the Dynamic FPS mod which limits
        Minecraft to 1fps when it's in the background to reduce CPU and GPU usage considerably. You can remove it if you want.
      </p>

      <p><strong>Ok Zoomer has a really weird zoom. How do I get OptiFine style zoom?</strong><br /> I feel you.
        Mods -&gt; Ok Zoomer -&gt; Settings button -&gt; Set Reset to Preset to Classic -&gt; Apply.
      </p>

      <p><strong>Are there any other mods I should consider?</strong><br /> WorldEdit is useful for single-player terrain editing.
        In-Game Account Switcher allows switching between accounts while in-game, and CraftPresence allows showing Minecraft
        in your Discord status, however, make sure you follow the rules of any Discord MC servers you are a member of when
        using this mod. These mods are not included for certain reasons, but they may be of use to you.
      </p>
    </div>
  </div>
)

export default Faq

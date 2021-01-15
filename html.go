package main

// Faq ... The HTML for the FAQ page.
const Faq = `
<!DOCTYPE html>
<html>
<head>
  <meta charset='utf-8'>
  <meta http-equiv='X-UA-Compatible' content='IE=edge'>
  <title>ibu's mod installer</title>
  <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/materialize/1.0.0/css/materialize.min.css">
  <link rel="stylesheet" href="https://fonts.googleapis.com/css2?family=Roboto:wght@200;300;400;500;700;900&display=swap">
  <meta name='viewport' content='width=device-width, initial-scale=1'>
  <!--
    <link rel='stylesheet' type='text/css' media='screen' href='main.css'>
    <script src='main.js'></script>
  -->
  <script>
    window.addEventListener('load', function () {
      document.getElementById('gui').addEventListener('click', function () {
        showGui()
      })
    })
  </script>
  <style>
    .header { /* Materialize specific styles */
      display: flex;
      padding-top: 1em;
      padding-bottom: 1em;
      align-items: center;
    }
    .header > * {
      display: inline-block;
    }
    body {
      padding: 8px;
      font-family: Roboto, sans-serif;
    }
  </style>
</head>
<body>
  <div class="header">
    <button class="waves-effect waves-light btn-small" id='gui'>back</button>
    <!-- h3 -->
    <h5 style="margin: 0 0 0 1rem;">ibu's mods - FAQ</h5>
  </div>
  <hr />
  <p><strong>Fabric or Forge?</strong> Fabric.</p>
  <p><strong>What Minecraft versions are supported?</strong> Currently, 1.16 is being actively updated while 1.15 is only getting new mods
  (not existing mod updates, these are mostly not very useful either). 1.14 is not being updated at all and the mods are
  out of date and some mods may be missing in the future. You may need to update Fabric API for 1.14/1.15 if using your own Fabric ver.</p>
  <p><strong>How to use these mods?</strong> You can check the Controls and the mod menu on how to use most of them. For FMap, MobCountMod,
  Watson and MiniHUD, you will need special keybinds (Y%2BC, P%2BC, L%2BC and H%2BC respectively).</p>
  <p><strong>OptiFine and 1.16?</strong> OptiFine is bundled with the 1.14 and 1.15 zips, while for 1.16, due to OptiFabric/OptiFine being
  problematic/slow respectively, it has been replaced in the zip with Sodium %2B replacements for zoom, show fps and dynamic lights. A full
  list of replacement mods can be found on the OptiFabric GitHub in case you want an OptiFine feature not in my mods. If
  wanted, you can also replace Sodium with Canvas for reduced performance but more quality adjustments, future shader
  support and compatibility with Connected Textures mod (Sodium will work with it soon, see its pinned issue). For showing
  fps, use H%2BC and enable fps in the MiniHUD mod, then use H to show it.</p>
  <p><strong>Ok Zoomer has a really weird zoom. How do I get OptiFine style zoom?</strong> I feel you.
  Mods -> Ok Zoomer -> Settings button -> Set Reset to Preset to Classic -> Apply.</p>
  <p><strong>Can I still use OptiFine on 1.16.5?</strong> Yes, you can download it from <a>https://optifine.net</a> and
  <a>https://www.curseforge.com/minecraft/mc-mods/optifabric</a>, however it can cause conflicts with other mods (the
  original reason OptiFabric was discontinued), hence it's recommended to stick to Sodium/Canvas <i>unless you really
  want shaders, apart from which there's no reason to use OptiFine anymore.</i> Make sure to disable lambDynamicLights,
  Sodium, Lithium and Phosphor before using OptiFabric 1.16 with my modpack. Please report any incompatibilities to me.</p>
  <p><strong>Where is VoxelMap data stored?</strong> <code>mods/mamiyaotaru</code> folder.</p>
  <p><strong>Full list:</strong> ChunkBorders, Chat Macros, EasierChests, Fabric API, FMapOverlayMod, VoxelMap,
  Hwyla, (lambDynamicLights/Ok Zoomer/Phosphor/Sodium/Capes) on 1.16, OptiFine%2BOptiFabric on 1.14/1.15,
  LightOverlay, Lithium, MaLiLib, MiniHUD, MobCountMod, MyBrightness, ShulkerBoxTooltip, Splash,
  ToroHealth, Watson, WorldEditCUI</p>
  <p><strong>What happens to my existing mods?</strong> They will be renamed to <code>oldmodfolder</code>.
  Currently there is a bug that if both <code>mods</code> and <code>oldmodfolder</code> exist, the mods will be duplicated
  in the <code>mods</code> folder. <strong>Hence, make sure your old mods are well clear before installing these mods!</strong></p>
  <p><strong>Can I update my mods when you update yours?</strong> Yes! However, ONLY if your mod folder is of the same MC
  version you selected, will you be able to update! Else, you need to move them out of the way. The version is determined by
  <code>modsversion.txt</code>. This method allows you to choose what mods you want and remove those you don't want and
  still be able to update. If you had mods of a different MC version, then they will be handled as in the previous question and
  hence is best to move them out of the way yourself.</p>
  <p><strong>Who is this for?</strong> Someone who wants most basic mods but not way too many.</p>
  <p><strong>Are there any other mods I should consider?</strong> WorldEdit is useful for single-player terrain editing.
  The Dynamic FPS mod limits Minecraft to 1fps when it's in the background to conserve your system resources, which can be
  very helpful. The Command Macros mod allows you to set keybinds to run commands of your choice in-game. You can also use
  the Rich Presence mod which is available at <a>https://github.com/HotLava03/rich-presence-mod/releases</a> to show
  Minecraft in your Discord status, however, make sure you follow the rules of any Discord MC servers you are a member
  of when using this mod. These mods are not included for certain reasons, but they may be of use to you.</p>
</body>
</html>
`

// HTML ... The HTML for the main page.
const HTML = `
<!DOCTYPE html>
<html>
<head>
  <meta charset='utf-8'>
  <meta http-equiv='X-UA-Compatible' content='IE=edge'>
  <title>ibu's mod installer</title>
  <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/materialize/1.0.0/css/materialize.min.css">
  <link rel="stylesheet" href="https://fonts.googleapis.com/css2?family=Roboto:wght@200;300;400;500;700;900&display=swap">
  <meta name='viewport' content='width=device-width, initial-scale=1'>
  <!--
    <link rel='stylesheet' type='text/css' media='screen' href='main.css'>
    <script src='main.js'></script>
  -->
  <script>
    window.addEventListener('load', function () {
      document.getElementById('select-version').addEventListener('change', function (event) {
        event.preventDefault()
        event.stopPropagation()
        changeVersion(event.target.value)
        event.srcElement.value = event.target.value
      })
      document.getElementById('install-fabric').addEventListener('change', function (event) {
        event.preventDefault()
        event.stopPropagation()
        toggleInstallFabric()
        event.srcElement.value = event.target.value
      })
      document.getElementById('install').addEventListener('click', function () {
        installMods()
      })
      document.getElementById('faq').addEventListener('click', function () {
        showFaq()
      })
    })
  </script>
  <style>
    #error {
      color: #ff4444;
    }
    body {
      text-align: center;
      font-family: Roboto, sans-serif;
    }
    /* Materialize specific styles */
    #progress-display {
      width: 12em;
      display: inline-block;
    }
    select {
      width: auto;
      display: inline-block !important;
    }
  </style>
</head>
<body>
  <!-- h2 -->
  <h4>installer for ibu's mods (Fabric only)</h4>
  <label for='select-version'>Minecraft Version:</label>
  <select id='select-version' name="Minecraft Version" class="browser-default">
    <option value="1.14.4">1.14.4 (see FAQ)</option>
    <option value="1.15.2">1.15.2 (see FAQ)</option>
    <option value="1.16.5" selected>1.16.5</option>
  </select>
  <br />
  <br />
  <!-- <label for="install-fabric">install fabric</label> -->
  <label>
    <input type="checkbox" class="filled-in" id="install-fabric" checked="checked" />
    <span>install fabric</span>
  </label>
  <br />
  <br />
  <button class="waves-effect waves-light btn-small" id='faq'>FAQ (read me)</button>
  <button class="waves-effect waves-light btn-small" id='install'>install</button>
  <br />
  <p style="display: none;" id="message">Done!</p>
  <p style="display: none;" id="error" />
  <p style="display: none;" id="progress">Working...</p>
  <div style="display: none;" id="progress-display" class="progress"><div class="indeterminate" /></div>
</body>
</html>
`

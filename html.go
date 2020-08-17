package main

// Faq ... The HTML for the FAQ page.
const Faq = `
<!DOCTYPE html>
<html>
<head>
  <meta charset='utf-8'>
  <meta http-equiv='X-UA-Compatible' content='IE=edge'>
  <title>ibu's mod installer</title>
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
    .header > * {
      display: inline-block;
    }
    body {
      font-family: sans-serif;
    }
  </style>
</head>
<body>
  <div class="header">
    <button id='gui'>back</button>  <h3>ibu's mods - FAQ</h3>
  </div>
  <hr />
  <p><strong>Fabric or Forge?</strong> Fabric.</p>
  <p><strong>What Minecraft versions are supported?</strong> Currently, 1.16 is being actively updated while 1.15 is only getting new mods
  (not existing mod updates, these are mostly not very useful either). 1.14 is not being updated at all and the mods are
  out of date and some mods may be missing in the future. You may need to update Fabric API for 1.14/1.15.</p>
  <p><strong>How to use these mods?</strong> You can check the Controls and the mod menu on how to use most of them. For MobCountMod,
  Watson and MiniHUD, you will need special keybinds (P%2BC, L%2BC and H%2BC respectively).</p>
  <p><strong>OptiFine and 1.16?</strong> OptiFine is bundled with the 1.14 and 1.15 zips, while for 1.16, due to OptiFabric being
  discontinued, it has been replaced in the zip with Sodium %2B replacements for zoom, show fps and dynamic lights. A full
  list of replacement mods can be found on the OptiFabric page in case you want an OptiFine feature not in my mods. If
  wanted, you can also replace Sodium with Canvas for reduced performance but more quality adjustments, future shader
  support and compatibility with Connected Textures mod (Sodium will work with it soon, see its pinned issue). For showing
  fps, use H%2BC and enable fps in the MiniHUD mod, then use H to show it.</p>
  <p><strong>1.16.1: Whenever I break a block while looking straight down, the animation doesn't display.</strong> Disable Compact Vertex
  Format or use the included dev build that fixes it (disable lambDynamicLights before doing so).</p>
  <p><strong>Full list:</strong> ChunkBorders, EasierChests, Fabric API, VoxelMap, Hwyla, (lambDynamicLights/Logical
  Zoom/Phosphor/Sodium)/OptiFine%2BOptiFabric, LightOverlay, MaLiLib, MiniHUD, MobCountMod, MyBrightness, ShulkerBoxTooltip,
  ToroHealth, Watson, WorldEditCUI</p>
  <p><strong>Who is this for?</strong> Someone who wants most basic mods but not way too many.</p>
  <p><strong>Faster single-player performance.</strong> Lithium mod.</p>
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
  <meta name='viewport' content='width=device-width, initial-scale=1'>
  <!--
    <link rel='stylesheet' type='text/css' media='screen' href='main.css'>
    <script src='main.js'></script>
  -->
  <script>
    window.addEventListener('load', function () {
      document.getElementById('select-version').addEventListener('change', function (event) {
        changeVersion(event.target.value)
      })
      document.getElementById('install-fabric').addEventListener('change', function (event) {
        toggleInstallFabric()
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
      font-family: sans-serif;
    }
  </style>
</head>
<body>
  <h2>installer for ibu's mods (Fabric only)</h2>
  <label for='select-version'>Minecraft Version:</label>
  <select id='select-version' name="Minecraft Version">
    <option value="1.14.4">1.14.4 (see FAQ)</option>
    <option value="1.15.2">1.15.2 (see FAQ)</option>
    <option value="1.16.1" selected>1.16.1</option>
  </select>
  <br />
  <br />
  <label for="install-fabric">install fabric</label>
  <input type="checkbox" id="install-fabric" checked />
  <br />
  <br />
  <button id='faq'>FAQ (read me)</button>
  <button id='install'>install</button>
  <br />
  <p style="display: none;" id="message">Done!</p>
  <p style="display: none;" id="error" />
  <p style="display: none;" id="progress">Working...</p>
  <progress style="display: none;" id="progress-display" />
</body>
</html>
`

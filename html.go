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
      font-family: sans-serif;
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
  <p><strong>What happens to my existing mods?</strong> They will be renamed to <code>oldmodfolder</code>oldmodfolder.
  Currently there is a bug that if both <code>mods</code> and <code>oldmodfolder</code> exist, the mods will be duplicated
  in the <code>mods</code> folder. Hence, make sure your old mods are well clear before installing mods.</p>
  <p><strong>Can I update my mods when you update yours?</strong> Currently, no. You have to move your old mod folder and
  backup your VoxelMap data, then restore after new mods are installed.</p>
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
  <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/materialize/1.0.0/css/materialize.min.css">
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
      font-family: sans-serif;
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
    <option value="1.16.1" selected>1.16.1</option>
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

import React from 'react'
import { createRoot } from 'react-dom/client'
import { ThemeProvider, CssBaseline, createTheme } from '@mui/material'
import { teal, green } from '@mui/material/colors'
import App from './App'

// TODO:
// - WelcomeScreen
// - VersionSelectionScreen
// - ModSelectionScreen
// - InstallConfirmScreen
// - ModpackFaq
// - Wire up all this logic with the actual Golang code.
// - Bundle Roboto font, don't depend on the internet to work properly.

const theme = createTheme({
  palette: {
    mode: 'dark',
    primary: { main: teal[700] },
    secondary: { main: green.A400 }
  }
})

createRoot(document.getElementById('app') as HTMLElement).render((
  <React.StrictMode>
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <App />
    </ThemeProvider>
  </React.StrictMode>
))

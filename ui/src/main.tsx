import React from 'react'
import { createRoot } from 'react-dom/client'
import { ThemeProvider, CssBaseline, createTheme } from '@mui/material'
import { teal, green } from '@mui/material/colors'
import App from './App'

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

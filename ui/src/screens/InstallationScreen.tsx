import { useEffect, useRef, useState } from 'react'
import { Button, Typography, LinearProgress, Dialog, DialogTitle, DialogContent, DialogContentText, DialogActions } from '@mui/material'

const InstallationScreen = ({ inProgress, setCurrentStep }: {
  inProgress: boolean
  setCurrentStep: React.Dispatch<React.SetStateAction<number>>
}): JSX.Element => {
  const [message, setMessage] = useState('Working...')
  const [error, setError] = useState('')
  const [query, setQuery] = useState('')
  window.setMessageState = setMessage
  window.setErrorState = setError
  window.setQueryState = setQuery

  const handleQueryResponse = (response: boolean): void => {
    setQuery('')
    window.respondQuery(response)
  }

  const firstLoadRef = useRef(true)
  useEffect(() => {
    if (firstLoadRef.current) {
      firstLoadRef.current = false
      window.installMods()
    }
  }, [])

  return (
    <>
      <Dialog open={query !== ''} maxWidth='xs' onClose={() => handleQueryResponse(false)}>
        <DialogTitle>Hey!</DialogTitle>
        <DialogContent>
          <DialogContentText>{query}</DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => handleQueryResponse(false)} color='inherit'>No</Button>
          <Button onClick={() => handleQueryResponse(true)} color='secondary' autoFocus>Yes</Button>
        </DialogActions>
      </Dialog>
      <Typography variant='h5' textAlign='center' gutterBottom>ibu's mod installer</Typography>
      <Typography variant='h6' color={error !== '' ? 'error' : 'inherit'} gutterBottom>
        {error === '' ? message : 'Error: ' + error}
      </Typography>
      {inProgress && <LinearProgress />}
      <div css={{ flex: 1 }} />
      {!inProgress && (
        <Button variant='outlined' color='secondary' onClick={() => setCurrentStep(1)}>
          Return to Welcome Screen
        </Button>
      )}
    </>
  )
}

export default InstallationScreen

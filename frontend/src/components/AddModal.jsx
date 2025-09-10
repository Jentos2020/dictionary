import { useState } from 'react';
import { Dialog, DialogTitle, DialogContent, DialogActions, Button, Typography } from '@mui/material';
import axios from 'axios';

function AddModal({ open, onClose, word, dictionary, onAddComplete }) {
  const [success, setSuccess] = useState(false);

  const handleAdd = async (e) => {
    e.preventDefault();
    try {
      await axios.post('/api/words', { 
        data: word, 
        dictionary: dictionary 
      }, {
        headers: { 'Content-Type': 'application/json' }
      });
      setSuccess(true);
    } catch (err) {
      console.log('Axios error:', err, err.response);
      alert('Ошибка при добавлении слова: ' + (err.response?.data || err.message));
    }
  };

  const handleClose = () => {
    if (success) {
      onAddComplete();
    }
    setSuccess(false);
    onClose();
  };

  return (
    <Dialog open={open} onClose={handleClose}>
      {success ? (
        <>
          <DialogTitle>Успех</DialogTitle>
          <DialogContent>
            <Typography>Слово "{word}" добавлено в словарь "{dictionary}".</Typography>
          </DialogContent>
          <DialogActions>
            <Button onClick={handleClose} color="primary">ОК</Button>
          </DialogActions>
        </>
      ) : (
        <>
          <DialogTitle>Добавить слово</DialogTitle>
          <DialogContent>
            <Typography>Добавить слово "{word}" в словарь "{dictionary}"?</Typography>
          </DialogContent>
          <DialogActions>
            <Button onClick={handleClose}>Отмена</Button>
            <Button onClick={handleAdd} color="primary">Добавить</Button>
          </DialogActions>
        </>
      )}
    </Dialog>
  );
}

export default AddModal;
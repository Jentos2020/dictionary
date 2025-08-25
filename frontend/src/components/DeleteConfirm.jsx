import { useState } from 'react';
import { Dialog, DialogTitle, DialogContent, DialogActions, Button, Typography } from '@mui/material';
import axios from 'axios';

function DeleteConfirm({ open, onClose, word, dictionary, onDeleteComplete }) {
  const [success, setSuccess] = useState(false);

  const handleDelete = async (e) => {
    e.preventDefault();
    try {
      await axios.delete(`/api/words/${word}`, {
        data: { dictionary },
      });
      setSuccess(true);
    } catch (err) {
      console.log('Axios error:', err, err.response);
      alert('Error deleting word: ' + (err.response?.data || err.message));
    }
  };

  const handleClose = () => {
    if (success) {
      onDeleteComplete();
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
            <Typography>Слово "{word}" было успешно удалено из словаря "{dictionary}".</Typography>
          </DialogContent>
          <DialogActions>
            <Button onClick={handleClose} color="primary">OK</Button>
          </DialogActions>
        </>
      ) : (
        <>
          <DialogTitle>Удалить слово</DialogTitle>
          <DialogContent>
            <Typography>Удалить "{word}" из словаря "{dictionary}"? Это действие нельзя отменить.</Typography>
          </DialogContent>
          <DialogActions>
            <Button onClick={handleClose}>Отмена</Button>
            <Button onClick={handleDelete} color="error">Удалить</Button>
          </DialogActions>
        </>
      )}
    </Dialog>
  );
}

export default DeleteConfirm;
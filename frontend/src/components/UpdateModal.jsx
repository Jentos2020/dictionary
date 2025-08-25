import { useState } from 'react';
import { Dialog, DialogTitle, DialogContent, DialogActions, Button, TextField, Typography } from '@mui/material';
import axios from 'axios';

function UpdateModal({ open, onClose, oldWord, dictionary, onUpdateComplete }) {
  const [newWord, setNewWord] = useState(oldWord);
  const [success, setSuccess] = useState(false);

  const handleUpdate = async (e) => {
    e.preventDefault();
    if (!newWord) {
      alert('Новое слово обязательно');
      return;
    }
    try {
      await axios.put(`/api/words/${oldWord}`, { data: newWord, dictionary });
      setSuccess(true);
    } catch (err) {
      console.log('Axios error:', err, err.response);
      alert('Error updating word: ' + (err.response?.data || err.message));
    }
  };

  const handleClose = () => {
    if (success) {
      onUpdateComplete();
    }
    setSuccess(false);
    setNewWord(oldWord); // Сбрасываем на исходное
    onClose();
  };

  return (
    <Dialog open={open} onClose={handleClose}>
      {success ? (
        <>
          <DialogTitle>Успех</DialogTitle>
          <DialogContent>
            <Typography>Слово "{oldWord}" было успешно изменено на "{newWord}" в словаре "{dictionary}".</Typography>
          </DialogContent>
          <DialogActions>
            <Button onClick={handleClose} color="primary">OK</Button>
          </DialogActions>
        </>
      ) : (
        <>
          <DialogTitle>Изменить слово</DialogTitle>
          <DialogContent>
            <TextField
              fullWidth
              label="Новое слово"
              value={newWord}
              onChange={(e) => setNewWord(e.target.value)}
              variant="outlined"
              sx={{ mt: 1 }}
            />
          </DialogContent>
          <DialogActions>
            <Button onClick={handleClose}>Отмена</Button>
            <Button onClick={handleUpdate} color="primary">Изменить</Button>
          </DialogActions>
        </>
      )}
    </Dialog>
  );
}

export default UpdateModal;
import { useState } from 'react';
import { Dialog, DialogTitle, DialogContent, DialogActions, Button, Typography } from '@mui/material';
import axios from 'axios';

function AddModal({ open, onClose, word, dictionary, onAddComplete }) {
  const [success, setSuccess] = useState(false);

  const handleAdd = async (e) => {
    e.preventDefault();
    try {
      await axios.post('/api/words', { data: word, dictionary });
      setSuccess(true); // Переключаем на сообщение об успехе
    } catch (err) {
      console.log('Axios error:', err, err.response);
      alert('Error adding word: ' + (err.response?.data || err.message)); // Оставляем alert для ошибок
    }
  };

  const handleClose = () => {
    if (success) {
      onAddComplete(); // Вызываем, только если успех
    }
    setSuccess(false); // Сбрасываем состояние
    onClose();
  };

  return (
    <Dialog open={open} onClose={handleClose}>
      {success ? (
        <>
          <DialogTitle>Успех</DialogTitle>
          <DialogContent>
            <Typography>Слово "{word}" было успешно добавлено в словарь "{dictionary}".</Typography>
          </DialogContent>
          <DialogActions>
            <Button onClick={handleClose} color="primary">OK</Button>
          </DialogActions>
        </>
      ) : (
        <>
          <DialogTitle>Добавить слово</DialogTitle>
          <DialogContent>
            <Typography>Добавить "{word}" в словарь "{dictionary}"?</Typography>
          </DialogContent>
          <DialogActions>
            <Button onClick={handleClose}>Отмена</Button>
            <Button onClick={handleAdd} color="primary">Подтвердить</Button>
          </DialogActions>
        </>
      )}
    </Dialog>
  );
}

export default AddModal;
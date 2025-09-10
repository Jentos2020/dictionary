import { useState } from 'react';
import { Dialog, DialogTitle, DialogContent, DialogActions, Button, Typography } from '@mui/material';
import axios from 'axios';

function DeleteConfirm({ open, onClose, word, dictionary, onDeleteComplete, isDefinition }) {
  const [success, setSuccess] = useState(false);

  const handleDelete = async (e) => {
    e.preventDefault();
    try {
      const endpoint = isDefinition 
        ? `/api/definitions/${word}?dictionary=${dictionary}`
        : `/api/words/${word}`;
      const config = {
        headers: { 'Content-Type': 'application/json' }
      };
      if (!isDefinition) {
        config.data = { dictionary: dictionary };
      }
      await axios.delete(endpoint, config);
      setSuccess(true);
    } catch (err) {
      console.log('Axios error:', err, err.response);
      alert(`Ошибка при удалении ${isDefinition ? 'определения' : 'слова'}: ` + (err.response?.data || err.message));
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
            <Typography>
              {isDefinition ? `Определение для "${word}" удалено из словаря "${dictionary}".` : `Слово "${word}" удалено из словаря "${dictionary}".`}
            </Typography>
          </DialogContent>
          <DialogActions>
            <Button onClick={handleClose} color="primary">ОК</Button>
          </DialogActions>
        </>
      ) : (
        <>
          <DialogTitle>{isDefinition ? 'Удалить определение' : 'Удалить слово'}</DialogTitle>
          <DialogContent>
            <Typography>
              {isDefinition
                ? `Удалить определение для "${word}" из словаря "${dictionary}"?`
                : `Удалить слово "${word}" из словаря "${dictionary}"?`} Это действие нельзя отменить.
            </Typography>
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
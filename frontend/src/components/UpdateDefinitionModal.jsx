import { useState } from 'react';
import { Dialog, DialogTitle, DialogContent, DialogActions, Button, TextField, Typography } from '@mui/material';
import axios from 'axios';

function UpdateDefinitionModal({ open, onClose, word, dictionary, onUpdateComplete }) {
  const [newDefinition, setNewDefinition] = useState('');
  const [success, setSuccess] = useState(false);

  const handleUpdate = async (e) => {
    e.preventDefault();
    if (!newDefinition) {
      alert('Определение обязательно');
      return;
    }
    try {
      await axios.put(`/api/definitions/${word}?dictionary=${dictionary}`, { definition: newDefinition });
      setSuccess(true);
    } catch (err) {
      console.log('Axios error:', err, err.response);
      alert('Ошибка при обновлении определения: ' + (err.response?.data || err.message));
    }
  };

  const handleClose = () => {
    if (success) {
      onUpdateComplete();
    }
    setSuccess(false);
    setNewDefinition('');
    onClose();
  };

  return (
    <Dialog open={open} onClose={handleClose}>
      {success ? (
        <>
          <DialogTitle>Успех</DialogTitle>
          <DialogContent>
            <Typography>Определение для "{word}" обновлено в словаре "{dictionary}".</Typography>
          </DialogContent>
          <DialogActions>
            <Button onClick={handleClose} color="primary">ОК</Button>
          </DialogActions>
        </>
      ) : (
        <>
          <DialogTitle>Изменить определение</DialogTitle>
          <DialogContent>
            <TextField
              fullWidth
              label="Новое определение"
              value={newDefinition}
              onChange={(e) => setNewDefinition(e.target.value)}
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

export default UpdateDefinitionModal;
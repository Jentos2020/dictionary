import { useState } from 'react';
import { Dialog, DialogTitle, DialogContent, DialogActions, Button, TextField, Typography } from '@mui/material';
import axios from 'axios';

function AddDefinitionModal({ open, onClose, word, dictionary, onAddComplete }) {
  const [definition, setDefinition] = useState('');
  const [success, setSuccess] = useState(false);

  const handleAdd = async (e) => {
    e.preventDefault();
    if (!definition) {
      alert('Определение обязательно');
      return;
    }
    try {
      await axios.post('/api/definitions', { word, definition, dictionary });
      setSuccess(true);
    } catch (err) {
      console.log('Axios error:', err, err.response);
      alert('Ошибка при добавлении определения: ' + (err.response?.data || err.message));
    }
  };

  const handleClose = () => {
    if (success) {
      onAddComplete();
    }
    setSuccess(false);
    setDefinition('');
    onClose();
  };

  return (
    <Dialog open={open} onClose={handleClose}>
      {success ? (
        <>
          <DialogTitle>Успех</DialogTitle>
          <DialogContent>
            <Typography>Определение для "{word}" добавлено в словарь "{dictionary}".</Typography>
          </DialogContent>
          <DialogActions>
            <Button onClick={handleClose} color="primary">ОК</Button>
          </DialogActions>
        </>
      ) : (
        <>
          <DialogTitle>Добавить определение</DialogTitle>
          <DialogContent>
            <Typography>Добавить определение для "{word}" в словарь "{dictionary}":</Typography>
            <TextField
              fullWidth
              label="Определение"
              value={definition}
              onChange={(e) => setDefinition(e.target.value)}
              variant="outlined"
              sx={{ mt: 1 }}
            />
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

export default AddDefinitionModal;
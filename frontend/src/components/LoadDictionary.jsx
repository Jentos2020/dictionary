import { useState } from 'react';
import { TextField, Button, Box, FormControl, InputLabel, Select, MenuItem } from '@mui/material';
import axios from 'axios';

function LoadDictionary() {
  const [name, setName] = useState('');

  const handleLoad = async () => {
    if (!name) return alert('Название словаря обязательно');
    try {
      await axios.post(`/api/dictionaries/add/${name}`);
      alert('Словарь успешно загружен!');
      setName('');
    } catch (err) {
      console.log('Axios error:', err, err.response);
      alert('Ошибка при загрузке словаря: ' + (err.response?.data || err.message));
    }
  };

  return (
    <Box display="flex" alignItems="center" gap={1}>
      <FormControl sx={{ flexGrow: 1 }}>
        <TextField
          label="Название словаря"
          value={name}
          onChange={(e) => setName(e.target.value)}
          variant="outlined"
        />
      </FormControl>
      <FormControl sx={{ minWidth: 120 }}>
        <InputLabel id="dictionary-select-label">Выбрать словарь</InputLabel>
        <Select
          labelId="dictionary-select-label"
          value=""
          label="Выбрать словарь"
          onChange={(e) => setName(e.target.value)}
        >
          <MenuItem value="">Ручной ввод</MenuItem>
          <MenuItem value="russian">Русский</MenuItem>
          <MenuItem value="english">Английский</MenuItem>
        </Select>
      </FormControl>
      <Button variant="contained" onClick={handleLoad}>
        Загрузить
      </Button>
    </Box>
  );
}

export default LoadDictionary;
import { useState } from 'react';
import { TextField, Button, Box } from '@mui/material';
import axios from 'axios';

function LoadDictionary() {
  const [name, setName] = useState('');

  const handleLoad = async () => {
    if (!name) return alert('Name required');
    try {
      // API call (update port if needed)
      await axios.post(`/api/dictionaries/add/${name}`);
      alert('Dictionary loaded successfully!');
      setName('');
    } catch (err) {
      console.log('Axios error:', err, err.response);
      alert('Error loading dictionary: ' + (err.response?.data || err.message));
    }
  };

  return (
    <Box display="flex" alignItems="center">
      <TextField
        label="Название словаря (например, russian)"
        value={name}
        onChange={(e) => setName(e.target.value)}
        variant="outlined"
        sx={{ mr: 1, flexGrow: 1 }}
      />
      <Button variant="contained" onClick={handleLoad}>
        Обновить
      </Button>
    </Box>
  );
}

export default LoadDictionary;
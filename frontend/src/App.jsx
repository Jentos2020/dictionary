import { useState } from 'react';
import { Container, Typography, FormControl, InputLabel, Select, MenuItem, Box } from '@mui/material';
import SearchInput from './components/SearchInput.jsx';
import LoadDictionary from './components/LoadDictionary.jsx';

function App() {
  const [dictionary, setDictionary] = useState('russian');

  return (
    <Box
      sx={{
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        minHeight: '100vh',
        width: '100vw',
        bgcolor: 'background.default',
        textAlign: 'center',
      }}
    >
      <Container maxWidth="md" sx={{ py: 4, mx: 'auto' }}>
        <Typography variant="h4" gutterBottom>
          Онлайн-словарь
        </Typography>

        <FormControl fullWidth sx={{ mb: 2 }}>
          <InputLabel id="dictionary-select-label">Словарь</InputLabel>
          <Select
            labelId="dictionary-select-label"
            value={dictionary}
            label="Словарь"
            onChange={(e) => setDictionary(e.target.value)}
          >
            <MenuItem value="russian">Русский</MenuItem>
            <MenuItem value="english">Английский</MenuItem>
          </Select>
        </FormControl>

        <SearchInput dictionary={dictionary} />
        <Box sx={{ mt: 4 }}>
          <LoadDictionary />
        </Box>
      </Container>
    </Box>
  );
}

export default App;
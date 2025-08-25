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
        justifyContent: 'center', // Центрирование по горизонтали
        alignItems: 'center',     // Центрирование по вертикали
        minHeight: '100vh',       // Полная высота экрана
        width: '100vw',           // Полная ширина экрана (добавлено для фикса в некоторых браузерах)
        bgcolor: 'background.default',
        textAlign: 'center',      // Центрирование текста внутри
      }}
    >
      <Container maxWidth="xs" sx={{ py: 4, mx: 'auto' }}> {/* mx: 'auto' явно центрирует контейнер */}
        <Typography variant="h4" gutterBottom>
          Online Dictionary
        </Typography>

        <FormControl fullWidth sx={{ mb: 2 }}>
          <InputLabel id="dictionary-select-label">Словарь</InputLabel>
          <Select
            labelId="dictionary-select-label"
            value={dictionary}
            label="Dictionary"
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
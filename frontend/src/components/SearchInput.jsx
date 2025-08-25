import { useState, useEffect, useRef } from 'react';
import { TextField, Paper, List, ListItem, ListItemText, Button, Box } from '@mui/material';
import AddModal from './AddModal.jsx';
import WordActions from './WordActions.jsx';

function SearchInput({ dictionary }) {
  const [input, setInput] = useState('');
  const [suggestions, setSuggestions] = useState([]);
  const [selectedWord, setSelectedWord] = useState(null);
  const [showAddModal, setShowAddModal] = useState(false);
  const ws = useRef(null);
  const reconnectTimeout = useRef(null);

  useEffect(() => {
    const connectWebSocket = () => {
      // Проверяем, нет ли уже активного соединения
      if (ws.current && ws.current.readyState === WebSocket.OPEN) {
        console.log('WebSocket already connected');
        return;
      }

      ws.current = new WebSocket('ws://127.0.0.1:5173/ws/search');

      ws.current.onopen = () => {
        console.log('WebSocket connected');
        clearTimeout(reconnectTimeout.current); // Очищаем таймер переподключения
      };

      ws.current.onclose = () => {
        console.log('WebSocket closed, retrying in 1s...');
        reconnectTimeout.current = setTimeout(connectWebSocket, 1000);
      };

      ws.current.onerror = (err) => {
        console.error('WebSocket error:', err);
        ws.current.close(); // Закрываем при ошибке, чтобы сработал onclose
      };

      ws.current.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data);
          setSuggestions(data.words || []);
        } catch (err) {
          console.error('Invalid WS message:', err);
        }
      };
    };

    connectWebSocket();

    return () => {
      if (ws.current) {
        ws.current.close();
      }
      clearTimeout(reconnectTimeout.current);
    };
  }, []);

  const handleInputChange = (e) => {
    const value = e.target.value;
    setInput(value);
    setSelectedWord(null);

    if (value && ws.current?.readyState === WebSocket.OPEN) {
      ws.current.send(JSON.stringify({ prefix: value }));
    } else {
      setSuggestions([]);
    }
  };

  const handleSelectSuggestion = (word) => {
    setInput(word.data);
    setSuggestions([]);
    setSelectedWord(word.data);
  };

  const handleKeyDown = (e) => {
    if (e.key === 'Enter' && input) {
      setSuggestions([]);
      setSelectedWord(input);
    } else if (e.key === 'Escape') {
      setSuggestions([]);
    }
  };

  return (
    <Box>
      <TextField
        fullWidth
        label="Поиск слова"
        value={input}
        onChange={handleInputChange}
        onKeyDown={handleKeyDown}
        variant="outlined"
      />
      {suggestions.length > 0 && (
        <Paper elevation={3} sx={{ mt: 1, maxHeight: 200, overflow: 'auto' }}>
          <List>
            {suggestions.map((word, index) => (
              <ListItem
                key={index}
                onClick={() => handleSelectSuggestion(word)}
              >
                <ListItemText primary={word.data} />
              </ListItem>
            ))}
          </List>
        </Paper>
      )}
      <Button
        variant="contained"
        color="primary"
        sx={{ mt: 1 }}
        disabled={!input}
        onClick={() => setShowAddModal(true)}
      >
        Добавить слово
      </Button>
      {selectedWord && (
        <WordActions
          word={selectedWord}
          dictionary={dictionary}
          onActionComplete={() => {
            setInput('');
            setSelectedWord(null);
          }}
        />
      )}
      <AddModal
        open={showAddModal}
        onClose={() => setShowAddModal(false)}
        word={input}
        dictionary={dictionary}
        onAddComplete={() => {
          setInput('');
          setSelectedWord(null);
        }}
      />
    </Box>
  );
}

export default SearchInput;
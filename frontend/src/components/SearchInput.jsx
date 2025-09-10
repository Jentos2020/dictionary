import { useState, useEffect, useRef } from 'react';
import { TextField, Paper, List, ListItem, ListItemText, Button, Box, Typography } from '@mui/material';
import AddDefinitionModal from './AddDefinitionModal.jsx';
import AddModal from './AddModal.jsx';
import WordActions from './WordActions.jsx';

function SearchInput({ dictionary }) {
  const [input, setInput] = useState('');
  const [suggestions, setSuggestions] = useState([]);
  const [selectedWord, setSelectedWord] = useState(null);
  const [definition, setDefinition] = useState(null);
  const [definitionError, setDefinitionError] = useState(null);
  const [showAddDefinitionModal, setShowAddDefinitionModal] = useState(false);
  const [showAddWordModal, setShowAddWordModal] = useState(false);
  const ws = useRef(null);
  const reconnectTimeout = useRef(null);

  useEffect(() => {
    const connectWebSocket = () => {
      if (ws.current && ws.current.readyState === WebSocket.OPEN) {
        console.log('WebSocket already connected');
        return;
      }

      ws.current = new WebSocket('ws://127.0.0.1:5173/ws/search');

      ws.current.onopen = () => {
        console.log('WebSocket connected');
        clearTimeout(reconnectTimeout.current);
      };

      ws.current.onclose = () => {
        console.log('WebSocket closed, retrying in 1s...');
        reconnectTimeout.current = setTimeout(connectWebSocket, 1000);
      };

      ws.current.onerror = (err) => {
        console.error('WebSocket error:', err);
        ws.current.close();
      };

      ws.current.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data);
          if (data.type === 'definition') {
            setDefinitionError(null);
            if (data.error) {
              setDefinitionError(data.error);
              setDefinition(null);
            } else {
              setDefinition(data.definition);
            }
          } else {
            setSuggestions(data.words || []);
          }
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
    setDefinition(null);
    setDefinitionError(null);

    if (value && ws.current?.readyState === WebSocket.OPEN) {
      ws.current.send(JSON.stringify({ type: 'search', prefix: value }));
    } else {
      setSuggestions([]);
    }
  };

  const handleSelectSuggestion = (word) => {
    setInput(word.data);
    setSuggestions([]);
    setSelectedWord(word.data);
    if (ws.current?.readyState === WebSocket.OPEN) {
      ws.current.send(JSON.stringify({ type: 'get_definition', word: word.data, dictionary }));
    }
  };

  const handleKeyDown = (e) => {
    if (e.key === 'Enter' && input) {
      setSuggestions([]);
      setSelectedWord(input);
      if (ws.current?.readyState === WebSocket.OPEN) {
        ws.current.send(JSON.stringify({ type: 'get_definition', word: input, dictionary }));
      }
    } else if (e.key === 'Escape') {
      setSuggestions([]);
      setDefinition(null);
      setDefinitionError(null);
    }
  };

  const isWordFound = suggestions.length > 0 || selectedWord;

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
                button
                onClick={() => handleSelectSuggestion(word)}
              >
                <ListItemText primary={word.data} />
              </ListItem>
            ))}
          </List>
        </Paper>
      )}
      {definition && (
        <Box sx={{ mt: 2, textAlign: 'left' }}>
          <Typography variant="subtitle1" fontWeight="bold">Определение:</Typography>
          <Typography variant="body1">{definition}</Typography>
        </Box>
      )}
      {definitionError && (
        <Box sx={{ mt: 2, textAlign: 'left' }}>
          <Typography variant="body1" color="error">
            Определение не найдено
          </Typography>
        </Box>
      )}
      <Box sx={{ mt: 1, display: 'flex', gap: 1 }}>
        {!isWordFound && (
          <Button
            variant="contained"
            color="primary"
            disabled={!input}
            onClick={() => setShowAddWordModal(true)}
          >
            Добавить слово
          </Button>
        )}
        {definitionError === 'Definition not found' && (
          <Button
            variant="contained"
            color="primary"
            disabled={!input}
            onClick={() => setShowAddDefinitionModal(true)}
          >
            Добавить определение
          </Button>
        )}
      </Box>
      {selectedWord && (
        <WordActions
          word={selectedWord}
          dictionary={dictionary}
          onActionComplete={() => {
            setInput('');
            setSelectedWord(null);
            setDefinition(null);
            setDefinitionError(null);
          }}
        />
      )}
      <AddModal
        open={showAddWordModal}
        onClose={() => setShowAddWordModal(false)}
        word={input}
        dictionary={dictionary}
        onAddComplete={() => {
          setInput('');
          setSelectedWord(null);
          setDefinition(null);
          setDefinitionError(null);
        }}
      />
      <AddDefinitionModal
        open={showAddDefinitionModal}
        onClose={() => setShowAddDefinitionModal(false)}
        word={input}
        dictionary={dictionary}
        onAddComplete={() => {
          setInput('');
          setSelectedWord(null);
          setDefinition(null);
          setDefinitionError(null);
        }}
      />
    </Box>
  );
}

export default SearchInput;
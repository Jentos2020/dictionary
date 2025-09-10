import { useState } from 'react';
import { Button, Box } from '@mui/material';
import UpdateDefinitionModal from './UpdateDefinitionModal.jsx';
import DeleteConfirm from './DeleteConfirm.jsx';
import UpdateModal from './UpdateModal.jsx';

function WordActions({ word, dictionary, onActionComplete }) {
  const [showUpdateDefinitionModal, setShowUpdateDefinitionModal] = useState(false);
  const [showDeleteDefinitionConfirm, setShowDeleteDefinitionConfirm] = useState(false);
  const [showUpdateWordModal, setShowUpdateWordModal] = useState(false);
  const [showDeleteWordConfirm, setShowDeleteWordConfirm] = useState(false);

  return (
    <Box sx={{ mt: 2, display: 'flex', gap: 1 }}>
      <Button
        variant="outlined"
        color="primary"
        onClick={() => setShowUpdateWordModal(true)}
      >
        Изменить слово
      </Button>
      <Button
        variant="outlined"
        color="error"
        onClick={() => setShowDeleteWordConfirm(true)}
      >
        Удалить слово
      </Button>
      <Button
        variant="outlined"
        color="primary"
        onClick={() => setShowUpdateDefinitionModal(true)}
      >
        Изменить определение
      </Button>
      <Button
        variant="outlined"
        color="error"
        onClick={() => setShowDeleteDefinitionConfirm(true)}
      >
        Удалить определение
      </Button>
      <UpdateModal
        open={showUpdateWordModal}
        onClose={() => setShowUpdateWordModal(false)}
        word={word}
        dictionary={dictionary}
        onUpdateComplete={onActionComplete}
      />
      <DeleteConfirm
        open={showDeleteWordConfirm}
        onClose={() => setShowDeleteWordConfirm(false)}
        word={word}
        dictionary={dictionary}
        onDeleteComplete={onActionComplete}
        isDefinition={false}
      />
      <UpdateDefinitionModal
        open={showUpdateDefinitionModal}
        onClose={() => setShowUpdateDefinitionModal(false)}
        word={word}
        dictionary={dictionary}
        onUpdateComplete={onActionComplete}
      />
      <DeleteConfirm
        open={showDeleteDefinitionConfirm}
        onClose={() => setShowDeleteDefinitionConfirm(false)}
        word={word}
        dictionary={dictionary}
        onDeleteComplete={onActionComplete}
        isDefinition={true}
      />
    </Box>
  );
}

export default WordActions;
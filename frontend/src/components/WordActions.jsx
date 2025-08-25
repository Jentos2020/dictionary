import { useState } from 'react';
import { Button, Box } from '@mui/material';
import UpdateModal from './UpdateModal.jsx';
import DeleteConfirm from './DeleteConfirm.jsx';

function WordActions({ word, dictionary, onActionComplete }) {
  const [showUpdateModal, setShowUpdateModal] = useState(false);
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false);

  return (
    <Box sx={{ mt: 2 }}>
      <Button
        variant="outlined"
        color="primary"
        sx={{ mr: 1 }}
        onClick={() => setShowUpdateModal(true)}
      >
        Update Word
      </Button>
      <Button
        variant="outlined"
        color="error"
        onClick={() => setShowDeleteConfirm(true)}
      >
        Delete Word
      </Button>
      <UpdateModal
        open={showUpdateModal}
        onClose={() => setShowUpdateModal(false)}
        oldWord={word}
        dictionary={dictionary}
        onUpdateComplete={onActionComplete}
      />
      <DeleteConfirm
        open={showDeleteConfirm}
        onClose={() => setShowDeleteConfirm(false)}
        word={word}
        dictionary={dictionary}
        onDeleteComplete={onActionComplete}
      />
    </Box>
  );
}

export default WordActions;
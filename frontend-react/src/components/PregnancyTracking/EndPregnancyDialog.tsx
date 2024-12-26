import React, { useState } from 'react';
import { Modal, Button, Stack, Select, Textarea } from '@mantine/core';

interface EndPregnancyDialogProps {
  opened: boolean;
  onClose: () => void;
  onSubmit: (outcome: string) => void;
}

const outcomes = [
  { value: 'SUCCESSFUL_BIRTH', label: 'Successful Birth' },
  { value: 'STILLBIRTH', label: 'Stillbirth' },
  { value: 'MISCARRIAGE', label: 'Miscarriage' },
  { value: 'MEDICAL_TERMINATION', label: 'Medical Termination' },
  { value: 'OTHER', label: 'Other' }
];

const EndPregnancyDialog: React.FC<EndPregnancyDialogProps> = ({
  opened,
  onClose,
  onSubmit
}) => {
  const [outcome, setOutcome] = useState<string | null>(null);
  const [notes, setNotes] = useState('');

  const handleSubmit = () => {
    if (outcome) {
      onSubmit(outcome);
      setOutcome(null);
      setNotes('');
    }
  };

  return (
    <Modal
      opened={opened}
      onClose={onClose}
      title="End Pregnancy Tracking"
      size="sm"
    >
      <Stack>
        <Select
          label="Pregnancy Outcome"
          placeholder="Select outcome"
          data={outcomes}
          value={outcome}
          onChange={setOutcome}
          required
        />

        <Textarea
          label="Additional Notes"
          placeholder="Enter any additional notes"
          value={notes}
          onChange={(event) => setNotes(event.currentTarget.value)}
          minRows={3}
        />

        <Button
          onClick={handleSubmit}
          disabled={!outcome}
          color="red"
          fullWidth
        >
          End Tracking
        </Button>
        <Button onClick={onClose} fullWidth>
          Cancel
        </Button>
      </Stack>
    </Modal>
  );
};

export default EndPregnancyDialog;

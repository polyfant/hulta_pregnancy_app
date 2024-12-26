import React, { useState } from 'react';
import { Modal, Button, Stack } from '@mantine/core';
import { DateInput } from '@mantine/dates';

interface StartPregnancyDialogProps {
  opened: boolean;
  onClose: () => void;
  onSubmit: (conceptionDate: string) => void;
}

const StartPregnancyDialog: React.FC<StartPregnancyDialogProps> = ({
  opened,
  onClose,
  onSubmit
}) => {
  const [conceptionDate, setConceptionDate] = useState<Date | null>(null);

  const handleSubmit = () => {
    if (conceptionDate) {
      onSubmit(conceptionDate.toISOString().split('T')[0]);
    }
  };

  return (
    <Modal 
      opened={opened}
      onClose={onClose}
      title="Start Pregnancy Tracking"
      size="sm"
    >
      <Stack>
        <DateInput
          label="Conception Date"
          placeholder="Select date"
          value={conceptionDate}
          onChange={setConceptionDate}
          maxDate={new Date()}
        />

        <Button
          onClick={onClose}
          variant="subtle"
          fullWidth
        >
          Cancel
        </Button>

        <Button
          onClick={handleSubmit}
          disabled={!conceptionDate}
          fullWidth
        >
          Start Tracking
        </Button>
      </Stack>
    </Modal>
  );
};

export default StartPregnancyDialog;

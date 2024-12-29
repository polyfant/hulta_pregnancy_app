import { Modal, Stack, Button, Text, Select, Textarea } from '@mantine/core';
import { useState } from 'react';

interface EndPregnancyDialogProps {
  opened: boolean;
  onClose: () => void;
  onConfirm: () => void;
  isLoading: boolean;
}

export function EndPregnancyDialog({
  opened,
  onClose,
  onConfirm,
  isLoading
}: EndPregnancyDialogProps) {
  const [reason, setReason] = useState<string | null>(null);
  const [notes, setNotes] = useState('');

  const handleConfirm = () => {
    onConfirm();
  };

  return (
    <Modal
      opened={opened}
      onClose={onClose}
      title="End Pregnancy Tracking"
      size="md"
    >
      <Stack>
        <Text size="sm" c="dimmed">
          Are you sure you want to end pregnancy tracking for this mare?
          This action cannot be undone.
        </Text>

        <Select
          label="Reason"
          placeholder="Select reason"
          data={[
            { value: 'foaling', label: 'Successful Foaling' },
            { value: 'loss', label: 'Pregnancy Loss' },
            { value: 'error', label: 'Tracking Error' },
            { value: 'other', label: 'Other' }
          ]}
          value={reason}
          onChange={setReason}
          required
        />

        <Textarea
          label="Additional Notes"
          placeholder="Enter any additional notes or observations"
          value={notes}
          onChange={(event) => setNotes(event.currentTarget.value)}
          minRows={3}
        />

        <Button.Group>
          <Button
            variant="light"
            onClick={onClose}
            style={{ flex: 1 }}
          >
            Cancel
          </Button>
          <Button
            color="red"
            onClick={handleConfirm}
            loading={isLoading}
            disabled={!reason}
            style={{ flex: 1 }}
          >
            End Tracking
          </Button>
        </Button.Group>
      </Stack>
    </Modal>
  );
}

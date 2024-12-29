import {
  Modal,
  Text,
  Stack,
  Button,
  Group
} from '@mantine/core';
import { Horse } from '../types/horse';
import { HorsePreviewCard } from './HorsePreviewCard';

interface ParentChangeDialogProps {
  opened: boolean;
  onClose: () => void;
  onConfirm: () => void;
  parentType: 'mother' | 'father';
  currentParent?: Horse | null;
  newParent?: Horse | null;
  currentExternalParent?: string;
  newExternalParent?: string;
}

export function ParentChangeDialog({
  opened,
  onClose,
  onConfirm,
  parentType,
  currentParent,
  newParent,
  currentExternalParent,
  newExternalParent,
}: ParentChangeDialogProps) {
  const parentTitle = parentType === 'mother' ? 'Mother' : 'Father';

  return (
    <Modal 
      opened={opened} 
      onClose={onClose} 
      title={`Change ${parentTitle}`}
      size="lg"
    >
      <Stack gap="md">
        <Text>
          Are you sure you want to change this horse's {parentType}?
        </Text>

        <Stack gap="xs">
          <Text fw={500} size="sm" c="dimmed">
            Current {parentTitle}:
          </Text>
          <HorsePreviewCard
            horse={currentParent}
            externalName={currentExternalParent}
            role={parentType}
          />
        </Stack>

        <Stack gap="xs">
          <Text fw={500} size="sm" c="dimmed">
            New {parentTitle}:
          </Text>
          <HorsePreviewCard
            horse={newParent}
            externalName={newExternalParent}
            role={parentType}
          />
        </Stack>

        <Group justify="flex-end" mt="xl">
          <Button variant="default" onClick={onClose}>
            Cancel
          </Button>
          <Button onClick={onConfirm}>
            Confirm Change
          </Button>
        </Group>
      </Stack>
    </Modal>
  );
}

import { Card, Checkbox, Text, Stack, Button, TextInput } from '@mantine/core';
import { Plus } from '@phosphor-icons/react';
import { FC, useState } from 'react';

interface ChecklistItem {
  id: string;
  text: string;
  completed: boolean;
}

export const PrefoalingChecklist: FC = () => {
  const [items, setItems] = useState<ChecklistItem[]>([
    { id: '1', text: 'Prepare foaling kit', completed: false },
    { id: '2', text: 'Clean stall thoroughly', completed: false },
    { id: '3', text: 'Check veterinarian availability', completed: false },
    { id: '4', text: 'Prepare fresh bedding', completed: false },
  ]);

  const [newItemText, setNewItemText] = useState('');

  const handleToggleItem = (itemId: string) => {
    setItems(prevItems =>
      prevItems.map(item =>
        item.id === itemId ? { ...item, completed: !item.completed } : item
      )
    );
  };

  const handleAddItem = () => {
    if (newItemText.trim()) {
      setItems(prevItems => [
        ...prevItems,
        {
          id: Date.now().toString(),
          text: newItemText.trim(),
          completed: false
        }
      ]);
      setNewItemText('');
    }
  };

  return (
    <Card shadow="sm" padding="lg" radius="md">
      <Stack gap="md">
        <Text fw={700} size="xl">Pre-foaling Checklist</Text>
        
        {items.map((item) => (
          <Checkbox
            key={item.id}
            label={item.text}
            checked={item.completed}
            onChange={() => handleToggleItem(item.id)}
            styles={(theme) => ({
              label: {
                textDecoration: item.completed ? 'line-through' : 'none',
                color: item.completed ? theme.colors.dark[3] : theme.colors.dark[0],
                transition: 'all 0.2s ease',
              }
            })}
          />
        ))}
        
        <Stack gap="sm">
          <TextInput
            placeholder="Add new checklist item"
            value={newItemText}
            onChange={(e) => setNewItemText(e.currentTarget.value)}
            onKeyPress={(e) => e.key === 'Enter' && handleAddItem()}
          />
          <Button 
            variant="filled" 
            color="brand"
            leftSection={<Plus size={20} />}
            onClick={handleAddItem}
            disabled={!newItemText.trim()}
          >
            Add New Item
          </Button>
        </Stack>
      </Stack>
    </Card>
  );
};
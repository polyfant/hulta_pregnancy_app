import React from 'react';
import { Card, Stack, Text, Checkbox, Group, Button, Skeleton } from '@mantine/core';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useHorsesApi } from '../../api/horses';
import { notifications } from '@mantine/notifications';

interface ChecklistItem {
  id: number;
  title: string;
  description: string;
  completed: boolean;
}

interface PrefoalingChecklistProps {
  horseId: number;
}

export const PrefoalingChecklist: React.FC<PrefoalingChecklistProps> = ({ horseId }) => {
  const api = useHorsesApi();
  const queryClient = useQueryClient();

  const { data: checklistItems, isLoading, error } = useQuery<ChecklistItem[]>({
    queryKey: ['prefoalingChecklist', horseId],
    queryFn: () => api.getPrefoalingChecklist(horseId),
    staleTime: 5 * 60 * 1000, // 5 minutes
    refetchOnWindowFocus: false,
  });

  const updateItemMutation = useMutation({
    mutationFn: (item: ChecklistItem) => api.updatePrefoalingChecklistItem(item),
    onMutate: async (updatedItem) => {
      // Cancel any outgoing refetches
      await queryClient.cancelQueries({ queryKey: ['prefoalingChecklist', horseId] });

      // Snapshot the previous value
      const previousItems = queryClient.getQueryData<ChecklistItem[]>(['prefoalingChecklist', horseId]);

      // Optimistically update to the new value
      queryClient.setQueryData(['prefoalingChecklist', horseId], (old: ChecklistItem[] | undefined) => 
        old?.map(item => item.id === updatedItem.id ? updatedItem : item) || []
      );

      // Return a context object with the snapshotted value
      return { previousItems };
    },
    onError: (err, newTodo, context) => {
      queryClient.setQueryData(['prefoalingChecklist', horseId], context?.previousItems);
      notifications.show({
        title: 'Error',
        message: 'Failed to update checklist item',
        color: 'red',
      });
    },
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: ['prefoalingChecklist', horseId] });
    },
  });

  const handleItemToggle = (item: ChecklistItem) => {
    updateItemMutation.mutate({ ...item, completed: !item.completed });
  };

  if (isLoading) {
    return (
      <Card shadow="sm" padding="lg" radius="md">
        <Stack gap="lg">
          <Skeleton height={30} width="60%" />
          <Stack gap="md">
            {[1, 2, 3, 4].map((item) => (
              <Skeleton key={item} height={40} />
            ))}
          </Stack>
        </Stack>
      </Card>
    );
  }

  if (error || !checklistItems || checklistItems.length === 0) {
    return (
      <Card shadow="sm" padding="lg" radius="md">
        <Text c="red">No pre-foaling checklist items found</Text>
      </Card>
    );
  }

  return (
    <Card shadow="sm" padding="lg" radius="md">
      <Stack gap="lg">
        <Text fw={700} size="xl">Pre-foaling Checklist</Text>
        
        <Stack gap="xs">
          {checklistItems.map((item) => (
            <Group key={item.id} justify="space-between" align="center">
              <Checkbox
                label={item.title}
                description={item.description}
                checked={item.completed}
                onChange={() => handleItemToggle(item)}
              />
            </Group>
          ))}
        </Stack>
      </Stack>
    </Card>
  );
};

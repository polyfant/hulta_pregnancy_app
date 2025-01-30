import { Stack, Text, Title, Button } from '@mantine/core';
import { Horse, Plus } from '@phosphor-icons/react';

interface EmptyStateProps {
    title?: string;
    message?: string;
    actionLabel?: string;
    onAction?: () => void;
    icon?: React.ReactNode;
}

export function EmptyState({ 
    title = 'No Data Found',
    message = 'There\'s nothing here yet.',
    actionLabel,
    onAction,
    icon = <Horse size="3rem" weight="thin" />
}: EmptyStateProps) {
    return (
        <Stack align="center" spacing="md" py="xl">
            {icon}
            <Title order={3}>{title}</Title>
            <Text c="dimmed" ta="center" maw={400}>
                {message}
            </Text>
            {actionLabel && onAction && (
                <Button
                    variant="light"
                    leftSection={<Plus size="1rem" />}
                    onClick={onAction}
                >
                    {actionLabel}
                </Button>
            )}
        </Stack>
    );
} 
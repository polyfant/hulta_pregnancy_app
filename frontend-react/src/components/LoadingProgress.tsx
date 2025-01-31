import { Progress, Text, Stack } from '@mantine/core';

export function LoadingProgress({ message }: { message: string }) {
    return (
        <Stack spacing="xs">
            <Text size="sm">{message}</Text>
            <Progress 
                value={100} 
                animated 
                color="blue"
                size="sm"
            />
        </Stack>
    );
} 
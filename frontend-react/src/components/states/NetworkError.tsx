import { Alert, Button, Stack, Text } from '@mantine/core';
import { Warning, ArrowClockwise } from '@phosphor-icons/react';

interface NetworkErrorProps {
    message?: string;
    onRetry?: () => void;
}

export function NetworkError({ message, onRetry }: NetworkErrorProps) {
    return (
        <Alert 
            icon={<Warning size="1.5rem" />} 
            title="Connection Error" 
            color="red"
            variant="light"
        >
            <Stack spacing="md">
                <Text size="sm">
                    {message || 'Unable to connect to the server. Please check your connection.'}
                </Text>
                {onRetry && (
                    <Button 
                        onClick={onRetry}
                        variant="light"
                        color="red"
                        leftSection={<ArrowClockwise size="1rem" />}
                    >
                        Try Again
                    </Button>
                )}
            </Stack>
        </Alert>
    );
} 
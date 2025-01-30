import { Button, Modal, Text } from '@mantine/core';
import { useDisclosure } from '@mantine/hooks';
import { notifications } from '@mantine/notifications';
import { Trash, Check, X } from '@phosphor-icons/react';

export function DeleteHorseButton({ horseId, horseName }: { horseId: string, horseName: string }) {
    const [opened, { open, close }] = useDisclosure(false);

    const handleDelete = async () => {
        try {
            await apiClient.delete(`/api/horses/${horseId}`);
            notifications.show({
                title: 'Success',
                message: `${horseName} has been deleted`,
                color: 'green',
                icon: <Check weight="bold" />,
            });
            close();
        } catch (error) {
            notifications.show({
                title: 'Error',
                message: 'Failed to delete horse',
                color: 'red',
                icon: <X weight="bold" />,
            });
        }
    };

    return (
        <>
            <Button 
                color="red" 
                leftSection={<Trash />}
                onClick={open}
            >
                Delete Horse
            </Button>

            <Modal opened={opened} onClose={close} title="Confirm Deletion">
                <Text mb="xl">
                    Are you sure you want to delete {horseName}? This action cannot be undone.
                </Text>
                <Button 
                    color="red" 
                    onClick={handleDelete}
                    mr="sm"
                >
                    Yes, Delete
                </Button>
                <Button variant="default" onClick={close}>
                    Cancel
                </Button>
            </Modal>
        </>
    );
} 
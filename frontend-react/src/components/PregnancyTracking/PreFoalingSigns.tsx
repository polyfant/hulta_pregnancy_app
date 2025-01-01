import {
	Button,
	Card,
	Group,
	Modal,
	Stack,
	Text,
	TextInput,
	Textarea,
	Title,
} from '@mantine/core';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { useState } from 'react';

import { notifications } from '@mantine/notifications';
import { Plus } from '@phosphor-icons/react';

interface PreFoalingSignsProps {
	horseId: string;
}

interface PreFoalingSign {
	id: string;
	description: string;
	date: string;
	notes?: string;
}

export default function PreFoalingSigns({ horseId }: PreFoalingSignsProps) {
	const [isModalOpen, setIsModalOpen] = useState(false);
	const [newSign, setNewSign] = useState<Partial<PreFoalingSign>>({});
	const queryClient = useQueryClient();

	const { data: signs = [], isLoading } = useQuery<PreFoalingSign[]>({
		queryKey: ['preFoalingSigns', horseId],
		queryFn: async () => {
			const response = await fetch(
				`/api/horses/${horseId}/pre-foaling-signs`
			);
			if (!response.ok)
				throw new Error('Failed to fetch pre-foaling signs');
			return response.json();
		},
	});

	const addSignMutation = useMutation({
		mutationFn: async (sign: Partial<PreFoalingSign>) => {
			const response = await fetch(
				`/api/horses/${horseId}/pre-foaling-signs`,
				{
					method: 'POST',
					headers: { 'Content-Type': 'application/json' },
					body: JSON.stringify(sign),
				}
			);
			if (!response.ok) throw new Error('Failed to add pre-foaling sign');
			return response.json();
		},
		onSuccess: () => {
			queryClient.invalidateQueries({
				queryKey: ['preFoalingSigns', horseId],
			});
			notifications.show({
				title: 'Success',
				message: 'Pre-foaling sign added',
				color: 'green',
			});
			setIsModalOpen(false);
			setNewSign({});
		},
	});

	const handleSubmit = () => {
		if (!newSign.description) {
			notifications.show({
				title: 'Error',
				message: 'Description is required',
				color: 'red',
			});
			return;
		}
		addSignMutation.mutate(newSign);
	};

	return (
		<Stack gap='lg'>
			<Group justify='space-between' mb='md'>
				<Title order={3}>Pre-Foaling Signs</Title>
				<Button
					onClick={() => setIsModalOpen(true)}
					leftSection={<Plus size={16} />}
					variant='light'
				>
					Add Sign
				</Button>
			</Group>

			{signs.map((sign) => (
				<Card key={sign.id} withBorder>
					<Stack gap='md'>
						<Group justify='space-between'>
							<Text fw={500}>{sign.description}</Text>
							<Text size='sm' c='dimmed'>
								{new Date(sign.date).toLocaleDateString()}
							</Text>
						</Group>
						{sign.notes && <Text size='sm'>{sign.notes}</Text>}
					</Stack>
				</Card>
			))}

			<Modal
				opened={isModalOpen}
				onClose={() => setIsModalOpen(false)}
				title='Add Pre-Foaling Sign'
			>
				<Stack gap='md'>
					<TextInput
						label='Description'
						placeholder='Enter sign description'
						required
						value={newSign.description || ''}
						onChange={(e) =>
							setNewSign({
								...newSign,
								description: e.target.value,
							})
						}
					/>

					<Textarea
						label='Notes'
						placeholder='Additional notes (optional)'
						value={newSign.notes || ''}
						onChange={(e) =>
							setNewSign({ ...newSign, notes: e.target.value })
						}
					/>

					<Button
						onClick={handleSubmit}
						loading={addSignMutation.isPending}
						mt='md'
					>
						Add Sign
					</Button>
				</Stack>
			</Modal>
		</Stack>
	);
}

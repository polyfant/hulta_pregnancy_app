import {
	ActionIcon,
	Alert,
	Badge,
	Button,
	Card,
	Group,
	LoadingOverlay,
	SimpleGrid,
	Stack,
	Text,
	TextInput,
	Title,
} from '@mantine/core';
import {
	Horse,
	MagnifyingGlass,
	Plus,
	User,
	Warning,
} from '@phosphor-icons/react';
import { useQuery } from '@tanstack/react-query';
import { useState } from 'react';
import { Link } from 'react-router-dom';
import { useApiClient } from '../api/client';

interface Horse {
	id: string;
	name: string;
	breed?: string;
	color?: string;
	gender: 'male' | 'female';
	birthDate?: string;
	isPregnant?: boolean;
}

const HorseList = () => {
	const [searchQuery, setSearchQuery] = useState('');
	const apiClient = useApiClient();

	const { data, isLoading, error } = useQuery<Horse[]>({
		queryKey: ['horses'],
		queryFn: async () => {
			console.log('Fetching horses...');
			try {
				const data = await apiClient.get<Horse[]>('/horses');
				console.log('Horses data:', data);
				return data;
			} catch (error) {
				console.error('Error fetching horses:', error);
				throw error;
			}
		},
		retry: 1,
		staleTime: 30000,
		refetchOnWindowFocus: false,
	});

	// Ensure data is never null
	const horses = data || [];

	// Filter horses based on search query
	const filteredHorses = horses.filter(
		(horse) =>
			horse.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
			(horse.breed &&
				horse.breed.toLowerCase().includes(searchQuery.toLowerCase()))
	);

	if (error) {
		return (
			<Alert icon={<Warning size='1rem' />} title='Error' color='red'>
				Failed to load horses. Please try again later.
			</Alert>
		);
	}

	return (
		<Stack>
			<Group justify='space-between' align='center'>
				<Title order={2}>Horses</Title>
				<Button
					component={Link}
					to='/add-horse'
					variant='filled'
					styles={(theme) => ({
						root: {
							color: theme.colors.green[4],
							backgroundColor: theme.colors.dark[7],
							'&:hover': {
								backgroundColor: theme.colors.dark[8],
							},
						},
					})}
					leftSection={<Plus size='1rem' />}
				>
					Add Horse
				</Button>
			</Group>

			<TextInput
				placeholder='Search horses...'
				leftSection={<MagnifyingGlass size='1rem' />}
				value={searchQuery}
				onChange={(e) => setSearchQuery(e.target.value)}
			/>

			<div style={{ position: 'relative' }}>
				<LoadingOverlay visible={isLoading} />
				<SimpleGrid cols={{ base: 1, sm: 2, md: 3 }} spacing='md'>
					{filteredHorses.map((horse) => (
						<Card
							key={horse.id}
							shadow='sm'
							padding='lg'
							radius='md'
							withBorder
						>
							<Group justify='space-between' mb='xs'>
								<Text fw={500}>{horse.name}</Text>
								<ActionIcon
									variant='light'
									color={
										horse.gender === 'male'
											? 'blue'
											: 'pink'
									}
									title={horse.gender}
								>
									{horse.gender === 'male' ? (
										<User color='blue' size='1.2rem' />
									) : (
										<User color='pink' size='1.2rem' />
									)}
								</ActionIcon>
							</Group>

							<Group gap='xs'>
								{horse.breed && (
									<Badge color='blue' variant='light'>
										{horse.breed}
									</Badge>
								)}
								{horse.isPregnant && (
									<Badge color='grape' variant='light'>
										Pregnant
									</Badge>
								)}
							</Group>

							<Button
								component={Link}
								to={`/horses/${horse.id}`}
								variant='light'
								color='blue'
								fullWidth
								mt='md'
								radius='md'
								leftSection={<Horse size='1rem' />}
							>
								View Details
							</Button>
						</Card>
					))}
				</SimpleGrid>
			</div>
		</Stack>
	);
};

export default HorseList;

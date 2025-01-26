import {
	ActionIcon,
	Alert,
	Badge,
	Button,
	Card,
	Group,
	LoadingOverlay,
	Progress,
	SimpleGrid,
	Stack,
	Text,
	TextInput,
	Title
} from '@mantine/core';
import {
	GenderFemale,
	GenderMale,
	Horse,
	MagnifyingGlass,
	Plus,
	Warning
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
	gender: 'MARE' | 'STALLION' | 'GELDING';
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

	// Fetch pregnancy status for all pregnant horses at once
	const pregnancyQueries = useQuery({
		queryKey: ['pregnancy'],
		queryFn: async () => {
			const pregnantHorses = filteredHorses.filter(horse => horse.isPregnant);
			const promises = pregnantHorses.map(horse => 
				fetch(`/api/horses/${horse.id}/pregnancy`)
					.then(response => {
						if (!response.ok) {
							throw new Error('Failed to fetch pregnancy status');
						}
						return response.json();
					})
			);
			return Promise.all(promises);
		},
		enabled: filteredHorses.some(horse => horse.isPregnant),
	});

	// Create a map of pregnancy statuses for easy lookup
	const pregnancyStatusMap = {};
	if (pregnancyQueries.data) {
		const pregnantHorses = filteredHorses.filter(horse => horse.isPregnant);
		pregnantHorses.forEach((horse, index) => {
			pregnancyStatusMap[horse.id] = pregnancyQueries.data[index];
		});
	}

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
				<Title order={2} c="white">Horses</Title>
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
				styles={(theme) => ({
					input: {
						backgroundColor: theme.colors.dark[7],
						color: theme.white,
						'&::placeholder': {
							color: theme.colors.dark[2],
						},
					},
				})}
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
							bg="dark.7"
						>
							<Group justify='space-between' mb='xs'>
								<Text fw={500} c="white">{horse.name}</Text>
								<ActionIcon
									variant='light'
									color={horse.gender === 'STALLION' || horse.gender === 'GELDING' ? 'blue' : 'pink'}
									title={horse.gender}
								>
									{horse.gender === 'STALLION' || horse.gender === 'GELDING' ? (
										<GenderMale color='blue' size='1.2rem' />
									) : (
										<GenderFemale color='pink' size='1.2rem' />
									)}
								</ActionIcon>
							</Group>

							<Group gap='xs' mb="xs">
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

							{horse.isPregnant && pregnancyStatusMap[horse.id] && (
								<Progress 
									value={pregnancyStatusMap[horse.id].progress} 
									color="grape" 
									size="sm" 
									mb="md"
								/>
							)}

							<Button
								component={Link}
								to={`/horses/${horse.id}`}
								variant='light'
								color='blue'
								fullWidth
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

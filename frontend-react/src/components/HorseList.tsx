import {
	ActionIcon,
	Badge,
	Button,
	Card,
	Group,
	LoadingOverlay,
	Progress,
	SimpleGrid,
	Skeleton,
	Stack,
	Text,
	TextInput,
	Title,
} from '@mantine/core';
import {
	GenderFemale,
	GenderMale,
	Horse,
	MagnifyingGlass,
	Plus,
} from '@phosphor-icons/react';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import { useMemo, useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useApiClient } from '../api/client';
import { PregnancyStage } from '../types/pregnancy';
import { EmptyState } from './states/EmptyState';
import { NetworkError } from './states/NetworkError';
import { vi } from 'vitest';

interface Horse {
	id: string;
	name: string;
	breed?: string;
	color?: string;
	gender: 'MARE' | 'STALLION' | 'GELDING';
	birthDate?: string;
	isPregnant?: boolean;
}

interface PregnancyStatusMap {
	[key: string]: {
		currentStage: PregnancyStage;
		daysRemaining: number;
		progress: number;
	};
}

interface PregnancyResponse {
	currentStage: PregnancyStage;
	daysRemaining: number;
}

const getStageColor = (stage: string) => {
	const colors = {
		EARLY: 'blue',
		MIDDLE: 'cyan',
		LATE: 'teal',
		NEARTERM: 'indigo',
		FOALING: 'grape',
	};
	return colors[stage as keyof typeof colors] || 'gray';
};

const HorseCardSkeleton = () => (
	<Card shadow='sm' padding='lg' radius='md' withBorder bg='dark.7'>
		<Skeleton height={20} width='60%' mb='xs' />
		<Group gap='xs' mb='xs'>
			<Skeleton height={20} width={60} />
			<Skeleton height={20} width={80} />
		</Group>
		<Skeleton height={8} width='100%' mb='xl' />
		<Skeleton height={36} width='100%' />
	</Card>
);

const mockNavigate = vi.fn();
vi.mock('react-router-dom', async () => ({
	...(await vi.importActual('react-router-dom')),
	useNavigate: () => mockNavigate,
}));

export function HorseList() {
	const [searchQuery, setSearchQuery] = useState('');
	const apiClient = useApiClient();
	const navigate = useNavigate();
	const queryClient = useQueryClient();

	const { data, isLoading, error, refetch } = useQuery<Horse[]>({
		queryKey: ['horses'],
		queryFn: async () => {
			console.log('Fetching horses...');
			try {
				const data = await apiClient.get<Horse[]>('/api/horses');
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

	const horses = data || [];

	const filteredHorses = horses.filter(
		(horse) =>
			horse.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
			(horse.breed &&
				horse.breed.toLowerCase().includes(searchQuery.toLowerCase()))
	);

	const pregnancyQueries = useQuery<PregnancyResponse[]>({
		queryKey: ['pregnancy'],
		queryFn: async () => {
			const pregnantHorses = filteredHorses.filter(
				(horse) => horse.isPregnant
			);
			const promises = pregnantHorses.map((horse) =>
				fetch(`/api/horses/${horse.id}/pregnancy`).then((response) => {
					if (!response.ok) {
						throw new Error('Failed to fetch pregnancy status');
					}
					return response.json();
				})
			);
			return Promise.all(promises);
		},
		enabled: filteredHorses.some((horse) => horse.isPregnant),
	});

	const pregnancyStatusMap = useMemo(() => {
		const statusMap: PregnancyStatusMap = {};
		if (pregnancyQueries.data) {
			const pregnantHorses = filteredHorses.filter(
				(horse) => horse.isPregnant
			);
			pregnantHorses.forEach((horse, index) => {
				const status = pregnancyQueries.data[index];
				if (status) {
					statusMap[horse.id] = {
						currentStage: status.currentStage,
						daysRemaining: status.daysRemaining,
						progress: Math.min(
							Math.max(
								((340 - status.daysRemaining) / 340) * 100,
								0
							),
							100
						),
					};
				}
			});
		}
		return statusMap;
	}, [pregnancyQueries.data, filteredHorses]);

	const getPregnancyStatus = (horseId: string) => {
		if (!hasPregnancyStatus(horseId)) return null;
		return pregnancyStatusMap[horseId];
	};

	const _getStatusBadge = (horse: Horse) => {
		if (!horse.isPregnant) return null;
		const status = getPregnancyStatus(horse.id);
		if (!status) return null;
		return (
			<Badge color={getStageColor(status.currentStage)}>
				{status.daysRemaining} days to foaling
			</Badge>
		);
	};

	const hasPregnancyStatus = (horseId: string): boolean => {
		return Boolean(
			pregnancyStatusMap[horseId] &&
				pregnancyStatusMap[horseId].currentStage &&
				pregnancyStatusMap[horseId].daysRemaining
		);
	};

	if (error) {
		return (
			<NetworkError
				message='Failed to load horses. Please try again.'
				onRetry={() => refetch()}
			/>
		);
	}

	if (isLoading) {
		return (
			<SimpleGrid cols={{ base: 1, sm: 2, md: 3 }} spacing='md'>
				{[...Array(6)].map((_, i) => (
					<HorseCardSkeleton key={i} />
				))}
			</SimpleGrid>
		);
	}

	if (!filteredHorses.length) {
		return (
			<EmptyState
				title='No Horses Found'
				message={
					searchQuery
						? 'No horses match your search criteria.'
						: 'Start by adding your first horse!'
				}
				actionLabel='Add Horse'
				onAction={() => navigate('/add-horse')}
			/>
		);
	}

	return (
		<Stack>
			<Group justify='space-between' align='center'>
				<Title order={2} c='white'>
					Horses
				</Title>
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
							bg='dark.7'
						>
							<Group justify='space-between' mb='xs'>
								<Text fw={500} c='white'>
									{horse.name}
								</Text>
								<ActionIcon
									variant='light'
									color={
										horse.gender === 'STALLION' ||
										horse.gender === 'GELDING'
											? 'blue'
											: 'pink'
									}
									title={horse.gender}
								>
									{horse.gender === 'STALLION' ||
									horse.gender === 'GELDING' ? (
										<GenderMale
											color='blue'
											size='1.2rem'
										/>
									) : (
										<GenderFemale
											color='pink'
											size='1.2rem'
										/>
									)}
								</ActionIcon>
							</Group>

							<Group gap='xs' mb='xs'>
								{horse.breed && (
									<Badge color='blue' variant='light'>
										{horse.breed}
									</Badge>
								)}
								{horse.isPregnant && _getStatusBadge(horse)}
							</Group>

							{horse.isPregnant &&
								hasPregnancyStatus(horse.id) && (
									<Progress
										value={
											pregnancyStatusMap[horse.id]
												?.progress ?? 0
										}
										color='grape'
										size='sm'
										mb='md'
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
}

export default HorseList;

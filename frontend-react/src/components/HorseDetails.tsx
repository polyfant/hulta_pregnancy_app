import {
	ActionIcon,
	Box,
	Button,
	Card,
	Grid,
	Group,
	Image,
	LoadingOverlay,
	Paper,
	Stack,
	Tabs,
	Text,
	Title,
} from '@mantine/core';
import {
	Activity,
	Calendar,
	GenderFemale,
	GenderMale,
	Heart,
	Horse,
	Pencil,
	Scales,
	Syringe,
	Tag,
	Trash
} from '@phosphor-icons/react';
import { Link, useNavigate, useParams } from 'react-router-dom';

import { notifications } from '@mantine/notifications';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { lazy, Suspense, useCallback } from 'react';

// Lazy load FamilyTree component
const FamilyTree = lazy(() => import('./FamilyTree/FamilyTree'));
const PregnancyStatus = lazy(
	() => import('./PregnancyTracking/PregnancyStatus')
);

interface Horse {
	id: number;
	name: string;
	breed: string;
	gender: 'MARE' | 'STALLION' | 'GELDING';
	dateOfBirth: string;
	weight?: number;
	age?: string;
	conceptionDate?: string;
	motherId?: number;
	fatherId?: number;
	externalMother?: string;
	externalFather?: string;
	created_at?: string;
	updated_at?: string;
	imageUrl?: string;
}

const HorseDetails = () => {
	const { id } = useParams<{ id: string }>();
	const navigate = useNavigate();
	const queryClient = useQueryClient();

	const {
		data: horse,
		isLoading,
		error,
	} = useQuery<Horse>({
		queryKey: ['horse', id],
		queryFn: async () => {
			const response = await fetch(`/api/horses/${id}`);
			if (!response.ok) throw new Error('Failed to fetch horse details');
			return response.json();
		},
		staleTime: 30000,
		gcTime: 5 * 60 * 1000,
		refetchOnWindowFocus: false,
	});

	const deleteMutation = useMutation({
		mutationFn: async () => {
			const response = await fetch(`/api/horses/${id}`, {
				method: 'DELETE',
			});
			if (!response.ok) throw new Error('Failed to delete horse');
		},
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ['horses'] });
			notifications.show({
				title: 'Success',
				message: 'Horse deleted successfully',
				color: 'green',
			});
			navigate('/');
		},
		onError: (error: Error) => {
			notifications.show({
				title: 'Error',
				message: error.message,
				color: 'red',
			});
		},
	});

	const handleDelete = useCallback(() => {
		if (
			window.confirm(
				'Are you sure you want to delete this horse? This action cannot be undone.'
			)
		) {
			deleteMutation.mutate();
		}
	}, [deleteMutation]);

	if (isLoading) {
		return (
			<Paper p='xl' pos='relative'>
				<LoadingOverlay visible />
			</Paper>
		);
	}

	if (error || !horse) {
		return (
			<Paper p='xl'>
				<Text c='red'>Error loading horse details</Text>
			</Paper>
		);
	}

	return (
		<Stack gap='lg'>
			<Card withBorder bg="dark.7">
				<Group justify='space-between' mb='md'>
					<Group gap='sm'>
						{horse.gender === 'STALLION' || horse.gender === 'GELDING' ? (
							<GenderMale size={30} color="var(--mantine-color-blue-6)" />
						) : (
							<GenderFemale size={30} color="var(--mantine-color-pink-6)" />
						)}
						<Title order={2} c="white">
							{horse.name}
						</Title>
					</Group>
					<Group gap='sm'>
						<Button
							component={Link}
							to={`/horses/${id}/edit`}
							variant='filled'
							leftSection={<Pencil size='1rem' />}
						>
							Edit Horse
						</Button>
						<ActionIcon
							color='red'
							variant='light'
							onClick={handleDelete}
							loading={deleteMutation.isPending}
						>
							<Trash size='1rem' />
						</ActionIcon>
					</Group>
				</Group>

				<Tabs defaultValue='details' styles={(theme) => ({
					tab: {
						color: theme.white,
						'&:hover': {
							backgroundColor: theme.colors.dark[5],
							color: theme.white,
						},
						'&[data-active="true"]': {
							backgroundColor: theme.colors.dark[6],
							color: theme.white,
						},
					},
					panel: {
						color: theme.white,
					}
				})}>
					<Tabs.List>
						<Tabs.Tab
							value='details'
							leftSection={<Calendar size='1rem' />}
						>
							Details
						</Tabs.Tab>
						<Tabs.Tab
							value='health'
							leftSection={<Heart size='1rem' />}
						>
							Health
						</Tabs.Tab>
						{horse.gender === 'MARE' && horse.conceptionDate && (
							<Tabs.Tab
								value='pregnancy'
								leftSection={<Activity size='1rem' />}
							>
								Pregnancy
							</Tabs.Tab>
						)}
						<Tabs.Tab
							value='family'
							leftSection={<Horse size='1rem' />}
						>
							Family Tree
						</Tabs.Tab>
					</Tabs.List>

					<Box mt='md'>
						<Tabs.Panel value='details'>
							<Paper p='md' withBorder bg="dark.8">
								<Grid>
									<Grid.Col span={{ base: 12, md: 6 }}>
										<Stack>
											<Group>
												<Calendar size='1rem' />
												<Text>
													Born: {new Date(horse.dateOfBirth).toLocaleDateString()}
												</Text>
											</Group>
											<Group>
												{horse.gender === 'STALLION' || horse.gender === 'GELDING' ? (
													<GenderMale size='1rem' color="var(--mantine-color-blue-6)" />
												) : (
													<GenderFemale size='1rem' color="var(--mantine-color-pink-6)" />
												)}
												<Text>
													Gender: {horse.gender}
												</Text>
											</Group>
											{horse.breed && (
												<Group>
													<Tag size='1rem' />
													<Text>
														Breed: {horse.breed}
													</Text>
												</Group>
											)}
											{horse.weight && (
												<Group>
													<Scales size='1rem' />
													<Text>
														Weight: {horse.weight} kg
													</Text>
												</Group>
											)}
										</Stack>
									</Grid.Col>
									<Grid.Col span={{ base: 12, md: 6 }}>
										{horse.imageUrl && (
											<Image
												src={horse.imageUrl}
												alt={horse.name}
												radius='md'
												fit='cover'
											/>
										)}
									</Grid.Col>
								</Grid>
							</Paper>
						</Tabs.Panel>

						<Tabs.Panel value='health'>
							<Paper p='md' withBorder bg="dark.8">
								<Stack>
									<Group>
										<Heart size='1rem' />
										<Text c="white">Health Status: Healthy</Text>
									</Group>
									<Group>
										<Syringe size='1rem' />
										<Text c="white">Last Vaccination: Up to date</Text>
									</Group>
								</Stack>
							</Paper>
						</Tabs.Panel>

						{horse.gender === 'MARE' && horse.conceptionDate && (
							<Tabs.Panel value='pregnancy'>
								<Paper p='md' withBorder bg="dark.8">
									<Suspense fallback={<LoadingOverlay visible />}>
										<PregnancyStatus horseId={horse.id} />
									</Suspense>
								</Paper>
							</Tabs.Panel>
						)}

						<Tabs.Panel value='family'>
							<Paper p='md' withBorder>
								<Suspense fallback={<LoadingOverlay visible />}>
									<FamilyTree horseId={horse.id} />
								</Suspense>
							</Paper>
						</Tabs.Panel>
					</Box>
				</Tabs>
			</Card>
		</Stack>
	);
};

export default HorseDetails;

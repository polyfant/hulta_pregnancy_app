import {
	ActionIcon,
	Badge,
	Card,
	Collapse,
	Divider,
	Grid,
	Group,
	Paper,
	Progress,
	SegmentedControl,
	Select,
	Stack,
	Switch,
	Tabs,
	Text,
	TextInput,
	ThemeIcon,
	Timeline,
} from '@mantine/core';
import {
	Baby,
	Bell,
	CaretDown,
	CaretRight,
	Horse,
	MagnifyingGlass,
	NotePencil,
	Scales,
	Syringe,
} from '@phosphor-icons/react';
import { useQuery } from '@tanstack/react-query';
import { useState } from 'react';
import { Link } from 'react-router-dom';

interface PregnancyNote {
	id: number;
	date: string;
	type: 'VET_CHECK' | 'WEIGHT' | 'VACCINATION' | 'OBSERVATION';
	content: string;
}

interface PregnantHorse {
	id: number;
	name: string;
	progress: number;
	daysUntilDue: number;
	stage: string;
	isInDueWindow: boolean;
	isOverdue: boolean;
	imageUrl?: string;
	breed?: string;
	weight?: number;
	lastVetCheck?: string;
	nextVetCheck?: string;
	notes: PregnancyNote[];
	notifications: {
		dueDate: boolean;
		vetChecks: boolean;
		weightChecks: boolean;
	};
}

interface Foal {
	id: number;
	name: string;
	motherId: number;
	motherName: string;
	dateOfBirth: string;
	breed?: string;
	gender: 'MALE' | 'FEMALE';
	ageInDays: number;
}

export function PregnancyOverview() {
	const [searchQuery, setSearchQuery] = useState('');
	const [sortBy, setSortBy] = useState('dueDate');
	const [filterStage, setFilterStage] = useState('all');

	const { data: pregnantHorses } = useQuery<PregnantHorse[]>({
		queryKey: ['pregnant-horses'],
		queryFn: async () => {
			const response = await fetch('/api/horses/pregnant');
			if (!response.ok)
				throw new Error('Failed to fetch pregnant horses');
			return response.json();
		},
	});

	const { data: recentFoals } = useQuery<Foal[]>({
		queryKey: ['recent-foals'],
		queryFn: async () => {
			const response = await fetch('/api/horses/foals/recent');
			if (!response.ok) throw new Error('Failed to fetch recent foals');
			return response.json();
		},
	});

	const filteredAndSortedHorses = pregnantHorses
		?.filter((horse) => {
			const matchesSearch =
				horse.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
				horse.breed?.toLowerCase().includes(searchQuery.toLowerCase());
			const matchesStage =
				filterStage === 'all' ||
				horse.stage.toLowerCase() === filterStage;
			return matchesSearch && matchesStage;
		})
		.sort((a, b) => {
			switch (sortBy) {
				case 'dueDate':
					return a.daysUntilDue - b.daysUntilDue;
				case 'progress':
					return b.progress - a.progress;
				case 'name':
					return a.name.localeCompare(b.name);
				default:
					return 0;
			}
		});

	return (
		<Stack>
			<Tabs defaultValue='pregnant'>
				<Tabs.List>
					<Tabs.Tab value='pregnant' leftSection={<Baby size={16} />}>
						Pregnant Mares ({pregnantHorses?.length || 0})
					</Tabs.Tab>
					<Tabs.Tab value='foals' leftSection={<Horse size={16} />}>
						Recent Foals ({recentFoals?.length || 0})
					</Tabs.Tab>
				</Tabs.List>

				<Tabs.Panel value='pregnant'>
					<Stack mt='md'>
						{/* Filters and Controls */}
						<Group>
							<TextInput
								placeholder='Search mares...'
								value={searchQuery}
								onChange={(e) => setSearchQuery(e.target.value)}
								leftSection={<MagnifyingGlass size={16} />}
								style={{ flex: 1 }}
							/>
							<Select
								value={sortBy}
								onChange={(value) =>
									setSortBy(value || 'dueDate')
								}
								data={[
									{
										value: 'dueDate',
										label: 'Sort by Due Date',
									},
									{
										value: 'progress',
										label: 'Sort by Progress',
									},
									{ value: 'name', label: 'Sort by Name' },
								]}
								style={{ width: 200 }}
							/>
							<SegmentedControl
								value={filterStage}
								onChange={setFilterStage}
								data={[
									{ value: 'all', label: 'All' },
									{ value: 'early', label: 'Early' },
									{ value: 'mid', label: 'Mid' },
									{ value: 'late', label: 'Late' },
									{ value: 'overdue', label: 'Overdue' },
								]}
							/>
						</Group>

						{/* Pregnant Mares Grid */}
						<Grid>
							{filteredAndSortedHorses?.map((horse) => (
								<Grid.Col
									key={horse.id}
									span={{ base: 12, sm: 6, md: 4 }}
								>
									<PregnancyCard horse={horse} />
								</Grid.Col>
							))}
						</Grid>
					</Stack>
				</Tabs.Panel>

				<Tabs.Panel value='foals'>
					<Grid mt='md'>
						{recentFoals?.map((foal) => (
							<Grid.Col
								key={foal.id}
								span={{ base: 12, sm: 6, md: 4 }}
							>
								<Card withBorder shadow='sm'>
									<Group>
										<Horse size={24} />
										<div>
											<Text fw={500}>{foal.name}</Text>
											<Text size='sm' c='dimmed'>
												{foal.ageInDays} days old
											</Text>
										</div>
									</Group>
									<Stack mt='md' spacing='xs'>
										<Text size='sm'>
											Born:{' '}
											{new Date(
												foal.dateOfBirth
											).toLocaleDateString()}
										</Text>
										<Text size='sm'>
											Dam: {foal.motherName}
										</Text>
										<Badge>{foal.gender}</Badge>
									</Stack>
								</Card>
							</Grid.Col>
						))}
					</Grid>
				</Tabs.Panel>
			</Tabs>
		</Stack>
	);
}

export function PregnancyCard({ horse }: { horse: PregnantHorse }) {
	const [expanded, setExpanded] = useState(false);

	return (
		<Card withBorder shadow='sm'>
			<Card.Section
				p='md'
				bg={horse.isOverdue ? 'red.1' : 'blue.1'}
				onClick={() => setExpanded(!expanded)}
				style={{ cursor: 'pointer' }}
			>
				<Group justify='space-between'>
					<Group>
						<Baby size={24} />
						<div>
							<Text fw={500}>{horse.name}</Text>
							<Text size='sm' c='dimmed'>
								{horse.isOverdue
									? 'Overdue'
									: horse.isInDueWindow
									? 'Due Soon'
									: horse.stage}
							</Text>
						</div>
					</Group>
					<ActionIcon
						component={Link}
						to={`/horses/${horse.id}/pregnancy`}
						variant='subtle'
					>
						<CaretRight size={16} />
					</ActionIcon>
					<ActionIcon variant='subtle'>
						{expanded ? (
							<CaretDown size={16} />
						) : (
							<CaretRight size={16} />
						)}
					</ActionIcon>
				</Group>
			</Card.Section>

			<Stack gap='xs' mt='md'>
				<Group justify='space-between'>
					<Text size='sm'>Progress</Text>
					<Badge
						color={horse.isOverdue ? 'red' : 'blue'}
						variant='light'
					>
						{Math.round(horse.progress)}%
					</Badge>
				</Group>
				<Progress
					value={horse.progress}
					color={horse.isOverdue ? 'red' : 'blue'}
					size='sm'
				/>
				<Text size='sm' c='dimmed' ta='right'>
					{horse.isOverdue
						? `${Math.abs(horse.daysUntilDue)} days overdue`
						: `${horse.daysUntilDue} days until due`}
				</Text>
			</Stack>

			<Collapse in={expanded}>
				<Divider my='md' />

				<Group grow mb='md'>
					<Paper p='xs' withBorder>
						<Text size='sm' fw={500}>
							Last Weight
						</Text>
						<Text>{horse.weight || 'N/A'} kg</Text>
					</Paper>
					<Paper p='xs' withBorder>
						<Text size='sm' fw={500}>
							Next Vet Check
						</Text>
						<Text>
							{horse.nextVetCheck
								? new Date(
										horse.nextVetCheck
								  ).toLocaleDateString()
								: 'N/A'}
						</Text>
					</Paper>
				</Group>

				<Timeline active={-1} bulletSize={24} lineWidth={2}>
					{horse.notes.map((note, index) => (
						<Timeline.Item
							key={note.id}
							bullet={
								<ThemeIcon
									size={24}
									radius='xl'
									color={
										note.type === 'VET_CHECK'
											? 'blue'
											: note.type === 'VACCINATION'
											? 'green'
											: note.type === 'WEIGHT'
											? 'orange'
											: 'gray'
									}
								>
									{note.type === 'VET_CHECK' && (
										<Syringe size={12} />
									)}
									{note.type === 'WEIGHT' && (
										<Scales size={12} />
									)}
									{note.type === 'OBSERVATION' && (
										<NotePencil size={12} />
									)}
								</ThemeIcon>
							}
							title={note.type.replace('_', ' ')}
						>
							<Text size='sm' mt={4}>
								{note.content}
							</Text>
							<Text size='xs' c='dimmed'>
								{new Date(note.date).toLocaleDateString()}
							</Text>
						</Timeline.Item>
					))}
				</Timeline>

				<Paper withBorder p='sm' mt='md'>
					<Text fw={500} mb='xs'>
						Notifications
					</Text>
					<Stack gap='xs'>
						<Group>
							<Switch
								checked={horse.notifications.dueDate}
								label='Due Date Alerts'
								leftSection={<Bell size={16} />}
							/>
						</Group>
						<Group>
							<Switch
								checked={horse.notifications.vetChecks}
								label='Vet Check Reminders'
								leftSection={<Syringe size={16} />}
							/>
						</Group>
						<Group>
							<Switch
								checked={horse.notifications.weightChecks}
								label='Weight Check Reminders'
								leftSection={<Scales size={16} />}
							/>
						</Group>
					</Stack>
				</Paper>
			</Collapse>
		</Card>
	);
}

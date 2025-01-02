import {
	Accordion,
	Badge,
	Card,
	Group,
	List,
	LoadingOverlay,
	Paper,
	Progress,
	Stack,
	Text,
	ThemeIcon,
	Timeline,
	Title,
} from '@mantine/core';
import {
	Baby,
	Calendar,
	Clock,
	Heart,
	Info,
	Syringe
} from '@phosphor-icons/react';

import { useQuery } from '@tanstack/react-query';

interface Horse {
	id: number;
	name: string;
	conceptionDate?: string;
}

interface PregnancyStatusProps {
	horseId: number;
}

const PREGNANCY_DURATION = 340; // Average horse pregnancy duration in days

const PregnancyStatus = ({ horseId }: PregnancyStatusProps) => {
	const { data: horse, isLoading } = useQuery<Horse>({
		queryKey: ['horse', horseId],
		queryFn: async () => {
			const response = await fetch(`/api/horses/${horseId}`);
			if (!response.ok) throw new Error('Failed to fetch horse details');
			return response.json();
		},
	});

	if (isLoading) {
		return (
			<Paper p='md' pos='relative' h={200}>
				<LoadingOverlay visible />
			</Paper>
		);
	}

	if (!horse?.conceptionDate) {
		return (
			<Paper p='md'>
				<Text>
					No pregnancy information available for {horse?.name}.
				</Text>
			</Paper>
		);
	}

	const startDate = new Date(horse.conceptionDate);
	const today = new Date();
	const daysPregnant = Math.floor(
		(today.getTime() - startDate.getTime()) / (1000 * 60 * 60 * 24)
	);
	const dueDate = new Date(
		startDate.getTime() + PREGNANCY_DURATION * 24 * 60 * 60 * 1000
	);
	const progressPercentage = Math.min(
		100,
		Math.round((daysPregnant / PREGNANCY_DURATION) * 100)
	);

	const trimester = Math.floor(daysPregnant / (PREGNANCY_DURATION / 3)) + 1;

	const getPregnancyStage = () => {
		if (daysPregnant < 114) return { stage: 'Early Stage', color: 'blue' };
		if (daysPregnant < 226) return { stage: 'Mid Stage', color: 'yellow' };
		if (daysPregnant < 310) return { stage: 'Late Stage', color: 'orange' };
		return { stage: 'Due Soon', color: 'red' };
	};

	const { stage, color } = getPregnancyStage();

	const getRecommendedActions = () => {
		if (daysPregnant < 114) {
			return [
				'Schedule initial vet checkup',
				'Maintain regular exercise routine',
				'Monitor appetite and weight',
				'Consider vaccinations if needed',
			];
		}
		if (daysPregnant < 226) {
			return [
				'Schedule mid-term ultrasound',
				'Adjust feed for increased nutritional needs',
				'Continue moderate exercise',
				'Monitor for any unusual behavior',
			];
		}
		if (daysPregnant < 310) {
			return [
				'Prepare foaling area',
				'Monitor udder development',
				'Reduce exercise intensity',
				'Schedule pre-foaling checkup',
			];
		}
		return [
			'Monitor for signs of impending labor',
			'Have vet on standby',
			'Check mare frequently',
			'Ensure foaling kit is ready',
		];
	};

	return (
		<Stack gap='lg'>
			<Card withBorder bg="dark.7">
				<Group justify='space-between' mb='md'>
					<Title order={3} c="white">Pregnancy Progress</Title>
					<Badge size='lg' variant='filled' color={color}>
						{stage}
					</Badge>
				</Group>

				<Progress
					value={progressPercentage}
					size='xl'
					color={color}
					mb='sm'
				/>

				<Group gap='lg'>
					<div>
						<Text size='sm' c='dimmed'>
							Days Pregnant
						</Text>
						<Text fw={500} c="white">{daysPregnant} days</Text>
					</div>
					<div>
						<Text size='sm' c='dimmed'>
							Due Date
						</Text>
						<Text fw={500} c="white">{dueDate.toLocaleDateString()}</Text>
					</div>
					<div>
						<Text size='sm' c='dimmed'>
							Trimester
						</Text>
						<Text fw={500} c="white">{trimester}</Text>
					</div>
				</Group>
			</Card>

			<Card withBorder bg="dark.7">
				<Title order={4} mb='md' c="white">
					Recommended Actions
				</Title>
				<List 
					spacing='xs'
					styles={(theme) => ({
						item: {
							color: theme.white,
						},
					})}
				>
					{getRecommendedActions().map((action, index) => (
						<List.Item key={index}>{action}</List.Item>
						))}
				</List>
			</Card>

			<Accordion variant='contained' styles={(theme) => ({
				label: {
					color: theme.white,
				},
				content: {
					color: theme.white,
					backgroundColor: theme.colors.dark[7],
				},
				control: {
					backgroundColor: theme.colors.dark[8],
					'&:hover': {
						backgroundColor: theme.colors.dark[6],
					},
				},
				item: {
					backgroundColor: theme.colors.dark[8],
					border: `1px solid ${theme.colors.dark[6]}`,
				}
			})}>
				<Accordion.Item value='currentStage'>
					<Accordion.Control icon={<Info size='1rem' color="white" />}>
						Current Stage Information
					</Accordion.Control>
					<Accordion.Panel>
						<Timeline
							active={Math.min(3, trimester)}
							bulletSize={24}
							styles={(theme) => ({
								item: {
									color: theme.white,
								},
								itemTitle: {
									color: theme.white,
								}
							})}
						>
							<Timeline.Item
								bullet={<Calendar size={12} />}
								title='First Trimester (0-114 days)'
								styles={(theme) => ({
									title: {
										color: theme.white,
									}
								})}
							>
								<Text size='sm' c="white">
									Early development stage. Regular check-ups important.
								</Text>
								<List size='sm' mt='xs' styles={(theme) => ({
									item: {
										color: theme.white,
									}
								})}>
									<List.Item>Schedule regular vet check-ups</List.Item>
									<List.Item>Maintain normal exercise routine</List.Item>
									<List.Item>Monitor appetite and weight</List.Item>
								</List>
							</Timeline.Item>

							<Timeline.Item
								bullet={<Heart size={12} />}
								title='Second Trimester (115-225 days)'
								styles={(theme) => ({
									title: {
										color: theme.white,
									}
								})}
							>
								<Text size='sm' c="white">
									Growth and development phase.
								</Text>
								<List size='sm' mt='xs' styles={(theme) => ({
									item: {
										color: theme.white,
									}
								})}>
									<List.Item>Continue moderate exercise</List.Item>
									<List.Item>Adjust feed as needed</List.Item>
									<List.Item>Monitor for any complications</List.Item>
								</List>
							</Timeline.Item>

							<Timeline.Item
								bullet={<Baby size={12} />}
								title='Third Trimester (226-340 days)'
								styles={(theme) => ({
									title: {
										color: theme.white,
									}
								})}
							>
								<Text size='sm' c="white">Final preparation stage.</Text>
								<List size='sm' mt='xs' styles={(theme) => ({
									item: {
										color: theme.white,
									}
								})}>
									<List.Item>Prepare foaling area</List.Item>
									<List.Item>Reduce exercise intensity</List.Item>
									<List.Item>Monitor closely for signs of labor</List.Item>
								</List>
							</Timeline.Item>
						</Timeline>
					</Accordion.Panel>
				</Accordion.Item>

				<Accordion.Item value='care'>
					<Accordion.Control icon={<Heart size='1rem' color="white" />}>
						Care Tips
					</Accordion.Control>
					<Accordion.Panel>
						<List spacing='sm'>
							<List.Item
								icon={
									<ThemeIcon color='blue' size={24}>
										<Heart size={16} color="white" />
									</ThemeIcon>
								}
							>
								<Text c="white">Schedule regular veterinary check-ups</Text>
							</List.Item>
							<List.Item
								icon={
									<ThemeIcon color='green' size={24}>
										<Syringe size={16} color="white" />
									</ThemeIcon>
								}
							>
								<Text c="white">Keep vaccinations up to date</Text>
							</List.Item>
							<List.Item
								icon={
									<ThemeIcon color='yellow' size={24}>
										<Clock size={16} color="white" />
									</ThemeIcon>
								}
							>
								<Text c="white">Watch for warning signs</Text>
							</List.Item>
						</List>
					</Accordion.Panel>
				</Accordion.Item>
			</Accordion>
		</Stack>
	);
};

export default PregnancyStatus;

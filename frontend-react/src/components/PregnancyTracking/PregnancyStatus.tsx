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
	Title,
} from '@mantine/core';
import { Calendar, Heart, Info } from '@phosphor-icons/react';

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
			<Card withBorder>
				<Group justify='space-between' mb='md'>
					<Title order={3}>Pregnancy Progress</Title>
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
						<Text fw={500}>{daysPregnant} days</Text>
					</div>
					<div>
						<Text size='sm' c='dimmed'>
							Due Date
						</Text>
						<Text fw={500}>{dueDate.toLocaleDateString()}</Text>
					</div>
					<div>
						<Text size='sm' c='dimmed'>
							Trimester
						</Text>
						<Text fw={500}>{trimester}</Text>
					</div>
				</Group>
			</Card>

			<Card withBorder>
				<Title order={4} mb='md'>
					Recommended Actions
				</Title>
				<List spacing='xs'>
					{getRecommendedActions().map((action, index) => (
						<List.Item key={index}>{action}</List.Item>
					))}
				</List>
			</Card>

			<Accordion variant='contained'>
				<Accordion.Item value='currentStage'>
					<Accordion.Control icon={<Info size='1rem' />}>
						Current Stage Information
					</Accordion.Control>
					<Accordion.Panel>
						<Timeline
							active={Math.min(3, trimester)}
							bulletSize={24}
						>
							<Timeline.Item
								bullet={<Calendar size={12} />}
								title='First Trimester (0-114 days)'
							>
								<Text size='sm'>
									Early development stage. Regular check-ups
									important.
								</Text>
								<List size='sm' mt='xs'>
									<List.Item>
										Schedule regular vet check-ups
									</List.Item>
									<List.Item>
										Maintain normal exercise routine
									</List.Item>
									<List.Item>
										Monitor appetite and weight
									</List.Item>
								</List>
							</Timeline.Item>

							<Timeline.Item
								bullet={<Heart size={12} />}
								title='Second Trimester (115-225 days)'
							>
								<Text size='sm'>
									Growth and development phase.
								</Text>
								<List size='sm' mt='xs'>
									<List.Item>
										Continue moderate exercise
									</List.Item>
									<List.Item>Adjust feed as needed</List.Item>
									<List.Item>
										Monitor for any complications
									</List.Item>
								</List>
							</Timeline.Item>

							<Timeline.Item
								bullet={<FiBaby size={12} />}
								title='Third Trimester (226-340 days)'
							>
								<Text size='sm'>Final preparation stage.</Text>
								<List size='sm' mt='xs'>
									<List.Item>Prepare foaling area</List.Item>
									<List.Item>
										Reduce exercise intensity
									</List.Item>
									<List.Item>
										Monitor closely for signs of labor
									</List.Item>
								</List>
							</Timeline.Item>
						</Timeline>
					</Accordion.Panel>
				</Accordion.Item>

				<Accordion.Item value='care'>
					<Accordion.Control icon={<FiHeart size='1rem' />}>
						Care Tips
					</Accordion.Control>
					<Accordion.Panel>
						<List spacing='sm'>
							<List.Item
								icon={
									<ThemeIcon color='blue' size={24}>
										<FiHeart size={16} />
									</ThemeIcon>
								}
							>
								Schedule regular veterinary check-ups
							</List.Item>
							<List.Item
								icon={
									<ThemeIcon color='green' size={24}>
										<FiVaccine size={16} />
									</ThemeIcon>
								}
							>
								Keep vaccinations up to date as recommended by
								your vet
							</List.Item>
							<List.Item
								icon={
									<ThemeIcon color='yellow' size={24}>
										<FiAlarm size={16} />
									</ThemeIcon>
								}
							>
								Watch for warning signs: decreased appetite,
								fever, discharge
							</List.Item>
						</List>
					</Accordion.Panel>
				</Accordion.Item>
			</Accordion>
		</Stack>
	);
};

export default PregnancyStatus;

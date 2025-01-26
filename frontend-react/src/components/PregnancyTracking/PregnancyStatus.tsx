import React from 'react';
import {
	Accordion,
	Badge,
	Card,
	Group,
	List,
	Progress,
	Stack,
	Stat,
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

interface PregnancyStatusProps {
	horseId: string;
}

const PREGNANCY_STAGES = [
	{
		title: 'Early Stage',
		weeks: '1-14',
		icon: <Baby size={12} />,
		description: 'Initial development phase',
		keyPoints: ['Regular vet checks', 'Normal exercise routine'],
	},
	{
		title: 'Mid Stage',
		weeks: '15-28',
		icon: <Heart size={12} />,
		description: 'Growth and development',
		keyPoints: ['Increased nutrition', 'Moderate exercise'],
	},
	{
		title: 'Late Stage',
		weeks: '29-48',
		icon: <Clock size={12} />,
		description: 'Final preparation',
		keyPoints: ['Prepare foaling area', 'Monitor closely'],
	},
];

export function PregnancyStatus({ horseId }: PregnancyStatusProps) {
	const { data: pregnancyInfo, isLoading, error } = useQuery({
		queryKey: ['pregnancyStatus', horseId],
		queryFn: async () => {
			const response = await fetch(`/api/horses/${horseId}/pregnancy`);
			if (!response.ok) {
				throw new Error('Failed to fetch pregnancy status');
			}
			return response.json();
		}
	});

	if (isLoading) return <Text>Loading pregnancy status...</Text>;
	if (error) return <Text color="red">Error fetching pregnancy status</Text>;
	if (!pregnancyInfo) return null;

	const { 
		daysPregnant, 
		dueDate, 
		trimester, 
		stage, 
		progress, 
		isOverdue 
	} = pregnancyInfo;

	const color = isOverdue ? 'red' : 'blue';

	const getRecommendedActions = () => {
		switch (trimester) {
			case 1: return [
				'Schedule regular vet check-ups',
				'Maintain normal exercise routine',
				'Monitor appetite and weight'
			];
			case 2: return [
				'Continue moderate exercise',
				'Adjust feed as needed',
				'Monitor for any complications'
			];
			case 3: return [
				'Prepare foaling area',
				'Reduce exercise intensity',
				'Monitor closely for signs of labor'
			];
			default: return [];
		}
	};

	return (
		<Stack gap='lg'>
			<Card withBorder bg="dark.7">
				<Group justify='space-between' mb='md'>
					<Title order={3} c="white">Pregnancy Progress</Title>
					<Badge size='lg' variant='filled' color={color}>
						{isOverdue ? 'Overdue' : `${Math.round(progress)}%`}
					</Badge>
				</Group>
				<Progress
					value={progress}
					color={isOverdue ? 'red' : 'blue'}
				/>

				<Group gap='lg' mt='md'>
					<div>
						<Text size='sm' c='dimmed'>Days Pregnant</Text>
						<Text fw={500} c="white">{daysPregnant} days</Text>
					</div>
					<div>
						<Text size='sm' c='dimmed'>Due Date</Text>
						<Text fw={500} c="white">{dueDate.toLocaleDateString()}</Text>
					</div>
					<div>
						<Text size='sm' c='dimmed'>Trimester</Text>
						<Text fw={500} c="white">{trimester}</Text>
					</div>
				</Group>
			</Card>

			<Card withBorder bg="dark.7">
				<Title order={4} mb='md' c="white">Recommended Actions</Title>
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

			<Accordion 
				variant='contained' 
				styles={(theme) => ({
					label: { color: theme.white },
					content: { 
						color: theme.white, 
						backgroundColor: theme.colors.dark[7] 
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
				})}
			>
				<Accordion.Item value='currentStage'>
					<Accordion.Control icon={<Info size='1rem' color="white" />}>
						Current Stage Information
					</Accordion.Control>
					<Accordion.Panel>
						<Timeline
							active={Math.min(3, trimester)}
							bulletSize={24}
							styles={(theme) => ({
								item: { color: theme.white },
								itemTitle: { color: theme.white }
							})}
						>
							{PREGNANCY_STAGES.map((stage, index) => (
								<Timeline.Item
									key={index}
									bullet={stage.icon}
									title={stage.title}
								>
									<Text size='sm' c="white">{stage.description}</Text>
									<List size='sm' mt='xs'>
										{stage.keyPoints.map((point, i) => (
											<List.Item key={i}>{point}</List.Item>
										))}
									</List>
								</Timeline.Item>
							))}
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
}

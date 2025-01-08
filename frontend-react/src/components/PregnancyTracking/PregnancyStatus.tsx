import {
	Alert,
	Badge,
	Card,
	Group,
	List,
	Progress,
	Stack,
	Stat,
	Text,
	Timeline,
} from '@mantine/core';
import { Baby, Calendar, Clock, Heart } from '@phosphor-icons/react';
import { useQuery } from '@tanstack/react-query';

interface PregnancyStatusProps {
	horseId: number;
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
	const { data: pregnancyInfo } = useQuery({
		queryKey: ['pregnancy', horseId],
		queryFn: async () => {
			const response = await fetch(`/api/horses/${horseId}/pregnancy`);
			if (!response.ok) throw new Error('Failed to fetch pregnancy info');
			return response.json();
		},
	});

	if (!pregnancyInfo) return null;

	return (
		<Stack>
			{/* Progress Card - Enhanced */}
			<Card withBorder>
				<Group position='apart' mb='md'>
					<Text size='lg' fw={500}>
						Pregnancy Progress
					</Text>
					<Badge
						color={pregnancyInfo.isOverdue ? 'red' : 'blue'}
						variant='filled'
					>
						{pregnancyInfo.isOverdue
							? 'Overdue'
							: `${Math.round(pregnancyInfo.progress)}%`}
					</Badge>
				</Group>
				<Progress
					value={pregnancyInfo.progress}
					size='xl'
					color={pregnancyInfo.isOverdue ? 'red' : 'blue'}
				/>

				{/* Key Stats */}
				<Group mt='md' grow>
					<Stat
						label='Days Pregnant'
						value={pregnancyInfo.daysSoFar}
					/>
					<Stat
						label='Days Until Due'
						value={pregnancyInfo.daysUntilDue}
						color={pregnancyInfo.isOverdue ? 'red' : undefined}
					/>
					<Stat
						label='Current Stage'
						value={pregnancyInfo.currentStage}
					/>
				</Group>
			</Card>

			{/* Timeline - Enhanced */}
			<Timeline active={pregnancyInfo.currentStageIndex} bulletSize={24}>
				{PREGNANCY_STAGES.map((stage, index) => (
					<Timeline.Item
						key={index}
						bullet={stage.icon}
						title={stage.title}
					>
						<Text size='sm' c='dimmed'>
							Weeks {stage.weeks}
						</Text>
						<Text size='sm' mt='xs'>
							{stage.description}
						</Text>
						<List size='sm' mt='xs'>
							{stage.keyPoints.map((point, i) => (
								<List.Item key={i}>{point}</List.Item>
							))}
						</List>
					</Timeline.Item>
				))}
			</Timeline>

			{/* Due Date Card - Enhanced */}
			<Card withBorder>
				<Stack>
					<Group>
						<Calendar size={20} />
						<div>
							<Text fw={500}>Expected Due Date</Text>
							<Text size='sm'>
								{new Date(
									pregnancyInfo.expectedDueDate
								).toLocaleDateString()}
							</Text>
						</div>
					</Group>

					{pregnancyInfo.isInDueWindow && (
						<Alert color='orange' icon={<Clock size={16} />}>
							Mare is in due window (Days{' '}
							{pregnancyInfo.minGestationDays}-
							{pregnancyInfo.maxGestationDays}) - Monitor closely
							for signs of foaling
						</Alert>
					)}
				</Stack>
			</Card>
		</Stack>
	);
}

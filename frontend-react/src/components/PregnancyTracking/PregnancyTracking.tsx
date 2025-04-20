import {
	Alert,
	Button,
	Card,
	Grid,
	Group,
	LoadingOverlay,
	Progress,
	Stack,
	Text,
	Title,
} from '@mantine/core';
import { notifications } from '@mantine/notifications';
import { Plus, Warning, X } from '@phosphor-icons/react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { useState } from 'react';
import { useParams } from 'react-router-dom';
import { PregnancyStage, PregnancyStatus } from '../../types';
import { Horse } from '../../types/horse';
import { formatDate } from '../../utils/dateUtils';
import { CriticalAlerts } from './CriticalAlerts';
import { EndPregnancyDialog } from './EndPregnancyDialog';
import { GrowthCharts } from './GrowthCharts';
import { PreFoalingSigns } from './PreFoalingSigns';
import { PregnancyChecklist } from './PregnancyChecklist';
import { PregnancyEvents } from './PregnancyEvents';
import { PregnancyGuidelines } from './PregnancyGuidelines';
import { PregnancyTimeline } from './PregnancyTimeline';
import { QuickMeasurement } from './QuickMeasurement';
import { StageChecklist } from './StageChecklist';
import { StageVisualization } from './StageVisualization';
import { StartPregnancyDialog } from './StartPregnancyDialog';
import { SyncDashboard } from './SyncDashboard';

const STAGES = {
	EARLY: { label: 'Early Stage', progress: 25, days: 110 },
	MIDDLE: { label: 'Middle Stage', progress: 50, days: 240 },
	LATE: { label: 'Late Stage', progress: 75, days: 320 },
	NEARTERM: { label: 'Near Term', progress: 90, days: 335 },
	FOALING: { label: 'Foaling', progress: 100, days: 340 },
} as const;

const STAGE_TIPS = {
	EARLY: [
		'Schedule initial vet check',
		'Maintain regular exercise routine',
		'Monitor for early pregnancy complications',
		'Ensure proper nutrition',
	],
	MIDDLE: [
		'Adjust feed for growing foal',
		'Continue moderate exercise',
		'Schedule vaccination updates',
		'Monitor weight gain',
	],
	LATE: [
		'Prepare foaling area',
		'Watch for udder development',
		'Reduce exercise intensity',
		'Begin monitoring for signs of discomfort',
	],
	NEARTERM: [
		'Monitor temperature twice daily',
		'Watch for waxing teats',
		'Have foaling kit ready',
		'Ensure 24/7 monitoring capability',
	],
	FOALING: [
		'Check mare frequently',
		'Be ready for foaling',
		'Have vet contact handy',
		'Monitor for signs of labor',
	],
};

const getStageColor = (stage: PregnancyStage) => {
	switch (stage) {
		case 'EARLY':
			return 'blue';
		case 'MIDDLE':
			return 'cyan';
		case 'LATE':
			return 'teal';
		case 'NEARTERM':
			return 'indigo';
		case 'FOALING':
			return 'grape';
		default:
			return 'gray';
	}
};

export default function PregnancyTracking() {
	const { id } = useParams<{ id: string }>();
	const [startDialogOpened, setStartDialogOpen] = useState(false);
	const [endDialogOpened, setEndDialogOpen] = useState(false);
	const queryClient = useQueryClient();

	const startPregnancyMutation = useMutation({
		mutationFn: async (date: string) => {
			const response = await fetch(`/api/horses/${id}/pregnancy/start`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ conceptionDate: date }),
			});
			if (!response.ok)
				throw new Error('Failed to start pregnancy tracking');
		},
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ['horse', id] });
			queryClient.invalidateQueries({
				queryKey: ['pregnancyStatus', id],
			});
			notifications.show({
				title: 'Success',
				message: 'Started pregnancy tracking',
				color: 'green',
			});
			setStartDialogOpen(false);
		},
	});

	const endPregnancyMutation = useMutation({
		mutationFn: async (data: { outcome: string; foalingDate: string }) => {
			const response = await fetch(`/api/horses/${id}/pregnancy/end`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify(data),
			});
			if (!response.ok)
				throw new Error('Failed to end pregnancy tracking');
		},
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ['horse', id] });
			queryClient.invalidateQueries({
				queryKey: ['pregnancyStatus', id],
			});
			notifications.show({
				title: 'Success',
				message: 'Ended pregnancy tracking',
				color: 'green',
			});
			setEndDialogOpen(false);
		},
	});

	const { data: horse, isLoading: horseLoading } = useQuery({
		queryKey: ['horse', id],
		queryFn: async () => {
			const response = await fetch(`/api/horses/${id}`);
			if (!response.ok) throw new Error('Failed to fetch horse');
			return response.json() as Promise<Horse>;
		},
		enabled: !!id,
	});

	const { data: pregnancyStatus, isLoading: statusLoading } = useQuery({
		queryKey: ['pregnancyStatus', id],
		queryFn: async () => {
			const response = await fetch(`/api/horses/${id}/pregnancy`);
			if (!response.ok)
				throw new Error('Failed to fetch pregnancy status');
			return response.json() as Promise<PregnancyStatus>;
		},
		enabled: !!id && !!horse?.isPregnant,
	});

	if (horseLoading || statusLoading) {
		return <LoadingOverlay visible />;
	}

	if (!horse) {
		return (
			<Alert icon={<Warning size='1rem' />} title='Error' color='red'>
				Horse not found
			</Alert>
		);
	}

	return (
		<Stack gap='lg'>
			{horse && (
				<StageVisualization
					horseId={parseInt(horse.id.toString(), 10)}
				/>
			)}
			<Card withBorder>
				<Stack>
					<Group justify='space-between'>
						<Title order={2}>Pregnancy Tracking</Title>
						{horse.isPregnant ? (
							<Button
								color='red'
								leftSection={<X size={16} />}
								onClick={() => setEndDialogOpen(true)}
							>
								End Pregnancy
							</Button>
						) : (
							<Button
								color='blue'
								leftSection={<Plus size={16} />}
								onClick={() => setStartDialogOpen(true)}
							>
								Start Pregnancy
							</Button>
						)}
					</Group>

					{pregnancyStatus && (
						<>
							<Progress
								value={pregnancyStatus.progress}
								size='xl'
								radius='xl'
								color={getStageColor(
									pregnancyStatus.currentStage
								)}
							/>
							<Grid>
								<Grid.Col span={6}>
									<Stack gap='xs'>
										<Text size='sm' c='dimmed'>
											Stage
										</Text>
										<Text fw={500}>
											{
												STAGES[
													pregnancyStatus.currentStage
												].label
											}
										</Text>
									</Stack>
								</Grid.Col>
								<Grid.Col span={6}>
									<Stack gap='xs'>
										<Text size='sm' c='dimmed'>
											Days in Pregnancy
										</Text>
										<Text>
											{pregnancyStatus.currentDay} days
										</Text>
									</Stack>
								</Grid.Col>
								<Grid.Col span={6}>
									<Stack gap='xs'>
										<Text size='sm' c='dimmed'>
											Days Remaining
										</Text>
										<Text>
											{pregnancyStatus.totalDays -
												pregnancyStatus.currentDay}{' '}
											days
										</Text>
									</Stack>
								</Grid.Col>
								<Grid.Col span={6}>
									<Stack gap='xs'>
										<Text size='sm' c='dimmed'>
											Conception Date
										</Text>
										<Text>
											{formatDate(
												pregnancyStatus.conceptionDate
											)}
										</Text>
									</Stack>
								</Grid.Col>
								<Grid.Col span={6}>
									<Stack gap='xs'>
										<Text size='sm' c='dimmed'>
											Expected Due Date
										</Text>
										<Text>
											{formatDate(
												pregnancyStatus.expectedDueDate
											)}
										</Text>
									</Stack>
								</Grid.Col>
								<Grid.Col span={12}>
									<Stack gap='xs'>
										<Text size='sm' c='dimmed'>
											Current Stage Tips
										</Text>
										{STAGE_TIPS[
											pregnancyStatus.currentStage
										].map((tip, index) => (
											<Text key={index} size='sm'>
												• {tip}
											</Text>
										))}
									</Stack>
								</Grid.Col>
							</Grid>

							<CriticalAlerts
								horseId={id!}
								currentStage={pregnancyStatus.currentStage}
							/>

							<StageChecklist
								horseId={id!}
								currentStage={pregnancyStatus.currentStage}
							/>
						</>
					)}
				</Stack>
			</Card>

			{horse && (
				<>
					<PregnancyEvents horseId={horse.id.toString()} />
					<QuickMeasurement
						foalId={parseInt(horse.id.toString(), 10)}
					/>
					<GrowthCharts foalId={parseInt(horse.id.toString(), 10)} />
					<PregnancyTimeline
						horseId={parseInt(horse.id.toString(), 10)}
					/>
					<PregnancyGuidelines />
					<PreFoalingSigns horseId={horse.id.toString()} />
					<SyncDashboard />
					<PregnancyChecklist />
				</>
			)}

			<StartPregnancyDialog
				opened={startDialogOpened}
				onClose={() => setStartDialogOpen(false)}
				onSubmit={(date) => startPregnancyMutation.mutate(date)}
				isLoading={startPregnancyMutation.isPending}
			/>

			<EndPregnancyDialog
				opened={endDialogOpened}
				onClose={() => setEndDialogOpen(false)}
				onSubmit={(data) => endPregnancyMutation.mutate(data)}
				isLoading={endPregnancyMutation.isPending}
			/>
		</Stack>
	);
}

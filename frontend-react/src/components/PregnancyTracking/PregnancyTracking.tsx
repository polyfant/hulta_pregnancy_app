import {
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
import { Plus, X } from '@phosphor-icons/react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { useState } from 'react';
import { useParams } from 'react-router-dom';
import { Horse } from '../../types/horse';
import { PregnancyStatus } from '../../types/pregnancy';
import { formatDate } from '../../utils/dateUtils';
import { EndPregnancyDialog } from './EndPregnancyDialog';
import { StartPregnancyDialog } from './StartPregnancyDialog';

const STAGES = {
	EARLY: { label: 'Early Stage', progress: 25 },
	MIDDLE: { label: 'Middle Stage', progress: 50 },
	LATE: { label: 'Late Stage', progress: 75 },
	NEARTERM: { label: 'Near Term', progress: 90 },
	FOALING: { label: 'Foaling', progress: 100 },
} as const;

export default function PregnancyTracking() {
	const { id } = useParams<{ id: string }>();
	const [startDialogOpened, setStartDialogOpen] = useState(false);
	const [endDialogOpened, setEndDialogOpen] = useState(false);
	const queryClient = useQueryClient();

	const { data: horse, isLoading: horseLoading } = useQuery({
		queryKey: ['horse', id],
		queryFn: async () => {
			const response = await fetch(`/api/horses/${id}`);
			if (!response.ok) throw new Error('Failed to fetch horse');
			return response.json() as Promise<Horse>;
		},
		enabled: !!id,
	});

	const { data: pregnancyStatus } = useQuery({
		queryKey: ['pregnancyStatus', id],
		queryFn: async () => {
			const response = await fetch(`/api/horses/${id}/pregnancy`);
			if (!response.ok)
				throw new Error('Failed to fetch pregnancy status');
			return response.json() as Promise<PregnancyStatus>;
		},
		enabled: !!id,
	});

	if (horseLoading) {
		return <LoadingOverlay visible />;
	}

	if (!horse) {
		return <Text>Horse not found</Text>;
	}

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

	const getStageProgress = (stage: keyof typeof STAGES) => {
		return STAGES[stage]?.progress || 0;
	};

	return (
		<Stack gap='lg'>
			<Card withBorder>
				<Stack>
					<Group justify='space-between'>
						<Title order={2}>Pregnancy Tracking</Title>
						{pregnancyStatus?.isPregnant ? (
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

					{pregnancyStatus?.isPregnant && (
						<>
							<Progress
								value={getStageProgress(
									pregnancyStatus.currentStage
								)}
								// Remove label prop since Progress component doesn't accept it
								size='xl'
								radius='xl'
							/>
							<Grid>
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
											Days in Pregnancy
										</Text>
										<Text>
											{pregnancyStatus.daysInPregnancy}{' '}
											days
										</Text>
									</Stack>
								</Grid.Col>
							</Grid>
						</>
					)}
				</Stack>
			</Card>

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

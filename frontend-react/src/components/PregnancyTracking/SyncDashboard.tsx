import {
	Badge,
	Button,
	Card,
	Group,
	Stack,
	Text,
	Timeline,
} from '@mantine/core';
import { ArrowsClockwise, Check, Warning, X } from '@phosphor-icons/react';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import { syncService } from '../../services/syncService';

export function SyncDashboard() {
	const queryClient = useQueryClient();
	const { data: syncStatus } = useQuery({
		queryKey: ['sync-status'],
		queryFn: () => syncService.getSyncStatus(),
		refetchInterval: 30000, // Check every 30 seconds
	});

	const handleManualSync = async () => {
		await syncService.syncPendingMeasurements();
		queryClient.invalidateQueries({ queryKey: ['sync-status'] });
	};

	return (
		<Card withBorder>
			<Stack>
				<Group position='apart'>
					<Text fw={500}>Sync Status</Text>
					<Badge
						color={
							syncStatus?.pendingCount > 0 ? 'yellow' : 'green'
						}
						leftSection={
							syncStatus?.pendingCount > 0 ? (
								<Warning size={12} />
							) : (
								<Check size={12} />
							)
						}
					>
						{syncStatus?.pendingCount > 0
							? `${syncStatus.pendingCount} Pending`
							: 'All Synced'}
					</Badge>
				</Group>

				<Timeline active={syncStatus?.syncErrors.length ?? 0}>
					{syncStatus?.syncErrors.map((error, index) => (
						<Timeline.Item
							key={index}
							bullet={<X size={12} />}
							title={new Date(error.timestamp).toLocaleString()}
						>
							<Text size='sm' color='red'>
								{error.error}
							</Text>
						</Timeline.Item>
					))}
				</Timeline>

				<Group position='apart'>
					<Text size='sm' c='dimmed'>
						Last synced:{' '}
						{syncStatus?.lastSuccessfulSync?.toLocaleString() ??
							'Never'}
					</Text>
					<Button
						variant='light'
						leftSection={<ArrowsClockwise size={16} />}
						onClick={handleManualSync}
						loading={syncStatus?.pendingCount > 0}
					>
						Sync Now
					</Button>
				</Group>
			</Stack>
		</Card>
	);
}

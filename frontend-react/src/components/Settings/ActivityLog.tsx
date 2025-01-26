import { Card, Stack, Text, Timeline } from '@mantine/core';
import { useQuery } from '@tanstack/react-query';

export function ActivityLog() {
	const { data: logs } = useQuery({
		queryKey: ['activity-log'],
		queryFn: async () => {
			// Only fetch local logs
			const stored = localStorage.getItem('activity-log');
			return stored ? JSON.parse(stored) : [];
		},
	});

	return (
		<Card withBorder>
			<Stack>
				<Text size='lg' fw={500}>
					Activity Log
				</Text>
				<Text size='sm' c='dimmed'>
					This log is stored only on your device and helps you track
					feature usage.
				</Text>

				<Timeline>
					{logs?.map((log, i) => (
						<Timeline.Item
							key={i}
							title={log.action}
							bullet={log.success ? '✓' : '✗'}
							color={log.success ? 'green' : 'red'}
						>
							<Text size='sm' c='dimmed'>
								{new Date(log.timestamp).toLocaleString()}
							</Text>
							<Text size='sm'>
								{log.component} - {log.type}
							</Text>
						</Timeline.Item>
					))}
				</Timeline>
			</Stack>
		</Card>
	);
}

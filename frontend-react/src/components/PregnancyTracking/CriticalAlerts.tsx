import { Alert, Stack, Title } from '@mantine/core';
import { Warning, Info, Horse } from '@phosphor-icons/react';
import { useQuery } from '@tanstack/react-query';
import { PregnancyStage } from '../../types/pregnancy';

interface PregnancyAlert {
	id: string;
	type: 'warning' | 'info' | 'critical';
	message: string;
	stage: PregnancyStage;
}

interface CriticalAlertsProps {
	horseId: string;
	currentStage: PregnancyStage;
}

export function CriticalAlerts({ horseId, currentStage }: CriticalAlertsProps) {
	const { data: alerts } = useQuery({
		queryKey: ['pregnancy-alerts', horseId, currentStage],
		queryFn: async () => {
			const response = await fetch(
				`/api/horses/${horseId}/pregnancy/alerts?stage=${currentStage}`
			);
			if (!response.ok) throw new Error('Failed to fetch alerts');
			return response.json() as Promise<PregnancyAlert[]>;
		},
	});

	if (!alerts?.length) return null;

	const getAlertIcon = (type: PregnancyAlert['type']) => {
		switch (type) {
			case 'warning':
				return <Warning size={18} />;
			case 'critical':
				return <Horse size={18} />;
			default:
				return <Info size={18} />;
		}
	};

	const getAlertColor = (type: PregnancyAlert['type']) => {
		switch (type) {
			case 'warning':
				return 'yellow';
			case 'critical':
				return 'red';
			default:
				return 'blue';
		}
	};

	return (
		<Stack>
			<Title order={3}>Important Alerts</Title>
			{alerts.map((alert) => (
				<Alert
					key={alert.id}
					icon={getAlertIcon(alert.type)}
					color={getAlertColor(alert.type)}
					title={alert.type.toUpperCase()}
				>
					{alert.message}
				</Alert>
			))}
		</Stack>
	);
}

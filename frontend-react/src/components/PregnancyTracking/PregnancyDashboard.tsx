import { Stack, Tabs } from '@mantine/core';
import { Checklist, Notes, Stethoscope } from '@phosphor-icons/react';
import { HealthMonitoring } from './HealthMonitoring';
import { PregnancyChecklist } from './PregnancyChecklist';
import { PregnancyStatus } from './PregnancyStatus';
import { VetNotes } from './VetNotes';

interface PregnancyDashboardProps {
	horseId: number;
}

export function PregnancyDashboard({ horseId }: PregnancyDashboardProps) {
	return (
		<Stack>
			{/* Main Status */}
			<PregnancyStatus horseId={horseId} />

			{/* Detailed Information Tabs */}
			<Tabs defaultValue='checklist'>
				<Tabs.List>
					<Tabs.Tab
						value='checklist'
						leftSection={<Checklist size={16} />}
					>
						Pre-Foaling Checklist
					</Tabs.Tab>
					<Tabs.Tab
						value='health'
						leftSection={<Stethoscope size={16} />}
					>
						Health Monitoring
					</Tabs.Tab>
					<Tabs.Tab value='notes' leftSection={<Notes size={16} />}>
						Veterinary Notes
					</Tabs.Tab>
				</Tabs.List>

				<Tabs.Panel value='checklist'>
					<PregnancyChecklist horseId={horseId} />
				</Tabs.Panel>

				<Tabs.Panel value='health'>
					<HealthMonitoring horseId={horseId} />
				</Tabs.Panel>

				<Tabs.Panel value='notes'>
					<VetNotes horseId={horseId} />
				</Tabs.Panel>
			</Tabs>
		</Stack>
	);
}

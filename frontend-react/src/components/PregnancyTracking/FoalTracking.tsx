import {
	Badge,
	Card,
	Group,
	Progress,
	Stack,
	Text,
	Timeline,
} from '@mantine/core';
import { Horse, Ruler, Scales } from '@phosphor-icons/react';
import { useState } from 'react';

interface FoalMilestone {
	id: number;
	date: string;
	type: 'WEIGHT' | 'HEIGHT' | 'VACCINATION' | 'DEVELOPMENT';
	measurement?: number;
	notes: string;
}

interface FoalGrowthData {
	age: number;
	weight: number;
	height: number;
	expectedWeight: number;
	expectedHeight: number;
}

export function FoalTracking({ foalId }: { foalId: number }) {
	const [activeTab, setActiveTab] = useState('overview');

	const { data: foal } = useQuery(['foal', foalId], async () => {
		const response = await fetch(`/api/foals/${foalId}`);
		return response.json();
	});

	const { data: growthData } = useQuery(['foal-growth', foalId], async () => {
		const response = await fetch(`/api/foals/${foalId}/growth`);
		return response.json();
	});

	return (
		<Stack>
			{/* Growth Progress */}
			<Card withBorder>
				<Text fw={500} mb='md'>
					Growth Progress
				</Text>
				<Group grow>
					<div>
						<Group justify='space-between' mb='xs'>
							<Text size='sm'>Weight Progress</Text>
							<Badge>
								{Math.round(
									(growthData.weight /
										growthData.expectedWeight) *
										100
								)}
								%
							</Badge>
						</Group>
						<Progress
							value={
								(growthData.weight /
									growthData.expectedWeight) *
								100
							}
							size='md'
						/>
					</div>
					<div>
						<Group justify='space-between' mb='xs'>
							<Text size='sm'>Height Progress</Text>
							<Badge>
								{Math.round(
									(growthData.height /
										growthData.expectedHeight) *
										100
								)}
								%
							</Badge>
						</Group>
						<Progress
							value={
								(growthData.height /
									growthData.expectedHeight) *
								100
							}
							size='md'
						/>
					</div>
				</Group>
			</Card>

			{/* Development Timeline */}
			<Timeline active={-1}>
				{foal.milestones.map((milestone) => (
					<Timeline.Item
						key={milestone.id}
						bullet={
							milestone.type === 'WEIGHT' ? (
								<Scales size={12} />
							) : milestone.type === 'HEIGHT' ? (
								<Ruler size={12} />
							) : (
								<Horse size={12} />
							)
						}
						title={milestone.type}
					>
						<Text size='sm'>{milestone.notes}</Text>
						{milestone.measurement && (
							<Text size='sm' fw={500}>
								{milestone.measurement}
								{milestone.type === 'WEIGHT' ? 'kg' : 'cm'}
							</Text>
						)}
						<Text size='xs' c='dimmed'>
							{new Date(milestone.date).toLocaleDateString()}
						</Text>
					</Timeline.Item>
				))}
			</Timeline>
		</Stack>
	);
}

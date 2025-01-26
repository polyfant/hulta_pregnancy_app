import { Card, Group, RingProgress, Stack, Text } from '@mantine/core';
import { useQuery } from '@tanstack/react-query';
import { Line } from 'react-chartjs-2';

interface ComparisonData {
	foalId: number;
	age: number;
	weight: number;
	height: number;
	breed: string;
}

export function ComparativeAnalysis({
	foalId,
	breed,
}: {
	foalId: number;
	breed: string;
}) {
	const { data: comparisons } = useQuery({
		queryKey: ['foal-comparisons', breed],
		queryFn: async () => {
			const response = await fetch(
				`/api/analytics/breed-comparisons/${breed}`
			);
			return response.json();
		},
	});

	const { data: currentFoal } = useQuery({
		queryKey: ['foal', foalId],
		queryFn: async () => {
			const response = await fetch(`/api/foals/${foalId}`);
			return response.json();
		},
	});

	const calculatePercentiles = (data: ComparisonData[]) => {
		const weights = data.map((d) => d.weight).sort((a, b) => a - b);
		const heights = data.map((d) => d.height).sort((a, b) => a - b);

		return {
			weightPercentile: getPercentile(currentFoal.weight, weights),
			heightPercentile: getPercentile(currentFoal.height, heights),
		};
	};

	const percentiles = comparisons ? calculatePercentiles(comparisons) : null;

	return (
		<Stack>
			<Card withBorder>
				<Text size='lg' fw={500} mb='md'>
					Breed Comparison
				</Text>
				<Group position='apart'>
					<RingProgress
						size={120}
						roundCaps
						thickness={8}
						sections={[
							{
								value: percentiles?.weightPercentile || 0,
								color: 'blue',
							},
						]}
						label={
							<Text size='xs' ta='center'>
								Weight
								<br />
								{percentiles?.weightPercentile}%
							</Text>
						}
					/>
					<RingProgress
						size={120}
						roundCaps
						thickness={8}
						sections={[
							{
								value: percentiles?.heightPercentile || 0,
								color: 'green',
							},
						]}
						label={
							<Text size='xs' ta='center'>
								Height
								<br />
								{percentiles?.heightPercentile}%
							</Text>
						}
					/>
				</Group>
			</Card>

			<Card withBorder>
				<Text size='lg' fw={500} mb='md'>
					Growth Trajectory
				</Text>
				<Line
					data={{
						labels: comparisons?.map((d) => d.age),
						datasets: [
							{
								label: 'Breed Average',
								data: comparisons?.map((d) => d.weight),
								borderColor: '#2196F3',
								fill: false,
							},
							{
								label: 'Your Foal',
								data: currentFoal?.growthHistory?.map(
									(d) => d.weight
								),
								borderColor: '#4CAF50',
								fill: false,
							},
						],
					}}
				/>
			</Card>

			<Card withBorder>
				<Text size='lg' fw={500} mb='md'>
					Growth Recommendations
				</Text>
				<Stack spacing='xs'>
					{getGrowthRecommendations(percentiles).map((rec, i) => (
						<Text key={i} size='sm'>
							• {rec}
						</Text>
					))}
				</Stack>
			</Card>
		</Stack>
	);
}

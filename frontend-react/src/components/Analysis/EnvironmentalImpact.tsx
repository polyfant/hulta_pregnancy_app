import { Card, Group, Stack, Text, ThemeIcon } from '@mantine/core';
import { Cloud, Warning } from '@phosphor-icons/react';
import { useQuery } from '@tanstack/react-query';
import { environmentalFactors } from '../../services/environmentalFactors';

export function EnvironmentalImpact({ location }: { location: string }) {
	const { data: impact } = useQuery({
		queryKey: ['environmental-impact', location],
		queryFn: async () => {
			const weatherData =
				await environmentalFactors.analyzeEnvironmentalImpact(location);
			return environmentalFactors.calculateEnvironmentalImpact(
				weatherData
			);
		},
		refetchInterval: 1800000, // 30 minutes
	});

	return (
		<Card withBorder>
			<Stack>
				<Group position='apart'>
					<Text size='lg' fw={500}>
						Environmental Conditions
					</Text>
					<Group spacing={4}>
						<ThemeIcon
							size='sm'
							color={
								impact?.growthImpact < -0.2 ? 'red' : 'green'
							}
							variant='light'
						>
							{impact?.growthImpact < -0.2 ? (
								<Warning />
							) : (
								<Cloud />
							)}
						</ThemeIcon>
						<Text size='sm' c='dimmed'>
							Impact Score:{' '}
							{((1 + (impact?.growthImpact || 0)) * 100).toFixed(
								0
							)}
							%
						</Text>
					</Group>
				</Group>

				{impact?.riskFactors.length > 0 && (
					<Stack spacing='xs'>
						<Text size='sm' fw={500} color='red'>
							Risk Factors:
						</Text>
						{impact.riskFactors.map((risk, i) => (
							<Text key={i} size='sm'>
								• {risk}
							</Text>
						))}
					</Stack>
				)}

				<Stack spacing='xs'>
					<Text size='sm' fw={500}>
						Recommendations:
					</Text>
					{impact?.recommendations.map((rec, i) => (
						<Text key={i} size='sm'>
							• {rec}
						</Text>
					))}
				</Stack>
			</Stack>
		</Card>
	);
}

import { Card, Text } from '@mantine/core';
import {
	CartesianGrid,
	Line,
	LineChart,
	ReferenceLine,
	ResponsiveContainer,
	Tooltip,
	XAxis,
	YAxis,
} from 'recharts';
import { calculateVelocity } from '../../utils/growthModels';

interface GrowthVelocityChartProps {
	data: {
		age: number;
		weight: number;
		height: number;
	}[];
}

export function GrowthVelocityChart({ data }: GrowthVelocityChartProps) {
	const weightVelocity = calculateVelocity(
		data.map((d) => ({ age: d.age, value: d.weight }))
	);
	const heightVelocity = calculateVelocity(
		data.map((d) => ({ age: d.age, value: d.height }))
	);

	const combinedData = weightVelocity.map((point, i) => ({
		age: point.age,
		weightGain: point.velocity,
		heightGain: heightVelocity[i].velocity,
	}));

	return (
		<Card withBorder>
			<Text fw={500} mb='md'>
				Growth Velocity
			</Text>
			<ResponsiveContainer width='100%' height={300}>
				<LineChart data={combinedData}>
					<CartesianGrid strokeDasharray='3 3' />
					<XAxis
						dataKey='age'
						label={{ value: 'Age (days)', position: 'bottom' }}
					/>
					<YAxis
						yAxisId='weight'
						label={{
							value: 'Weight Gain (kg/day)',
							angle: -90,
							position: 'insideLeft',
						}}
					/>
					<YAxis
						yAxisId='height'
						orientation='right'
						label={{
							value: 'Height Gain (cm/day)',
							angle: 90,
							position: 'insideRight',
						}}
					/>
					<Tooltip />
					<Line
						yAxisId='weight'
						type='monotone'
						dataKey='weightGain'
						stroke='#2196F3'
						name='Weight Gain'
						dot={false}
					/>
					<Line
						yAxisId='height'
						type='monotone'
						dataKey='heightGain'
						stroke='#4CAF50'
						name='Height Gain'
						dot={false}
					/>
					<ReferenceLine
						y={0.8}
						yAxisId='weight'
						label='Optimal Weight Gain'
						stroke='#FF9800'
						strokeDasharray='3 3'
					/>
				</LineChart>
			</ResponsiveContainer>
		</Card>
	);
}

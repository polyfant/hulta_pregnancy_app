import {
	Alert,
	Badge,
	Button,
	Card,
	Group,
	SimpleGrid,
	Stack,
	Text,
} from '@mantine/core';
import { useMediaQuery } from '@mantine/hooks';
import {
	Calculator,
	FileExport,
	NotePencil,
	Warning,
} from '@phosphor-icons/react';
import { useQuery } from '@tanstack/react-query';
import { Link } from 'react-router-dom';
import {
	CartesianGrid,
	Legend,
	Line,
	LineChart,
	PolarAngleAxis,
	PolarGrid,
	PolarRadiusAxis,
	Radar,
	RadarChart,
	ReferenceLine,
	ResponsiveContainer,
	Tooltip,
	XAxis,
	YAxis,
} from 'recharts';
import { analyzeGrowthTrends } from '../../utils/growthCalculations';

interface GrowthData {
	age: number;
	weight: number;
	height: number;
	expectedWeight: number;
	expectedHeight: number;
}

interface GrowthChartsProps {
	foalId: number;
}

interface GrowthRate {
	weightGainPerDay: number;
	heightGainPerMonth: number;
	isWeightGainHealthy: boolean;
	isHeightGainHealthy: boolean;
	recommendations?: string[];
}

interface BodyCondition {
	score: number; // 1-9 scale
	areas: {
		neck: number;
		withers: number;
		loin: number;
		tailhead: number;
		ribs: number;
		shoulder: number;
	};
	lastUpdated: string;
}

interface GrowthProjection {
	ageInMonths: number;
	projectedWeight: number;
	projectedHeight: number;
	targetWeight: number;
	targetHeight: number;
}

function calculateGrowthRate(data: GrowthData[]): GrowthRate {
	if (data.length < 2) return null;

	// Calculate recent growth (last 30 days)
	const recentData = data.slice(-30);
	const weightGainPerDay =
		(recentData[recentData.length - 1].weight - recentData[0].weight) /
		recentData.length;
	const heightGainPerMonth =
		(recentData[recentData.length - 1].height - recentData[0].height) *
		(30 / recentData.length);

	// Define healthy ranges (these would come from breed standards)
	const healthyWeightGain = { min: 0.7, max: 1.5 }; // kg per day
	const healthyHeightGain = { min: 2.5, max: 5 }; // cm per month

	const isWeightGainHealthy =
		weightGainPerDay >= healthyWeightGain.min &&
		weightGainPerDay <= healthyWeightGain.max;
	const isHeightGainHealthy =
		heightGainPerMonth >= healthyHeightGain.min &&
		heightGainPerMonth <= healthyHeightGain.max;

	// Generate recommendations
	const recommendations = [];
	if (weightGainPerDay < healthyWeightGain.min) {
		recommendations.push(
			'Weight gain is below target. Consider:',
			'• Increasing feed quantity',
			'• Adding quality forage',
			'• Checking for parasites',
			'• Consulting with veterinarian'
		);
	} else if (weightGainPerDay > healthyWeightGain.max) {
		recommendations.push(
			'Weight gain is above target. Consider:',
			'• Adjusting feed portions',
			'• Increasing exercise',
			'• Checking feed quality'
		);
	}

	return {
		weightGainPerDay,
		heightGainPerMonth,
		isWeightGainHealthy,
		isHeightGainHealthy,
		recommendations,
	};
}

function calculateProjections(data: GrowthData[]): GrowthProjection[] {
	// Calculate growth trajectory based on recent data
	const recentGrowth = data.slice(-30);
	const weightGainRate = calculateGrowthRate(data).weightGainPerDay;
	const heightGainRate = calculateGrowthRate(data).heightGainPerMonth / 30;

	const projections: GrowthProjection[] = [];
	const currentAge = data[data.length - 1].age;
	const projectForDays = 90; // Project 3 months ahead

	for (let i = 1; i <= projectForDays; i++) {
		const ageInDays = currentAge + i;
		const ageInMonths = Math.floor(ageInDays / 30);

		projections.push({
			ageInMonths,
			projectedWeight: data[data.length - 1].weight + weightGainRate * i,
			projectedHeight: data[data.length - 1].height + heightGainRate * i,
			targetWeight: data[data.length - 1].expectedWeight + 1.2 * i, // Example target rate
			targetHeight: data[data.length - 1].expectedHeight + 0.3 * i, // Example target rate
		});
	}

	return projections;
}

export function GrowthCharts({ foalId }: GrowthChartsProps) {
	const isMobile = useMediaQuery('(max-width: 768px)');
	const { data: growthData } = useQuery<GrowthData[]>({
		queryKey: ['foal-growth', foalId],
		queryFn: async () => {
			const response = await fetch(`/api/foals/${foalId}/growth-data`);
			return response.json();
		},
	});

	const { data: bodyCondition } = useQuery<BodyCondition>({
		queryKey: ['foal-condition', foalId],
		queryFn: async () => {
			const response = await fetch(`/api/foals/${foalId}/condition`);
			return response.json();
		},
	});

	const exportChart = () => {
		// Implementation for PDF export
	};

	const growthRate = growthData ? calculateGrowthRate(growthData) : null;
	const projections = growthData ? calculateProjections(growthData) : null;
	const growthTrends = growthData ? analyzeGrowthTrends(growthData) : null;

	return (
		<Stack>
			<Group position={isMobile ? 'center' : 'apart'} wrap='nowrap'>
				<Text size={isMobile ? 'lg' : 'xl'} fw={500}>
					Growth Progress
				</Text>
				{!isMobile && (
					<Button
						leftSection={<FileExport size={16} />}
						variant='light'
						onClick={exportChart}
					>
						Export PDF
					</Button>
				)}
			</Group>

			{/* Quick Stats for Mobile */}
			{isMobile && growthRate && (
				<Card withBorder>
					<SimpleGrid cols={2}>
						<div>
							<Text size='sm' c='dimmed'>
								Weight Gain
							</Text>
							<Text fw={500}>
								{growthRate.weightGainPerDay.toFixed(1)} kg/day
							</Text>
						</div>
						<div>
							<Text size='sm' c='dimmed'>
								Height Gain
							</Text>
							<Text fw={500}>
								{growthRate.heightGainPerMonth.toFixed(1)}{' '}
								cm/month
							</Text>
						</div>
					</SimpleGrid>
				</Card>
			)}

			{/* Growth Rate Summary */}
			{growthRate && (
				<Card withBorder>
					<Stack>
						<Group position='apart'>
							<Text fw={500}>
								<Calculator size={20} /> Recent Growth Rate
							</Text>
							<Group>
								<Badge
									color={
										growthRate.isWeightGainHealthy
											? 'green'
											: 'orange'
									}
									variant='light'
								>
									{growthRate.weightGainPerDay.toFixed(2)}{' '}
									kg/day
								</Badge>
								<Badge
									color={
										growthRate.isHeightGainHealthy
											? 'green'
											: 'orange'
									}
									variant='light'
								>
									{growthRate.heightGainPerMonth.toFixed(1)}{' '}
									cm/month
								</Badge>
							</Group>
						</Group>

						{growthRate.recommendations?.length > 0 && (
							<Alert
								icon={<Warning size={16} />}
								color='orange'
								title='Growth Recommendations'
							>
								{growthRate.recommendations.map((rec, i) => (
									<Text
										key={i}
										size='sm'
										mt={i === 0 ? 0 : 'xs'}
									>
										{rec}
									</Text>
								))}
							</Alert>
						)}
					</Stack>
				</Card>
			)}

			<Card withBorder>
				<ResponsiveContainer width='100%' height={isMobile ? 250 : 400}>
					<LineChart
						data={growthData}
						margin={
							isMobile
								? { top: 5, right: 20, left: 0, bottom: 5 }
								: { top: 20, right: 30, left: 20, bottom: 10 }
						}
					>
						<CartesianGrid strokeDasharray='3 3' />
						<XAxis
							dataKey='age'
							label={{ value: 'Age (days)', position: 'bottom' }}
						/>
						<YAxis
							yAxisId='weight'
							label={{
								value: 'Weight (kg)',
								angle: -90,
								position: 'insideLeft',
							}}
						/>
						<YAxis
							yAxisId='height'
							orientation='right'
							label={{
								value: 'Height (cm)',
								angle: 90,
								position: 'insideRight',
							}}
						/>
						<Tooltip
							labelStyle={{ fontSize: isMobile ? 12 : 14 }}
							contentStyle={{ fontSize: isMobile ? 12 : 14 }}
						/>
						<Legend />

						{/* Actual measurements */}
						<Line
							yAxisId='weight'
							type='monotone'
							dataKey='weight'
							stroke='#2196F3'
							name='Weight'
							strokeWidth={2}
							dot={{ r: 4 }}
						/>
						<Line
							yAxisId='height'
							type='monotone'
							dataKey='height'
							stroke='#4CAF50'
							name='Height'
							strokeWidth={2}
							dot={{ r: 4 }}
						/>

						{/* Expected growth curves */}
						<Line
							yAxisId='weight'
							type='monotone'
							dataKey='expectedWeight'
							stroke='#90CAF9'
							strokeDasharray='5 5'
							name='Expected Weight'
						/>
						<Line
							yAxisId='height'
							type='monotone'
							dataKey='expectedHeight'
							stroke='#A5D6A7'
							strokeDasharray='5 5'
							name='Expected Height'
						/>

						{/* Reference lines for breed standards */}
						<ReferenceLine
							y={500}
							yAxisId='weight'
							label='Breed Avg.'
							stroke='#FF9800'
							strokeDasharray='3 3'
						/>
					</LineChart>
				</ResponsiveContainer>
			</Card>

			{/* Growth Milestones */}
			<Card withBorder>
				<Stack>
					<Text fw={500}>Key Milestones</Text>
					<Group grow>
						{growthData?.map((point, index) => {
							if (index % 30 === 0) {
								// Show monthly milestones
								return (
									<Card key={point.age} withBorder>
										<Text fw={500}>{point.age} Days</Text>
										<Text size='sm'>
											Weight: {point.weight}kg
										</Text>
										<Text size='sm'>
											Height: {point.height}cm
										</Text>
									</Card>
								);
							}
							return null;
						})}
					</Group>
				</Stack>
			</Card>

			<Group grow>
				{/* Body Condition Score Radar */}
				<Card withBorder>
					<Text fw={500} mb='md'>
						Body Condition Score
					</Text>
					<ResponsiveContainer width='100%' height={300}>
						<RadarChart data={[bodyCondition?.areas]}>
							<PolarGrid />
							<PolarAngleAxis dataKey='name' />
							<PolarRadiusAxis angle={30} domain={[0, 9]} />
							<Radar
								name='Score'
								dataKey='value'
								stroke='#2196F3'
								fill='#2196F3'
								fillOpacity={0.6}
							/>
						</RadarChart>
					</ResponsiveContainer>
					<Group position='apart' mt='md'>
						<Text>Overall Score: {bodyCondition?.score}/9</Text>
						<Text size='sm' c='dimmed'>
							Last updated:{' '}
							{new Date(
								bodyCondition?.lastUpdated
							).toLocaleDateString()}
						</Text>
					</Group>
				</Card>

				{/* Growth Projections - Enhanced */}
				<Card withBorder>
					<Stack>
						<Group position='apart'>
							<Text fw={500}>Growth Analysis</Text>
							<Group>
								<Badge
									color={
										growthTrends?.weightTrend === 'steady'
											? 'green'
											: 'orange'
									}
									variant='light'
								>
									Weight: {growthTrends?.weightTrend}
								</Badge>
								<Badge
									color={
										growthTrends?.heightTrend === 'steady'
											? 'green'
											: 'orange'
									}
									variant='light'
								>
									Height: {growthTrends?.heightTrend}
								</Badge>
							</Group>
						</Group>

						<Group grow>
							<Card withBorder>
								<Text fw={500}>Current Percentiles</Text>
								<Text>
									Weight: {growthTrends?.percentiles.weight}th
								</Text>
								<Text>
									Height: {growthTrends?.percentiles.height}th
								</Text>
							</Card>

							<Card withBorder>
								<Text fw={500}>Projected Mature Size</Text>
								<Text>
									Weight:{' '}
									{growthTrends?.projectedMaturity.weight.toFixed(
										0
									)}{' '}
									kg
								</Text>
								<Text>
									Height:{' '}
									{growthTrends?.projectedMaturity.height.toFixed(
										0
									)}{' '}
									cm
								</Text>
								<Text size='sm' c='dimmed'>
									Est.{' '}
									{Math.round(
										growthTrends?.projectedMaturity
											.timeToMaturity || 0
									)}{' '}
									days to maturity
								</Text>
							</Card>
						</Group>

						{/* Existing projection chart */}
					</Stack>
				</Card>
			</Group>

			{/* Mobile-friendly recommendations */}
			{isMobile && growthRate?.recommendations?.length > 0 && (
				<Alert
					icon={<Warning size={16} />}
					color='orange'
					title='Recommendations'
					styles={{
						title: { fontSize: 14 },
						message: { fontSize: 12 },
					}}
				>
					{growthRate.recommendations.map((rec, i) => (
						<Text key={i} size='sm' mt={i === 0 ? 0 : 'xs'}>
							{rec}
						</Text>
					))}
				</Alert>
			)}

			{/* Mobile Action Buttons */}
			{isMobile && (
				<Group grow mt='md'>
					<Button
						leftSection={<FileExport size={16} />}
						variant='light'
						onClick={exportChart}
					>
						Export
					</Button>
					<Button
						leftSection={<NotePencil size={16} />}
						variant='light'
						component={Link}
						to={`/foals/${foalId}/notes`}
					>
						Add Note
					</Button>
				</Group>
			)}
		</Stack>
	);
}

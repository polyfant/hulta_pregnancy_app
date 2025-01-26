import { Card, Text } from '@mantine/core';
import { Line } from 'react-chartjs-2';
import { growthPredictions } from '../../services/mlPredictions';

export function GrowthPredictionChart({ historicalData }) {
	const [prediction, setPrediction] = useState(null);

	useEffect(() => {
		const loadPrediction = async () => {
			const result = await growthPredictions.predictSeasonalGrowth(
				historicalData
			);
			setPrediction(result);
		};
		loadPrediction();
	}, [historicalData]);

	const chartData = {
		labels: [...historicalData.map((d) => d.date), 'Next Month'],
		datasets: [
			{
				label: 'Actual Growth',
				data: historicalData.map((d) => d.weight),
				borderColor: '#2196F3',
				fill: false,
			},
			{
				label: 'Predicted Growth',
				data: [
					...Array(historicalData.length - 1).fill(null),
					historicalData[historicalData.length - 1].weight,
					prediction?.nextMonth,
				],
				borderColor: '#4CAF50',
				borderDash: [5, 5],
				fill: false,
			},
		],
	};

	return (
		<Card withBorder>
			<Text size='lg' fw={500} mb='md'>
				Growth Prediction
			</Text>
			<Line data={chartData} />
			{prediction && (
				<Text size='sm' c='dimmed' mt='md'>
					Seasonal adjustment factor:{' '}
					{prediction.seasonalAdjustment.toFixed(2)}
				</Text>
			)}
		</Card>
	);
}

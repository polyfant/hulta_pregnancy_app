import { logger, MLService } from './mlService';

interface GrowthData {
	age: number;
	weight: number;
	height: number;
	date: string;
	temperature: number;
}

interface GrowthPrediction {
	nextMonth: [number, number]; // [weight, height]
	confidence: number;
	seasonalAdjustment: number;
}

interface GrowthPredictionResult {
	forecast: [number, number];
	confidence: number;
	seasonalFactor: number;
}

// Update ModelType to include training type and prediction type
type ModelType = 'GROWTH' | 'HEALTH' | 'BREEDING' | 'TRAIN' | 'PREDICT';

// Define a type for training data
interface TrainingData {
	age: number;
	weight: number;
	height: number;
	temperature: number;
	// Add other relevant fields
}

export class MLPredictions {
	private static readonly SEASONAL_FACTORS = [
		0.8, 0.9, 1.1, 1.2, 1.3, 1.2, 1.0, 0.9, 0.8, 0.7, 0.7, 0.8,
	];

	static async predictSeasonalGrowth(
		historicalData: GrowthData[]
	): Promise<GrowthPrediction> {
		const recentData = this.prepareGrowthData(historicalData);

		try {
			const prediction = (await MLService.predict({
				type: 'growth',
				data: recentData,
				features: [
					'age',
					'weight',
					'height',
					'seasonal',
					'temperature',
				],
			})) as GrowthPredictionResult;

			return {
				nextMonth: prediction.forecast,
				confidence: prediction.confidence,
				seasonalAdjustment: prediction.seasonalFactor,
			};
		} catch (error) {
			logger.error('Growth prediction failed:', error);
			throw new Error('Failed to predict growth pattern');
		}
	}

	private static prepareGrowthData(historicalData: GrowthData[]) {
		return historicalData
			.slice(-90) // Last 90 days
			.map((d) => ({
				age: d.age,
				weight: d.weight,
				height: d.height,
				seasonal: this.getSeasonalFactor(new Date(d.date)),
				temperature: this.getTemperatureFactor(d.temperature),
			}));
	}
	private static getSeasonalFactor(date: Date): number {
		const factor = this.SEASONAL_FACTORS[date.getMonth()];
		if (factor === undefined) {
			throw new Error('Invalid month index');
		}
		return factor;
	}

	private static getTemperatureFactor(temperature: number): number {
		// Normalize temperature to factor between 0.7 and 1.3
		return Math.max(0.7, Math.min(1.3, temperature / 20));
	}
}

export class MLPredictionService {
	private mlService: MLService;
	private trainingData: Map<string, TrainingData[]> = new Map();

	constructor(mlService: MLService) {
		this.mlService = mlService;
	}

	async trainModel(horseId: string, actualData: TrainingData) {
		const existingData = this.trainingData.get(horseId) || [];
		this.trainingData.set(horseId, [...existingData, actualData]);

		// Send training data to backend
		await this.mlService.predict({
			type: 'TRAIN',
			horseId,
			data: actualData,
		});
	}

	async getPrediction(horseId: string) {
		return this.mlService.predict({
			type: 'PREDICT',
			horseId,
		});
	}
}

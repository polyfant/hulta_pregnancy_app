import { mlService } from './mlService';

export const growthPredictions = {
	async predictSeasonalGrowth(historicalData: GrowthData[]) {
		// Extract seasonal patterns and recent growth
		const recentData = historicalData
			.slice(-90)
			.map((d) => [
				d.age,
				d.weight,
				d.height,
				this.getSeasonalFactor(new Date(d.date)),
				this.getTemperatureFactor(d.temperature),
			]);

		const prediction = await mlService.predictGrowth(recentData);
		return {
			nextMonth: prediction[0],
			confidence: prediction[1],
			seasonalAdjustment: prediction[2],
		};
	},

	getSeasonalFactor(date: Date): number {
		const month = date.getMonth();
		// Spring and early summer typically show increased growth
		const seasonalFactors = [
			0.8, 0.9, 1.1, 1.2, 1.3, 1.2, 1.0, 0.9, 0.8, 0.7, 0.7, 0.8,
		];
		return seasonalFactors[month];
	},
};

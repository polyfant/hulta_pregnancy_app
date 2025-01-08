interface DataPoint {
	value: number;
	timestamp: string;
}

export const anomalyDetection = {
	// Z-score based anomaly detection
	detectOutliers(data: DataPoint[], threshold = 2) {
		const values = data.map((d) => d.value);
		const mean = values.reduce((a, b) => a + b) / values.length;
		const std = Math.sqrt(
			values.reduce((sq, n) => sq + Math.pow(n - mean, 2), 0) /
				values.length
		);

		return data.map((point) => ({
			...point,
			isAnomaly: Math.abs((point.value - mean) / std) > threshold,
			confidence: this.calculateConfidence(point, mean, std),
		}));
	},

	// Growth rate analysis
	analyzeGrowthRate(data: DataPoint[]) {
		const sorted = [...data].sort(
			(a, b) =>
				new Date(a.timestamp).getTime() -
				new Date(b.timestamp).getTime()
		);

		return sorted.map((point, i) => {
			if (i === 0) return { ...point, growthRate: 0 };

			const prev = sorted[i - 1];
			const timeDiff =
				(new Date(point.timestamp).getTime() -
					new Date(prev.timestamp).getTime()) /
				(1000 * 60 * 60 * 24);
			const growthRate = (point.value - prev.value) / timeDiff;

			return {
				...point,
				growthRate,
				isAbnormalGrowth: Math.abs(growthRate) > 2, // More than 2 units per day
			};
		});
	},

	// Seasonal pattern detection
	detectSeasonalPatterns(data: DataPoint[], periodDays = 30) {
		const patterns = [];
		for (let i = 0; i < data.length - periodDays; i++) {
			const period = data.slice(i, i + periodDays);
			const periodMean =
				period.reduce((sum, p) => sum + p.value, 0) / periodDays;
			patterns.push({
				startDate: period[0].timestamp,
				endDate: period[periodDays - 1].timestamp,
				meanValue: periodMean,
				variance: this.calculateVariance(period.map((p) => p.value)),
			});
		}
		return patterns;
	},

	calculateConfidence(point: DataPoint, mean: number, std: number): number {
		const zScore = Math.abs((point.value - mean) / std);
		return Math.max(0, 1 - zScore / 3); // Scale confidence 0-1
	},

	calculateVariance(values: number[]): number {
		const mean = values.reduce((a, b) => a + b) / values.length;
		return (
			values.reduce((sq, n) => sq + Math.pow(n - mean, 2), 0) /
			values.length
		);
	},
};

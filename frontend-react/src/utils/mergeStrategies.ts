interface WeightedMeasurement {
	value: number;
	confidence: number;
	source: 'local' | 'server';
	timestamp: string;
}

export const measurementMergeStrategies = {
	weightedAverage(measurements: WeightedMeasurement[]) {
		const grouped = measurements.reduce((acc, m) => {
			if (!acc[m.timestamp]) acc[m.timestamp] = [];
			acc[m.timestamp].push(m);
			return acc;
		}, {});

		return Object.entries(grouped).map(([timestamp, values]) => {
			const totalWeight = values.reduce(
				(sum, m) => sum + m.confidence,
				0
			);
			const weightedValue =
				values.reduce((sum, m) => sum + m.value * m.confidence, 0) /
				totalWeight;

			return {
				timestamp,
				value: weightedValue,
				confidence: Math.max(...values.map((v) => v.confidence)),
			};
		});
	},

	// Smart conflict detection
	detectAnomalies(measurements: WeightedMeasurement[]) {
		const sorted = [...measurements].sort(
			(a, b) =>
				new Date(a.timestamp).getTime() -
				new Date(b.timestamp).getTime()
		);

		return sorted.map((m, i) => {
			if (i === 0) return { ...m, confidence: 1 };

			const prev = sorted[i - 1];
			const timeDiff =
				(new Date(m.timestamp).getTime() -
					new Date(prev.timestamp).getTime()) /
				(1000 * 60 * 60 * 24);
			const changeRate = (m.value - prev.value) / timeDiff;

			// Adjust confidence based on rate of change
			const isUnusual = Math.abs(changeRate) > 2; // More than 2 units per day
			return {
				...m,
				confidence: isUnusual ? 0.5 : 1,
			};
		});
	},
};

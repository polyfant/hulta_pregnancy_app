interface GrowthTrend {
	weightTrend: 'accelerating' | 'steady' | 'slowing';
	heightTrend: 'accelerating' | 'steady' | 'slowing';
	percentiles: {
		weight: number;
		height: number;
	};
	projectedMaturity: {
		weight: number;
		height: number;
		timeToMaturity: number;
	};
}

export function analyzeGrowthTrends(data: GrowthData[]): GrowthTrend {
	// Analyze last 3 months of data for trends
	const recentData = data.slice(-90);
	const periods = chunk(recentData, 30); // Split into 30-day periods

	// Calculate growth rates for each period
	const rates = periods.map((period, i) => ({
		weight:
			(period[period.length - 1].weight - period[0].weight) /
			period.length,
		height:
			(period[period.length - 1].height - period[0].height) /
			period.length,
	}));

	// Determine if growth is accelerating, steady, or slowing
	const weightTrend = determineTrend(rates.map((r) => r.weight));
	const heightTrend = determineTrend(rates.map((r) => r.height));

	// Calculate percentiles based on breed standards
	const currentAge = data[data.length - 1].age;
	const percentiles = calculatePercentiles(
		data[data.length - 1].weight,
		data[data.length - 1].height,
		currentAge
	);

	// Project mature size based on current growth curve
	const maturity = projectMaturity(data);

	return {
		weightTrend,
		heightTrend,
		percentiles,
		projectedMaturity: maturity,
	};
}

function determineTrend(
	rates: number[]
): 'accelerating' | 'steady' | 'slowing' {
	const changes = rates.slice(1).map((rate, i) => rate - rates[i]);
	const avgChange = average(changes);

	if (avgChange > 0.1) return 'accelerating';
	if (avgChange < -0.1) return 'slowing';
	return 'steady';
}

function projectMaturity(data: GrowthData[]): {
	weight: number;
	height: number;
	timeToMaturity: number;
} {
	// Use growth curve modeling to project final size
	const weightCurve = fitGrowthCurve(
		data.map((d) => ({ age: d.age, value: d.weight }))
	);
	const heightCurve = fitGrowthCurve(
		data.map((d) => ({ age: d.age, value: d.height }))
	);

	return {
		weight: weightCurve.asymptote,
		height: heightCurve.asymptote,
		timeToMaturity: Math.max(
			weightCurve.timeToMaturity,
			heightCurve.timeToMaturity
		),
	};
}

function fitGrowthCurve(data: { age: number; value: number }[]) {
	// Implement Gompertz or von Bertalanffy growth curve fitting
	// Returns asymptote and time to reach 95% of asymptote
	// This is a simplified example
	const maxValue = Math.max(...data.map((d) => d.value));
	const growthRate = calculateGrowthRate(data);

	return {
		asymptote: maxValue * 1.2, // Estimate mature size
		timeToMaturity:
			(maxValue * 1.2 - data[data.length - 1].value) / growthRate,
	};
}

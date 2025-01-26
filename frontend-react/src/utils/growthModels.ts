interface GrowthParameters {
	maxValue: number; // Asymptotic value (mature size)
	growthRate: number; // Growth rate parameter
	inflectionPoint: number; // Age at inflection point
}

export function gompertzGrowth(age: number, params: GrowthParameters): number {
	const { maxValue, growthRate, inflectionPoint } = params;
	return (
		maxValue * Math.exp(-Math.exp(-growthRate * (age - inflectionPoint)))
	);
}

export function calculateVelocity(
	data: { age: number; value: number }[]
): { age: number; velocity: number }[] {
	return data.slice(1).map((point, i) => ({
		age: point.age,
		velocity: (point.value - data[i].value) / (point.age - data[i].age),
	}));
}

export function fitGompertzCurve(
	data: { age: number; value: number }[]
): GrowthParameters {
	// Simple parameter estimation
	const maxObserved = Math.max(...data.map((d) => d.value));
	const maxValue = maxObserved * 1.2; // Estimate mature size

	// Find point of maximum growth rate
	const velocities = calculateVelocity(data);
	const maxVelocityPoint = velocities.reduce((max, point) =>
		point.velocity > max.velocity ? point : max
	);

	return {
		maxValue,
		growthRate: 0.015, // Typical value for horses
		inflectionPoint: maxVelocityPoint.age,
	};
}

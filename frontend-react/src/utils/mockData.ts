export function generateMockHorseData(count: number) {
	return Array.from({ length: count }, () => ({
		age: Math.floor(Math.random() * 20) + 1,
		weight: 400 + Math.random() * 200, // 400-600kg
		height: 1.4 + Math.random() * 0.4, // 1.4-1.8m
		temperature: 37 + Math.random(), // 37-38°C
		heartRate: 28 + Math.floor(Math.random() * 12), // 28-40 bpm
		pregnancyData: {
			gestationDay: Math.floor(Math.random() * 340),
			fetalHeartRate: 100 + Math.floor(Math.random() * 50),
			placentalThickness: 1 + Math.random() * 0.5,
		},
	}));
}

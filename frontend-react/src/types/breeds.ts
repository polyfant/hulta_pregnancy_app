export interface BreedStandard {
	id: string;
	name: string;
	category: 'Light' | 'Draft' | 'Pony';
	growthStandards: {
		ageInMonths: number;
		weightRange: { min: number; max: number };
		heightRange: { min: number; max: number };
	}[];
	bodyConditionGuidelines: {
		ideal: { min: number; max: number };
		warning: { min: number; max: number };
	};
}

export const COMMON_BREEDS: BreedStandard[] = [
	{
		id: 'thoroughbred',
		name: 'Thoroughbred',
		category: 'Light',
		growthStandards: [
			{
				ageInMonths: 1,
				weightRange: { min: 70, max: 90 },
				heightRange: { min: 105, max: 115 },
			},
			// ... more age ranges
		],
		bodyConditionGuidelines: {
			ideal: { min: 4, max: 6 },
			warning: { min: 3, max: 7 },
		},
	},
	// Add more breeds as needed
];

import { gompertzGrowth } from '../utils/growthModels';

const createGrowthData = (params: {
	weightParams: GrowthParameters;
	heightParams: GrowthParameters;
}) => {
	return Array.from({ length: 180 }, (_, i) => ({
		age: i,
		weight: gompertzGrowth(i, params.weightParams),
		height: gompertzGrowth(i, params.heightParams),
		expectedWeight: gompertzGrowth(i, STANDARD_GROWTH.weightParams),
		expectedHeight: gompertzGrowth(i, STANDARD_GROWTH.heightParams),
	}));
};

const STANDARD_GROWTH = {
	weightParams: {
		maxValue: 500,
		growthRate: 0.015,
		inflectionPoint: 90,
	},
	heightParams: {
		maxValue: 160,
		growthRate: 0.02,
		inflectionPoint: 60,
	},
};

export const MOCK_GROWTH_DATA = {
	thoroughbred: {
		normal: createGrowthData(STANDARD_GROWTH),
		slow: createGrowthData({
			weightParams: {
				...STANDARD_GROWTH.weightParams,
				growthRate: 0.012,
			},
			heightParams: {
				...STANDARD_GROWTH.heightParams,
				growthRate: 0.015,
			},
		}),
		rapid: createGrowthData({
			weightParams: {
				...STANDARD_GROWTH.weightParams,
				growthRate: 0.018,
			},
			heightParams: {
				...STANDARD_GROWTH.heightParams,
				growthRate: 0.025,
			},
		}),
	},
};

export const MOCK_BODY_CONDITION = {
	normal: {
		score: 5,
		areas: {
			neck: 5,
			withers: 5,
			loin: 5,
			tailhead: 5,
			ribs: 5,
			shoulder: 5,
		},
		lastUpdated: new Date().toISOString(),
	},
	thin: {
		score: 3,
		areas: {
			neck: 3,
			withers: 2,
			loin: 3,
			tailhead: 3,
			ribs: 2,
			shoulder: 3,
		},
		lastUpdated: new Date().toISOString(),
	},
};

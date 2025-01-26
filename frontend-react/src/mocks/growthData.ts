import { gompertzGrowth } from '../utils/growthModels';

// Growth parameters representing different breed characteristics
const BREED_GROWTH_PARAMS = {
    thoroughbred: {
        weightParams: {
            maxValue: 550,  // Typical max weight for thoroughbred
            growthRate: 0.015,
            inflectionPoint: 90,
        },
        heightParams: {
            maxValue: 170,  // Typical max height for thoroughbred
            growthRate: 0.02,
            inflectionPoint: 60,
        },
    },
    warmblood: {
        weightParams: {
            maxValue: 650,  // Larger breed, higher max weight
            growthRate: 0.013,
            inflectionPoint: 100,
        },
        heightParams: {
            maxValue: 180,  // Taller breed
            growthRate: 0.018,
            inflectionPoint: 70,
        },
    },
    arabian: {
        weightParams: {
            maxValue: 450,  // Typically lighter breed
            growthRate: 0.016,
            inflectionPoint: 85,
        },
        heightParams: {
            maxValue: 150,  // Generally shorter
            growthRate: 0.019,
            inflectionPoint: 55,
        },
    }
};

const createGrowthData = (params: {
    weightParams: GrowthParameters;
    heightParams: GrowthParameters;
}) => {
    return Array.from({ length: 180 }, (_, i) => ({
        age: i,
        weight: gompertzGrowth(i, params.weightParams),
        height: gompertzGrowth(i, params.heightParams),
        expectedWeight: gompertzGrowth(i, params.weightParams),
        expectedHeight: gompertzGrowth(i, params.heightParams),
    }));
};

// Comprehensive mock growth data scenarios
export const MOCK_GROWTH_DATA = {
    thoroughbred: {
        // Normal growth trajectory for a healthy thoroughbred foal
        normal: createGrowthData(BREED_GROWTH_PARAMS.thoroughbred),
        
        // Slower growth scenario (potential nutrition or health issues)
        slow: createGrowthData({
            weightParams: {
                ...BREED_GROWTH_PARAMS.thoroughbred.weightParams,
                growthRate: 0.012,  // Reduced growth rate
                maxValue: 500,      // Lower max weight
            },
            heightParams: {
                ...BREED_GROWTH_PARAMS.thoroughbred.heightParams,
                growthRate: 0.015,  // Slower height increase
                maxValue: 160,      // Lower max height
            },
        }),
        
        // Rapid growth scenario (excellent nutrition, genetics)
        rapid: createGrowthData({
            weightParams: {
                ...BREED_GROWTH_PARAMS.thoroughbred.weightParams,
                growthRate: 0.018,  // Accelerated growth
                maxValue: 600,      // Higher potential max weight
            },
            heightParams: {
                ...BREED_GROWTH_PARAMS.thoroughbred.heightParams,
                growthRate: 0.025,  // Faster height increase
                maxValue: 175,      // Higher potential max height
            },
        }),
    },
    warmblood: {
        normal: createGrowthData(BREED_GROWTH_PARAMS.warmblood),
        slow: createGrowthData({
            weightParams: {
                ...BREED_GROWTH_PARAMS.warmblood.weightParams,
                growthRate: 0.011,
                maxValue: 600,
            },
            heightParams: {
                ...BREED_GROWTH_PARAMS.warmblood.heightParams,
                growthRate: 0.016,
                maxValue: 170,
            },
        }),
        rapid: createGrowthData({
            weightParams: {
                ...BREED_GROWTH_PARAMS.warmblood.weightParams,
                growthRate: 0.016,
                maxValue: 700,
            },
            heightParams: {
                ...BREED_GROWTH_PARAMS.warmblood.heightParams,
                growthRate: 0.022,
                maxValue: 190,
            },
        }),
    },
    arabian: {
        normal: createGrowthData(BREED_GROWTH_PARAMS.arabian),
        slow: createGrowthData({
            weightParams: {
                ...BREED_GROWTH_PARAMS.arabian.weightParams,
                growthRate: 0.013,
                maxValue: 420,
            },
            heightParams: {
                ...BREED_GROWTH_PARAMS.arabian.heightParams,
                growthRate: 0.016,
                maxValue: 145,
            },
        }),
        rapid: createGrowthData({
            weightParams: {
                ...BREED_GROWTH_PARAMS.arabian.weightParams,
                growthRate: 0.019,
                maxValue: 480,
            },
            heightParams: {
                ...BREED_GROWTH_PARAMS.arabian.heightParams,
                growthRate: 0.022,
                maxValue: 160,
            },
        }),
    }
};

// Mock body condition scenarios
export const MOCK_BODY_CONDITION = {
    normal: {
        bodyConditionScore: 5,  // Ideal condition
        muscleScore: 6,
        fatScore: 5,
        notes: 'Healthy and well-proportioned'
    },
    underweight: {
        bodyConditionScore: 3,  // Below ideal
        muscleScore: 4,
        fatScore: 3,
        notes: 'Potential nutritional concerns, monitor closely'
    },
    overweight: {
        bodyConditionScore: 7,  // Above ideal
        muscleScore: 5,
        fatScore: 7,
        notes: 'Potential risk of metabolic issues, adjust diet'
    }
};

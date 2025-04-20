// Growth-related predictions
export interface GrowthPrediction {
	predictedWeight: number; // in kg
	predictedHeight: number; // in meters
	confidence: number; // 0-1 scale
	growthCurve: number[]; // daily predictions
	seasonalFactors: {
		temperature: number; // impact coefficient
		season: number; // seasonal adjustment
		nutrition: number; // nutrition factor
	};
	milestones: {
		date: Date;
		milestone: string;
		expectedValue: number;
	}[];
	nextMilestone: string;
}

// Health-related predictions
export interface HealthPrediction {
	healthScore: number; // 0-100 scale
	riskFactors: string[];
	recommendations: string[];
	vitalSigns: {
		temperature: number; // in Celsius
		heartRate: number; // BPM
		respiratoryRate: number; // breaths per minute
	};
	nutritionalStatus: {
		bodyCondition: number; // 1-9 scale
		dietaryNeeds: string[];
		supplements: string[];
	};
	riskLevel: 'LOW' | 'MEDIUM' | 'HIGH';
}

// Breeding-related suggestions
export interface BreedingSuggestion {
	optimalDates: Date[];
	successProbability: number; // 0-1 scale
	geneticCompatibility: number; // 0-1 scale
	recommendedPartners: {
		horseId: string;
		compatibility: number;
		traits: string[];
	}[];
	healthConsiderations: string[];
	seasonalFactors: {
		optimal: boolean;
		considerations: string[];
	};
	compatibilityScore: number;
	geneticDiversity: number;
	potentialRisks: string[];
}

// Risk analysis
export interface RiskAnalysis {
	overallRisk: number; // 0-1 scale
	factors: {
		genetic: number;
		environmental: number;
		medical: number;
	};
	recommendations: string[];
}

// Environmental impact assessment
export interface EnvironmentalImpact {
	carbonFootprint: number; // in kg CO2
	waterUsage: number; // in liters
	landUse: number; // in square meters
	sustainabilityScore: number; // 0-100 scale
	recommendations: {
		action: string;
		potentialImpact: number;
		timeframe: string;
	}[];
	seasonalVariation: {
		season: string;
		impact: number;
	}[];
	stressFactors: string[];
	optimalConditions: {
		temperature: [number, number];
		humidity: [number, number];
	};
}

// Common types used across predictions
export interface PredictionMetadata {
	timestamp: Date;
	confidence: number;
	dataPoints: number;
	modelVersion: string;
}

export interface MLResponse<T> {
	data: T;
	metadata: PredictionMetadata;
}

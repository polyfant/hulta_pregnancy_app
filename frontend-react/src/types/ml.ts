export type ModelType = 'GROWTH' | 'HEALTH' | 'BREEDING';

export interface MLConfig {
	endpoint: string;
	apiKey: string;
}

export interface ModelMetadata {
	version: string;
	accuracy: number;
	lastUpdated: Date;
}

export interface TrainingData {
	features: Record<string, number | string>;
	labels: {
		pregnancySuccess?: boolean;
		healthScore?: number;
		growthRate?: number;
	};
}

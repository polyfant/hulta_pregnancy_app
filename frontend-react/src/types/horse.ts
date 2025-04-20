export interface Horse {
	id: string;
	name: string;
	breed?: string;
	gender?: string;
	age?: number;
	status: 'active' | 'retired' | 'sold';
	isPregnant?: boolean;
	isPremium?: boolean;
	isChampion?: boolean;
	nextCheckup?: string;
	lastBreedingDate?: string;
	dueDate?: string;
	healthStatus?: 'excellent' | 'good' | 'fair' | 'poor';
	notes?: string;
}

export interface CreateHorseInput {
	name: string;
	breed: string;
	gender: 'MARE' | 'STALLION' | 'GELDING';
	dateOfBirth: string;
	weight?: number;
	isPregnant?: boolean;
	conceptionDate?: string;
	motherId?: number;
	fatherId?: number;
	externalMother?: string;
	externalFather?: string;
}

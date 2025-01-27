
export * from './pregnancy';
export * from './horse';

export interface GrowthData {
    timestamp: number;
    weight: number;
    height: number;
    temperature?: number;
  }
  
  export interface HealthIndicators {
    heartRate: number;
    bodyTemperature: number;
    // Add more as needed
  }
  
  export interface FeedingProgram {
    type: string;
    quantity: number;
    frequency: number;
  }

  // Core data models for the Hulta Pregnancy App

export interface GrowthData {
    timestamp: number;
    weight: number;
    height: number;
    temperature?: number;
  }
  
  export interface HealthIndicators {
    heartRate: number;
    bodyTemperature: number;
    respiratoryRate?: number;
  }
  
  export interface FeedingProgram {
    type: string;
    quantity: number;
    frequency: number;
    nutritionalContent?: Record<string, number>;
  }
  
  export interface PregnancyStatus {
    isPregnant: boolean
    currentStage: 'EARLY' | 'MIDDLE' | 'LATE' | 'NEARTERM' | 'FOALING'
    conceptionDate: string
    expectedDueDate: string
    daysInPregnancy: number
    healthRisk: 'low' | 'medium' | 'high'  
}
  
  export interface EnvironmentalFactors {
    temperature: number;
    humidity: number;
    airQuality?: number;
  }
  export interface Horse {
    id: string;
    name: string;
    breed: string;
    dateOfBirth: string;
    sex: 'mare' | 'stallion' | 'gelding';
    color: string;
    owner: {
        id: string;
        name: string;
    };
    // Optional fields
    registrationNumber?: string;
    microchipNumber?: string;
    notes?: string;
}
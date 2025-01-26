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
    startDate: Date;
    expectedDeliveryDate: Date;
    currentStage: 'early' | 'mid' | 'late';
    healthRisk: 'low' | 'medium' | 'high';
}
  
  export interface EnvironmentalFactors {
    temperature: number;
    humidity: number;
    airQuality?: number;
  }

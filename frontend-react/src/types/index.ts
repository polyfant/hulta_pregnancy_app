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
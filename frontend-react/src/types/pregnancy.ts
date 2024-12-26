export interface PregnancyStatus {
  isPregnant: boolean;
  conceptionDate: string;
  currentStage: 'EARLY' | 'MIDDLE' | 'LATE' | 'NEARTERM' | 'FOALING';
  daysInPregnancy: number;
  expectedDueDate: string;
  lastEvent?: PregnancyEvent;
}

export interface PregnancyEvent {
  id: number;
  date: string;
  eventType: string;
  description: string;
  notes?: string;
}

export interface PregnancyGuideline {
  stage: string;
  title: string;
  description: string;
  recommendations: string[];
  warnings: string[];
  checkpoints: string[];
}

export interface PreFoalingSign {
  id: number;
  signName: string;
  observed: boolean;
  dateObserved?: string;
  notes?: string;
}

export type PregnancyStage = 'Early' | 'Mid' | 'Late' | 'Pre-foaling';

export interface PregnancyStatus {
  horseId: number;
  currentDay: number;
  totalDays: number;
  stage: PregnancyStage;
  dueDate: string;
  lastCheckup?: string;
  nextCheckup?: string;
  notes?: string;
}

export interface PregnancyEvent {
  id: string;
  horseId: number;
  date: string;
  type: 'checkup' | 'milestone' | 'note';
  title: string;
  description: string;
}

export interface ChecklistItem {
  id: string;
  text: string;
  completed: boolean;
  createdAt: string;
  completedAt?: string;
}

export interface PregnancyChecklist {
  id: string;
  horseId: number;
  items: ChecklistItem[];
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

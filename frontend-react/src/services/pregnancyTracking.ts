import { PregnancyStatus } from '../types';
import { HealthIndicators } from '../types';
import AuditLogger from './auditLogger';

export class PregnancyTrackingService {
    private pregnancyData: PregnancyStatus | null = null;
    private healthRecords: HealthIndicators[] = [];
    startPregnancy(conceptionDate: string): void {
        const dueDate = new Date(conceptionDate);
        dueDate.setDate(dueDate.getDate() + 340); // ~11 months for horses
    
        this.pregnancyData = {
            isPregnant: true,
            currentStage: 'EARLY',
            conceptionDate: conceptionDate,
            expectedDueDate: dueDate.toISOString(),
            daysInPregnancy: this.calculateDaysInPregnancy(conceptionDate),
            healthRisk: 'low'
        };
        
        AuditLogger.log('INFO', 'Pregnancy tracking started', {
            conceptionDate: conceptionDate,
            action: 'START_PREGNANCY'
        });
    }

    private calculateDaysInPregnancy(conceptionDate: string): number {
        const start = new Date(conceptionDate);
        const today = new Date();
        const diffTime = Math.abs(today.getTime() - start.getTime());
        return Math.ceil(diffTime / (1000 * 60 * 60 * 24));
    }

    // Update stage calculation to match our enum
    private calculatePregnancyStage(daysInPregnancy: number): PregnancyStatus['currentStage'] {
        if (daysInPregnancy < 114) return 'EARLY';      // First trimester
        if (daysInPregnancy < 228) return 'MIDDLE';     // Second trimester
        if (daysInPregnancy < 310) return 'LATE';       // Third trimester
        if (daysInPregnancy < 330) return 'NEARTERM';   // Near term
        return 'FOALING';                               // Ready to foal
    }

    assessHealthRisk(indicators: HealthIndicators): void {
        if (!this.pregnancyData) return;
    
        let risk: 'low' | 'medium' | 'high' = 'low';  // Default to low
    
        // Example risk assessment logic based on health indicators
        if (indicators.heartRate > 80 || indicators.bodyTemperature > 38.5) {
            risk = 'high';
        } else if (indicators.heartRate > 60 || indicators.bodyTemperature > 38.0) {
            risk = 'medium';
        }
    
        this.pregnancyData.healthRisk = risk;
        this.healthRecords.push(indicators);  // Store the health record
    
        AuditLogger.log('INFO', 'Health risk assessed', {
            risk,
            indicators
        });
    }

    getCurrentStatus(): PregnancyStatus | null {
        if (!this.pregnancyData) return null;
        
        // Update stage based on current days
        const daysInPregnancy = this.calculateDaysInPregnancy(this.pregnancyData.conceptionDate);
        this.pregnancyData.daysInPregnancy = daysInPregnancy;
        this.pregnancyData.currentStage = this.calculatePregnancyStage(daysInPregnancy);
        
        return this.pregnancyData;
    }


    // Get health history
    getHealthHistory(): HealthIndicators[] {
        return this.healthRecords;
    }

    // Estimate delivery probability
    estimateDeliveryProbability(): number {
        if (!this.pregnancyData) return 0;
    
        const daysRemaining = (new Date(this.pregnancyData.expectedDueDate).getTime() - new Date().getTime()) / (1000 * 3600 * 24);
        const totalPregnancyDays = (new Date(this.pregnancyData.expectedDueDate).getTime() - new Date(this.pregnancyData.conceptionDate).getTime()) / (1000 * 3600 * 24);
    
        return Math.max(0, Math.min(1, 1 - (daysRemaining / totalPregnancyDays)));
    }
}

export default new PregnancyTrackingService();
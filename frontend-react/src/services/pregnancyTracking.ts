// src/services/pregnancyTrackingService.ts
import { PregnancyStatus, HealthIndicators } from '../types';
import AuditLogger from './auditLogger';

export class PregnancyTrackingService {
    // Core pregnancy tracking data
    private pregnancyData: PregnancyStatus | null = null;
    private healthRecords: HealthIndicators[] = [];

    // Initialize pregnancy tracking
    startPregnancy(startDate: Date, expectedDeliveryDate: Date): void {
        this.pregnancyData = {
            startDate,
            expectedDeliveryDate,
            currentStage: this.calculatePregnancyStage(startDate),
            healthRisk: 'low'  // Default initial risk
        };

        AuditLogger.log('PREGNANCY', 'Tracking Started', {
            startDate: startDate.toISOString(),
            expectedDelivery: expectedDeliveryDate.toISOString()
        });
    }

    // Record health indicators
    recordHealthIndicators(indicators: HealthIndicators): void {
        if (!this.pregnancyData) {
            throw new Error('Pregnancy tracking not initialized');
        }

        this.healthRecords.push(indicators);
        this.assessHealthRisk(indicators);

        AuditLogger.log('HEALTH', 'Indicators Recorded', indicators);
    }

    // Calculate current pregnancy stage
    private calculatePregnancyStage(startDate: Date): 'early' | 'mid' | 'late' {
        const currentDate = new Date();
        const daysSinceStart = (currentDate.getTime() - startDate.getTime()) / (1000 * 3600 * 24);

        if (daysSinceStart < 90) return 'early';
        if (daysSinceStart < 270) return 'mid';
        return 'late';
    }

    // Assess health risk based on indicators
    private assessHealthRisk(indicators: HealthIndicators): void {
        const { heartRate, bodyTemperature } = indicators;

        // Basic risk assessment logic
        if (
            heartRate > 100 || 
            heartRate < 60 || 
            bodyTemperature > 38 || 
            bodyTemperature < 36
        ) {
            this.pregnancyData!.healthRisk = 'high';
            AuditLogger.log('ALERT', 'Potential Health Risk Detected', indicators);
        }
    }

    // Get current pregnancy status
    getCurrentStatus(): PregnancyStatus | null {
        return this.pregnancyData;
    }

    // Get health history
    getHealthHistory(): HealthIndicators[] {
        return this.healthRecords;
    }

    // Estimate delivery probability
    estimateDeliveryProbability(): number {
        if (!this.pregnancyData) return 0;

        const daysRemaining = (this.pregnancyData.expectedDeliveryDate.getTime() - new Date().getTime()) / (1000 * 3600 * 24);
        const totalPregnancyDays = (this.pregnancyData.expectedDeliveryDate.getTime() - this.pregnancyData.startDate.getTime()) / (1000 * 3600 * 24);

        return Math.max(0, Math.min(1, 1 - (daysRemaining / totalPregnancyDays)));
    }
}

export default new PregnancyTrackingService();
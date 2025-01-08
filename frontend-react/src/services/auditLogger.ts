interface AuditEvent {
    type: string;
    timestamp: number;
    component: string;
    action: string;
    success: boolean;
    anonymousMetrics?: {
        duration?: number;
        errorType?: string;
        featureUsed?: string;
    };
}

export const auditLogger = {
    private queue: AuditEvent[] = [],
    private processingInterval: number = 1800000, // 30 minutes

    async logEvent(event: Omit<AuditEvent, 'timestamp'>) {
        // Never log personal data
        const sanitizedEvent = this.sanitizeEvent(event);
        
        this.queue.push({
            ...sanitizedEvent,
            timestamp: Date.now()
        });

        // Process queue if it gets too large
        if (this.queue.length > 50) {
            await this.processQueue();
        }
    },

    private sanitizeEvent(event: any) {
        // Remove any potential PII
        const sanitized = { ...event };
        delete sanitized.userId;
        delete sanitized.horseId;
        delete sanitized.location;
        delete sanitized.deviceId;
        return sanitized;
    },

    async processQueue() {
        if (this.queue.length === 0) return;

        // Aggregate events by type
        const aggregated = this.queue.reduce((acc, event) => {
            const key = `${event.component}-${event.action}`;
            if (!acc[key]) {
                acc[key] = {
                    count: 0,
                    successes: 0,
                    averageDuration: 0,
                    errorTypes: {}
                };
            }
            acc[key].count++;
            if (event.success) acc[key].successes++;
            if (event.anonymousMetrics?.errorType) {
                acc[key].errorTypes[event.anonymousMetrics.errorType] = 
                    (acc[key].errorTypes[event.anonymousMetrics.errorType] || 0) + 1;
            }
            return acc;
        }, {});

        // Send aggregated data
        if (navigator.onLine) {
            try {
                await fetch('/api/analytics/anonymous', {
                    method: 'POST',
                    body: JSON.stringify({
                        timestamp: Date.now(),
                        data: aggregated,
                        // Add noise to make it harder to identify users
                        noise: Math.random()
                    })
                });
                this.queue = [];
            } catch (error) {
                console.error('Failed to send analytics');
            }
        }
    },

    startPeriodicProcessing() {
        setInterval(() => this.processQueue(), this.processingInterval);
    }
}; 
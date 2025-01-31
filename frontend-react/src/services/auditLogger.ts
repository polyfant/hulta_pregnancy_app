export interface LogEntry {
    category: string;
    action: string;
    details?: Record<string, any>;
    timestamp: number;
}

export class AuditLogger {
    private static logs: LogEntry[] = [];

    static log(category: string, action: string, details?: Record<string, any>): void {
        const entry: LogEntry = {
            category,
            action,
            details: details || {},
            timestamp: Date.now()
        };
        
        this.logs.push(entry);
        console.log('Audit Log:', entry);
    }

    static getLogs(): LogEntry[] {
        return this.logs;
    }

    static clearLogs(): void {
        this.logs = [];
    }
}

export default AuditLogger;
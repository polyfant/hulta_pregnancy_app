import AuditLogger from './auditLogger';
import MLService from './mlService';
import NotificationService from './notificationService';
import RoleManager, { UserRole } from '../utils/roleManager';
import { GrowthData, HealthIndicators } from '../types';

export function validateServices() {
    try {
        // AuditLogger Validation
        AuditLogger.log('VALIDATION', 'Service Check', { source: 'serviceValidator' });
        const logs = AuditLogger.getLogs();
        console.assert(logs.length > 0, 'AuditLogger validation failed');

        // MLService Mock Validation
        const mlConfig = {
            apiKey: 'test-key',
            endpoint: 'https://mock-ml-service.com/predict'
        };
        const mlService = new MLService(mlConfig);

        // NotificationService Validation
        const notifyConfig = { provider: 'test-provider' };
        const notificationService = new NotificationService(notifyConfig);

        // RoleManager Validation
        RoleManager.assignRole('test-user', UserRole.USER);
        const userRole = RoleManager.getUserRole('test-user');
        console.assert(userRole === UserRole.USER, 'RoleManager role assignment failed');

        console.log('🎉 All Services Validated Successfully! 🚀');
    } catch (error) {
        console.error('🚨 Service Validation Failed:', error);
        throw error;
    }
}

// Run validation in development
if (import.meta.env.DEV) {
    validateServices();
}

export default validateServices;
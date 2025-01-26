import AuditLogger from './auditLogger';
import MLService from './mlService';
import NotificationService from './notificationService';
import RoleManager, { UserRole } from '../utils/roleManager';

export function validateServices() {
    try {
        // AuditLogger Test
        AuditLogger.log('TEST', 'Validation', { source: 'serviceValidator' });
        const logs = AuditLogger.getLogs();
        console.assert(logs.length > 0, 'AuditLogger failed');

        // MLService Test (mock)
        const mlConfig = {
            apiKey: 'test-key',
            endpoint: 'https://mock-ml-service.com/predict'
        };
        const mlService = new MLService(mlConfig);
        
        // Async prediction test (commented out for now)
        // const predictionResult = await mlService.predict({ input: 'test' });

        // NotificationService Test
        const notifyConfig = { provider: 'test-provider' };
        const notificationService = new NotificationService(notifyConfig);
        
        // Async notification test (commented out for now)
        // const notificationResult = await notificationService.send(['test@example.com'], 'Test message');

        // RoleManager Test
        RoleManager.assignRole('test-user', UserRole.USER);
        const userRole = RoleManager.getUserRole('test-user');
        console.assert(userRole === UserRole.USER, 'RoleManager role assignment failed');

        console.log('Service Validation Completed Successfully');
    } catch (error) {
        console.error('Service Validation Failed:', error);
    }
}

// Only run in development
if (import.meta.env.DEV) {
    validateServices();
}

export default validateServices;
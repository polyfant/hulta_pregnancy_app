export interface MLConfig {
    apiKey: string;
    endpoint: string;
}

export class MLService {
    private config: MLConfig;

    constructor(config: MLConfig) {
        this.config = config;
    }

    async predict(data: any): Promise<any> {
        try {
            const response = await fetch(this.config.endpoint, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${this.config.apiKey}`
                },
                body: JSON.stringify(data)
            });

            return await response.json();
        } catch (error) {
            console.error('ML Service Error:', error);
            throw error;
        }
    }
}

export default MLService;
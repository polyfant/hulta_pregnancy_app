type ModelType = 'growth' | 'health' | 'nutrition';

export const mlService = {
    private models: Map<ModelType, any> = new Map(),
    private tf: typeof import('@tensorflow/tfjs') | null = null,

    async initialize(type: ModelType) {
        if (!this.tf) {
            // Dynamic import - only loads when needed
            this.tf = await import('@tensorflow/tfjs/dist/tf.min.js');
        }
        
        if (!this.models.has(type)) {
            // Load specific model on demand
            const modelPath = `/models/${type}-model.json`;
            const model = await this.tf.loadLayersModel(modelPath);
            this.models.set(type, model);
        }
    },

    async predictGrowth(data: number[][]): Promise<number[]> {
        await this.initialize('growth');
        const model = this.models.get('growth');
        
        // Convert to tensor, predict, and cleanup
        const tensor = this.tf.tensor2d(data);
        const prediction = model.predict(tensor);
        const result = await prediction.array();
        
        // Cleanup to prevent memory leaks
        tensor.dispose();
        prediction.dispose();
        
        return result;
    }
}; 
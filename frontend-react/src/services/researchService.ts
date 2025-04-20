import { EQUINE_DATABASES, ResearchData } from '../types/research';

export class ResearchService {
    async fetchResearchData(source: keyof typeof EQUINE_DATABASES): Promise<ResearchData[]> {
        const url = EQUINE_DATABASES[source];
        const response = await fetch(url);
        return response.json();
    }
}

export const researchService = new ResearchService(); 
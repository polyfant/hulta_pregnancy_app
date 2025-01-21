import { Horse, CreateHorseInput } from '../types/horse';
import { PregnancyStatus, PregnancyEvent, PregnancyChecklist } from '../types/pregnancy';
import { useApiClient } from './client';

export const useHorsesApi = () => {
    const api = useApiClient();

    return {
        // Existing horse endpoints
        getAll: async (): Promise<Horse[]> => {
            return api.get<Horse[]>('/horses');
        },

        getById: async (id: number): Promise<Horse> => {
            return api.get<Horse>(`/horses/${id}`);
        },

        create: async (horse: CreateHorseInput): Promise<Horse> => {
            return api.post<Horse>('/horses', horse);
        },

        delete: async (id: number): Promise<void> => {
            return api.delete(`/horses/${id}`);
        },

        update: async (id: number, horse: CreateHorseInput): Promise<Horse> => {
            return api.put<Horse>(`/horses/${id}`, horse);
        },

        // New pregnancy endpoints
        getPregnancyStatus: async (horseId: number): Promise<PregnancyStatus> => {
            return api.get<PregnancyStatus>(`/horses/${horseId}/pregnancy`);
        },

        getPregnancyEvents: async (horseId: number): Promise<PregnancyEvent[]> => {
            return api.get<PregnancyEvent[]>(`/horses/${horseId}/pregnancy/events`);
        },

        getPregnancyChecklist: async (horseId: number): Promise<PregnancyChecklist> => {
            return api.get<PregnancyChecklist>(`/horses/${horseId}/pregnancy/checklist`);
        },

        updateChecklistItem: async (
            horseId: number, 
            itemId: string, 
            completed: boolean
        ): Promise<PregnancyChecklist> => {
            return api.patch<PregnancyChecklist>(
                `/horses/${horseId}/pregnancy/checklist/${itemId}`,
                { completed }
            );
        },

        addChecklistItem: async (
            horseId: number,
            text: string
        ): Promise<PregnancyChecklist> => {
            return api.post<PregnancyChecklist>(
                `/horses/${horseId}/pregnancy/checklist`,
                { text }
            );
        },

        addPregnancyEvent: async (
            horseId: number,
            event: Omit<PregnancyEvent, 'id' | 'horseId'>
        ): Promise<PregnancyEvent> => {
            return api.post<PregnancyEvent>(
                `/horses/${horseId}/pregnancy/events`,
                event
            );
        },
    };
};

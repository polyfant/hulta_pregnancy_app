import { Horse, CreateHorseInput } from '../types/horse';
import { useApiClient } from './client';

export const useHorsesApi = () => {
    const api = useApiClient();

    return {
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
    };
};

import axios from 'axios';
import { Horse, CreateHorseInput } from '../types/horse';

const API_URL = 'http://localhost:8080/api';

export const horsesApi = {
    getAll: async (): Promise<Horse[]> => {
        const response = await axios.get(`${API_URL}/horses`);
        return response.data;
    },

    getById: async (id: number): Promise<Horse> => {
        const response = await axios.get(`${API_URL}/horses/${id}`);
        return response.data;
    },

    create: async (horse: CreateHorseInput): Promise<Horse> => {
        const response = await axios.post(`${API_URL}/horses`, horse);
        return response.data;
    },

    delete: async (id: number): Promise<void> => {
        await axios.delete(`${API_URL}/horses/${id}`);
    },

    update: async (id: number, horse: CreateHorseInput): Promise<Horse> => {
        const response = await axios.put(`${API_URL}/horses/${id}`, horse);
        return response.data;
    },
};

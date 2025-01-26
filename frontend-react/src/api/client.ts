import axios from 'axios';
import { useAuth0 } from '@auth0/auth0-react';

const API_URL = `${import.meta.env.VITE_API_URL || 'http://localhost:8080'}/api/v1`;

export const createApiClient = () => {
    const { getAccessTokenSilently } = useAuth0();

    const client = axios.create({
        baseURL: API_URL,
    });

    client.interceptors.request.use(async (config) => {
        try {
            const token = await getAccessTokenSilently();
            if (token) {
                config.headers.Authorization = `Bearer ${token}`;
            }
        } catch (error) {
            console.error('Error getting access token:', error);
        }
        return config;
    });

    return {
        get: async <T>(url: string): Promise<T> => {
            const response = await client.get<T>(url);
            return response.data;
        },
        post: async <T>(url: string, data: unknown): Promise<T> => {
            const response = await client.post<T>(url, data);
            return response.data;
        },
        put: async <T>(url: string, data: unknown): Promise<T> => {
            const response = await client.put<T>(url, data);
            return response.data;
        },
        patch: async <T>(url: string, data: unknown): Promise<T> => {
            const response = await client.patch<T>(url, data);
            return response.data;
        },
        delete: async (url: string): Promise<void> => {
            await client.delete(url);
        },
    };
};

export const useApiClient = () => {
    return createApiClient();
};

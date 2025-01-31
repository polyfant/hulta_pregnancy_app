import axios, { AxiosInstance } from 'axios';
import { useAuth0 } from '@auth0/auth0-react';

// Type for environment variables
declare global {
    interface ImportMetaEnv {
        VITE_API_URL: string;
    }
}

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

export interface ApiClient {
    get: <T>(url: string) => Promise<T>;
    post: <T>(url: string, data: unknown) => Promise<T>;
    put: <T>(url: string, data: unknown) => Promise<T>;
    patch: <T>(url: string, data: unknown) => Promise<T>;
    delete: (url: string) => Promise<void>;
}

export const useApiClient = (): ApiClient => {
    const { getAccessTokenSilently } = useAuth0();

    const client: AxiosInstance = axios.create({
        baseURL: '/api',
        headers: {
            'Content-Type': 'application/json',
        },
        withCredentials: true,
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
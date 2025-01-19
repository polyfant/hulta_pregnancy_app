import axios from 'axios';
import { useAuth0 } from '@auth0/auth0-react';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api';

export const createApiClient = () => {
    const { getAccessTokenSilently } = useAuth0();

    const client = axios.create({
        baseURL: API_URL,
    });

    client.interceptors.request.use(async (config) => {
        try {
            const token = await getAccessTokenSilently();
            config.headers.Authorization = `Bearer ${token}`;
        } catch (error) {
            console.error('Error getting access token:', error);
        }
        return config;
    });

    return client;
};

export const useApiClient = () => {
    const client = createApiClient();
    
    return {
        get: <T>(url: string) => client.get<T>(url).then(res => res.data),
        post: <T>(url: string, data: any) => client.post<T>(url, data).then(res => res.data),
        put: <T>(url: string, data: any) => client.put<T>(url, data).then(res => res.data),
        delete: (url: string) => client.delete(url).then(res => res.data),
    };
};

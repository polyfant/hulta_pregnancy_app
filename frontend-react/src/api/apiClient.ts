import { useAuth } from '../auth/useAuth';

const API_URL = '/api/v1';

interface FetchOptions extends RequestInit {
	params?: Record<string, string>;
}

export const createApiClient = (getToken: () => Promise<string | null>) => {
	const fetchWithAuth = async <T>(
		endpoint: string,
		options: FetchOptions = {}
	): Promise<T> => {
		const { params, ...fetchOptions } = options;

		// Add query parameters if provided
		let url = `${API_URL}${endpoint}`;
		if (params) {
			const queryParams = new URLSearchParams();
			Object.entries(params).forEach(([key, value]) => {
				queryParams.append(key, value);
			});
			url += `?${queryParams.toString()}`;
		}

		// Get the auth token
		const token = await getToken();

		// Prepare headers
		const headers = new Headers(options.headers);

		// Add content type if not present and we're sending a body
		if (
			options.body &&
			!headers.has('Content-Type') &&
			!(options.body instanceof FormData)
		) {
			headers.append('Content-Type', 'application/json');
		}

		// Add authorization header if we have a token
		if (token) {
			headers.append('Authorization', `Bearer ${token}`);
		}

		const response = await fetch(url, {
			...fetchOptions,
			headers,
		});

		// Handle errors
		if (!response.ok) {
			const errorData = await response.json().catch(() => ({
				message: response.statusText,
			}));

			throw new Error(
				errorData.message ||
					`API error: ${response.status} ${response.statusText}`
			);
		}

		// Parse response based on content type
		const contentType = response.headers.get('content-type');
		if (contentType?.includes('application/json')) {
			return response.json();
		} else if (contentType?.includes('text/')) {
			return (await response.text()) as unknown as T;
		} else {
			return response as unknown as T;
		}
	};

	return {
		get: <T>(endpoint: string, options?: FetchOptions) =>
			fetchWithAuth<T>(endpoint, { ...options, method: 'GET' }),

		post: <T>(endpoint: string, data: any, options?: FetchOptions) =>
			fetchWithAuth<T>(endpoint, {
				...options,
				method: 'POST',
				body: JSON.stringify(data),
			}),

		put: <T>(endpoint: string, data: any, options?: FetchOptions) =>
			fetchWithAuth<T>(endpoint, {
				...options,
				method: 'PUT',
				body: JSON.stringify(data),
			}),

		patch: <T>(endpoint: string, data: any, options?: FetchOptions) =>
			fetchWithAuth<T>(endpoint, {
				...options,
				method: 'PATCH',
				body: JSON.stringify(data),
			}),

		delete: <T>(endpoint: string, options?: FetchOptions) =>
			fetchWithAuth<T>(endpoint, { ...options, method: 'DELETE' }),

		upload: <T>(
			endpoint: string,
			formData: FormData,
			options?: FetchOptions
		) =>
			fetchWithAuth<T>(endpoint, {
				...options,
				method: 'POST',
				body: formData,
			}),
	};
};

// Hook to use the API client with current auth token
export const useApiClient = () => {
	const { getToken } = useAuth();
	return createApiClient(getToken);
};

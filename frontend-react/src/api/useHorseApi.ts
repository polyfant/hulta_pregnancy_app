import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';

const API_URL = '/api/v1';

// Type definitions for API requests
export interface CreateHorseRequest {
	name: string;
	breed?: string;
	gender: 'Male' | 'Female' | 'Gelding';
	birthDate?: string;
	color?: string;
	height?: number;
	weight?: number;
	registrationNumber?: string;
	microchipNumber?: string;
	status: 'active' | 'retired' | 'sold' | 'deceased';
	isPregnant?: boolean;
	lastBreedingDate?: string;
	healthStatus?: 'excellent' | 'good' | 'fair' | 'poor';
	notes?: string;
	sire?: string;
	dam?: string;
	isPremium?: boolean;
	isChampion?: boolean;
}

export type UpdateHorseRequest = Partial<CreateHorseRequest> & { id: string };

// Error handling helper
const handleApiError = async (response: Response) => {
	if (!response.ok) {
		const errorData = await response.json().catch(() => ({}));
		throw new Error(
			errorData.message ||
				`API error: ${response.status} ${response.statusText}`
		);
	}
	return response.json();
};

// Hook for fetching all horses
export const useHorses = () => {
	return useQuery({
		queryKey: ['horses'],
		queryFn: async () => {
			const response = await fetch(`${API_URL}/horses`);
			return handleApiError(response);
		},
	});
};

// Hook for fetching a single horse by ID
export const useHorse = (id: string) => {
	return useQuery({
		queryKey: ['horse', id],
		queryFn: async () => {
			const response = await fetch(`${API_URL}/horses/${id}`);
			return handleApiError(response);
		},
		enabled: !!id, // Only run query if ID is provided
	});
};

// Hook for creating a new horse
export const useCreateHorse = () => {
	const queryClient = useQueryClient();

	return useMutation({
		mutationFn: async (horse: CreateHorseRequest) => {
			const response = await fetch(`${API_URL}/horses`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
				},
				body: JSON.stringify(horse),
			});
			return handleApiError(response);
		},
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ['horses'] });
		},
	});
};

// Hook for updating an existing horse
export const useUpdateHorse = () => {
	const queryClient = useQueryClient();

	return useMutation({
		mutationFn: async (horse: UpdateHorseRequest) => {
			const { id, ...horseData } = horse;
			const response = await fetch(`${API_URL}/horses/${id}`, {
				method: 'PUT',
				headers: {
					'Content-Type': 'application/json',
				},
				body: JSON.stringify(horseData),
			});
			return handleApiError(response);
		},
		onSuccess: (_, variables) => {
			queryClient.invalidateQueries({ queryKey: ['horses'] });
			queryClient.invalidateQueries({
				queryKey: ['horse', variables.id],
			});
		},
	});
};

// Hook for deleting a horse
export const useDeleteHorse = () => {
	const queryClient = useQueryClient();

	return useMutation({
		mutationFn: async (id: string) => {
			const response = await fetch(`${API_URL}/horses/${id}`, {
				method: 'DELETE',
			});
			return handleApiError(response);
		},
		onSuccess: (_, id) => {
			queryClient.invalidateQueries({ queryKey: ['horses'] });
			queryClient.invalidateQueries({ queryKey: ['horse', id] });
		},
	});
};

// Hook for pregnancy-related operations
export const useHorsePregnancy = () => {
	const queryClient = useQueryClient();

	return {
		markAsPregnant: useMutation({
			mutationFn: async ({
				horseId,
				breedingDate,
			}: {
				horseId: string;
				breedingDate: string;
			}) => {
				const response = await fetch(
					`${API_URL}/horses/${horseId}/pregnancy`,
					{
						method: 'POST',
						headers: {
							'Content-Type': 'application/json',
						},
						body: JSON.stringify({ breedingDate }),
					}
				);
				return handleApiError(response);
			},
			onSuccess: (_, variables) => {
				queryClient.invalidateQueries({
					queryKey: ['horse', variables.horseId],
				});
				queryClient.invalidateQueries({ queryKey: ['horses'] });
			},
		}),

		updatePregnancyStatus: useMutation({
			mutationFn: async ({
				horseId,
				stage,
				notes,
			}: {
				horseId: string;
				stage: string;
				notes?: string;
			}) => {
				const response = await fetch(
					`${API_URL}/horses/${horseId}/pregnancy/status`,
					{
						method: 'PUT',
						headers: {
							'Content-Type': 'application/json',
						},
						body: JSON.stringify({ stage, notes }),
					}
				);
				return handleApiError(response);
			},
			onSuccess: (_, variables) => {
				queryClient.invalidateQueries({
					queryKey: ['horse', variables.horseId],
				});
			},
		}),

		endPregnancy: useMutation({
			mutationFn: async ({
				horseId,
				outcome,
				foalingDate,
				notes,
			}: {
				horseId: string;
				outcome: 'successful' | 'unsuccessful';
				foalingDate: string;
				notes?: string;
			}) => {
				const response = await fetch(
					`${API_URL}/horses/${horseId}/pregnancy/end`,
					{
						method: 'PUT',
						headers: {
							'Content-Type': 'application/json',
						},
						body: JSON.stringify({ outcome, foalingDate, notes }),
					}
				);
				return handleApiError(response);
			},
			onSuccess: (_, variables) => {
				queryClient.invalidateQueries({
					queryKey: ['horse', variables.horseId],
				});
				queryClient.invalidateQueries({ queryKey: ['horses'] });
			},
		}),
	};
};

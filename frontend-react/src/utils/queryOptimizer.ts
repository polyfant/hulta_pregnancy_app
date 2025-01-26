import { QueryClient, QueryClientProvider, useQuery, UseQueryOptions } from '@tanstack/react-query';
import { useCallback } from 'react';

// Create a custom query client with optimized settings
export const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      // Stale-while-revalidate strategy
      staleTime: 1000 * 60 * 5, // 5 minutes
      
      // Keep previous data while refetching
      keepPreviousData: true,
      
      // Retry failed queries 3 times with increasing intervals
      retry: (failureCount, error) => {
        if ((error as any)?.status === 404) return false;
        return failureCount < 3;
      },
      retryDelay: (attemptIndex) => Math.min(1000 * 2 ** attemptIndex, 30000),
      
      // Aggressive background refetching
      refetchOnWindowFocus: true,
      refetchInterval: 1000 * 60 * 10, // 10 minutes
    },
  },
});

// Advanced query hook with additional optimization features
export function useOptimizedQuery<TQueryFnData, TError, TData = TQueryFnData>(
  queryKey: string[],
  queryFn: () => Promise<TQueryFnData>,
  options?: Omit<UseQueryOptions<TQueryFnData, TError, TData>, 'queryKey' | 'queryFn'>
) {
  return useQuery<TQueryFnData, TError, TData>(
    queryKey,
    queryFn,
    {
      // Merge custom options with defaults
      ...options,
      
      // Optional: Add custom error handling
      onError: (error) => {
        console.error('Query failed:', error);
        options?.onError?.(error);
      },
      
      // Optional: Add custom success handling
      onSuccess: (data) => {
        // Potential cache warming or related data prefetching
        options?.onSuccess?.(data);
      },
    }
  );
}

// Smart cache invalidation utility
export const invalidateQueries = {
  // Intelligently invalidate related queries
  byPrefix: (prefix: string) => {
    queryClient.invalidateQueries({
      predicate: (query) => 
        query.queryKey[0].toString().startsWith(prefix)
    });
  },
  
  // Selective cache clearing
  selective: (queryKeys: string[][]) => {
    queryKeys.forEach(key => queryClient.invalidateQueries({ queryKey: key }));
  }
};

// Prefetching utility for performance optimization
export const prefetchData = {
  // Prefetch data for anticipated user actions
  anticipate: (queryKey: string[], queryFn: () => Promise<any>) => {
    queryClient.prefetchQuery(queryKey, queryFn, {
      staleTime: 1000 * 60 * 5, // 5 minutes
    });
  }
};

// Example of how to use these utilities
export function ExampleComponent() {
  // Optimized query with advanced features
  const horseQuery = useOptimizedQuery(
    ['horses', 'list'], 
    fetchHorses,
    {
      // Additional per-query customization
      staleTime: 1000 * 60 * 10, // 10 minutes for this specific query
    }
  );

  // Prefetch related data on hover or anticipated user action
  const prefetchHorseDetails = useCallback((horseId: string) => {
    prefetchData.anticipate(
      ['horses', 'details', horseId], 
      () => fetchHorseDetails(horseId)
    );
  }, []);

  return null; // Placeholder
}

// Placeholder functions for demonstration
function fetchHorses() { return Promise.resolve([]); }
function fetchHorseDetails(id: string) { return Promise.resolve({}); }

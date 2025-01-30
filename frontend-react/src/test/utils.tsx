import { render } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { MemoryRouter, Routes, Route } from 'react-router-dom';
import { MantineProvider } from '@mantine/core';
import { Notifications } from '@mantine/notifications';

const queryClient = new QueryClient({
    defaultOptions: {
        queries: {
            retry: false,
        },
    },
});

// Create a custom render function
export function renderWithProviders(ui: React.ReactElement, { initialRoute = '/' } = {}) {
    return render(
        <QueryClientProvider client={queryClient}>
            <MantineProvider>
                <Notifications />
                <MemoryRouter initialEntries={[initialRoute]}>
                    <Routes>
                        <Route path="*" element={ui} />
                    </Routes>
                </MemoryRouter>
            </MantineProvider>
        </QueryClientProvider>
    );
}

// Re-export everything
export * from '@testing-library/react';

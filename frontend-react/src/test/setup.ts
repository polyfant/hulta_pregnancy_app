import '@testing-library/jest-dom';
import { expect, afterEach, vi } from 'vitest';
import { cleanup } from '@testing-library/react';
import * as matchers from '@testing-library/jest-dom/matchers';
import { http, HttpResponse } from 'msw';
import { setupServer } from 'msw/node';

// Add this before other imports
Object.defineProperty(window, 'matchMedia', {
    writable: true,
    value: vi.fn().mockImplementation(query => ({
        matches: false,
        media: query,
        onchange: null,
        addListener: vi.fn(),
        removeListener: vi.fn(),
        addEventListener: vi.fn(),
        removeEventListener: vi.fn(),
        dispatchEvent: vi.fn(),
    })),
});

// Mock Auth0
vi.mock('@auth0/auth0-react', () => ({
    useAuth0: () => ({
        isAuthenticated: true,
        user: { sub: 'user123', email: 'test@example.com' },
        getAccessTokenSilently: () => Promise.resolve('fake_token'),
    }),
}));

// Mock environment variables
vi.mock('../env', () => ({
    VITE_AUTH0_DOMAIN: 'test.auth0.com',
    VITE_AUTH0_CLIENT_ID: 'test_client_id',
    VITE_AUTH0_AUDIENCE: 'https://api.hulta-foaltracker.app',
    VITE_API_URL: 'http://localhost:8080',
}));

// Mock ResizeObserver
global.ResizeObserver = class ResizeObserver {
    observe() {}
    unobserve() {}
    disconnect() {}
};

expect.extend(matchers);

// Mock server setup
export const handlers = [
  http.get('/api/horses', () => {
    return HttpResponse.json([
      {
        id: '1',
        name: 'Thunder',
        breed: 'Arabian',
        gender: 'STALLION',
        dateOfBirth: '2020-01-01',
      },
      {
        id: '2',
        name: 'Storm',
        breed: 'Thoroughbred',
        gender: 'MARE',
        dateOfBirth: '2019-05-15',
      },
    ]);
  }),
  http.post('http://localhost:8080/api/horses', async ({ request }) => {
    const newHorse = await request.json();
    return HttpResponse.json({ id: 3, ...newHorse });
  }),
  http.put('http://localhost:8080/api/horses/:id', async ({ request }) => {
    const updatedHorse = await request.json();
    return HttpResponse.json(updatedHorse);
  }),
  http.delete('http://localhost:8080/api/horses/:id', () => {
    return new HttpResponse(null, { status: 204 });
  }),
];

export const server = setupServer(...handlers);

// Start server before all tests
beforeAll(() => server.listen({ onUnhandledRequest: 'bypass' }));

// Reset handlers after each test
afterEach(() => {
  server.resetHandlers();
  cleanup();
});

// Clean up after all tests are done
afterAll(() => server.close());

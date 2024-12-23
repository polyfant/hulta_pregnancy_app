import '@testing-library/jest-dom';
import { expect, afterEach } from 'vitest';
import { cleanup } from '@testing-library/react';
import * as matchers from '@testing-library/jest-dom/matchers';
import { http, HttpResponse } from 'msw';
import { setupServer } from 'msw/node';

expect.extend(matchers);

// Mock server setup
export const handlers = [
  http.get('http://localhost:8080/api/horses', () => {
    return HttpResponse.json([
      {
        id: 1,
        name: 'Thunder',
        breed: 'Arabian',
        gender: 'STALLION',
        dateOfBirth: '2020-01-01',
      },
      {
        id: 2,
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
beforeAll(() => server.listen({ onUnhandledRequest: 'error' }));

// Reset handlers after each test
afterEach(() => {
  server.resetHandlers();
  cleanup();
});

// Clean up after all tests are done
afterAll(() => server.close());

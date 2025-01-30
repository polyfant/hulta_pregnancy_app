import { screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { HorseList } from '../HorseList';
import { server } from '../../test/setup';
import { http, HttpResponse } from 'msw';
import { renderWithProviders } from '../../test/utils';

describe('HorseList', () => {
    it('shows loading skeleton initially', () => {
        renderWithProviders(<HorseList />);
        expect(screen.getAllByTestId('horse-card-skeleton')).toHaveLength(6);
    });

    it('handles network errors gracefully', async () => {
        server.use(
            http.get('/api/horses', () => {
                return new HttpResponse(null, { status: 500 });
            })
        );

        renderWithProviders(<HorseList />);
        expect(await screen.findByText(/failed to load horses/i)).toBeInTheDocument();
        expect(screen.getByRole('button', { name: /try again/i })).toBeInTheDocument();
    });

    it('shows pregnancy status for pregnant mares', async () => {
        server.use(
            http.get('/api/horses', () => {
                return HttpResponse.json([{
                    id: '1',
                    name: 'Luna',
                    gender: 'MARE',
                    isPregnant: true
                }]);
            }),
            http.get('/api/horses/1/pregnancy', () => {
                return HttpResponse.json({
                    currentStage: 'MIDDLE',
                    daysRemaining: 180
                });
            })
        );

        renderWithProviders(<HorseList />);
        expect(await screen.findByText(/180 days to foaling/i)).toBeInTheDocument();
        expect(screen.getByRole('progressbar')).toBeInTheDocument();
    });

    it('shows empty state when no horses', async () => {
        renderWithProviders(<HorseList />);
        expect(await screen.findByText('No Horses Found')).toBeInTheDocument();
    });

    it('displays horses when data is loaded', async () => {
        const mockHorses = [
            {
                id: '1',
                name: 'Thunder',
                breed: 'Arabian',
                gender: 'STALLION',
            },
            {
                id: '2',
                name: 'Storm',
                breed: 'Thoroughbred',
                gender: 'MARE',
                isPregnant: true,
            },
        ];

        // Mock the API call
        vi.spyOn(global, 'fetch').mockImplementation(() =>
            Promise.resolve({
                ok: true,
                json: () => Promise.resolve(mockHorses),
            } as Response)
        );

        renderWithProviders(<HorseList />);

        // Check if horses are displayed
        expect(await screen.findByText('Thunder')).toBeInTheDocument();
        expect(await screen.findByText('Storm')).toBeInTheDocument();
    });

    it('filters horses by search query', async () => {
        server.use(
            http.get('/api/horses', () => {
                return HttpResponse.json([
                    { id: '1', name: 'Thunder', breed: 'Arabian' },
                    { id: '2', name: 'Storm', breed: 'Thoroughbred' }
                ]);
            })
        );

        renderWithProviders(<HorseList />);
        
        const searchInput = screen.getByPlaceholderText(/search horses/i);
        await userEvent.type(searchInput, 'thunder');
        
        expect(screen.getByText('Thunder')).toBeInTheDocument();
        expect(screen.queryByText('Storm')).not.toBeInTheDocument();
    });

    it('shows different pregnancy stages with correct colors', async () => {
        server.use(
            http.get('/api/horses', () => {
                return HttpResponse.json([{
                    id: '1',
                    name: 'Luna',
                    gender: 'MARE',
                    isPregnant: true
                }]);
            }),
            http.get('/api/horses/1/pregnancy', () => {
                return HttpResponse.json({
                    currentStage: 'LATE',
                    daysRemaining: 90,
                    progress: 75
                });
            })
        );

        renderWithProviders(<HorseList />);
        
        const badge = await screen.findByText(/90 days to foaling/i);
        expect(badge).toHaveStyle({ backgroundColor: expect.stringContaining('teal') });
        
        const progressBar = screen.getByRole('progressbar');
        expect(progressBar).toHaveAttribute('aria-valuenow', '75');
    });

    it('handles pregnancy data loading errors', async () => {
        server.use(
            http.get('/api/horses', () => {
                return HttpResponse.json([{
                    id: '1',
                    name: 'Luna',
                    gender: 'MARE',
                    isPregnant: true
                }]);
            }),
            http.get('/api/horses/1/pregnancy', () => {
                return new HttpResponse(null, { status: 500 });
            })
        );

        renderWithProviders(<HorseList />);
        expect(await screen.findByText(/failed to fetch pregnancy status/i)).toBeInTheDocument();
    });
}); 
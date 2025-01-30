import { DeleteHorseButton } from '../DeleteHorseButton';
import { screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { HttpResponse } from 'msw';
import { server } from '../../test/setup';
import { http } from 'msw';
import { renderWithProviders } from '../../test/utils';

describe('DeleteHorseButton', () => {
    it('shows confirmation dialog before deleting', async () => {
        renderWithProviders(<DeleteHorseButton horseId="1" horseName="Thunder" />);
        
        await userEvent.click(screen.getByRole('button', { name: /delete horse/i }));
        expect(screen.getByText(/are you sure/i)).toBeInTheDocument();
    });

    it('shows success notification after deletion', async () => {
        renderWithProviders(<DeleteHorseButton horseId="1" horseName="Thunder" />);
        
        await userEvent.click(screen.getByRole('button', { name: /delete horse/i }));
        await userEvent.click(screen.getByRole('button', { name: /yes, delete/i }));
        
        expect(await screen.findByText(/thunder has been deleted/i)).toBeInTheDocument();
    });

    it('handles deletion errors', async () => {
        server.use(
            http.delete('/api/horses/1', () => {
                return new HttpResponse(null, { status: 500 });
            })
        );

        renderWithProviders(<DeleteHorseButton horseId="1" horseName="Thunder" />);
        
        await userEvent.click(screen.getByRole('button', { name: /delete horse/i }));
        await userEvent.click(screen.getByRole('button', { name: /yes, delete/i }));
        
        expect(await screen.findByText(/failed to delete horse/i)).toBeInTheDocument();
    });

    it('closes modal on cancel', async () => {
        renderWithProviders(<DeleteHorseButton horseId="1" horseName="Thunder" />);
        
        await userEvent.click(screen.getByRole('button', { name: /delete horse/i }));
        await userEvent.click(screen.getByRole('button', { name: /cancel/i }));
        
        expect(screen.queryByText(/are you sure/i)).not.toBeInTheDocument();
    });
}); 
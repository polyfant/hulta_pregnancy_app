import { screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { HorseForm } from '../HorseForm';
import dayjs from 'dayjs';
import { HttpResponse } from 'msw';
import { server } from '../../test/setup';
import { http } from 'msw';
import { renderWithProviders } from '../../test/utils';

describe('HorseForm', () => {
    it('validates required fields', async () => {
        const onSubmit = vi.fn();
        renderWithProviders(<HorseForm onSubmit={onSubmit} />);

        // Try to submit empty form
        await userEvent.click(screen.getByRole('button', { name: /add horse/i }));

        // Check for validation messages
        expect(await screen.findByText(/name is required/i)).toBeInTheDocument();
        expect(await screen.findByText(/breed is required/i)).toBeInTheDocument();
        expect(onSubmit).not.toHaveBeenCalled();
    });

    it('submits form with valid data', async () => {
        const onSubmit = vi.fn();
        renderWithProviders(<HorseForm onSubmit={onSubmit} />);

        // Fill form
        await userEvent.type(screen.getByLabelText(/horse name/i), 'Thunder');
        await userEvent.type(screen.getByLabelText(/breed/i), 'Arabian');
        await userEvent.selectOptions(screen.getByLabelText(/gender/i), ['STALLION']);

        // Submit form
        await userEvent.click(screen.getByRole('button', { name: /add horse/i }));

        await waitFor(() => {
            expect(onSubmit).toHaveBeenCalledWith({
                name: 'Thunder',
                breed: 'Arabian',
                gender: 'STALLION',
                // ... other fields
            });
        });
    });

    it('shows pregnancy fields only for mares', async () => {
        renderWithProviders(<HorseForm />);
        
        // Select mare
        await userEvent.selectOptions(screen.getByLabelText(/gender/i), ['MARE']);
        expect(screen.getByText(/pregnancy information/i)).toBeInTheDocument();
        
        // Change to stallion
        await userEvent.selectOptions(screen.getByLabelText(/gender/i), ['STALLION']);
        expect(screen.queryByText(/pregnancy information/i)).not.toBeInTheDocument();
    });

    it('shows success notification on successful submission', async () => {
        renderWithProviders(<HorseForm />);
        
        // Fill and submit form
        await userEvent.type(screen.getByLabelText(/name/i), 'Thunder');
        await userEvent.type(screen.getByLabelText(/breed/i), 'Arabian');
        await userEvent.click(screen.getByRole('button', { name: /add horse/i }));
        
        expect(await screen.findByText(/horse added successfully/i)).toBeInTheDocument();
    });

    it('shows loading progress during submission', async () => {
        renderWithProviders(<HorseForm />);
        
        await userEvent.type(screen.getByLabelText(/name/i), 'Thunder');
        await userEvent.click(screen.getByRole('button', { name: /add horse/i }));
        
        expect(screen.getByText(/saving horse details/i)).toBeInTheDocument();
        expect(screen.getByRole('progressbar')).toBeInTheDocument();
    });

    it('validates date of birth is not in future', async () => {
        renderWithProviders(<HorseForm />);
        
        const futureDate = dayjs().add(1, 'year').format('YYYY-MM-DD');
        await userEvent.type(screen.getByLabelText(/date of birth/i), futureDate);
        
        expect(screen.getByText(/date cannot be in the future/i)).toBeInTheDocument();
    });

    it('handles parent selection validation', async () => {
        server.use(
            http.get('/api/horses', () => {
                return HttpResponse.json([
                    { id: '1', name: 'Mother', gender: 'MARE' },
                    { id: '2', name: 'Father', gender: 'STALLION' }
                ]);
            })
        );

        renderWithProviders(<HorseForm />);
        
        // Test mother selection
        const motherSelect = screen.getByPlaceholderText(/select mother/i);
        await userEvent.click(motherSelect);
        await userEvent.click(screen.getByText('Mother'));
        
        // Test father selection
        const fatherSelect = screen.getByPlaceholderText(/select father/i);
        await userEvent.click(fatherSelect);
        await userEvent.click(screen.getByText('Father'));
        
        expect(screen.queryByText(/invalid parent selection/i)).not.toBeInTheDocument();
    });

    it('shows conception date field only when mare is pregnant', async () => {
        renderWithProviders(<HorseForm />);
        
        // Select mare
        await userEvent.selectOptions(screen.getByLabelText(/gender/i), ['MARE']);
        
        // Toggle pregnancy
        const pregnantCheckbox = screen.getByLabelText(/mare is pregnant/i);
        await userEvent.click(pregnantCheckbox);
        
        expect(screen.getByLabelText(/conception date/i)).toBeInTheDocument();
        
        // Toggle off pregnancy
        await userEvent.click(pregnantCheckbox);
        expect(screen.queryByLabelText(/conception date/i)).not.toBeInTheDocument();
    });
}); 
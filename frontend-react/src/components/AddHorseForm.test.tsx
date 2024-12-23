import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent, waitFor } from '../test/utils';
import { AddHorseForm } from './AddHorseForm';
import userEvent from '@testing-library/user-event';

describe('AddHorseForm', () => {
  it('renders all form fields', () => {
    render(<AddHorseForm />);
    
    expect(screen.getByLabelText(/name/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/breed/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/gender/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/date of birth/i)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /add horse/i })).toBeInTheDocument();
  });

  it('shows validation errors for required fields', async () => {
    render(<AddHorseForm />);
    
    const submitButton = screen.getByRole('button', { name: /add horse/i });
    fireEvent.click(submitButton);

    await waitFor(() => {
      expect(screen.getByText(/name must be at least 2 characters/i)).toBeInTheDocument();
      expect(screen.getByText(/breed must be at least 2 characters/i)).toBeInTheDocument();
    });
  });

  it('submits form with valid data', async () => {
    const onSubmit = vi.fn();
    render(<AddHorseForm onSubmit={onSubmit} />);
    
    await userEvent.type(screen.getByLabelText(/name/i), 'Thunder');
    await userEvent.type(screen.getByLabelText(/breed/i), 'Arabian');
    
    const submitButton = screen.getByRole('button', { name: /add horse/i });
    fireEvent.click(submitButton);

    await waitFor(() => {
      expect(onSubmit).toHaveBeenCalledWith(expect.objectContaining({
        name: 'Thunder',
        breed: 'Arabian',
        gender: 'MARE', // Default value
      }));
    });
  });

  it('pre-fills form when initialValues are provided', () => {
    const initialValues = {
      id: 1,
      name: 'Storm',
      breed: 'Thoroughbred',
      gender: 'STALLION',
      dateOfBirth: '2020-01-01',
    };
    
    render(<AddHorseForm initialValues={initialValues} submitButtonText="Update Horse" />);
    
    expect(screen.getByLabelText(/name/i)).toHaveValue('Storm');
    expect(screen.getByLabelText(/breed/i)).toHaveValue('Thoroughbred');
    expect(screen.getByRole('button', { name: /update horse/i })).toBeInTheDocument();
  });
});

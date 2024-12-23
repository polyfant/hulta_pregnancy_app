import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent, waitFor } from '../test/utils';
import { HorseList } from './HorseList';
import { server, handlers } from '../test/setup';
import { http, HttpResponse } from 'msw';

describe('HorseList', () => {
  it('renders loading state initially', () => {
    render(<HorseList />);
    expect(screen.getByText(/loading horses/i)).toBeInTheDocument();
  });

  it('renders horses after loading', async () => {
    render(<HorseList />);
    
    await waitFor(() => {
      expect(screen.getByText('Thunder')).toBeInTheDocument();
      expect(screen.getByText('Storm')).toBeInTheDocument();
    });
  });

  it('opens edit modal when edit button is clicked', async () => {
    render(<HorseList />);
    
    await waitFor(() => {
      expect(screen.getByText('Thunder')).toBeInTheDocument();
    });

    const editButtons = await screen.findAllByRole('button', { name: /edit/i });
    fireEvent.click(editButtons[0]);

    expect(screen.getByText(/edit horse/i)).toBeInTheDocument();
    expect(screen.getByDisplayValue('Thunder')).toBeInTheDocument();
  });

  it('deletes a horse when delete button is clicked', async () => {
    render(<HorseList />);
    
    await waitFor(() => {
      expect(screen.getByText('Thunder')).toBeInTheDocument();
    });

    // Mock the delete request
    server.use(
      http.delete('http://localhost:8080/api/horses/1', () => {
        return new HttpResponse(null, { status: 204 });
      })
    );

    const deleteButtons = await screen.findAllByRole('button', { name: /delete/i });
    fireEvent.click(deleteButtons[0]);

    // Wait for the delete request to complete and the list to update
    await waitFor(() => {
      expect(screen.queryByText('Thunder')).not.toBeInTheDocument();
    });
  });

  it('updates a horse when edit form is submitted', async () => {
    render(<HorseList />);
    
    await waitFor(() => {
      expect(screen.getByText('Thunder')).toBeInTheDocument();
    });

    // Mock the update request
    server.use(
      http.put('http://localhost:8080/api/horses/1', async ({ request }) => {
        const updatedHorse = await request.json();
        return HttpResponse.json({ ...updatedHorse, id: 1 });
      })
    );

    // Click edit button
    const editButtons = await screen.findAllByRole('button', { name: /edit/i });
    fireEvent.click(editButtons[0]);

    // Update the name
    const nameInput = screen.getByDisplayValue('Thunder');
    fireEvent.change(nameInput, { target: { value: 'Thunder Storm' } });

    // Submit the form
    const updateButton = screen.getByRole('button', { name: /update horse/i });
    fireEvent.click(updateButton);

    // Wait for the update to complete and the modal to close
    await waitFor(() => {
      expect(screen.queryByText(/edit horse/i)).not.toBeInTheDocument();
    });
  });

  it('handles error state', async () => {
    // Mock an error response
    server.use(
      http.get('http://localhost:8080/api/horses', () => {
        return new HttpResponse(null, { status: 500 });
      })
    );

    render(<HorseList />);
    
    await waitFor(() => {
      expect(screen.getByText(/error loading horses/i)).toBeInTheDocument();
    });
  });
});

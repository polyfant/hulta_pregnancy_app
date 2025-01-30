import { screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { EmptyState } from '../states/EmptyState';
import { renderWithProviders } from '../../test/utils';

describe('EmptyState', () => {
    it('renders custom message and action', () => {
        const onAction = vi.fn();
        renderWithProviders(
            <EmptyState 
                title="No Horses"
                message="Start by adding your first horse"
                actionLabel="Add Horse"
                onAction={onAction}
            />
        );
        
        expect(screen.getByText(/no horses/i)).toBeInTheDocument();
        expect(screen.getByText(/start by adding/i)).toBeInTheDocument();
        
        userEvent.click(screen.getByRole('button', { name: /add horse/i }));
        expect(onAction).toHaveBeenCalled();
    });
}); 
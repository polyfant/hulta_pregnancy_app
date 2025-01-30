import { screen } from '@testing-library/react';
import { LoadingProgress } from '../LoadingProgress';
import { renderWithProviders } from '../../test/utils';

describe('LoadingProgress', () => {
    it('renders with custom message', () => {
        renderWithProviders(<LoadingProgress message="Custom loading message" />);
        
        expect(screen.getByText('Custom loading message')).toBeInTheDocument();
        expect(screen.getByRole('progressbar')).toBeInTheDocument();
    });

    it('shows animated progress bar', () => {
        renderWithProviders(<LoadingProgress message="Loading..." />);
        
        const progressBar = screen.getByRole('progressbar');
        expect(progressBar).toHaveAttribute('data-animated', 'true');
    });
}); 
import { MantineProvider } from '@mantine/core';
import { render, screen } from '@testing-library/react';
import { EmptyState } from '../states/EmptyState';

describe('EmptyState', () => {
	it('renders basic empty state', () => {
		render(
			<MantineProvider>
				<EmptyState
					title='Test Title'
					message='Test Message'
					actionLabel='Test Action'
					onAction={() => {}}
				/>
			</MantineProvider>
		);

		expect(screen.getByText('Test Title')).toBeInTheDocument();
		expect(screen.getByText('Test Message')).toBeInTheDocument();
	});
});

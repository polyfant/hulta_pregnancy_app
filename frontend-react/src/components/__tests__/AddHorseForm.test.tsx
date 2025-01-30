import { render, screen } from '@testing-library/react';
import { TestWrapper } from '../../test/wrapper';
import AddHorseForm from '../AddHorseForm';

describe('AddHorseForm', () => {
	it('renders without crashing', () => {
		render(
			<TestWrapper>
				<AddHorseForm submitButtonText='Add Horse' />
			</TestWrapper>
		);

		expect(screen.getByText('Add Horse')).toBeInTheDocument();
		expect(screen.getByLabelText(/Name/i)).toBeInTheDocument();
		expect(screen.getByLabelText(/Gender/i)).toBeInTheDocument();
	});
});

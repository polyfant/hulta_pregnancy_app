import { render, screen, waitFor } from '@testing-library/react';
import { http, HttpResponse } from 'msw';
import { server } from '../../test/setup';
import { TestWrapper } from '../../test/wrapper';
import PregnancyTracking from '../PregnancyTracking/PregnancyTracking';

describe('PregnancyTracking', () => {
	it('renders without crashing', async () => {
		server.use(
			http.get('/api/horses/:id/pregnancy', () => {
				return HttpResponse.json({
					currentStage: 'EARLY',
					daysInPregnancy: 30,
					daysRemaining: 310,
					conceptionDate: '2024-01-01',
					expectedDueDate: '2024-12-01',
					progress: 10,
				});
			})
		);

		render(
			<TestWrapper>
				<PregnancyTracking />
			</TestWrapper>
		);

		await waitFor(() => {
			expect(screen.getByText(/Pregnancy Status/i)).toBeInTheDocument();
		});
	});
});

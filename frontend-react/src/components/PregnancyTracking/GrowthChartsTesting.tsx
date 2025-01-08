import { Select, Stack } from '@mantine/core';
import { useState } from 'react';
import { MOCK_BODY_CONDITION, MOCK_GROWTH_DATA } from '../../mocks/growthData';
import { GrowthCharts } from './GrowthCharts';

export function GrowthChartsTesting() {
	const [scenario, setScenario] = useState('normal');
	const [breed, setBreed] = useState('thoroughbred');

	// Mock the API responses
	const mockApi = {
		'/api/foals/1/growth-data': () => MOCK_GROWTH_DATA[breed][scenario],
		'/api/foals/1/condition': () => MOCK_BODY_CONDITION[scenario],
	};

	// Override fetch for testing
	const originalFetch = window.fetch;
	window.fetch = async (url) => {
		if (url in mockApi) {
			return new Promise((resolve) => {
				setTimeout(() => {
					resolve({
						ok: true,
						json: async () => mockApi[url](),
					});
				}, 500); // Simulate network delay
			});
		}
		return originalFetch(url);
	};

	return (
		<Stack>
			<Select
				label='Test Scenario'
				value={scenario}
				onChange={(value) => setScenario(value)}
				data={[
					{ value: 'normal', label: 'Normal Growth' },
					{ value: 'slow', label: 'Slow Growth' },
					{ value: 'rapid', label: 'Rapid Growth' },
				]}
			/>
			<GrowthCharts foalId={1} />
		</Stack>
	);
}

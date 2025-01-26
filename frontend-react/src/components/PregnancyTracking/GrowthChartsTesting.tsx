import { Select, Stack } from '@mantine/core';
import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { MOCK_GROWTH_DATA, MOCK_BODY_CONDITION } from '../../mocks/growthData';
import { GrowthCharts } from './GrowthCharts';

export function GrowthChartsTesting() {
    const [breed, setBreed] = useState<'thoroughbred' | 'warmblood' | 'arabian'>('thoroughbred');
    const [scenario, setScenario] = useState<'normal' | 'slow' | 'rapid'>('normal');

    // Simulate API call with mock data
    const { data: growthData, isLoading, error } = useQuery({
        queryKey: ['foal-growth', breed, scenario],
        queryFn: async () => {
            // Simulate network delay and API behavior
            await new Promise(resolve => setTimeout(resolve, 500));
            
            // Return mock data based on breed and scenario
            return MOCK_GROWTH_DATA[breed][scenario];
        }
    });

    const { data: bodyCondition } = useQuery({
        queryKey: ['body-condition', scenario],
        queryFn: async () => {
            await new Promise(resolve => setTimeout(resolve, 300));
            
            // Map scenarios to body condition
            const conditionMap = {
                normal: MOCK_BODY_CONDITION.normal,
                slow: MOCK_BODY_CONDITION.underweight,
                rapid: MOCK_BODY_CONDITION.overweight
            };
            
            return conditionMap[scenario];
        }
    });

    return (
        <Stack>
            <div style={{ display: 'flex', gap: '10px' }}>
                <Select
                    label="Breed"
                    value={breed}
                    onChange={(value) => setBreed(value as any)}
                    data={[
                        { value: 'thoroughbred', label: 'Thoroughbred' },
                        { value: 'warmblood', label: 'Warmblood' },
                        { value: 'arabian', label: 'Arabian' }
                    ]}
                />
                <Select
                    label="Growth Scenario"
                    value={scenario}
                    onChange={(value) => setScenario(value as any)}
                    data={[
                        { value: 'normal', label: 'Normal Growth' },
                        { value: 'slow', label: 'Slow Growth' },
                        { value: 'rapid', label: 'Rapid Growth' }
                    ]}
                />
            </div>

            {isLoading && <div>Loading growth data...</div>}
            {error && <div>Error fetching growth data</div>}

            {growthData && bodyCondition && (
                <>
                    <div>
                        <strong>Body Condition:</strong> {bodyCondition.notes}
                        <br />
                        Score: {bodyCondition.bodyConditionScore}
                        <br />
                        Muscle Score: {bodyCondition.muscleScore}
                        <br />
                        Fat Score: {bodyCondition.fatScore}
                    </div>
                    <GrowthCharts 
                        data={growthData} 
                        breed={breed} 
                        scenario={scenario} 
                    />
                </>
            )}
        </Stack>
    );
}

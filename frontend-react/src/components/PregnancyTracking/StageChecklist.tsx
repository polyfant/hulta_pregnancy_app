import { Checkbox, Paper, Stack, Text, Title } from '@mantine/core';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { PregnancyStage } from '../../types/pregnancy';

interface ChecklistItem {
	id: string;
	text: string;
	isCompleted: boolean;
	stage: PregnancyStage;
	isRequired: boolean;
}

interface StageChecklistProps {
	horseId: string;
	currentStage: PregnancyStage;
}

export function StageChecklist({ horseId, currentStage }: StageChecklistProps) {
	const queryClient = useQueryClient();

	const { data: checklist, isLoading } = useQuery({
		queryKey: ['pregnancy-checklist', horseId, currentStage],
		queryFn: async () => {
			const response = await fetch(
				`/api/horses/${horseId}/pregnancy/checklist?stage=${currentStage}`
			);
			if (!response.ok) throw new Error('Failed to fetch checklist');
			return response.json() as Promise<ChecklistItem[]>;
		},
	});

	const toggleMutation = useMutation({
		mutationFn: async (itemId: string) => {
			const response = await fetch(
				`/api/horses/${horseId}/pregnancy/checklist/${itemId}/toggle`,
				{
					method: 'POST',
				}
			);
			if (!response.ok)
				throw new Error('Failed to toggle checklist item');
		},
		onSuccess: () => {
			queryClient.invalidateQueries({
				queryKey: ['pregnancy-checklist', horseId],
			});
		},
	});

	if (isLoading) return <Text>Loading checklist...</Text>;

	return (
		<Paper p='md' radius='md' withBorder>
			<Stack>
				<Title order={3}>Stage Checklist</Title>
				{checklist?.map((item) => (
					<Checkbox
						key={item.id}
						label={item.text}
						checked={item.isCompleted}
						onChange={() => toggleMutation.mutate(item.id)}
						styles={(theme) => ({
							label: {
								color: item.isRequired
									? theme.colors.red[6]
									: undefined,
								fontWeight: item.isRequired ? 500 : undefined,
							},
						})}
					/>
				))}
			</Stack>
		</Paper>
	);
}

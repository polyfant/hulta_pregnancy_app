import { Checkbox, Paper, Stack, Title } from '@mantine/core';

interface ChecklistItem {
	id: string;
	text: string;
	completed: boolean;
	dueDate?: Date;
}

export function PregnancyChecklist() {
	const checklistItems: ChecklistItem[] = [
		{ id: '1', text: 'Prepare foaling kit', completed: false },
		{ id: '2', text: 'Clean and prepare foaling stall', completed: false },
		{ id: '3', text: 'Check emergency vet contacts', completed: false },
		// Add more items
	];

	return (
		<Paper p='md' withBorder>
			<Title order={3} mb='md'>
				Pre-Foaling Checklist
			</Title>
			<Stack>
				{checklistItems.map((item) => (
					<Checkbox
						key={item.id}
						label={item.text}
						checked={item.completed}
						styles={{
							label: { fontSize: '1rem' },
						}}
					/>
				))}
			</Stack>
		</Paper>
	);
}

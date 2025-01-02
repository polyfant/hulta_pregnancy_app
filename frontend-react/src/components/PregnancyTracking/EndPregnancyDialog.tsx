import { Button, Modal, Select, Stack, Textarea } from '@mantine/core';
import { DatePickerInput } from '@mantine/dates';
import { useState } from 'react';

interface EndPregnancyDialogProps {
	opened: boolean;
	onClose: () => void;
	onSubmit: (data: { outcome: string; foalingDate: string }) => void;
	isLoading: boolean;
}

export function EndPregnancyDialog({
	opened,
	onClose,
	onSubmit,
	isLoading,
}: EndPregnancyDialogProps) {
	const [outcome, setOutcome] = useState<string | null>(null);
	const [foalingDate, setFoalingDate] = useState<Date | null>(new Date());
	const [notes, setNotes] = useState('');

	const handleSubmit = () => {
		if (outcome && foalingDate) {
			onSubmit({
				outcome,
				foalingDate: foalingDate.toISOString().split('T')[0],
			});
		}
	};

	return (
		<Modal
			opened={opened}
			onClose={onClose}
			title='End Pregnancy Tracking'
			size='md'
		>
			<Stack>
				<Text size='sm' c='dimmed'>
					Please provide the details about the end of the pregnancy
					tracking. This action cannot be undone.
				</Text>

				<Select
					label='Outcome'
					placeholder='Select outcome'
					data={[
						{ value: 'SUCCESSFUL', label: 'Successful Foaling' },
						{ value: 'LOSS', label: 'Pregnancy Loss' },
						{ value: 'ERROR', label: 'Tracking Error' },
						{ value: 'OTHER', label: 'Other' },
					]}
					value={outcome}
					onChange={setOutcome}
					required
				/>

				<DatePickerInput
					label='End Date'
					placeholder='Select date'
					value={foalingDate}
					onChange={setFoalingDate}
					maxDate={new Date()}
					required
				/>

				<Textarea
					label='Additional Notes'
					placeholder='Enter any additional notes...'
					value={notes}
					onChange={(event) => setNotes(event.currentTarget.value)}
					minRows={3}
				/>

				<Button
					onClick={handleSubmit}
					loading={isLoading}
					fullWidth
					mt='md'
				>
					Confirm
				</Button>
			</Stack>
		</Modal>
	);
}

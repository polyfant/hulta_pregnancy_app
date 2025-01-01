import { Button, Modal, Stack } from '@mantine/core';
import { DatePickerInput } from '@mantine/dates';
import { useState } from 'react';

interface StartPregnancyDialogProps {
	opened: boolean;
	onClose: () => void;
	onSubmit: (date: string) => void;
	isLoading: boolean;
}

export function StartPregnancyDialog({
	opened,
	onClose,
	onSubmit,
	isLoading,
}: StartPregnancyDialogProps) {
	const [conceptionDate, setConceptionDate] = useState<Date | null>(
		new Date()
	);

	const handleSubmit = () => {
		if (conceptionDate) {
			// Format date as YYYY-MM-DD
			const formattedDate = conceptionDate.toISOString().split('T')[0];
			onSubmit(formattedDate);
		}
	};

	return (
		<Modal
			opened={opened}
			onClose={onClose}
			title='Start Pregnancy Tracking'
			size='md'
		>
			<Stack>
				<Text size='sm'>
					Please select the conception date to begin tracking this
					mare's pregnancy. This will help us calculate important
					milestones and the expected due date.
				</Text>

				<DatePickerInput
					label='Conception Date'
					placeholder='Select date'
					value={conceptionDate}
					onChange={setConceptionDate}
					maxDate={new Date()}
					required
				/>

				<Button
					onClick={handleSubmit}
					loading={isLoading}
					disabled={!conceptionDate}
					mt='md'
				>
					Start Tracking
				</Button>
			</Stack>
		</Modal>
	);
}

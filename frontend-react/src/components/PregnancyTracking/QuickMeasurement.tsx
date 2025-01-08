import { Badge, Button, Group, NumberInput, Stack, Text } from '@mantine/core';
import { useForm } from '@mantine/form';
import { notifications } from '@mantine/notifications';
import { CloudCheck, CloudSlash, Ruler, Scales } from '@phosphor-icons/react';
import { useEffect, useState } from 'react';
import { syncService } from '../../services/syncService';

export function QuickMeasurement({ foalId }: { foalId: number }) {
	const [isOnline, setIsOnline] = useState(navigator.onLine);
	const form = useForm({
		initialValues: {
			weight: '',
			height: '',
		},
	});

	useEffect(() => {
		const handleOnline = () => setIsOnline(true);
		const handleOffline = () => setIsOnline(false);

		window.addEventListener('online', handleOnline);
		window.addEventListener('offline', handleOffline);

		return () => {
			window.removeEventListener('online', handleOnline);
			window.removeEventListener('offline', handleOffline);
		};
	}, []);

	const handleSubmit = async (values) => {
		try {
			await syncService.saveMeasurement({
				foalId,
				...values,
				date: new Date().toISOString(),
			});

			notifications.show({
				title: isOnline
					? 'Measurement Added'
					: 'Measurement Saved Offline',
				message: isOnline
					? 'Growth data has been updated'
					: 'Data will sync when connection is restored',
				color: isOnline ? 'green' : 'yellow',
			});

			form.reset();
		} catch (error) {
			notifications.show({
				title: 'Error',
				message: 'Failed to save measurement',
				color: 'red',
			});
		}
	};

	return (
		<form onSubmit={form.onSubmit(handleSubmit)}>
			<Stack>
				<Group position='apart'>
					<Text fw={500}>Quick Measurement</Text>
					{!isOnline && (
						<Badge color='yellow' variant='light'>
							Offline Mode
						</Badge>
					)}
				</Group>

				<NumberInput
					{...form.getInputProps('weight')}
					label='Weight (kg)'
					placeholder='Enter weight'
					leftSection={<Scales size={16} />}
					hideControls
					inputMode='decimal'
				/>

				<NumberInput
					{...form.getInputProps('height')}
					label='Height (cm)'
					placeholder='Enter height'
					leftSection={<Ruler size={16} />}
					hideControls
					inputMode='decimal'
				/>

				<Button
					type='submit'
					fullWidth
					leftSection={isOnline ? <CloudCheck /> : <CloudSlash />}
				>
					Save Measurement
					{!isOnline && ' (Offline)'}
				</Button>
			</Stack>
		</form>
	);
}

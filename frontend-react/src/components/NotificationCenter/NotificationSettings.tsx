import { Button, Card, Stack, Switch, Text, TextInput } from '@mantine/core';
import { EnvelopeSimple } from '@phosphor-icons/react';
import { useState } from 'react';

export function NotificationSettings() {
	const [emailEnabled, setEmailEnabled] = useState(false);
	const [email, setEmail] = useState('');

	const handleSave = async () => {
		await fetch('/api/notifications/settings', {
			method: 'PUT',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({
				emailNotifications: emailEnabled,
				emailAddress: email,
				notificationTypes: {
					dueDate: true,
					vetChecks: true,
					weightChecks: true,
					stageChanges: true,
				},
			}),
		});
	};

	return (
		<Card withBorder>
			<Stack>
				<Text fw={500} size='lg'>
					Notification Settings
				</Text>

				<Switch
					label='Email Notifications'
					description='Receive important updates via email'
					checked={emailEnabled}
					onChange={(e) => setEmailEnabled(e.currentTarget.checked)}
					leftSection={<EnvelopeSimple size={16} />}
				/>

				{emailEnabled && (
					<TextInput
						label='Email Address'
						value={email}
						onChange={(e) => setEmail(e.currentTarget.value)}
						placeholder='your@email.com'
					/>
				)}

				<Button onClick={handleSave}>Save Settings</Button>
			</Stack>
		</Card>
	);
}

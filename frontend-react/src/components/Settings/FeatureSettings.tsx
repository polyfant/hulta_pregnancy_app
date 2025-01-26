import { Alert, Card, Select, Stack, Switch, Text } from '@mantine/core';
import { InfoCircle } from '@phosphor-icons/react';
import { useEffect, useState } from 'react';
import { userSettings } from '../../services/userSettings';

export function FeatureSettings() {
	const [settings, setSettings] = useState(null);

	useEffect(() => {
		userSettings.getSettings().then(setSettings);
	}, []);

	const handleToggle = async (feature: string, value: boolean) => {
		const updated = {
			...settings,
			[feature]: {
				...settings[feature],
				enabled: value,
			},
		};
		await userSettings.updateSettings(updated);
		setSettings(updated);
	};

	if (!settings) return null;

	return (
		<Card withBorder>
			<Stack>
				<Text size='lg' fw={500}>
					Feature Settings
				</Text>

				<Alert icon={<InfoCircle />} color='blue'>
					When features are disabled, all related data is permanently
					deleted from your device.
				</Alert>

				<Stack spacing='xs'>
					<Switch
						label='Environmental Monitoring'
						description='Track weather and environmental conditions'
						checked={settings.environmentalMonitoring.enabled}
						onChange={(e) =>
							handleToggle(
								'environmentalMonitoring',
								e.currentTarget.checked
							)
						}
					/>

					{settings.environmentalMonitoring.enabled && (
						<Stack spacing='xs' ml='md'>
							<Switch
								label='Air Quality Monitoring'
								size='sm'
								checked={
									settings.environmentalMonitoring.airQuality
								}
								onChange={(e) =>
									handleToggle(
										'environmentalMonitoring.airQuality',
										e.currentTarget.checked
									)
								}
							/>
							<Switch
								label='Weather Alerts'
								size='sm'
								checked={
									settings.environmentalMonitoring
										.weatherAlerts
								}
								onChange={(e) =>
									handleToggle(
										'environmentalMonitoring.weatherAlerts',
										e.currentTarget.checked
									)
								}
							/>
						</Stack>
					)}

					<Switch
						label='Location Services'
						description='Required for local weather data'
						checked={settings.locationTracking.enabled}
						onChange={(e) =>
							handleToggle(
								'locationTracking',
								e.currentTarget.checked
							)
						}
					/>

					{settings.locationTracking.enabled && (
						<Select
							label='Location Precision'
							size='sm'
							ml='md'
							value={settings.locationTracking.precision}
							onChange={(value) =>
								handleToggle(
									'locationTracking.precision',
									value
								)
							}
							data={[
								{ value: 'off', label: 'Disabled' },
								{ value: 'city', label: 'City Level Only' },
								{ value: 'precise', label: 'Precise Location' },
							]}
						/>
					)}
				</Stack>
			</Stack>
		</Card>
	);
}

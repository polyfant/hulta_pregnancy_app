import {
	Card,
	Grid,
	RingProgress,
	Stack,
	Text,
	ThemeIcon,
} from '@mantine/core';
import { Database, Eye, Lock, Shield } from '@phosphor-icons/react';
import { useQuery } from '@tanstack/react-query';
import { privacyAssessment } from '../../services/privacyAssessment';
import { privacyControls } from '../../services/privacyControls';

export function PrivacyDashboard() {
	const { data: settings } = useQuery({
		queryKey: ['privacy-settings'],
		queryFn: () => privacyControls.getPrivacySettings(),
	});

	const { data: assessment } = useQuery({
		queryKey: ['privacy-assessment'],
		queryFn: async () => {
			const features = ['location', 'measurements', 'environmental'];
			const impacts = await Promise.all(
				features.map((f) => privacyAssessment.assessFeaturePrivacy(f))
			);
			return impacts.flat();
		},
	});

	const privacyScore =
		assessment?.reduce((score, impact) => {
			const scores = { low: 1, medium: 0.5, high: 0 };
			return score + scores[impact.risk];
		}, 0) ?? 0;

	return (
		<Stack spacing='xl'>
			<Grid>
				<Grid.Col span={6}>
					<Card withBorder>
						<Stack align='center'>
							<RingProgress
								size={120}
								thickness={12}
								sections={[
									{
										value:
											(privacyScore /
												assessment?.length) *
											100,
										color: 'blue',
									},
								]}
								label={
									<Text size='xl' ta='center'>
										{Math.round(
											(privacyScore /
												assessment?.length) *
												100
										)}
										%
									</Text>
								}
							/>
							<Text>Privacy Score</Text>
						</Stack>
					</Card>
				</Grid.Col>
				<Grid.Col span={6}>
					<Card withBorder>
						<Stack>
							<Text size='lg' fw={500}>
								Active Protections
							</Text>
							<Grid>
								{[
									{
										icon: Shield,
										label: 'Data Masking',
										active: settings?.dataMasking
											.maskHorseNames,
									},
									{
										icon: Lock,
										label: 'Encryption',
										active: true,
									},
									{
										icon: Eye,
										label: 'Privacy Mode',
										active:
											settings?.dataSharing
												.shareAnonymized === false,
									},
									{
										icon: Database,
										label: 'Auto Cleanup',
										active: settings?.dataRetention
											.autoDeleteEnabled,
									},
								].map((protection, i) => (
									<Grid.Col span={6} key={i}>
										<Stack align='center' spacing={4}>
											<ThemeIcon
												size='lg'
												color={
													protection.active
														? 'green'
														: 'gray'
												}
												variant='light'
											>
												<protection.icon />
											</ThemeIcon>
											<Text size='sm'>
												{protection.label}
											</Text>
										</Stack>
									</Grid.Col>
								))}
							</Grid>
						</Stack>
					</Card>
				</Grid.Col>
			</Grid>

			<Card withBorder>
				<Stack>
					<Text size='lg' fw={500}>
						Privacy Impacts
					</Text>
					{assessment?.map((impact, i) => (
						<Card key={i} withBorder p='sm'>
							<Grid align='center'>
								<Grid.Col span={8}>
									<Stack spacing={4}>
										<Text fw={500}>{impact.category}</Text>
										<Text size='sm' c='dimmed'>
											{impact.description}
										</Text>
									</Stack>
								</Grid.Col>
								<Grid.Col span={4}>
									<Text
										ta='right'
										c={
											impact.risk === 'high'
												? 'red'
												: impact.risk === 'medium'
												? 'yellow'
												: 'green'
										}
									>
										{impact.risk.toUpperCase()} RISK
									</Text>
								</Grid.Col>
							</Grid>
						</Card>
					))}
				</Stack>
			</Card>
		</Stack>
	);
}

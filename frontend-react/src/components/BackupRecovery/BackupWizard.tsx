import { Button, Card, Group, Stack, Stepper, Text } from '@mantine/core';
import { useState } from 'react';
import { backupService } from '../../services/backupService';
import { encryptionService } from '../../services/encryptionService';

export function BackupWizard() {
	const [active, setActive] = useState(0);
	const [loading, setLoading] = useState(false);
	const [error, setError] = useState<string | null>(null);
	const [backupFile, setBackupFile] = useState<File | null>(null);
	const [password, setPassword] = useState('');

	const handleFileSelect = (event: React.ChangeEvent<HTMLInputElement>) => {
		const file = event.target.files?.[0];
		if (file) setBackupFile(file);
	};

	const handleRestore = async () => {
		try {
			setLoading(true);
			setError(null);

			const key = await encryptionService.generateKey(password);
			await backupService.restore(backupFile, key);

			setActive(3); // Success step
		} catch (err) {
			setError(err.message);
		} finally {
			setLoading(false);
		}
	};

	return (
		<Card withBorder>
			<Stack>
				<Stepper active={active} onStepClick={setActive}>
					<Stepper.Step
						label='Select Backup'
						description='Choose backup file'
					>
						<input
							type='file'
							accept='.json'
							onChange={handleFileSelect}
						/>
					</Stepper.Step>

					<Stepper.Step
						label='Verify'
						description='Enter recovery password'
					>
						<PasswordInput
							value={password}
							onChange={(e) => setPassword(e.target.value)}
							placeholder='Enter backup password'
						/>
					</Stepper.Step>

					<Stepper.Step
						label='Restore'
						description='Restore your data'
					>
						<Text>Ready to restore your data?</Text>
						{error && (
							<Text color='red' size='sm'>
								{error}
							</Text>
						)}
					</Stepper.Step>

					<Stepper.Completed>
						<Text color='green'>Backup restored successfully!</Text>
					</Stepper.Completed>
				</Stepper>

				<Group position='right'>
					<Button
						variant='default'
						onClick={() => setActive((a) => Math.max(0, a - 1))}
						disabled={active === 0}
					>
						Back
					</Button>
					<Button
						onClick={
							active === 2
								? handleRestore
								: () => setActive((a) => a + 1)
						}
						loading={loading}
						disabled={
							(active === 0 && !backupFile) ||
							(active === 1 && !password)
						}
					>
						{active === 2 ? 'Restore' : 'Next'}
					</Button>
				</Group>
			</Stack>
		</Card>
	);
}

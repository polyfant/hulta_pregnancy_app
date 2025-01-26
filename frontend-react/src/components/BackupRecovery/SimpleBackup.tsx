import {
	Button,
	Card,
	CopyButton,
	Stack,
	Text,
	TextInput,
} from '@mantine/core';
import { securityService } from '../../services/securityService';

export function SimpleBackup() {
	const [phrase, setPhrase] = useState('');

	const handleGenerateBackup = async () => {
		// Generate new recovery phrase if none exists
		if (!phrase) {
			const newPhrase = securityService.generateRecoveryPhrase();
			setPhrase(newPhrase);
		}

		try {
			await backupService.createBackup(phrase);
			notifications.show({
				title: 'Backup Created',
				message: 'Please save your recovery phrase somewhere safe!',
				color: 'green',
			});
		} catch (error) {
			notifications.show({
				title: 'Backup Failed',
				message: error.message,
				color: 'red',
			});
		}
	};

	return (
		<Card withBorder>
			<Stack>
				<Text size='sm' c='dimmed'>
					Your recovery phrase is like a password that helps protect
					your data. Write it down and keep it safe!
				</Text>

				<TextInput
					value={phrase}
					readOnly
					label='Recovery Phrase'
					placeholder='Generate a backup to get your phrase'
					rightSection={
						<CopyButton value={phrase}>
							{({ copied, copy }) => (
								<Button
									variant='subtle'
									size='xs'
									onClick={copy}
								>
									{copied ? 'Copied!' : 'Copy'}
								</Button>
							)}
						</CopyButton>
					}
				/>

				<Button onClick={handleGenerateBackup}>
					Create Protected Backup
				</Button>
			</Stack>
		</Card>
	);
}

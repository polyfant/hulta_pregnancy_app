import { motion } from 'framer-motion';
import {
	Clock,
	Download,
	EyeOff,
	FileText,
	Lock,
	RefreshCw,
	Shield,
	Trash2,
} from 'lucide-react';
import React, { useState } from 'react';
import { Badge } from './ui/badge';
import { Button } from './ui/button';
import {
	Card,
	CardContent,
	CardDescription,
	CardFooter,
	CardHeader,
	CardTitle,
} from './ui/card';
import { Progress } from './ui/progress';
import { Separator } from './ui/separator';
import { Switch } from './ui/switch';
import { Tabs, TabsContent, TabsList, TabsTrigger } from './ui/tabs';

interface PrivacySettings {
	dataSharing: {
		shareHorseProfiles: boolean;
		shareBreedingOutcomes: boolean;
		participateInResearch: boolean;
		shareAnonymousStatistics: boolean;
	};
	dataRetention: {
		keepHealthRecords: number; // months
		keepWeatherData: number; // days
		keepLocationData: number; // days
		autoDeleteOldRecords: boolean;
	};
	dataMasking: {
		maskHorseNames: boolean;
		maskOwnerInfo: boolean;
		maskFinancialData: boolean;
		maskBreedingOutcomes: boolean;
	};
}

interface PrivacyScore {
	overall: number;
	dataSharing: number;
	dataRetention: number;
	dataMasking: number;
	lastAudit: Date;
}

export const PrivacyDashboard: React.FC = () => {
	const [settings, setSettings] = useState<PrivacySettings>({
		dataSharing: {
			shareHorseProfiles: false,
			shareBreedingOutcomes: false,
			participateInResearch: true,
			shareAnonymousStatistics: true,
		},
		dataRetention: {
			keepHealthRecords: 36, // months
			keepWeatherData: 30, // days
			keepLocationData: 7, // days
			autoDeleteOldRecords: true,
		},
		dataMasking: {
			maskHorseNames: true,
			maskOwnerInfo: true,
			maskFinancialData: true,
			maskBreedingOutcomes: false,
		},
	});

	const [privacyScore, setPrivacyScore] = useState<PrivacyScore>({
		overall: 85,
		dataSharing: 90,
		dataRetention: 75,
		dataMasking: 90,
		lastAudit: new Date(Date.now() - 2 * 24 * 60 * 60 * 1000), // 2 days ago
	});

	const [recentChanges, setRecentChanges] = useState([
		{
			id: '1',
			timestamp: new Date(Date.now() - 10 * 60 * 1000), // 10 minutes ago
			setting: 'Data Retention: Weather Data',
			oldValue: '14 days',
			newValue: '30 days',
			user: 'Admin',
		},
		{
			id: '2',
			timestamp: new Date(Date.now() - 2 * 24 * 60 * 60 * 1000), // 2 days ago
			setting: 'Data Masking: Mask Owner Info',
			oldValue: 'Off',
			newValue: 'On',
			user: 'Admin',
		},
		{
			id: '3',
			timestamp: new Date(Date.now() - 5 * 24 * 60 * 60 * 1000), // 5 days ago
			setting: 'Data Sharing: Share Breeding Outcomes',
			oldValue: 'On',
			newValue: 'Off',
			user: 'System Audit',
		},
	]);

	// Calculate privacy score based on settings
	const calculatePrivacyScore = (
		newSettings: PrivacySettings
	): PrivacyScore => {
		// Data sharing score (higher is more private)
		const sharingScore =
			(!newSettings.dataSharing.shareHorseProfiles ? 25 : 0) +
			(!newSettings.dataSharing.shareBreedingOutcomes ? 25 : 0) +
			(!newSettings.dataSharing.participateInResearch ? 25 : 10) +
			(!newSettings.dataSharing.shareAnonymousStatistics ? 25 : 15);

		// Data retention score (lower retention periods are more private)
		const retentionScore =
			(newSettings.dataRetention.keepHealthRecords <= 12
				? 30
				: newSettings.dataRetention.keepHealthRecords <= 24
				? 25
				: newSettings.dataRetention.keepHealthRecords <= 36
				? 20
				: 15) +
			(newSettings.dataRetention.keepWeatherData <= 7
				? 30
				: newSettings.dataRetention.keepWeatherData <= 14
				? 25
				: newSettings.dataRetention.keepWeatherData <= 30
				? 20
				: 15) +
			(newSettings.dataRetention.keepLocationData <= 1
				? 30
				: newSettings.dataRetention.keepLocationData <= 7
				? 25
				: newSettings.dataRetention.keepLocationData <= 14
				? 20
				: 15) +
			(newSettings.dataRetention.autoDeleteOldRecords ? 10 : 0);

		// Data masking score
		const maskingScore =
			(newSettings.dataMasking.maskHorseNames ? 25 : 0) +
			(newSettings.dataMasking.maskOwnerInfo ? 25 : 0) +
			(newSettings.dataMasking.maskFinancialData ? 25 : 0) +
			(newSettings.dataMasking.maskBreedingOutcomes ? 25 : 0);

		// Overall score (weighted average)
		const overall = Math.round(
			sharingScore * 0.35 + retentionScore * 0.35 + maskingScore * 0.3
		);

		return {
			overall,
			dataSharing: sharingScore,
			dataRetention: retentionScore,
			dataMasking: maskingScore,
			lastAudit: new Date(),
		};
	};

	const handleSettingChange = (
		category: keyof PrivacySettings,
		setting: string,
		value: boolean | number
	) => {
		const newSettings = {
			...settings,
			[category]: {
				...settings[category],
				[setting]: value,
			},
		};

		setSettings(newSettings);

		// Log the change
		const newChange = {
			id: Date.now().toString(),
			timestamp: new Date(),
			setting: `${
				category.charAt(0).toUpperCase() + category.slice(1)
			}: ${
				setting.charAt(0).toUpperCase() +
				setting.slice(1).replace(/([A-Z])/g, ' $1')
			}`,
			oldValue:
				typeof settings[category][setting] === 'boolean'
					? settings[category][setting]
						? 'On'
						: 'Off'
					: `${settings[category][setting]} ${
							typeof settings[category][setting] === 'number' &&
							setting.includes('keep')
								? setting.includes('Health')
									? 'months'
									: 'days'
								: ''
					  }`,
			newValue:
				typeof value === 'boolean'
					? value
						? 'On'
						: 'Off'
					: `${value} ${
							setting.includes('keep')
								? setting.includes('Health')
									? 'months'
									: 'days'
								: ''
					  }`,
			user: 'Admin',
		};

		setRecentChanges([newChange, ...recentChanges.slice(0, 9)]);

		// Update privacy score
		setPrivacyScore(calculatePrivacyScore(newSettings));
	};

	const runPrivacyAudit = () => {
		// In a real app, this would trigger a comprehensive privacy audit
		// For this demo, we'll just recalculate the score and update the audit timestamp
		setPrivacyScore({
			...calculatePrivacyScore(settings),
			lastAudit: new Date(),
		});
	};

	// Helper function to determine score badge variant
	const getScoreBadgeVariant = (score: number) => {
		if (score >= 90) return 'success';
		if (score >= 70) return 'default';
		if (score >= 50) return 'warning';
		return 'destructive';
	};

	// Helper function to format retention period
	const formatRetentionPeriod = (value: number, unit: string) => {
		return `${value} ${unit}${value !== 1 ? 's' : ''}`;
	};

	return (
		<motion.div
			initial={{ opacity: 0, y: 20 }}
			animate={{ opacity: 1, y: 0 }}
			transition={{ duration: 0.5 }}
		>
			<Card className='w-full'>
				<CardHeader className='pb-2'>
					<div className='flex justify-between items-center'>
						<div>
							<CardTitle className='text-xl font-bold'>
								Privacy Dashboard
							</CardTitle>
							<CardDescription>
								Manage data privacy and retention settings
							</CardDescription>
						</div>
						<div className='flex items-center gap-2'>
							<Badge
								variant={getScoreBadgeVariant(
									privacyScore.overall
								)}
								className='px-3 py-1 rounded-full'
							>
								<Lock className='h-3 w-3 mr-1' />
								Privacy Score: {privacyScore.overall}
							</Badge>
							<Button
								variant='outline'
								size='sm'
								onClick={runPrivacyAudit}
							>
								<RefreshCw className='h-4 w-4 mr-2' />
								Run Audit
							</Button>
						</div>
					</div>
				</CardHeader>

				<CardContent className='space-y-6'>
					<div className='grid grid-cols-1 md:grid-cols-3 gap-4'>
						<div className='p-4 rounded-lg border flex flex-col items-center'>
							<Shield className='h-8 w-8 text-primary mb-2' />
							<p className='text-sm font-medium'>
								Overall Privacy Rating
							</p>
							<div className='w-full mt-2'>
								<Progress
									value={privacyScore.overall}
									className='h-2'
								/>
							</div>
							<p className='text-xl font-bold mt-2'>
								{privacyScore.overall}%
							</p>
							<p className='text-xs text-muted-foreground mt-1'>
								Last audited:{' '}
								{privacyScore.lastAudit.toLocaleDateString()}
							</p>
						</div>

						<div className='p-4 rounded-lg border flex flex-col'>
							<div className='flex items-center gap-2 mb-2'>
								<EyeOff className='h-5 w-5 text-blue-500' />
								<p className='text-sm font-medium'>
									Data Sharing
								</p>
							</div>
							<div className='w-full'>
								<Progress
									value={privacyScore.dataSharing}
									className='h-2'
								/>
							</div>
							<p className='text-xs text-muted-foreground mt-1 flex-1'>
								Controls who can see your horse breeding data
							</p>
							<p className='text-right text-sm font-medium'>
								{privacyScore.dataSharing}%
							</p>
						</div>

						<div className='p-4 rounded-lg border flex flex-col'>
							<div className='flex items-center gap-2 mb-2'>
								<Clock className='h-5 w-5 text-purple-500' />
								<p className='text-sm font-medium'>
									Data Retention
								</p>
							</div>
							<div className='w-full'>
								<Progress
									value={privacyScore.dataRetention}
									className='h-2'
								/>
							</div>
							<p className='text-xs text-muted-foreground mt-1 flex-1'>
								Defines how long your data is stored
							</p>
							<p className='text-right text-sm font-medium'>
								{privacyScore.dataRetention}%
							</p>
						</div>
					</div>

					<Tabs defaultValue='dataSharing' className='w-full'>
						<TabsList className='grid w-full grid-cols-3'>
							<TabsTrigger value='dataSharing'>
								<EyeOff className='h-4 w-4 mr-2' />
								Data Sharing
							</TabsTrigger>
							<TabsTrigger value='dataRetention'>
								<Clock className='h-4 w-4 mr-2' />
								Data Retention
							</TabsTrigger>
							<TabsTrigger value='dataMasking'>
								<FileText className='h-4 w-4 mr-2' />
								Data Masking
							</TabsTrigger>
						</TabsList>

						{/* Data Sharing */}
						<TabsContent
							value='dataSharing'
							className='space-y-4 p-4 border rounded-lg mt-4'
						>
							<div className='flex items-center justify-between'>
								<div>
									<p className='font-medium'>
										Share Horse Profiles
									</p>
									<p className='text-sm text-muted-foreground'>
										Allow other farms to view basic
										information about your horses
									</p>
								</div>
								<Switch
									checked={
										settings.dataSharing.shareHorseProfiles
									}
									onCheckedChange={(checked) =>
										handleSettingChange(
											'dataSharing',
											'shareHorseProfiles',
											checked
										)
									}
								/>
							</div>

							<Separator />

							<div className='flex items-center justify-between'>
								<div>
									<p className='font-medium'>
										Share Breeding Outcomes
									</p>
									<p className='text-sm text-muted-foreground'>
										Make your breeding success data
										accessible to other breeders
									</p>
								</div>
								<Switch
									checked={
										settings.dataSharing
											.shareBreedingOutcomes
									}
									onCheckedChange={(checked) =>
										handleSettingChange(
											'dataSharing',
											'shareBreedingOutcomes',
											checked
										)
									}
								/>
							</div>

							<Separator />

							<div className='flex items-center justify-between'>
								<div>
									<p className='font-medium'>
										Participate in Research
									</p>
									<p className='text-sm text-muted-foreground'>
										Contribute anonymized data to equine
										research initiatives
									</p>
								</div>
								<Switch
									checked={
										settings.dataSharing
											.participateInResearch
									}
									onCheckedChange={(checked) =>
										handleSettingChange(
											'dataSharing',
											'participateInResearch',
											checked
										)
									}
								/>
							</div>

							<Separator />

							<div className='flex items-center justify-between'>
								<div>
									<p className='font-medium'>
										Share Anonymous Statistics
									</p>
									<p className='text-sm text-muted-foreground'>
										Contribute to industry-wide breeding
										statistics (no identifiable data)
									</p>
								</div>
								<Switch
									checked={
										settings.dataSharing
											.shareAnonymousStatistics
									}
									onCheckedChange={(checked) =>
										handleSettingChange(
											'dataSharing',
											'shareAnonymousStatistics',
											checked
										)
									}
								/>
							</div>
						</TabsContent>

						{/* Data Retention */}
						<TabsContent
							value='dataRetention'
							className='space-y-4 p-4 border rounded-lg mt-4'
						>
							<div>
								<div className='flex items-center justify-between mb-1'>
									<p className='font-medium'>
										Health Records Retention
									</p>
									<p className='text-sm font-medium'>
										{formatRetentionPeriod(
											settings.dataRetention
												.keepHealthRecords,
											'month'
										)}
									</p>
								</div>
								<p className='text-sm text-muted-foreground mb-2'>
									How long to keep detailed health records
								</p>
								<div className='flex space-x-2'>
									{[12, 24, 36, 48].map((months) => (
										<Button
											key={months}
											variant={
												settings.dataRetention
													.keepHealthRecords ===
												months
													? 'default'
													: 'outline'
											}
											size='sm'
											onClick={() =>
												handleSettingChange(
													'dataRetention',
													'keepHealthRecords',
													months
												)
											}
										>
											{months} Months
										</Button>
									))}
								</div>
							</div>

							<Separator />

							<div>
								<div className='flex items-center justify-between mb-1'>
									<p className='font-medium'>
										Weather Data Retention
									</p>
									<p className='text-sm font-medium'>
										{formatRetentionPeriod(
											settings.dataRetention
												.keepWeatherData,
											'day'
										)}
									</p>
								</div>
								<p className='text-sm text-muted-foreground mb-2'>
									How long to keep environmental data
								</p>
								<div className='flex space-x-2'>
									{[7, 14, 30, 60].map((days) => (
										<Button
											key={days}
											variant={
												settings.dataRetention
													.keepWeatherData === days
													? 'default'
													: 'outline'
											}
											size='sm'
											onClick={() =>
												handleSettingChange(
													'dataRetention',
													'keepWeatherData',
													days
												)
											}
										>
											{days} Days
										</Button>
									))}
								</div>
							</div>

							<Separator />

							<div>
								<div className='flex items-center justify-between mb-1'>
									<p className='font-medium'>
										Location Data Retention
									</p>
									<p className='text-sm font-medium'>
										{formatRetentionPeriod(
											settings.dataRetention
												.keepLocationData,
											'day'
										)}
									</p>
								</div>
								<p className='text-sm text-muted-foreground mb-2'>
									How long to keep horse location tracking
									data
								</p>
								<div className='flex space-x-2'>
									{[1, 7, 14, 30].map((days) => (
										<Button
											key={days}
											variant={
												settings.dataRetention
													.keepLocationData === days
													? 'default'
													: 'outline'
											}
											size='sm'
											onClick={() =>
												handleSettingChange(
													'dataRetention',
													'keepLocationData',
													days
												)
											}
										>
											{days} Day{days !== 1 ? 's' : ''}
										</Button>
									))}
								</div>
							</div>

							<Separator />

							<div className='flex items-center justify-between'>
								<div>
									<p className='font-medium'>
										Auto-Delete Old Records
									</p>
									<p className='text-sm text-muted-foreground'>
										Automatically purge data when it reaches
										retention limits
									</p>
								</div>
								<Switch
									checked={
										settings.dataRetention
											.autoDeleteOldRecords
									}
									onCheckedChange={(checked) =>
										handleSettingChange(
											'dataRetention',
											'autoDeleteOldRecords',
											checked
										)
									}
								/>
							</div>
						</TabsContent>

						{/* Data Masking */}
						<TabsContent
							value='dataMasking'
							className='space-y-4 p-4 border rounded-lg mt-4'
						>
							<div className='flex items-center justify-between'>
								<div>
									<p className='font-medium'>
										Mask Horse Names
									</p>
									<p className='text-sm text-muted-foreground'>
										Replace real horse names with aliases in
										shared reports
									</p>
								</div>
								<Switch
									checked={
										settings.dataMasking.maskHorseNames
									}
									onCheckedChange={(checked) =>
										handleSettingChange(
											'dataMasking',
											'maskHorseNames',
											checked
										)
									}
								/>
							</div>

							<Separator />

							<div className='flex items-center justify-between'>
								<div>
									<p className='font-medium'>
										Mask Owner Information
									</p>
									<p className='text-sm text-muted-foreground'>
										Hide owner contact and personal details
										from shared data
									</p>
								</div>
								<Switch
									checked={settings.dataMasking.maskOwnerInfo}
									onCheckedChange={(checked) =>
										handleSettingChange(
											'dataMasking',
											'maskOwnerInfo',
											checked
										)
									}
								/>
							</div>

							<Separator />

							<div className='flex items-center justify-between'>
								<div>
									<p className='font-medium'>
										Mask Financial Data
									</p>
									<p className='text-sm text-muted-foreground'>
										Hide costs, prices, and other financial
										information
									</p>
								</div>
								<Switch
									checked={
										settings.dataMasking.maskFinancialData
									}
									onCheckedChange={(checked) =>
										handleSettingChange(
											'dataMasking',
											'maskFinancialData',
											checked
										)
									}
								/>
							</div>

							<Separator />

							<div className='flex items-center justify-between'>
								<div>
									<p className='font-medium'>
										Mask Breeding Outcomes
									</p>
									<p className='text-sm text-muted-foreground'>
										Hide specific details of breeding
										successes and failures
									</p>
								</div>
								<Switch
									checked={
										settings.dataMasking
											.maskBreedingOutcomes
									}
									onCheckedChange={(checked) =>
										handleSettingChange(
											'dataMasking',
											'maskBreedingOutcomes',
											checked
										)
									}
								/>
							</div>
						</TabsContent>
					</Tabs>

					{/* Recent Changes */}
					<div className='space-y-2 p-4 border rounded-lg'>
						<h3 className='text-sm font-medium mb-2'>
							Privacy Change Log
						</h3>

						<div className='space-y-3'>
							{recentChanges.slice(0, 3).map((change) => (
								<div
									key={change.id}
									className='flex items-start gap-3 text-sm border-b pb-2'
								>
									<div className='text-xs text-muted-foreground min-w-[80px]'>
										{change.timestamp.toLocaleDateString()}
									</div>
									<div className='flex-1'>
										<p className='font-medium'>
											{change.setting}
										</p>
										<p className='text-xs text-muted-foreground'>
											Changed from{' '}
											<span className='text-amber-600'>
												{change.oldValue}
											</span>{' '}
											to{' '}
											<span className='text-green-600'>
												{change.newValue}
											</span>
										</p>
									</div>
									<div className='text-xs text-muted-foreground'>
										{change.user}
									</div>
								</div>
							))}

							{recentChanges.length > 3 && (
								<Button
									variant='link'
									size='sm'
									className='w-full mt-2'
								>
									View all {recentChanges.length} changes
								</Button>
							)}
						</div>
					</div>
				</CardContent>

				<CardFooter className='border-t p-4 flex flex-col sm:flex-row gap-2 justify-between items-center'>
					<div className='flex gap-2'>
						<Button
							variant='outline'
							size='sm'
							className='flex-1 sm:flex-none'
						>
							<Download className='h-4 w-4 mr-2' />
							Export Data
						</Button>
						<Button
							variant='destructive'
							size='sm'
							className='flex-1 sm:flex-none'
						>
							<Trash2 className='h-4 w-4 mr-2' />
							Purge Data
						</Button>
					</div>
					<p className='text-xs text-muted-foreground'>
						Privacy settings are compliant with equine industry best
						practices
					</p>
				</CardFooter>
			</Card>
		</motion.div>
	);
};

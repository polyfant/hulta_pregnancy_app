import { motion } from 'framer-motion';
import {
	CheckCircle2,
	ChevronDown,
	ChevronUp,
	Dna,
	Heart,
	Trophy,
	X,
} from 'lucide-react';
import React, { useMemo } from 'react';
import { Button } from '../components/ui/button';
import { Horse } from '../types/horse';
import { Badge } from './ui/badge';
import { Card, CardContent, CardHeader, CardTitle } from './ui/card';
import { Progress } from './ui/progress';

interface BreedingStatsProps {
	horses: Horse[];
	breedingHistory?: BreedingRecord[];
}

interface BreedingRecord {
	id: string;
	mareId: string;
	mareName: string;
	stallionId: string;
	stallionName: string;
	breedingDate: string;
	successful: boolean;
	foalId?: string;
	foalName?: string;
	foalBirthDate?: string;
	notes?: string;
}

export const BreedingStats: React.FC<BreedingStatsProps> = ({
	horses,
	// Mock data for demonstration
	breedingHistory = [
		{
			id: '1',
			mareId: '101',
			mareName: 'Bella',
			stallionId: '201',
			stallionName: 'Thunder',
			breedingDate: '2024-01-15',
			successful: true,
			foalId: '301',
			foalName: 'Storm',
			foalBirthDate: '2024-12-20',
			notes: 'Healthy foal, uncomplicated birth',
		},
		{
			id: '2',
			mareId: '102',
			mareName: 'Luna',
			stallionId: '202',
			stallionName: 'Apollo',
			breedingDate: '2024-02-10',
			successful: true,
			foalId: '302',
			foalName: 'Nova',
			foalBirthDate: '2025-01-15',
			notes: 'Slight complications but healthy foal',
		},
		{
			id: '3',
			mareId: '103',
			mareName: 'Star',
			stallionId: '203',
			stallionName: 'Blaze',
			breedingDate: '2024-03-05',
			successful: false,
			notes: 'Failed to conceive',
		},
		{
			id: '4',
			mareId: '104',
			mareName: 'Misty',
			stallionId: '204',
			stallionName: 'Shadow',
			breedingDate: '2024-03-20',
			successful: true,
			foalId: '303',
			foalName: 'Phantom',
			foalBirthDate: '2025-02-25',
			notes: 'Healthy foal',
		},
		{
			id: '5',
			mareId: '101',
			mareName: 'Bella',
			stallionId: '205',
			stallionName: 'Duke',
			breedingDate: '2023-04-10',
			successful: true,
			foalId: '304',
			foalName: 'Princess',
			foalBirthDate: '2024-03-15',
			notes: 'Healthy filly',
		},
	],
}) => {
	// Calculate overall statistics
	const stats = useMemo(() => {
		const totalBreedingAttempts = breedingHistory.length;
		const successfulBreedings = breedingHistory.filter(
			(b) => b.successful
		).length;
		const successRate =
			totalBreedingAttempts > 0
				? (successfulBreedings / totalBreedingAttempts) * 100
				: 0;

		// Current year statistics
		const currentYear = new Date().getFullYear();
		const currentYearBreedings = breedingHistory.filter(
			(b) => new Date(b.breedingDate).getFullYear() === currentYear
		);
		const currentYearSuccessful = currentYearBreedings.filter(
			(b) => b.successful
		).length;
		const currentYearRate =
			currentYearBreedings.length > 0
				? (currentYearSuccessful / currentYearBreedings.length) * 100
				: 0;

		// Last year statistics for comparison
		const lastYear = currentYear - 1;
		const lastYearBreedings = breedingHistory.filter(
			(b) => new Date(b.breedingDate).getFullYear() === lastYear
		);
		const lastYearSuccessful = lastYearBreedings.filter(
			(b) => b.successful
		).length;
		const lastYearRate =
			lastYearBreedings.length > 0
				? (lastYearSuccessful / lastYearBreedings.length) * 100
				: 0;

		// Change from previous year
		const successRateChange =
			lastYearRate > 0 ? currentYearRate - lastYearRate : 0;

		// Top stallions by success rate (with minimum 2 attempts)
		const stallionStats = breedingHistory.reduce((acc, record) => {
			if (!acc[record.stallionId]) {
				acc[record.stallionId] = {
					id: record.stallionId,
					name: record.stallionName,
					attempts: 0,
					successful: 0,
				};
			}

			acc[record.stallionId].attempts += 1;
			if (record.successful) {
				acc[record.stallionId].successful += 1;
			}

			return acc;
		}, {} as Record<string, { id: string; name: string; attempts: number; successful: number }>);

		const topStallions = Object.values(stallionStats)
			.filter((s) => s.attempts >= 2)
			.map((s) => ({
				...s,
				successRate: (s.successful / s.attempts) * 100,
			}))
			.sort((a, b) => b.successRate - a.successRate)
			.slice(0, 3);

		// Most productive mares
		const mareStats = breedingHistory.reduce((acc, record) => {
			if (!acc[record.mareId]) {
				acc[record.mareId] = {
					id: record.mareId,
					name: record.mareName,
					attempts: 0,
					successful: 0,
				};
			}

			acc[record.mareId].attempts += 1;
			if (record.successful) {
				acc[record.mareId].successful += 1;
			}

			return acc;
		}, {} as Record<string, { id: string; name: string; attempts: number; successful: number }>);

		const topMares = Object.values(mareStats)
			.sort((a, b) => b.successful - a.successful)
			.slice(0, 3);

		return {
			totalBreedingAttempts,
			successfulBreedings,
			successRate: Math.round(successRate),
			currentYearBreedings: currentYearBreedings.length,
			currentYearSuccessful,
			currentYearRate: Math.round(currentYearRate),
			successRateChange: Math.round(successRateChange),
			topStallions,
			topMares,
		};
	}, [breedingHistory]);

	return (
		<motion.div
			initial={{ opacity: 0, y: 20 }}
			animate={{ opacity: 1, y: 0 }}
			transition={{ duration: 0.5 }}
		>
			<Card className='w-full'>
				<CardHeader className='pb-2'>
					<CardTitle className='text-xl font-bold'>
						Breeding Statistics
					</CardTitle>
				</CardHeader>

				<CardContent className='space-y-6'>
					{/* Summary Statistics */}
					<div className='grid grid-cols-1 md:grid-cols-3 gap-4'>
						<div className='p-4 rounded-lg border'>
							<div className='flex items-center justify-between'>
								<div>
									<p className='text-sm text-muted-foreground'>
										Success Rate
									</p>
									<p className='text-2xl font-bold'>
										{stats.successRate}%
									</p>
								</div>
								<div
									className={`p-2 rounded-full ${
										stats.successRate >= 70
											? 'bg-green-100'
											: stats.successRate >= 50
											? 'bg-amber-100'
											: 'bg-red-100'
									}`}
								>
									<Trophy
										className={`h-6 w-6 ${
											stats.successRate >= 70
												? 'text-green-500'
												: stats.successRate >= 50
												? 'text-amber-500'
												: 'text-red-500'
										}`}
									/>
								</div>
							</div>
							<Progress
								value={stats.successRate}
								className='h-1.5 mt-2'
							/>
							<div className='flex items-center gap-1 mt-2'>
								{stats.successRateChange > 0 ? (
									<>
										<ChevronUp className='h-4 w-4 text-green-500' />
										<span className='text-xs text-green-500'>
											{Math.abs(stats.successRateChange)}%
											from last year
										</span>
									</>
								) : stats.successRateChange < 0 ? (
									<>
										<ChevronDown className='h-4 w-4 text-red-500' />
										<span className='text-xs text-red-500'>
											{Math.abs(stats.successRateChange)}%
											from last year
										</span>
									</>
								) : (
									<span className='text-xs text-muted-foreground'>
										No change from last year
									</span>
								)}
							</div>
						</div>

						<div className='p-4 rounded-lg border'>
							<div className='flex items-center justify-between'>
								<div>
									<p className='text-sm text-muted-foreground'>
										Total Breedings
									</p>
									<p className='text-2xl font-bold'>
										{stats.totalBreedingAttempts}
									</p>
								</div>
								<div className='p-2 rounded-full bg-blue-100'>
									<Heart className='h-6 w-6 text-blue-500' />
								</div>
							</div>
							<div className='flex justify-between mt-2'>
								<div>
									<p className='text-xs text-muted-foreground'>
										Successful
									</p>
									<p className='text-sm font-medium'>
										{stats.successfulBreedings}
									</p>
								</div>
								<div>
									<p className='text-xs text-muted-foreground'>
										Current Year
									</p>
									<p className='text-sm font-medium'>
										{stats.currentYearBreedings}
									</p>
								</div>
								<div>
									<p className='text-xs text-muted-foreground'>
										Success Rate
									</p>
									<p className='text-sm font-medium'>
										{stats.currentYearRate}%
									</p>
								</div>
							</div>
						</div>

						<div className='p-4 rounded-lg border'>
							<div className='flex items-center justify-between'>
								<div>
									<p className='text-sm text-muted-foreground'>
										Breeding Pairs
									</p>
									<p className='text-2xl font-bold'>
										{breedingHistory.length}
									</p>
								</div>
								<div className='p-2 rounded-full bg-purple-100'>
									<Dna className='h-6 w-6 text-purple-500' />
								</div>
							</div>
							<div className='mt-2'>
								<div className='flex justify-between text-xs text-muted-foreground mb-1'>
									<span>Mares</span>
									<span>Stallions</span>
								</div>
								<div className='flex justify-between'>
									<p className='text-sm font-medium'>
										{Object.keys(stats.topMares).length}
									</p>
									<p className='text-sm font-medium'>
										{Object.keys(stats.topStallions).length}
									</p>
								</div>
							</div>
						</div>
					</div>

					{/* Top Performers */}
					<div className='grid grid-cols-1 md:grid-cols-2 gap-6'>
						{/* Top Stallions */}
						<div>
							<h3 className='font-semibold mb-3'>
								Top Stallions
							</h3>
							<div className='space-y-3'>
								{stats.topStallions.length > 0 ? (
									stats.topStallions.map(
										(stallion, index) => (
											<div
												key={stallion.id}
												className='flex items-center gap-3 p-3 rounded-lg border'
											>
												<div className='flex items-center justify-center h-8 w-8 rounded-full bg-amber-100 text-amber-500 font-medium'>
													{index + 1}
												</div>
												<div className='flex-1'>
													<div className='flex justify-between items-center'>
														<p className='font-medium'>
															{stallion.name}
														</p>
														<Badge
															variant={
																index === 0
																	? 'default'
																	: 'outline'
															}
														>
															{Math.round(
																stallion.successRate
															)}
															%
														</Badge>
													</div>
													<div className='flex justify-between items-center mt-1'>
														<p className='text-xs text-muted-foreground'>
															{
																stallion.successful
															}
															/{stallion.attempts}{' '}
															successful breedings
														</p>
													</div>
													<Progress
														value={
															stallion.successRate
														}
														className='h-1 mt-2'
													/>
												</div>
											</div>
										)
									)
								) : (
									<div className='text-center py-4 text-muted-foreground'>
										No stallion data available
									</div>
								)}
							</div>
						</div>

						{/* Top Mares */}
						<div>
							<h3 className='font-semibold mb-3'>
								Most Productive Mares
							</h3>
							<div className='space-y-3'>
								{stats.topMares.length > 0 ? (
									stats.topMares.map((mare, index) => (
										<div
											key={mare.id}
											className='flex items-center gap-3 p-3 rounded-lg border'
										>
											<div className='flex items-center justify-center h-8 w-8 rounded-full bg-pink-100 text-pink-500 font-medium'>
												{index + 1}
											</div>
											<div className='flex-1'>
												<div className='flex justify-between items-center'>
													<p className='font-medium'>
														{mare.name}
													</p>
													<Badge
														variant={
															index === 0
																? 'default'
																: 'outline'
														}
													>
														{mare.successful} foals
													</Badge>
												</div>
												<div className='flex justify-between items-center mt-1'>
													<p className='text-xs text-muted-foreground'>
														{mare.successful}/
														{mare.attempts}{' '}
														successful breedings
													</p>
													<p className='text-xs font-medium'>
														{Math.round(
															(mare.successful /
																mare.attempts) *
																100
														)}
														% success
													</p>
												</div>
												<Progress
													value={
														(mare.successful /
															mare.attempts) *
														100
													}
													className='h-1 mt-2'
												/>
											</div>
										</div>
									))
								) : (
									<div className='text-center py-4 text-muted-foreground'>
										No mare data available
									</div>
								)}
							</div>
						</div>
					</div>

					{/* Breeding History Summary */}
					<div>
						<h3 className='font-semibold mb-3'>
							Recent Breeding Outcomes
						</h3>
						<div className='space-y-3'>
							{breedingHistory.slice(0, 3).map((record) => (
								<div
									key={record.id}
									className='flex items-start gap-3 p-3 rounded-lg border'
								>
									<div className='mt-1'>
										<div
											className={`h-8 w-8 rounded-full flex items-center justify-center ${
												record.successful
													? 'bg-green-100 text-green-500'
													: 'bg-red-100 text-red-500'
											}`}
										>
											{record.successful ? (
												<CheckCircle2 className='h-5 w-5' />
											) : (
												<X className='h-5 w-5' />
											)}
										</div>
									</div>
									<div className='flex-1'>
										<div className='flex justify-between items-start'>
											<div>
												<p className='font-medium'>
													{record.mareName} ×{' '}
													{record.stallionName}
												</p>
												<p className='text-sm text-muted-foreground'>
													Bred:{' '}
													{new Date(
														record.breedingDate
													).toLocaleDateString()}
												</p>
												{record.successful &&
													record.foalName && (
														<div className='flex items-center gap-2 mt-1'>
															<Badge
																variant='outline'
																className='text-green-500 border-green-200 bg-green-50'
															>
																Foal:{' '}
																{
																	record.foalName
																}
															</Badge>
															{record.foalBirthDate && (
																<span className='text-xs text-muted-foreground'>
																	Born:{' '}
																	{new Date(
																		record.foalBirthDate
																	).toLocaleDateString()}
																</span>
															)}
														</div>
													)}
											</div>
										</div>

										{record.notes && (
											<p className='text-xs text-muted-foreground mt-2'>
												{record.notes}
											</p>
										)}
									</div>
								</div>
							))}

							{breedingHistory.length > 3 && (
								<div className='text-center'>
									<Button variant='link' size='sm'>
										View all {breedingHistory.length}{' '}
										breeding records
									</Button>
								</div>
							)}
						</div>
					</div>
				</CardContent>
			</Card>
		</motion.div>
	);
};

import { addDays, differenceInDays, format } from 'date-fns';
import { motion } from 'framer-motion';
import {
	AlertCircle,
	Calendar,
	CheckCircle2,
	Clock,
	Heart,
} from 'lucide-react';
import React, { useState } from 'react';
import { Horse } from '../types/horse';
import { Badge } from './ui/badge';
import { Card, CardContent, CardHeader, CardTitle } from './ui/card';
import { Progress } from './ui/progress';

interface PregnancyTrackerProps {
	horse: Horse;
}

// Pregnancy stages for horses (typically 340 days)
const STAGES = [
	{ name: 'Early', days: [0, 114], color: 'bg-blue-500' },
	{ name: 'Mid', days: [115, 225], color: 'bg-indigo-500' },
	{ name: 'Late', days: [226, 310], color: 'bg-purple-500' },
	{ name: 'Pre-foaling', days: [311, 340], color: 'bg-pink-500' },
];

// Milestones during pregnancy
const MILESTONES = [
	{
		day: 30,
		name: 'Heartbeat detectable',
		icon: <Heart className='h-4 w-4 text-red-500' />,
	},
	{
		day: 60,
		name: 'Fetus development begins',
		icon: <Clock className='h-4 w-4 text-blue-500' />,
	},
	{
		day: 150,
		name: 'Mid-term checkup recommended',
		icon: <AlertCircle className='h-4 w-4 text-amber-500' />,
	},
	{
		day: 270,
		name: 'Begin foaling preparations',
		icon: <CheckCircle2 className='h-4 w-4 text-green-500' />,
	},
	{
		day: 320,
		name: 'Foaling imminent',
		icon: <AlertCircle className='h-4 w-4 text-red-500' />,
	},
];

export const PregnancyTracker: React.FC<PregnancyTrackerProps> = ({
	horse,
}) => {
	const [expandedSection, setExpandedSection] = useState<string | null>(null);

	// If the horse isn't pregnant, display appropriate message
	if (!horse.isPregnant || !horse.lastBreedingDate) {
		return (
			<Card className='w-full'>
				<CardHeader>
					<CardTitle className='text-xl font-bold'>
						Pregnancy Tracker
					</CardTitle>
				</CardHeader>
				<CardContent>
					<div className='flex items-center justify-center h-32 text-muted-foreground'>
						Not currently pregnant
					</div>
				</CardContent>
			</Card>
		);
	}

	// Calculate pregnancy details
	const breedingDate = new Date(horse.lastBreedingDate);
	const dueDate = horse.dueDate
		? new Date(horse.dueDate)
		: addDays(breedingDate, 340);
	const today = new Date();

	// Calculate days pregnant and progress percentage
	const daysPregnant = differenceInDays(today, breedingDate);
	const totalDays = 340; // Average gestation period for horses
	const progressPercentage = Math.min(
		Math.round((daysPregnant / totalDays) * 100),
		100
	);

	// Determine current stage
	const currentStage =
		STAGES.find(
			(stage) =>
				daysPregnant >= stage.days[0] && daysPregnant <= stage.days[1]
		) || STAGES[STAGES.length - 1];

	// Get upcoming milestones
	const upcomingMilestones = MILESTONES.filter(
		(milestone) => milestone.day > daysPregnant
	)
		.sort((a, b) => a.day - b.day)
		.slice(0, 3);

	// Get completed milestones
	const completedMilestones = MILESTONES.filter(
		(milestone) => milestone.day <= daysPregnant
	).sort((a, b) => b.day - a.day);

	const toggleSection = (section: string) => {
		if (expandedSection === section) {
			setExpandedSection(null);
		} else {
			setExpandedSection(section);
		}
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
						<CardTitle className='text-xl font-bold'>
							Pregnancy Tracker
						</CardTitle>
						<Badge variant='outline' className='font-medium'>
							{currentStage.name} Stage
						</Badge>
					</div>
				</CardHeader>
				<CardContent className='space-y-6'>
					{/* Progress bar */}
					<div className='space-y-2'>
						<div className='flex justify-between text-sm'>
							<span className='text-muted-foreground'>
								Day {daysPregnant} of {totalDays}
							</span>
							<span className='font-medium'>
								{progressPercentage}%
							</span>
						</div>
						<Progress value={progressPercentage} className='h-2' />
					</div>

					{/* Important dates */}
					<div className='grid grid-cols-1 md:grid-cols-2 gap-4'>
						<div className='flex items-center gap-3 p-3 rounded-lg border'>
							<Calendar className='h-5 w-5 text-blue-500' />
							<div>
								<p className='text-sm font-medium'>
									Breeding Date
								</p>
								<p className='text-sm text-muted-foreground'>
									{format(breedingDate, 'MMM d, yyyy')}
								</p>
							</div>
						</div>
						<div className='flex items-center gap-3 p-3 rounded-lg border'>
							<Calendar className='h-5 w-5 text-pink-500' />
							<div>
								<p className='text-sm font-medium'>
									Expected Due Date
								</p>
								<p className='text-sm text-muted-foreground'>
									{format(dueDate, 'MMM d, yyyy')}
								</p>
							</div>
						</div>
					</div>

					{/* Stage timeline */}
					<div>
						<button
							className='flex justify-between items-center w-full py-2 text-left font-medium'
							onClick={() => toggleSection('stages')}
						>
							<span>Pregnancy Stages</span>
							<span className='text-muted-foreground text-sm'>
								{expandedSection === 'stages' ? 'Hide' : 'Show'}
							</span>
						</button>

						{expandedSection === 'stages' && (
							<motion.div
								initial={{ opacity: 0, height: 0 }}
								animate={{ opacity: 1, height: 'auto' }}
								transition={{ duration: 0.3 }}
								className='mt-2 space-y-3'
							>
								{STAGES.map((stage, index) => (
									<div
										key={index}
										className='flex items-center gap-3'
									>
										<div
											className={`w-3 h-3 rounded-full ${stage.color}`}
										></div>
										<div className='flex-1'>
											<div className='flex justify-between'>
												<span className='font-medium'>
													{stage.name}
												</span>
												<span className='text-sm text-muted-foreground'>
													Days {stage.days[0]}-
													{stage.days[1]}
												</span>
											</div>
											<Progress
												value={
													currentStage.name ===
													stage.name
														? ((daysPregnant -
																stage.days[0]) /
																(stage.days[1] -
																	stage
																		.days[0])) *
														  100
														: daysPregnant >
														  stage.days[1]
														? 100
														: 0
												}
												className='h-1 mt-1'
											/>
										</div>
									</div>
								))}
							</motion.div>
						)}
					</div>

					{/* Upcoming milestones */}
					<div>
						<button
							className='flex justify-between items-center w-full py-2 text-left font-medium'
							onClick={() => toggleSection('upcoming')}
						>
							<span>Upcoming Milestones</span>
							<span className='text-muted-foreground text-sm'>
								{expandedSection === 'upcoming'
									? 'Hide'
									: 'Show'}
							</span>
						</button>

						{expandedSection === 'upcoming' && (
							<motion.div
								initial={{ opacity: 0, height: 0 }}
								animate={{ opacity: 1, height: 'auto' }}
								transition={{ duration: 0.3 }}
								className='mt-2 space-y-2'
							>
								{upcomingMilestones.length > 0 ? (
									upcomingMilestones.map(
										(milestone, index) => (
											<div
												key={index}
												className='flex items-start gap-3 p-2 rounded-lg border'
											>
												{milestone.icon}
												<div>
													<p className='text-sm font-medium'>
														{milestone.name}
													</p>
													<p className='text-xs text-muted-foreground'>
														Day {milestone.day} (
														{milestone.day -
															daysPregnant}{' '}
														days remaining)
													</p>
												</div>
											</div>
										)
									)
								) : (
									<p className='text-sm text-muted-foreground'>
										No upcoming milestones
									</p>
								)}
							</motion.div>
						)}
					</div>

					{/* Completed milestones */}
					<div>
						<button
							className='flex justify-between items-center w-full py-2 text-left font-medium'
							onClick={() => toggleSection('completed')}
						>
							<span>Completed Milestones</span>
							<span className='text-muted-foreground text-sm'>
								{expandedSection === 'completed'
									? 'Hide'
									: 'Show'}
							</span>
						</button>

						{expandedSection === 'completed' && (
							<motion.div
								initial={{ opacity: 0, height: 0 }}
								animate={{ opacity: 1, height: 'auto' }}
								transition={{ duration: 0.3 }}
								className='mt-2 space-y-2'
							>
								{completedMilestones.length > 0 ? (
									completedMilestones.map(
										(milestone, index) => (
											<div
												key={index}
												className='flex items-start gap-3 p-2 rounded-lg border bg-muted/10'
											>
												<CheckCircle2 className='h-4 w-4 text-green-500' />
												<div>
													<p className='text-sm font-medium'>
														{milestone.name}
													</p>
													<p className='text-xs text-muted-foreground'>
														Completed on day{' '}
														{milestone.day} (
														{daysPregnant -
															milestone.day}{' '}
														days ago)
													</p>
												</div>
											</div>
										)
									)
								) : (
									<p className='text-sm text-muted-foreground'>
										No completed milestones yet
									</p>
								)}
							</motion.div>
						)}
					</div>
				</CardContent>
			</Card>
		</motion.div>
	);
};

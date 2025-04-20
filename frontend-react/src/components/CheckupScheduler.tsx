import { addDays, format, isBefore, isToday, parseISO } from 'date-fns';
import { motion } from 'framer-motion';
import { AlertCircle, CheckCircle2, Plus, Trash2 } from 'lucide-react';
import React, { useState } from 'react';
import { Button } from '../components/ui/button';
import {
	SelectAdapter,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from '../components/ui/select';
import { Horse } from '../types/horse';
import { Badge } from './ui/badge';
import {
	Card,
	CardContent,
	CardFooter,
	CardHeader,
	CardTitle,
} from './ui/card';
import { Input } from './ui/input';

interface CheckupSchedulerProps {
	horses: Horse[];
	onScheduleCheckup?: (
		horseId: string,
		date: string,
		notes: string,
		type: string
	) => void;
	onCancelCheckup?: (horseId: string, date: string) => void;
}

interface Checkup {
	id: string;
	horseId: string;
	horseName: string;
	date: string;
	notes: string;
	type: string;
	completed: boolean;
}

// Types of veterinary checkups
const CHECKUP_TYPES = [
	'Routine',
	'Pregnancy',
	'Vaccination',
	'Dental',
	'Hoof Care',
	'Emergency',
];

export const CheckupScheduler: React.FC<CheckupSchedulerProps> = ({
	horses,
	onScheduleCheckup = () => {},
	onCancelCheckup = () => {},
}) => {
	const [newCheckup, setNewCheckup] = useState<{
		horseId: string;
		date: string;
		notes: string;
		type: string;
	}>({
		horseId: '',
		date: format(new Date(), 'yyyy-MM-dd'),
		notes: '',
		type: 'Routine',
	});

	// Mock data for checkups (in a real app, this would come from API)
	const [checkups, setCheckups] = useState<Checkup[]>([
		{
			id: '1',
			horseId: horses[0]?.id || '',
			horseName: horses[0]?.name || 'Unknown',
			date: format(addDays(new Date(), 3), 'yyyy-MM-dd'),
			notes: 'Regular pregnancy check',
			type: 'Pregnancy',
			completed: false,
		},
		{
			id: '2',
			horseId: horses[1]?.id || '',
			horseName: horses[1]?.name || 'Unknown',
			date: format(addDays(new Date(), -2), 'yyyy-MM-dd'),
			notes: 'Hoof trimming',
			type: 'Hoof Care',
			completed: true,
		},
	]);

	const handleAddCheckup = () => {
		if (!newCheckup.horseId || !newCheckup.date) return;

		// Find horse name
		const horse = horses.find((h) => h.id === newCheckup.horseId);

		const checkup: Checkup = {
			id: Date.now().toString(),
			horseId: newCheckup.horseId,
			horseName: horse?.name || 'Unknown',
			date: newCheckup.date,
			notes: newCheckup.notes,
			type: newCheckup.type,
			completed: false,
		};

		setCheckups([...checkups, checkup]);
		onScheduleCheckup(
			checkup.horseId,
			checkup.date,
			checkup.notes,
			checkup.type
		);

		// Reset form
		setNewCheckup({
			horseId: '',
			date: format(new Date(), 'yyyy-MM-dd'),
			notes: '',
			type: 'Routine',
		});
	};

	const handleCancelCheckup = (id: string) => {
		const checkup = checkups.find((c) => c.id === id);
		if (checkup) {
			onCancelCheckup(checkup.horseId, checkup.date);
			setCheckups(checkups.filter((c) => c.id !== id));
		}
	};

	const handleToggleComplete = (id: string) => {
		setCheckups(
			checkups.map((c) => {
				if (c.id === id) {
					return { ...c, completed: !c.completed };
				}
				return c;
			})
		);
	};

	// Separate checkups into upcoming and past
	const today = new Date();
	const upcomingCheckups = checkups
		.filter(
			(c) =>
				!c.completed &&
				(isToday(parseISO(c.date)) || isBefore(today, parseISO(c.date)))
		)
		.sort(
			(a, b) => new Date(a.date).getTime() - new Date(b.date).getTime()
		);

	const pastCheckups = checkups
		.filter((c) => c.completed || isBefore(parseISO(c.date), today))
		.sort(
			(a, b) => new Date(b.date).getTime() - new Date(a.date).getTime()
		);

	// Check if there are urgent checkups (today or overdue)
	const urgentCheckups = checkups.filter(
		(c) =>
			!c.completed &&
			(isToday(parseISO(c.date)) || isBefore(parseISO(c.date), today))
	);

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
							Checkup Scheduler
						</CardTitle>
						{urgentCheckups.length > 0 && (
							<Badge variant='destructive'>
								{urgentCheckups.length} Urgent
							</Badge>
						)}
					</div>
				</CardHeader>

				<CardContent className='space-y-6'>
					{/* New Checkup Form */}
					<div className='space-y-4 rounded-lg border p-4'>
						<h3 className='font-semibold'>Schedule New Checkup</h3>

						<div className='space-y-2'>
							<label className='text-sm font-medium'>Horse</label>
							<SelectAdapter
								value={newCheckup.horseId}
								onValueChange={(value) =>
									setNewCheckup({
										...newCheckup,
										horseId: value,
									})
								}
							>
								<SelectTrigger className='w-full'>
									<SelectValue placeholder='Select horse' />
								</SelectTrigger>
								<SelectContent>
									{horses.map((horse) => (
										<SelectItem
											key={horse.id}
											value={horse.id}
										>
											{horse.name}
										</SelectItem>
									))}
								</SelectContent>
							</SelectAdapter>
						</div>

						<div className='grid grid-cols-1 md:grid-cols-2 gap-4'>
							<div className='space-y-2'>
								<label className='text-sm font-medium'>
									Date
								</label>
								<Input
									type='date'
									value={newCheckup.date}
									onChange={(e) =>
										setNewCheckup({
											...newCheckup,
											date: e.target.value,
										})
									}
								/>
							</div>

							<div className='space-y-2'>
								<label className='text-sm font-medium'>
									Type
								</label>
								<SelectAdapter
									value={newCheckup.type}
									onValueChange={(value) =>
										setNewCheckup({
											...newCheckup,
											type: value,
										})
									}
								>
									<SelectTrigger className='w-full'>
										<SelectValue placeholder='Select type' />
									</SelectTrigger>
									<SelectContent>
										{CHECKUP_TYPES.map((type) => (
											<SelectItem key={type} value={type}>
												{type}
											</SelectItem>
										))}
									</SelectContent>
								</SelectAdapter>
							</div>
						</div>

						<div className='space-y-2'>
							<label className='text-sm font-medium'>Notes</label>
							<Input
								value={newCheckup.notes}
								onChange={(e) =>
									setNewCheckup({
										...newCheckup,
										notes: e.target.value,
									})
								}
								placeholder='Add any special instructions or details'
							/>
						</div>

						<Button onClick={handleAddCheckup} className='w-full'>
							<Plus className='h-4 w-4 mr-2' />
							Schedule Checkup
						</Button>
					</div>

					{/* Upcoming Checkups */}
					<div>
						<h3 className='font-semibold mb-3'>
							Upcoming Checkups
						</h3>
						{upcomingCheckups.length > 0 ? (
							<div className='space-y-3'>
								{upcomingCheckups.map((checkup) => (
									<div
										key={checkup.id}
										className='flex items-start gap-3 p-3 rounded-lg border'
									>
										<div className='mt-1'>
											{isToday(parseISO(checkup.date)) ? (
												<Badge
													variant='destructive'
													className='h-8 w-8 rounded-full flex items-center justify-center p-0 text-xs'
												>
													Today
												</Badge>
											) : (
												<div className='h-8 w-8 rounded-full bg-primary/10 flex items-center justify-center text-xs font-medium'>
													{format(
														parseISO(checkup.date),
														'd'
													)}
												</div>
											)}
										</div>

										<div className='flex-1'>
											<div className='flex justify-between items-start'>
												<div>
													<p className='font-medium'>
														{checkup.horseName}
													</p>
													<p className='text-sm text-muted-foreground'>
														{format(
															parseISO(
																checkup.date
															),
															'MMMM d, yyyy'
														)}
													</p>
													<Badge
														variant='outline'
														className='mt-1'
													>
														{checkup.type}
													</Badge>
												</div>

												<div className='flex gap-2'>
													<Button
														variant='outline'
														size='sm'
														onClick={() =>
															handleToggleComplete(
																checkup.id
															)
														}
													>
														<CheckCircle2 className='h-4 w-4 mr-1' />
														Complete
													</Button>
													<Button
														variant='destructive'
														size='sm'
														onClick={() =>
															handleCancelCheckup(
																checkup.id
															)
														}
													>
														<Trash2 className='h-4 w-4' />
													</Button>
												</div>
											</div>

											{checkup.notes && (
												<p className='text-sm mt-2 text-muted-foreground'>
													{checkup.notes}
												</p>
											)}
										</div>
									</div>
								))}
							</div>
						) : (
							<div className='text-center py-6 text-muted-foreground'>
								No upcoming checkups scheduled
							</div>
						)}
					</div>

					{/* Past Checkups */}
					{pastCheckups.length > 0 && (
						<div>
							<h3 className='font-semibold mb-3'>
								Past Checkups
							</h3>
							<div className='space-y-3'>
								{pastCheckups.slice(0, 3).map((checkup) => (
									<div
										key={checkup.id}
										className='flex items-start gap-3 p-3 rounded-lg border bg-muted/10'
									>
										<div className='mt-1'>
											<div className='h-8 w-8 rounded-full bg-muted flex items-center justify-center text-xs font-medium text-muted-foreground'>
												{format(
													parseISO(checkup.date),
													'd'
												)}
											</div>
										</div>

										<div className='flex-1'>
											<div className='flex justify-between items-start'>
												<div>
													<p className='font-medium'>
														{checkup.horseName}
													</p>
													<p className='text-sm text-muted-foreground'>
														{format(
															parseISO(
																checkup.date
															),
															'MMMM d, yyyy'
														)}
													</p>
													<div className='flex items-center gap-2 mt-1'>
														<Badge
															variant='outline'
															className='text-muted-foreground'
														>
															{checkup.type}
														</Badge>
														{checkup.completed ? (
															<span className='text-xs flex items-center text-green-500'>
																<CheckCircle2 className='h-3 w-3 mr-1' />
																Completed
															</span>
														) : (
															<span className='text-xs flex items-center text-amber-500'>
																<AlertCircle className='h-3 w-3 mr-1' />
																Missed
															</span>
														)}
													</div>
												</div>
											</div>
										</div>
									</div>
								))}

								{pastCheckups.length > 3 && (
									<div className='text-center'>
										<Button variant='link' size='sm'>
											View all {pastCheckups.length} past
											checkups
										</Button>
									</div>
								)}
							</div>
						</div>
					)}
				</CardContent>

				<CardFooter className='justify-between text-sm text-muted-foreground'>
					<span>
						<span className='font-medium'>
							{upcomingCheckups.length}
						</span>{' '}
						upcoming
					</span>
					<span>
						<span className='font-medium'>
							{pastCheckups.length}
						</span>{' '}
						past checkups
					</span>
				</CardFooter>
			</Card>
		</motion.div>
	);
};

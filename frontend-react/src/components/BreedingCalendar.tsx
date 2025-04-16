import { format } from 'date-fns';
import { motion } from 'framer-motion';
import { AlertCircle, Calendar as CalendarIcon, Heart } from 'lucide-react';
import React, { useState } from 'react';
import { Horse } from '../types/horse';
import { Badge } from './ui/badge';
import { Calendar } from './ui/calendar';
import { Card, CardContent, CardHeader, CardTitle } from './ui/card';

interface BreedingCalendarProps {
	horses: Horse[];
}

export const BreedingCalendar: React.FC<BreedingCalendarProps> = ({
	horses,
}) => {
	const [selectedDate, setSelectedDate] = useState<Date>(new Date());

	// Calculate breeding and due dates for all horses
	const breedingEvents = horses.flatMap((horse) => {
		const events = [];
		if (horse.lastBreedingDate) {
			events.push({
				date: new Date(horse.lastBreedingDate),
				type: 'breeding',
				horseName: horse.name,
				horseId: horse.id,
			});
		}
		if (horse.dueDate) {
			events.push({
				date: new Date(horse.dueDate),
				type: 'due',
				horseName: horse.name,
				horseId: horse.id,
			});
		}
		return events;
	});

	// Custom date formatter to show events
	const dateFormatter = (date: Date) => {
		const events = breedingEvents.filter(
			(event) =>
				format(event.date, 'yyyy-MM-dd') === format(date, 'yyyy-MM-dd')
		);

		if (events.length === 0) return format(date, 'd');

		return (
			<div className='flex flex-col items-center gap-1'>
				<span className='text-sm font-medium'>{format(date, 'd')}</span>
				<div className='flex gap-1'>
					{events.map((event, index) => (
						<Badge
							key={`${event.horseId}-${event.type}-${index}`}
							variant={
								event.type === 'breeding'
									? 'default'
									: 'destructive'
							}
							className='text-xs'
						>
							{event.type === 'breeding' ? 'B' : 'D'}
						</Badge>
					))}
				</div>
			</div>
		);
	};

	// Get events for selected date
	const selectedDateEvents = breedingEvents.filter(
		(event) =>
			format(event.date, 'yyyy-MM-dd') ===
			format(selectedDate, 'yyyy-MM-dd')
	);

	return (
		<motion.div
			initial={{ opacity: 0, y: 20 }}
			animate={{ opacity: 1, y: 0 }}
			transition={{ duration: 0.5 }}
		>
			<Card className='w-full'>
				<CardHeader className='flex flex-row items-center justify-between space-y-0 pb-2'>
					<CardTitle className='text-2xl font-bold'>
						Breeding Calendar
					</CardTitle>
					<CalendarIcon className='h-6 w-6 text-muted-foreground' />
				</CardHeader>
				<CardContent>
					<div className='grid gap-6 md:grid-cols-2'>
						<div>
							<Calendar
								mode='single'
								selected={selectedDate}
								onSelect={setSelectedDate}
								className='rounded-md border'
								formatters={{
									formatDay: dateFormatter,
								}}
							/>
						</div>
						<div className='space-y-4'>
							<div className='flex items-center gap-2'>
								<Badge variant='default'>B</Badge>
								<span className='text-sm text-muted-foreground'>
									Breeding Date
								</span>
							</div>
							<div className='flex items-center gap-2'>
								<Badge variant='destructive'>D</Badge>
								<span className='text-sm text-muted-foreground'>
									Due Date
								</span>
							</div>

							{selectedDateEvents.length > 0 ? (
								<div className='mt-6 space-y-4'>
									<h3 className='text-lg font-semibold'>
										Events on{' '}
										{format(selectedDate, 'MMMM d, yyyy')}
									</h3>
									{selectedDateEvents.map((event, index) => (
										<motion.div
											key={`${event.horseId}-${event.type}-${index}`}
											initial={{ opacity: 0, x: -20 }}
											animate={{ opacity: 1, x: 0 }}
											transition={{ delay: index * 0.1 }}
											className='flex items-center gap-3 rounded-lg border p-3'
										>
											{event.type === 'breeding' ? (
												<Heart className='h-5 w-5 text-primary' />
											) : (
												<AlertCircle className='h-5 w-5 text-destructive' />
											)}
											<div>
												<p className='font-medium'>
													{event.horseName}
												</p>
												<p className='text-sm text-muted-foreground'>
													{event.type === 'breeding'
														? 'Breeding Date'
														: 'Due Date'}
												</p>
											</div>
										</motion.div>
									))}
								</div>
							) : (
								<div className='mt-6 text-center text-muted-foreground'>
									No events scheduled for this date
								</div>
							)}
						</div>
					</div>
				</CardContent>
			</Card>
		</motion.div>
	);
};

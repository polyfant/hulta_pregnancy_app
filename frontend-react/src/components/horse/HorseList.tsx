import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { cn } from '@/lib/utils';
import { Filter, Horse, Search } from 'lucide-react';
import { useState } from 'react';

interface Horse {
	id: string;
	name: string;
	breed: string;
	status: 'healthy' | 'warning' | 'critical';
	pregnancyWeek: number;
	dueDate: string;
	image: string;
}

export function HorseList() {
	const [searchQuery, setSearchQuery] = useState('');
	const [horses] = useState<Horse[]>([
		// Mock data - replace with actual data
		{
			id: '1',
			name: 'Stella',
			breed: 'Swedish Warmblood',
			status: 'healthy',
			pregnancyWeek: 24,
			dueDate: '2025-03-15',
			image: '/horses/stella.jpg',
		},
		// Add more horses...
	]);

	const filteredHorses = horses.filter(
		(horse) =>
			horse.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
			horse.breed.toLowerCase().includes(searchQuery.toLowerCase())
	);

	return (
		<div className='space-y-6'>
			<div className='flex items-center gap-4'>
				<div className='relative flex-1'>
					<Search className='absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted' />
					<Input
						placeholder='Search horses...'
						value={searchQuery}
						onChange={(e) => setSearchQuery(e.target.value)}
						className='pl-9'
					/>
				</div>
				<Button variant='outline' size='icon'>
					<Filter className='h-4 w-4' />
				</Button>
			</div>

			<div className='grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6'>
				{filteredHorses.map((horse) => (
					<Card
						key={horse.id}
						className='bg-background-DEFAULT hover:shadow-lg transition-all'
					>
						<CardHeader className='border-b border-primary/10'>
							<div className='flex items-center justify-between'>
								<div className='flex items-center gap-3'>
									<img
										src={horse.image}
										alt={horse.name}
										className='w-12 h-12 rounded-full object-cover border-2 border-primary'
									/>
									<div>
										<h3 className='text-lg font-semibold text-primary-dark'>
											{horse.name}
										</h3>
										<p className='text-sm text-muted'>
											{horse.breed}
										</p>
									</div>
								</div>
								<Badge
									className={cn('capitalize', {
										'bg-green-100 text-green-800':
											horse.status === 'healthy',
										'bg-yellow-100 text-yellow-800':
											horse.status === 'warning',
										'bg-red-100 text-red-800':
											horse.status === 'critical',
									})}
								>
									{horse.status}
								</Badge>
							</div>
						</CardHeader>
						<CardContent className='pt-4'>
							<div className='grid grid-cols-2 gap-4'>
								<div className='flex items-center gap-2'>
									<Horse className='h-4 w-4 text-primary' />
									<div>
										<p className='text-sm text-muted'>
											Pregnancy Stage
										</p>
										<p className='font-medium'>
											{horse.pregnancyWeek} weeks
										</p>
									</div>
								</div>
								<div className='flex items-center gap-2'>
									<CalendarDays className='h-4 w-4 text-primary' />
									<div>
										<p className='text-sm text-muted'>
											Due Date
										</p>
										<p className='font-medium'>
											{horse.dueDate}
										</p>
									</div>
								</div>
							</div>
						</CardContent>
					</Card>
				))}
			</div>
		</div>
	);
}

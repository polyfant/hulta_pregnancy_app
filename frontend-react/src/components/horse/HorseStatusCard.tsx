import { Badge } from '@/components/ui/badge';
import { Card, CardContent, CardHeader } from '@/components/ui/card';
import { cn } from '@/lib/utils';
import { CalendarDays, Horse } from 'lucide-react';

interface HorseStatusCardProps {
	horse: {
		name: string;
		breed: string;
		status: 'healthy' | 'warning' | 'critical';
		pregnancyWeek: number;
		dueDate: string;
		image: string;
	};
}

export function HorseStatusCard({ horse }: HorseStatusCardProps) {
	const statusColors = {
		healthy: 'bg-green-100 text-green-800',
		warning: 'bg-yellow-100 text-yellow-800',
		critical: 'bg-red-100 text-red-800',
	};

	return (
		<Card className='bg-background-DEFAULT hover:shadow-lg transition-all'>
			<CardHeader className='border-b border-primary/10'>
				<div className='flex items-center justify-between'>
					<div className='flex items-center gap-3'>
						<img
							src={horse.image}
							alt={horse.name}
							className='w-16 h-16 rounded-full object-cover border-2 border-primary'
						/>
						<div>
							<h3 className='text-xl font-semibold text-primary-dark'>
								{horse.name}
							</h3>
							<p className='text-sm text-muted'>{horse.breed}</p>
						</div>
					</div>
					<Badge
						className={cn('capitalize', statusColors[horse.status])}
					>
						{horse.status}
					</Badge>
				</div>
			</CardHeader>
			<CardContent className='pt-4'>
				<div className='grid grid-cols-2 gap-4'>
					<div className='flex items-center gap-2'>
						<Horse className='h-5 w-5 text-primary' />
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
						<CalendarDays className='h-5 w-5 text-primary' />
						<div>
							<p className='text-sm text-muted'>Due Date</p>
							<p className='font-medium'>{horse.dueDate}</p>
						</div>
					</div>
				</div>
			</CardContent>
		</Card>
	);
}

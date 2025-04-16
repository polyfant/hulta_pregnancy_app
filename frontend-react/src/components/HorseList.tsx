import { motion } from 'framer-motion';
import {
	Calendar,
	Filter,
	Heart,
	Info,
	Search,
	Trophy,
	User,
} from 'lucide-react';
import React, { useEffect, useMemo, useState } from 'react';
import { Button } from '../components/ui/button';
import {
	SelectAdapter,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from '../components/ui/select';
import { useToast } from '../components/ui/use-toast';
import { Horse } from '../types/horse';
import { StatsDashboard } from './StatsDashboard';
import { Badge } from './ui/badge';
import {
	Card,
	CardContent,
	CardFooter,
	CardHeader,
	CardTitle,
} from './ui/card';
import { Input } from './ui/input';

const HorseList: React.FC = () => {
	const [horses, setHorses] = useState<Horse[]>([]);
	const [loading, setLoading] = useState(true);
	const [error, setError] = useState<string | null>(null);
	const [searchTerm, setSearchTerm] = useState('');
	const [filterStatus, setFilterStatus] = useState<string>('all');
	const { toast } = useToast();

	// Calculate statistics
	const stats = useMemo(() => {
		const totalHorses = horses.length;
		const activeHorses = horses.filter((h) => h.status === 'active').length;
		const pregnantHorses = horses.filter((h) => h.isPregnant).length;
		const upcomingCheckups = horses.filter((h) => {
			if (!h.nextCheckup) return false;
			const nextCheckup = new Date(h.nextCheckup);
			const thirtyDaysFromNow = new Date();
			thirtyDaysFromNow.setDate(thirtyDaysFromNow.getDate() + 30);
			return nextCheckup <= thirtyDaysFromNow;
		}).length;

		const averageAge =
			horses.reduce((acc, horse) => acc + (horse.age || 0), 0) /
				totalHorses || 0;

		const healthStatus = {
			excellent: horses.filter((h) => h.healthStatus === 'excellent')
				.length,
			good: horses.filter((h) => h.healthStatus === 'good').length,
			fair: horses.filter((h) => h.healthStatus === 'fair').length,
			poor: horses.filter((h) => h.healthStatus === 'poor').length,
		};

		// ML Predictions (mock data for now)
		const mlPredictions = {
			nextPregnancyProbability: Math.round(Math.random() * 100),
			optimalBreedingTime: new Date(
				Date.now() + 30 * 24 * 60 * 60 * 1000
			).toLocaleDateString(),
			foalHealthPrediction: 'Excellent',
		};

		// Environmental Data (mock data for now)
		const environmentalData = {
			temperature: 22.5,
			humidity: 65,
			weatherCondition: 'Sunny',
			impactScore: 8,
		};

		return {
			totalHorses,
			activeHorses,
			pregnantHorses,
			upcomingCheckups,
			averageAge: Math.round(averageAge * 10) / 10,
			healthStatus,
			mlPredictions,
			environmentalData,
		};
	}, [horses]);

	useEffect(() => {
		const fetchHorses = async () => {
			try {
				const response = await fetch(
					'http://localhost:8080/api/horses'
				);
				if (!response.ok) {
					throw new Error('Failed to fetch horses');
				}
				const data = await response.json();
				setHorses(data);
			} catch (err) {
				setError(
					err instanceof Error ? err.message : 'An error occurred'
				);
				toast({
					title: 'Error',
					description: 'Failed to fetch horses',
					variant: 'destructive',
				});
			} finally {
				setLoading(false);
			}
		};

		fetchHorses();
	}, [toast]);

	const filteredHorses = horses.filter((horse) => {
		const matchesSearch =
			horse.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
			(horse.breed?.toLowerCase() || '').includes(
				searchTerm.toLowerCase()
			);
		const matchesStatus =
			filterStatus === 'all' || horse.status === filterStatus;
		return matchesSearch && matchesStatus;
	});

	if (loading) {
		return (
			<div className='flex items-center justify-center min-h-screen'>
				<div className='animate-spin rounded-full h-12 w-12 border-b-2 border-primary'></div>
			</div>
		);
	}

	if (error) {
		return (
			<div className='flex items-center justify-center min-h-screen'>
				<Card className='w-full max-w-md'>
					<CardHeader>
						<CardTitle>Error</CardTitle>
					</CardHeader>
					<CardContent>
						<p className='text-destructive'>{error}</p>
					</CardContent>
				</Card>
			</div>
		);
	}

	return (
		<div className='container mx-auto px-4 py-8'>
			<div className='mb-8'>
				<h1 className='text-3xl font-bold tracking-tight mb-2'>
					Horse Management
				</h1>
				<p className='text-muted-foreground'>
					Manage and monitor your horses with precision and care
				</p>
			</div>

			{/* Stats Dashboard */}
			<StatsDashboard {...stats} />

			<div className='flex flex-col md:flex-row gap-4 mb-8'>
				<div className='relative flex-1'>
					<Search className='absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground' />
					<Input
						type='text'
						placeholder='Search by name or breed...'
						value={searchTerm}
						onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
							setSearchTerm(e.target.value)
						}
						className='w-full pl-10'
					/>
				</div>
				<div className='w-full md:w-48'>
					<SelectAdapter
						value={filterStatus}
						onValueChange={setFilterStatus}
					>
						<SelectTrigger className='w-full'>
							<Filter className='h-4 w-4 mr-2' />
							<SelectValue placeholder='Filter by status' />
						</SelectTrigger>
						<SelectContent>
							<SelectItem value='all'>All Horses</SelectItem>
							<SelectItem value='active'>Active</SelectItem>
							<SelectItem value='retired'>Retired</SelectItem>
							<SelectItem value='sold'>Sold</SelectItem>
						</SelectContent>
					</SelectAdapter>
				</div>
			</div>

			<div className='grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6'>
				{filteredHorses.map((horse, index) => (
					<motion.div
						key={horse.id}
						initial={{ opacity: 0, y: 20 }}
						animate={{ opacity: 1, y: 0 }}
						transition={{ duration: 0.3, delay: index * 0.1 }}
					>
						<Card className='group hover:shadow-xl transition-all duration-300 border-border/50 hover:border-primary/50 relative overflow-hidden'>
							{/* Premium Ribbon */}
							{horse.isPremium && (
								<div className='absolute top-0 right-0 bg-gradient-to-r from-amber-500 to-yellow-400 text-white px-3 py-1 text-xs font-semibold transform rotate-45 translate-x-8 -translate-y-2'>
									Premium
								</div>
							)}

							<CardHeader className='pb-4'>
								<div className='flex items-start justify-between'>
									<div>
										<div className='flex items-center gap-2'>
											<CardTitle className='text-xl font-semibold tracking-tight group-hover:text-primary transition-colors'>
												{horse.name}
											</CardTitle>
											{horse.isChampion && (
												<Trophy className='h-4 w-4 text-amber-500' />
											)}
										</div>
										<p className='text-sm text-muted-foreground mt-1'>
											{horse.breed || 'Unknown breed'}
										</p>
									</div>
									<Badge
										variant={
											// Type-safe variant selection
											horse.status === 'active'
												? 'success'
												: horse.status === 'retired'
												? 'warning'
												: 'secondary'
										}
										className='capitalize'
									>
										{horse.status}
									</Badge>
								</div>
							</CardHeader>
							<CardContent className='pb-4'>
								<div className='space-y-3'>
									<div className='flex items-center gap-2 text-sm'>
										<User className='h-4 w-4 text-muted-foreground' />
										<span className='font-medium'>
											Gender:
										</span>
										<span className='text-muted-foreground'>
											{horse.gender}
										</span>
									</div>
									<div className='flex items-center gap-2 text-sm'>
										<Calendar className='h-4 w-4 text-muted-foreground' />
										<span className='font-medium'>
											Age:
										</span>
										<span className='text-muted-foreground'>
											{horse.age || 'Unknown'} years
										</span>
									</div>
									{horse.notes && (
										<div className='flex items-start gap-2 text-sm'>
											<Info className='h-4 w-4 text-muted-foreground mt-0.5' />
											<span className='text-muted-foreground'>
												{horse.notes}
											</span>
										</div>
									)}
								</div>
							</CardContent>
							<CardFooter className='flex justify-between items-center pt-4 border-t'>
								<Button
									variant='ghost'
									size='sm'
									className='hover:bg-primary/10 group/favorite'
								>
									<Heart className='h-4 w-4 mr-2 group-hover/favorite:text-red-500 transition-colors' />
									Favorite
								</Button>
								<div className='flex gap-2'>
									<Button
										variant='outline'
										size='sm'
										className='hover:bg-primary/10 hover:text-primary'
									>
										Edit
									</Button>
									<Button
										variant='destructive'
										size='sm'
										className='hover:bg-destructive/90'
									>
										Delete
									</Button>
								</div>
							</CardFooter>
						</Card>
					</motion.div>
				))}
			</div>
		</div>
	);
};

export default HorseList;

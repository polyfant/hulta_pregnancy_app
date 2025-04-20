import {
	Activity,
	Brain,
	Calendar,
	CloudSun,
	Droplets,
	Horse,
	Thermometer,
	TrendingUp,
} from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from './ui/card';
import { Progress } from './ui/progress';

interface StatsDashboardProps {
	totalHorses: number;
	activeHorses: number;
	pregnantHorses: number;
	upcomingCheckups: number;
	averageAge: number;
	healthStatus: {
		excellent: number;
		good: number;
		fair: number;
		poor: number;
	};
	mlPredictions?: {
		nextPregnancyProbability: number;
		optimalBreedingTime: string;
		foalHealthPrediction: string;
	};
	environmentalData?: {
		temperature: number;
		humidity: number;
		weatherCondition: string;
		impactScore: number;
	};
}

export function StatsDashboard({
	totalHorses,
	activeHorses,
	pregnantHorses,
	upcomingCheckups,
	averageAge,
	healthStatus,
	mlPredictions,
	environmentalData,
}: StatsDashboardProps) {
	const healthPercentage = {
		excellent: (healthStatus.excellent / totalHorses) * 100,
		good: (healthStatus.good / totalHorses) * 100,
		fair: (healthStatus.fair / totalHorses) * 100,
		poor: (healthStatus.poor / totalHorses) * 100,
	};

	return (
		<div className='space-y-6'>
			{/* Main Stats Row */}
			<div className='grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4'>
				{/* Total Horses Card */}
				<Card className='hover:shadow-lg transition-shadow'>
					<CardHeader className='flex flex-row items-center justify-between space-y-0 pb-2'>
						<CardTitle className='text-sm font-medium'>
							Total Horses
						</CardTitle>
						<Horse className='h-4 w-4 text-muted-foreground' />
					</CardHeader>
					<CardContent>
						<div className='text-2xl font-bold'>{totalHorses}</div>
						<p className='text-xs text-muted-foreground'>
							{activeHorses} active, {pregnantHorses} pregnant
						</p>
					</CardContent>
				</Card>

				{/* Health Overview Card */}
				<Card className='hover:shadow-lg transition-shadow'>
					<CardHeader className='flex flex-row items-center justify-between space-y-0 pb-2'>
						<CardTitle className='text-sm font-medium'>
							Health Overview
						</CardTitle>
						<Activity className='h-4 w-4 text-muted-foreground' />
					</CardHeader>
					<CardContent>
						<div className='space-y-2'>
							<div className='flex items-center justify-between'>
								<span className='text-sm'>Excellent</span>
								<span className='text-sm font-medium'>
									{healthStatus.excellent}
								</span>
							</div>
							<Progress
								value={healthPercentage.excellent}
								className='h-2'
							/>
							<div className='flex items-center justify-between'>
								<span className='text-sm'>Good</span>
								<span className='text-sm font-medium'>
									{healthStatus.good}
								</span>
							</div>
							<Progress
								value={healthPercentage.good}
								className='h-2'
							/>
						</div>
					</CardContent>
				</Card>

				{/* Upcoming Checkups Card */}
				<Card className='hover:shadow-lg transition-shadow'>
					<CardHeader className='flex flex-row items-center justify-between space-y-0 pb-2'>
						<CardTitle className='text-sm font-medium'>
							Upcoming Checkups
						</CardTitle>
						<Calendar className='h-4 w-4 text-muted-foreground' />
					</CardHeader>
					<CardContent>
						<div className='text-2xl font-bold'>
							{upcomingCheckups}
						</div>
						<p className='text-xs text-muted-foreground'>
							Next 30 days
						</p>
					</CardContent>
				</Card>

				{/* Average Age Card */}
				<Card className='hover:shadow-lg transition-shadow'>
					<CardHeader className='flex flex-row items-center justify-between space-y-0 pb-2'>
						<CardTitle className='text-sm font-medium'>
							Average Age
						</CardTitle>
						<TrendingUp className='h-4 w-4 text-muted-foreground' />
					</CardHeader>
					<CardContent>
						<div className='text-2xl font-bold'>
							{averageAge} years
						</div>
						<p className='text-xs text-muted-foreground'>
							Across all horses
						</p>
					</CardContent>
				</Card>
			</div>

			{/* ML Predictions Row */}
			{mlPredictions && (
				<div className='grid grid-cols-1 md:grid-cols-3 gap-4'>
					<Card className='hover:shadow-lg transition-shadow bg-gradient-to-br from-purple-50 to-purple-100'>
						<CardHeader className='flex flex-row items-center justify-between space-y-0 pb-2'>
							<CardTitle className='text-sm font-medium'>
								Next Pregnancy Probability
							</CardTitle>
							<Brain className='h-4 w-4 text-purple-600' />
						</CardHeader>
						<CardContent>
							<div className='text-2xl font-bold text-purple-600'>
								{mlPredictions.nextPregnancyProbability}%
							</div>
							<p className='text-xs text-muted-foreground'>
								Based on historical data and current conditions
							</p>
						</CardContent>
					</Card>

					<Card className='hover:shadow-lg transition-shadow bg-gradient-to-br from-blue-50 to-blue-100'>
						<CardHeader className='flex flex-row items-center justify-between space-y-0 pb-2'>
							<CardTitle className='text-sm font-medium'>
								Optimal Breeding Time
							</CardTitle>
							<Calendar className='h-4 w-4 text-blue-600' />
						</CardHeader>
						<CardContent>
							<div className='text-2xl font-bold text-blue-600'>
								{mlPredictions.optimalBreedingTime}
							</div>
							<p className='text-xs text-muted-foreground'>
								AI-recommended breeding window
							</p>
						</CardContent>
					</Card>

					<Card className='hover:shadow-lg transition-shadow bg-gradient-to-br from-green-50 to-green-100'>
						<CardHeader className='flex flex-row items-center justify-between space-y-0 pb-2'>
							<CardTitle className='text-sm font-medium'>
								Foal Health Prediction
							</CardTitle>
							<Activity className='h-4 w-4 text-green-600' />
						</CardHeader>
						<CardContent>
							<div className='text-2xl font-bold text-green-600'>
								{mlPredictions.foalHealthPrediction}
							</div>
							<p className='text-xs text-muted-foreground'>
								Predicted health status of next foal
							</p>
						</CardContent>
					</Card>
				</div>
			)}

			{/* Environmental Monitoring Row */}
			{environmentalData && (
				<div className='grid grid-cols-1 md:grid-cols-4 gap-4'>
					<Card className='hover:shadow-lg transition-shadow bg-gradient-to-br from-orange-50 to-orange-100'>
						<CardHeader className='flex flex-row items-center justify-between space-y-0 pb-2'>
							<CardTitle className='text-sm font-medium'>
								Temperature
							</CardTitle>
							<Thermometer className='h-4 w-4 text-orange-600' />
						</CardHeader>
						<CardContent>
							<div className='text-2xl font-bold text-orange-600'>
								{environmentalData.temperature}°C
							</div>
							<p className='text-xs text-muted-foreground'>
								Current stable temperature
							</p>
						</CardContent>
					</Card>

					<Card className='hover:shadow-lg transition-shadow bg-gradient-to-br from-blue-50 to-blue-100'>
						<CardHeader className='flex flex-row items-center justify-between space-y-0 pb-2'>
							<CardTitle className='text-sm font-medium'>
								Humidity
							</CardTitle>
							<Droplets className='h-4 w-4 text-blue-600' />
						</CardHeader>
						<CardContent>
							<div className='text-2xl font-bold text-blue-600'>
								{environmentalData.humidity}%
							</div>
							<p className='text-xs text-muted-foreground'>
								Current stable humidity
							</p>
						</CardContent>
					</Card>

					<Card className='hover:shadow-lg transition-shadow bg-gradient-to-br from-sky-50 to-sky-100'>
						<CardHeader className='flex flex-row items-center justify-between space-y-0 pb-2'>
							<CardTitle className='text-sm font-medium'>
								Weather
							</CardTitle>
							<CloudSun className='h-4 w-4 text-sky-600' />
						</CardHeader>
						<CardContent>
							<div className='text-2xl font-bold text-sky-600'>
								{environmentalData.weatherCondition}
							</div>
							<p className='text-xs text-muted-foreground'>
								Current weather conditions
							</p>
						</CardContent>
					</Card>

					<Card className='hover:shadow-lg transition-shadow bg-gradient-to-br from-emerald-50 to-emerald-100'>
						<CardHeader className='flex flex-row items-center justify-between space-y-0 pb-2'>
							<CardTitle className='text-sm font-medium'>
								Environmental Impact
							</CardTitle>
							<Activity className='h-4 w-4 text-emerald-600' />
						</CardHeader>
						<CardContent>
							<div className='text-2xl font-bold text-emerald-600'>
								{environmentalData.impactScore}/10
							</div>
							<p className='text-xs text-muted-foreground'>
								Current environmental impact score
							</p>
						</CardContent>
					</Card>
				</div>
			)}
		</div>
	);
}

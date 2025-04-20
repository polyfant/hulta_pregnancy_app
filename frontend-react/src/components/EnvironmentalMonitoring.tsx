import { motion } from 'framer-motion';
import {
	Cloud,
	CloudDrizzle,
	CloudSnow,
	Droplets,
	Leaf,
	RefreshCw,
	Sun,
	Wind,
} from 'lucide-react';
import React, { useEffect, useState } from 'react';
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
import { Skeleton } from './ui/skeleton';

interface WeatherData {
	temperature: number;
	humidity: number;
	condition: 'sunny' | 'cloudy' | 'rainy' | 'snowy' | 'windy';
	windSpeed: number;
	pressure: number;
	uvIndex: number;
	updatedAt: Date;
}

interface ForecastDay {
	day: string;
	condition: 'sunny' | 'cloudy' | 'rainy' | 'snowy' | 'windy';
	temperature: number;
	humidity: number;
	risk: 'low' | 'moderate' | 'high' | 'severe';
}

interface EnvironmentalImpact {
	current: {
		overall: number;
		thermal: number;
		humidity: number;
		stability: number;
	};
	recommendations: string[];
}

export const EnvironmentalMonitoring: React.FC = () => {
	const [loading, setLoading] = useState(true);
	const [weatherData, setWeatherData] = useState<WeatherData | null>(null);
	const [forecast, setForecast] = useState<ForecastDay[]>([]);
	const [impact, setImpact] = useState<EnvironmentalImpact | null>(null);

	// Mock fetching weather data
	useEffect(() => {
		const fetchWeatherData = async () => {
			setLoading(true);
			// Simulate API call delay
			await new Promise((resolve) => setTimeout(resolve, 1500));

			// Mock data that would normally come from a weather API
			const mockWeatherData: WeatherData = {
				temperature: 22.5,
				humidity: 65,
				condition: 'sunny',
				windSpeed: 8,
				pressure: 1015,
				uvIndex: 6,
				updatedAt: new Date(),
			};

			setWeatherData(mockWeatherData);

			// Generate mock forecast data
			const conditions: Array<
				'sunny' | 'cloudy' | 'rainy' | 'snowy' | 'windy'
			> = ['sunny', 'cloudy', 'rainy', 'snowy', 'windy'];
			const days = ['Today', 'Tomorrow', 'Day 3', 'Day 4', 'Day 5'];

			const mockForecast: ForecastDay[] = days.map((day) => {
				const randomTemp = Math.round(15 + Math.random() * 15);
				const randomHumidity = Math.round(40 + Math.random() * 50);
				const randomCondition =
					conditions[Math.floor(Math.random() * conditions.length)];

				// Calculate risk based on conditions
				let risk: 'low' | 'moderate' | 'high' | 'severe' = 'low';
				if (randomTemp > 28 && randomHumidity > 70) {
					risk = 'severe';
				} else if (randomTemp > 25 && randomHumidity > 60) {
					risk = 'high';
				} else if (randomTemp < 5 || randomHumidity > 85) {
					risk = 'moderate';
				}

				return {
					day,
					condition: randomCondition,
					temperature: randomTemp,
					humidity: randomHumidity,
					risk,
				};
			});

			setForecast(mockForecast);

			// Calculate environmental impact
			calculateEnvironmentalImpact(mockWeatherData);

			setLoading(false);
		};

		fetchWeatherData();
	}, []);

	// Calculate environmental impact based on weather data
	const calculateEnvironmentalImpact = (weather: WeatherData) => {
		// These calculations would normally be based on veterinary research
		// and specific to horse pregnancy, but we're using simplified formulas for the demo

		// Thermal impact: 0-100 score based on how close temperature is to ideal (15-20°C)
		const thermalImpact =
			100 - Math.min(100, Math.abs(weather.temperature - 17.5) * 5);

		// Humidity impact: 0-100 score (lower is better for high humidity)
		const humidityImpact = 100 - Math.max(0, weather.humidity - 50);

		// Stability score: 0-100 based on pressure and wind (higher is more stable)
		const stabilityImpact = Math.min(
			100,
			((weather.pressure - 990) / 40) * 100 - weather.windSpeed * 2
		);

		// Overall impact: weighted average
		const overallImpact =
			thermalImpact * 0.4 + humidityImpact * 0.3 + stabilityImpact * 0.3;

		// Generate recommendations based on impacts
		const recommendations: string[] = [];

		if (thermalImpact < 60) {
			if (weather.temperature > 25) {
				recommendations.push(
					'Provide extra shade and cooling for pregnant mares'
				);
				recommendations.push('Ensure constant access to fresh water');
			} else if (weather.temperature < 10) {
				recommendations.push('Provide warm shelter for pregnant mares');
				recommendations.push(
					'Consider adding extra bedding for insulation'
				);
			}
		}

		if (humidityImpact < 70) {
			recommendations.push(
				'Monitor for signs of heat stress in pregnant mares'
			);
			recommendations.push(
				'Consider reducing exercise during peak humidity hours'
			);
		}

		if (stabilityImpact < 60) {
			recommendations.push(
				'Ensure shelter is available during unstable weather'
			);
			recommendations.push(
				'Monitor for stress behaviors in pregnant mares'
			);
		}

		if (recommendations.length === 0) {
			recommendations.push(
				'Current conditions are generally favorable for pregnant mares'
			);
			recommendations.push('Maintain regular care and monitoring');
		}

		setImpact({
			current: {
				overall: Math.round(overallImpact),
				thermal: Math.round(thermalImpact),
				humidity: Math.round(humidityImpact),
				stability: Math.round(stabilityImpact),
			},
			recommendations,
		});
	};

	// Helper function to get weather icon based on condition
	const getWeatherIcon = (
		condition: 'sunny' | 'cloudy' | 'rainy' | 'snowy' | 'windy'
	) => {
		switch (condition) {
			case 'sunny':
				return <Sun className='h-6 w-6 text-amber-500' />;
			case 'cloudy':
				return <Cloud className='h-6 w-6 text-gray-500' />;
			case 'rainy':
				return <CloudDrizzle className='h-6 w-6 text-blue-500' />;
			case 'snowy':
				return <CloudSnow className='h-6 w-6 text-sky-500' />;
			case 'windy':
				return <Wind className='h-6 w-6 text-teal-500' />;
		}
	};

	// Get impact color based on score
	const getImpactColor = (score: number) => {
		if (score >= 80) return 'bg-green-500';
		if (score >= 60) return 'bg-lime-500';
		if (score >= 40) return 'bg-amber-500';
		if (score >= 20) return 'bg-orange-500';
		return 'bg-red-500';
	};

	// Handle refresh button click
	const handleRefresh = () => {
		setLoading(true);
		// This would normally trigger a new API call
		// For demo purposes, we'll just set a timeout and reuse the same data
		setTimeout(() => {
			if (weatherData) {
				setWeatherData({
					...weatherData,
					updatedAt: new Date(),
				});
			}
			setLoading(false);
		}, 1000);
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
								Environmental Monitoring
							</CardTitle>
							<CardDescription>
								Weather conditions and impact on pregnant horses
							</CardDescription>
						</div>
						<Button
							variant='outline'
							size='sm'
							onClick={handleRefresh}
							disabled={loading}
						>
							<RefreshCw
								className={`h-4 w-4 mr-2 ${
									loading ? 'animate-spin' : ''
								}`}
							/>
							Refresh
						</Button>
					</div>
				</CardHeader>

				<CardContent className='space-y-6'>
					{/* Current Weather */}
					<div className='grid grid-cols-1 md:grid-cols-2 gap-4'>
						<div className='flex flex-col p-4 rounded-lg border'>
							<h3 className='text-sm font-medium mb-2'>
								Current Conditions
							</h3>

							{loading ? (
								<div className='space-y-3'>
									<Skeleton className='h-12 w-full' />
									<Skeleton className='h-4 w-3/4' />
									<Skeleton className='h-4 w-1/2' />
								</div>
							) : weatherData ? (
								<div className='space-y-3'>
									<div className='flex items-center justify-between'>
										<div className='flex items-center gap-3'>
											{getWeatherIcon(
												weatherData.condition
											)}
											<div>
												<p className='text-2xl font-bold'>
													{weatherData.temperature}°C
												</p>
												<p className='text-xs text-muted-foreground'>
													Feels like{' '}
													{Math.round(
														weatherData.temperature +
															(weatherData.humidity >
															70
																? 2
																: -1)
													)}
													°C
												</p>
											</div>
										</div>
										<div className='text-right'>
											<p className='text-muted-foreground text-sm'>
												Last updated
											</p>
											<p className='text-xs'>
												{weatherData.updatedAt.toLocaleTimeString()}
											</p>
										</div>
									</div>

									<div className='grid grid-cols-3 gap-2 mt-2'>
										<div className='flex items-center'>
											<Droplets className='h-4 w-4 text-blue-500 mr-1' />
											<span className='text-sm'>
												{weatherData.humidity}%
											</span>
										</div>
										<div className='flex items-center'>
											<Wind className='h-4 w-4 text-teal-500 mr-1' />
											<span className='text-sm'>
												{weatherData.windSpeed} km/h
											</span>
										</div>
										<div className='flex items-center'>
											<Sun className='h-4 w-4 text-amber-500 mr-1' />
											<span className='text-sm'>
												UV {weatherData.uvIndex}
											</span>
										</div>
									</div>
								</div>
							) : null}
						</div>

						{/* Environmental Impact */}
						<div className='flex flex-col p-4 rounded-lg border'>
							<h3 className='text-sm font-medium mb-2'>
								Pregnancy Impact Score
							</h3>

							{loading ? (
								<div className='space-y-3'>
									<Skeleton className='h-6 w-full' />
									<Skeleton className='h-4 w-3/4' />
									<Skeleton className='h-8 w-full' />
								</div>
							) : impact ? (
								<div className='space-y-3'>
									<div className='flex items-center justify-between'>
										<div className='flex items-center gap-2'>
											<Leaf className='h-5 w-5 text-green-500' />
											<p className='font-medium'>
												Overall:{' '}
												{impact.current.overall}/100
											</p>
										</div>
										<Badge
											variant={
												impact.current.overall >= 60
													? 'success'
													: impact.current.overall >=
													  40
													? 'warning'
													: 'destructive'
											}
										>
											{impact.current.overall >= 80
												? 'Excellent'
												: impact.current.overall >= 60
												? 'Good'
												: impact.current.overall >= 40
												? 'Fair'
												: impact.current.overall >= 20
												? 'Poor'
												: 'Critical'}
										</Badge>
									</div>

									<div className='space-y-2'>
										<div>
											<div className='flex justify-between text-xs mb-1'>
												<span>Thermal Comfort</span>
												<span>
													{impact.current.thermal}%
												</span>
											</div>
											<Progress
												value={impact.current.thermal}
												className={`h-1.5 ${getImpactColor(
													impact.current.thermal
												)}`}
											/>
										</div>
										<div>
											<div className='flex justify-between text-xs mb-1'>
												<span>Humidity Impact</span>
												<span>
													{impact.current.humidity}%
												</span>
											</div>
											<Progress
												value={impact.current.humidity}
												className={`h-1.5 ${getImpactColor(
													impact.current.humidity
												)}`}
											/>
										</div>
										<div>
											<div className='flex justify-between text-xs mb-1'>
												<span>Weather Stability</span>
												<span>
													{impact.current.stability}%
												</span>
											</div>
											<Progress
												value={impact.current.stability}
												className={`h-1.5 ${getImpactColor(
													impact.current.stability
												)}`}
											/>
										</div>
									</div>
								</div>
							) : null}
						</div>
					</div>

					{/* Recommendations */}
					{!loading && impact && (
						<div className='space-y-2 p-4 rounded-lg border'>
							<h3 className='text-sm font-medium'>
								Recommendations
							</h3>
							<ul className='space-y-1'>
								{impact.recommendations.map((rec, index) => (
									<li key={index} className='flex text-sm'>
										<span className='text-green-500 mr-2'>
											•
										</span>
										<span>{rec}</span>
									</li>
								))}
							</ul>
						</div>
					)}

					{/* 5-Day Forecast */}
					<div>
						<h3 className='text-sm font-medium mb-3'>
							5-Day Forecast & Impact Prediction
						</h3>

						{loading ? (
							<div className='grid grid-cols-5 gap-2'>
								{[...Array(5)].map((_, i) => (
									<Skeleton key={i} className='h-24 w-full' />
								))}
							</div>
						) : (
							<div className='grid grid-cols-2 sm:grid-cols-3 md:grid-cols-5 gap-2'>
								{forecast.map((day, index) => (
									<div
										key={index}
										className='flex flex-col items-center p-3 rounded-lg border'
									>
										<p className='text-xs font-medium'>
											{day.day}
										</p>
										{getWeatherIcon(day.condition)}
										<p className='text-sm font-medium mt-1'>
											{day.temperature}°C
										</p>
										<div className='flex items-center mt-1'>
											<Droplets className='h-3 w-3 text-blue-500 mr-1' />
											<span className='text-xs'>
												{day.humidity}%
											</span>
										</div>
										<Badge
											variant={
												day.risk === 'low'
													? 'success'
													: day.risk === 'moderate'
													? 'outline'
													: day.risk === 'high'
													? 'warning'
													: 'destructive'
											}
											className='mt-2 text-xs'
										>
											{day.risk} risk
										</Badge>
									</div>
								))}
							</div>
						)}
					</div>
				</CardContent>

				<CardFooter className='border-t p-4 text-xs text-muted-foreground'>
					<div className='w-full flex justify-between items-center'>
						<p>Data source: Local Hulta Farm Weather Station</p>
						<Button variant='link' size='sm' className='h-auto p-0'>
							View detailed forecast
						</Button>
					</div>
				</CardFooter>
			</Card>
		</motion.div>
	);
};

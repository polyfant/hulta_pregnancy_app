import { motion } from 'framer-motion';
import {
	Activity,
	AlertCircle,
	Brain,
	Calendar,
	ChevronUp,
	Clock,
	Dna,
	LineChart,
	ThermometerSnowflake,
	Zap,
} from 'lucide-react';
import React, { useState } from 'react';
import { Horse } from '../types/horse';
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
import { Tabs, TabsContent, TabsList, TabsTrigger } from './ui/tabs';
import {
	Tooltip,
	TooltipContent,
	TooltipProvider,
	TooltipTrigger,
} from './ui/tooltip';

interface MLPredictionsProps {
	horse: Horse;
	environmentalData?: {
		temperature: number;
		humidity: number;
		weatherCondition: string;
	};
}

export const MLPredictions: React.FC<MLPredictionsProps> = ({
	horse,
	environmentalData = {
		temperature: 22.5,
		humidity: 65,
		weatherCondition: 'Sunny',
	},
}) => {
	const [activeTab, setActiveTab] = useState<string>('growth');
	const [showDetails, setShowDetails] = useState<boolean>(false);

	// Mock ML predictions data
	// In a real application, this would come from an API call to a ML model
	const predictions = {
		growth: {
			weightPrediction: horse.isPregnant ? 'Above average' : 'Normal',
			growthRate: horse.isPregnant ? 8.2 : 5.5,
			nextCheckupRecommendation: '2 weeks',
			nutritionSuggestion: horse.isPregnant
				? 'Increased protein intake recommended'
				: 'Maintain current diet',
			confidenceScore: 87,
			environmentalImpact: determineEnvironmentalImpact(
				environmentalData.temperature,
				environmentalData.humidity
			),
		},
		health: {
			vitalSignsPrediction: 'Normal',
			riskFactors: horse.age > 15 ? ['Age-related concerns'] : [],
			recommendedTests:
				horse.age > 15
					? ['Full blood panel', 'Joint mobility assessment']
					: ['Routine checkup'],
			confidenceScore: 92,
			longevityEstimate:
				horse.age > 0
					? Math.round(30 - horse.age + Math.random() * 5)
					: 'Unknown',
		},
		breeding: {
			fertilityScore: horse.gender === 'Female' ? 84 : 91,
			optimalBreedingTime: getRandomFutureDate(30, 60),
			geneticCompatibilityFactors: [
				'Good temperament',
				'Strong conformation',
			],
			predictedOffspringTraits: [
				'Athletic ability',
				'Good health disposition',
			],
			confidenceScore: 78,
		},
	};

	// Helper function to determine environmental impact based on conditions
	function determineEnvironmentalImpact(
		temperature: number,
		humidity: number
	): string {
		if (temperature > 30 && humidity > 70) {
			return 'High stress conditions - monitor closely';
		} else if (temperature < 5) {
			return 'Cold stress - ensure adequate shelter';
		} else if (temperature > 25 && humidity > 60) {
			return 'Mild heat stress - increase water availability';
		} else {
			return 'Optimal conditions - maintain regular care';
		}
	}

	// Helper function to get a random future date within a range
	function getRandomFutureDate(minDays: number, maxDays: number): string {
		const today = new Date();
		const randomDays =
			Math.floor(Math.random() * (maxDays - minDays + 1)) + minDays;
		const futureDate = new Date(
			today.setDate(today.getDate() + randomDays)
		);
		return futureDate.toLocaleDateString();
	}

	// Get confidence level description based on score
	function getConfidenceLevel(score: number): string {
		if (score >= 90) return 'Very High';
		if (score >= 80) return 'High';
		if (score >= 70) return 'Good';
		if (score >= 60) return 'Moderate';
		return 'Low';
	}

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
								ML Predictions
								<Badge
									variant='outline'
									className='ml-2 bg-purple-50 text-purple-700 border-purple-200'
								>
									<Brain className='h-3 w-3 mr-1' />
									AI
								</Badge>
							</CardTitle>
							<CardDescription>
								Machine learning insights and predictions
							</CardDescription>
						</div>
						<TooltipProvider>
							<Tooltip>
								<TooltipTrigger asChild>
									<div className='flex items-center rounded-full bg-purple-100 px-2 py-1 text-xs font-medium text-purple-700'>
										<Activity className='h-3 w-3 mr-1' />
										Model v2.3.1
									</div>
								</TooltipTrigger>
								<TooltipContent>
									<p>
										Using latest ML model (updated 14 days
										ago)
									</p>
								</TooltipContent>
							</Tooltip>
						</TooltipProvider>
					</div>
				</CardHeader>

				<CardContent className='pb-2'>
					<Tabs
						defaultValue='growth'
						className='w-full'
						onValueChange={setActiveTab}
					>
						<TabsList className='grid w-full grid-cols-3'>
							<TabsTrigger value='growth'>
								<LineChart className='h-4 w-4 mr-2' />
								Growth
							</TabsTrigger>
							<TabsTrigger value='health'>
								<Activity className='h-4 w-4 mr-2' />
								Health
							</TabsTrigger>
							<TabsTrigger value='breeding'>
								<Dna className='h-4 w-4 mr-2' />
								Breeding
							</TabsTrigger>
						</TabsList>

						{/* Growth Predictions */}
						<TabsContent value='growth' className='space-y-4 mt-4'>
							<div className='grid grid-cols-1 md:grid-cols-2 gap-4'>
								<div className='space-y-2'>
									<div className='flex justify-between items-center'>
										<div className='flex items-center'>
											<p className='text-sm font-medium'>
												Weight Trajectory
											</p>
											<Badge
												variant='outline'
												className='ml-2'
											>
												{
													predictions.growth
														.weightPrediction
												}
											</Badge>
										</div>
										<span className='text-xs text-muted-foreground'>
											{predictions.growth.growthRate}%
											monthly
										</span>
									</div>
									<Progress
										value={
											predictions.growth.growthRate * 10
										}
										className='h-2'
									/>
								</div>

								<div className='space-y-2'>
									<div className='flex justify-between items-center'>
										<div className='flex items-center'>
											<p className='text-sm font-medium'>
												Next Checkup
											</p>
										</div>
										<span className='text-xs text-muted-foreground'>
											<Clock className='h-3 w-3 inline mr-1' />
											{
												predictions.growth
													.nextCheckupRecommendation
											}
										</span>
									</div>
									<div className='p-2 rounded-md bg-blue-50 border border-blue-100 text-xs text-blue-700'>
										{predictions.growth.nutritionSuggestion}
									</div>
								</div>
							</div>

							<div className='flex items-center gap-2 p-3 rounded-lg border mt-2'>
								<ThermometerSnowflake className='h-5 w-5 text-blue-500 flex-shrink-0' />
								<div>
									<p className='text-sm font-medium'>
										Environmental Impact
									</p>
									<p className='text-xs text-muted-foreground'>
										{predictions.growth.environmentalImpact}
									</p>
								</div>
							</div>

							<div className='flex items-center justify-between'>
								<div className='flex items-center text-sm'>
									<AlertCircle className='h-4 w-4 mr-1 text-muted-foreground' />
									<span className='text-muted-foreground'>
										Prediction confidence:
									</span>
									<span className='font-medium ml-1'>
										{getConfidenceLevel(
											predictions.growth.confidenceScore
										)}
									</span>
								</div>
								<Badge variant='outline'>
									{predictions.growth.confidenceScore}%
								</Badge>
							</div>
						</TabsContent>

						{/* Health Predictions */}
						<TabsContent value='health' className='space-y-4 mt-4'>
							<div className='grid grid-cols-1 md:grid-cols-2 gap-4'>
								<div className='p-3 rounded-lg border'>
									<div className='flex items-center gap-2'>
										<Activity className='h-5 w-5 text-green-500' />
										<div>
											<p className='text-sm font-medium'>
												Vital Signs Prediction
											</p>
											<p className='text-xs text-muted-foreground'>
												{
													predictions.health
														.vitalSignsPrediction
												}
											</p>
										</div>
									</div>
								</div>

								<div className='p-3 rounded-lg border'>
									<div className='flex items-center gap-2'>
										<Calendar className='h-5 w-5 text-purple-500' />
										<div>
											<p className='text-sm font-medium'>
												Estimated Longevity
											</p>
											<p className='text-xs text-muted-foreground'>
												~
												{
													predictions.health
														.longevityEstimate
												}{' '}
												years
											</p>
										</div>
									</div>
								</div>
							</div>

							<div className='space-y-2'>
								<p className='text-sm font-medium'>
									Recommended Health Tests
								</p>
								<div className='flex flex-wrap gap-2'>
									{predictions.health.recommendedTests.map(
										(test, index) => (
											<Badge
												key={index}
												variant='outline'
												className='bg-green-50 text-green-700 border-green-200'
											>
												{test}
											</Badge>
										)
									)}
								</div>
							</div>

							{predictions.health.riskFactors.length > 0 && (
								<div className='space-y-2'>
									<p className='text-sm font-medium'>
										Risk Factors
									</p>
									<div className='flex flex-wrap gap-2'>
										{predictions.health.riskFactors.map(
											(factor, index) => (
												<Badge
													key={index}
													variant='outline'
													className='bg-amber-50 text-amber-700 border-amber-200'
												>
													{factor}
												</Badge>
											)
										)}
									</div>
								</div>
							)}

							<div className='flex items-center justify-between'>
								<div className='flex items-center text-sm'>
									<AlertCircle className='h-4 w-4 mr-1 text-muted-foreground' />
									<span className='text-muted-foreground'>
										Prediction confidence:
									</span>
									<span className='font-medium ml-1'>
										{getConfidenceLevel(
											predictions.health.confidenceScore
										)}
									</span>
								</div>
								<Badge variant='outline'>
									{predictions.health.confidenceScore}%
								</Badge>
							</div>
						</TabsContent>

						{/* Breeding Predictions */}
						<TabsContent
							value='breeding'
							className='space-y-4 mt-4'
						>
							<div className='grid grid-cols-1 md:grid-cols-2 gap-4'>
								<div className='space-y-2'>
									<p className='text-sm font-medium'>
										Fertility Score
									</p>
									<div className='flex items-center'>
										<Progress
											value={
												predictions.breeding
													.fertilityScore
											}
											className='h-2 flex-1'
										/>
										<span className='ml-2 text-sm font-medium'>
											{
												predictions.breeding
													.fertilityScore
											}
											%
										</span>
									</div>
								</div>

								<div className='p-3 rounded-lg border'>
									<div className='flex items-center gap-2'>
										<Calendar className='h-5 w-5 text-blue-500' />
										<div>
											<p className='text-sm font-medium'>
												Optimal Breeding Time
											</p>
											<p className='text-xs text-muted-foreground'>
												{
													predictions.breeding
														.optimalBreedingTime
												}
											</p>
										</div>
									</div>
								</div>
							</div>

							<div className='space-y-2'>
								<p className='text-sm font-medium'>
									Genetic Compatibility Factors
								</p>
								<div className='flex flex-wrap gap-2'>
									{predictions.breeding.geneticCompatibilityFactors.map(
										(factor, index) => (
											<Badge
												key={index}
												variant='outline'
												className='bg-purple-50 text-purple-700 border-purple-200'
											>
												{factor}
											</Badge>
										)
									)}
								</div>
							</div>

							<div className='space-y-2'>
								<p className='text-sm font-medium'>
									Predicted Offspring Traits
								</p>
								<div className='flex flex-wrap gap-2'>
									{predictions.breeding.predictedOffspringTraits.map(
										(trait, index) => (
											<Badge
												key={index}
												variant='outline'
												className='bg-blue-50 text-blue-700 border-blue-200'
											>
												{trait}
											</Badge>
										)
									)}
								</div>
							</div>

							<div className='flex items-center justify-between'>
								<div className='flex items-center text-sm'>
									<AlertCircle className='h-4 w-4 mr-1 text-muted-foreground' />
									<span className='text-muted-foreground'>
										Prediction confidence:
									</span>
									<span className='font-medium ml-1'>
										{getConfidenceLevel(
											predictions.breeding.confidenceScore
										)}
									</span>
								</div>
								<Badge variant='outline'>
									{predictions.breeding.confidenceScore}%
								</Badge>
							</div>
						</TabsContent>
					</Tabs>
				</CardContent>

				<CardFooter className='pt-1 flex flex-col'>
					<Button
						variant='link'
						size='sm'
						className='ml-auto mb-2'
						onClick={() => setShowDetails(!showDetails)}
					>
						{showDetails ? 'Hide details' : 'Show model details'}
						<ChevronUp
							className={`ml-1 h-4 w-4 transition-transform ${
								showDetails ? 'rotate-0' : 'rotate-180'
							}`}
						/>
					</Button>

					{showDetails && (
						<motion.div
							initial={{ opacity: 0, height: 0 }}
							animate={{ opacity: 1, height: 'auto' }}
							transition={{ duration: 0.3 }}
							className='w-full text-xs text-muted-foreground space-y-2 border-t pt-2'
						>
							<p>
								<span className='font-medium'>ML Model:</span>{' '}
								TensorFlow.js Neural Network
							</p>
							<p>
								<span className='font-medium'>
									Training Data:
								</span>{' '}
								Based on 15,000+ horses with similar
								characteristics
							</p>
							<p>
								<span className='font-medium'>
									Last Updated:
								</span>{' '}
								14 days ago
							</p>
							<p>
								<span className='font-medium'>
									Key Factors:
								</span>{' '}
								Age, breed, gender, medical history,
								environmental conditions
							</p>
							<div className='flex items-center mt-2'>
								<Zap className='h-4 w-4 text-amber-500 mr-1' />
								<p className='text-amber-600'>
									Predictions are for guidance only and not a
									substitute for veterinary care
								</p>
							</div>
						</motion.div>
					)}
				</CardFooter>
			</Card>
		</motion.div>
	);
};

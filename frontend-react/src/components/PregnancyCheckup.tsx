import { AlertCircle, Check, Heart, Scale, Thermometer } from 'lucide-react';
import React, { useState } from 'react';
import { toast } from 'sonner';
import { Horse } from '../types/horse';
import { Button } from './ui/button';
import { Calendar } from './ui/calendar';
import { Card, CardContent, CardHeader, CardTitle } from './ui/card';
import { Input } from './ui/input';
import { Label } from './ui/label';
import { Textarea } from './ui/textarea';

interface PregnancyCheckupProps {
	horse: Horse;
	onCheckupComplete: (checkup: CheckupData) => void;
}

interface CheckupData {
	date: Date;
	temperature: number;
	heartRate: number;
	weight: number;
	notes: string;
	nextCheckupDate: Date;
}

export const PregnancyCheckup: React.FC<PregnancyCheckupProps> = ({
	horse,
	onCheckupComplete,
}) => {
	const [checkupData, setCheckupData] = useState<CheckupData>({
		date: new Date(),
		temperature: 0,
		heartRate: 0,
		weight: 0,
		notes: '',
		nextCheckupDate: new Date(
			new Date().setDate(new Date().getDate() + 30)
		),
	});

	const handleSubmit = (e: React.FormEvent) => {
		e.preventDefault();

		// Validate data
		if (checkupData.temperature < 37 || checkupData.temperature > 38.5) {
			toast.error('Temperature must be between 37°C and 38.5°C');
			return;
		}

		if (checkupData.heartRate < 30 || checkupData.heartRate > 50) {
			toast.error('Heart rate must be between 30 and 50 BPM');
			return;
		}

		if (checkupData.weight < 400 || checkupData.weight > 1000) {
			toast.error('Weight must be between 400kg and 1000kg');
			return;
		}

		onCheckupComplete(checkupData);
		toast.success('Checkup recorded successfully! 🎉');
	};

	return (
		<Card className='w-full'>
			<CardHeader>
				<CardTitle className='flex items-center gap-2'>
					<Heart className='h-5 w-5 text-pink-500' />
					Pregnancy Checkup for {horse.name}
				</CardTitle>
			</CardHeader>
			<CardContent>
				<form onSubmit={handleSubmit} className='space-y-6'>
					<div className='grid grid-cols-1 md:grid-cols-2 gap-4'>
						<div className='space-y-2'>
							<Label htmlFor='date'>Checkup Date</Label>
							<Calendar
								mode='single'
								selected={checkupData.date}
								onSelect={(date) =>
									date &&
									setCheckupData({ ...checkupData, date })
								}
								className='rounded-md border'
							/>
						</div>

						<div className='space-y-2'>
							<Label htmlFor='nextCheckup'>
								Next Checkup Date
							</Label>
							<Calendar
								mode='single'
								selected={checkupData.nextCheckupDate}
								onSelect={(date) =>
									date &&
									setCheckupData({
										...checkupData,
										nextCheckupDate: date,
									})
								}
								className='rounded-md border'
							/>
						</div>
					</div>

					<div className='grid grid-cols-1 md:grid-cols-3 gap-4'>
						<div className='space-y-2'>
							<Label
								htmlFor='temperature'
								className='flex items-center gap-2'
							>
								<Thermometer className='h-4 w-4' />
								Temperature (°C)
							</Label>
							<Input
								id='temperature'
								type='number'
								step='0.1'
								value={checkupData.temperature}
								onChange={(e) =>
									setCheckupData({
										...checkupData,
										temperature: parseFloat(e.target.value),
									})
								}
								placeholder='37.5'
							/>
						</div>

						<div className='space-y-2'>
							<Label
								htmlFor='heartRate'
								className='flex items-center gap-2'
							>
								<Heart className='h-4 w-4 text-red-500' />
								Heart Rate (BPM)
							</Label>
							<Input
								id='heartRate'
								type='number'
								value={checkupData.heartRate}
								onChange={(e) =>
									setCheckupData({
										...checkupData,
										heartRate: parseInt(e.target.value),
									})
								}
								placeholder='40'
							/>
						</div>

						<div className='space-y-2'>
							<Label
								htmlFor='weight'
								className='flex items-center gap-2'
							>
								<Scale className='h-4 w-4' />
								Weight (kg)
							</Label>
							<Input
								id='weight'
								type='number'
								value={checkupData.weight}
								onChange={(e) =>
									setCheckupData({
										...checkupData,
										weight: parseInt(e.target.value),
									})
								}
								placeholder='600'
							/>
						</div>
					</div>

					<div className='space-y-2'>
						<Label
							htmlFor='notes'
							className='flex items-center gap-2'
						>
							<AlertCircle className='h-4 w-4' />
							Notes
						</Label>
						<Textarea
							id='notes'
							value={checkupData.notes}
							onChange={(e) =>
								setCheckupData({
									...checkupData,
									notes: e.target.value,
								})
							}
							placeholder='Enter any observations or concerns...'
							className='min-h-[100px]'
						/>
					</div>

					<div className='flex justify-end'>
						<Button type='submit' className='gap-2'>
							<Check className='h-4 w-4' />
							Record Checkup
						</Button>
					</div>
				</form>
			</CardContent>
		</Card>
	);
};

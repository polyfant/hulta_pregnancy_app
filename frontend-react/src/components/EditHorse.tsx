import { zodResolver } from '@hookform/resolvers/zod';
import { motion } from 'framer-motion';
import {
	Calendar,
	Check,
	Heart,
	HorseshoeIcon,
	Save,
	User,
	X,
} from 'lucide-react';
import React, { useEffect, useState } from 'react';
import { useForm } from 'react-hook-form';
import { z } from 'zod';
import { Horse } from '../types/horse';
import { Badge } from './ui/badge';
import { Button } from './ui/button';
import {
	Card,
	CardContent,
	CardFooter,
	CardHeader,
	CardTitle,
} from './ui/card';
import {
	Form,
	FormControl,
	FormDescription,
	FormField,
	FormItem,
	FormLabel,
	FormMessage,
} from './ui/form';
import { Input } from './ui/input';
import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from './ui/select';
import { Switch } from './ui/switch';
import { Tabs, TabsContent, TabsList, TabsTrigger } from './ui/tabs';
import { Textarea } from './ui/textarea';
import { useToast } from './ui/use-toast';

// Form schema using Zod for validation
const formSchema = z.object({
	name: z.string().min(2, 'Name must be at least 2 characters').max(50),
	breed: z.string().min(2, 'Breed must be at least 2 characters').optional(),
	gender: z.enum(['Male', 'Female', 'Gelding']),
	birthDate: z.string().optional(),
	color: z.string().optional(),
	height: z.string().optional(),
	weight: z.string().optional(),
	registrationNumber: z.string().optional(),
	microchipNumber: z.string().optional(),
	status: z.enum(['active', 'retired', 'sold', 'deceased']),
	isPregnant: z.boolean().default(false),
	lastBreedingDate: z.string().optional(),
	healthStatus: z.enum(['excellent', 'good', 'fair', 'poor']).optional(),
	notes: z.string().optional(),
	sire: z.string().optional(),
	dam: z.string().optional(),
	isPremium: z.boolean().default(false),
	isChampion: z.boolean().default(false),
});

type FormValues = z.infer<typeof formSchema>;

interface EditHorseProps {
	horse: Horse;
	onSave: (updatedHorse: Horse) => Promise<void>;
	onCancel: () => void;
}

export const EditHorse: React.FC<EditHorseProps> = ({
	horse,
	onSave,
	onCancel,
}) => {
	const [activeTab, setActiveTab] = useState('details');
	const [isSubmitting, setIsSubmitting] = useState(false);
	const [hasChanges, setHasChanges] = useState(false);
	const { toast } = useToast();

	// Set up form with default values from the horse prop
	const form = useForm<FormValues>({
		resolver: zodResolver(formSchema),
		defaultValues: {
			name: horse.name,
			breed: horse.breed || '',
			gender: horse.gender as 'Male' | 'Female' | 'Gelding',
			birthDate: horse.birthDate || '',
			color: horse.color || '',
			height: horse.height?.toString() || '',
			weight: horse.weight?.toString() || '',
			registrationNumber: horse.registrationNumber || '',
			microchipNumber: horse.microchipNumber || '',
			status: horse.status as 'active' | 'retired' | 'sold' | 'deceased',
			isPregnant: horse.isPregnant || false,
			lastBreedingDate: horse.lastBreedingDate || '',
			healthStatus:
				(horse.healthStatus as
					| 'excellent'
					| 'good'
					| 'fair'
					| 'poor') || 'good',
			notes: horse.notes || '',
			sire: horse.sire || '',
			dam: horse.dam || '',
			isPremium: horse.isPremium || false,
			isChampion: horse.isChampion || false,
		},
	});

	// Watch for form changes to enable/disable the save button
	useEffect(() => {
		const subscription = form.watch(() => {
			setHasChanges(true);
		});
		return () => subscription.unsubscribe();
	}, [form.watch]);

	const onSubmit = async (data: FormValues) => {
		setIsSubmitting(true);
		try {
			// Create updated horse object
			const updatedHorse: Horse = {
				...horse,
				...data,
				height: data.height ? parseFloat(data.height) : undefined,
				weight: data.weight ? parseFloat(data.weight) : undefined,
			};

			await onSave(updatedHorse);
			toast({
				title: 'Horse updated',
				description: `${data.name} has been successfully updated.`,
				variant: 'success',
			});
			setHasChanges(false);
		} catch (error) {
			toast({
				title: 'Error updating horse',
				description:
					'There was a problem updating this horse. Please try again.',
				variant: 'destructive',
			});
			console.error('Error updating horse:', error);
		} finally {
			setIsSubmitting(false);
		}
	};

	const confirmCancel = () => {
		if (hasChanges) {
			// Simple confirmation - in a real app you might want a dialog component
			if (
				window.confirm(
					'You have unsaved changes. Are you sure you want to cancel?'
				)
			) {
				onCancel();
			}
		} else {
			onCancel();
		}
	};

	return (
		<motion.div
			initial={{ opacity: 0, y: 20 }}
			animate={{ opacity: 1, y: 0 }}
			transition={{ duration: 0.3 }}
			className='w-full max-w-4xl mx-auto'
		>
			<Card>
				<CardHeader>
					<CardTitle className='text-2xl flex items-center gap-2'>
						<HorseshoeIcon className='h-6 w-6 text-primary' />
						Edit Horse: {horse.name}
						{horse.isPremium && (
							<Badge
								variant='outline'
								className='ml-2 bg-amber-50 text-amber-700 border-amber-200'
							>
								Premium
							</Badge>
						)}
						{horse.isChampion && (
							<Badge
								variant='outline'
								className='ml-1 bg-purple-50 text-purple-700 border-purple-200'
							>
								Champion
							</Badge>
						)}
					</CardTitle>
				</CardHeader>
				<CardContent>
					<Tabs value={activeTab} onValueChange={setActiveTab}>
						<TabsList className='grid grid-cols-3 mb-6'>
							<TabsTrigger value='details'>
								<User className='h-4 w-4 mr-2' />
								Basic Details
							</TabsTrigger>
							<TabsTrigger value='health'>
								<Heart className='h-4 w-4 mr-2' />
								Health & Breeding
							</TabsTrigger>
							<TabsTrigger value='additional'>
								<Calendar className='h-4 w-4 mr-2' />
								Additional Info
							</TabsTrigger>
						</TabsList>

						<Form {...form}>
							<form
								onSubmit={form.handleSubmit(onSubmit)}
								className='space-y-6'
							>
								<TabsContent
									value='details'
									className='space-y-4'
								>
									<div className='grid grid-cols-1 md:grid-cols-2 gap-4'>
										<FormField
											control={form.control}
											name='name'
											render={({ field }) => (
												<FormItem>
													<FormLabel>
														Horse Name*
													</FormLabel>
													<FormControl>
														<Input
															placeholder='Enter horse name'
															{...field}
														/>
													</FormControl>
													<FormMessage />
												</FormItem>
											)}
										/>

										<FormField
											control={form.control}
											name='breed'
											render={({ field }) => (
												<FormItem>
													<FormLabel>Breed</FormLabel>
													<FormControl>
														<Input
															placeholder='Enter breed'
															{...field}
														/>
													</FormControl>
													<FormMessage />
												</FormItem>
											)}
										/>

										<FormField
											control={form.control}
											name='gender'
											render={({ field }) => (
												<FormItem>
													<FormLabel>
														Gender*
													</FormLabel>
													<Select
														onValueChange={
															field.onChange
														}
														defaultValue={
															field.value
														}
													>
														<FormControl>
															<SelectTrigger>
																<SelectValue placeholder='Select gender' />
															</SelectTrigger>
														</FormControl>
														<SelectContent>
															<SelectItem value='Male'>
																Male
															</SelectItem>
															<SelectItem value='Female'>
																Female
															</SelectItem>
															<SelectItem value='Gelding'>
																Gelding
															</SelectItem>
														</SelectContent>
													</Select>
													<FormMessage />
												</FormItem>
											)}
										/>

										<FormField
											control={form.control}
											name='birthDate'
											render={({ field }) => (
												<FormItem>
													<FormLabel>
														Birth Date
													</FormLabel>
													<FormControl>
														<Input
															type='date'
															{...field}
														/>
													</FormControl>
													<FormMessage />
												</FormItem>
											)}
										/>

										<FormField
											control={form.control}
											name='color'
											render={({ field }) => (
												<FormItem>
													<FormLabel>Color</FormLabel>
													<FormControl>
														<Input
															placeholder='Enter color'
															{...field}
														/>
													</FormControl>
													<FormMessage />
												</FormItem>
											)}
										/>

										<FormField
											control={form.control}
											name='status'
											render={({ field }) => (
												<FormItem>
													<FormLabel>
														Status*
													</FormLabel>
													<Select
														onValueChange={
															field.onChange
														}
														defaultValue={
															field.value
														}
													>
														<FormControl>
															<SelectTrigger>
																<SelectValue placeholder='Select status' />
															</SelectTrigger>
														</FormControl>
														<SelectContent>
															<SelectItem value='active'>
																Active
															</SelectItem>
															<SelectItem value='retired'>
																Retired
															</SelectItem>
															<SelectItem value='sold'>
																Sold
															</SelectItem>
															<SelectItem value='deceased'>
																Deceased
															</SelectItem>
														</SelectContent>
													</Select>
													<FormMessage />
												</FormItem>
											)}
										/>
									</div>
								</TabsContent>

								<TabsContent
									value='health'
									className='space-y-4'
								>
									<div className='grid grid-cols-1 md:grid-cols-2 gap-4'>
										<FormField
											control={form.control}
											name='height'
											render={({ field }) => (
												<FormItem>
													<FormLabel>
														Height (cm)
													</FormLabel>
													<FormControl>
														<Input
															type='number'
															placeholder='Enter height'
															{...field}
														/>
													</FormControl>
													<FormMessage />
												</FormItem>
											)}
										/>

										<FormField
											control={form.control}
											name='weight'
											render={({ field }) => (
												<FormItem>
													<FormLabel>
														Weight (kg)
													</FormLabel>
													<FormControl>
														<Input
															type='number'
															placeholder='Enter weight'
															{...field}
														/>
													</FormControl>
													<FormMessage />
												</FormItem>
											)}
										/>

										<FormField
											control={form.control}
											name='healthStatus'
											render={({ field }) => (
												<FormItem>
													<FormLabel>
														Health Status
													</FormLabel>
													<Select
														onValueChange={
															field.onChange
														}
														defaultValue={
															field.value
														}
													>
														<FormControl>
															<SelectTrigger>
																<SelectValue placeholder='Select health status' />
															</SelectTrigger>
														</FormControl>
														<SelectContent>
															<SelectItem value='excellent'>
																Excellent
															</SelectItem>
															<SelectItem value='good'>
																Good
															</SelectItem>
															<SelectItem value='fair'>
																Fair
															</SelectItem>
															<SelectItem value='poor'>
																Poor
															</SelectItem>
														</SelectContent>
													</Select>
													<FormMessage />
												</FormItem>
											)}
										/>

										{form.watch('gender') === 'Female' && (
											<>
												<FormField
													control={form.control}
													name='isPregnant'
													render={({ field }) => (
														<FormItem className='flex flex-row items-center justify-between rounded-lg border p-3 shadow-sm'>
															<div className='space-y-0.5'>
																<FormLabel>
																	Pregnant
																</FormLabel>
																<FormDescription>
																	Is this mare
																	currently
																	pregnant?
																</FormDescription>
															</div>
															<FormControl>
																<Switch
																	checked={
																		field.value
																	}
																	onCheckedChange={
																		field.onChange
																	}
																/>
															</FormControl>
														</FormItem>
													)}
												/>

												<FormField
													control={form.control}
													name='lastBreedingDate'
													render={({ field }) => (
														<FormItem>
															<FormLabel>
																Last Breeding
																Date
															</FormLabel>
															<FormControl>
																<Input
																	type='date'
																	{...field}
																/>
															</FormControl>
															<FormMessage />
														</FormItem>
													)}
												/>
											</>
										)}
									</div>
								</TabsContent>

								<TabsContent
									value='additional'
									className='space-y-4'
								>
									<div className='grid grid-cols-1 md:grid-cols-2 gap-4'>
										<FormField
											control={form.control}
											name='registrationNumber'
											render={({ field }) => (
												<FormItem>
													<FormLabel>
														Registration Number
													</FormLabel>
													<FormControl>
														<Input
															placeholder='Enter registration number'
															{...field}
														/>
													</FormControl>
													<FormMessage />
												</FormItem>
											)}
										/>

										<FormField
											control={form.control}
											name='microchipNumber'
											render={({ field }) => (
												<FormItem>
													<FormLabel>
														Microchip Number
													</FormLabel>
													<FormControl>
														<Input
															placeholder='Enter microchip number'
															{...field}
														/>
													</FormControl>
													<FormMessage />
												</FormItem>
											)}
										/>

										<FormField
											control={form.control}
											name='sire'
											render={({ field }) => (
												<FormItem>
													<FormLabel>
														Sire (Father)
													</FormLabel>
													<FormControl>
														<Input
															placeholder="Enter sire's name"
															{...field}
														/>
													</FormControl>
													<FormMessage />
												</FormItem>
											)}
										/>

										<FormField
											control={form.control}
											name='dam'
											render={({ field }) => (
												<FormItem>
													<FormLabel>
														Dam (Mother)
													</FormLabel>
													<FormControl>
														<Input
															placeholder="Enter dam's name"
															{...field}
														/>
													</FormControl>
													<FormMessage />
												</FormItem>
											)}
										/>

										<FormField
											control={form.control}
											name='isPremium'
											render={({ field }) => (
												<FormItem className='flex flex-row items-center justify-between rounded-lg border p-3 shadow-sm'>
													<div className='space-y-0.5'>
														<FormLabel>
															Premium Horse
														</FormLabel>
														<FormDescription>
															Mark this horse as
															premium
														</FormDescription>
													</div>
													<FormControl>
														<Switch
															checked={
																field.value
															}
															onCheckedChange={
																field.onChange
															}
														/>
													</FormControl>
												</FormItem>
											)}
										/>

										<FormField
											control={form.control}
											name='isChampion'
											render={({ field }) => (
												<FormItem className='flex flex-row items-center justify-between rounded-lg border p-3 shadow-sm'>
													<div className='space-y-0.5'>
														<FormLabel>
															Champion Horse
														</FormLabel>
														<FormDescription>
															Mark this horse as a
															champion
														</FormDescription>
													</div>
													<FormControl>
														<Switch
															checked={
																field.value
															}
															onCheckedChange={
																field.onChange
															}
														/>
													</FormControl>
												</FormItem>
											)}
										/>
									</div>

									<FormField
										control={form.control}
										name='notes'
										render={({ field }) => (
											<FormItem>
												<FormLabel>Notes</FormLabel>
												<FormControl>
													<Textarea
														placeholder='Enter any additional notes'
														className='min-h-32'
														{...field}
													/>
												</FormControl>
												<FormMessage />
											</FormItem>
										)}
									/>
								</TabsContent>

								<div className='flex justify-end space-x-2 pt-4 border-t'>
									<Button
										type='button'
										variant='outline'
										onClick={confirmCancel}
										disabled={isSubmitting}
									>
										<X className='h-4 w-4 mr-2' />
										Cancel
									</Button>
									<Button
										type='submit'
										disabled={isSubmitting || !hasChanges}
										className='gap-2'
									>
										{isSubmitting ? (
											<div className='h-4 w-4 border-2 border-current border-t-transparent animate-spin rounded-full'></div>
										) : (
											<Save className='h-4 w-4' />
										)}
										Save Changes
									</Button>
								</div>
							</form>
						</Form>
					</Tabs>
				</CardContent>
				<CardFooter className='text-sm text-muted-foreground bg-muted/50'>
					<div className='flex items-center'>
						<Check className='h-4 w-4 mr-2 text-green-500' />
						Last edited: {new Date().toLocaleString()}
					</div>
				</CardFooter>
			</Card>
		</motion.div>
	);
};

export default EditHorse;

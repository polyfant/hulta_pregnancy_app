import { zodResolver } from '@hookform/resolvers/zod';
import { format } from 'date-fns';
import { motion } from 'framer-motion';
import {
	Calendar,
	Check,
	Heart,
	Horse as HorseIcon,
	Info,
	User,
} from 'lucide-react';
import React, { useState } from 'react';
import { useForm } from 'react-hook-form';
import { z } from 'zod';
import { cn } from '../lib/utils';
import { Button } from './ui/button';
import { Calendar as CalendarComponent } from './ui/calendar';
import {
	Card,
	CardContent,
	CardDescription,
	CardFooter,
	CardHeader,
	CardTitle,
} from './ui/card';
import { Checkbox } from './ui/checkbox';
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
import { Popover, PopoverContent, PopoverTrigger } from './ui/popover';
import { RadioGroup, RadioGroupItem } from './ui/radio-group';
import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from './ui/select';
import { Separator } from './ui/separator';
import { Textarea } from './ui/textarea';
import { useToast } from './ui/use-toast';

const formSchema = z.object({
	name: z
		.string()
		.min(2, {
			message: 'Horse name must be at least 2 characters.',
		})
		.max(50),
	breed: z
		.string()
		.min(2, {
			message: 'Breed must be at least 2 characters.',
		})
		.optional(),
	gender: z.enum(['Male', 'Female', 'Gelding'], {
		required_error: 'Please select a gender.',
	}),
	birthDate: z.date({
		required_error: 'Birth date is required.',
	}),
	color: z.string().optional(),
	height: z.coerce.number().positive().optional(),
	weight: z.coerce.number().positive().optional(),
	registrationNumber: z.string().optional(),
	microchipNumber: z.string().optional(),
	status: z.enum(['active', 'retired', 'sold'], {
		required_error: 'Please select a status.',
	}),
	isPregnant: z.boolean().default(false),
	lastBreedingDate: z.date().optional(),
	healthStatus: z
		.enum(['excellent', 'good', 'fair', 'poor'], {
			required_error: 'Please select a health status.',
		})
		.default('good'),
	notes: z.string().optional(),
	sire: z.string().optional(),
	dam: z.string().optional(),
	isPremium: z.boolean().default(false),
	isChampion: z.boolean().default(false),
});

type FormValues = z.infer<typeof formSchema>;

export const AddHorse: React.FC = () => {
	const [step, setStep] = useState(1);
	const [isSubmitting, setIsSubmitting] = useState(false);
	const { toast } = useToast();

	const form = useForm<FormValues>({
		resolver: zodResolver(formSchema),
		defaultValues: {
			name: '',
			breed: '',
			gender: undefined,
			color: '',
			status: 'active',
			isPregnant: false,
			healthStatus: 'good',
			notes: '',
			isPremium: false,
			isChampion: false,
		},
	});

	const { watch, setValue } = form;
	const selectedGender = watch('gender');
	const isPregnant = watch('isPregnant');

	// Function to go to next step
	const nextStep = () => {
		setStep(step + 1);
	};

	// Function to go to previous step
	const prevStep = () => {
		setStep(step - 1);
	};

	const onSubmit = async (data: FormValues) => {
		setIsSubmitting(true);

		try {
			// This would be an API call in a real application
			await new Promise((resolve) => setTimeout(resolve, 1500));

			// Convert the data to match the Horse type expected by the API
			const horseData = {
				...data,
				id: Date.now().toString(),
				age: new Date().getFullYear() - data.birthDate.getFullYear(),
				dueDate:
					data.isPregnant && data.lastBreedingDate
						? new Date(
								data.lastBreedingDate.getTime() +
									340 * 24 * 60 * 60 * 1000
						  )
						: undefined,
			};

			console.log('Submitting horse data:', horseData);

			// Show success notification
			toast({
				title: 'Horse Added Successfully',
				description: `${data.name} has been added to your stable.`,
			});

			// Reset form
			form.reset();
			setStep(1);
		} catch (error) {
			console.error('Error adding horse:', error);

			// Show error notification
			toast({
				title: 'Error Adding Horse',
				description:
					'There was a problem adding the horse. Please try again.',
				variant: 'destructive',
			});
		} finally {
			setIsSubmitting(false);
		}
	};

	return (
		<motion.div
			initial={{ opacity: 0, y: 20 }}
			animate={{ opacity: 1, y: 0 }}
			transition={{ duration: 0.5 }}
		>
			<Card className='w-full'>
				<CardHeader>
					<CardTitle className='text-xl font-bold'>
						Add New Horse
					</CardTitle>
					<CardDescription>
						Enter details to add a new horse to your stable
					</CardDescription>
				</CardHeader>

				<CardContent>
					<Form {...form}>
						<form
							onSubmit={form.handleSubmit(onSubmit)}
							className='space-y-6'
						>
							{/* Step 1: Basic Information */}
							{step === 1 && (
								<motion.div
									initial={{ opacity: 0 }}
									animate={{ opacity: 1 }}
									transition={{ duration: 0.3 }}
								>
									<div className='space-y-4'>
										<div className='flex items-center gap-2 text-lg font-semibold'>
											<HorseIcon className='h-5 w-5' />
											<h3>Basic Information</h3>
										</div>

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
														<FormLabel>
															Breed
														</FormLabel>
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
										</div>

										<div className='grid grid-cols-1 md:grid-cols-2 gap-4'>
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
													<FormItem className='flex flex-col'>
														<FormLabel>
															Birth Date*
														</FormLabel>
														<Popover>
															<PopoverTrigger
																asChild
															>
																<FormControl>
																	<Button
																		variant={
																			'outline'
																		}
																		className={cn(
																			'w-full pl-3 text-left font-normal',
																			!field.value &&
																				'text-muted-foreground'
																		)}
																	>
																		{field.value ? (
																			format(
																				field.value,
																				'PPP'
																			)
																		) : (
																			<span>
																				Pick
																				a
																				date
																			</span>
																		)}
																		<Calendar className='ml-auto h-4 w-4 opacity-50' />
																	</Button>
																</FormControl>
															</PopoverTrigger>
															<PopoverContent
																className='w-auto p-0'
																align='start'
															>
																<CalendarComponent
																	mode='single'
																	selected={
																		field.value
																	}
																	onSelect={
																		field.onChange
																	}
																	disabled={(
																		date
																	) =>
																		date >
																			new Date() ||
																		date <
																			new Date(
																				'1990-01-01'
																			)
																	}
																	initialFocus
																/>
															</PopoverContent>
														</Popover>
														<FormMessage />
													</FormItem>
												)}
											/>
										</div>

										<div className='grid grid-cols-1 md:grid-cols-2 gap-4'>
											<FormField
												control={form.control}
												name='color'
												render={({ field }) => (
													<FormItem>
														<FormLabel>
															Color
														</FormLabel>
														<FormControl>
															<Input
																placeholder='e.g., Bay, Chestnut, etc.'
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
															</SelectContent>
														</Select>
														<FormMessage />
													</FormItem>
												)}
											/>
										</div>

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
																min='0'
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
																min='0'
																{...field}
															/>
														</FormControl>
														<FormMessage />
													</FormItem>
												)}
											/>
										</div>
									</div>

									<div className='flex justify-end mt-6'>
										<Button
											type='button'
											onClick={nextStep}
										>
											Next Step
										</Button>
									</div>
								</motion.div>
							)}

							{/* Step 2: Health & Breeding Information */}
							{step === 2 && (
								<motion.div
									initial={{ opacity: 0 }}
									animate={{ opacity: 1 }}
									transition={{ duration: 0.3 }}
								>
									<div className='space-y-4'>
										<div className='flex items-center gap-2 text-lg font-semibold'>
											<Heart className='h-5 w-5' />
											<h3>Health & Breeding</h3>
										</div>

										<FormField
											control={form.control}
											name='healthStatus'
											render={({ field }) => (
												<FormItem className='space-y-3'>
													<FormLabel>
														Health Status*
													</FormLabel>
													<FormControl>
														<RadioGroup
															onValueChange={
																field.onChange
															}
															defaultValue={
																field.value
															}
															className='flex space-x-4'
														>
															<FormItem className='flex items-center space-x-2 space-y-0'>
																<FormControl>
																	<RadioGroupItem value='excellent' />
																</FormControl>
																<FormLabel className='font-normal'>
																	Excellent
																</FormLabel>
															</FormItem>
															<FormItem className='flex items-center space-x-2 space-y-0'>
																<FormControl>
																	<RadioGroupItem value='good' />
																</FormControl>
																<FormLabel className='font-normal'>
																	Good
																</FormLabel>
															</FormItem>
															<FormItem className='flex items-center space-x-2 space-y-0'>
																<FormControl>
																	<RadioGroupItem value='fair' />
																</FormControl>
																<FormLabel className='font-normal'>
																	Fair
																</FormLabel>
															</FormItem>
															<FormItem className='flex items-center space-x-2 space-y-0'>
																<FormControl>
																	<RadioGroupItem value='poor' />
																</FormControl>
																<FormLabel className='font-normal'>
																	Poor
																</FormLabel>
															</FormItem>
														</RadioGroup>
													</FormControl>
													<FormMessage />
												</FormItem>
											)}
										/>

										<Separator />

										{selectedGender === 'Female' && (
											<>
												<FormField
													control={form.control}
													name='isPregnant'
													render={({ field }) => (
														<FormItem className='flex flex-row items-start space-x-3 space-y-0 rounded-md border p-4'>
															<FormControl>
																<Checkbox
																	checked={
																		field.value
																	}
																	onCheckedChange={
																		field.onChange
																	}
																/>
															</FormControl>
															<div className='space-y-1 leading-none'>
																<FormLabel>
																	Currently
																	Pregnant
																</FormLabel>
																<FormDescription>
																	Indicate if
																	this mare is
																	currently
																	pregnant
																</FormDescription>
															</div>
														</FormItem>
													)}
												/>

												{isPregnant && (
													<FormField
														control={form.control}
														name='lastBreedingDate'
														render={({ field }) => (
															<FormItem className='flex flex-col'>
																<FormLabel>
																	Breeding
																	Date*
																</FormLabel>
																<Popover>
																	<PopoverTrigger
																		asChild
																	>
																		<FormControl>
																			<Button
																				variant={
																					'outline'
																				}
																				className={cn(
																					'w-full pl-3 text-left font-normal',
																					!field.value &&
																						'text-muted-foreground'
																				)}
																			>
																				{field.value ? (
																					format(
																						field.value,
																						'PPP'
																					)
																				) : (
																					<span>
																						Pick
																						a
																						date
																					</span>
																				)}
																				<Calendar className='ml-auto h-4 w-4 opacity-50' />
																			</Button>
																		</FormControl>
																	</PopoverTrigger>
																	<PopoverContent
																		className='w-auto p-0'
																		align='start'
																	>
																		<CalendarComponent
																			mode='single'
																			selected={
																				field.value
																			}
																			onSelect={
																				field.onChange
																			}
																			disabled={(
																				date
																			) =>
																				date >
																					new Date() ||
																				date <
																					new Date(
																						Date.now() -
																							365 *
																								24 *
																								60 *
																								60 *
																								1000
																					)
																			}
																			initialFocus
																		/>
																	</PopoverContent>
																</Popover>
																<FormDescription>
																	This will be
																	used to
																	calculate
																	the expected
																	due date
																</FormDescription>
																<FormMessage />
															</FormItem>
														)}
													/>
												)}

												<Separator />
											</>
										)}

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
															<Input {...field} />
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
															<Input {...field} />
														</FormControl>
														<FormMessage />
													</FormItem>
												)}
											/>
										</div>

										<FormField
											control={form.control}
											name='notes'
											render={({ field }) => (
												<FormItem>
													<FormLabel>
														Health Notes
													</FormLabel>
													<FormControl>
														<Textarea
															placeholder='Enter any health conditions, medications, or special care instructions'
															className='min-h-[100px]'
															{...field}
														/>
													</FormControl>
													<FormMessage />
												</FormItem>
											)}
										/>
									</div>

									<div className='flex justify-between mt-6'>
										<Button
											type='button'
											variant='outline'
											onClick={prevStep}
										>
											Previous Step
										</Button>
										<Button
											type='button'
											onClick={nextStep}
										>
											Next Step
										</Button>
									</div>
								</motion.div>
							)}

							{/* Step 3: Lineage & Extra Information */}
							{step === 3 && (
								<motion.div
									initial={{ opacity: 0 }}
									animate={{ opacity: 1 }}
									transition={{ duration: 0.3 }}
								>
									<div className='space-y-4'>
										<div className='flex items-center gap-2 text-lg font-semibold'>
											<User className='h-5 w-5' />
											<h3>
												Lineage & Additional Details
											</h3>
										</div>

										<div className='grid grid-cols-1 md:grid-cols-2 gap-4'>
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
										</div>

										<Separator />

										<div className='space-y-4'>
											<FormField
												control={form.control}
												name='isPremium'
												render={({ field }) => (
													<FormItem className='flex flex-row items-start space-x-3 space-y-0 rounded-md border p-4'>
														<FormControl>
															<Checkbox
																checked={
																	field.value
																}
																onCheckedChange={
																	field.onChange
																}
															/>
														</FormControl>
														<div className='space-y-1 leading-none'>
															<FormLabel>
																Premium Horse
															</FormLabel>
															<FormDescription>
																Mark this horse
																as premium
																quality
															</FormDescription>
														</div>
													</FormItem>
												)}
											/>

											<FormField
												control={form.control}
												name='isChampion'
												render={({ field }) => (
													<FormItem className='flex flex-row items-start space-x-3 space-y-0 rounded-md border p-4'>
														<FormControl>
															<Checkbox
																checked={
																	field.value
																}
																onCheckedChange={
																	field.onChange
																}
															/>
														</FormControl>
														<div className='space-y-1 leading-none'>
															<FormLabel>
																Champion Status
															</FormLabel>
															<FormDescription>
																Indicate if this
																horse has
																champion status
																or competition
																wins
															</FormDescription>
														</div>
													</FormItem>
												)}
											/>
										</div>

										<div className='rounded-md border p-4 bg-muted/20'>
											<div className='flex items-center gap-2'>
												<Info className='h-5 w-5 text-blue-500' />
												<p className='text-sm font-medium'>
													Completing Horse Profile
												</p>
											</div>
											<p className='text-sm text-muted-foreground mt-2'>
												After adding the horse, you can
												further customize the profile by
												adding photos, health records,
												competition history, and
												breeding records from the
												horse's detail page.
											</p>
										</div>
									</div>

									<div className='flex justify-between mt-6'>
										<Button
											type='button'
											variant='outline'
											onClick={prevStep}
										>
											Previous Step
										</Button>
										<Button
											type='submit'
											disabled={isSubmitting}
										>
											{isSubmitting ? (
												<>
													<span className='mr-2'>
														Saving...
													</span>
												</>
											) : (
												<>
													<Check className='mr-2 h-4 w-4' />
													Add Horse
												</>
											)}
										</Button>
									</div>
								</motion.div>
							)}
						</form>
					</Form>
				</CardContent>

				<CardFooter className='flex justify-between border-t pt-4'>
					<div className='flex items-center space-x-2'>
						<div
							className={`w-3 h-3 rounded-full ${
								step >= 1 ? 'bg-primary' : 'bg-gray-300'
							}`}
						></div>
						<div
							className={`w-3 h-3 rounded-full ${
								step >= 2 ? 'bg-primary' : 'bg-gray-300'
							}`}
						></div>
						<div
							className={`w-3 h-3 rounded-full ${
								step >= 3 ? 'bg-primary' : 'bg-gray-300'
							}`}
						></div>
					</div>
					<div className='text-sm text-muted-foreground'>
						Step {step} of 3:{' '}
						{step === 1
							? 'Basic Information'
							: step === 2
							? 'Health & Breeding'
							: 'Lineage & Details'}
					</div>
				</CardFooter>
			</Card>
		</motion.div>
	);
};

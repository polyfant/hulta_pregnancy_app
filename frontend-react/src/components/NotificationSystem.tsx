import { motion } from 'framer-motion';
import { BellRing, Mail, MessageSquare, Settings, Wifi } from 'lucide-react';
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
import { Switch } from './ui/switch';
import { Tabs, TabsContent, TabsList, TabsTrigger } from './ui/tabs';
import { useToast } from './ui/use-toast';

interface NotificationPreferences {
	email: {
		enabled: boolean;
		pregnancyMilestones: boolean;
		healthChecks: boolean;
		weatherAlerts: boolean;
	};
	webSocket: {
		enabled: boolean;
		pregnancyMilestones: boolean;
		healthChecks: boolean;
		weatherAlerts: boolean;
	};
	sms: {
		enabled: boolean;
		pregnancyMilestones: boolean;
		healthChecks: boolean;
		criticalAlertsOnly: boolean;
	};
}

interface Notification {
	id: string;
	type: 'email' | 'webSocket' | 'sms';
	title: string;
	message: string;
	timestamp: Date;
	read: boolean;
	category: 'pregnancyMilestone' | 'healthCheck' | 'weatherAlert';
	priority: 'low' | 'medium' | 'high' | 'critical';
}

export const NotificationSystem: React.FC = () => {
	const [preferences, setPreferences] = useState<NotificationPreferences>({
		email: {
			enabled: true,
			pregnancyMilestones: true,
			healthChecks: true,
			weatherAlerts: false,
		},
		webSocket: {
			enabled: true,
			pregnancyMilestones: true,
			healthChecks: true,
			weatherAlerts: true,
		},
		sms: {
			enabled: true,
			pregnancyMilestones: false,
			healthChecks: false,
			criticalAlertsOnly: true,
		},
	});

	const [notifications, setNotifications] = useState<Notification[]>([
		{
			id: '1',
			type: 'email',
			title: 'Pregnancy Milestone Alert',
			message:
				'Bella has reached the Mid-term stage of pregnancy. Consider scheduling a checkup.',
			timestamp: new Date(Date.now() - 24 * 60 * 60 * 1000),
			read: false,
			category: 'pregnancyMilestone',
			priority: 'medium',
		},
		{
			id: '2',
			type: 'webSocket',
			title: 'Health Check Reminder',
			message: 'Thunder is due for vaccination in 2 days.',
			timestamp: new Date(Date.now() - 2 * 60 * 60 * 1000),
			read: false,
			category: 'healthCheck',
			priority: 'high',
		},
		{
			id: '3',
			type: 'sms',
			title: 'Critical Weather Alert',
			message:
				'Extreme temperature warning: Take measures to protect pregnant mares.',
			timestamp: new Date(),
			read: false,
			category: 'weatherAlert',
			priority: 'critical',
		},
	]);

	const [activeTab, setActiveTab] = useState('all');
	const [connectionStatus, setConnectionStatus] = useState<
		'connected' | 'disconnected'
	>('connected');
	const { toast } = useToast();

	// Simulate WebSocket connection
	useEffect(() => {
		// Connection simulation
		const timer = setTimeout(() => {
			setConnectionStatus('connected');
			toast({
				title: 'Notification System Connected',
				description: 'You will receive real-time updates',
			});
		}, 1500);

		// Cleanup
		return () => clearTimeout(timer);
	}, [toast]);

	// Mock WebSocket message reception
	useEffect(() => {
		if (connectionStatus === 'connected' && preferences.webSocket.enabled) {
			const interval = setInterval(() => {
				const shouldSendNotification = Math.random() > 0.7;

				if (shouldSendNotification) {
					const newNotification: Notification = {
						id: Date.now().toString(),
						type: 'webSocket',
						title: 'Real-time Horse Update',
						message: `Luna's vital signs have been updated. All readings normal.`,
						timestamp: new Date(),
						read: false,
						category: 'healthCheck',
						priority: 'low',
					};

					setNotifications((prev) => [newNotification, ...prev]);

					toast({
						title: newNotification.title,
						description: newNotification.message,
					});
				}
			}, 45000); // Every 45 seconds for demo purposes

			return () => clearInterval(interval);
		}
	}, [connectionStatus, preferences.webSocket.enabled, toast]);

	const handleTogglePreference = (
		channel: keyof NotificationPreferences,
		setting: string,
		value: boolean
	) => {
		setPreferences((prev) => ({
			...prev,
			[channel]: {
				...prev[channel],
				[setting]: value,
			},
		}));

		toast({
			title: 'Preferences Updated',
			description: `${channel} ${setting} is now ${
				value ? 'enabled' : 'disabled'
			}`,
		});
	};

	const markAsRead = (id: string) => {
		setNotifications(
			notifications.map((notification) =>
				notification.id === id
					? { ...notification, read: true }
					: notification
			)
		);
	};

	const markAllAsRead = () => {
		setNotifications(
			notifications.map((notification) => ({
				...notification,
				read: true,
			}))
		);

		toast({
			title: 'All Notifications Marked as Read',
		});
	};

	const deleteNotification = (id: string) => {
		setNotifications(
			notifications.filter((notification) => notification.id !== id)
		);
	};

	const filteredNotifications = notifications.filter((notification) => {
		if (activeTab === 'all') return true;
		if (activeTab === 'unread') return !notification.read;
		return notification.type === activeTab;
	});

	const unreadCount = notifications.filter((n) => !n.read).length;

	return (
		<motion.div
			initial={{ opacity: 0, y: 20 }}
			animate={{ opacity: 1, y: 0 }}
			transition={{ duration: 0.5 }}
		>
			<Card className='w-full'>
				<CardHeader>
					<div className='flex items-center justify-between'>
						<div>
							<CardTitle className='text-xl font-bold'>
								Notification Center
							</CardTitle>
							<CardDescription>
								Manage your notifications and preferences
							</CardDescription>
						</div>
						<Badge
							variant={
								connectionStatus === 'connected'
									? 'success'
									: 'destructive'
							}
						>
							{connectionStatus === 'connected'
								? 'Connected'
								: 'Disconnected'}
						</Badge>
					</div>
				</CardHeader>

				<CardContent className='space-y-6'>
					<Tabs defaultValue='notifications' className='w-full'>
						<TabsList className='grid w-full grid-cols-2'>
							<TabsTrigger value='notifications'>
								<BellRing className='h-4 w-4 mr-2' />
								Notifications{' '}
								{unreadCount > 0 && `(${unreadCount})`}
							</TabsTrigger>
							<TabsTrigger value='preferences'>
								<Settings className='h-4 w-4 mr-2' />
								Preferences
							</TabsTrigger>
						</TabsList>

						<TabsContent
							value='notifications'
							className='space-y-4 mt-4'
						>
							<div className='flex items-center justify-between'>
								<TabsList className='inline-flex'>
									<TabsTrigger
										value='all'
										onClick={() => setActiveTab('all')}
										className={
											activeTab === 'all'
												? 'bg-primary text-primary-foreground'
												: ''
										}
									>
										All
									</TabsTrigger>
									<TabsTrigger
										value='unread'
										onClick={() => setActiveTab('unread')}
										className={
											activeTab === 'unread'
												? 'bg-primary text-primary-foreground'
												: ''
										}
									>
										Unread
									</TabsTrigger>
									<TabsTrigger
										value='email'
										onClick={() => setActiveTab('email')}
										className={
											activeTab === 'email'
												? 'bg-primary text-primary-foreground'
												: ''
										}
									>
										<Mail className='h-4 w-4 mr-1' />
										Email
									</TabsTrigger>
									<TabsTrigger
										value='webSocket'
										onClick={() =>
											setActiveTab('webSocket')
										}
										className={
											activeTab === 'webSocket'
												? 'bg-primary text-primary-foreground'
												: ''
										}
									>
										<Wifi className='h-4 w-4 mr-1' />
										Real-time
									</TabsTrigger>
									<TabsTrigger
										value='sms'
										onClick={() => setActiveTab('sms')}
										className={
											activeTab === 'sms'
												? 'bg-primary text-primary-foreground'
												: ''
										}
									>
										<MessageSquare className='h-4 w-4 mr-1' />
										SMS
									</TabsTrigger>
								</TabsList>

								<Button
									variant='outline'
									size='sm'
									onClick={markAllAsRead}
									disabled={
										!filteredNotifications.some(
											(n) => !n.read
										)
									}
								>
									Mark all as read
								</Button>
							</div>

							<div className='space-y-3 mt-4'>
								{filteredNotifications.length > 0 ? (
									filteredNotifications.map(
										(notification) => (
											<motion.div
												key={notification.id}
												initial={{ opacity: 0, x: -20 }}
												animate={{ opacity: 1, x: 0 }}
												className={`p-4 rounded-lg border ${
													!notification.read
														? 'bg-primary/5 border-primary/20'
														: ''
												}`}
											>
												<div className='flex items-start justify-between'>
													<div className='flex items-start gap-3'>
														{notification.type ===
															'email' && (
															<Mail
																className={`h-5 w-5 ${
																	!notification.read
																		? 'text-primary'
																		: 'text-muted-foreground'
																}`}
															/>
														)}
														{notification.type ===
															'webSocket' && (
															<Wifi
																className={`h-5 w-5 ${
																	!notification.read
																		? 'text-primary'
																		: 'text-muted-foreground'
																}`}
															/>
														)}
														{notification.type ===
															'sms' && (
															<MessageSquare
																className={`h-5 w-5 ${
																	!notification.read
																		? 'text-primary'
																		: 'text-muted-foreground'
																}`}
															/>
														)}

														<div className='flex-1'>
															<div className='flex items-center gap-2'>
																<p className='font-medium'>
																	{
																		notification.title
																	}
																</p>
																<Badge
																	variant={
																		notification.priority ===
																		'critical'
																			? 'destructive'
																			: notification.priority ===
																			  'high'
																			? 'default'
																			: 'outline'
																	}
																	className='ml-2'
																>
																	{
																		notification.priority
																	}
																</Badge>
															</div>
															<p className='text-sm text-muted-foreground mt-1'>
																{
																	notification.message
																}
															</p>
															<p className='text-xs text-muted-foreground mt-2'>
																{notification.timestamp.toLocaleString()}
															</p>
														</div>
													</div>

													<div className='flex gap-2'>
														{!notification.read && (
															<Button
																variant='outline'
																size='sm'
																onClick={() =>
																	markAsRead(
																		notification.id
																	)
																}
															>
																Mark as read
															</Button>
														)}
														<Button
															variant='ghost'
															size='sm'
															onClick={() =>
																deleteNotification(
																	notification.id
																)
															}
														>
															Delete
														</Button>
													</div>
												</div>
											</motion.div>
										)
									)
								) : (
									<div className='text-center py-12 text-muted-foreground'>
										No notifications to display
									</div>
								)}
							</div>
						</TabsContent>

						<TabsContent
							value='preferences'
							className='space-y-6 mt-4'
						>
							{/* Email Preferences */}
							<div className='rounded-lg border p-4'>
								<div className='flex items-center justify-between mb-4'>
									<div className='flex items-center gap-2'>
										<Mail className='h-5 w-5' />
										<h3 className='font-semibold'>
											Email Notifications
										</h3>
									</div>
									<Switch
										checked={preferences.email.enabled}
										onCheckedChange={(checked) =>
											handleTogglePreference(
												'email',
												'enabled',
												checked
											)
										}
									/>
								</div>

								{preferences.email.enabled && (
									<div className='space-y-3 pl-7'>
										<div className='flex items-center justify-between'>
											<label className='text-sm'>
												Pregnancy Milestones
											</label>
											<Switch
												checked={
													preferences.email
														.pregnancyMilestones
												}
												onCheckedChange={(checked) =>
													handleTogglePreference(
														'email',
														'pregnancyMilestones',
														checked
													)
												}
											/>
										</div>
										<div className='flex items-center justify-between'>
											<label className='text-sm'>
												Health Check Reminders
											</label>
											<Switch
												checked={
													preferences.email
														.healthChecks
												}
												onCheckedChange={(checked) =>
													handleTogglePreference(
														'email',
														'healthChecks',
														checked
													)
												}
											/>
										</div>
										<div className='flex items-center justify-between'>
											<label className='text-sm'>
												Weather Alerts
											</label>
											<Switch
												checked={
													preferences.email
														.weatherAlerts
												}
												onCheckedChange={(checked) =>
													handleTogglePreference(
														'email',
														'weatherAlerts',
														checked
													)
												}
											/>
										</div>
									</div>
								)}
							</div>

							{/* WebSocket Preferences */}
							<div className='rounded-lg border p-4'>
								<div className='flex items-center justify-between mb-4'>
									<div className='flex items-center gap-2'>
										<Wifi className='h-5 w-5' />
										<h3 className='font-semibold'>
											Real-time Notifications
										</h3>
									</div>
									<Switch
										checked={preferences.webSocket.enabled}
										onCheckedChange={(checked) =>
											handleTogglePreference(
												'webSocket',
												'enabled',
												checked
											)
										}
									/>
								</div>

								{preferences.webSocket.enabled && (
									<div className='space-y-3 pl-7'>
										<div className='flex items-center justify-between'>
											<label className='text-sm'>
												Pregnancy Milestones
											</label>
											<Switch
												checked={
													preferences.webSocket
														.pregnancyMilestones
												}
												onCheckedChange={(checked) =>
													handleTogglePreference(
														'webSocket',
														'pregnancyMilestones',
														checked
													)
												}
											/>
										</div>
										<div className='flex items-center justify-between'>
											<label className='text-sm'>
												Health Check Reminders
											</label>
											<Switch
												checked={
													preferences.webSocket
														.healthChecks
												}
												onCheckedChange={(checked) =>
													handleTogglePreference(
														'webSocket',
														'healthChecks',
														checked
													)
												}
											/>
										</div>
										<div className='flex items-center justify-between'>
											<label className='text-sm'>
												Weather Alerts
											</label>
											<Switch
												checked={
													preferences.webSocket
														.weatherAlerts
												}
												onCheckedChange={(checked) =>
													handleTogglePreference(
														'webSocket',
														'weatherAlerts',
														checked
													)
												}
											/>
										</div>
									</div>
								)}
							</div>

							{/* SMS Preferences */}
							<div className='rounded-lg border p-4'>
								<div className='flex items-center justify-between mb-4'>
									<div className='flex items-center gap-2'>
										<MessageSquare className='h-5 w-5' />
										<h3 className='font-semibold'>
											SMS Notifications
										</h3>
									</div>
									<Switch
										checked={preferences.sms.enabled}
										onCheckedChange={(checked) =>
											handleTogglePreference(
												'sms',
												'enabled',
												checked
											)
										}
									/>
								</div>

								{preferences.sms.enabled && (
									<div className='space-y-3 pl-7'>
										<div className='flex items-center justify-between'>
											<label className='text-sm'>
												Critical Alerts Only
											</label>
											<Switch
												checked={
													preferences.sms
														.criticalAlertsOnly
												}
												onCheckedChange={(checked) =>
													handleTogglePreference(
														'sms',
														'criticalAlertsOnly',
														checked
													)
												}
											/>
										</div>
										<div className='flex items-center justify-between'>
											<label className='text-sm'>
												Pregnancy Milestones
											</label>
											<Switch
												checked={
													preferences.sms
														.pregnancyMilestones
												}
												disabled={
													preferences.sms
														.criticalAlertsOnly
												}
												onCheckedChange={(checked) =>
													handleTogglePreference(
														'sms',
														'pregnancyMilestones',
														checked
													)
												}
											/>
										</div>
										<div className='flex items-center justify-between'>
											<label className='text-sm'>
												Health Check Reminders
											</label>
											<Switch
												checked={
													preferences.sms.healthChecks
												}
												disabled={
													preferences.sms
														.criticalAlertsOnly
												}
												onCheckedChange={(checked) =>
													handleTogglePreference(
														'sms',
														'healthChecks',
														checked
													)
												}
											/>
										</div>

										<div className='mt-3 p-3 rounded bg-muted/50 text-xs text-muted-foreground'>
											<p>
												SMS notifications may incur
												charges. With "Critical Alerts
												Only" enabled, you'll only
												receive SMS for urgent matters
												like extreme weather or health
												emergencies.
											</p>
										</div>
									</div>
								)}
							</div>
						</TabsContent>
					</Tabs>
				</CardContent>

				<CardFooter className='justify-between text-xs text-muted-foreground border-t pt-4'>
					<p>
						Notification delivery is subject to network conditions
					</p>
					<Button variant='link' size='sm' className='h-auto p-0'>
						Privacy Policy
					</Button>
				</CardFooter>
			</Card>
		</motion.div>
	);
};

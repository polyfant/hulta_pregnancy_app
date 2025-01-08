import { notifications } from '@mantine/notifications';
import { Baby, Syringe } from '@phosphor-icons/react';

interface NotificationEvent {
	id: string;
	type:
		| 'DUE_DATE'
		| 'VET_CHECK'
		| 'WEIGHT_CHECK'
		| 'STAGE_CHANGE'
		| 'VACCINATION';
	horseName: string;
	daysUntil?: number;
	date?: string;
	details?: string;
	priority: 'high' | 'medium' | 'low';
}

export const notificationService = {
	async checkAndNotify() {
		const response = await fetch('/api/notifications/pending');
		const events: NotificationEvent[] = await response.json();

		events.forEach((event) => {
			const notification = this.createNotification(event);
			notifications.show(notification);

			// If it's a high priority notification, also trigger system notification
			if (event.priority === 'high') {
				this.triggerSystemNotification(notification);
			}
		});
	},

	createNotification(event: NotificationEvent) {
		const baseNotification = {
			id: event.id,
			autoClose: event.priority === 'high' ? false : 5000,
			withCloseButton: true,
		};

		switch (event.type) {
			case 'DUE_DATE':
				return {
					...baseNotification,
					title: 'Due Date Alert',
					message: `${event.horseName} is due in ${event.daysUntil} days`,
					icon: <Baby />,
					color: 'blue',
					styles: { root: { borderLeft: '4px solid blue' } },
				};
			case 'VET_CHECK':
				return {
					...baseNotification,
					title: 'Veterinary Check Required',
					message: `Schedule check for ${event.horseName} - ${event.details}`,
					icon: <Syringe />,
					color: 'green',
				};
			// ... other cases
		}
	},

	async triggerSystemNotification({
		title,
		message,
	}: {
		title: string;
		message: string;
	}) {
		if ('Notification' in window && Notification.permission === 'granted') {
			new Notification(title, { body: message });
		}
	},

	async exportNotifications(horseId: number, format: 'PDF' | 'CSV' = 'PDF') {
		const response = await fetch(
			`/api/horses/${horseId}/notifications/export?format=${format}`
		);
		const blob = await response.blob();

		const url = window.URL.createObjectURL(blob);
		const a = document.createElement('a');
		a.href = url;
		a.download = `notifications_${horseId}_${
			new Date().toISOString().split('T')[0]
		}.${format.toLowerCase()}`;
		document.body.appendChild(a);
		a.click();
		document.body.removeChild(a);
		window.URL.revokeObjectURL(url);
	},
};

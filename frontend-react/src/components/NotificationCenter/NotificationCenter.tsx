import {
	ActionIcon,
	Badge,
	Card,
	Drawer,
	Group,
	Menu,
	SegmentedControl,
	Stack,
	Text,
} from '@mantine/core';
import { Bell, FileExport } from '@phosphor-icons/react';
import { useQuery } from '@tanstack/react-query';
import { useState } from 'react';
import { notificationService } from '../../services/notificationService';

export function NotificationCenter() {
	const [opened, setOpened] = useState(false);
	const [filter, setFilter] = useState('all');

	const { data: notifications } = useQuery({
		queryKey: ['notifications'],
		queryFn: async () => {
			const response = await fetch('/api/notifications');
			return response.json();
		},
		refetchInterval: 60000, // Refresh every minute
	});

	const filteredNotifications = notifications?.filter(
		(n) => filter === 'all' || n.priority === filter
	);

	return (
		<>
			<ActionIcon
				variant='light'
				onClick={() => setOpened(true)}
				pos='relative'
			>
				<Bell size={20} />
				{notifications?.some((n) => n.priority === 'high') && (
					<Badge
						size='xs'
						variant='filled'
						color='red'
						pos='absolute'
						top={-2}
						right={-2}
					/>
				)}
			</ActionIcon>

			<Drawer
				opened={opened}
				onClose={() => setOpened(false)}
				title='Notifications'
				position='right'
				size='md'
			>
				<Stack>
					<Group position='apart'>
						<SegmentedControl
							value={filter}
							onChange={setFilter}
							data={[
								{ label: 'All', value: 'all' },
								{ label: 'High', value: 'high' },
								{ label: 'Medium', value: 'medium' },
								{ label: 'Low', value: 'low' },
							]}
						/>
						<Menu>
							<Menu.Target>
								<ActionIcon>
									<FileExport size={20} />
								</ActionIcon>
							</Menu.Target>
							<Menu.Dropdown>
								<Menu.Item
									onClick={() =>
										notificationService.exportNotifications(
											0,
											'PDF'
										)
									}
								>
									Export as PDF
								</Menu.Item>
								<Menu.Item
									onClick={() =>
										notificationService.exportNotifications(
											0,
											'CSV'
										)
									}
								>
									Export as CSV
								</Menu.Item>
							</Menu.Dropdown>
						</Menu>
					</Group>

					{filteredNotifications?.map((notification) => (
						<Card key={notification.id} withBorder>
							<Group position='apart'>
								<div>
									<Text fw={500}>{notification.title}</Text>
									<Text size='sm' c='dimmed'>
										{notification.message}
									</Text>
								</div>
								<Badge
									color={
										notification.priority === 'high'
											? 'red'
											: notification.priority === 'medium'
											? 'yellow'
											: 'blue'
									}
								>
									{notification.priority}
								</Badge>
							</Group>
						</Card>
					))}
				</Stack>
			</Drawer>
		</>
	);
}

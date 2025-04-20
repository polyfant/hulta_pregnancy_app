import { LogIn, LogOut, Settings, Shield, UserCircle } from 'lucide-react';
import React from 'react';
import { useAuth } from '../../auth/useAuth';
import { Avatar, AvatarFallback, AvatarImage } from '../ui/avatar';
import { Badge } from '../ui/badge';
import { Button } from '../ui/button';
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuGroup,
	DropdownMenuItem,
	DropdownMenuLabel,
	DropdownMenuSeparator,
	DropdownMenuTrigger,
} from '../ui/dropdown-menu';

export const LoginButton: React.FC = () => {
	const { isAuthenticated, isLoading, user, login, logout } = useAuth();

	// Helper function to get initials from name
	const getInitials = (name: string) => {
		return name
			.split(' ')
			.map((n) => n[0])
			.join('')
			.toUpperCase();
	};

	// Get role badge variant
	const getRoleBadgeVariant = (role: string) => {
		switch (role) {
			case 'admin':
				return 'destructive';
			case 'owner':
				return 'default';
			case 'farm_manager':
				return 'secondary';
			default:
				return 'outline';
		}
	};

	// If loading, show loading state
	if (isLoading) {
		return (
			<Button variant='outline' disabled>
				<span className='h-4 w-4 mr-2 rounded-full border-2 border-current border-t-transparent animate-spin' />
				Loading...
			</Button>
		);
	}

	// If not authenticated, show login button
	if (!isAuthenticated) {
		return (
			<Button onClick={() => login()}>
				<LogIn className='h-4 w-4 mr-2' />
				Login
			</Button>
		);
	}

	// If authenticated, show user profile dropdown
	return (
		<DropdownMenu>
			<DropdownMenuTrigger asChild>
				<Button variant='outline' className='p-1 px-2 h-10'>
					<Avatar className='h-7 w-7 mr-2'>
						<AvatarImage
							src={user?.picture}
							alt={user?.name || 'User'}
						/>
						<AvatarFallback>
							{user?.name ? getInitials(user.name) : 'U'}
						</AvatarFallback>
					</Avatar>
					<span className='text-sm hidden md:inline-block'>
						{user?.name}
					</span>
				</Button>
			</DropdownMenuTrigger>
			<DropdownMenuContent align='end' className='w-56'>
				<DropdownMenuLabel className='font-normal'>
					<div className='flex flex-col space-y-1'>
						<p className='text-sm font-medium leading-none'>
							{user?.name}
						</p>
						<p className='text-xs leading-none text-muted-foreground'>
							{user?.email}
						</p>
						<div className='flex items-center mt-1.5'>
							<Shield className='h-3 w-3 mr-1 text-muted-foreground' />
							<Badge
								variant={getRoleBadgeVariant(
									user?.role || 'user'
								)}
								className='text-[10px] h-4'
							>
								{user?.role}
							</Badge>
						</div>
					</div>
				</DropdownMenuLabel>
				<DropdownMenuSeparator />
				<DropdownMenuGroup>
					<DropdownMenuItem>
						<UserCircle className='h-4 w-4 mr-2' />
						<span>My Profile</span>
					</DropdownMenuItem>
					<DropdownMenuItem>
						<Settings className='h-4 w-4 mr-2' />
						<span>Settings</span>
					</DropdownMenuItem>
				</DropdownMenuGroup>
				<DropdownMenuSeparator />
				<DropdownMenuItem onClick={() => logout()}>
					<LogOut className='h-4 w-4 mr-2' />
					<span>Log out</span>
				</DropdownMenuItem>
			</DropdownMenuContent>
		</DropdownMenu>
	);
};

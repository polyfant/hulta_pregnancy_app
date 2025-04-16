import { Loader2 } from 'lucide-react';
import React, { useEffect } from 'react';
import { useLocation } from 'react-router-dom';
import { useAuth } from './useAuth';

interface ProtectedRouteProps {
	children: React.ReactNode;
	requiredRoles?: Array<'user' | 'admin' | 'owner' | 'farm_manager'>;
}

export const ProtectedRoute: React.FC<ProtectedRouteProps> = ({
	children,
	requiredRoles,
}) => {
	const { isAuthenticated, isLoading, user, login } = useAuth();
	const location = useLocation();

	useEffect(() => {
		// If not loading and not authenticated, redirect to login
		if (!isLoading && !isAuthenticated) {
			login(location.pathname);
		}
	}, [isLoading, isAuthenticated, login, location.pathname]);

	// Show loading state
	if (isLoading) {
		return (
			<div className='flex items-center justify-center min-h-screen'>
				<div className='flex flex-col items-center'>
					<Loader2 className='h-12 w-12 animate-spin text-primary' />
					<p className='mt-4 text-lg font-medium text-muted-foreground'>
						Loading...
					</p>
				</div>
			</div>
		);
	}

	// Check if authenticated
	if (!isAuthenticated) {
		return null; // Don't render anything while redirecting to login
	}

	// Check for required roles if specified
	if (requiredRoles && user && !requiredRoles.includes(user.role)) {
		return (
			<div className='flex items-center justify-center min-h-screen'>
				<div className='p-6 bg-destructive/10 border border-destructive/30 rounded-lg max-w-lg text-center'>
					<h2 className='text-2xl font-semibold text-destructive mb-2'>
						Access Denied
					</h2>
					<p className='text-muted-foreground'>
						You don't have permission to access this page. This area
						requires
						{requiredRoles.length > 1
							? ` one of these roles: ${requiredRoles.join(', ')}`
							: ` the ${requiredRoles[0]} role`}
						.
					</p>
					<p className='mt-4 text-sm'>
						Your current role:{' '}
						<span className='font-medium'>{user.role}</span>
					</p>
				</div>
			</div>
		);
	}

	// If all checks pass, render the children
	return <>{children}</>;
};

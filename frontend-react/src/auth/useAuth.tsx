import { useAuth0 } from '@auth0/auth0-react';
import { useEffect, useState } from 'react';

interface User {
	id: string;
	email: string;
	name: string;
	role: 'user' | 'admin' | 'owner' | 'farm_manager';
	picture?: string;
}

export const useAuth = () => {
	const {
		isAuthenticated,
		isLoading,
		user,
		loginWithRedirect,
		logout,
		getAccessTokenSilently,
	} = useAuth0();

	const [token, setToken] = useState<string | null>(null);
	const [userProfile, setUserProfile] = useState<User | null>(null);
	const [tokenLoading, setTokenLoading] = useState(false);

	// Get and store the access token whenever auth state changes
	useEffect(() => {
		const getToken = async () => {
			if (isAuthenticated && user) {
				setTokenLoading(true);
				try {
					const accessToken = await getAccessTokenSilently();
					setToken(accessToken);

					// Extract user role from Auth0 user metadata or claims
					// This depends on how your Auth0 is configured
					const role =
						user['https://hulta-app.com/roles']?.[0] || 'user';

					setUserProfile({
						id: user.sub!,
						email: user.email!,
						name: user.name || 'User',
						role: role as
							| 'user'
							| 'admin'
							| 'owner'
							| 'farm_manager',
						picture: user.picture,
					});
				} catch (error) {
					console.error('Error getting access token', error);
				} finally {
					setTokenLoading(false);
				}
			} else {
				setToken(null);
				setUserProfile(null);
			}
		};

		getToken();
	}, [isAuthenticated, user, getAccessTokenSilently]);

	// Function to get a fresh token for API calls
	const getToken = async (): Promise<string | null> => {
		if (!isAuthenticated) return null;

		try {
			const accessToken = await getAccessTokenSilently();
			setToken(accessToken);
			return accessToken;
		} catch (error) {
			console.error('Error getting token', error);
			return null;
		}
	};

	const loginWithRedirectAndParams = (returnTo?: string) => {
		const options: any = {};
		if (returnTo) {
			options.appState = { returnTo };
		}
		return loginWithRedirect(options);
	};

	return {
		isAuthenticated,
		isLoading: isLoading || tokenLoading,
		user: userProfile,
		token,
		getToken,
		login: loginWithRedirectAndParams,
		logout: () =>
			logout({ logoutParams: { returnTo: window.location.origin } }),
		hasRole: (role: string | string[]) => {
			if (!userProfile) return false;
			if (Array.isArray(role)) {
				return role.includes(userProfile.role);
			}
			return userProfile.role === role;
		},
	};
};

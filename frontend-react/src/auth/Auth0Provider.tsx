import { Auth0Provider as Auth0ProviderBase } from '@auth0/auth0-react';
import React from 'react';
import { useNavigate } from 'react-router-dom';

export const Auth0Provider: React.FC<React.PropsWithChildren> = ({
	children,
}) => {
	const navigate = useNavigate();
	const domain = import.meta.env.VITE_AUTH0_DOMAIN;
	const clientId = import.meta.env.VITE_AUTH0_CLIENT_ID;
	const audience = import.meta.env.VITE_AUTH0_AUDIENCE;
	const redirectUri = window.location.origin;

	const onRedirectCallback = (appState: any) => {
		navigate(appState?.returnTo || window.location.pathname);
	};

	if (!(domain && clientId && audience)) {
		return (
			<div className='flex items-center justify-center min-h-screen'>
				<div className='p-4 rounded-md bg-red-50 text-red-700 border border-red-200 max-w-lg'>
					<h3 className='text-lg font-semibold'>
						Authentication Configuration Error
					</h3>
					<p className='mt-2'>
						Missing Auth0 configuration. Please check your
						environment variables:
					</p>
					<ul className='list-disc list-inside mt-2'>
						{!domain && <li>VITE_AUTH0_DOMAIN</li>}
						{!clientId && <li>VITE_AUTH0_CLIENT_ID</li>}
						{!audience && <li>VITE_AUTH0_AUDIENCE</li>}
					</ul>
				</div>
			</div>
		);
	}

	return (
		<Auth0ProviderBase
			domain={domain}
			clientId={clientId}
			authorizationParams={{
				redirect_uri: redirectUri,
				audience: audience,
			}}
			onRedirectCallback={onRedirectCallback}
			cacheLocation='localstorage'
		>
			{children}
		</Auth0ProviderBase>
	);
};

import { Auth0Provider } from '@auth0/auth0-react';
import { createRoot } from 'react-dom/client';
import { BrowserRouter } from 'react-router-dom';
import App from './App';

import '@mantine/core/styles.css';
import '@mantine/notifications/styles.css';

const root = createRoot(document.getElementById('root')!);

root.render(
	<Auth0Provider
		domain='dev-r083cwkcv0pgz20x.eu.auth0.com'
		clientId='OBmEol1z4U49r3YI3priDdGbvF5i4O7d'
		authorizationParams={{
			redirect_uri: window.location.origin,
			scope: 'openid profile email',
			response_type: 'code',
			audience: 'https://api.hulta-pregnancy.app'
		}}
		cacheLocation='localstorage'
		useRefreshTokens={true}
		onRedirectCallback={(appState) => {
			window.history.replaceState(
				{},
				document.title,
				appState?.returnTo || window.location.origin
			);
		}}
	>
		<BrowserRouter>
			<App />
		</BrowserRouter>
	</Auth0Provider>
);

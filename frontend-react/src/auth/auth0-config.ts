import { Auth0Config } from '../types/auth0-config';

export const auth0Config: Auth0Config = {
  domain: 'dev-r083cwkcv0pgz20x.eu.auth0.com',
  clientId: 'OBmEol1z4U49r3YI3priDdGbvF5i4O7d',
  authorizationParams: {
    redirect_uri: `${window.location.origin}/callback`, // Dynamic redirect URI
    scope: 'openid profile email',
    audience: 'https://dev-r083cwkcv0pgz20x.eu.auth0.com/api/v2/'
  },
  cacheLocation: 'localstorage',
  useRefreshTokens: true,
  onRedirectCallback: (appState) => {
    console.log('Redirect callback executed');
    window.history.replaceState(
      {},
      document.title,
      appState?.returnTo || window.location.pathname
    );
  },
};

export interface Auth0Config {
  domain: string;
  clientId: string;
  authorizationParams: {
    redirect_uri: string;
    scope?: string;
    audience?: string;
  };
  cacheLocation: 'localstorage';
  useRefreshTokens?: boolean;
  onRedirectCallback: (appState: { returnTo?: string }) => void;
  returnTo?: string;
}

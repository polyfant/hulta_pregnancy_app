import React from 'react';
import { BrowserRouter as Router } from 'react-router-dom';
import { MantineProvider } from '@mantine/core';
import { Notifications } from '@mantine/notifications';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { AppShell, LoadingOverlay, Center, Text, Button } from '@mantine/core';
import { theme } from './theme';
import { Auth0Provider, useAuth0 } from '@auth0/auth0-react';
import { auth0Config } from './auth/auth0-config';
import Navbar from './components/Navbar/Navbar';
import AppRoutes from './routes';

// Import Mantine styles
import '@mantine/core/styles.css';
import '@mantine/notifications/styles.css';
import '@mantine/dates/styles.css';

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 1000 * 60 * 5, // 5 minutes
      retry: 1,
      refetchOnWindowFocus: false,
    },
  },
});

function AppContent() {
  const { 
    isLoading, 
    isAuthenticated, 
    loginWithRedirect, 
    error,
    logout 
  } = useAuth0();

  console.log('AppContent rendered');
  console.log('isLoading:', isLoading);
  console.log('isAuthenticated:', isAuthenticated);
  
  // Log any Auth0 errors
  if (error) {
    console.error('Auth0 Error:', error);
    return (
      <Center h="100vh" p="md">
        <div style={{ textAlign: 'center' }}>
          <Text c="red" size="xl" mb="md">Authentication Error</Text>
          <Text mb="md">{error.message}</Text>
          <Button 
            color="red" 
            onClick={() => {
              logout({ logoutParams: { returnTo: window.location.origin } });
              loginWithRedirect();
            }}
          >
            Retry Login
          </Button>
        </div>
      </Center>
    );
  }

  if (isLoading) {
    console.log('Loading authentication state...');
    return (
      <Center h="100vh">
        <LoadingOverlay visible={true} zIndex={1000} />
      </Center>
    );
  }

  if (!isAuthenticated) {
    console.log('User not authenticated, redirecting to login...');
    try {
      loginWithRedirect({
        appState: { returnTo: window.location.pathname }
      });
    } catch (redirectError) {
      console.error('Login redirect error:', redirectError);
      return (
        <Center h="100vh" p="md">
          <div style={{ textAlign: 'center' }}>
            <Text c="red" size="xl" mb="md">Login Redirect Failed</Text>
            <Text mb="md">{redirectError instanceof Error ? redirectError.message : 'Unknown error'}</Text>
            <Button 
              color="red" 
              onClick={() => loginWithRedirect()}
            >
              Retry Login
            </Button>
          </div>
        </Center>
      );
    }
    return (
      <Center h="100vh">
        <LoadingOverlay visible={true} zIndex={1000} />
      </Center>
    );
  }

  console.log('User authenticated, rendering app content');
  return (
    <AppShell
      header={{ height: 60 }}
      padding="md"
    >
      <AppShell.Header>
        <Navbar />
      </AppShell.Header>

      <AppShell.Main>
        <AppRoutes />
      </AppShell.Main>
    </AppShell>
  );
}

function App() {
  return (
    <Auth0Provider
      domain={auth0Config.domain}
      clientId={auth0Config.clientId}
      authorizationParams={auth0Config.authorizationParams}
      cacheLocation={auth0Config.cacheLocation}
      useRefreshTokens={auth0Config.useRefreshTokens}
      onRedirectCallback={(appState, user) => {
        if (appState && appState.returnTo) {
          window.history.replaceState(
            {},
            document.title,
            appState.returnTo
          );
        }
      }}
    >
      <QueryClientProvider client={queryClient}>
        <MantineProvider theme={theme}>
          <Notifications />
          <Router>
            <AppContent />
          </Router>
        </MantineProvider>
      </QueryClientProvider>
    </Auth0Provider>
  );
}

export default App;

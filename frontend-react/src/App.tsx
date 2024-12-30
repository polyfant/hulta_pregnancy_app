import { BrowserRouter as Router } from 'react-router-dom';
import { MantineProvider } from '@mantine/core';
import { Notifications } from '@mantine/notifications';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { AppShell } from '@mantine/core';
import { theme } from './theme';
import { Suspense, lazy } from 'react';

// Lazy load components
const Navbar = lazy(() => import('./components/Navbar/Navbar').then(module => ({ default: module.default })));
const AppRoutes = lazy(() => import('./routes').then(module => ({ default: module.default })));

// Import Mantine styles
import '@mantine/core/styles.css';
import '@mantine/notifications/styles.css';
import '@mantine/dates/styles.css';

// Loading component
const Loading = () => (
  <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
    Loading...
  </div>
);

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 1000 * 60 * 5, // 5 minutes
      retry: 1,
      refetchOnWindowFocus: false,
      refetchOnMount: false,
      refetchOnReconnect: false
    },
    mutations: {
      retry: 1
    }
  }
});

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <MantineProvider theme={theme}>
        <Notifications limit={5} />
        <Router>
          <AppShell
            header={{ height: 60 }}
            padding="md"
          >
            <AppShell.Header style={{ backgroundColor: 'var(--mantine-color-dark-7)', borderBottom: '1px solid var(--mantine-color-dark-4)' }}>
              <Suspense fallback={<Loading />}>
                <Navbar />
              </Suspense>
            </AppShell.Header>

            <AppShell.Main>
              <Suspense fallback={<Loading />}>
                <AppRoutes />
              </Suspense>
            </AppShell.Main>
          </AppShell>
        </Router>
      </MantineProvider>
    </QueryClientProvider>
  );
}

export default App;

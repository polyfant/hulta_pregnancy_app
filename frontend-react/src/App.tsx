import { BrowserRouter as Router } from 'react-router-dom';
import { MantineProvider } from '@mantine/core';
import { Notifications } from '@mantine/notifications';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { AppShell } from '@mantine/core';

import { Navbar } from './components/Navbar/Navbar';
import { AppRoutes } from './routes';

// Import Mantine styles
import '@mantine/core/styles.css';
import '@mantine/notifications/styles.css';
import '@mantine/dates/styles.css';

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 1000 * 60 * 5, // 5 minutes
      retry: 1
    }
  }
});

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <MantineProvider>
        <Notifications />
        <Router>
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
        </Router>
      </MantineProvider>
    </QueryClientProvider>
  );
}

export default App;

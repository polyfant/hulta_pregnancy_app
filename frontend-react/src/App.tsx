import { MantineProvider } from '@mantine/core';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { AddHorseForm } from './components/AddHorseForm';
import { HorseList } from './components/HorseList';
import '@mantine/core/styles.css';
import '@mantine/dates/styles.css';

const queryClient = new QueryClient();

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <MantineProvider>
        <div style={{ padding: '2rem' }}>
          <h1>Horse Management</h1>
          <div style={{ marginBottom: '2rem' }}>
            <h2>Add New Horse</h2>
            <AddHorseForm />
          </div>
          <div>
            <h2>Horses</h2>
            <HorseList />
          </div>
        </div>
      </MantineProvider>
    </QueryClientProvider>
  );
}

export default App;

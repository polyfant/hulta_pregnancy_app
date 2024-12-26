import { MantineProvider } from '@mantine/core';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { AddHorseForm } from './components/AddHorseForm';
import { HorseList } from './components/HorseList';
import PregnancyTracking from './components/PregnancyTracking/PregnancyTracking';
import '@mantine/core/styles.css';
import '@mantine/dates/styles.css';

const queryClient = new QueryClient();

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <MantineProvider>
        <Router>
          <div style={{ padding: '2rem' }}>
            <h1>Horse Management</h1>
            <Routes>
              <Route
                path="/"
                element={
                  <>
                    <div style={{ marginBottom: '2rem' }}>
                      <h2>Add New Horse</h2>
                      <AddHorseForm />
                    </div>
                    <div>
                      <h2>Horses</h2>
                      <HorseList />
                    </div>
                  </>
                }
              />
              <Route
                path="/horses/:id/pregnancy"
                element={<PregnancyTracking />}
              />
            </Routes>
          </div>
        </Router>
      </MantineProvider>
    </QueryClientProvider>
  );
}

export default App;

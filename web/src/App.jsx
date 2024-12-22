import { MantineProvider } from '@mantine/core';
import { DatabaseProvider } from './contexts/DatabaseContext';
import { Layout } from './components/Layout';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { Notifications } from '@mantine/notifications';

function App() {
  return (
    <MantineProvider withGlobalStyles withNormalizeCSS>
      <Notifications />
      <DatabaseProvider>
        <Router>
          <Layout>
            <Routes>
              <Route path="/" element={<div>Horses List (Coming Soon)</div>} />
              <Route path="/health" element={<div>Health Records (Coming Soon)</div>} />
              <Route path="/pregnancy" element={<div>Pregnancy Events (Coming Soon)</div>} />
            </Routes>
          </Layout>
        </Router>
      </DatabaseProvider>
    </MantineProvider>
  );
}

export default App;

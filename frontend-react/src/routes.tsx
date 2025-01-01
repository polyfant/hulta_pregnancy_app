import { Routes, Route, Navigate } from 'react-router-dom';
import { useAuth0 } from '@auth0/auth0-react';
import { Center, LoadingOverlay } from '@mantine/core';

import HorseList from './components/HorseList';
import AddHorse from './components/AddHorse';
import HorseDetails from './components/HorseDetails';
import EditHorse from './components/EditHorse';
import PregnancyTracking from './components/PregnancyTracking/PregnancyTracking';
import PregnancyGuidelines from './components/PregnancyTracking/PregnancyGuidelines';
import PreFoalingSigns from './components/PregnancyTracking/PreFoalingSigns';

function Callback() {
  const { isLoading } = useAuth0();

  if (isLoading) {
    return (
      <Center h="100vh">
        <LoadingOverlay visible={true} zIndex={1000} />
      </Center>
    );
  }

  return <Navigate to="/" replace />;
}

export default function AppRoutes() {
  return (
    <Routes>
      <Route path="/callback" element={<Callback />} />
      <Route path="/" element={<HorseList />} />
      <Route path="/add-horse" element={<AddHorse />} />
      <Route path="/horses/:id" element={<HorseDetails />} />
      <Route path="/horses/:id/edit" element={<EditHorse />} />
      <Route path="/horses/:id/pregnancy" element={<PregnancyTracking />} />
      <Route path="/horses/:id/pregnancy/guidelines" element={<PregnancyGuidelines />} />
      <Route path="/horses/:id/pregnancy/pre-foaling" element={<PreFoalingSigns horseId={":id"} />} />
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  );
}

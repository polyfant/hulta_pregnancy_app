import { Routes, Route, useParams } from 'react-router-dom';
import { lazy, Suspense } from 'react';
import { LoadingOverlay } from '@mantine/core';

// Lazy load components
const HorseList = lazy(() => import('./components/HorseList').then(module => ({ default: module.default })));
const AddHorse = lazy(() => import('./components/AddHorse').then(module => ({ default: module.default })));
const EditHorse = lazy(() => import('./components/EditHorse').then(module => ({ default: module.default })));
const HorseDetails = lazy(() => import('./components/HorseDetails').then(module => ({ default: module.default })));
const PregnancyTracking = lazy(() => import('./components/PregnancyTracking/PregnancyTracking').then(module => ({ default: module.default })));
const PregnancyGuidelines = lazy(() => import('./components/PregnancyTracking/PregnancyGuidelines').then(module => ({ default: module.default })));
const PreFoalingSigns = lazy(() => import('./components/PregnancyTracking/PreFoalingSigns').then(module => ({ default: module.default })));

// Loading component
const Loading = () => (
  <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
    Loading...
  </div>
);

// Wrapper component to handle URL parameters
const PreFoalingSignsWrapper = () => {
  const { id } = useParams();
  return <PreFoalingSigns horseId={id || ''} />;
};

const AppRoutes = () => {
  return (
    <Suspense fallback={<LoadingOverlay visible />}>
      <Routes>
        {/* Main routes */}
        <Route path="/" element={<HorseList />} />
        <Route path="/add-horse" element={<AddHorse />} />
        
        {/* Horse details and management */}
        <Route path="/horses/:id" element={<HorseDetails />} />
        <Route path="/horses/:id/edit" element={<EditHorse />} />
        
        {/* Pregnancy tracking routes */}
        <Route path="/horses/:id/pregnancy" element={<PregnancyTracking />} />
        <Route path="/horses/:id/pregnancy/guidelines" element={<PregnancyGuidelines />} />
        <Route path="/horses/:id/pregnancy/pre-foaling" element={<PreFoalingSignsWrapper />} />
      </Routes>
    </Suspense>
  );
};

export default AppRoutes;

import { Routes, Route } from 'react-router-dom';
import { HorseList } from './components/HorseList';
import { AddHorse } from './components/AddHorse';
import { HorseDetails } from './components/HorseDetails';
import { PregnancyTracking } from './components/PregnancyTracking/PregnancyTracking';
import { PregnancyGuidelines } from './components/PregnancyTracking/PregnancyGuidelines';
import { PreFoalingSigns } from './components/PregnancyTracking/PreFoalingSigns';

export function AppRoutes() {
  return (
    <Routes>
      {/* Main routes */}
      <Route path="/" element={<HorseList />} />
      <Route path="/add-horse" element={<AddHorse />} />
      
      {/* Horse details and management */}
      <Route path="/horses/:id" element={<HorseDetails />} />
      
      {/* Pregnancy tracking routes */}
      <Route path="/horses/:id/pregnancy" element={<PregnancyTracking />} />
      <Route path="/horses/:id/pregnancy/guidelines" element={<PregnancyGuidelines />} />
      <Route path="/horses/:id/pregnancy/pre-foaling" element={<PreFoalingSigns />} />
    </Routes>
  );
}

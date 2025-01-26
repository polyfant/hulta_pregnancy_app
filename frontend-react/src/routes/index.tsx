import { LoadingOverlay } from '@mantine/core';
import { lazy, Suspense } from 'react';
import { Route, Routes } from 'react-router-dom';
import { GrowthChartsTesting } from '../components/PregnancyTracking/GrowthChartsTesting';
import { PregnancyOverview } from '../components/PregnancyTracking/PregnancyOverview';

// Lazy load components
const HorseList = lazy(() => import('../components/HorseList'));
const HorseDetails = lazy(() => import('../components/HorseDetails'));
const AddHorse = lazy(() => import('../components/AddHorse'));
const EditHorse = lazy(() => import('../components/EditHorse'));
const PregnancyTracking = lazy(
	() => import('../components/PregnancyTracking/PregnancyTracking')
);

// Callback component
const Callback = () => (
	<div style={{ height: '100vh', position: 'relative' }}>
		<LoadingOverlay visible={true} />
	</div>
);

const AppRoutes = () => {
	return (
		<Suspense fallback={<LoadingOverlay visible />}>
			<Routes>
				<Route path='/' element={<HorseList />} />
				<Route path='/horses/:id' element={<HorseDetails />} />
				<Route path='/add-horse' element={<AddHorse />} />
				<Route path='/horses/:id/edit' element={<EditHorse />} />
				<Route
					path='/horses/:id/pregnancy'
					element={<PregnancyTracking />}
				/>
				<Route path='/callback' element={<Callback />} />
				<Route path='/pregnancies' element={<PregnancyOverview />} />
				{process.env.NODE_ENV === 'development' && (
					<Route
						path='/testing/growth-charts'
						element={<GrowthChartsTesting />}
					/>
				)}
			</Routes>
		</Suspense>
	);
};

export default AppRoutes;

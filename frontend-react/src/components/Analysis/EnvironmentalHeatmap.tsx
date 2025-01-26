import { Card, Text } from '@mantine/core';
import { HeatMap } from '@nivo/heatmap';

interface HeatmapData {
	hour: number;
	temperature: number;
	humidity: number;
	risk: number;
}

export function EnvironmentalHeatmap({ data }: { data: HeatmapData[] }) {
	const heatmapData = Array.from({ length: 7 }, (_, day) => ({
		id: `Day ${day + 1}`,
		data: Array.from({ length: 24 }, (_, hour) => ({
			x: `${hour}:00`,
			y: data[day * 24 + hour]?.risk || 0,
		})),
	}));

	return (
		<Card withBorder>
			<Text size='lg' fw={500} mb='md'>
				Environmental Risk Heatmap
			</Text>
			<div style={{ height: 300 }}>
				<HeatMap
					data={heatmapData}
					margin={{ top: 20, right: 50, bottom: 60, left: 60 }}
					axisTop={null}
					axisRight={null}
					axisBottom={{
						tickSize: 5,
						tickRotation: -45,
						legend: 'Hour',
						legendPosition: 'middle',
						legendOffset: 40,
					}}
					axisLeft={{
						tickSize: 5,
						legend: 'Day',
						legendPosition: 'middle',
						legendOffset: -40,
					}}
					colors={{
						type: 'sequential',
						scheme: 'YlOrRd',
					}}
					emptyColor='#ffffff'
					legends={[
						{
							anchor: 'bottom',
							translateX: 0,
							translateY: 30,
							length: 400,
							thickness: 8,
							title: 'Risk Level',
						},
					]}
				/>
			</div>
		</Card>
	);
}

import { animated, useSprings } from '@react-spring/web';
import { useDrag } from '@use-gesture/react';
import { useRef, useState } from 'react';

export function SwipeableChart({ charts }) {
	const [currentIndex, setCurrentIndex] = useState(0);
	const dragRef = useRef(null);

	const [springs, api] = useSprings(charts.length, (i) => ({
		x: i * window.innerWidth,
		scale: 1,
		display: 'block',
	}));

	const bind = useDrag(
		({ active, movement: [mx], direction: [xDir], cancel }) => {
			if (active && Math.abs(mx) > window.innerWidth / 2) {
				const newIndex = clamp(
					currentIndex + (xDir > 0 ? -1 : 1),
					0,
					charts.length - 1
				);
				setCurrentIndex(newIndex);
				cancel();
			}

			api.start((i) => {
				const x =
					(i - currentIndex) * window.innerWidth + (active ? mx : 0);
				const scale = active
					? 1 - Math.abs(mx) / window.innerWidth / 2
					: 1;
				return { x, scale, display: 'block' };
			});
		}
	);

	return (
		<div ref={dragRef} style={{ height: '100%', touchAction: 'none' }}>
			{springs.map(({ x, display, scale }, i) => (
				<animated.div
					{...bind()}
					key={i}
					style={{
						position: 'absolute',
						width: '100%',
						height: '100%',
						x,
						display,
						scale,
					}}
				>
					{charts[i]}
				</animated.div>
			))}
		</div>
	);
}

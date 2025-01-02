import {
    ActionIcon,
    Card,
    Group,
    LoadingOverlay,
    Paper,
    Text
} from '@mantine/core';
import {
    CaretDown,
    CaretRight,
    GenderFemale,
    GenderMale,
    Horse
} from '@phosphor-icons/react';
import { useState } from 'react';

import { useQuery } from '@tanstack/react-query';

interface Horse {
	id: number;
	name: string;
	breed?: string;
	gender: 'MARE' | 'STALLION' | 'GELDING';
	dateOfBirth: string;
	weight?: number;
	motherId?: number;
	fatherId?: number;
	externalMother?: string;
	externalFather?: string;
}

interface TreeNodeProps {
	horse: Horse;
	level: number;
	maxLevel?: number;
}

const TreeNode = ({ horse, level, maxLevel = 3 }: TreeNodeProps) => {
	const [expanded, setExpanded] = useState(level < 2);
	const hasParents =
		horse.fatherId ||
		horse.motherId ||
		horse.externalFather ||
		horse.externalMother;
	const canExpand = level < maxLevel && hasParents;

	const { data: parents, isLoading } = useQuery({
		queryKey: ['horseParents', horse.id],
		queryFn: async () => {
			if (!hasParents) return null;
			const response = await fetch(`/api/horses/${horse.id}/family`);
			if (!response.ok) throw new Error('Failed to fetch family tree');
			return response.json();
		},
		enabled: expanded && hasParents ? true : false,
	});

	const handleToggle = () => {
		if (canExpand) {
			setExpanded(!expanded);
		}
	};

	return (
		<div style={{ marginLeft: level > 0 ? 40 : 0 }}>
			<Card withBorder shadow='sm' mt='xs' bg="dark.7">
				<Group justify='space-between'>
					<Group>
						{canExpand && (
							<ActionIcon
								onClick={handleToggle}
								variant='subtle'
								loading={isLoading}
								c="white"
							>
								{expanded ? (
									<CaretDown size='1rem' />
								) : (
									<CaretRight size='1rem' />
								)}
							</ActionIcon>
						)}
						<Group>
							{horse.gender === 'STALLION' || horse.gender === 'GELDING' ? (
								<GenderMale size='1rem' color="var(--mantine-color-blue-6)" />
							) : (
								<GenderFemale size='1rem' color="var(--mantine-color-pink-6)" />
							)}
							<Text c="white" fw={500}>{horse.name}</Text>
						</Group>
					</Group>
				</Group>

				{expanded && hasParents && (
					<div>
						{parents?.mother && (
							<TreeNode
								horse={parents.mother}
								level={level + 1}
								maxLevel={maxLevel}
							/>
						)}
						{parents?.father && (
							<TreeNode
								horse={parents.father}
								level={level + 1}
								maxLevel={maxLevel}
							/>
						)}
					</div>
				)}
			</Card>
		</div>
	);
};

interface FamilyTreeProps {
	horseId: number;
}

const FamilyTree = ({ horseId }: FamilyTreeProps) => {
	const {
		data: horse,
		isLoading,
		error,
	} = useQuery({
		queryKey: ['horse', horseId],
		queryFn: async () => {
			const response = await fetch(`/api/horses/${horseId}`);
			if (!response.ok) throw new Error('Failed to fetch horse details');
			return response.json();
		},
	});

	if (isLoading) {
		return <LoadingOverlay visible />;
	}

	if (error || !horse) {
		return <Text c='red'>Error loading family tree</Text>;
	}

	return (
		<Paper p='md' bg="dark.8">
			<TreeNode horse={horse} level={0} />
		</Paper>
	);
};

export default FamilyTree;

import { format } from 'date-fns';
import React from 'react';
import { Horse } from '../types/horse';
import { Badge } from './ui/badge';
import { Card, CardContent, CardHeader, CardTitle } from './ui/card';
import {
	Table,
	TableBody,
	TableCell,
	TableHead,
	TableHeader,
	TableRow,
} from './ui/table';

interface BreedingHistoryProps {
	horses: Horse[];
}

export const BreedingHistory: React.FC<BreedingHistoryProps> = ({ horses }) => {
	// Filter horses with breeding history
	const horsesWithBreeding = horses.filter((horse) => horse.breedingDate);

	return (
		<Card className='w-full'>
			<CardHeader>
				<CardTitle>Breeding History</CardTitle>
			</CardHeader>
			<CardContent>
				<Table>
					<TableHeader>
						<TableRow>
							<TableHead>Horse</TableHead>
							<TableHead>Breeding Date</TableHead>
							<TableHead>Due Date</TableHead>
							<TableHead>Status</TableHead>
							<TableHead>Notes</TableHead>
						</TableRow>
					</TableHeader>
					<TableBody>
						{horsesWithBreeding.map((horse) => (
							<TableRow key={horse.id}>
								<TableCell className='font-medium'>
									{horse.name}
								</TableCell>
								<TableCell>
									{horse.breedingDate
										? format(
												new Date(horse.breedingDate),
												'MMM dd, yyyy'
										  )
										: '-'}
								</TableCell>
								<TableCell>
									{horse.breedingDate
										? format(
												new Date(
													horse.breedingDate
												).setDate(
													new Date(
														horse.breedingDate
													).getDate() + 340
												),
												'MMM dd, yyyy'
										  )
										: '-'}
								</TableCell>
								<TableCell>
									<Badge
										variant={
											horse.status === 'pregnant'
												? 'success'
												: 'secondary'
										}
									>
										{horse.status}
									</Badge>
								</TableCell>
								<TableCell className='text-muted-foreground'>
									{horse.notes || '-'}
								</TableCell>
							</TableRow>
						))}
						{horsesWithBreeding.length === 0 && (
							<TableRow>
								<TableCell
									colSpan={5}
									className='text-center text-muted-foreground'
								>
									No breeding history available
								</TableCell>
							</TableRow>
						)}
					</TableBody>
				</Table>
			</CardContent>
		</Card>
	);
};

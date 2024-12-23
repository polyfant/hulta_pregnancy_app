import React from 'react';
import {
    Box,
    Card,
    CardContent,
    Typography,
    Chip,
} from '@mui/material';
import {
    Female as FemaleIcon,
    Male as MaleIcon,
} from '@mui/icons-material';
import { Horse } from '../types/horse';

interface HorsePreviewCardProps {
    horse: Horse | null;
    externalName?: string;
    role: 'mother' | 'father';
}

export const HorsePreviewCard: React.FC<HorsePreviewCardProps> = ({
    horse,
    externalName,
    role,
}) => {
    if (!horse && !externalName) return null;

    const isExternal = !horse && externalName;
    const gender = role === 'mother' ? 'MARE' : 'STALLION';

    return (
        <Card variant="outlined" sx={{ 
            mt: 1,
            backgroundColor: role === 'mother' ? 'rgba(233, 30, 99, 0.08)' : 'rgba(33, 150, 243, 0.08)',
        }}>
            <CardContent>
                <Box display="flex" alignItems="center" gap={1}>
                    {gender === 'MARE' ? (
                        <FemaleIcon color="primary" />
                    ) : (
                        <MaleIcon color="primary" />
                    )}
                    <Typography variant="subtitle1">
                        {isExternal ? externalName : horse?.name}
                    </Typography>
                    {isExternal && (
                        <Chip
                            label="External"
                            size="small"
                            variant="outlined"
                            color="primary"
                        />
                    )}
                </Box>
                {!isExternal && horse && (
                    <Box mt={1}>
                        <Typography variant="body2" color="text.secondary">
                            Breed: {horse.breed}
                        </Typography>
                        {horse.age && (
                            <Typography variant="body2" color="text.secondary">
                                Age: {horse.age}
                            </Typography>
                        )}
                    </Box>
                )}
            </CardContent>
        </Card>
    );
};

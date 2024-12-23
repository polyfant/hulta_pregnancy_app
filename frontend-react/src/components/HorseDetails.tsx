import React, { useState } from 'react';
import {
    Box,
    Card,
    CardContent,
    Typography,
    IconButton,
    Collapse,
    Divider,
} from '@mui/material';
import {
    ExpandMore as ExpandMoreIcon,
    Female as FemaleIcon,
    Male as MaleIcon,
} from '@mui/icons-material';
import { FamilyTree } from './FamilyTree';
import { Horse } from '../types/horse';

interface HorseDetailsProps {
    horse: Horse;
    onNavigateToHorse?: (id: number) => void;
}

export const HorseDetails: React.FC<HorseDetailsProps> = ({ horse, onNavigateToHorse }) => {
    const [showFamilyTree, setShowFamilyTree] = useState(false);

    return (
        <Card>
            <CardContent>
                <Box display="flex" alignItems="center" gap={1}>
                    {horse.gender === 'MARE' ? (
                        <FemaleIcon color="primary" />
                    ) : (
                        <MaleIcon color="primary" />
                    )}
                    <Typography variant="h5" component="div">
                        {horse.name}
                    </Typography>
                </Box>
                
                <Box mt={2}>
                    <Typography variant="body1">
                        <strong>Breed:</strong> {horse.breed}
                    </Typography>
                    <Typography variant="body1">
                        <strong>Age:</strong> {horse.age}
                    </Typography>
                    {horse.weight && (
                        <Typography variant="body1">
                            <strong>Weight:</strong> {horse.weight} kg
                        </Typography>
                    )}
                </Box>

                <Divider sx={{ my: 2 }} />

                <Box display="flex" alignItems="center" gap={1}>
                    <Typography variant="h6">Family Tree</Typography>
                    <IconButton
                        onClick={() => setShowFamilyTree(!showFamilyTree)}
                        sx={{
                            transform: showFamilyTree ? 'rotate(180deg)' : 'none',
                            transition: 'transform 0.3s',
                        }}
                    >
                        <ExpandMoreIcon />
                    </IconButton>
                </Box>

                <Collapse in={showFamilyTree}>
                    <Box mt={2}>
                        <FamilyTree
                            horseId={horse.id}
                            onMemberClick={onNavigateToHorse}
                        />
                    </Box>
                </Collapse>
            </CardContent>
        </Card>
    );
};

import React from 'react';
import {
    Dialog,
    DialogTitle,
    DialogContent,
    DialogActions,
    Button,
    Typography,
    Box,
} from '@mui/material';
import { Horse } from '../types/horse';
import { HorsePreviewCard } from './HorsePreviewCard';

interface ParentChangeDialogProps {
    open: boolean;
    onClose: () => void;
    onConfirm: () => void;
    parentType: 'mother' | 'father';
    currentParent?: Horse | null;
    newParent?: Horse | null;
    currentExternalParent?: string;
    newExternalParent?: string;
}

export const ParentChangeDialog: React.FC<ParentChangeDialogProps> = ({
    open,
    onClose,
    onConfirm,
    parentType,
    currentParent,
    newParent,
    currentExternalParent,
    newExternalParent,
}) => {
    const parentTitle = parentType === 'mother' ? 'Mother' : 'Father';

    return (
        <Dialog open={open} onClose={onClose} maxWidth="sm" fullWidth>
            <DialogTitle>Change {parentTitle}</DialogTitle>
            <DialogContent>
                <Typography variant="body1" gutterBottom>
                    Are you sure you want to change this horse's {parentType}?
                </Typography>

                <Box mt={2}>
                    <Typography variant="subtitle2" color="text.secondary">
                        Current {parentTitle}:
                    </Typography>
                    <HorsePreviewCard
                        horse={currentParent}
                        externalName={currentExternalParent}
                        role={parentType}
                    />
                </Box>

                <Box mt={2}>
                    <Typography variant="subtitle2" color="text.secondary">
                        New {parentTitle}:
                    </Typography>
                    <HorsePreviewCard
                        horse={newParent}
                        externalName={newExternalParent}
                        role={parentType}
                    />
                </Box>
            </DialogContent>
            <DialogActions>
                <Button onClick={onClose} color="inherit">
                    Cancel
                </Button>
                <Button onClick={onConfirm} color="primary" variant="contained">
                    Confirm Change
                </Button>
            </DialogActions>
        </Dialog>
    );
};

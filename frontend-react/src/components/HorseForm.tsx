import React, { useState, useEffect } from 'react';
import {
    Box,
    TextField,
    MenuItem,
    Button,
    FormControl,
    InputLabel,
    Select,
    Typography,
    Autocomplete,
    FormControlLabel,
    Switch,
    Collapse,
} from '@mui/material';
import { DatePicker } from '@mui/x-date-pickers/DatePicker';
import { Horse, CreateHorseInput } from '../types/horse';
import { HorsePreviewCard } from './HorsePreviewCard';
import { ParentChangeDialog } from './ParentChangeDialog';
import dayjs from 'dayjs';

interface HorseFormProps {
    onSubmit: (horse: CreateHorseInput) => void;
    initialValues?: Horse;
}

const validateParentSelection = (availableHorses: Horse[], horseId: number | undefined, parentId: number) => {
    if (!horseId) return { isValid: true };

    const parent = availableHorses.find(h => h.id === parentId);
    if (!parent) return { isValid: true };

    const isOffspring = parent.offspring?.some(o => o.id === horseId);
    if (isOffspring) return { isValid: false, error: 'Cannot select offspring as parent' };

    const isCircular = availableHorses.some(h => h.id === horseId && h.motherId === parentId || h.id === horseId && h.fatherId === parentId);
    if (isCircular) return { isValid: false, error: 'Circular relationship detected' };

    return { isValid: true };
};

export const HorseForm: React.FC<HorseFormProps> = ({ onSubmit, initialValues }) => {
    const [formData, setFormData] = useState<CreateHorseInput>({
        name: '',
        breed: '',
        gender: 'MARE',
        dateOfBirth: dayjs().format('YYYY-MM-DD'),
        weight: undefined,
        motherId: undefined,
        fatherId: undefined,
        externalMother: '',
        externalFather: '',
        ...initialValues,
    });

    const [availableHorses, setAvailableHorses] = useState<Horse[]>([]);
    const [useExternalMother, setUseExternalMother] = useState(!!formData.externalMother);
    const [useExternalFather, setUseExternalFather] = useState(!!formData.externalFather);

    const [dialogState, setDialogState] = useState<{
        open: boolean;
        parentType: 'mother' | 'father';
        currentParent: Horse | null;
        newParent: Horse | null;
        currentExternal: string;
        newExternal: string;
        onConfirm: () => void;
    }>({
        open: false,
        parentType: 'mother',
        currentParent: null,
        newParent: null,
        currentExternal: '',
        newExternal: '',
        onConfirm: () => {},
    });

    const [validationErrors, setValidationErrors] = useState<{
        mother?: string;
        father?: string;
    }>({});

    // Fetch available horses for parent selection
    useEffect(() => {
        const fetchHorses = async () => {
            try {
                const response = await fetch('/api/horses');
                if (!response.ok) throw new Error('Failed to fetch horses');
                const horses = await response.json();
                setAvailableHorses(horses);
            } catch (error) {
                console.error('Error fetching horses:', error);
            }
        };
        fetchHorses();
    }, []);

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        onSubmit(formData);
    };

    const handleDateChange = (date: dayjs.Dayjs | null) => {
        if (date) {
            setFormData(prev => ({
                ...prev,
                dateOfBirth: date.format('YYYY-MM-DD'),
            }));
        }
    };

    // Filter potential mothers (only mares)
    const potentialMothers = availableHorses.filter(
        horse => horse.gender === 'MARE' && horse.id !== initialValues?.id
    );

    // Filter potential fathers (only stallions)
    const potentialFathers = availableHorses.filter(
        horse => horse.gender === 'STALLION' && horse.id !== initialValues?.id
    );

    // Function to handle parent changes with validation
    const handleParentChange = (
        parentType: 'mother' | 'father',
        newParent: Horse | null,
        newExternal: string = ''
    ) => {
        // Clear previous validation errors
        setValidationErrors(prev => ({ ...prev, [parentType]: undefined }));

        if (newParent) {
            const validation = validateParentSelection(
                availableHorses,
                initialValues?.id,
                newParent.id
            );

            if (!validation.isValid) {
                setValidationErrors(prev => ({
                    ...prev,
                    [parentType]: validation.error
                }));
                return;
            }
        }

        if (!initialValues) {
            // If it's a new horse, just update the form
            setFormData(prev => ({
                ...prev,
                [parentType === 'mother' ? 'motherId' : 'fatherId']: newParent?.id,
                [parentType === 'mother' ? 'externalMother' : 'externalFather']: newExternal,
            }));
            return;
        }

        // For existing horses, show confirmation dialog
        const currentParent = parentType === 'mother'
            ? availableHorses.find(h => h.id === initialValues.motherId)
            : availableHorses.find(h => h.id === initialValues.fatherId);

        const currentExternal = parentType === 'mother'
            ? initialValues.externalMother || ''
            : initialValues.externalFather || '';

        setDialogState({
            open: true,
            parentType,
            currentParent: currentParent || null,
            newParent: newParent,
            currentExternal,
            newExternal,
            onConfirm: () => {
                setFormData(prev => ({
                    ...prev,
                    [parentType === 'mother' ? 'motherId' : 'fatherId']: newParent?.id,
                    [parentType === 'mother' ? 'externalMother' : 'externalFather']: newExternal,
                }));
                setDialogState(prev => ({ ...prev, open: false }));
            },
        });
    };

    return (
        <Box component="form" onSubmit={handleSubmit} sx={{ maxWidth: 600, mx: 'auto', p: 2 }}>
            <Typography variant="h5" gutterBottom>
                {initialValues ? 'Edit Horse' : 'Add New Horse'}
            </Typography>

            <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
                <TextField
                    required
                    label="Name"
                    value={formData.name}
                    onChange={(e) => setFormData(prev => ({ ...prev, name: e.target.value }))}
                />

                <TextField
                    required
                    label="Breed"
                    value={formData.breed}
                    onChange={(e) => setFormData(prev => ({ ...prev, breed: e.target.value }))}
                />

                <FormControl fullWidth>
                    <InputLabel>Gender</InputLabel>
                    <Select
                        required
                        value={formData.gender}
                        label="Gender"
                        onChange={(e) => setFormData(prev => ({ ...prev, gender: e.target.value as 'MARE' | 'STALLION' | 'GELDING' }))}
                    >
                        <MenuItem value="MARE">Mare</MenuItem>
                        <MenuItem value="STALLION">Stallion</MenuItem>
                        <MenuItem value="GELDING">Gelding</MenuItem>
                    </Select>
                </FormControl>

                <DatePicker
                    label="Date of Birth"
                    value={dayjs(formData.dateOfBirth)}
                    onChange={handleDateChange}
                />

                <TextField
                    label="Weight (kg)"
                    type="number"
                    value={formData.weight || ''}
                    onChange={(e) => setFormData(prev => ({ ...prev, weight: parseFloat(e.target.value) || undefined }))}
                />

                {/* Mother Selection */}
                <Box>
                    <FormControlLabel
                        control={
                            <Switch
                                checked={useExternalMother}
                                onChange={(e) => {
                                    setUseExternalMother(e.target.checked);
                                    handleParentChange('mother', null, '');
                                }}
                            />
                        }
                        label="External Mother"
                    />

                    <Collapse in={!useExternalMother}>
                        <Autocomplete
                            options={potentialMothers}
                            getOptionLabel={(option) => option.name}
                            value={potentialMothers.find(m => m.id === formData.motherId) || null}
                            onChange={(_, newValue) => {
                                handleParentChange('mother', newValue);
                            }}
                            renderInput={(params) => (
                                <TextField
                                    {...params}
                                    label="Select Mother"
                                    variant="outlined"
                                    error={!!validationErrors.mother}
                                    helperText={validationErrors.mother}
                                />
                            )}
                        />
                        {formData.motherId && (
                            <HorsePreviewCard
                                horse={potentialMothers.find(m => m.id === formData.motherId) || null}
                                role="mother"
                            />
                        )}
                    </Collapse>

                    <Collapse in={useExternalMother}>
                        <TextField
                            fullWidth
                            label="External Mother Name"
                            value={formData.externalMother}
                            onChange={(e) => handleParentChange('mother', null, e.target.value)}
                            sx={{ mt: 2 }}
                        />
                        {formData.externalMother && (
                            <HorsePreviewCard
                                externalName={formData.externalMother}
                                horse={null}
                                role="mother"
                            />
                        )}
                    </Collapse>
                </Box>

                {/* Father Selection */}
                <Box>
                    <FormControlLabel
                        control={
                            <Switch
                                checked={useExternalFather}
                                onChange={(e) => {
                                    setUseExternalFather(e.target.checked);
                                    handleParentChange('father', null, '');
                                }}
                            />
                        }
                        label="External Father"
                    />

                    <Collapse in={!useExternalFather}>
                        <Autocomplete
                            options={potentialFathers}
                            getOptionLabel={(option) => option.name}
                            value={potentialFathers.find(f => f.id === formData.fatherId) || null}
                            onChange={(_, newValue) => {
                                handleParentChange('father', newValue);
                            }}
                            renderInput={(params) => (
                                <TextField
                                    {...params}
                                    label="Select Father"
                                    variant="outlined"
                                    error={!!validationErrors.father}
                                    helperText={validationErrors.father}
                                />
                            )}
                        />
                        {formData.fatherId && (
                            <HorsePreviewCard
                                horse={potentialFathers.find(f => f.id === formData.fatherId) || null}
                                role="father"
                            />
                        )}
                    </Collapse>

                    <Collapse in={useExternalFather}>
                        <TextField
                            fullWidth
                            label="External Father Name"
                            value={formData.externalFather}
                            onChange={(e) => handleParentChange('father', null, e.target.value)}
                            sx={{ mt: 2 }}
                        />
                        {formData.externalFather && (
                            <HorsePreviewCard
                                externalName={formData.externalFather}
                                horse={null}
                                role="father"
                            />
                        )}
                    </Collapse>
                </Box>

                <Button
                    type="submit"
                    variant="contained"
                    color="primary"
                    size="large"
                    sx={{ mt: 2 }}
                >
                    {initialValues ? 'Update Horse' : 'Add Horse'}
                </Button>

                <ParentChangeDialog
                    open={dialogState.open}
                    onClose={() => setDialogState(prev => ({ ...prev, open: false }))}
                    onConfirm={dialogState.onConfirm}
                    parentType={dialogState.parentType}
                    currentParent={dialogState.currentParent}
                    newParent={dialogState.newParent}
                    currentExternalParent={dialogState.currentExternal}
                    newExternalParent={dialogState.newExternal}
                />
            </Box>
        </Box>
    );
};

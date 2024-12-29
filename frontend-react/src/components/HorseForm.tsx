import { useState, useEffect } from 'react';
import {
  TextInput,
  Select,
  Button,
  Stack,
  Group,
  Text,
  Switch,
  NumberInput,
  Collapse,
  Box
} from '@mantine/core';
import { DatePickerInput } from '@mantine/dates';
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

  const isCircular = availableHorses.some(h => 
    (h.id === horseId && h.motherId === parentId) || 
    (h.id === horseId && h.fatherId === parentId)
  );
  if (isCircular) return { isValid: false, error: 'Circular relationship detected' };

  return { isValid: true };
};

export function HorseForm({ onSubmit, initialValues }: HorseFormProps) {
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
    opened: boolean;
    parentType: 'mother' | 'father';
    currentParent: Horse | null;
    newParent: Horse | null;
    currentExternal: string;
    newExternal: string;
    onConfirm: () => void;
  }>({
    opened: false,
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

  const handleMotherChange = (motherId: string | null) => {
    const newMotherId = motherId ? parseInt(motherId) : undefined;
    const validation = validateParentSelection(availableHorses, initialValues?.id, newMotherId!);

    if (!validation.isValid) {
      setValidationErrors(prev => ({ ...prev, mother: validation.error }));
      return;
    }

    setValidationErrors(prev => ({ ...prev, mother: undefined }));
    setFormData(prev => ({ ...prev, motherId: newMotherId }));
  };

  const handleFatherChange = (fatherId: string | null) => {
    const newFatherId = fatherId ? parseInt(fatherId) : undefined;
    const validation = validateParentSelection(availableHorses, initialValues?.id, newFatherId!);

    if (!validation.isValid) {
      setValidationErrors(prev => ({ ...prev, father: validation.error }));
      return;
    }

    setValidationErrors(prev => ({ ...prev, father: undefined }));
    setFormData(prev => ({ ...prev, fatherId: newFatherId }));
  };

  const handleExternalMotherChange = (value: string) => {
    setFormData(prev => ({ ...prev, externalMother: value }));
  };

  const handleExternalFatherChange = (value: string) => {
    setFormData(prev => ({ ...prev, externalFather: value }));
  };

  return (
    <form onSubmit={handleSubmit}>
      <Stack gap="md">
        <TextInput
          label="Name"
          required
          value={formData.name}
          onChange={(e) => setFormData(prev => ({ ...prev, name: e.target.value }))}
        />

        <TextInput
          label="Breed"
          required
          value={formData.breed}
          onChange={(e) => setFormData(prev => ({ ...prev, breed: e.target.value }))}
        />

        <Select
          label="Gender"
          required
          value={formData.gender}
          onChange={(value) => setFormData(prev => ({ ...prev, gender: value as Horse['gender'] }))}
          data={[
            { value: 'MARE', label: 'Mare' },
            { value: 'STALLION', label: 'Stallion' },
            { value: 'GELDING', label: 'Gelding' }
          ]}
        />

        <DatePickerInput
          label="Date of Birth"
          required
          value={dayjs(formData.dateOfBirth).toDate()}
          onChange={(date) => setFormData(prev => ({ 
            ...prev, 
            dateOfBirth: date ? dayjs(date).format('YYYY-MM-DD') : '' 
          }))}
          maxDate={new Date()}
        />

        <NumberInput
          label="Weight (kg)"
          value={formData.weight}
          onChange={(value) => setFormData(prev => ({ ...prev, weight: value }))}
          min={0}
          max={1000}
        />

        <Box>
          <Group mb="xs">
            <Text size="sm" fw={500}>Mother</Text>
            <Switch
              label="External Mother"
              checked={useExternalMother}
              onChange={(event) => setUseExternalMother(event.currentTarget.checked)}
            />
          </Group>

          <Collapse in={!useExternalMother}>
            <Select
              placeholder="Select mother"
              data={availableHorses
                .filter(h => h.gender === 'MARE')
                .map(h => ({ value: h.id.toString(), label: h.name }))}
              value={formData.motherId?.toString()}
              onChange={handleMotherChange}
              error={validationErrors.mother}
              clearable
            />
          </Collapse>

          <Collapse in={useExternalMother}>
            <TextInput
              placeholder="Enter external mother's name"
              value={formData.externalMother}
              onChange={(e) => handleExternalMotherChange(e.target.value)}
            />
          </Collapse>
        </Box>

        <Box>
          <Group mb="xs">
            <Text size="sm" fw={500}>Father</Text>
            <Switch
              label="External Father"
              checked={useExternalFather}
              onChange={(event) => setUseExternalFather(event.currentTarget.checked)}
            />
          </Group>

          <Collapse in={!useExternalFather}>
            <Select
              placeholder="Select father"
              data={availableHorses
                .filter(h => h.gender === 'STALLION')
                .map(h => ({ value: h.id.toString(), label: h.name }))}
              value={formData.fatherId?.toString()}
              onChange={handleFatherChange}
              error={validationErrors.father}
              clearable
            />
          </Collapse>

          <Collapse in={useExternalFather}>
            <TextInput
              placeholder="Enter external father's name"
              value={formData.externalFather}
              onChange={(e) => handleExternalFatherChange(e.target.value)}
            />
          </Collapse>
        </Box>

        <Group justify="flex-end" mt="xl">
          <Button type="submit">
            {initialValues ? 'Update Horse' : 'Add Horse'}
          </Button>
        </Group>
      </Stack>

      <ParentChangeDialog
        opened={dialogState.opened}
        onClose={() => setDialogState(prev => ({ ...prev, opened: false }))}
        onConfirm={dialogState.onConfirm}
        parentType={dialogState.parentType}
        currentParent={dialogState.currentParent}
        newParent={dialogState.newParent}
        currentExternalParent={dialogState.currentExternal}
        newExternalParent={dialogState.newExternal}
      />
    </form>
  );
}

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

  // Check if the potential parent is actually an offspring by checking if the horse is its parent
  const isOffspring = availableHorses.some(h => 
    h.id === parentId && (h.motherId === horseId || h.fatherId === horseId)
  );
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
    isPregnant: false,
    conceptionDate: undefined,
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

  const [showPregnancyFields, setShowPregnancyFields] = useState(formData.gender === 'MARE');

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

  const handleInputChange = (field: keyof CreateHorseInput, value: CreateHorseInput[keyof CreateHorseInput]) => {
    setFormData(prev => ({ ...prev, [field]: value }));
  };

  return (
    <form onSubmit={handleSubmit}>
      <Stack gap="xl">
        <Text size="sm" c="dimmed" mb="md">
          Fields marked with an asterisk (*) are required
        </Text>

        <TextInput
          label="Horse Name"
          description="Enter the horse's registered name"
          placeholder="e.g., Thunder Spirit"
          required
          withAsterisk
          styles={(theme) => ({
            label: {
              marginBottom: '0.2rem',
            },
            required: {
              color: theme.colors.red[6],
              marginLeft: '0.2rem',
            },
            description: {
              marginTop: '0.2rem',
            },
          })}
          value={formData.name}
          onChange={(e) => handleInputChange('name', e.target.value)}
        />

        <TextInput
          label="Breed"
          description="Enter the horse's breed"
          placeholder="e.g., Arabian, Thoroughbred"
          required
          withAsterisk
          styles={(theme) => ({
            label: {
              marginBottom: '0.2rem',
            },
            required: {
              color: theme.colors.red[6],
              marginLeft: '0.2rem',
            },
            description: {
              marginTop: '0.2rem',
            },
          })}
          value={formData.breed}
          onChange={(e) => handleInputChange('breed', e.target.value)}
        />

        <Select
          label="Gender"
          description="Select the horse's gender"
          placeholder="Select gender"
          required
          withAsterisk
          styles={(theme) => ({
            label: {
              marginBottom: '0.2rem',
            },
            required: {
              color: theme.colors.red[6],
              marginLeft: '0.2rem',
            },
            description: {
              marginTop: '0.2rem',
            },
          })}
          value={formData.gender ?? undefined}
          onChange={(value) => {
            handleInputChange('gender', value);
            setShowPregnancyFields(value === 'MARE');
            if (value !== 'MARE') {
              setFormData(prev => ({
                ...prev,
                isPregnant: false,
                conceptionDate: undefined
              }));
            }
          }}
          data={[
            { value: 'MARE', label: 'Mare (Female)' },
            { value: 'STALLION', label: 'Stallion (Male)' },
            { value: 'GELDING', label: 'Gelding (Castrated Male)' },
          ]}
        />

        <DatePickerInput
          label="Date of Birth"
          description="Select the horse's birth date"
          placeholder="Pick a date"
          required
          withAsterisk
          maxDate={new Date()}
          styles={(theme) => ({
            label: {
              marginBottom: '0.2rem',
            },
            required: {
              color: theme.colors.red[6],
              marginLeft: '0.2rem',
            },
            description: {
              marginTop: '0.2rem',
            },
          })}
          value={dayjs(formData.dateOfBirth).toDate()}
          onChange={(date) => handleInputChange('dateOfBirth', date ? dayjs(date).format('YYYY-MM-DD') : '')}
        />

        <NumberInput
          label="Weight"
          description="Enter the horse's weight in kilograms"
          placeholder="e.g., 450"
          suffix=" kg"
          min={0}
          max={1000}
          styles={(theme) => ({
            label: {
              marginBottom: '0.2rem',
            },
            description: {
              marginTop: '0.2rem',
            },
          })}
          value={formData.weight || ''}
          onChange={(value) => handleInputChange('weight', value)}
        />

        {showPregnancyFields && (
          <Box>
            <Text fw={500} size="sm" mb="xs">Pregnancy Information</Text>
            <Stack gap="md">
              <Group align="center" gap="xl">
                <Switch
                  label="Mare is Pregnant"
                  description="Toggle pregnancy status"
                  size="md"
                  labelPosition="left"
                  styles={(theme) => ({
                    root: {
                      width: '100%',
                    },
                    body: {
                      display: 'flex',
                      justifyContent: 'space-between',
                      alignItems: 'center',
                      width: '100%',
                    },
                    labelWrapper: {
                      marginRight: 'auto',
                    },
                  })}
                  checked={formData.isPregnant}
                  onChange={(event) => {
                    handleInputChange('isPregnant', event.currentTarget.checked);
                    if (!event.currentTarget.checked) {
                      handleInputChange('conceptionDate', undefined);
                    }
                  }}
                />
              </Group>

              {formData.isPregnant && (
                <DatePickerInput
                  label="Conception Date"
                  description="Select the date of conception"
                  placeholder="Pick a date"
                  maxDate={new Date()}
                  styles={(theme) => ({
                    label: {
                      marginBottom: '0.2rem',
                    },
                    description: {
                      marginTop: '0.2rem',
                    },
                  })}
                  value={formData.conceptionDate ? dayjs(formData.conceptionDate).toDate() : null}
                  onChange={(date) => handleInputChange('conceptionDate', date ? dayjs(date).format('YYYY-MM-DD') : undefined)}
                />
              )}
            </Stack>
          </Box>
        )}

        <Box>
          <Text fw={500} size="sm" mb="md">Parent Information</Text>
          <Stack gap="xl">
            <Box>
              <Group justify="apart" mb="md">
                <Group gap="xl">
                  <Text size="sm" fw={500}>Mother</Text>
                  <Switch
                    label="External Mother"
                    labelPosition="left"
                    size="md"
                    styles={(theme) => ({
                      root: {
                        display: 'flex',
                        alignItems: 'center',
                      },
                      label: {
                        paddingRight: '1rem',
                      },
                    })}
                    checked={useExternalMother}
                    onChange={(e) => {
                      setUseExternalMother(e.currentTarget.checked);
                      if (e.currentTarget.checked) {
                        handleInputChange('motherId', undefined);
                      } else {
                        handleInputChange('externalMother', '');
                      }
                    }}
                  />
                </Group>
              </Group>
              {useExternalMother ? (
                <TextInput
                  placeholder="Enter external mother's name"
                  value={formData.externalMother || ''}
                  onChange={(e) => handleInputChange('externalMother', e.target.value)}
                />
              ) : (
                <Select
                  placeholder="Select mother from registered horses"
                  data={availableHorses
                    .filter(h => h.gender === 'MARE')
                    .map(h => ({ value: h.id.toString(), label: h.name }))}
                  value={formData.motherId?.toString()}
                  onChange={handleMotherChange}
                  error={validationErrors.mother}
                  clearable
                />
              )}
            </Box>

            <Box>
              <Group justify="apart" mb="md">
                <Group gap="xl">
                  <Text size="sm" fw={500}>Father</Text>
                  <Switch
                    label="External Father"
                    labelPosition="left"
                    size="md"
                    styles={(theme) => ({
                      root: {
                        display: 'flex',
                        alignItems: 'center',
                      },
                      label: {
                        paddingRight: '1rem',
                      },
                    })}
                    checked={useExternalFather}
                    onChange={(e) => {
                      setUseExternalFather(e.currentTarget.checked);
                      if (e.currentTarget.checked) {
                        handleInputChange('fatherId', undefined);
                      } else {
                        handleInputChange('externalFather', '');
                      }
                    }}
                  />
                </Group>
              </Group>
              {useExternalFather ? (
                <TextInput
                  placeholder="Enter external father's name"
                  value={formData.externalFather || ''}
                  onChange={(e) => handleInputChange('externalFather', e.target.value)}
                />
              ) : (
                <Select
                  placeholder="Select father from registered horses"
                  data={availableHorses
                    .filter(h => h.gender === 'STALLION')
                    .map(h => ({ value: h.id.toString(), label: h.name }))}
                  value={formData.fatherId?.toString()}
                  onChange={handleFatherChange}
                  error={validationErrors.father}
                  clearable
                />
              )}
            </Box>
          </Stack>
        </Box>

        <Button 
          type="submit" 
          size="lg"
          variant="filled"
          styles={(theme) => ({
            root: {
              marginTop: theme.spacing.xl,
            },
          })}
        >
          {initialValues ? 'Update Horse' : 'Add Horse'}
        </Button>
      </Stack>
    </form>
  );
}

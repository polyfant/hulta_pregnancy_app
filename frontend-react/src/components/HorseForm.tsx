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
import { useForm, zodResolver } from '@mantine/form';
import { notifications } from '@mantine/notifications';
import { useMutation } from '@tanstack/react-query';
import { horseSchema, type HorseFormValues } from '../validation/horseSchema';
import { Check, X } from '@phosphor-icons/react';
import { LoadingProgress } from './LoadingProgress';

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
  const form = useForm<HorseFormValues>({
    validate: zodResolver(horseSchema),
    initialValues: {
      name: initialValues?.name || '',
      breed: initialValues?.breed || '',
      gender: initialValues?.gender || 'MARE',
      dateOfBirth: initialValues?.dateOfBirth || dayjs().format('YYYY-MM-DD'),
      weight: initialValues?.weight || undefined,
      isPregnant: initialValues?.isPregnant || false,
      conceptionDate: initialValues?.conceptionDate || undefined,
      motherId: initialValues?.motherId || undefined,
      fatherId: initialValues?.fatherId || undefined,
      externalMother: initialValues?.externalMother || '',
      externalFather: initialValues?.externalFather || '',
    },
  });

  const [availableHorses, setAvailableHorses] = useState<Horse[]>([]);
  const [useExternalMother, setUseExternalMother] = useState(!!form.values.externalMother);
  const [useExternalFather, setUseExternalFather] = useState(!!form.values.externalFather);

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

  const [showPregnancyFields, setShowPregnancyFields] = useState(form.values.gender === 'MARE');

  const mutation = useMutation({
    mutationFn: (values: HorseFormValues) => {
      return apiClient.post('/api/horses', values);
    },
    onSuccess: () => {
      notifications.show({
        title: 'Success',
        message: 'Horse added successfully',
        color: 'green',
        icon: <Check weight="bold" />,
        autoClose: 3000,
      });
      navigate('/horses');
    },
    onError: (error) => {
      notifications.show({
        title: 'Error',
        message: 'Failed to add horse. Please try again.',
        color: 'red',
        icon: <X weight="bold" />,
        autoClose: 5000,
      });
    },
  });

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

  const handleSubmit = form.onSubmit((values) => {
    mutation.mutate(values);
  });

  const handleMotherChange = (motherId: string | null) => {
    const newMotherId = motherId ? parseInt(motherId) : undefined;
    const validation = validateParentSelection(availableHorses, initialValues?.id, newMotherId!);

    if (!validation.isValid) {
      setValidationErrors(prev => ({ ...prev, mother: validation.error }));
      return;
    }

    setValidationErrors(prev => ({ ...prev, mother: undefined }));
    form.setFieldValue('motherId', newMotherId);
  };

  const handleFatherChange = (fatherId: string | null) => {
    const newFatherId = fatherId ? parseInt(fatherId) : undefined;
    const validation = validateParentSelection(availableHorses, initialValues?.id, newFatherId!);

    if (!validation.isValid) {
      setValidationErrors(prev => ({ ...prev, father: validation.error }));
      return;
    }

    setValidationErrors(prev => ({ ...prev, father: undefined }));
    form.setFieldValue('fatherId', newFatherId);
  };

  const handleExternalMotherChange = (value: string) => {
    form.setFieldValue('externalMother', value);
  };

  const handleExternalFatherChange = (value: string) => {
    form.setFieldValue('externalFather', value);
  };

  return (
    <form onSubmit={handleSubmit}>
      <Stack spacing="md">
        {mutation.isPending && (
          <LoadingProgress message="Saving horse details..." />
        )}

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
          {...form.getInputProps('name')}
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
          {...form.getInputProps('breed')}
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
          {...form.getInputProps('gender')}
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
          {...form.getInputProps('dateOfBirth')}
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
          {...form.getInputProps('weight')}
        />

        {showPregnancyFields && (
          <Box>
            <Text fw={500} size="sm" mb="xs">Pregnancy Information</Text>
            <Stack gap="md">
              <Switch
                label="Mare is Pregnant"
                description="Toggle pregnancy status"
                size="md"
                labelPosition="left"
                styles={(theme) => ({
                  root: {
                    width: '100%'
                  },
                  body: {
                    display: 'flex',
                    justifyContent: 'space-between',
                    width: '100%'
                  }
                })}
                {...form.getInputProps('isPregnant', { type: 'checkbox' })}
              />
              {form.values.isPregnant && (
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
                  {...form.getInputProps('conceptionDate')}
                />
              )}
            </Stack>
          </Box>
        )}

        <Box>
          <Text fw={500} size="sm" mb="xs">Parent Information</Text>
          <Stack gap="md">
            <Switch
              label="External Mother"
              labelPosition="left"
              size="md"
              styles={(theme) => ({
                root: {
                  width: '100%'
                },
                body: {
                  display: 'flex',
                  justifyContent: 'space-between',
                  width: '100%'
                }
              })}
              {...form.getInputProps('useExternalMother', { type: 'checkbox' })}
            />

            {form.values.useExternalMother ? (
              <TextInput
                placeholder="Enter external mother's name"
                {...form.getInputProps('externalMother')}
              />
            ) : (
              <Select
                placeholder="Select mother from registered horses"
                data={availableHorses
                  .filter(h => h.gender === 'MARE')
                  .map(h => ({ value: h.id.toString(), label: h.name }))}
                {...form.getInputProps('motherId')}
                onChange={handleMotherChange}
                error={validationErrors.mother}
                clearable
              />
            )}

            <Switch
              label="External Father"
              labelPosition="left"
              size="md"
              styles={(theme) => ({
                root: {
                  width: '100%'
                },
                body: {
                  display: 'flex',
                  justifyContent: 'space-between',
                  width: '100%'
                }
              })}
              {...form.getInputProps('useExternalFather', { type: 'checkbox' })}
            />

            {form.values.useExternalFather ? (
              <TextInput
                placeholder="Enter external father's name"
                {...form.getInputProps('externalFather')}
              />
            ) : (
              <Select
                placeholder="Select father from registered horses"
                data={availableHorses
                  .filter(h => h.gender === 'STALLION')
                  .map(h => ({ value: h.id.toString(), label: h.name }))}
                {...form.getInputProps('fatherId')}
                onChange={handleFatherChange}
                error={validationErrors.father}
                clearable
              />
            )}
          </Stack>
        </Box>

        <Button 
          type="submit"
          loading={mutation.isPending}
          disabled={mutation.isPending}
        >
          {initialValues ? 'Update Horse' : 'Add Horse'}
        </Button>
      </Stack>
    </form>
  );
}

import { useForm } from '@mantine/form';
import { TextInput, Select, Button, Stack } from '@mantine/core';
import { DateInput } from '@mantine/dates';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { horsesApi } from '../api/horses';
import { CreateHorseInput, Horse } from '../types/horse';

interface AddHorseFormProps {
    initialValues?: Horse;
    onSubmit?: (values: CreateHorseInput) => void;
    submitButtonText?: string;
}

export function AddHorseForm({ initialValues, onSubmit, submitButtonText = 'Add Horse' }: AddHorseFormProps) {
    const queryClient = useQueryClient();
    const form = useForm<CreateHorseInput>({
        initialValues: initialValues ?? {
            name: '',
            breed: '',
            gender: 'MARE',
            dateOfBirth: '',
        },
        validate: {
            name: (value) => (value.length < 2 ? 'Name must be at least 2 characters' : null),
            breed: (value) => (value.length < 2 ? 'Breed must be at least 2 characters' : null),
        },
    });

    const mutation = useMutation({
        mutationFn: horsesApi.create,
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['horses'] });
            form.reset();
        },
    });

    const handleSubmit = (values: CreateHorseInput) => {
        if (onSubmit) {
            onSubmit(values);
        } else {
            mutation.mutate(values);
        }
    };

    return (
        <form onSubmit={form.onSubmit(handleSubmit)}>
            <Stack>
                <TextInput
                    label="Name"
                    placeholder="Horse name"
                    required
                    {...form.getInputProps('name')}
                />
                <TextInput
                    label="Breed"
                    placeholder="Horse breed"
                    required
                    {...form.getInputProps('breed')}
                />
                <Select
                    label="Gender"
                    required
                    data={[
                        { value: 'MARE', label: 'Mare' },
                        { value: 'STALLION', label: 'Stallion' },
                        { value: 'GELDING', label: 'Gelding' },
                    ]}
                    {...form.getInputProps('gender')}
                />
                <DateInput
                    label="Date of Birth"
                    placeholder="Pick a date"
                    {...form.getInputProps('dateOfBirth')}
                />
                <Button type="submit" loading={mutation.isPending}>
                    {submitButtonText}
                </Button>
            </Stack>
        </form>
    );
}

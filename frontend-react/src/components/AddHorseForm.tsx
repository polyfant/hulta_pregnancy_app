import { useForm } from '@mantine/form';
import { TextInput, Select, Button, Stack } from '@mantine/core';
import { DateInput } from '@mantine/dates';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { horsesApi } from '../api/horses';
import { CreateHorseInput } from '../types/horse';

export function AddHorseForm() {
    const queryClient = useQueryClient();
    const form = useForm<CreateHorseInput>({
        initialValues: {
            name: '',
            breed: '',
            gender: 'Mare',
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

    return (
        <form onSubmit={form.onSubmit((values) => mutation.mutate(values))}>
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
                        { value: 'Mare', label: 'Mare' },
                        { value: 'Stallion', label: 'Stallion' },
                        { value: 'Gelding', label: 'Gelding' },
                    ]}
                    {...form.getInputProps('gender')}
                />
                <DateInput
                    label="Date of Birth"
                    placeholder="Pick a date"
                    {...form.getInputProps('dateOfBirth')}
                />
                <Button type="submit" loading={mutation.isPending}>
                    Add Horse
                </Button>
            </Stack>
        </form>
    );
}

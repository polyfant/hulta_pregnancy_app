import { z } from 'zod';

export const horseSchema = z.object({
    name: z.string()
        .min(2, 'Name must be at least 2 characters')
        .max(50, 'Name must be less than 50 characters'),
    breed: z.string().optional(),
    color: z.string().optional(),
    gender: z.enum(['MARE', 'STALLION', 'GELDING'], {
        required_error: 'Please select a gender',
    }),
    birthDate: z.string().optional(),
    isPregnant: z.boolean().optional(),
});

export type HorseFormValues = z.infer<typeof horseSchema>; 
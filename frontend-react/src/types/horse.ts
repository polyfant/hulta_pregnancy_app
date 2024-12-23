export interface Horse {
    id: number;
    name: string;
    breed: string;
    gender: 'Mare' | 'Stallion' | 'Gelding';
    dateOfBirth?: string;
    created_at?: string;
    updated_at?: string;
}

export interface CreateHorseInput {
    name: string;
    breed: string;
    gender: 'Mare' | 'Stallion' | 'Gelding';
    dateOfBirth?: string;
}

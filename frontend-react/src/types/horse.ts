export interface Horse {
    id: string;
    name: string;
    breed?: string;
    color?: string;
    gender: 'MARE' | 'STALLION' | 'GELDING';
    birthDate?: string;
    isPregnant?: boolean;
    dateOfBirth: string;
    weight?: number;
    age?: string;
    conceptionDate?: string;
    motherId?: number;
    fatherId?: number;
    externalMother?: string;
    externalFather?: string;
    created_at?: string;
    updated_at?: string;
}

export interface CreateHorseInput {
    name: string;
    breed: string;
    gender: 'MARE' | 'STALLION' | 'GELDING';
    dateOfBirth: string;
    weight?: number;
    isPregnant?: boolean;
    conceptionDate?: string;
    motherId?: number;
    fatherId?: number;
    externalMother?: string;
    externalFather?: string;
}

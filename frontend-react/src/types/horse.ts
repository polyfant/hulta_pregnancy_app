export interface Horse {
    id: number;
    name: string;
    breed: string;
    gender: 'MARE' | 'STALLION' | 'GELDING';
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
    conceptionDate?: string;
    motherId?: number;
    fatherId?: number;
    externalMother?: string;
    externalFather?: string;
}

export interface FamilyMember {
    id?: number;
    name: string;
    breed?: string;
    gender?: string;
    dateOfBirth?: string;
    age?: string;
    isExternal: boolean;
    externalSource?: string;
}

export interface FamilyTree {
    horse: {
        id: number;
        name: string;
        breed: string;
        gender: string;
        dateOfBirth: string;
        age: string;
    };
    mother?: FamilyMember;
    father?: FamilyMember;
    offspring?: FamilyMember[];
    siblings?: FamilyMember[];
}

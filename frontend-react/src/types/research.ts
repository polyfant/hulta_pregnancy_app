export interface ResearchData {
    source: 'USDA' | 'NAHMS' | 'EquineScience';
    datasets: {
        pregnancy: {
            url: string;
            description: string;
            variables: string[];
        };
        foal: {
            url: string;
            description: string;
            variables: string[];
        };
    };
}

export const EQUINE_DATABASES = {
    NAHMS: 'https://www.aphis.usda.gov/aphis/ourfocus/animalhealth/monitoring-and-surveillance/nahms/nahms_equine_studies',
    UCD: 'https://www.vetmed.ucdavis.edu/research/equine-research',
    AAEP: 'https://aaep.org/research',
} as const; 
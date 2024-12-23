import { Horse } from '../types/horse';

interface HorseWithParents extends Horse {
    motherId?: number | null;
    fatherId?: number | null;
}

export const isCircularRelationship = (
    horses: HorseWithParents[],
    horseId: number,
    potentialParentId: number,
    maxDepth: number = 10
): boolean => {
    if (maxDepth === 0) return false;
    if (horseId === potentialParentId) return true;

    const parent = horses.find(h => h.id === potentialParentId);
    if (!parent) return false;

    if (parent.motherId) {
        if (isCircularRelationship(horses, horseId, parent.motherId, maxDepth - 1)) {
            return true;
        }
    }

    if (parent.fatherId) {
        if (isCircularRelationship(horses, horseId, parent.fatherId, maxDepth - 1)) {
            return true;
        }
    }

    return false;
};

export const isOffspring = (
    horses: HorseWithParents[],
    horseId: number,
    potentialParentId: number
): boolean => {
    const potentialParent = horses.find(h => h.id === potentialParentId);
    if (!potentialParent) return false;

    // Check if the horse is a direct offspring
    return horses.some(h => 
        h.id === horseId && 
        (h.motherId === potentialParentId || h.fatherId === potentialParentId)
    );
};

export const validateParentSelection = (
    horses: HorseWithParents[],
    horseId: number | undefined,
    potentialParentId: number
): { isValid: boolean; error?: string } => {
    if (!horseId) return { isValid: true }; // New horse, no validation needed yet

    // Check if trying to set itself as parent
    if (horseId === potentialParentId) {
        return {
            isValid: false,
            error: "A horse cannot be its own parent"
        };
    }

    // Check for circular relationships
    if (isCircularRelationship(horses, horseId, potentialParentId)) {
        return {
            isValid: false,
            error: "This would create a circular family relationship"
        };
    }

    // Check if selecting offspring as parent
    if (isOffspring(horses, potentialParentId, horseId)) {
        return {
            isValid: false,
            error: "Cannot select offspring as parent"
        };
    }

    return { isValid: true };
};

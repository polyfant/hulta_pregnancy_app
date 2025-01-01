# Icon Usage Guide

## Icon Library: Phosphor Icons

We use Phosphor Icons (@phosphor-icons/react) as our primary icon library for its:

-   Modern, consistent design
-   Tree-shakeable imports
-   Comprehensive icon set
-   TypeScript support
-   Active maintenance

## Implementation

All icons are centralized in `src/utils/icons.ts` with consistent naming:

```typescript
import { Plus, Pencil } from '@phosphor-icons/react';

export {
	Plus as FiPlus, // For adding items
	Pencil as FiEdit, // For editing items
	// etc...
};
```

## Common Icons and Their Uses

-   `FiPlus` - Adding new items
-   `FiEdit` - Editing existing items
-   `FiTrash2` - Deleting items
-   `FiSearch` - Search functionality
-   `FiPets` (Horse) - Horse-related actions
-   `FiBaby` - Pregnancy tracking
-   `FiHeart` - Health tracking
-   `FiCalendar` - Date-related actions
-   `FiMale`/`FiFemale` - Gender indicators

## Usage in Components

```typescript
import { FiPlus, FiEdit } from '@/utils/icons';

function MyComponent() {
	return <Button leftSection={<FiPlus size={16} />}>Add New</Button>;
}
```

## Icon Properties

All icons accept these common props:

-   `size`: number | string
-   `color`: string
-   `weight`: 'thin' | 'light' | 'regular' | 'bold' | 'fill'
-   `mirrored`: boolean

## Migration Note

Previously using Feather Icons (react-icons/fi), migrated to Phosphor Icons for:

-   Better tree-shaking
-   More comprehensive icon set
-   Better TypeScript support
-   Active maintenance

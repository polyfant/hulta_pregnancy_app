# Icon Usage Guide

## Icon Library: Phosphor Icons

We use Phosphor Icons (@phosphor-icons/react) as our primary icon library for:

- Modern, consistent design
- Tree-shakeable imports
- Comprehensive icon set
- TypeScript support
- Active maintenance

## Implementation

All icons are centralized in `src/utils/icons.ts`:

```typescript
import {
	Activity,
	Baby,
	Calendar,
	CaretDown,
	CaretRight,
	Check,
	Clock,
	GenderFemale,
	GenderMale,
	Heart,
	Horse,
	Info,
	MagnifyingGlass,
	Plus,
	Syringe,
	Trash,
	Warning,
	X
} from '@phosphor-icons/react';

export {
	Activity,          // Health activity tracking
	Baby,             // Pregnancy/foal related
	Calendar,         // Date selection/events
	CaretDown,        // Expandable sections
	CaretRight,       // Navigation indicators
	Check,            // Completion/success
	Clock,            // Time-related
	GenderFemale,     // Mare indicator
	GenderMale,       // Stallion/Gelding indicator
	Heart,            // Health status
	Horse,            // General horse actions
	Info,             // Information/help
	MagnifyingGlass,  // Search functionality
	Plus,            // Add new items
	Syringe,         // Medical/vaccination
	Trash,           // Delete actions
	Warning,         // Alert/warning states
	X               // Close/cancel actions
};
```

## Common Use Cases

### Horse Management
- `Horse` - Horse profile/details
- `GenderFemale`/`GenderMale` - Gender indicators
- `Plus` - Add new horse
- `Trash` - Delete horse

### Health Tracking
- `Activity` - Health records
- `Heart` - Health status
- `Syringe` - Vaccinations

### Pregnancy Monitoring
- `Baby` - Pregnancy tracking
- `Calendar` - Due dates
- `Clock` - Timing/schedules

### UI Elements
- `CaretDown`/`CaretRight` - Expandable sections
- `Check`/`X` - Success/cancel actions
- `Info` - Help tooltips
- `MagnifyingGlass` - Search
- `Warning` - Alert messages

## Usage Example

```typescript
import { GenderFemale, GenderMale, Horse } from '@phosphor-icons/react';

function HorseCard({ horse }) {
	return (
		<Card>
			<Group>
				<Horse size="1.2rem" />
				<Text>{horse.name}</Text>
				{horse.gender === 'MARE' ? (
					<GenderFemale size="1.2rem" color="var(--mantine-color-pink-6)" />
				) : (
					<GenderMale size="1.2rem" color="var(--mantine-color-blue-6)" />
				)}
			</Group>
		</Card>
	);
}
```

## Icon Properties

All icons accept these common props:

- `size`: number | string
- `color`: string
- `weight`: 'thin' | 'light' | 'regular' | 'bold' | 'fill'
- `mirrored`: boolean

## Best Practices

1. Always import icons from `@phosphor-icons/react`
2. Use consistent sizes within similar contexts
3. Use semantic colors from the theme
4. Consider accessibility when using icons alone
5. Add tooltips for icon-only buttons

import * as React from 'react';
import { cn } from '../../lib/utils';

export interface SelectProps
	extends React.SelectHTMLAttributes<HTMLSelectElement> {}

export const Select = React.forwardRef<HTMLSelectElement, SelectProps>(
	({ className, children, ...props }, ref) => {
		return (
			<select
				className={cn(
					'flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50',
					className
				)}
				ref={ref}
				{...props}
			>
				{children}
			</select>
		);
	}
);
Select.displayName = 'Select';

// Add a SelectAdapter component that supports onValueChange
export interface SelectAdapterProps extends Omit<SelectProps, 'onChange'> {
	onValueChange?: (value: string) => void;
}

export const SelectAdapter = React.forwardRef<
	HTMLSelectElement,
	SelectAdapterProps
>(({ onValueChange, ...props }, ref) => {
	const handleChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
		onValueChange?.(e.target.value);
	};

	return <Select onChange={handleChange} ref={ref} {...props} />;
});
SelectAdapter.displayName = 'SelectAdapter';

export interface SelectTriggerProps
	extends React.HTMLAttributes<HTMLDivElement> {}

export const SelectTrigger = React.forwardRef<
	HTMLDivElement,
	SelectTriggerProps
>(({ className, children, ...props }, ref) => {
	return (
		<div
			ref={ref}
			className={cn(
				'flex h-10 w-full items-center justify-between rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50',
				className
			)}
			{...props}
		>
			{children}
		</div>
	);
});
SelectTrigger.displayName = 'SelectTrigger';

export interface SelectValueProps {
	placeholder?: string;
	children?: React.ReactNode;
}

export const SelectValue: React.FC<SelectValueProps> = ({
	placeholder,
	children,
}) => {
	return <span>{children || placeholder}</span>;
};

export interface SelectContentProps
	extends React.HTMLAttributes<HTMLDivElement> {}

export const SelectContent = React.forwardRef<
	HTMLDivElement,
	SelectContentProps
>(({ className, children, ...props }, ref) => {
	return (
		<div
			ref={ref}
			className={cn(
				'absolute z-50 min-w-[8rem] overflow-hidden rounded-md border border-input bg-popover text-popover-foreground shadow-md animate-in fade-in-80',
				className
			)}
			{...props}
		>
			<div className='p-1'>{children}</div>
		</div>
	);
});
SelectContent.displayName = 'SelectContent';

export interface SelectItemProps extends React.HTMLAttributes<HTMLDivElement> {
	value: string;
}

export const SelectItem = React.forwardRef<HTMLDivElement, SelectItemProps>(
	({ className, children, ...props }, ref) => {
		return (
			<div
				ref={ref}
				className={cn(
					'relative flex w-full cursor-default select-none items-center rounded-sm py-1.5 pl-8 pr-2 text-sm outline-none focus:bg-accent focus:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50',
					className
				)}
				{...props}
			>
				{children}
			</div>
		);
	}
);
SelectItem.displayName = 'SelectItem';

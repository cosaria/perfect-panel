import { type VariantProps, cva } from 'class-variance-authority';
import React, { forwardRef, type HTMLAttributes } from 'react';

import { cn } from '../../lib/utils';

const sidebarVariants = cva(
	'relative flex h-full min-h-0 flex-col overflow-hidden border-r bg-background text-foreground shadow-sm',
	{
		variants: {
			variant: {
				default: 'w-64',
				compact: 'w-56',
				wide: 'w-72',
			},
		},
		defaultVariants: {
			variant: 'default',
		},
	},
);

export interface SidebarProps
	extends HTMLAttributes<HTMLElement>,
		VariantProps<typeof sidebarVariants> {}

export const Sidebar = forwardRef<HTMLElement, SidebarProps>(
	({ className, variant, ...props }, ref) => {
		return (
			<aside
				ref={ref}
				data-slot="sidebar"
				data-variant={variant ?? 'default'}
				className={cn(sidebarVariants({ variant }), className)}
				{...props}
			/>
		);
	},
);

Sidebar.displayName = 'Sidebar';

import React, { forwardRef, type HTMLAttributes } from 'react';

import { cn } from '../../lib/utils';

export interface SidebarProps extends HTMLAttributes<HTMLElement> {}

export const Sidebar = forwardRef<HTMLElement, SidebarProps>(({ className, ...props }, ref) => {
	return (
		<aside
			ref={ref}
			data-slot="sidebar"
			className={cn(
				'flex h-full w-64 flex-col border-r border-border bg-background text-foreground',
				className,
			)}
			{...props}
		/>
	);
});

Sidebar.displayName = 'Sidebar';

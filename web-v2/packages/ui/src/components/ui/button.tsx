import React, { forwardRef, type ButtonHTMLAttributes } from 'react';

import { cn } from '../../lib/utils';

export interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {}

export const Button = forwardRef<HTMLButtonElement, ButtonProps>(
	({ className, type = 'button', ...props }, ref) => {
		return (
			<button
				ref={ref}
				type={type}
				className={cn(
					'inline-flex items-center justify-center rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground transition-colors hover:opacity-90 disabled:pointer-events-none disabled:opacity-50',
					className,
				)}
				{...props}
			/>
		);
	},
);

Button.displayName = 'Button';

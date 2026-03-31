'use client';

import { ExternalLink } from 'lucide-react';
import React from 'react';

interface IpLinkProps {
  ip: string;
  children?: React.ReactNode;
  className?: string;
  target?: '_blank' | '_self';
}

export function IpLink({ ip, children, className = '', target = '_blank' }: IpLinkProps) {
  const url = `https://ipinfo.io/${ip}`;

  return (
    <a
      href={url}
      target={target}
      rel={target === '_blank' ? 'noopener noreferrer' : undefined}
      className={`text-primary hover:text-primary/80 inline-flex items-center gap-1 font-mono transition-colors hover:underline ${className}`}
    >
      {children || ip}
      <ExternalLink className='h-3 w-3' />
    </a>
  );
}

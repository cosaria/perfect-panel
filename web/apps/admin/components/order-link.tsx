import { Button } from '@workspace/ui/components/button';
import Link from 'next/link';

interface OrderLinkProps {
  orderId?: string | number;
  className?: string;
}

export function OrderLink({ orderId, className }: OrderLinkProps) {
  if (!orderId) return <span>--</span>;

  return (
    <Button variant='link' className={`p-0 ${className || ''}`} asChild>
      <Link href={`/dashboard/order?search=${orderId}`}>{orderId}</Link>
    </Button>
  );
}

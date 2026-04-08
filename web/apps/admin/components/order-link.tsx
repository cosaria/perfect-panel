import { Button } from "@workspace/ui/components/button";
import { AdminLink } from "./admin-link";

interface OrderLinkProps {
  orderId?: string | number;
  className?: string;
}

export function OrderLink({ orderId, className }: OrderLinkProps) {
  if (!orderId) return <span>--</span>;

  return (
    <Button variant="link" className={`p-0 ${className || ""}`} asChild>
      <AdminLink href={`/dashboard/order?search=${orderId}`}>{orderId}</AdminLink>
    </Button>
  );
}

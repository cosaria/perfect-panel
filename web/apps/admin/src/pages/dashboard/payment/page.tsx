import Billing from "@/components/billing";
import PaymentTable from "./payment-table";

export default function Page() {
  return (
    <>
      <PaymentTable />
      <div className="mt-5 flex flex-col gap-3">
        <Billing type="payment" />
      </div>
    </>
  );
}

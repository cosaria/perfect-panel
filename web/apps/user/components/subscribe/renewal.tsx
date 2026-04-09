"use client";

import { useQuery } from "@tanstack/react-query";
import { Button } from "@workspace/ui/components/button";
import { Card, CardContent } from "@workspace/ui/components/card";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@workspace/ui/components/dialog";
import { Separator } from "@workspace/ui/components/separator";
import { LoaderCircle } from "lucide-react";
import { useRouter } from "@/src/compat/app-navigation";
import { useTranslations } from "@workspace/ui/components/i18n-provider";
import { useCallback, useEffect, useRef, useState, useTransition } from "react";
import CouponInput from "@/components/subscribe/coupon-input";
import DurationSelector from "@/components/subscribe/duration-selector";
import PaymentMethods from "@/components/subscribe/payment-methods";
import useGlobalStore from "@/config/use-global";
import { preCreateOrder, renewal } from "@/services/user-api/sdk.gen";
import type {
  PurchaseOrderRequest,
  RenewalOrderRequest,
  Subscribe,
} from "@/services/user-api/types.gen";
import { SubscribeBilling } from "./billing";
import { SubscribeDetail } from "./detail";

interface RenewalProps {
  id: number;
  subscribe: Subscribe;
}

type PreCreateOrderData = Awaited<ReturnType<typeof preCreateOrder>>["data"];

export default function Renewal({ id, subscribe }: Readonly<RenewalProps>) {
  const t = useTranslations("subscribe");
  const { getUserInfo } = useGlobalStore();
  const [open, setOpen] = useState<boolean>(false);
  const router = useRouter();
  const [params, setParams] = useState<Partial<RenewalOrderRequest>>({
    quantity: 1,
    payment: -1,
    coupon: "",
    user_subscribe_id: id,
  });
  const [loading, startTransition] = useTransition();
  const lastSuccessOrderRef = useRef<PreCreateOrderData | null>(null);

  const { data: order } = useQuery({
    enabled: !!subscribe.id && open,
    queryKey: ["preCreateOrder", params],
    queryFn: async () => {
      try {
        const { data } = await preCreateOrder({
          body: {
            ...params,
            subscribe_id: subscribe.id,
          } as PurchaseOrderRequest,
        });
        const result = data;
        if (result) {
          lastSuccessOrderRef.current = result;
        }
        return result;
      } catch (_error) {
        if (lastSuccessOrderRef.current) {
          return lastSuccessOrderRef.current;
        }
      }
    },
  });

  useEffect(() => {
    if (subscribe.id && id) {
      setParams((prev) => ({
        ...prev,
        quantity: 1,
        subscribe_id: subscribe.id,
        user_subscribe_id: id,
      }));
    }
  }, [subscribe.id, id]);

  const handleChange = useCallback((field: keyof typeof params, value: string | number) => {
    setParams((prev) => ({
      ...prev,
      [field]: value,
    }));
  }, []);

  const handleSubmit = useCallback(async () => {
    startTransition(async () => {
      try {
        const { data: response } = await renewal({ body: params as RenewalOrderRequest });
        const orderNo = response?.order_no;
        if (orderNo) {
          getUserInfo();
          router.push(`/payment?order_no=${orderNo}`);
        }
      } catch (_error) {
        /* empty */
      }
    });
  }, [params, router, getUserInfo]);

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button size="sm">{t("renew")}</Button>
      </DialogTrigger>
      <DialogContent className="flex h-full max-w-screen-lg flex-col overflow-hidden md:h-auto">
        <DialogHeader>
          <DialogTitle>{t("renewSubscription")}</DialogTitle>
        </DialogHeader>
        <div className="grid w-full gap-3 lg:grid-cols-2">
          <Card className="border-transparent shadow-none md:border-inherit md:shadow">
            <CardContent className="grid gap-3 p-0 text-sm md:p-6">
              <SubscribeDetail
                subscribe={{
                  ...subscribe,
                  quantity: params.quantity,
                }}
              />
              <Separator />
              <SubscribeBilling
                order={{
                  ...order,
                  quantity: params.quantity,
                  unit_price: subscribe?.unit_price,
                }}
              />
            </CardContent>
          </Card>
          <div className="flex flex-col justify-between text-sm">
            <div className="mb-6 grid gap-3">
              <DurationSelector
                quantity={params.quantity ?? 1}
                unitTime={subscribe?.unit_time}
                discounts={subscribe?.discount ?? undefined}
                onChange={(value) => {
                  handleChange("quantity", value);
                }}
              />
              <CouponInput
                coupon={params.coupon}
                onChange={(value) => handleChange("coupon", value)}
              />
              <PaymentMethods
                value={params.payment ?? -1}
                onChange={(value) => {
                  handleChange("payment", value);
                }}
              />
            </div>
            <Button
              className="fixed bottom-0 left-0 w-full rounded-none md:relative md:mt-6"
              disabled={loading}
              onClick={handleSubmit}
            >
              {loading && <LoaderCircle className="mr-2 animate-spin" />}
              {t("buyNow")}
            </Button>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
}

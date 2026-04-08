import { Card, CardContent, CardHeader, CardTitle } from "@workspace/ui/components/card";
import { type ComponentType, type LazyExoticComponent, lazy, Suspense } from "react";
import { Link, Navigate, Outlet } from "react-router-dom";
import DashboardLayout from "@/app/dashboard/layout";
import AppShell from "./app-shell";

const AuthPage = lazy(() => import("@/app/(auth)/page"));
const AnnouncementPage = lazy(() => import("@/app/dashboard/announcement/page"));
const AuthControlPage = lazy(() => import("@/app/dashboard/auth-control/page"));
const CouponPage = lazy(() => import("@/app/dashboard/coupon/page"));
const DashboardPage = lazy(() => import("@/app/dashboard/page"));
const DocumentPage = lazy(() => import("@/app/dashboard/document/page"));
const BalanceLogPage = lazy(() => import("@/app/dashboard/log/balance/page"));
const CommissionLogPage = lazy(() => import("@/app/dashboard/log/commission/page"));
const EmailLogPage = lazy(() => import("@/app/dashboard/log/email/page"));
const GiftLogPage = lazy(() => import("@/app/dashboard/log/gift/page"));
const LoginLogPage = lazy(() => import("@/app/dashboard/log/login/page"));
const MobileLogPage = lazy(() => import("@/app/dashboard/log/mobile/page"));
const RegisterLogPage = lazy(() => import("@/app/dashboard/log/register/page"));
const ResetSubscribeLogPage = lazy(() => import("@/app/dashboard/log/reset-subscribe/page"));
const ServerTrafficLogPage = lazy(() => import("@/app/dashboard/log/server-traffic/page"));
const SubscribeLogPage = lazy(() => import("@/app/dashboard/log/subscribe/page"));
const SubscribeTrafficLogPage = lazy(() => import("@/app/dashboard/log/subscribe-traffic/page"));
const TrafficDetailsLogPage = lazy(() => import("@/app/dashboard/log/traffic-details/page"));
const MarketingPage = lazy(() => import("@/app/dashboard/marketing/page"));
const NodesPage = lazy(() => import("@/app/dashboard/nodes/page"));
const OrderPage = lazy(() => import("@/app/dashboard/order/page"));
const PaymentPage = lazy(() => import("@/app/dashboard/payment/page"));
const ProductPage = lazy(() => import("@/app/dashboard/product/page"));
const ServersPage = lazy(() => import("@/app/dashboard/servers/page"));
const SubscribePage = lazy(() => import("@/app/dashboard/subscribe/page"));
const SystemPage = lazy(() => import("@/app/dashboard/system/page"));
const TicketPage = lazy(() => import("@/app/dashboard/ticket/page"));
const UserPage = lazy(() => import("@/app/dashboard/user/page"));

function DashboardShell() {
  return (
    <DashboardLayout>
      <Outlet />
    </DashboardLayout>
  );
}

function RoutePlaceholder({
  title,
  description,
  actionLabel,
  actionTo,
}: {
  title: string;
  description: string;
  actionLabel?: string;
  actionTo?: string;
}) {
  return (
    <div className="mx-auto flex h-full max-w-3xl items-center justify-center">
      <Card className="w-full">
        <CardHeader>
          <CardTitle>{title}</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4 text-sm">
          <p className="text-muted-foreground">{description}</p>
          {actionLabel && actionTo ? (
            <Link
              className="text-primary inline-flex font-medium underline-offset-4 hover:underline"
              to={actionTo}
            >
              {actionLabel}
            </Link>
          ) : null}
        </CardContent>
      </Card>
    </div>
  );
}

function RouteLoadingFallback() {
  return <RoutePlaceholder title="Loading" description="正在加载管理端页面，请稍候。" />;
}

function renderLazyPage(Page: LazyExoticComponent<ComponentType>) {
  return (
    <Suspense fallback={<RouteLoadingFallback />}>
      <Page />
    </Suspense>
  );
}

export const routes = [
  {
    path: "/",
    element: (
      <AppShell>
        <Outlet />
      </AppShell>
    ),
    children: [
      {
        index: true,
        element: renderLazyPage(AuthPage),
      },
      {
        path: "dashboard",
        element: <DashboardShell />,
        children: [
          {
            index: true,
            element: renderLazyPage(DashboardPage),
          },
          {
            path: "servers",
            element: renderLazyPage(ServersPage),
          },
          {
            path: "auth-control",
            element: renderLazyPage(AuthControlPage),
          },
          {
            path: "announcement",
            element: renderLazyPage(AnnouncementPage),
          },
          {
            path: "coupon",
            element: renderLazyPage(CouponPage),
          },
          {
            path: "document",
            element: renderLazyPage(DocumentPage),
          },
          {
            path: "marketing",
            element: renderLazyPage(MarketingPage),
          },
          {
            path: "nodes",
            element: renderLazyPage(NodesPage),
          },
          {
            path: "order",
            element: renderLazyPage(OrderPage),
          },
          {
            path: "payment",
            element: renderLazyPage(PaymentPage),
          },
          {
            path: "product",
            element: renderLazyPage(ProductPage),
          },
          {
            path: "subscribe",
            element: renderLazyPage(SubscribePage),
          },
          {
            path: "system",
            element: renderLazyPage(SystemPage),
          },
          {
            path: "ticket",
            element: renderLazyPage(TicketPage),
          },
          {
            path: "user",
            element: renderLazyPage(UserPage),
          },
          {
            path: "log/login",
            element: renderLazyPage(LoginLogPage),
          },
          {
            path: "log/register",
            element: renderLazyPage(RegisterLogPage),
          },
          {
            path: "log/email",
            element: renderLazyPage(EmailLogPage),
          },
          {
            path: "log/mobile",
            element: renderLazyPage(MobileLogPage),
          },
          {
            path: "log/subscribe",
            element: renderLazyPage(SubscribeLogPage),
          },
          {
            path: "log/reset-subscribe",
            element: renderLazyPage(ResetSubscribeLogPage),
          },
          {
            path: "log/subscribe-traffic",
            element: renderLazyPage(SubscribeTrafficLogPage),
          },
          {
            path: "log/server-traffic",
            element: renderLazyPage(ServerTrafficLogPage),
          },
          {
            path: "log/traffic-details",
            element: renderLazyPage(TrafficDetailsLogPage),
          },
          {
            path: "log/balance",
            element: renderLazyPage(BalanceLogPage),
          },
          {
            path: "log/commission",
            element: renderLazyPage(CommissionLogPage),
          },
          {
            path: "log/gift",
            element: renderLazyPage(GiftLogPage),
          },
          {
            path: "*",
            element: (
              <RoutePlaceholder
                title="Admin Route Placeholder"
                description="当前阶段只验证管理端壳与关键路由，完整页面树会在后续 unit 接入。"
              />
            ),
          },
        ],
      },
      {
        path: "*",
        element: <Navigate replace to="/" />,
      },
    ],
  },
];

import { Card, CardContent, CardHeader, CardTitle } from "@workspace/ui/components/card";
import { type ComponentType, type LazyExoticComponent, lazy, Suspense } from "react";
import { Link, Outlet } from "react-router-dom";
import DashboardLayout from "@/app/dashboard/layout";
import { ADMIN_HOME_PATH } from "@/utils/admin-path";
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

function RootShell() {
  return (
    <AppShell>
      <Outlet />
    </AppShell>
  );
}

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

function RouteNotFound({ actionLabel, actionTo }: { actionLabel: string; actionTo: string }) {
  return (
    <RoutePlaceholder
      title="404"
      description="页面不存在，你访问的管理端地址无效或已经下线。"
      actionLabel={actionLabel}
      actionTo={actionTo}
    />
  );
}

function createLazyRouteComponent(Page: LazyExoticComponent<ComponentType>) {
  return function LazyRouteComponent() {
    return (
      <Suspense fallback={<RouteLoadingFallback />}>
        <Page />
      </Suspense>
    );
  };
}

export const routes = [
  {
    path: "/",
    Component: RootShell,
    children: [
      {
        index: true,
        Component: createLazyRouteComponent(AuthPage),
      },
      {
        path: "dashboard",
        Component: DashboardShell,
        children: [
          {
            index: true,
            element: <RouteNotFound actionLabel="前往工作台" actionTo={ADMIN_HOME_PATH} />,
          },
          {
            path: "workplace",
            Component: createLazyRouteComponent(DashboardPage),
          },
          {
            path: "servers",
            Component: createLazyRouteComponent(ServersPage),
          },
          {
            path: "auth-control",
            Component: createLazyRouteComponent(AuthControlPage),
          },
          {
            path: "announcement",
            Component: createLazyRouteComponent(AnnouncementPage),
          },
          {
            path: "coupon",
            Component: createLazyRouteComponent(CouponPage),
          },
          {
            path: "document",
            Component: createLazyRouteComponent(DocumentPage),
          },
          {
            path: "marketing",
            Component: createLazyRouteComponent(MarketingPage),
          },
          {
            path: "nodes",
            Component: createLazyRouteComponent(NodesPage),
          },
          {
            path: "order",
            Component: createLazyRouteComponent(OrderPage),
          },
          {
            path: "payment",
            Component: createLazyRouteComponent(PaymentPage),
          },
          {
            path: "product",
            Component: createLazyRouteComponent(ProductPage),
          },
          {
            path: "subscribe",
            Component: createLazyRouteComponent(SubscribePage),
          },
          {
            path: "system",
            Component: createLazyRouteComponent(SystemPage),
          },
          {
            path: "ticket",
            Component: createLazyRouteComponent(TicketPage),
          },
          {
            path: "user",
            Component: createLazyRouteComponent(UserPage),
          },
          {
            path: "log/login",
            Component: createLazyRouteComponent(LoginLogPage),
          },
          {
            path: "log/register",
            Component: createLazyRouteComponent(RegisterLogPage),
          },
          {
            path: "log/email",
            Component: createLazyRouteComponent(EmailLogPage),
          },
          {
            path: "log/mobile",
            Component: createLazyRouteComponent(MobileLogPage),
          },
          {
            path: "log/subscribe",
            Component: createLazyRouteComponent(SubscribeLogPage),
          },
          {
            path: "log/reset-subscribe",
            Component: createLazyRouteComponent(ResetSubscribeLogPage),
          },
          {
            path: "log/subscribe-traffic",
            Component: createLazyRouteComponent(SubscribeTrafficLogPage),
          },
          {
            path: "log/server-traffic",
            Component: createLazyRouteComponent(ServerTrafficLogPage),
          },
          {
            path: "log/traffic-details",
            Component: createLazyRouteComponent(TrafficDetailsLogPage),
          },
          {
            path: "log/balance",
            Component: createLazyRouteComponent(BalanceLogPage),
          },
          {
            path: "log/commission",
            Component: createLazyRouteComponent(CommissionLogPage),
          },
          {
            path: "log/gift",
            Component: createLazyRouteComponent(GiftLogPage),
          },
          {
            path: "*",
            element: <RouteNotFound actionLabel="返回仪表盘" actionTo={ADMIN_HOME_PATH} />,
          },
        ],
      },
      {
        path: "*",
        element: <RouteNotFound actionLabel="返回登录页" actionTo="/" />,
      },
    ],
  },
];

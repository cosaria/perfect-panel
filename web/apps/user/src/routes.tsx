import { Card, CardContent, CardHeader, CardTitle } from "@workspace/ui/components/card";
import { type ComponentType, type LazyExoticComponent, lazy, Suspense } from "react";
import { Link, Outlet } from "react-router-dom";
import DashboardLayout from "@/src/pages/(main)/(user)/layout";
import MainLayout from "@/src/pages/(main)/layout";
import AppShell from "./app-shell";

const AffiliatePage = lazy(() => import("@/src/pages/(main)/(user)/affiliate/page"));
const AnnouncementPage = lazy(() => import("@/src/pages/(main)/(user)/announcement/page"));
const DocumentPage = lazy(() => import("@/src/pages/(main)/(user)/document/page"));
const HomePage = lazy(() => import("@/src/pages/(main)/page"));
const OrderPage = lazy(() => import("@/src/pages/(main)/(user)/order/page"));
const PaymentPage = lazy(() => import("@/src/pages/(main)/(user)/payment/page"));
const ProfilePage = lazy(() => import("@/src/pages/(main)/(user)/profile/page"));
const SubscribePage = lazy(() => import("@/src/pages/(main)/(user)/subscribe/page"));
const TicketPage = lazy(() => import("@/src/pages/(main)/(user)/ticket/page"));
const WalletPage = lazy(() => import("@/src/pages/(main)/(user)/wallet/page"));
const PrivacyPolicyPage = lazy(() => import("@/src/pages/(main)/privacy-policy/page"));
const PurchasingOrderPage = lazy(() => import("@/src/pages/(main)/purchasing/order/page"));
const PurchasingPage = lazy(() => import("@/src/pages/(main)/purchasing/page"));
const TosPage = lazy(() => import("@/src/pages/(main)/tos/page"));
const AuthPage = lazy(() => import("@/src/pages/auth/page"));
const BindPage = lazy(() => import("@/src/pages/bind/[platform]/page"));
const DashboardPage = lazy(() => import("@/src/pages/(main)/(user)/dashboard/page"));
const OAuthPage = lazy(() => import("@/src/pages/oauth/[platform]/page"));

function MainShell() {
  return (
    <MainLayout>
      <Outlet />
    </MainLayout>
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
  return <RoutePlaceholder title="Loading" description="正在加载用户端页面，请稍候。" />;
}

function RouteNotFound({ actionLabel, actionTo }: { actionLabel: string; actionTo: string }) {
  return (
    <RoutePlaceholder
      title="404"
      description="页面不存在，你访问的地址无效或已经下线。"
      actionLabel={actionLabel}
      actionTo={actionTo}
    />
  );
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
        path: "auth",
        element: renderLazyPage(AuthPage),
      },
      {
        path: "bind/:platform",
        element: renderLazyPage(BindPage),
      },
      {
        path: "oauth/:platform",
        element: renderLazyPage(OAuthPage),
      },
      {
        element: <MainShell />,
        children: [
          {
            index: true,
            element: renderLazyPage(HomePage),
          },
          {
            path: "privacy-policy",
            element: renderLazyPage(PrivacyPolicyPage),
          },
          {
            path: "tos",
            element: renderLazyPage(TosPage),
          },
          {
            path: "purchasing",
            element: renderLazyPage(PurchasingPage),
          },
          {
            path: "purchasing/order",
            element: renderLazyPage(PurchasingOrderPage),
          },
          {
            element: <DashboardShell />,
            children: [
              {
                path: "dashboard",
                element: renderLazyPage(DashboardPage),
              },
              {
                path: "profile",
                element: renderLazyPage(ProfilePage),
              },
              {
                path: "subscribe",
                element: renderLazyPage(SubscribePage),
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
                path: "wallet",
                element: renderLazyPage(WalletPage),
              },
              {
                path: "affiliate",
                element: renderLazyPage(AffiliatePage),
              },
              {
                path: "document",
                element: renderLazyPage(DocumentPage),
              },
              {
                path: "announcement",
                element: renderLazyPage(AnnouncementPage),
              },
              {
                path: "ticket",
                element: renderLazyPage(TicketPage),
              },
              {
                path: "*",
                element: <RouteNotFound actionLabel="返回仪表盘" actionTo="/dashboard" />,
              },
            ],
          },
          {
            path: "*",
            element: <RouteNotFound actionLabel="返回首页" actionTo="/" />,
          },
        ],
      },
      {
        path: "*",
        element: <RouteNotFound actionLabel="返回首页" actionTo="/" />,
      },
    ],
  },
];

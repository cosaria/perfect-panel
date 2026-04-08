import { Card, CardContent, CardHeader, CardTitle } from "@workspace/ui/components/card";
import { type ComponentType, type LazyExoticComponent, lazy, Suspense } from "react";
import { Navigate, Outlet } from "react-router-dom";
import DashboardLayout from "@/app/(main)/(user)/layout";
import MainLayout from "@/app/(main)/layout";
import AppShell from "./app-shell";

const AffiliatePage = lazy(() => import("@/app/(main)/(user)/affiliate/page"));
const AnnouncementPage = lazy(() => import("@/app/(main)/(user)/announcement/page"));
const DocumentPage = lazy(() => import("@/app/(main)/(user)/document/page"));
const HomePage = lazy(() => import("@/app/(main)/page"));
const OrderPage = lazy(() => import("@/app/(main)/(user)/order/page"));
const PaymentPage = lazy(() => import("@/app/(main)/(user)/payment/page"));
const ProfilePage = lazy(() => import("@/app/(main)/(user)/profile/page"));
const SubscribePage = lazy(() => import("@/app/(main)/(user)/subscribe/page"));
const TicketPage = lazy(() => import("@/app/(main)/(user)/ticket/page"));
const WalletPage = lazy(() => import("@/app/(main)/(user)/wallet/page"));
const PrivacyPolicyPage = lazy(() => import("@/app/(main)/privacy-policy/page"));
const PurchasingOrderPage = lazy(() => import("@/app/(main)/purchasing/order/page"));
const PurchasingPage = lazy(() => import("@/app/(main)/purchasing/page"));
const TosPage = lazy(() => import("@/app/(main)/tos/page"));
const AuthPage = lazy(() => import("@/app/auth/page"));
const BindPage = lazy(() => import("@/app/bind/[platform]/page"));
const DashboardPage = lazy(() => import("@/app/(main)/(user)/dashboard/page"));
const OAuthPage = lazy(() => import("@/app/oauth/[platform]/page"));

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

function RoutePlaceholder({ title, description }: { title: string; description: string }) {
  return (
    <div className="mx-auto flex h-full max-w-3xl items-center justify-center">
      <Card className="w-full">
        <CardHeader>
          <CardTitle>{title}</CardTitle>
        </CardHeader>
        <CardContent className="text-muted-foreground text-sm">{description}</CardContent>
      </Card>
    </div>
  );
}

function RouteLoadingFallback() {
  return <RoutePlaceholder title="Loading" description="正在加载用户端页面，请稍候。" />;
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
            ],
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

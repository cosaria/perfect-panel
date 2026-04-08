import { Card, CardContent, CardHeader, CardTitle } from "@workspace/ui/components/card";
import { type ComponentType, type LazyExoticComponent, lazy, Suspense } from "react";
import { Navigate, Outlet } from "react-router-dom";
import DashboardLayout from "@/app/(main)/(user)/layout";
import MainLayout from "@/app/(main)/layout";
import AppShell from "./app-shell";

const HomePage = lazy(() => import("@/app/(main)/page"));
const AuthPage = lazy(() => import("@/app/auth/page"));
const DashboardPage = lazy(() => import("@/app/(main)/(user)/dashboard/page"));

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
        element: <MainShell />,
        children: [
          {
            index: true,
            element: renderLazyPage(HomePage),
          },
          {
            path: "dashboard",
            element: <DashboardShell />,
            children: [
              {
                index: true,
                element: renderLazyPage(DashboardPage),
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

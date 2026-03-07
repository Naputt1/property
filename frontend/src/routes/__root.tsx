import {
  createRootRouteWithContext,
  Link,
  Outlet,
} from "@tanstack/react-router";
import { TanStackRouterDevtools } from "@tanstack/react-router-devtools";
import type { QueryClient } from "@tanstack/react-query";
import { Toaster } from "sonner";

export interface MyRouterContext {
  queryClient: QueryClient;
}

export const Route = createRootRouteWithContext<MyRouterContext>()({
  component: () => (
    <div className="min-h-screen flex flex-col bg-muted/30 text-secondary-foreground">
      <header className="border-b border-border bg-white sticky top-0 z-10">
        <div className="container mx-auto px-4 h-16 flex items-center justify-between">
          <Link to="/properties" className="text-xl font-bold text-primary">
            Property Market AI
          </Link>
          <nav className="flex items-center gap-6">
            <Link
              to="/properties"
              className="text-sm font-medium hover:text-primary transition-colors [&.active]:text-primary"
            >
              Properties
            </Link>
            <Link
              to="/analytics"
              className="text-sm font-medium hover:text-primary transition-colors [&.active]:text-primary"
            >
              Analytics
            </Link>
            <Link
              to="/admin"
              className="text-sm font-medium hover:text-primary transition-colors [&.active]:text-primary"
            >
              Admin
            </Link>
          </nav>
        </div>
      </header>
      <main className="flex-1 container mx-auto px-4">
        <Outlet />
      </main>
      <Toaster position="top-right" richColors />
      <TanStackRouterDevtools />
    </div>
  ),
});

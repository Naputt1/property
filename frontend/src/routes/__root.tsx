import { createRootRouteWithContext, Link, Outlet } from "@tanstack/react-router";
import { TanStackRouterDevtools } from "@tanstack/react-router-devtools";
import type { QueryClient } from "@tanstack/react-query";
import { Toaster } from "sonner";

export interface MyRouterContext {
  queryClient: QueryClient;
}

export const Route = createRootRouteWithContext<MyRouterContext>()({
  component: () => (
    <div className="min-h-screen flex flex-col bg-gray-50 text-gray-900">
      <header className="border-b bg-white sticky top-0 z-10">
        <div className="container mx-auto px-4 h-16 flex items-center justify-between">
          <Link to="/" className="text-xl font-bold text-blue-600">
            Property Market AI
          </Link>
          <nav className="flex items-center gap-6">
            <Link to="/" className="text-sm font-medium hover:text-blue-600 [&.active]:text-blue-600">
              Home
            </Link>
            <Link to="/properties" className="text-sm font-medium hover:text-blue-600 [&.active]:text-blue-600">
              Properties
            </Link>
            <Link to="/analytics" className="text-sm font-medium hover:text-blue-600 [&.active]:text-blue-600">
              Analytics
            </Link>
            <Link to="/admin" className="text-sm font-medium hover:text-blue-600 [&.active]:text-blue-600">
              Admin
            </Link>
          </nav>
        </div>
      </header>
      <main className="flex-1 container mx-auto px-4 py-8">
        <Outlet />
      </main>
      <Toaster position="top-right" richColors />
      <TanStackRouterDevtools />
    </div>
  ),
});

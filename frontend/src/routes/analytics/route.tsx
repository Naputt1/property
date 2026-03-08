import { createFileRoute, Outlet } from "@tanstack/react-router";

export const Route = createFileRoute("/analytics")({
  component: RouteComponent,
});

function RouteComponent() {
  return (
    <div className="py-8">
      <Outlet />
    </div>
  );
}

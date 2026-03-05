import { createFileRoute, Link } from "@tanstack/react-router";

export const Route = createFileRoute("/")({
  component: Index,
});

function Index() {
  return (
    <div className="p-8">
      <h1 className="text-4xl font-bold mb-4">Welcome to Property Market AI</h1>
      <Link to="/properties" className="text-blue-600 hover:underline">
        View UK Properties
      </Link>
    </div>
  );
}

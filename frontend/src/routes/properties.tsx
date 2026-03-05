import { createFileRoute } from "@tanstack/react-router";
import { useQuery } from "@tanstack/react-query";
import { useState } from "react";
import { propertyListQuery } from "@/query/property";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Button } from "@/components/ui/button";

export const Route = createFileRoute("/properties")({
  component: PropertiesPage,
});

function PropertiesPage() {
  const [page, setPage] = useState(1);
  const pageSize = 10;

  const { data, isLoading, error } = useQuery(
    propertyListQuery.getOptions({
      param: {
        page,
        pageSize,
      },
    }),
  );

  if (isLoading) return <div className="p-8">Loading properties...</div>;
  if (error)
    return <div className="p-8 text-destructive">Error: {error.message}</div>;

  const properties = data?.data || [];

  return (
    <div className="container mx-auto py-10 px-4">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold">UK Property Market</h1>
      </div>

      <div className="border rounded-md">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Address</TableHead>
              <TableHead>Town/City</TableHead>
              <TableHead>Postcode</TableHead>
              <TableHead>Price</TableHead>
              <TableHead>Date</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {properties.length > 0 ? (
              properties.map((property: any) => (
                <TableRow key={property.id}>
                  <TableCell>
                    {property.paon} {property.saon} {property.street}
                  </TableCell>
                  <TableCell>{property.town_city}</TableCell>
                  <TableCell className="font-mono">
                    {property.postcode}
                  </TableCell>
                  <TableCell>£{property.price.toLocaleString()}</TableCell>
                  <TableCell>
                    {new Date(property.date_of_transfer).toLocaleDateString()}
                  </TableCell>
                </TableRow>
              ))
            ) : (
              <TableRow>
                <TableCell colSpan={5} className="text-center h-24">
                  No properties found.
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>

      <div className="flex items-center justify-end space-x-2 py-4">
        <Button
          variant="outline"
          size="sm"
          onClick={() => setPage((p) => Math.max(1, p - 1))}
          disabled={page === 1}
        >
          Previous
        </Button>
        <div className="text-sm font-medium">Page {page}</div>
        <Button
          variant="outline"
          size="sm"
          onClick={() => setPage((p) => p + 1)}
          disabled={properties.length < pageSize}
        >
          Next
        </Button>
      </div>
    </div>
  );
}

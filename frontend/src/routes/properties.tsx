import { createFileRoute } from "@tanstack/react-router";
import { useState } from "react";
import { useGetPropertiesQuery } from "@/gen-gql/graphql";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Button } from "@/components/ui/button";

import {
  Pagination,
  PaginationContent,
  PaginationEllipsis,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination";

export const Route = createFileRoute("/properties")({
  component: PropertiesPage,
});

function PropertiesPage() {
  const [page, setPage] = useState(1);
  const [town, setTown] = useState("");
  const [county, setCounty] = useState("");
  const [propertyType, setPropertyType] = useState("");
  const [minPrice, setMinPrice] = useState<string>("");
  const [maxPrice, setMaxPrice] = useState<string>("");

  const pageSize = 10;

  const { data, isLoading } = useGetPropertiesQuery({
    limit: pageSize,
    offset: (page - 1) * pageSize,
    townCity: town || undefined,
    county: county || undefined,
    propertyType: propertyType || undefined,
    minPrice: minPrice ? parseInt(minPrice) : undefined,
    maxPrice: maxPrice ? parseInt(maxPrice) : undefined,
  });

  const formatPrice = (value: number) => {
    return new Intl.NumberFormat("en-GB", {
      style: "currency",
      currency: "GBP",
      maximumFractionDigits: 0,
    }).format(value);
  };

  const getPropertyTypeLabel = (type: string) => {
    const types: Record<string, string> = {
      D: "Detached",
      S: "Semi-Detached",
      T: "Terraced",
      F: "Flat/Maisonette",
      O: "Other",
    };
    return types[type] || type;
  };

  const handleFilterChange = () => {
    setPage(1); // Reset to first page on filter change
  };

  const properties = data?.properties.items || [];
  const totalItems = data?.properties.total || 0;
  const totalPages = Math.ceil(totalItems / pageSize);

  const renderPaginationItems = () => {
    const items = [];
    const maxVisiblePages = 5;

    if (totalPages <= maxVisiblePages) {
      for (let i = 1; i <= totalPages; i++) {
        items.push(
          <PaginationItem key={i}>
            <PaginationLink
              onClick={() => setPage(i)}
              isActive={page === i}
              className="cursor-pointer"
            >
              {i}
            </PaginationLink>
          </PaginationItem>,
        );
      }
    } else {
      // Always show first page
      items.push(
        <PaginationItem key={1}>
          <PaginationLink
            onClick={() => setPage(1)}
            isActive={page === 1}
            className="cursor-pointer"
          >
            1
          </PaginationLink>
        </PaginationItem>,
      );

      if (page > 3) {
        items.push(
          <PaginationItem key="ellipsis-start">
            <PaginationEllipsis />
          </PaginationItem>,
        );
      }

      // Show pages around current page
      const start = Math.max(2, page - 1);
      const end = Math.min(totalPages - 1, page + 1);

      for (let i = start; i <= end; i++) {
        if (i === 1 || i === totalPages) continue;
        items.push(
          <PaginationItem key={i}>
            <PaginationLink
              onClick={() => setPage(i)}
              isActive={page === i}
              className="cursor-pointer"
            >
              {i}
            </PaginationLink>
          </PaginationItem>,
        );
      }

      if (page < totalPages - 2) {
        items.push(
          <PaginationItem key="ellipsis-end">
            <PaginationEllipsis />
          </PaginationItem>,
        );
      }

      // Always show last page
      items.push(
        <PaginationItem key={totalPages}>
          <PaginationLink
            onClick={() => setPage(totalPages)}
            isActive={page === totalPages}
            className="cursor-pointer"
          >
            {totalPages}
          </PaginationLink>
        </PaginationItem>,
      );
    }

    return items;
  };

  return (
    <div className="container mx-auto py-10 px-4">
      <div className="flex flex-col md:flex-row justify-between items-start md:items-center mb-8 gap-4">
        <div>
          <h1 className="text-3xl font-bold">UK Property Market</h1>
          <p className="text-muted-foreground mt-1">
            Browse through historical housing transaction data.
          </p>
        </div>
        <div className="text-right">
          <div className="text-2xl font-bold text-primary">
            {totalItems.toLocaleString()}
          </div>
          <div className="text-xs text-muted-foreground uppercase tracking-wider">
            Total Transactions
          </div>
        </div>
      </div>

      {/* Filters */}
      <div className="bg-white p-4 rounded-lg border border-border shadow-sm mb-8 grid grid-cols-1 md:grid-cols-3 lg:grid-cols-6 gap-4">
        <div className="space-y-1.5">
          <label className="text-xs font-semibold uppercase text-muted-foreground">
            Town/City
          </label>
          <input
            type="text"
            placeholder="Search town..."
            className="w-full border border-input rounded px-3 py-1.5 text-sm"
            value={town}
            onChange={(e) => {
              setTown(e.target.value);
              handleFilterChange();
            }}
          />
        </div>
        <div className="space-y-1.5">
          <label className="text-xs font-semibold uppercase text-muted-foreground">
            County
          </label>
          <input
            type="text"
            placeholder="Search county..."
            className="w-full border border-input rounded px-3 py-1.5 text-sm"
            value={county}
            onChange={(e) => {
              setCounty(e.target.value);
              handleFilterChange();
            }}
          />
        </div>
        <div className="space-y-1.5">
          <label className="text-xs font-semibold uppercase text-muted-foreground">
            Type
          </label>
          <select
            className="w-full border border-input rounded px-3 py-1.5 text-sm bg-muted/50"
            value={propertyType}
            onChange={(e) => {
              setPropertyType(e.target.value);
              handleFilterChange();
            }}
          >
            <option value="">All Types</option>
            <option value="D">Detached</option>
            <option value="S">Semi-Detached</option>
            <option value="T">Terraced</option>
            <option value="F">Flat</option>
            <option value="O">Other</option>
          </select>
        </div>
        <div className="space-y-1.5">
          <label className="text-xs font-semibold uppercase text-muted-foreground">
            Min Price
          </label>
          <input
            type="number"
            placeholder="Min £"
            className="w-full border border-input rounded px-3 py-1.5 text-sm"
            value={minPrice}
            onChange={(e) => {
              setMinPrice(e.target.value);
              handleFilterChange();
            }}
          />
        </div>
        <div className="space-y-1.5">
          <label className="text-xs font-semibold uppercase text-muted-foreground">
            Max Price
          </label>
          <input
            type="number"
            placeholder="Max £"
            className="w-full border border-input rounded px-3 py-1.5 text-sm"
            value={maxPrice}
            onChange={(e) => {
              setMaxPrice(e.target.value);
              handleFilterChange();
            }}
          />
        </div>
        <div className="flex items-end">
          <Button
            variant="outline"
            className="w-full border-input"
            onClick={() => {
              setTown("");
              setCounty("");
              setPropertyType("");
              setMinPrice("");
              setMaxPrice("");
              setPage(1);
            }}
          >
            Reset
          </Button>
        </div>
      </div>

      <div className="bg-white rounded-xl border border-border shadow-sm overflow-hidden">
        <Table>
          <TableHeader className="bg-muted/50">
            <TableRow>
              <TableHead className="w-[40%]">Address</TableHead>
              <TableHead>Location</TableHead>
              <TableHead>Type</TableHead>
              <TableHead className="text-right">Price</TableHead>
              <TableHead className="text-right">Date</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {isLoading ? (
              Array.from({ length: pageSize }).map((_, i) => (
                <TableRow key={i} className="animate-pulse">
                  <TableCell>
                    <div className="h-4 bg-muted rounded w-3/4"></div>
                  </TableCell>
                  <TableCell>
                    <div className="h-4 bg-muted rounded w-1/2"></div>
                  </TableCell>
                  <TableCell>
                    <div className="h-4 bg-muted rounded w-1/4"></div>
                  </TableCell>
                  <TableCell>
                    <div className="h-4 bg-muted rounded w-1/3 ml-auto"></div>
                  </TableCell>
                  <TableCell>
                    <div className="h-4 bg-muted rounded w-1/4 ml-auto"></div>
                  </TableCell>
                </TableRow>
              ))
            ) : properties.length > 0 ? (
              properties.map((property: any) => (
                <TableRow
                  key={property.id}
                  className="group hover:bg-muted/50 transition-colors"
                >
                  <TableCell>
                    <div className="font-medium text-gray-900 line-clamp-1">
                      {property.address}
                    </div>
                    <div className="text-xs text-muted-foreground font-mono">
                      {property.postcodeOutward} {property.postcodeInward}
                    </div>
                  </TableCell>
                  <TableCell>
                    <div className="text-sm">{property.townCity}</div>
                    <div className="text-xs text-muted-foreground">
                      {property.county}
                    </div>
                  </TableCell>
                  <TableCell>
                    <span
                      className={`inline-flex items-center px-2 py-0.5 rounded text-xs font-medium ${
                        property.propertyType === "D"
                          ? "bg-primary/10 text-primary"
                          : property.propertyType === "F"
                            ? "bg-accent text-accent-foreground"
                            : "bg-secondary text-secondary-foreground"
                      }`}
                    >
                      {getPropertyTypeLabel(property.propertyType)}
                    </span>
                    {property.oldNew === "Y" && (
                      <span className="ml-1 inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-chart-2/10 text-chart-2">
                        New Build
                      </span>
                    )}
                  </TableCell>
                  <TableCell className="text-right font-semibold text-gray-900">
                    {formatPrice(property.price)}
                  </TableCell>
                  <TableCell className="text-right text-sm text-muted-foreground">
                    {new Date(property.dateOfTransfer).toLocaleDateString(
                      "en-GB",
                      {
                        day: "2-digit",
                        month: "short",
                        year: "numeric",
                      },
                    )}
                  </TableCell>
                </TableRow>
              ))
            ) : (
              <TableRow>
                <TableCell colSpan={5} className="text-center py-20">
                  <div className="flex flex-col items-center">
                    <div className="text-muted-foreground mb-2">
                      No properties found
                    </div>
                    <Button
                      variant="link"
                      onClick={() => {
                        setTown("");
                        setCounty("");
                        setPropertyType("");
                        setMinPrice("");
                        setMaxPrice("");
                      }}
                    >
                      Clear all filters
                    </Button>
                  </div>
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>

      <div className="flex items-center justify-between py-6">
        <div className="text-sm text-muted-foreground">
          Showing page <span className="font-medium">{page}</span> of{" "}
          <span className="font-medium">{totalPages || 1}</span>
        </div>
        <Pagination className="justify-end w-auto mx-0">
          <PaginationContent>
            <PaginationItem>
              <PaginationPrevious
                onClick={() => setPage((p) => Math.max(1, p - 1))}
                className={
                  page === 1 || isLoading
                    ? "pointer-events-none opacity-50"
                    : "cursor-pointer"
                }
              />
            </PaginationItem>

            {renderPaginationItems()}

            <PaginationItem>
              <PaginationNext
                onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
                className={
                  page === totalPages || totalPages === 0 || isLoading
                    ? "pointer-events-none opacity-50"
                    : "cursor-pointer"
                }
              />
            </PaginationItem>
          </PaginationContent>
        </Pagination>
      </div>
    </div>
  );
}

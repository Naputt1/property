import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useState, useEffect } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useAuthStore } from "@/store/auth";
import {
  propertyListQuery,
  createPropertyMutation,
  updatePropertyMutation,
  deletePropertyMutation,
} from "@/query/property";
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
  PaginationItem,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination";
import { toast } from "sonner";
import type { IProperty } from "@/types/property";

export const Route = createFileRoute("/admin/properties")({
  component: AdminProperties,
});

function AdminProperties() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { user } = useAuthStore();
  const [page, setPage] = useState(1);
  const [editingProperty, setEditingProperty] =
    useState<Partial<IProperty> | null>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);

  const pageSize = 10;

  useEffect(() => {
    if (!user || !user.is_admin) {
      navigate({ to: "/login" });
    }
  }, [user, navigate]);

  const { data, isLoading } = useQuery(
    propertyListQuery.getOptions({
      param: { page, pageSize },
    }),
  );

  const createMutation = useMutation(
    createPropertyMutation.getOptions({
      onSuccess: () => {
        toast.success("Property created successfully");
        setIsModalOpen(false);
        queryClient.invalidateQueries({ queryKey: ["property", "list"] });
      },
      onError: (error: any) => {
        toast.error(error.message || "Failed to create property");
      },
    }),
  );

  const updateMutation = useMutation(
    updatePropertyMutation.getOptions({
      onSuccess: () => {
        toast.success("Property updated successfully");
        setIsModalOpen(false);
        setEditingProperty(null);
        queryClient.invalidateQueries({ queryKey: ["property", "list"] });
      },
      onError: (error: any) => {
        toast.error(error.message || "Failed to update property");
      },
    }),
  );

  const deleteMutation = useMutation(
    deletePropertyMutation.getOptions({
      onSuccess: () => {
        toast.success("Property deleted successfully");
        queryClient.invalidateQueries({ queryKey: ["property", "list"] });
      },
      onError: (error: any) => {
        toast.error(error.message || "Failed to delete property");
      },
    }),
  );

  const properties = data?.data || [];
  const totalItems = data?.total || 0;
  const totalPages = Math.ceil(totalItems / pageSize);

  const handleEdit = (property: IProperty) => {
    setEditingProperty(property);
    setIsModalOpen(true);
  };

  const handleDelete = (id: string) => {
    if (window.confirm("Are you sure you want to delete this property?")) {
      deleteMutation.mutate(id);
    }
  };

  const handleAddNew = () => {
    setEditingProperty({
      date_of_transfer: new Date().toISOString(),
      property_type: "D",
      old_new: "N",
      duration: "F",
    });
    setIsModalOpen(true);
  };

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const formData = new FormData(e.currentTarget);
    const data: any = Object.fromEntries(formData.entries());

    // Convert price to number
    data.price = parseInt(data.price);

    if (editingProperty?.id) {
      updateMutation.mutate({ id: editingProperty.id, data });
    } else {
      createMutation.mutate(data);
    }
  };

  if (!user) return null;

  return (
    <div className="container mx-auto py-10 px-4">
      <div className="flex justify-between items-center mb-8">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">
            Manage Properties
          </h1>
          <p className="text-gray-500">
            Admin control panel for property records.
          </p>
        </div>
        <Button onClick={handleAddNew}>Add New Property</Button>
      </div>

      <div className="bg-white rounded-xl border shadow-sm overflow-hidden mb-6">
        <Table>
          <TableHeader className="bg-gray-50">
            <TableRow>
              <TableHead>Address</TableHead>
              <TableHead>Town/City</TableHead>
              <TableHead className="text-right">Price</TableHead>
              <TableHead className="text-right">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {isLoading ? (
              <TableRow>
                <TableCell colSpan={4} className="text-center py-10">
                  Loading...
                </TableCell>
              </TableRow>
            ) : properties.length > 0 ? (
              properties.map((p: IProperty) => (
                <TableRow key={p.id}>
                  <TableCell>
                    <div className="font-medium">
                      {p.paon} {p.saon} {p.street}
                    </div>
                    <div className="text-xs text-gray-400 font-mono">
                      {p.postcode}
                    </div>
                  </TableCell>
                  <TableCell>{p.town_city}</TableCell>
                  <TableCell className="text-right font-semibold">
                    £{p.price.toLocaleString()}
                  </TableCell>
                  <TableCell className="text-right">
                    <div className="flex justify-end gap-2">
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={() => handleEdit(p)}
                      >
                        Edit
                      </Button>
                      <Button
                        variant="destructive"
                        size="sm"
                        onClick={() => handleDelete(p.id)}
                      >
                        Delete
                      </Button>
                    </div>
                  </TableCell>
                </TableRow>
              ))
            ) : (
              <TableRow>
                <TableCell
                  colSpan={4}
                  className="text-center py-10 text-gray-400"
                >
                  No properties found.
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>

      <div className="flex items-center justify-between">
        <div className="text-sm text-gray-500">
          Total:{" "}
          <span className="font-medium">{totalItems.toLocaleString()}</span>{" "}
          properties
        </div>
        <Pagination className="justify-end w-auto mx-0">
          <PaginationContent>
            <PaginationItem>
              <PaginationPrevious
                onClick={() => setPage((p) => Math.max(1, p - 1))}
                className={
                  page === 1
                    ? "pointer-events-none opacity-50"
                    : "cursor-pointer"
                }
              />
            </PaginationItem>
            <PaginationItem>
              <span className="px-4 text-sm font-medium text-gray-700">
                Page {page} of {totalPages || 1}
              </span>
            </PaginationItem>
            <PaginationItem>
              <PaginationNext
                onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
                className={
                  page === totalPages || totalPages === 0
                    ? "pointer-events-none opacity-50"
                    : "cursor-pointer"
                }
              />
            </PaginationItem>
          </PaginationContent>
        </Pagination>
      </div>

      {/* Manual Modal (Standard HTML/CSS since we don't have Dialog UI component) */}
      {isModalOpen && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-100 p-4">
          <div className="bg-white rounded-xl shadow-xl w-full max-w-2xl max-h-[90vh] overflow-y-auto">
            <div className="p-6 border-b sticky top-0 bg-white z-10 flex justify-between items-center">
              <h2 className="text-xl font-bold">
                {editingProperty?.id ? "Edit Property" : "Add New Property"}
              </h2>
              <button
                onClick={() => setIsModalOpen(false)}
                className="text-gray-400 hover:text-gray-600"
              >
                &times;
              </button>
            </div>
            <form onSubmit={handleSubmit} className="p-6 space-y-4">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="space-y-1">
                  <label className="text-xs font-bold uppercase text-gray-500">
                    Street
                  </label>
                  <input
                    name="street"
                    defaultValue={editingProperty?.street}
                    className="w-full border rounded px-3 py-2"
                    required
                  />
                </div>
                <div className="space-y-1">
                  <label className="text-xs font-bold uppercase text-gray-500">
                    Postcode
                  </label>
                  <input
                    name="postcode"
                    defaultValue={editingProperty?.postcode}
                    className="w-full border rounded px-3 py-2"
                    required
                  />
                </div>
                <div className="space-y-1">
                  <label className="text-xs font-bold uppercase text-gray-500">
                    PAON (House No)
                  </label>
                  <input
                    name="paon"
                    defaultValue={editingProperty?.paon}
                    className="w-full border rounded px-3 py-2"
                    required
                  />
                </div>
                <div className="space-y-1">
                  <label className="text-xs font-bold uppercase text-gray-500">
                    SAON (Flat/Unit)
                  </label>
                  <input
                    name="saon"
                    defaultValue={editingProperty?.saon}
                    className="w-full border rounded px-3 py-2"
                  />
                </div>
                <div className="space-y-1">
                  <label className="text-xs font-bold uppercase text-gray-500">
                    Town/City
                  </label>
                  <input
                    name="town_city"
                    defaultValue={editingProperty?.town_city}
                    className="w-full border rounded px-3 py-2"
                    required
                  />
                </div>
                <div className="space-y-1">
                  <label className="text-xs font-bold uppercase text-gray-500">
                    District
                  </label>
                  <input
                    name="district"
                    defaultValue={editingProperty?.district}
                    className="w-full border rounded px-3 py-2"
                    required
                  />
                </div>
                <div className="space-y-1">
                  <label className="text-xs font-bold uppercase text-gray-500">
                    County
                  </label>
                  <input
                    name="county"
                    defaultValue={editingProperty?.county}
                    className="w-full border rounded px-3 py-2"
                    required
                  />
                </div>
                <div className="space-y-1">
                  <label className="text-xs font-bold uppercase text-gray-500">
                    Price (£)
                  </label>
                  <input
                    name="price"
                    type="number"
                    defaultValue={editingProperty?.price}
                    className="w-full border rounded px-3 py-2"
                    required
                  />
                </div>
                <div className="space-y-1">
                  <label className="text-xs font-bold uppercase text-gray-500">
                    Property Type
                  </label>
                  <select
                    name="property_type"
                    defaultValue={editingProperty?.property_type}
                    className="w-full border rounded px-3 py-2"
                  >
                    <option value="D">Detached</option>
                    <option value="S">Semi-Detached</option>
                    <option value="T">Terraced</option>
                    <option value="F">Flat/Maisonette</option>
                    <option value="O">Other</option>
                  </select>
                </div>
                <div className="space-y-1">
                  <label className="text-xs font-bold uppercase text-gray-500">
                    New Build?
                  </label>
                  <select
                    name="old_new"
                    defaultValue={editingProperty?.old_new}
                    className="w-full border rounded px-3 py-2"
                  >
                    <option value="Y">Yes</option>
                    <option value="N">No</option>
                  </select>
                </div>
              </div>
              <div className="pt-4 flex justify-end gap-3 border-t">
                <Button
                  type="button"
                  variant="outline"
                  onClick={() => setIsModalOpen(false)}
                >
                  Cancel
                </Button>
                <Button
                  type="submit"
                  disabled={
                    createMutation.isPending || updateMutation.isPending
                  }
                >
                  {editingProperty?.id ? "Update Property" : "Create Property"}
                </Button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}

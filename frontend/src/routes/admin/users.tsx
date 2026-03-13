import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useState, useEffect } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useAuthStore } from "@/store/auth";
import { initAxios } from "@/services/axios";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import { toast } from "sonner";
import type { User } from "@/types/user";

export const Route = createFileRoute("/admin/users")({
  component: AdminUsers,
});

function AdminUsers() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { user: currentUser } = useAuthStore();
  const [editingUser, setEditingUser] = useState<Partial<User> | null>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [password, setPassword] = useState("");

  const api = initAxios();

  useEffect(() => {
    if (!currentUser || !currentUser.is_admin) {
      navigate({ to: "/login" });
    }
  }, [currentUser, navigate]);

  const { data: users = [], isLoading } = useQuery({
    queryKey: ["admin-users"],
    queryFn: async () => {
      const res = await api.get<User[]>("/admin/users");
      return res.data;
    },
    enabled: !!currentUser?.is_admin,
  });

  const createMutation = useMutation({
    mutationFn: async (newUser: any) => {
      const res = await api.post("/admin/users", newUser);
      return res.data;
    },
    onSuccess: () => {
      toast.success("User created successfully");
      setIsModalOpen(false);
      queryClient.invalidateQueries({ queryKey: ["admin-users"] });
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error || "Failed to create user");
    },
  });

  const updateMutation = useMutation({
    mutationFn: async ({ id, data }: { id: number; data: any }) => {
      const res = await api.put(`/admin/users/${id}`, data);
      return res.data;
    },
    onSuccess: () => {
      toast.success("User updated successfully");
      setIsModalOpen(false);
      setEditingUser(null);
      queryClient.invalidateQueries({ queryKey: ["admin-users"] });
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error || "Failed to update user");
    },
  });

  const deleteMutation = useMutation({
    mutationFn: async (id: number) => {
      const res = await api.delete(`/admin/users/${id}`);
      return res.data;
    },
    onSuccess: () => {
      toast.success("User deleted successfully");
      queryClient.invalidateQueries({ queryKey: ["admin-users"] });
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error || "Failed to delete user");
    },
  });

  const handleEdit = (user: User) => {
    setEditingUser(user);
    setPassword("");
    setIsModalOpen(true);
  };

  const handleDelete = (id: number) => {
    if (id === currentUser?.id) {
      toast.error("You cannot delete yourself");
      return;
    }
    if (window.confirm("Are you sure you want to delete this user?")) {
      deleteMutation.mutate(id);
    }
  };

  const handleAddNew = () => {
    setEditingUser({
      username: "",
      name: "",
      is_admin: false,
    });
    setPassword("");
    setIsModalOpen(true);
  };

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const formData = new FormData(e.currentTarget);
    const data: any = Object.fromEntries(formData.entries());

    const userData = {
      username: data.username,
      name: data.name,
      is_admin: data.is_admin === "on",
    };

    if (editingUser?.id) {
      updateMutation.mutate({ id: editingUser.id, data: userData });
    } else {
      createMutation.mutate({ ...userData, password });
    }
  };

  if (!currentUser) return null;

  return (
    <div className="container mx-auto py-10 px-4">
      <div className="flex justify-between items-center mb-8">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">User Management</h1>
          <p className="text-gray-500">
            Manage administrative and standard users.
          </p>
        </div>
        <Button onClick={handleAddNew}>Add New User</Button>
      </div>

      <div className="bg-white rounded-xl border shadow-sm overflow-hidden mb-6">
        <Table>
          <TableHeader className="bg-gray-50">
            <TableRow>
              <TableHead>Username</TableHead>
              <TableHead>Name</TableHead>
              <TableHead>Role</TableHead>
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
            ) : users.length > 0 ? (
              users.map((u: User) => (
                <TableRow key={u.id}>
                  <TableCell className="font-medium">{u.username}</TableCell>
                  <TableCell>{u.name || "-"}</TableCell>
                  <TableCell>
                    <span
                      className={`px-2 py-1 rounded-full text-[10px] font-bold ${
                        u.is_admin
                          ? "bg-purple-100 text-purple-700"
                          : "bg-blue-100 text-blue-700"
                      }`}
                    >
                      {u.is_admin ? "ADMIN" : "USER"}
                    </span>
                  </TableCell>
                  <TableCell className="text-right">
                    <div className="flex justify-end gap-2">
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={() => handleEdit(u)}
                      >
                        Edit
                      </Button>
                      <Button
                        variant="destructive"
                        size="sm"
                        onClick={() => handleDelete(u.id)}
                        disabled={u.id === currentUser.id}
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
                  No users found.
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>

      {isModalOpen && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-xl shadow-xl w-full max-w-md overflow-hidden">
            <div className="p-6 border-b flex justify-between items-center">
              <h2 className="text-xl font-bold">
                {editingUser?.id ? "Edit User" : "Add New User"}
              </h2>
              <button
                onClick={() => setIsModalOpen(false)}
                className="text-gray-400 hover:text-gray-600"
              >
                &times;
              </button>
            </div>
            <form onSubmit={handleSubmit} className="p-6 space-y-4">
              <div className="space-y-1">
                <label className="text-xs font-bold uppercase text-gray-500">
                  Username
                </label>
                <input
                  name="username"
                  defaultValue={editingUser?.username}
                  className="w-full border rounded px-3 py-2"
                  required
                />
              </div>
              <div className="space-y-1">
                <label className="text-xs font-bold uppercase text-gray-500">
                  Full Name
                </label>
                <input
                  name="name"
                  defaultValue={editingUser?.name}
                  className="w-full border rounded px-3 py-2"
                />
              </div>
              {!editingUser?.id && (
                <div className="space-y-1">
                  <label className="text-xs font-bold uppercase text-gray-500">
                    Password
                  </label>
                  <input
                    type="password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    className="w-full border rounded px-3 py-2"
                    required
                    minLength={6}
                  />
                </div>
              )}
              <div className="flex items-center gap-2 py-2">
                <input
                  type="checkbox"
                  id="is_admin"
                  name="is_admin"
                  defaultChecked={editingUser?.is_admin}
                  className="w-4 h-4 text-blue-600"
                />
                <label htmlFor="is_admin" className="text-sm font-medium">
                  Administrator Privileges
                </label>
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
                  {editingUser?.id ? "Update User" : "Create User"}
                </Button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}

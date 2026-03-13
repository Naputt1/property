import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useState, useEffect } from "react";
import { useMutation } from "@tanstack/react-query";
import { useAuthStore } from "@/store/auth";
import { initAxios } from "@/services/axios";
import { Button } from "@/components/ui/button";
import { toast } from "sonner";

export const Route = createFileRoute("/profile/")({
  component: Profile,
});

function Profile() {
  const navigate = useNavigate();
  const { user, clearAuth } = useAuthStore();
  const [oldPassword, setOldPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");

  const api = initAxios();

  useEffect(() => {
    if (!user) {
      navigate({ to: "/login" });
    }
  }, [user, navigate]);

  const changePasswordMutation = useMutation({
    mutationFn: async (data: any) => {
      const res = await api.post("/user/change-password", data);
      return res.data;
    },
    onSuccess: () => {
      toast.success("Password changed successfully. Please log in again.");
      clearAuth();
      navigate({ to: "/login" });
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error || "Failed to change password");
    },
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (newPassword !== confirmPassword) {
      toast.error("New passwords do not match");
      return;
    }
    changePasswordMutation.mutate({
      old_password: oldPassword,
      new_password: newPassword,
    });
  };

  const handleLogout = () => {
    clearAuth();
    navigate({ to: "/login" });
  };

  if (!user) return null;

  return (
    <div className="container mx-auto py-10 px-4 max-w-md">
      <h1 className="text-3xl font-bold text-gray-900 mb-8">User Profile</h1>
      
      <div className="bg-white rounded-xl border shadow-sm p-6 mb-8">
        <h2 className="text-xl font-semibold mb-4">Account Information</h2>
        <div className="space-y-4">
          <div>
            <label className="text-xs font-bold uppercase text-gray-500">Username</label>
            <p className="text-gray-900 font-medium">{user.username}</p>
          </div>
          <div>
            <label className="text-xs font-bold uppercase text-gray-500">Role</label>
            <p className="text-gray-900 font-medium">{user.is_admin ? "Administrator" : "Standard User"}</p>
          </div>
        </div>
      </div>

      <div className="bg-white rounded-xl border shadow-sm p-6 mb-8">
        <h2 className="text-xl font-semibold mb-4">Change Password</h2>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-1">
            <label className="text-xs font-bold uppercase text-gray-500">Current Password</label>
            <input
              type="password"
              value={oldPassword}
              onChange={(e) => setOldPassword(e.target.value)}
              className="w-full border rounded px-3 py-2"
              required
            />
          </div>
          <div className="space-y-1">
            <label className="text-xs font-bold uppercase text-gray-500">New Password</label>
            <input
              type="password"
              value={newPassword}
              onChange={(e) => setNewPassword(e.target.value)}
              className="w-full border rounded px-3 py-2"
              required
              minLength={6}
            />
          </div>
          <div className="space-y-1">
            <label className="text-xs font-bold uppercase text-gray-500">Confirm New Password</label>
            <input
              type="password"
              value={confirmPassword}
              onChange={(e) => setConfirmPassword(e.target.value)}
              className="w-full border rounded px-3 py-2"
              required
              minLength={6}
            />
          </div>
          <Button 
            type="submit" 
            className="w-full"
            disabled={changePasswordMutation.isPending}
          >
            {changePasswordMutation.isPending ? "Updating..." : "Update Password"}
          </Button>
        </form>
      </div>

      <div className="bg-white rounded-xl border border-red-100 shadow-sm p-6">
        <h2 className="text-xl font-semibold mb-4 text-red-600">Danger Zone</h2>
        <p className="text-sm text-gray-500 mb-4">
          Once you log out, you will need to enter your credentials again to access the system.
        </p>
        <Button 
          variant="destructive" 
          className="w-full"
          onClick={handleLogout}
        >
          Logout from Account
        </Button>
      </div>
    </div>
  );
}

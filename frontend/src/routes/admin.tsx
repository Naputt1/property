import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useEffect, useState, useCallback } from "react";
import { useAuthStore } from "@/store/auth";
import { initAxios } from "@/services/axios";
import { toast } from "sonner";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { type IJob, JobStatus } from "@/types/job";
import { useWebSocket } from "@/hooks/useWebSocket";

export const Route = createFileRoute("/admin")({
  component: Admin,
});

const Progress = ({ value, label }: { value: number; label?: string }) => (
  <div className="w-full">
    <div className="flex justify-between mb-1">
      <span className="text-[10px] font-medium text-blue-700 dark:text-white">{label}</span>
      <span className="text-[10px] font-medium text-blue-700 dark:text-white">{Math.round(value)}%</span>
    </div>
    <div className="w-full bg-gray-200 rounded-full h-1.5 dark:bg-gray-700">
      <div 
        className="bg-blue-600 h-1.5 rounded-full transition-all duration-500 ease-out" 
        style={{ width: `${Math.min(100, Math.max(0, value))}%` }}
      ></div>
    </div>
  </div>
);

function Admin() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { user, clearAuth } = useAuthStore();
  const [file, setFile] = useState<File | null>(null);
  const [uploading, setUploading] = useState(false);
  const [uploadProgress, setUploadProgress] = useState(0);

  useEffect(() => {
    if (!user) {
      navigate({ to: "/login" });
    }
  }, [user, navigate]);

  const { data: jobs, refetch } = useQuery({
    queryKey: ["jobs"],
    queryFn: async () => {
      const api = initAxios();
      const res = await api.get("/admin/jobs?limit=10");
      return res.data.data as IJob[];
    },
    enabled: !!user,
  });

  // Handle real-time job updates via WebSocket
  const onWsMessage = useCallback((msg: any) => {
    if (msg.type === 'job_update') {
      const update = msg.data;
      
      // Update the cache directly for immediate UI feedback
      queryClient.setQueryData(["jobs"], (oldJobs: IJob[] | undefined) => {
        if (!oldJobs) return oldJobs;
        
        return oldJobs.map(job => {
          if (job.id === update.id) {
            return {
              ...job,
              status: update.status,
              progress: update.progress,
              total: update.total
            };
          }
          return job;
        });
      });

      // If job finished, refetch to get the latest message/result
      if (update.status === JobStatus.SUCCESS || update.status === JobStatus.FAILED) {
        refetch();
      }
    }
  }, [queryClient, refetch]);

  useWebSocket(onWsMessage);

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      setFile(e.target.files[0]);
    }
  };

  const handleUpload = async () => {
    if (!file) return;
    setUploading(true);
    setUploadProgress(0);

    try {
      const api = initAxios();
      const res = await api.post(`/admin/stream-upload?filename=${encodeURIComponent(file.name)}`, file, {
        headers: {
          "Content-Type": "application/octet-stream",
        },
        onUploadProgress: (progressEvent) => {
          const percentCompleted = Math.round((progressEvent.loaded * 100) / (progressEvent.total || 1));
          setUploadProgress(percentCompleted);
        },
        timeout: 0, 
      });

      if (res.data.status) {
        toast.success("Stream migration job queued!");
        setFile(null);
        setUploadProgress(0);
        // Add a small delay to ensure the backend has processed the job creation
        setTimeout(refetch, 500);
      } else {
        toast.error(res.data.error || "Upload failed");
      }
    } catch (err: any) {
      toast.error(err.response?.data?.error || "Upload failed");
    } finally {
      setUploading(false);
    }
  };

  const handleLogout = () => {
    clearAuth();
    navigate({ to: "/login" });
  };

  if (!user) return null;

  return (
    <div className="space-y-8">
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold">Admin Dashboard</h1>
          <p className="text-gray-500">Welcome back, {user?.username}</p>
        </div>
        <button
          onClick={handleLogout}
          className="text-sm font-medium text-red-600 hover:text-red-700"
        >
          Logout
        </button>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        {/* Upload Section */}
        <div className="lg:col-span-1 bg-white p-6 rounded-xl border shadow-sm h-fit">
          <h2 className="text-xl font-semibold mb-4">Upload Property Data</h2>
          <p className="text-sm text-gray-500 mb-6">
            Upload a CSV file containing UK Land Registry property data.
          </p>
          
          <div className="space-y-4">
            <div className="border-2 border-dashed border-gray-200 rounded-lg p-6 text-center">
              <input
                type="file"
                id="csv-upload"
                className="hidden"
                accept=".csv"
                onChange={handleFileChange}
              />
              <label htmlFor="csv-upload" className="cursor-pointer block">
                <div className="text-blue-600 font-medium mb-1 truncate px-2">
                  {file ? file.name : "Click to select file"}
                </div>
                <div className="text-xs text-gray-400">CSV files only</div>
              </label>
            </div>
            
            <button
              onClick={handleUpload}
              disabled={!file || uploading}
              className="w-full bg-blue-600 text-white font-semibold py-2 rounded-lg hover:bg-blue-700 disabled:opacity-50 transition-colors"
            >
              {uploading ? "Uploading..." : "Start Migration"}
            </button>

            {uploading && (
              <div className="mt-4">
                <Progress value={uploadProgress} label="Uploading to storage..." />
              </div>
            )}
          </div>
        </div>

        {/* Jobs List Section */}
        <div className="lg:col-span-2 bg-white p-6 rounded-xl border shadow-sm">
          <h2 className="text-xl font-semibold mb-4">Migration Jobs</h2>
          <div className="overflow-x-auto">
            <table className="w-full text-sm text-left">
              <thead className="bg-gray-50 text-gray-600 uppercase text-xs font-medium">
                <tr>
                  <th className="px-4 py-3">Job ID</th>
                  <th className="px-4 py-3">Status</th>
                  <th className="px-4 py-3">Progress</th>
                  <th className="px-4 py-3">Task</th>
                  <th className="px-4 py-3 text-right">Date</th>
                </tr>
              </thead>
              <tbody className="divide-y">
                {jobs && jobs.length > 0 ? (
                  jobs.map((job) => (
                    <tr key={job.id} className="hover:bg-gray-50">
                      <td className="px-4 py-3 font-mono text-xs truncate max-w-[120px]">
                        {job.id}
                      </td>
                      <td className="px-4 py-3">
                        <span className={`px-2 py-1 rounded-full text-[10px] font-bold ${
                          job.status === JobStatus.SUCCESS ? "bg-green-100 text-green-700" :
                          job.status === JobStatus.FAILED ? "bg-red-100 text-red-700" :
                          job.status === JobStatus.RUNNING ? "bg-blue-100 text-blue-700" :
                          "bg-gray-100 text-gray-700"
                        }`}>
                          {job.status}
                        </span>
                      </td>
                      <td className="px-4 py-3 min-w-[150px]">
                        {job.status === JobStatus.RUNNING ? (
                          <Progress value={job.progress} label={`${job.total.toLocaleString()} records`} />
                        ) : job.status === JobStatus.SUCCESS ? (
                          <div className="text-xs text-green-600 font-medium">
                            Completed: {job.total.toLocaleString()} records
                          </div>
                        ) : job.message ? (
                          <div className="text-xs text-gray-500 truncate max-w-[150px]" title={job.message}>
                            {job.message}
                          </div>
                        ) : "-"}
                      </td>
                      <td className="px-4 py-3 text-gray-500">{job.task_type}</td>
                      <td className="px-4 py-3 text-right text-gray-400 text-xs">
                        {new Date(job.created_at).toLocaleString()}
                      </td>
                    </tr>
                  ))
                ) : (
                  <tr>
                    <td colSpan={5} className="px-4 py-8 text-center text-gray-400">
                      No jobs found.
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  );
}

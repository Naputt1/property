export enum JobStatus {
  PENDING = "PENDING",
  RUNNING = "RUNNING",
  SUCCESS = "SUCCESS",
  FAILED = "FAILED",
}

export interface IJob {
  id: string;
  task_type: string;
  status: JobStatus;
  message: string;
  progress: number;
  total: number;
  created_at: string;
  updated_at: string;
}

export interface IJobListResponse {
  status: boolean;
  data: IJob[];
  total: number;
}

export type JobStatus = 'discovered' | 'qualified' | 'review' | 'applied' | 'outcome_interview' | 'outcome_rejected';
export type AgentMode = 'review-each' | 'routine-auto';

export interface Job {
  id: number;
  company: string;
  role: string;
  url?: string;
  status: JobStatus;
  score?: number;
  applied_at?: string;
  outcome?: string;
  notes?: string;
  created_at: string;
}

export interface Action {
  id: number;
  type: string;
  message: string;
  job_id?: number;
  created_at: string;
}

export interface Run {
  id: string;
  started_at: string;
  ended_at?: string;
  total_found: number;
  total_qualified: number;
  total_applied: number;
  total_skipped: number;
}

export interface AgentStatus {
  id: number;
  mode: AgentMode;
  ai_provider: string;
  is_running: number;
  last_action?: string;
  updated_at: string;
}

export interface Stats {
  discovered: number;
  qualified: number;
  review: number;
  applied: number;
  outcome_interview: number;
  outcome_rejected: number;
  total: number;
}

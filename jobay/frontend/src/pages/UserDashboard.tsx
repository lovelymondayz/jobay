import { useState, useEffect } from 'react';
import { Briefcase, Search, CheckCircle2, AlertTriangle, Zap, Activity, RefreshCw, ChevronDown, ExternalLink, FileText } from 'lucide-react';
import type { Job, Action, Stats } from '../types';

const statusColors: Record<string, string> = {
  discovered: 'text-blue-400 bg-blue-400/10',
  qualified: 'text-emerald-400 bg-emerald-400/10',
  review: 'text-amber-400 bg-amber-400/10',
  applied: 'text-purple-400 bg-purple-400/10',
  outcome_interview: 'text-cyan-400 bg-cyan-400/10',
  outcome_rejected: 'text-red-400 bg-red-400/10',
};

const statusLabels: Record<string, string> = {
  discovered: 'Discovered',
  qualified: 'Qualified',
  review: 'Review',
  applied: 'Applied',
  outcome_interview: 'Interview',
  outcome_rejected: 'Rejected',
};

function StatusBadge({ status }: { status: string }) {
  return (
    <span className={`px-2 py-0.5 rounded text-xs font-medium ${statusColors[status] || 'text-gray-400 bg-gray-400/10'}`}>
      {statusLabels[status] || status}
    </span>
  );
}

function StatsCard({ label, value, icon: Icon, color }: { label: string; value: number; icon: React.ElementType; color: string }) {
  return (
    <div className="bg-gray-900 rounded-xl p-4 flex items-center gap-4">
      <div className={`p-3 rounded-lg ${color}`}>
        <Icon size={20} />
      </div>
      <div>
        <div className="text-2xl font-bold">{value}</div>
        <div className="text-sm text-gray-400">{label}</div>
      </div>
    </div>
  );
}

function JobRow({ job }: { job: Job }) {
  const [expanded, setExpanded] = useState(false);
  return (
    <div className="bg-gray-900 rounded-lg border border-gray-800 hover:border-gray-700 transition">
      <div className="p-4 flex items-center gap-4 cursor-pointer" onClick={() => setExpanded(!expanded)}>
        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-2 mb-1">
            <span className="font-medium truncate">{job.company}</span>
            <StatusBadge status={job.status} />
          </div>
          <div className="text-sm text-gray-400 truncate">{job.role}</div>
        </div>
        {job.score !== null && job.score !== undefined && (
          <div className="text-right">
            <div className={`text-lg font-bold ${job.score >= 80 ? 'text-green-400' : job.score >= 60 ? 'text-amber-400' : 'text-red-400'}`}>
              {job.score}
            </div>
            <div className="text-xs text-gray-500">score</div>
          </div>
        )}
        <ChevronDown size={16} className={`text-gray-500 transition-transform ${expanded ? 'rotate-180' : ''}`} />
      </div>
      {expanded && (
        <div className="px-4 pb-4 border-t border-gray-800 pt-3 space-y-3">
          {job.url && (
            <a href={job.url} target="_blank" rel="noopener noreferrer" className="text-sm text-blue-400 hover:text-blue-300 flex items-center gap-1">
              <ExternalLink size={14} /> {job.url}
            </a>
          )}
          {job.location && (
            <div className="text-sm text-gray-400">
              <span className="text-gray-500 text-xs">Location:</span> {job.location}
            </div>
          )}
          {job.source && (
            <div className="text-sm text-gray-400">
              <span className="text-gray-500 text-xs">Source:</span> {job.source}
            </div>
          )}
          {job.apply_status && (
            <div className="text-sm">
              <span className="text-gray-500 text-xs">Apply Status:</span>{' '}
              <span className={job.apply_status === 'applied' ? 'text-green-400' : 'text-red-400'}>
                {job.apply_status}
              </span>
            </div>
          )}
          {job.match_reasons && (
            <div className="text-sm text-gray-400">
              <span className="text-gray-500 text-xs">Match:</span> {job.match_reasons}
            </div>
          )}
          {job.notes && (
            <div className="text-sm text-gray-400">
              <span className="text-gray-500 text-xs">Notes:</span> {job.notes}
            </div>
          )}
          <div className="text-xs text-gray-600">
            Added: {new Date(job.created_at).toLocaleDateString()}
          </div>
        </div>
      )}
    </div>
  );
}

function ActionsLog({ actions }: { actions: Action[] }) {
  const iconMap: Record<string, React.ElementType> = {
    discover: Search,
    score: Activity,
    apply: CheckCircle2,
    review: AlertTriangle,
    outcome: Activity,
  };
  return (
    <div className="space-y-2 max-h-96 overflow-y-auto">
      {actions.map((action) => {
        const Icon = iconMap[action.type] || Activity;
        return (
          <div key={action.id} className="flex items-start gap-3 p-2 rounded hover:bg-gray-900">
            <Icon size={16} className="mt-0.5 text-gray-500 shrink-0" />
            <div className="flex-1 min-w-0">
              <div className="text-sm">{action.message}</div>
              <div className="text-xs text-gray-500">{new Date(action.created_at).toLocaleString()}</div>
            </div>
          </div>
        );
      })}
    </div>
  );
}

interface UserData {
  slug: string;
  name: string;
  cv_path?: string;
  created_at: string;
}

export default function UserDashboard({ slug }: { slug: string }) {
  const [user, setUser] = useState<UserData | null>(null);
  const [jobs, setJobs] = useState<Job[]>([]);
  const [actions, setActions] = useState<Action[]>([]);
  const [stats, setStats] = useState<Stats | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [userRes, jobsRes, actionsRes, statsRes] = await Promise.all([
          fetch(`/api/users/${slug}`),
          fetch(`/api/users/${slug}/jobs`),
          fetch(`/api/users/${slug}/actions`),
          fetch(`/api/users/${slug}/stats`),
        ]);

        if (!userRes.ok) throw new Error('User not found');

        setUser(await userRes.json());
        setJobs((await jobsRes.json()).jobs);
        setActions(await actionsRes.json());
        setStats(await statsRes.json());
      } catch (err: unknown) {
        setError(err instanceof Error ? err.message : 'Failed to load');
      } finally {
        setLoading(false);
      }
    };
    fetchData();
  }, [slug]);

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-950 text-gray-100 flex items-center justify-center">
        <div className="text-center space-y-4">
          <RefreshCw size={32} className="animate-spin mx-auto" />
          <p className="text-gray-400">Loading...</p>
        </div>
      </div>
    );
  }

  if (error || !user) {
    return (
      <div className="min-h-screen bg-gray-950 text-gray-100 flex items-center justify-center">
        <div className="text-center space-y-4">
          <AlertTriangle size={48} className="text-red-400 mx-auto" />
          <h2 className="text-xl font-bold">User not found</h2>
          <p className="text-gray-400">{error || `No user with slug "${slug}"`}</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-950 text-gray-100">
      <header className="border-b border-gray-800 bg-gray-950 sticky top-0 z-50">
        <div className="max-w-7xl mx-auto px-4 py-3 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="w-8 h-8 bg-gradient-to-br from-blue-500 to-purple-600 rounded-lg flex items-center justify-center">
              <Briefcase size={18} className="text-white" />
            </div>
            <div>
              <h1 className="text-lg font-bold">{user.name}</h1>
              <p className="text-xs text-gray-500">Job Hunt Progress</p>
            </div>
          </div>
          <div className="flex items-center gap-3">
            {user.cv_path && (
              <a href={user.cv_path} target="_blank" rel="noopener noreferrer" className="p-2 rounded-lg hover:bg-gray-800 transition">
                <FileText size={16} />
              </a>
            )}
          </div>
        </div>
      </header>

      <main className="max-w-7xl mx-auto px-4 py-6 space-y-6">
        {stats && (
          <div className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-7 gap-3">
            <StatsCard label="Total" value={stats.total} icon={Briefcase} color="bg-blue-500/20 text-blue-400" />
            <StatsCard label="Discovered" value={stats.discovered} icon={Search} color="bg-blue-500/20 text-blue-400" />
            <StatsCard label="Qualified" value={stats.qualified} icon={CheckCircle2} color="bg-emerald-500/20 text-emerald-400" />
            <StatsCard label="Review" value={stats.review} icon={AlertTriangle} color="bg-amber-500/20 text-amber-400" />
            <StatsCard label="Applied" value={stats.applied} icon={Zap} color="bg-purple-500/20 text-purple-400" />
            <StatsCard label="Interview" value={stats.outcome_interview} icon={Activity} color="bg-cyan-500/20 text-cyan-400" />
            <StatsCard label="Rejected" value={stats.outcome_rejected} icon={AlertTriangle} color="bg-red-500/20 text-red-400" />
          </div>
        )}

        <div className="grid lg:grid-cols-2 gap-6">
          <div className="bg-gray-900 rounded-xl p-4">
            <h3 className="text-sm font-medium mb-4">Activity Log</h3>
            <ActionsLog actions={actions} />
          </div>

          <div className="bg-gray-900 rounded-xl p-4">
            <h3 className="text-sm font-medium mb-4">Jobs ({jobs.length})</h3>
            <div className="space-y-2">
              {jobs.length > 0 ? jobs.map(job => (
                <JobRow key={job.id} job={job} />
              )) : (
                <div className="text-center py-8 text-gray-600">No jobs yet</div>
              )}
            </div>
          </div>
        </div>
      </main>
    </div>
  );
}

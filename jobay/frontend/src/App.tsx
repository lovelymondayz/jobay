import { useState, useEffect, useRef, useCallback } from 'react';
import { BarChart3, Briefcase, CheckCircle2, Clock, AlertTriangle, Activity, Settings, Search, Plus, RefreshCw, Trash2, ExternalLink, Zap, Eye, Play, Square, ChevronDown, Filter } from 'lucide-react';
import { LineChart, Line, BarChart, Bar, XAxis, YAxis, Tooltip, ResponsiveContainer, PieChart, Pie, Cell } from 'recharts';
import type { Job, Action, Run, AgentStatus, Stats, JobStatus, AgentMode } from './types';

const API = '/api';

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

function AgentControl({ agent, onModeChange, onToggle }: { agent: AgentStatus | null; onModeChange: (mode: AgentMode) => void; onToggle: () => void }) {
  if (!agent) return null;
  const isRunning = agent.is_running === 1;
  return (
    <div className="bg-gray-900 rounded-xl p-4 flex items-center justify-between">
      <div className="flex items-center gap-3">
        <div className={`w-3 h-3 rounded-full ${isRunning ? 'bg-green-500 animate-pulse' : 'bg-gray-600'}`} />
        <div>
          <div className="font-medium">Agent: {isRunning ? 'Running' : 'Idle'}</div>
          <div className="text-xs text-gray-400">Mode: {agent.mode} · Provider: {agent.ai_provider}</div>
        </div>
      </div>
      <div className="flex items-center gap-2">
        <select
          value={agent.mode}
          onChange={(e) => onModeChange(e.target.value as AgentMode)}
          className="bg-gray-800 text-sm rounded px-2 py-1 border border-gray-700 focus:outline-none focus:border-blue-500"
        >
          <option value="review-each">Review Each</option>
          <option value="routine-auto">Routine Auto</option>
        </select>
        <button
          onClick={onToggle}
          className={`p-2 rounded-lg transition ${isRunning ? 'bg-red-500/20 text-red-400 hover:bg-red-500/30' : 'bg-green-500/20 text-green-400 hover:bg-green-500/30'}`}
        >
          {isRunning ? <Square size={16} /> : <Play size={16} />}
        </button>
      </div>
    </div>
  );
}

function JobRow({ job, onUpdate }: { job: Job; onUpdate: (id: number, data: Partial<Job>) => void }) {
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
          <div className="grid grid-cols-2 gap-3 text-sm">
            <div>
              <label className="text-gray-500 text-xs">Status</label>
              <select
                value={job.status}
                onChange={(e) => onUpdate(job.id, { status: e.target.value as JobStatus })}
                className="w-full mt-1 bg-gray-800 rounded px-2 py-1 border border-gray-700 text-sm focus:outline-none focus:border-blue-500"
              >
                {Object.keys(statusLabels).map(s => (
                  <option key={s} value={s}>{statusLabels[s]}</option>
                ))}
              </select>
            </div>
            <div>
              <label className="text-gray-500 text-xs">Score</label>
              <input
                type="number"
                value={job.score ?? ''}
                onChange={(e) => onUpdate(job.id, { score: e.target.value ? Number(e.target.value) : undefined })}
                className="w-full mt-1 bg-gray-800 rounded px-2 py-1 border border-gray-700 text-sm focus:outline-none focus:border-blue-500"
                placeholder="—"
              />
            </div>
            <div className="col-span-2">
              <label className="text-gray-500 text-xs">Notes</label>
              <textarea
                value={job.notes || ''}
                onChange={(e) => onUpdate(job.id, { notes: e.target.value })}
                className="w-full mt-1 bg-gray-800 rounded px-2 py-1 border border-gray-700 text-sm focus:outline-none focus:border-blue-500"
                rows={2}
                placeholder="Add notes..."
              />
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

function ActionsLog({ actions }: { actions: Action[] }) {
  const iconMap: Record<string, React.ElementType> = {
    discover: Search,
    score: BarChart3,
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
              <div className="text-xs text-gray-500">{new Date(action.created_at + 'Z').toLocaleString()}</div>
            </div>
          </div>
        );
      })}
    </div>
  );
}

export default function App() {
  const [jobs, setJobs] = useState<Job[]>([]);
  const [actions, setActions] = useState<Action[]>([]);
  const [agent, setAgent] = useState<AgentStatus | null>(null);
  const [runs, setRuns] = useState<Run[]>([]);
  const [stats, setStats] = useState<Stats | null>(null);
  const [wsConnected, setWsConnected] = useState(false);
  const [activeTab, setActiveTab] = useState<'overview' | 'jobs' | 'runs'>('overview');
  const [statusFilter, setStatusFilter] = useState<string>('all');
  const [searchQuery, setSearchQuery] = useState('');
  const wsRef = useRef<WebSocket | null>(null);

  const fetchData = useCallback(async () => {
    try {
      const [statusRes, actionsRes, agentRes, runsRes] = await Promise.all([
        fetch(`${API}/status`),
        fetch(`${API}/actions`),
        fetch(`${API}/agent`),
        fetch(`${API}/runs`),
      ]);
      const statusData = await statusRes.json();
      setJobs(statusData.jobs);
      setStats(statusData.stats);
      setAgent(statusData.agent);
      setRuns(statusData.runs);
      setActions(await actionsRes.json());
      setAgent(await agentRes.json());
    } catch (err) {
      console.error('Fetch error:', err);
    }
  }, []);

  useEffect(() => {
    fetchData();
    const ws = new WebSocket(`ws://${window.location.host}/ws`);
    wsRef.current = ws;
    ws.onopen = () => setWsConnected(true);
    ws.onclose = () => {
      setWsConnected(false);
      setTimeout(() => {
        if (wsRef.current === ws) {
          wsRef.current = null;
        }
      }, 5000);
    };
    ws.onmessage = (event) => {
      const msg = JSON.parse(event.data);
      if (msg.type === 'init') {
        setJobs(msg.data.jobs);
        setStats(msg.data.stats);
        setAgent(msg.data.agent);
        setRuns(msg.data.runs);
      } else if (msg.type === 'job_added' || msg.type === 'job_updated' || msg.type === 'job_deleted') {
        fetchData();
      } else if (msg.type === 'action_added') {
        fetchData();
      } else if (msg.type === 'agent_updated') {
        setAgent(msg.data);
      }
    };
    return () => ws.close();
  }, [fetchData]);

  const handleModeChange = async (mode: AgentMode) => {
    const res = await fetch(`${API}/agent/mode`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ mode }),
    });
    setAgent(await res.json());
  };

  const handleAgentToggle = async () => {
    const res = await fetch(`${API}/agent/toggle`, { method: 'POST' });
    setAgent(await res.json());
  };

  const handleJobUpdate = async (id: number, data: Partial<Job>) => {
    await fetch(`${API}/jobs/${id}`, {
      method: 'PATCH',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    });
    fetchData();
  };

  const handleJobDelete = async (id: number) => {
    await fetch(`${API}/jobs/${id}`, { method: 'DELETE' });
    fetchData();
  };

  const handleAddJob = async () => {
    const company = prompt('Company:');
    if (!company) return;
    const role = prompt('Role:');
    if (!role) return;
    await fetch(`${API}/jobs`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ company, role }),
    });
    fetchData();
  };

  const filteredJobs = jobs.filter((j) => {
    if (statusFilter !== 'all' && j.status !== statusFilter) return false;
    if (searchQuery) {
      const q = searchQuery.toLowerCase();
      return j.company.toLowerCase().includes(q) || j.role.toLowerCase().includes(q);
    }
    return true;
  });

  const chartData = stats ? [
    { name: 'Discovered', value: stats.discovered, color: '#60a5fa' },
    { name: 'Qualified', value: stats.qualified, color: '#34d399' },
    { name: 'Review', value: stats.review, color: '#fbbf24' },
    { name: 'Applied', value: stats.applied, color: '#a78bfa' },
    { name: 'Interview', value: stats.outcome_interview, color: '#22d3ee' },
    { name: 'Rejected', value: stats.outcome_rejected, color: '#f87171' },
  ].filter(d => d.value > 0) : [];

  const PIEColors = ['#60a5fa', '#34d399', '#fbbf24', '#a78bfa', '#22d3ee', '#f87171'];

  return (
    <div className="min-h-screen bg-gray-950 text-gray-100">
      {/* Header */}
      <header className="border-b border-gray-800 bg-gray-950 sticky top-0 z-50">
        <div className="max-w-7xl mx-auto px-4 py-3 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="w-8 h-8 bg-gradient-to-br from-blue-500 to-purple-600 rounded-lg flex items-center justify-center">
              <Briefcase size={18} className="text-white" />
            </div>
            <div>
              <h1 className="text-lg font-bold">Jobay</h1>
              <p className="text-xs text-gray-500">Job Application Agent Dashboard</p>
            </div>
          </div>
          <div className="flex items-center gap-3">
            <div className={`flex items-center gap-2 text-xs ${wsConnected ? 'text-green-400' : 'text-red-400'}`}>
              <div className={`w-2 h-2 rounded-full ${wsConnected ? 'bg-green-400' : 'bg-red-400'}`} />
              {wsConnected ? 'Live' : 'Disconnected'}
            </div>
            <button onClick={fetchData} className="p-2 rounded-lg hover:bg-gray-800 transition">
              <RefreshCw size={16} />
            </button>
          </div>
        </div>
      </header>

      <main className="max-w-7xl mx-auto px-4 py-6 space-y-6">
        {/* Agent Control */}
        <AgentControl agent={agent} onModeChange={handleModeChange} onToggle={handleAgentToggle} />

        {/* Stats Grid */}
        <div className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-7 gap-3">
          <StatsCard label="Total" value={stats?.total || 0} icon={Briefcase} color="bg-blue-500/20 text-blue-400" />
          <StatsCard label="Discovered" value={stats?.discovered || 0} icon={Search} color="bg-blue-500/20 text-blue-400" />
          <StatsCard label="Qualified" value={stats?.qualified || 0} icon={CheckCircle2} color="bg-emerald-500/20 text-emerald-400" />
          <StatsCard label="Review" value={stats?.review || 0} icon={AlertTriangle} color="bg-amber-500/20 text-amber-400" />
          <StatsCard label="Applied" value={stats?.applied || 0} icon={Zap} color="bg-purple-500/20 text-purple-400" />
          <StatsCard label="Interview" value={stats?.outcome_interview || 0} icon={Activity} color="bg-cyan-500/20 text-cyan-400" />
          <StatsCard label="Rejected" value={stats?.outcome_rejected || 0} icon={AlertTriangle} color="bg-red-500/20 text-red-400" />
        </div>

        {/* Tabs */}
        <div className="flex gap-1 border-b border-gray-800">
          {(['overview', 'jobs', 'runs'] as const).map(tab => (
            <button
              key={tab}
              onClick={() => setActiveTab(tab)}
              className={`px-4 py-2 text-sm font-medium capitalize transition border-b-2 ${activeTab === tab ? 'border-blue-500 text-blue-400' : 'border-transparent text-gray-500 hover:text-gray-300'}`}
            >
              {tab}
            </button>
          ))}
        </div>

        {/* Tab Content */}
        {activeTab === 'overview' && (
          <div className="grid lg:grid-cols-2 gap-6">
            {/* Pie Chart */}
            <div className="bg-gray-900 rounded-xl p-4">
              <h3 className="text-sm font-medium mb-4">Pipeline Distribution</h3>
              {chartData.length > 0 ? (
                <ResponsiveContainer width="100%" height={250}>
                  <PieChart>
                    <Pie data={chartData} dataKey="value" nameKey="name" cx="50%" cy="50%" outerRadius={80} label={({ name, value }) => `${name}: ${value}`}>
                      {chartData.map((_, i) => (
                        <Cell key={i} fill={PIEColors[i % PIEColors.length]} />
                      ))}
                    </Pie>
                    <Tooltip contentStyle={{ backgroundColor: '#111827', border: '1px solid #374151', borderRadius: '8px' }} />
                  </PieChart>
                </ResponsiveContainer>
              ) : (
                <div className="flex items-center justify-center h-[250px] text-gray-600">No data</div>
              )}
            </div>

            {/* Actions Log */}
            <div className="bg-gray-900 rounded-xl p-4">
              <h3 className="text-sm font-medium mb-4">Activity Log</h3>
              <ActionsLog actions={actions} />
            </div>
          </div>
        )}

        {activeTab === 'jobs' && (
          <div className="space-y-4">
            {/* Search & Filter */}
            <div className="flex flex-wrap gap-3 items-center">
              <div className="relative flex-1 min-w-[200px]">
                <Search size={16} className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-500" />
                <input
                  type="text"
                  placeholder="Search jobs..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="w-full bg-gray-900 border border-gray-800 rounded-lg pl-10 pr-4 py-2 text-sm focus:outline-none focus:border-blue-500"
                />
              </div>
              <div className="flex items-center gap-2">
                <Filter size={16} className="text-gray-500" />
                <select
                  value={statusFilter}
                  onChange={(e) => setStatusFilter(e.target.value)}
                  className="bg-gray-900 border border-gray-800 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-blue-500"
                >
                  <option value="all">All Status</option>
                  {Object.entries(statusLabels).map(([k, v]) => (
                    <option key={k} value={k}>{v}</option>
                  ))}
                </select>
              </div>
              <button onClick={handleAddJob} className="flex items-center gap-2 px-3 py-2 bg-blue-500/20 text-blue-400 rounded-lg hover:bg-blue-500/30 transition text-sm">
                <Plus size={16} /> Add
              </button>
            </div>

            {/* Job List */}
            <div className="space-y-2">
              {filteredJobs.length > 0 ? filteredJobs.map(job => (
                <div key={job.id} className="relative">
                  <JobRow job={job} onUpdate={handleJobUpdate} />
                  <button
                    onClick={() => handleJobDelete(job.id)}
                    className="absolute top-2 right-2 p-1 text-gray-600 hover:text-red-400 transition"
                    title="Delete"
                  >
                    <Trash2 size={14} />
                  </button>
                </div>
              )) : (
                <div className="text-center py-12 text-gray-600">No jobs found</div>
              )}
            </div>
          </div>
        )}

        {activeTab === 'runs' && (
          <div className="space-y-3">
            {runs.length > 0 ? runs.map(run => (
              <div key={run.id} className="bg-gray-900 rounded-lg p-4 flex items-center justify-between">
                <div>
                  <div className="font-medium">{run.id}</div>
                  <div className="text-xs text-gray-500">
                    Started: {new Date(run.started_at + 'Z').toLocaleString()}
                    {run.ended_at && ` · Ended: ${new Date(run.ended_at + 'Z').toLocaleString()}`}
                  </div>
                </div>
                <div className="flex gap-4 text-sm">
                  <div className="text-center">
                    <div className="font-bold">{run.total_found}</div>
                    <div className="text-xs text-gray-500">Found</div>
                  </div>
                  <div className="text-center">
                    <div className="font-bold text-emerald-400">{run.total_qualified}</div>
                    <div className="text-xs text-gray-500">Qualified</div>
                  </div>
                  <div className="text-center">
                    <div className="font-bold text-purple-400">{run.total_applied}</div>
                    <div className="text-xs text-gray-500">Applied</div>
                  </div>
                  <div className="text-center">
                    <div className="font-bold text-amber-400">{run.total_skipped}</div>
                    <div className="text-xs text-gray-500">Skipped</div>
                  </div>
                </div>
              </div>
            )) : (
              <div className="text-center py-12 text-gray-600">No runs yet</div>
            )}
          </div>
        )}
      </main>
    </div>
  );
}

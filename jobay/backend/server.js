const express = require('express');
const cors = require('cors');
const { WebSocketServer } = require('ws');
const http = require('http');
const path = require('path');
const fs = require('fs');

const app = express();
const server = http.createServer(app);
const wss = new WebSocketServer({ server, path: '/ws' });

const PORT = process.env.PORT || 3001;
const DB_PATH = process.env.DB_PATH || path.join(__dirname, 'jobay.json');

app.use(cors());
app.use(express.json());
app.use(express.static(path.join(__dirname, '..', 'frontend', 'dist')));

// ── JSON Database ─────────────────────────────────────────
function loadDB() {
  try {
    if (fs.existsSync(DB_PATH)) {
      return JSON.parse(fs.readFileSync(DB_PATH, 'utf8'));
    }
  } catch (e) {
    console.error('DB load error:', e.message);
  }
  return { jobs: [], actions: [], runs: [], agent: { id: 1, mode: 'review-each', ai_provider: '9router', is_running: 0, last_action: null, updated_at: new Date().toISOString() }, stats: {} };
}

function saveDB(data) {
  fs.writeFileSync(DB_PATH, JSON.stringify(data, null, 2));
}

let db = loadDB();

// Seed dummy data
if (!db.jobs || db.jobs.length === 0) {
  db.jobs = [
    { id: 1, company: 'TechCorp', role: 'Senior Backend Engineer', url: 'https://techcorp.com/jobs/123', status: 'applied', score: 88, notes: '', created_at: new Date(Date.now() - 86400000).toISOString() },
    { id: 2, company: 'StartupXYZ', role: 'Full Stack Developer', url: 'https://startupxyz.com/jobs/456', status: 'review', score: 75, notes: '', created_at: new Date(Date.now() - 72000000).toISOString() },
    { id: 3, company: 'BigTech Inc', role: 'Staff Engineer', url: 'https://bigtech.com/jobs/789', status: 'discovered', score: 62, notes: '', created_at: new Date(Date.now() - 50000000).toISOString() },
    { id: 4, company: 'RemoteCo', role: 'React Developer', url: 'https://remote.co/jobs/101', status: 'qualified', score: 82, notes: '', created_at: new Date(Date.now() - 40000000).toISOString() },
    { id: 5, company: 'AI Labs', role: 'ML Engineer', url: 'https://ailabs.com/jobs/202', status: 'outcome_rejected', score: 70, notes: '', created_at: new Date(Date.now() - 30000000).toISOString() },
    { id: 6, company: 'DataDriven', role: 'Data Engineer', url: 'https://datadriven.com/jobs/303', status: 'outcome_interview', score: 85, notes: '', created_at: new Date(Date.now() - 20000000).toISOString() },
    { id: 7, company: 'CloudFirst', role: 'DevOps Engineer', url: 'https://cloudfirst.com/jobs/404', status: 'discovered', score: 78, notes: '', created_at: new Date(Date.now() - 10000000).toISOString() },
    { id: 8, company: 'FinTech Pro', role: 'Backend Engineer', url: 'https://fintechpro.com/jobs/505', status: 'applied', score: 91, notes: '', created_at: new Date().toISOString() },
  ];
  db.actions = [
    { id: 1, type: 'discover', message: 'Found 3 new roles from career pages', job_id: null, created_at: new Date(Date.now() - 86400000).toISOString() },
    { id: 2, type: 'score', message: 'TechCorp Senior Backend Engineer scored 88/100', job_id: 1, created_at: new Date(Date.now() - 86300000).toISOString() },
    { id: 3, type: 'apply', message: 'Application submitted to TechCorp', job_id: 1, created_at: new Date(Date.now() - 86200000).toISOString() },
    { id: 4, type: 'score', message: 'StartupXYZ Full Stack Developer scored 75/100', job_id: 2, created_at: new Date(Date.now() - 72000000).toISOString() },
    { id: 5, type: 'review', message: 'Flagged for manual review: StartupXYZ', job_id: 2, created_at: new Date(Date.now() - 71900000).toISOString() },
    { id: 6, type: 'score', message: 'RemoteCo React Developer scored 82/100', job_id: 4, created_at: new Date(Date.now() - 40000000).toISOString() },
    { id: 7, type: 'outcome', message: 'Interview scheduled with DataDriven', job_id: 6, created_at: new Date(Date.now() - 20000000).toISOString() },
  ];
  db.runs = [
    { id: 'demo-run-1', started_at: new Date(Date.now() - 7200000).toISOString(), ended_at: new Date(Date.now() - 3600000).toISOString(), total_found: 8, total_qualified: 5, total_applied: 2, total_skipped: 1 },
  ];
  saveDB(db);
  console.log('✓ Seeded dummy data');
}

// ── Helpers ───────────────────────────────────────────────
function nextId(arr) { return arr.length === 0 ? 1 : Math.max(...arr.map(x => x.id)) + 1; }

function computeStats() {
  const s = {
    discovered: db.jobs.filter(j => j.status === 'discovered').length,
    qualified: db.jobs.filter(j => j.status === 'qualified').length,
    review: db.jobs.filter(j => j.status === 'review').length,
    applied: db.jobs.filter(j => j.status === 'applied').length,
    outcome_interview: db.jobs.filter(j => j.status === 'outcome_interview').length,
    outcome_rejected: db.jobs.filter(j => j.status === 'outcome_rejected').length,
    total: db.jobs.length,
  };
  db.stats = s;
  return s;
}

// ── WebSocket ─────────────────────────────────────────────
const clients = new Set();
wss.on('connection', (ws) => {
  clients.add(ws);
  ws.send(JSON.stringify({ type: 'init', data: getStatus() }));
  ws.on('close', () => clients.delete(ws));
});

function broadcast(data) {
  const msg = JSON.stringify(data);
  for (const c of clients) { if (c.readyState === 1) c.send(msg); }
}

function getStatus() {
  return {
    jobs: db.jobs,
    actions: db.actions.sort((a, b) => new Date(b.created_at) - new Date(a.created_at)).slice(0, 50),
    agent: db.agent,
    runs: db.runs,
    stats: computeStats(),
  };
}

// ── REST API ──────────────────────────────────────────────
app.get('/api/status', (req, res) => res.json(getStatus()));

app.get('/api/jobs', (req, res) => {
  let jobs = [...db.jobs];
  const { status, search } = req.query;
  if (status) jobs = jobs.filter(j => j.status === status);
  if (search) {
    const q = search.toLowerCase();
    jobs = jobs.filter(j => j.company.toLowerCase().includes(q) || j.role.toLowerCase().includes(q));
  }
  jobs.sort((a, b) => new Date(b.created_at) - new Date(a.created_at));
  res.json({ jobs, total: db.jobs.length });
});

app.get('/api/jobs/:id', (req, res) => {
  const job = db.jobs.find(j => j.id === parseInt(req.params.id));
  if (!job) return res.status(404).json({ error: 'Not found' });
  res.json(job);
});

app.post('/api/jobs', (req, res) => {
  const { company, role, url, status, score } = req.body;
  if (!company || !role) return res.status(400).json({ error: 'company and role required' });
  const job = { id: nextId(db.jobs), company, role, url: url || null, status: status || 'discovered', score: score || null, notes: '', created_at: new Date().toISOString() };
  db.jobs.push(job);
  saveDB(db);
  broadcast({ type: 'job_added', data: job });
  res.status(201).json(job);
});

app.patch('/api/jobs/:id', (req, res) => {
  const i = db.jobs.findIndex(j => j.id === parseInt(req.params.id));
  if (i === -1) return res.status(404).json({ error: 'Not found' });
  const { company, role, url, status, score, notes } = req.body;
  if (company !== undefined) db.jobs[i].company = company;
  if (role !== undefined) db.jobs[i].role = role;
  if (url !== undefined) db.jobs[i].url = url;
  if (status !== undefined) db.jobs[i].status = status;
  if (score !== undefined) db.jobs[i].score = score;
  if (notes !== undefined) db.jobs[i].notes = notes;
  saveDB(db);
  broadcast({ type: 'job_updated', data: db.jobs[i] });
  res.json(db.jobs[i]);
});

app.delete('/api/jobs/:id', (req, res) => {
  const id = parseInt(req.params.id);
  db.jobs = db.jobs.filter(j => j.id !== id);
  saveDB(db);
  broadcast({ type: 'job_deleted', data: { id } });
  res.json({ ok: true });
});

app.get('/api/actions', (req, res) => {
  res.json(db.actions.sort((a, b) => new Date(b.created_at) - new Date(a.created_at)).slice(0, 50));
});

app.post('/api/actions', (req, res) => {
  const { type, message, job_id } = req.body;
  if (!type || !message) return res.status(400).json({ error: 'type and message required' });
  const action = { id: nextId(db.actions), type, message, job_id: job_id || null, created_at: new Date().toISOString() };
  db.actions.push(action);
  saveDB(db);
  broadcast({ type: 'action_added', data: action });
  res.status(201).json(action);
});

app.get('/api/agent', (req, res) => res.json(db.agent));

app.post('/api/agent/mode', (req, res) => {
  const { mode } = req.body;
  if (!['review-each', 'routine-auto'].includes(mode)) return res.status(400).json({ error: 'Invalid mode' });
  db.agent.mode = mode;
  db.agent.updated_at = new Date().toISOString();
  saveDB(db);
  broadcast({ type: 'agent_updated', data: db.agent });
  res.json(db.agent);
});

app.post('/api/agent/toggle', (req, res) => {
  db.agent.is_running = db.agent.is_running ? 0 : 1;
  db.agent.updated_at = new Date().toISOString();
  saveDB(db);
  broadcast({ type: 'agent_updated', data: db.agent });
  res.json(db.agent);
});

app.get('/api/runs', (req, res) => res.json(db.runs));

app.get('/api/stats', (req, res) => res.json(computeStats()));

app.get('*', (req, res) => {
  res.sendFile(path.join(__dirname, '..', 'frontend', 'dist', 'index.html'));
});

server.listen(PORT, () => {
  console.log(`Jobay API on http://localhost:${PORT}`);
  console.log(`Dashboard: http://localhost:${PORT}`);
});

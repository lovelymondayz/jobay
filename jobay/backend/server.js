const express = require('express');
const cors = require('cors');
const { WebSocketServer } = require('ws');
const http = require('http');
const path = require('path');
const Database = require('better-sqlite3');

const app = express();
const server = http.createServer(app);
const wss = new WebSocketServer({ server, path: '/ws' });

const PORT = process.env.PORT || 3001;
const DB_PATH = process.env.DB_PATH || path.join(__dirname, 'jobay.db');

app.use(cors());
app.use(express.json());
app.use(express.static(path.join(__dirname, '..', 'frontend', 'dist')));

// ── Database ──────────────────────────────────────────────
const db = new Database(DB_PATH);
db.pragma('journal_mode = WAL');

db.exec(`
  CREATE TABLE IF NOT EXISTS jobs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    company TEXT NOT NULL,
    role TEXT NOT NULL,
    url TEXT,
    status TEXT DEFAULT 'discovered',
    score REAL,
    applied_at TEXT,
    outcome TEXT,
    notes TEXT,
    created_at TEXT DEFAULT (datetime('now'))
  );

  CREATE TABLE IF NOT EXISTS actions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    type TEXT NOT NULL,
    message TEXT,
    job_id INTEGER,
    created_at TEXT DEFAULT (datetime('now'))
  );

  CREATE TABLE IF NOT EXISTS agent_status (
    id INTEGER PRIMARY KEY,
    mode TEXT DEFAULT 'review-each',
    ai_provider TEXT DEFAULT '9router',
    is_running INTEGER DEFAULT 0,
    last_action TEXT,
    updated_at TEXT DEFAULT (datetime('now'))
  );

  CREATE TABLE IF NOT EXISTS runs (
    id TEXT PRIMARY KEY,
    started_at TEXT DEFAULT (datetime('now')),
    ended_at TEXT,
    total_found INTEGER DEFAULT 0,
    total_qualified INTEGER DEFAULT 0,
    total_applied INTEGER DEFAULT 0,
    total_skipped INTEGER DEFAULT 0
  );

  INSERT OR IGNORE INTO agent_status (id, mode, ai_provider, is_running) VALUES (1, 'review-each', '9router', 0);
`);

// ── Helpers ───────────────────────────────────────────────
function seedDummyData() {
  const count = db.prepare('SELECT COUNT(*) as c FROM jobs').get().c;
  if (count > 0) return;

  const dummyJobs = [
    { company: 'TechCorp', role: 'Senior Backend Engineer', url: 'https://techcorp.com/jobs/123', status: 'applied', score: 88 },
    { company: 'StartupXYZ', role: 'Full Stack Developer', url: 'https://startupxyz.com/jobs/456', status: 'review', score: 75 },
    { company: 'BigTech Inc', role: 'Staff Engineer', url: 'https://bigtech.com/jobs/789', status: 'discovered', score: 62 },
    { company: 'RemoteCo', role: 'React Developer', url: 'https://remote.co/jobs/101', status: 'qualified', score: 82 },
    { company: 'AI Labs', role: 'ML Engineer', url: 'https://ailabs.com/jobs/202', status: 'outcome_rejected', score: 70 },
    { company: 'DataDriven', role: 'Data Engineer', url: 'https://datadriven.com/jobs/303', status: 'outcome_interview', score: 85 },
    { company: 'CloudFirst', role: 'DevOps Engineer', url: 'https://cloudfirst.com/jobs/404', status: 'discovered', score: 78 },
    { company: 'FinTech Pro', role: 'Backend Engineer', url: 'https://fintechpro.com/jobs/505', status: 'applied', score: 91 },
  ];

  const insert = db.prepare(`
    INSERT INTO jobs (company, role, url, status, score)
    VALUES (@company, @role, @url, @status, @score)
  `);

  for (const job of dummyJobs) {
    insert.run(job);
  }

  const dummyActions = [
    { type: 'discover', message: 'Found 3 new roles from career pages', job_id: null },
    { type: 'score', message: 'TechCorp Senior Backend Engineer scored 88/100', job_id: 1 },
    { type: 'apply', message: 'Application submitted to TechCorp', job_id: 1 },
    { type: 'score', message: 'StartupXYZ Full Stack Developer scored 75/100', job_id: 2 },
    { type: 'review', message: 'Flagged for manual review: StartupXYZ', job_id: 2 },
    { type: 'score', message: 'RemoteCo React Developer scored 82/100', job_id: 4 },
    { type: 'outcome', message: 'Interview scheduled with DataDriven', job_id: 6 },
  ];

  const insertAction = db.prepare(`
    INSERT INTO actions (type, message, job_id)
    VALUES (@type, @message, @job_id)
  `);

  for (const action of dummyActions) {
    insertAction.run(action);
  }

  // Seed a run
  db.prepare(`
    INSERT INTO runs (id, started_at, ended_at, total_found, total_qualified, total_applied, total_skipped)
    VALUES (?, datetime('now', '-2 hours'), datetime('now', '-1 hour'), 8, 5, 2, 1)
  `).run('demo-run-1');

  console.log('✓ Seeded dummy data');
}

seedDummyData();

// ── WebSocket ─────────────────────────────────────────────
const clients = new Set();

wss.on('connection', (ws) => {
  clients.add(ws);
  console.log('WS client connected');
  ws.send(JSON.stringify({ type: 'init', data: getStatus() }));

  ws.on('close', () => clients.delete(ws));
});

function broadcast(data) {
  const msg = JSON.stringify(data);
  for (const client of clients) {
    if (client.readyState === 1) {
      client.send(msg);
    }
  }
}

function getStatus() {
  const jobs = db.prepare('SELECT * FROM jobs ORDER BY created_at DESC').all();
  const actions = db.prepare('SELECT * FROM actions ORDER BY created_at DESC LIMIT 50').all();
  const agent = db.prepare('SELECT * FROM agent_status WHERE id = 1').get();
  const runs = db.prepare('SELECT * FROM runs ORDER BY started_at DESC LIMIT 10').all();
  const stats = {
    discovered: db.prepare("SELECT COUNT(*) as c FROM jobs WHERE status = 'discovered'").get().c,
    qualified: db.prepare("SELECT COUNT(*) as c FROM jobs WHERE status = 'qualified'").get().c,
    review: db.prepare("SELECT COUNT(*) as c FROM jobs WHERE status = 'review'").get().c,
    applied: db.prepare("SELECT COUNT(*) as c FROM jobs WHERE status = 'applied'").get().c,
    outcome_interview: db.prepare("SELECT COUNT(*) as c FROM jobs WHERE status = 'outcome_interview'").get().c,
    outcome_rejected: db.prepare("SELECT COUNT(*) as c FROM jobs WHERE status = 'outcome_rejected'").get().c,
    total: db.prepare("SELECT COUNT(*) as c FROM jobs").get().c,
  };
  return { jobs, actions, agent, runs, stats };
}

// ── REST API ──────────────────────────────────────────────

// Dashboard status
app.get('/api/status', (req, res) => {
  res.json(getStatus());
});

// Jobs
app.get('/api/jobs', (req, res) => {
  const { status, search, page = 1, limit = 20 } = req.query;
  let sql = 'SELECT * FROM jobs';
  const params = [];
  const conditions = [];

  if (status) {
    conditions.push('status = ?');
    params.push(status);
  }
  if (search) {
    conditions.push('(company LIKE ? OR role LIKE ?)');
    params.push(`%${search}%`, `%${search}%`);
  }

  if (conditions.length > 0) {
    sql += ' WHERE ' + conditions.join(' AND ');
  }

  sql += ' ORDER BY created_at DESC';
  sql += ' LIMIT ? OFFSET ?';
  params.push(parseInt(limit), (parseInt(page) - 1) * parseInt(limit));

  const jobs = db.prepare(sql).all(...params);
  const total = db.prepare('SELECT COUNT(*) as c FROM jobs').get().c;
  res.json({ jobs, total, page: parseInt(page), limit: parseInt(limit) });
});

app.get('/api/jobs/:id', (req, res) => {
  const job = db.prepare('SELECT * FROM jobs WHERE id = ?').get(req.params.id);
  if (!job) return res.status(404).json({ error: 'Job not found' });
  res.json(job);
});

app.post('/api/jobs', (req, res) => {
  const { company, role, url, status, score } = req.body;
  if (!company || !role) return res.status(400).json({ error: 'company and role required' });

  const result = db.prepare(`
    INSERT INTO jobs (company, role, url, status, score)
    VALUES (?, ?, ?, ?, ?)
  `).run(company, role, url || null, status || 'discovered', score || null);

  const job = db.prepare('SELECT * FROM jobs WHERE id = ?').get(result.lastInsertRowid);
  broadcast({ type: 'job_added', data: job });
  res.status(201).json(job);
});

app.patch('/api/jobs/:id', (req, res) => {
  const existing = db.prepare('SELECT * FROM jobs WHERE id = ?').get(req.params.id);
  if (!existing) return res.status(404).json({ error: 'Job not found' });

  const { company, role, url, status, score, outcome, notes } = req.body;
  db.prepare(`
    UPDATE jobs SET
      company = COALESCE(?, company),
      role = COALESCE(?, role),
      url = COALESCE(?, url),
      status = COALESCE(?, status),
      score = COALESCE(?, score),
      outcome = COALESCE(?, outcome),
      notes = COALESCE(?, notes)
    WHERE id = ?
  `).run(company ?? null, role ?? null, url ?? null, status ?? null, score ?? null, outcome ?? null, notes ?? null, req.params.id);

  const job = db.prepare('SELECT * FROM jobs WHERE id = ?').get(req.params.id);
  broadcast({ type: 'job_updated', data: job });
  res.json(job);
});

app.delete('/api/jobs/:id', (req, res) => {
  db.prepare('DELETE FROM jobs WHERE id = ?').run(req.params.id);
  broadcast({ type: 'job_deleted', data: { id: parseInt(req.params.id) } });
  res.json({ ok: true });
});

// Actions log
app.get('/api/actions', (req, res) => {
  const limit = parseInt(req.query.limit) || 50;
  const actions = db.prepare('SELECT * FROM actions ORDER BY created_at DESC LIMIT ?').all(limit);
  res.json(actions);
});

app.post('/api/actions', (req, res) => {
  const { type, message, job_id } = req.body;
  if (!type || !message) return res.status(400).json({ error: 'type and message required' });

  const result = db.prepare(`
    INSERT INTO actions (type, message, job_id)
    VALUES (?, ?, ?)
  `).run(type, message, job_id || null);

  const action = db.prepare('SELECT * FROM actions WHERE id = ?').get(result.lastInsertRowid);
  broadcast({ type: 'action_added', data: action });
  res.status(201).json(action);
});

// Agent status
app.get('/api/agent', (req, res) => {
  const agent = db.prepare('SELECT * FROM agent_status WHERE id = 1').get();
  res.json(agent);
});

app.post('/api/agent/mode', (req, res) => {
  const { mode } = req.body;
  if (!['review-each', 'routine-auto'].includes(mode)) {
    return res.status(400).json({ error: 'mode must be review-each or routine-auto' });
  }
  db.prepare('UPDATE agent_status SET mode = ?, updated_at = datetime(\'now\') WHERE id = 1').run(mode);
  const agent = db.prepare('SELECT * FROM agent_status WHERE id = 1').get();
  broadcast({ type: 'agent_updated', data: agent });
  res.json(agent);
});

app.post('/api/agent/toggle', (req, res) => {
  const agent = db.prepare('SELECT * FROM agent_status WHERE id = 1').get();
  const newRunning = agent.is_running ? 0 : 1;
  db.prepare('UPDATE agent_status SET is_running = ?, updated_at = datetime(\'now\') WHERE id = 1').run(newRunning);
  const updated = db.prepare('SELECT * FROM agent_status WHERE id = 1').get();
  broadcast({ type: 'agent_updated', data: updated });
  res.json(updated);
});

// Runs
app.get('/api/runs', (req, res) => {
  const runs = db.prepare('SELECT * FROM runs ORDER BY started_at DESC LIMIT 20').all();
  res.json(runs);
});

app.post('/api/runs', (req, res) => {
  const { id, total_found, total_qualified, total_applied, total_skipped } = req.body;
  db.prepare(`
    INSERT INTO runs (id, started_at, ended_at, total_found, total_qualified, total_applied, total_skipped)
    VALUES (?, datetime('now'), NULL, ?, ?, ?, ?)
  `).run(id || `run-${Date.now()}`, total_found || 0, total_qualified || 0, total_applied || 0, total_skipped || 0);
  const run = db.prepare('SELECT * FROM runs WHERE id = ?').get(id || `run-${Date.now()}`);
  broadcast({ type: 'run_started', data: run });
  res.status(201).json(run);
});

app.patch('/api/runs/:id', (req, res) => {
  const { total_found, total_qualified, total_applied, total_skipped } = req.body;
  db.prepare(`
    UPDATE runs SET
      total_found = COALESCE(?, total_found),
      total_qualified = COALESCE(?, total_qualified),
      total_applied = COALESCE(?, total_applied),
      total_skipped = COALESCE(?, total_skipped),
      ended_at = COALESCE(?, ended_at)
    WHERE id = ?
  `).run(total_found ?? null, total_qualified ?? null, total_applied ?? null, total_skipped ?? null, new Date().toISOString(), req.params.id);
  const run = db.prepare('SELECT * FROM runs WHERE id = ?').get(req.params.id);
  broadcast({ type: 'run_updated', data: run });
  res.json(run);
});

// Stats
app.get('/api/stats', (req, res) => {
  const stats = {
    discovered: db.prepare("SELECT COUNT(*) as c FROM jobs WHERE status = 'discovered'").get().c,
    qualified: db.prepare("SELECT COUNT(*) as c FROM jobs WHERE status = 'qualified'").get().c,
    review: db.prepare("SELECT COUNT(*) as c FROM jobs WHERE status = 'review'").get().c,
    applied: db.prepare("SELECT COUNT(*) as c FROM jobs WHERE status = 'applied'").get().c,
    outcome_interview: db.prepare("SELECT COUNT(*) as c FROM jobs WHERE status = 'outcome_interview'").get().c,
    outcome_rejected: db.prepare("SELECT COUNT(*) as c FROM jobs WHERE status = 'outcome_rejected'").get().c,
    total: db.prepare("SELECT COUNT(*) as c FROM jobs").get().c,
  };
  res.json(stats);
});

// SPA fallback
app.get('*', (req, res) => {
  res.sendFile(path.join(__dirname, '..', 'frontend', 'dist', 'index.html'));
});

// ── Start ─────────────────────────────────────────────────
server.listen(PORT, () => {
  console.log(`Jobay API running on http://localhost:${PORT}`);
  console.log(`Dashboard: http://localhost:${PORT}`);
  console.log(`WebSocket: ws://localhost:${PORT}/ws`);
});

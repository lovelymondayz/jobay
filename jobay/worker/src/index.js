const express = require('express');
const { chromium } = require('playwright');
const path = require('path');
const fs = require('fs');

const app = express();
app.use(express.json());

const PORT = process.env.WORKER_PORT || 3011;
const SCREENSHOT_DIR = process.env.SCREENSHOT_DIR || '/app/data/screenshots';

// Ensure screenshot dir exists
fs.mkdirSync(SCREENSHOT_DIR, { recursive: true });

// POST /apply — apply to a job
app.post('/apply', async (req, res) => {
  const { jobUrl, userProfile } = req.body;

  if (!jobUrl || !userProfile) {
    return res.status(400).json({ ok: false, error: 'jobUrl and userProfile required' });
  }

  let browser;
  try {
    browser = await chromium.launch({
      headless: true,
      args: ['--no-sandbox', '--disable-setuid-sandbox'],
    });

    const context = await browser.newContext({
      userAgent: 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36',
    });

    const page = await context.newPage();
    await page.goto(jobUrl, { waitUntil: 'networkidle', timeout: 30000 });

    // Detect and fill form fields
    const fillResult = await fillForm(page, userProfile);

    // Click submit if found
    const submitResult = await clickSubmit(page);

    // Take screenshot
    const screenshotPath = path.join(SCREENSHOT_DIR, `${Date.now()}.png`);
    await page.screenshot({ path: screenshotPath, fullPage: true });

    await browser.close();

    res.json({
      ok: true,
      filled: fillResult.filled,
      submitted: submitResult.submitted,
      screenshot: screenshotPath,
      details: fillResult.details,
    });
  } catch (err) {
    if (browser) await browser.close();
    res.json({ ok: false, error: err.message });
  }
});

async function fillForm(page, profile) {
  const filled = [];
  const details = [];

  // Common name field selectors
  const nameSelectors = [
    'input[name="name"]', 'input[name="full_name"]', 'input[name="fullname"]',
    'input[id="name"]', 'input[placeholder*="name" i]', 'input[placeholder*="nama" i]',
    '[data-testid="name"]', 'input[type="text"]',
  ];

  for (const sel of nameSelectors) {
    const el = await page.$(sel);
    if (el) {
      await el.fill(profile.name || '');
      filled.push('name');
      details.push(`Filled name: ${profile.name}`);
      break;
    }
  }

  // Email
  const emailSelectors = [
    'input[name="email"]', 'input[type="email"]', 'input[id="email"]',
    'input[placeholder*="email" i]',
  ];
  for (const sel of emailSelectors) {
    const el = await page.$(sel);
    if (el) {
      await el.fill(profile.email || '');
      filled.push('email');
      details.push(`Filled email: ${profile.email}`);
      break;
    }
  }

  // Phone
  const phoneSelectors = [
    'input[name="phone"]', 'input[name="mobile"]', 'input[type="tel"]',
    'input[placeholder*="phone" i]', 'input[placeholder*="telepon" i]',
  ];
  for (const sel of phoneSelectors) {
    const el = await page.$(sel);
    if (el) {
      await el.fill(profile.phone || '');
      filled.push('phone');
      details.push(`Filled phone: ${profile.phone}`);
      break;
    }
  }

  // Cover letter / message
  const coverSelectors = [
    'textarea[name="cover_letter"]', 'textarea[name="message"]', 'textarea[name="about"]',
    'textarea[id="cover"]', 'textarea[placeholder*="cover" i]', 'textarea[placeholder*="tentang" i]',
  ];
  for (const sel of coverSelectors) {
    const el = await page.$(sel);
    if (el) {
      await el.fill(profile.summary || profile.coverLetter || '');
      filled.push('cover_letter');
      details.push('Filled cover letter');
      break;
    }
  }

  // Upload CV if file input exists
  const fileSelectors = [
    'input[type="file"]', 'input[name="cv"]', 'input[name="resume"]',
    'input[id="cv"]', 'input[id="resume"]',
  ];
  for (const sel of fileSelectors) {
    const el = await page.$(sel);
    if (el && profile.cvPath) {
      await el.setInputFiles(profile.cvPath);
      filled.push('cv');
      details.push(`Uploaded CV: ${profile.cvPath}`);
      break;
    }
  }

  return { filled, details };
}

async function clickSubmit(page) {
  const submitSelectors = [
    'button[type="submit"]', 'input[type="submit"]',
    'button:has-text("Apply")', 'button:has-text("Submit")',
    'button:has-text("Kirim")', 'button:has-text("Lamar")',
    '[data-testid="submit"]', '[data-testid="apply"]',
  ];

  for (const sel of submitSelectors) {
    const el = await page.$(sel);
    if (el) {
      try {
        await el.click({ timeout: 5000 });
        return { submitted: true, selector: sel };
      } catch {
        continue;
      }
    }
  }

  return { submitted: false, selector: null };
}

// Health check
app.get('/health', (req, res) => {
  res.json({ status: 'ok' });
});

app.listen(PORT, '0.0.0.0', () => {
  console.log(`Jobay worker listening on port ${PORT}`);
});

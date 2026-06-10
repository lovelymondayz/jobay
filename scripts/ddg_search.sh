#!/bin/bash
# ddg_search.sh — DuckDuckGo Lite search (no API key needed)
# Usage: ./ddg_search.sh "search query" [num_results]

QUERY="${1:?Usage: $0 \"search query\" [num_results]}"
NUM="${2:-10}"

ENCODED=$(python3 -c "import urllib.parse; print(urllib.parse.quote_plus('$QUERY'))")

curl -s --max-time 15 "https://lite.duckduckgo.com/lite/?q=${ENCODED}" \
  -H "User-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36" \
  2>/dev/null | python3 -c "
import re, sys
from urllib.parse import unquote

html = sys.stdin.read()
num = int(sys.argv[1]) if len(sys.argv) > 1 else 10

# DDG Lite: results are in <a href=\"/l/?uddg=ENCODED_URL...\">Title</a>
# followed by <td> with snippet text
results = re.findall(
    r'<a[^>]+href=\"(/l/\?uddg=[^\"]+)\"[^>]*>([^<]+)</a>.*?<td[^>]*class=\"result-snippet\"[^>]*>(.*?)</td>',
    html, re.DOTALL
)

count = 0
for link, title, snippet in results:
    if count >= num:
        break
    title = title.strip()
    snippet = re.sub(r'<[^>]+>', '', snippet).strip()
    # Decode the real URL from uddg=
    real_url = re.search(r'uddg=([^&]+)', link)
    url = unquote(real_url.group(1)) if real_url else ''
    # Skip ads
    if not title or not url or url.startswith('https://duckduckgo.com/y.js'):
        continue
    import html as htmlmod
    title = htmlmod.unescape(title)
    snippet = htmlmod.unescape(snippet)
    count += 1
    print(f'--- Result {count} ---')
    print(f'Title: {title}')
    print(f'URL: {url}')
    print(f'Snippet: {snippet}')
    print()
" "${NUM}"

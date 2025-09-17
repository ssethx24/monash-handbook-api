import re
import json
import time
import os
from urllib.parse import urlparse
import requests
from lxml import etree
from tqdm import tqdm

HAND_BOOK_BASE = "https://handbook.monash.edu"
YEAR = "2025"  # change to "current" or a specific year (2020–2025)
LOCAL_API = f"http://localhost:8080/v1/{YEAR}/units"

SITEMAP_INDEX = f"{HAND_BOOK_BASE}/sitemap.xml"
OUT_PATH = f"units_{YEAR}.ndjson"

def fetch(url, **kwargs):
    headers = {
        "User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X) AppleWebKit/537.36 "
                      "(KHTML, like Gecko) Chrome Safari"
    }
    return requests.get(url, headers=headers, timeout=30, **kwargs)

def parse_xml(content):
    return etree.fromstring(content.encode("utf-8") if isinstance(content, str) else content)

def get_all_sitemaps():
    # Get top-level sitemap index, then pull nested sitemaps (if any)
    r = fetch(SITEMAP_INDEX)
    r.raise_for_status()
    root = parse_xml(r.content)
    ns = {"sm": "http://www.sitemaps.org/schemas/sitemap/0.9"}
    sitemaps = [el.text for el in root.findall(".//sm:sitemap/sm:loc", ns)]
    if root.tag.endswith("urlset"):
        return [SITEMAP_INDEX] + sitemaps
    return sitemaps

def get_urls_from_sitemap(smap_url):
    r = fetch(smap_url)
    r.raise_for_status()
    root = parse_xml(r.content)
    ns = {"sm": "http://www.sitemaps.org/schemas/sitemap/0.9"}
    if root.tag.endswith("sitemapindex"):
        return [el.text for el in root.findall(".//sm:sitemap/sm:loc", ns)]
    return [el.text for el in root.findall(".//sm:url/sm:loc", ns)]

def is_unit_url(url):
    try:
        path = urlparse(url).path.strip("/")
        parts = path.split("/")
        return (
            len(parts) == 3 and
            parts[0] == YEAR and
            parts[1].lower() == "units" and
            re.match(r"^[A-Za-z]{2,4}\d{4}$", parts[2]) is not None
        )
    except Exception:
        return False

def extract_code(url):
    return url.rstrip("/").split("/")[-1]

def load_existing():
    """Load codes already fetched into a dict {code: data}"""
    done = {}
    if os.path.exists(OUT_PATH):
        with open(OUT_PATH, "r", encoding="utf-8") as f:
            for line in f:
                try:
                    data = json.loads(line)
                    code = data.get("common", {}).get("code") or data.get("code")
                    if code:
                        done[code.upper()] = data
                except json.JSONDecodeError:
                    continue
    return done

def main():
    print(f"Fetching sitemap index: {SITEMAP_INDEX}")
    to_visit = get_all_sitemaps()
    all_urls = []

    # Expand nested sitemap indexes
    stack = list(to_visit)
    visited = set()
    while stack:
        sm = stack.pop()
        if sm in visited:
            continue
        visited.add(sm)
        try:
            urls = get_urls_from_sitemap(sm)
        except Exception as e:
            print(f"Warning: failed to fetch {sm}: {e}")
            continue
        for u in urls:
            if u.endswith(".xml"):
                stack.append(u)
            else:
                all_urls.append(u)

    unit_urls = [u for u in all_urls if is_unit_url(u)]
    unit_codes = sorted({extract_code(u).upper() for u in unit_urls})

    print(f"Found {len(unit_codes)} unit codes for YEAR={YEAR}.")

    existing = load_existing()
    print(f"Already have {len(existing)} units in {OUT_PATH}")

    # Append mode so we don’t overwrite existing file
    with open(OUT_PATH, "a", encoding="utf-8") as f:
        for code in tqdm(unit_codes, desc="Fetching units from local API"):
            if code in existing and "error_status" not in existing[code]:
                continue  # already done successfully

            url = f"{LOCAL_API}/{code}"
            delay = 0.2
            for attempt in range(5):  # retry up to 5 times
                try:
                    resp = fetch(url)
                    if resp.status_code == 200:
                        f.write(json.dumps(resp.json(), ensure_ascii=False) + "\n")
                        f.flush()
                        break
                    elif resp.status_code == 403:
                        print(f"403 Forbidden for {code}, backing off {delay:.1f}s...")
                        time.sleep(delay)
                        delay *= 2
                    else:
                        f.write(json.dumps({
                            "code": code,
                            "error_status": resp.status_code,
                            "error_text": resp.text
                        }, ensure_ascii=False) + "\n")
                        f.flush()
                        break
                except Exception as e:
                    f.write(json.dumps({"code": code, "exception": str(e)}) + "\n")
                    f.flush()
                    break
            time.sleep(0.1)  # normal pacing

    print(f"Done. Updated {OUT_PATH}")

if __name__ == "__main__":
    main()


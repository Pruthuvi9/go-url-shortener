
# Functional requirements
- Submit a long URL => get short URL
    - Optional custom alias
    - Optional expiration time
- Go to short URL => redirect to long URL
- Delete a short URL
- Auth required to create/delete URLs (V2)

# NFR
- Shortcode uniqueness
- Minimal redirection delay (<100ms)
- 4 9s availability (99.99% uptime) and reliability (availability > consistency)
- Support 1B URLs and 100M DAU

# Important considerations
- Read and write requirements of this system are significantly different
  (many more visits to short URLs compared to creation of short URLs from long URLs)
- Hash-based generation means duplicate long URLs resolve to the same short code — deduplication is free

# The Core Entities
- Original URL
- Short URL
- User (the user who created the URL) — auth is a V2 feature

# API

// Shorten a URL
POST /urls
{
  "long_url": "https://www.example.com/some/very/long/url",
  "custom_alias": "optional_custom_alias",
  "expiration_date": "optional_expiration_date"
}
->
{
  "short_url": "http://short.ly/abc123"
}

// Redirect to long URL
GET /{short_code}
-> 302 Found  (Location: <long_url>)
-> 404 Not Found
-> 410 Gone   (URL existed but has expired)

// Delete a short URL
DELETE /urls/{short_code}
-> 204 No Content
-> 404 Not Found

# Shortcode Generation

Strategy: SHA-256 hash of the long URL, Base62-encoded, first 7 characters taken.
- 7 Base62 characters = 62^7 ≈ 3.5 trillion possible codes — sufficient for 1B URLs
- Same long URL always produces the same short code (built-in deduplication)

Collision handling:
- On write, check if the 7-char code already maps to a DIFFERENT long URL
- If collision: append an incrementing counter to the input before re-hashing (url+"1", url+"2", ...)
  and retry until a free slot is found (in practice, collisions are extremely rare)

Custom alias:
- Allowed characters: [a-zA-Z0-9_-]
- Length: 3–20 characters
- Stored as-is (case-sensitive)
- A blocklist of reserved paths prevents aliasing system routes
  (e.g. "api", "urls", "admin", "health")
- Conflict check against existing codes before saving

# URL Validation

Input long URL is validated with a regex before processing:
  ^https?://[^\s/$.?#].[^\s]*$

This checks scheme (http/https) and basic structure without making a network request,
keeping creation latency low.

# Database

**Suggestion:** PostgreSQL (single instance to start, read replicas when read load grows).

Schema (urls table):
  short_code   VARCHAR(20)   PRIMARY KEY
  long_url     TEXT          NOT NULL
  created_at   TIMESTAMPTZ   NOT NULL DEFAULT now()
  expires_at   TIMESTAMPTZ   NULL

- `short_code` as primary key gives O(log n) lookup and natural uniqueness enforcement.
- No user_id column until V2.

Sharding path (when single Postgres saturates writes):
- Shard by the first character of short_code (Base62 = 62 shards, or group into N ranges).
- Consistent hashing on short_code is an alternative to avoid rebalancing on shard count changes.
- At 1B rows × ~200 bytes avg row size ≈ 200 GB — fits a single large Postgres instance for a long time
  before sharding is actually necessary.

# Caching (Valkey)

All redirect lookups check Valkey before hitting Postgres.

Cache entry: short_code → long_url (string key-value)

TTL policy:
- If the URL has an expiration_date: TTL = expiration_date - now
- Otherwise: TTL = 24 hours (configurable)

On create:  write-through — populate Valkey immediately after the DB write
On delete:  DEL the key from Valkey
On miss:    fetch from Postgres, populate Valkey, serve response
On expired: Valkey entry has already TTL'd out; DB lookup returns the row,
            expires_at is checked, 410 is returned (entry is NOT re-cached)

# Main Flows

1. Create short URL

    Start

    Validate long_url with regex.
        Invalid => 400 Bad Request

    Is a custom_alias provided?
        Yes:
            Validate alias format ([a-zA-Z0-9_-], 3–20 chars).
                Invalid => 400 Bad Request
            Check alias against reserved blocklist.
                Reserved => 409 Conflict
            Check alias uniqueness in DB.
                Already exists => 409 Conflict
            short_code = alias
        No:
            Compute SHA-256(long_url), Base62-encode, take first 7 chars => short_code
            Check if short_code exists in DB AND maps to a different long_url (collision).
                Collision: retry with SHA-256(long_url + counter), increment counter, repeat.

    Save (short_code, long_url, expires_at) to Postgres.

    Write short_code => long_url to Valkey (with TTL).

    Return 201 Created { "short_url": "http://short.ly/<short_code>" }

    End


2. Redirect (GET /{short_code})

    Start

    Look up short_code in Valkey.
        Hit:
            Return 302 Found, Location: long_url
            End

        Miss:
            Look up short_code in Postgres.
                Not found => 404 Not Found, End

                Found:
                    Is expires_at set AND expires_at < now?
                        Yes => 410 Gone, End

                    Populate Valkey (short_code => long_url, TTL).
                    Return 302 Found, Location: long_url
                    End


3. Delete short URL (DELETE /urls/{short_code})

    Start

    Look up short_code in Postgres.
        Not found => 404 Not Found, End

    Delete row from Postgres.

    DEL short_code from Valkey.

    Return 204 No Content

    End

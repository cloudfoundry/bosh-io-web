# bosh.io API

bosh.io exposes a simple read-only JSON API for querying stemcells and releases.

**Base URL**: `https://bosh.io`

---

## Stemcells

### List stemcell versions

Returns a list of versions for a given stemcell name, sorted from latest to oldest.

```
GET /api/v1/stemcells/{name}
```

**Path parameters**

| Parameter | Description |
|-----------|-------------|
| `name` | Stemcell manifest name (e.g. `bosh-aws-xen-hvm-ubuntu-jammy-go_agent`) |

**Query parameters**

| Parameter | Description |
|-----------|-------------|
| `all` | Set to `1` or `true` to include all versions. By default only the latest unique versions are returned. |

**Response**

`200 OK` — JSON array of stemcell objects. The first element is the latest version. Returns an empty array `[]` when no stemcells are found.

```json
[
  {
    "name": "bosh-aws-xen-hvm-ubuntu-jammy-go_agent",
    "version": "1.234",
    "regular": {
      "url": "https://...",
      "size": 123456789,
      "md5": "abc123",
      "sha1": "abc123",
      "sha256": "abc123"
    },
    "light": {
      "url": "https://...",
      "size": 123456789,
      "md5": "abc123",
      "sha1": "abc123",
      "sha256": "abc123"
    }
  }
]
```

**Stemcell object fields**

| Field | Type | Description |
|-------|------|-------------|
| `name` | string | Stemcell manifest name |
| `version` | string | Stemcell version |
| `regular` | object \| null | Full (non-light) stemcell source. Present when available. |
| `light` | object \| null | Light stemcell source. Present when available. |
| `light_china` | object \| null | Light stemcell for China region. Present when available. |

**Stemcell source object fields** (`regular`, `light`, `light_china`)

| Field | Type | Description |
|-------|------|-------------|
| `url` | string | Direct download URL |
| `size` | number | File size in bytes |
| `md5` | string | MD5 checksum |
| `sha1` | string | SHA1 checksum (omitted if empty) |
| `sha256` | string | SHA256 checksum (omitted if empty) |

**Error responses**

| Status | Description |
|--------|-------------|
| `400` | `name` parameter is missing |
| `500` | Internal server error |

**Example**

```bash
# Get all versions of a stemcell
curl https://bosh.io/api/v1/stemcells/bosh-aws-xen-hvm-ubuntu-jammy-go_agent

# Get latest version only
curl "https://bosh.io/api/v1/stemcells/bosh-aws-xen-hvm-ubuntu-jammy-go_agent" | jq '.[0]'
```

---

### Download a stemcell

Redirects to the actual stemcell tarball download URL. When no version is specified, redirects to the latest version.

```
GET /d/stemcells/{name}
```

**Path parameters**

| Parameter | Description |
|-----------|-------------|
| `name` | Stemcell manifest name (e.g. `bosh-aws-xen-hvm-ubuntu-jammy-go_agent`) |

**Query parameters**

| Parameter | Default | Description |
|-----------|---------|-------------|
| `v` | latest | Stemcell version to download (e.g. `1.234`). Omit to get the latest. |
| `light` | `true` | Set to `1`/`true` to prefer the light stemcell. Set to `0`/`false` to force the full stemcell. |
| `china` | `false` | Set to `1`/`true` to get the light stemcell for the China region. |

**Response**

`302 Found` — Redirects to the direct download URL of the stemcell tarball.

**Error responses**

| Status | Description |
|--------|-------------|
| `400` | Invalid request parameters (for example, an empty stemcell name or an invalid `v` version string) |
| `404` | Stemcell or requested version not found |
| `500` | Internal server error |

**Examples**

```bash
# Download the latest stemcell (follows redirect)
curl -L https://bosh.io/d/stemcells/bosh-aws-xen-hvm-ubuntu-jammy-go_agent

# Download a specific version
curl -L "https://bosh.io/d/stemcells/bosh-aws-xen-hvm-ubuntu-jammy-go_agent?v=1.234"

# Download the full (non-light) stemcell
curl -L "https://bosh.io/d/stemcells/bosh-aws-xen-hvm-ubuntu-jammy-go_agent?light=false"

# Get the redirect URL without following it
curl -I https://bosh.io/d/stemcells/bosh-aws-xen-hvm-ubuntu-jammy-go_agent
```

---

## Releases

### List release versions

Returns a list of versions for a given release source, sorted from latest to oldest.

```
GET /api/v1/releases/{source}
```

**Path parameters**

| Parameter | Description |
|-----------|-------------|
| `source` | Release source path, typically a GitHub repository path (e.g. `github.com/cloudfoundry/cf-release`) |

**Response**

`200 OK` — JSON array of release objects. The first element is the latest version. Returns an empty array `[]` when no releases are found.

```json
[
  {
    "name": "github.com/cloudfoundry/cf-release",
    "version": "287",
    "url": "https://bosh.io/d/github.com/cloudfoundry/cf-release?v=287",
    "sha1": "abc123",
    "sha256": "abc123"
  }
]
```

**Release object fields**

| Field | Type | Description |
|-------|------|-------------|
| `name` | string | Release source (GitHub repository path) |
| `version` | string | Release version |
| `url` | string | bosh.io download URL for this version |
| `sha1` | string | SHA1 checksum of the release tarball |
| `sha256` | string | SHA256 checksum of the release tarball (omitted if empty) |

**Error responses**

| Status | Description |
|--------|-------------|
| `400` | `source` parameter is missing |
| `500` | Internal server error. On repository failures, this endpoint renders an HTML error page (`text/html`) rather than a JSON error body, unlike the stemcell API. |

**Example**

```bash
# Get all versions of a release
curl https://bosh.io/api/v1/releases/github.com/cloudfoundry/cf-release

# Get the latest version
curl https://bosh.io/api/v1/releases/github.com/cloudfoundry/cf-release | jq '.[0]'
```

---

### Download a release tarball

Redirects to the actual release tarball download URL. When no version is specified, redirects to the latest version.

```
GET /d/{source}
```

**Path parameters**

| Parameter | Description |
|-----------|-------------|
| `source` | Release source path (e.g. `github.com/cloudfoundry/cf-release`) |

**Query parameters**

| Parameter | Default | Description |
|-----------|---------|-------------|
| `v` | latest | Release version to download (e.g. `287`). Omit to get the latest. |

**Response**

`302 Found` — Redirects to the direct download URL of the release tarball.

> Note: Unlike the JSON API endpoints above, this download endpoint returns an HTML error page on failure rather than a JSON error body.

**Error responses**

| Status | Description |
|--------|-------------|
| `400` | Bad request; the release source path is empty |
| `500` | Release or version not found, or tarball unavailable |

**Examples**

```bash
# Download the latest release tarball (follows redirect)
curl -L https://bosh.io/d/github.com/cloudfoundry/cf-release

# Download a specific version
curl -L "https://bosh.io/d/github.com/cloudfoundry/cf-release?v=287"

# Get the redirect URL without following it
curl -I https://bosh.io/d/github.com/cloudfoundry/cf-release
```

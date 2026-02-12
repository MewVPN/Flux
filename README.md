# Flux

Flux is a lightweight edge agent for the Ampere control plane.

It runs on managed servers, controls WireGuard via wg-easy, performs external connectivity checks, and exposes a secure HTTP API for remote orchestration.

---

## Features

- Managed **wg-easy (15.2.2)**
- WireGuard peer management:
    - create
    - enable / disable
    - download client config
- Secure bootstrap with shared secret
- Protected API (X-Agent-Secret)
- Cached external service health checks (hourly background checks)
- Multi-platform static builds

---

## Architecture

Ampere (control plane) → Flux (edge agent) → wg-easy (WireGuard)

Flux runs on each edge server and exposes a secure HTTP API.  
Ampere communicates with Flux`s to manage WireGuard peers and monitor connectivity.
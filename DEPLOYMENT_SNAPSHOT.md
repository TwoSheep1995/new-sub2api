# Sub2API Deployment Snapshot

This repository is a source snapshot of the running `/opt/sub2api/app`
deployment on `2026-05-25`.

## Current Runtime Shape

- Application container: `sub2api`
- Image tag: `sub2api-local:model-square-group-rows-20260525`
- Entrypoint: `/app/docker-entrypoint.sh`
- Command: `/app/sub2api`
- Host port mapping: `${BIND_HOST:-0.0.0.0}:${SERVER_PORT:-8080}:8080`
- Database container: `sub2api-postgres` using `postgres:18-alpine`
- Redis container: `sub2api-redis` using `redis:8-alpine`
- External Docker networks expected by the current compose file:
  - `sub2api-ingress-network`
  - `sub2api-vpn-network`
  - `sub2api-cpa-network`

The compose file includes `build: .` and the current image tag, so a new
server can build the same local image name from source with
`docker compose up -d --build`.

## Not Included In Git

The following runtime state is intentionally not committed:

- `.env`
- `data/`
- `postgres_data/`
- `redis_data/`
- `backups/`
- local logs, caches, generated frontend build output, and database dumps

Restore PostgreSQL/Redis/database state from backups separately. Fill `.env`
from `.env.example` with the production values for the target server.

## New Server Bootstrap

```sh
git clone https://github.com/TwoSheep1995/new-sub2api.git
cd new-sub2api
cp .env.example .env
# edit .env with production values and restored secret values
docker network create sub2api-ingress-network || true
docker network create sub2api-vpn-network || true
docker network create sub2api-cpa-network || true
docker compose up -d --build
```

If the target server has different ingress/proxy/VPN topology, update the
external network section in `docker-compose.yml` before starting the service.

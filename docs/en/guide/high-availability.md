# High Availability

gocron supports multi-instance deployment with automatic leader election. Only the leader node runs the scheduler; followers stay on hot standby and take over within seconds if the leader goes down.

## How It Works

gocron uses database row locking (`SELECT ... FOR UPDATE`) for leader election, inspired by the proven approach used in [XXL-JOB](https://github.com/xuxueli/xxl-job).

- A `scheduler_lock` table is created automatically on first startup
- One instance acquires the lock and becomes the **Leader**, running all scheduled tasks
- Other instances are **Followers** — they serve the Web UI and API but do not execute tasks
- The Leader renews its lease every **5 seconds** (lease expires after **15 seconds**)
- If the Leader shuts down gracefully, the lock is released immediately — failover is instant
- If the Leader crashes, the lease expires and a Follower takes over within **15 seconds**

::: warning Database Requirement
HA mode requires **MySQL** or **PostgreSQL**. SQLite does not support concurrent multi-process access, so leader election is automatically skipped in SQLite mode (single-node).
:::

## Setup

### 1. Install the First Node

Start the first node and complete the web installation wizard. Choose **MySQL** or **PostgreSQL** as the database engine.

```bash
./gocron web --port 5920
# Open http://localhost:5920 and complete installation
```

### 2. Copy Configuration to Other Nodes

Copy the `.gocron/conf/` directory from the first node to each additional node:

```bash
scp -r .gocron/conf/ user@node2:/path/to/gocron/.gocron/conf/
```

The directory contains:
- `app.ini` — database and application settings
- `install.lock` — marks the installation as complete
- `.version` — current application version

### 3. Start All Nodes

```bash
# Node 1
./gocron web --port 5920

# Node 2
./gocron web --port 5921
```

Both nodes connect to the same database. One will become the Leader automatically.

### Verify

Check the logs:

```
# Leader node
This node elected as leader: node1:12345
Starting to load scheduled tasks (this node is leader)

# Follower node
Scheduler infrastructure initialized
# (no "elected as leader" message)
```

## Failover

### Graceful Shutdown (Ctrl+C / SIGTERM)

```
Releasing leader lock: node1:12345
Stopping scheduler (this node lost leadership)
```

The Follower takes over **within 1 second**.

### Crash / Kill -9

The lease expires after 15 seconds. A Follower detects the expired lock and takes over automatically.

## Kubernetes / Docker Deployment

For container environments, you don't need to copy config files manually. Use environment variables to override database settings:

| Variable | Description |
|---|---|
| `GOCRON_DB_ENGINE` | `mysql` or `postgres` |
| `GOCRON_DB_HOST` | Database host |
| `GOCRON_DB_PORT` | Database port |
| `GOCRON_DB_USER` | Database user |
| `GOCRON_DB_PASSWORD` | Database password |
| `GOCRON_DB_DATABASE` | Database name |
| `GOCRON_DB_PREFIX` | Table prefix |
| `GOCRON_AUTH_SECRET` | JWT auth secret |

Example Kubernetes deployment:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: gocron
spec:
  replicas: 2
  selector:
    matchLabels:
      app: gocron
  template:
    metadata:
      labels:
        app: gocron
    spec:
      containers:
        - name: gocron
          image: gocron:latest
          ports:
            - containerPort: 5920
          env:
            - name: GOCRON_DB_ENGINE
              value: "postgres"
            - name: GOCRON_DB_HOST
              value: "postgres-service"
            - name: GOCRON_DB_PORT
              value: "5432"
            - name: GOCRON_DB_USER
              valueFrom:
                secretKeyRef:
                  name: gocron-db
                  key: username
            - name: GOCRON_DB_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: gocron-db
                  key: password
            - name: GOCRON_DB_DATABASE
              value: "gocron"
```

Both replicas will compete for the leader lock. Only one runs the scheduler at any time.

## Architecture

```
┌─────────────┐     ┌─────────────┐
│   Node 1    │     │   Node 2    │
│  (Leader)   │     │ (Follower)  │
│             │     │             │
│ ┌─────────┐ │     │ ┌─────────┐ │
│ │Scheduler│ │     │ │  Standby│ │
│ │ (active)│ │     │ │  (idle) │ │
│ └─────────┘ │     │ └─────────┘ │
│ ┌─────────┐ │     │ ┌─────────┐ │
│ │ Web UI  │ │     │ │ Web UI  │ │
│ │  & API  │ │     │ │  & API  │ │
│ └─────────┘ │     │ └─────────┘ │
└──────┬──────┘     └──────┬──────┘
       │                   │
       └───────┬───────────┘
               │
        ┌──────┴──────┐
        │  MySQL / PG │
        │             │
        │ scheduler   │
        │ _lock table │
        └─────────────┘
```

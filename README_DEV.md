Quick dev setup

1) Start Postgres (using Docker Compose):

```bash
# from repository root
docker compose up -d
```

2) Ensure `DATABASE_URL` is set (we provided `.env`):

```
DATABASE_URL=postgres://dev:1234@localhost:5432/local_db?sslmode=disable
```

3) Start backend:

```bash
# from repository root
go mod tidy
go run main.go
```

The backend listens on http://localhost:8000

4) Start frontend:

```bash
cd frontend
npm install
npm run dev
```

5) Quick API checks:

```bash
# register
curl -X POST http://localhost:8000/api/register -H "Content-Type: application/json" -d '{"username":"test","email":"test@example.com","password":"pass"}'

# login
curl -X POST http://localhost:8000/api/login -H "Content-Type: application/json" -d '{"email":"test@example.com","password":"pass"}'

# leaderboard
curl http://localhost:8000/api/leaderboard
```

Notes:
- `main.go` auto-creates the Postgres tables at startup. Do not use the `schema.sql` (it's SQLite formatted) unless you convert the code to SQLite.
- If you prefer I can: (A) create a `docker-compose` + `Makefile` or `run` scripts, or (B) modify the backend to use SQLite instead of Postgres. Tell me which option you want next.

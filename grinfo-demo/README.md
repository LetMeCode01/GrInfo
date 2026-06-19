# GrInfo Demo Setup

Acest folder contine setup minim pentru baza de date GrInfo.

## 1. Pornire PostgreSQL local (port 5433)

```bash
cd grinfo-demo
docker compose up -d
```

## 2. Conectare backend la DB GrInfo

In fisierul `.env` din radacina proiectului setati:

```env
DATABASE_URL=postgres://grinfo:grinfo123@localhost:5433/grinfo_db?sslmode=disable
```

## 3. Pornire backend + frontend

```bash
# backend (radacina proiect)
go run main.go

# frontend
cd frontend
npm install
npm run dev
```

## API-uri GrInfo disponibile

- `GET /api/grinfo/categories`
- `GET /api/grinfo/questions?category=all&limit=50`
- `POST /api/grinfo/incident` (anti-cheat log)
- `POST /api/grinfo/session` (salvare sesiune, auth)
- `GET /api/grinfo/profile` (dashboard GrInfo, auth)

## Structura intrebare (frontend)

Fiecare intrebare intoarsa de API contine:

```json
{
  "id": 1,
  "category": "orientate",
  "dificultate": "medie",
  "eloRating": 1150,
  "enunt": "...",
  "optiuni": ["...", "...", "...", "..."],
  "raspunsCorect": 2,
  "explicatieRaspuns": "...",
  "graphData": {
    "nodes": [{"id": "1"}],
    "edges": [{"source": "1", "target": "2"}]
  }
}
```

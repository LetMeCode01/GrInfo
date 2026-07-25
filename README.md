## Despre aplicație

**GrInfo** este o platformă web educațională adaptivă dedicată studiului teoriei grafurilor. Aplicația urmărește să transforme procesul clasic de evaluare într-o experiență interactivă și personalizată, prin ajustarea automată a dificultății întrebărilor în funcție de nivelul utilizatorului.

Sistemul utilizează o variantă adaptată a algoritmului **ELO** pentru estimarea nivelului de pregătire al studentului și pentru selectarea întrebărilor potrivite competenței acestuia. Pe lângă mecanismul de evaluare adaptivă, platforma oferă vizualizarea interactivă a grafurilor, un spațiu de lucru de tip **Scratchpad** pentru rezolvarea problemelor și un modul de analiză bazat pe inteligență artificială, care generează feedback personalizat pentru răspunsurile incorecte.

Aplicația este construită pe o arhitectură de tip client-server, folosind **React (Vite)** pentru frontend, **Go** pentru backend și **PostgreSQL** pentru stocarea datelor.

---

## Demo și prezentare

<p align="center">
	<a href="./Demo%20platforma.mkv" target="_blank">
		<img src="https://img.shields.io/badge/Demo%20platforma-Deschide-0F766E?style=for-the-badge&logo=vlcmediaplayer&logoColor=white" alt="Demo platforma" />
	</a>
	<a href="./Lucrare%20de%20licen%C8%9B%C4%83.pptx" target="_blank">
		<img src="https://img.shields.io/badge/Prezentare%20platforma-Deschide-B91C1C?style=for-the-badge&logo=microsoftpowerpoint&logoColor=white" alt="Prezentare platforma" />
	</a>
</p>

---

## Cerințe

* Node.js (versiunea 18 sau mai nouă)
* Go (versiunea 1.22 sau mai nouă)
* PostgreSQL
* Cheie API Google Gemini (opțional, pentru funcționalitatea de feedback AI)

---

## Configurarea bazei de date

1. Instalați și porniți PostgreSQL local.
2. Creați baza de date utilizată de aplicație.
3. Actualizați parametrii de conectare din fișierul de configurare al backend-ului (host, port, user, password, database).

---

## Rularea backend-ului

Accesați directorul backend și porniți serverul:

```bash
cd backend
go run .
```

Serverul va porni pe portul configurat în aplicație.

---

## Rularea frontend-ului

Accesați directorul frontend și instalați dependențele:

```bash
cd frontend
npm install
```

Porniți aplicația în modul de dezvoltare:

```bash
npm run dev
```

Frontend-ul va fi disponibil la adresa afișată de Vite (de regulă `http://localhost:5173`).

---

## Funcționalități principale

* autentificare și gestionare cont utilizator;
* quiz-uri adaptive bazate pe scor ELO;
* actualizarea dinamică a nivelului de dificultate;
* vizualizarea interactivă a grafurilor;
* Scratchpad pentru rezolvarea problemelor;
* dashboard cu statistici și evoluția scorului ELO;
* clasament al utilizatorilor;
* feedback personalizat generat cu ajutorul inteligenței artificiale;
* sistem de recenzii pentru utilizatori.

---

## Tehnologii utilizate

### Frontend

* React
* Vite
* React Router
* VIS Network

### Backend

* Go
* net/http
* JWT Authentication

### Bază de date

* PostgreSQL

### Servicii externe

* Google Gemini API (pentru feedback AI)

---

## SIMA MIHAI

Proiect realizat în cadrul lucrării de licență **„Platformă web educațională adaptivă pentru teoria grafurilor”**.

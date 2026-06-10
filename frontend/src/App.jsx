import React, { useState, useEffect } from "react";
import { BrowserRouter, Routes, Route, Link } from "react-router-dom";
import Home from "./Home";
import Profile from "./Profile";
import Terms from "./Terms";
import Footer from "./Footer";
import Register from "./Register";
import Login from "./Login";
import Reviews from "./Reviews";
import Leaderboard from "./Leaderboard";
import GrInfoQuiz from "./GrInfoQuiz";
import "./assets/main.css";

export default function App() {
  const [token, setToken] = useState(localStorage.getItem("token"));
  
  const handleAuthChange = () => {
    setToken(localStorage.getItem("token"));
  };
  useEffect(() => {
    window.addEventListener("storage", handleAuthChange);
    return () => window.removeEventListener("storage", handleAuthChange);
  }, []);
  const logout = () => {
    localStorage.clear();
    setToken(null); 
    window.location.href = "/";
  };
  return (
    <BrowserRouter>
      <header>
        <title>GrInfo</title>
        <div
          className="header-container"
          style={{
            display: "flex",
            alignItems: "center",
            justifyContent: "space-between",
            gap: "16px",
          }}
        >
          <div style={{ display: "flex", alignItems: "center" }}>
            <Link
              to="/"
              className="logo-link"
            >
              <svg
                width="32"
                height="32"
                viewBox="0 0 24 24"
                fill="none"
                xmlns="http://www.w3.org/2000/svg"
                aria-hidden="true"
              >
                <defs>
                  <linearGradient id="logo-gradient" x1="0%" y1="0%" x2="100%" y2="100%">
                    <stop offset="0%" stopColor="#1d4ed8" />
                    <stop offset="100%" stopColor="#1e3a8a" />
                  </linearGradient>
                </defs>
                <path
                  d="M3 7l9-4 9 4v10a1 1 0 0 1-1 1h-5v-6H8v6H4a1 1 0 0 1-1-1V7z"
                  fill="url(#logo-gradient)"
                />
                <path d="M8 13h8v5H8z" fill="#fff" />
              </svg>
              <span className="logo-text">GrInfo</span>
            </Link>
          </div>

          <nav
            className="nav-subjects"
            style={{
              display: "flex",
              gap: 12,
              justifyContent: "center",
            }}
          >
            <Link className="subject-link" to="/grinfo/quiz">
              📊 Quiz Grafuri
            </Link>
            <Link className="subject-link" to="/leaderboard" style={{ borderLeft: "1px solid #ddd", paddingLeft: "12px" }}>
              🏆 Top
            </Link>
          </nav>

          <div style={{ display: "flex", gap: 12, justifyContent: "flex-end" }}>
          {token ? (
            <>
              <Link className="profile-link" to="/profile">👤 Profil</Link>
              <button 
                className="profile-link" 
                onClick={logout}
                style={{ cursor: "pointer", background: "none", border: "2px solid #667eea" }}
              >
                🚪 Deconectare
              </button>
            </>
          ) : (
            <>
              <Link className="profile-link" to="/register">✨ Înregistrare</Link>
              <Link className="profile-link" to="/login">🔐 Autentificare</Link>
            </>
          )}
          </div>
        </div>
      </header>
      <main style={{ padding: 10 }}>
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/grinfo/quiz" element={<GrInfoQuiz />} />
          <Route path="/leaderboard" element={<Leaderboard />} />
          <Route path="/login" element={<Login onLogin={handleAuthChange} />} />
          <Route path="/register" element={<Register onRegister={handleAuthChange} />} />
          <Route path="/profile" element={<Profile />} />
          <Route path="/terms" element={<Terms />} />
          <Route path="/reviews" element={<Reviews />} />
        </Routes>
      </main>
      <Footer />
    </BrowserRouter>
  );
}

//BrowserRouter: Componentele de navigare sunt învelite într-un BrowserRouter pentru a permite navigarea între pagini fără reîncărcarea întregii pagini.
//Link: Componentele Link sunt folosite pentru a crea link-uri de navigare către diferite rute definite în aplicație.
//Routes și Route: Componentele Routes și Route sunt folosite pentru a defini rutele aplicației și pentru a specifica ce componentă să fie redată pentru fiecare rută.

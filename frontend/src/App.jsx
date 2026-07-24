import React, { useState, useEffect, useRef } from "react";
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
  const [theme, setTheme] = useState(localStorage.getItem("theme") || "light");
  const [isThemeMenuOpen, setIsThemeMenuOpen] = useState(false);
  const themeMenuRef = useRef(null);
  
  const handleAuthChange = () => {
    setToken(localStorage.getItem("token"));
  };
  useEffect(() => {
    window.addEventListener("storage", handleAuthChange);
    return () => window.removeEventListener("storage", handleAuthChange);
  }, []);

  useEffect(() => {
    document.documentElement.setAttribute("data-theme", theme);
    localStorage.setItem("theme", theme);
  }, [theme]);

  useEffect(() => {
    const closeOnOutsideClick = (event) => {
      if (themeMenuRef.current && !themeMenuRef.current.contains(event.target)) {
        setIsThemeMenuOpen(false);
      }
    };

    document.addEventListener("mousedown", closeOnOutsideClick);
    return () => document.removeEventListener("mousedown", closeOnOutsideClick);
  }, []);

  const logout = () => {
    localStorage.clear();
    setToken(null); 
    window.location.href = "/";
  };

  const selectTheme = (nextTheme) => {
    setTheme(nextTheme);
    setIsThemeMenuOpen(false);
  };

  return (
    <BrowserRouter>
      <header>
        <title>GrInfo</title>
        <div className="header-container">
          <div className="header-left">
            <div className="theme-switcher" ref={themeMenuRef}>
              <button
                type="button"
                className="theme-switcher-btn"
                onClick={() => setIsThemeMenuOpen((open) => !open)}
                aria-haspopup="menu"
                aria-expanded={isThemeMenuOpen}
                aria-label="Alege tema"
              >
                💡
              </button>
              {isThemeMenuOpen && (
                <div className="theme-switcher-menu" role="menu" aria-label="Temă">
                  <button
                    type="button"
                    role="menuitem"
                    className={`theme-option ${theme === "dark" ? "active" : ""}`}
                    onClick={() => selectTheme("dark")}
                  >
                    Dark mode
                  </button>
                  <button
                    type="button"
                    role="menuitem"
                    className={`theme-option ${theme === "light" ? "active" : ""}`}
                    onClick={() => selectTheme("light")}
                  >
                    Light mode
                  </button>
                </div>
              )}
            </div>
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
          >
            <Link className="subject-link header-pill pill-quiz" to="/grinfo/quiz">
              📊 Quiz Grafuri
            </Link>
            <Link className="subject-link nav-chip-link header-pill pill-top" to="/leaderboard">
              🏆 Leaderboard
            </Link>
          </nav>

          <div className="header-actions">
          {token ? (
            <>
              <Link className="profile-link header-action-btn" to="/profile">👤 Profil</Link>
              <button 
                className="profile-link logout-link header-action-btn" 
                onClick={logout}
              >
                🚪 Deconectare
              </button>
            </>
          ) : (
            <>
              <Link className="profile-link header-action-btn" to="/register">✨ Înregistrare</Link>
              <Link className="profile-link header-action-btn" to="/login">🔐 Autentificare</Link>
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

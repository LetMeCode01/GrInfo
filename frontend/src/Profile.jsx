import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import EloChart from "./components/EloChart";
import "./styles/profile.css";

export default function Profile() {
  const navigate = useNavigate();
  const [profile, setProfile] = useState(null);
  const [error, setError] = useState(null);
  const [loading, setLoading] = useState(true);
  const [grInfoProfile, setGrInfoProfile] = useState(null);

  useEffect(() => {
    const token = localStorage.getItem("token");
    if (!token) {
      navigate("/login");
    } else {
      loadProfile();
      loadGrInfoProfile();
    }

    const refreshGrInfo = () => loadGrInfoProfile();
    window.addEventListener("grinfo:profile-updated", refreshGrInfo);

    return () => {
      window.removeEventListener("grinfo:profile-updated", refreshGrInfo);
    };
  }, [navigate]);

  const loadProfile = async () => {
    const token = localStorage.getItem("token");
    
    if (!token) {
      setError("Nu ești autentificat. Te rugăm să te conectezi.");
      setLoading(false);
      return;
    }

    try {
      const res = await fetch("/api/profile", {
        headers: {
          "Authorization": `Bearer ${token}`,
        },
      });
      
      if (!res.ok) {
        if (res.status === 401) {
          localStorage.clear();
          throw new Error("Sesiune expirată. Te rugăm să te autentifici din nou.");
        }
        throw new Error("Eroare la încărcarea profilului");
      }
      
      const data = await res.json();
      setProfile(data);
    } catch (e) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  };

  const loadGrInfoProfile = async () => {
    const token = localStorage.getItem("token");
    if (!token) return;

    try {
      const res = await fetch("/api/grinfo/profile", {
        headers: {
          "Authorization": `Bearer ${token}`,
        },
      });

      if (res.ok) {
        const data = await res.json();
        setGrInfoProfile(data);
      }
    } catch (e) {
      console.error("Error loading GrInfo profile:", e);
    }
  };

  const handleLogout = () => {
    localStorage.clear();
    navigate("/login");
  };

  const formatGrInfoDateTime = (item) => {
    const rawValue = item.endedAt || item.startedAt || item.createdAt;
    const date = new Date(rawValue);
    if (Number.isNaN(date.getTime())) {
      return "Data indisponibilă";
    }

    const datePart = date.toLocaleDateString("ro-RO", {
      year: "numeric",
      month: "long",
      day: "numeric",
    });
    const timePart = date.toLocaleTimeString("ro-RO", {
      hour: "2-digit",
      minute: "2-digit",
      second: "2-digit",
      hour12: false,
    });

    return `📅 ${datePart} • 🕒 ${timePart}`;
  };

  const formatChartDate = (rawValue) => {
    const date = new Date(rawValue);
    if (Number.isNaN(date.getTime())) {
      return "Data indisponibila";
    }

    return date.toLocaleDateString("ro-RO", {
      day: "numeric",
      month: "short",
    });
  };

  const rawEloHistory = Array.isArray(grInfoProfile?.eloHistory)
    ? grInfoProfile.eloHistory
    : Array.isArray(grInfoProfile?.history)
      ? grInfoProfile.history
      : [];

  const grInfoMatchesData = rawEloHistory
    ? rawEloHistory
        .filter((item) => item && (item.elo != null || item.finalElo != null))
        .map((item) => ({
          elo: Number(item.elo != null ? item.elo : item.finalElo),
          date: formatChartDate(item.at || item.endedAt || item.startedAt || item.createdAt),
        }))
        .filter((item) => Number.isFinite(item.elo))
    : [];
    
  if (loading) {
    return (
      <div className="profile-page" style={{ 
        display: "flex", 
        alignItems: "center", 
        justifyContent: "center",
        backgroundColor: "white", 
        background: "white" 
      }}>
        <div style={{ fontSize: 18, color: "#667eea", fontWeight: 600 }}>Se încarcă profilul...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="profile-page" style={{ 
        display: "flex", 
        alignItems: "center", 
        justifyContent: "center",
        flexDirection: "column",
        gap: 20
      }}>
        <div style={{ color: "white", fontSize: 18, fontWeight: 600 }}>{error}</div>
        <button className="logout-btn" onClick={() => navigate("/login")}>
          Autentifică-te
        </button>
      </div>
    );
  }

  return (
    <div className="profile-page">
      <div className="profile-container">
        {/* Header with username and logout */}
        <div className="profile-header">
          <div className="profile-title">
            <h1><span className="wave-emoji">👋</span> Bună, {profile.username}!</h1>
            <p>Profilul tău GrInfo și istoricul quiz-urilor.</p>
          </div>
          <button
            onClick={handleLogout}
            className="logout-btn"
          >
            Deconectare
          </button>
        </div>

        {/* GrInfo Dashboard */}
        {grInfoProfile && (
          <div className="history-section" style={{ marginTop: 24 }}>
            <h2 className="section-title">🧠 GrInfo Dashboard</h2>

            <div className="stats-grid" style={{ marginTop: 12 }}>
              <div className="stat-card level">
                <div className="stat-icon">📈</div>
                <div className="stat-value">{Number(grInfoProfile.currentElo || 1000).toFixed(1)}</div>
                <div className="stat-label">ELO curent GrInfo</div>
          
              </div>

              <div className="stat-card quizzes">
                <div className="stat-icon">🧪</div>
                <div className="stat-value">{grInfoProfile.totalSessions || 0}</div>
                <div className="stat-label">Sesiuni GrInfo</div>
                <div className="stat-sublabel">Quiz-uri finalizate</div>
              </div>

              <div className="stat-card streak">
                <div className="stat-icon">✅</div>
                <div className="stat-value stat-value-compact">
                  <span className="stat-value-line">{grInfoProfile.totalCorrectAnswers || 0} /</span>
                  <span className="stat-value-line">{grInfoProfile.totalQuestionsAnswered || 0}</span>
                </div>
                <div className="stat-label">Răspunsuri corecte / totale</div>
                <div className="stat-sublabel">Quiz-uri × 10 întrebări</div>
              </div>

              <div className="stat-card xp">
                <div className="stat-icon">🎯</div>
                <div className="stat-value">{Number(grInfoProfile.accuracy || 0).toFixed(1)}%</div>
                <div className="stat-label">Acuratețe GrInfo</div>
                <div className="stat-sublabel">Bazată pe răspunsuri corecte</div>
              </div>
            </div>

            <EloChart matchesData={grInfoMatchesData} />

            {Array.isArray(grInfoProfile.history) && grInfoProfile.history.length > 0 ? (
              <div className="history-list" style={{ marginTop: 16 }}>
                {grInfoProfile.history.map((item, idx) => (
                  <div key={idx} className={`history-card ${item.isInProgress ? "in-progress" : ""}`}>
                    <div className="history-header">
                      <div>
                        <div className="history-course-title">
                          <span>GrInfo</span>
                          <span> • {item.category} • {item.difficulty || "all"}</span>
                        </div>
                        <div className="history-date">
                          {formatGrInfoDateTime(item)}
                        </div>
                      </div>
                      <div className="history-score-badge">
                        <div className="score-text good">
                          ELO {Number(item.finalElo).toFixed(1)}
                        </div>
                        <div style={{ fontSize: 18 }}>📊</div>
                      </div>
                    </div>

                    <div className="progress-track">
                      <div
                        className="progress-fill good"
                        style={{
                          width: `${Math.min(100, ((item.correctAnswers || 0) / Math.max(1, item.totalQuestions || 10)) * 100)}%`,
                        }}
                      />
                    </div>

                    <div className="history-footer">
                      {item.correctAnswers}/{item.totalQuestions} corecte
                    </div>
                  </div>
                ))}
              </div>
            ) : (
              <div className="empty-state" style={{ marginTop: 16 }}>
                <div className="empty-icon">📘</div>
                <p style={{ margin: 0, fontSize: 16 }}>Nu ai încă sesiuni GrInfo salvate.</p>
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  );
}

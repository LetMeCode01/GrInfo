import { Link } from "react-router-dom";
import { useState } from "react";

export default function Home() {
  const [reviewName, setReviewName] = useState("");
  const [reviewContent, setReviewContent] = useState("");
  const [reviewRating, setReviewRating] = useState(0);
  const [isSubmittingReview, setIsSubmittingReview] = useState(false);

  const showToast = (message, type = "success") => {
    const toast = document.createElement("div");
    toast.className = `toast toast-${type}`;
    toast.textContent = message;
    toast.style.cssText = `
      position: fixed;
      top: 20px;
      right: 20px;
      background: ${type === "success" ? "#4caf50" : "#f44336"};
      color: white;
      padding: 16px 24px;
      border-radius: 8px;
      box-shadow: 0 4px 6px rgba(0,0,0,0.2);
      z-index: 9999;
      font-size: 14px;
      animation: slideIn 0.3s ease-out;
    `;
    document.body.appendChild(toast);
    setTimeout(() => toast.remove(), 4000);
  };

  const handleSubmitReview = async (e) => {
    e.preventDefault();

    if (!reviewContent.trim() || reviewRating === 0) {
      showToast("Completează toate câmpurile", "error");
      return;
    }

    setIsSubmittingReview(true);

    try {
      const token = localStorage.getItem("token");
      const headers = {
        "Content-Type": "application/json",
      };
      
      if (token) {
        headers["Authorization"] = `Bearer ${token}`;
      }

      const response = await fetch("http://localhost:8000/api/reviews", {
        method: "POST",
        headers: headers,
        body: JSON.stringify({
          name: reviewName.trim(),
          content: reviewContent,
          rating: reviewRating,
        }),
      });

      const data = await response.json();

      if (response.ok) {
        showToast("Review-ul a fost trimis cu succes!", "success");
        setReviewName("");
        setReviewContent("");
        setReviewRating(0);
      } else {
        showToast(data.error || "Eroare la trimiterea review-ului", "error");
      }
    } catch (error) {
      showToast("Eroare de conexiune: " + error.message, "error");
    } finally {
      setIsSubmittingReview(false);
    }
  };

  return (
    <div className="home-page">
      <style>{`
        @keyframes slideIn {
          from {
            transform: translateX(400px);
            opacity: 0;
          }
          to {
            transform: translateX(0);
            opacity: 1;
          }
        }
      `}</style>

      <header className="home-header">
        <h1 className="home-title">GrInfo</h1>
      </header>

      <section className="site-intro">
        <span className="intro-badge">Bun venit!</span>
        <h2 className="intro-title">Platformă demo pentru Grafuri</h2>
        <p className="intro-text">
          GrInfo este viitorul pentru învățarea de Algoritmi și structuri de date, axat momentan pe quizuri de grafuri
          orientate și neorientate. Quiz-urile includ selecție adaptivă pe baza algoritmului ELO, 
          întrebare unică pe pagină și monitorizare anti-cheat.
        </p>
      </section>

      <section className="subjects-section">
        <div className="subject-grid">
          <div className="subject-card">
            <div className="subject-head">
              <span className="subject-icon" aria-hidden>
                🎢
              </span>
              <h4>Grafuri neorientate</h4>
            </div>
            <p>
              Concepte de conectivitate, arbori, cicluri și parcurgeri clasice
              pentru grafuri fără sens pe muchii.
            </p>
          </div>

          <div className="subject-card">
            <div className="subject-head">
              <span className="subject-icon" aria-hidden>
                ➡️
              </span>
              <h4>Grafuri orientate</h4>
            </div>
            <p>
              DAG, grade interne/externe, sortare topologică și algoritmi de
              drum minim pentru grafuri cu arce direcționale.
            </p>
          </div>

          <div className="subject-card">
            <div className="subject-head">
              <span className="subject-icon" aria-hidden>
                💻
              </span>
              <h4>Quiz adaptiv ELO</h4>
            </div>
            <p>
              După fiecare răspuns, se recalculează ELO și se alege următoarea
              întrebare necompletată cu dificultatea cea mai apropiată.
            </p>
          </div>
        </div>
      </section>

      <section className="featured-section">
        <h3 style={{textAlign:"center"}}>Recomandate</h3>
        <div className="featured-grid">
          <Link className="featured-card" to="/grinfo/quiz?category=orientate">
            <strong>Quiz: grafuri orientate</strong>
            <div className="muted">10 întrebări secvențiale, explicații complete</div>
          </Link>

          <Link className="featured-card" to="/grinfo/quiz?category=neorientate">
            <strong>Quiz: grafuri neorientate</strong>
            <div className="muted">Set pentru demo și validare logică</div>
          </Link>

          <Link className="featured-card" to="/grinfo/quiz">
            <strong>Start demo GrInfo</strong>
            <div className="muted">ELO inițial 1000, anti-cheat activ</div>
          </Link>
        </div>
      </section>

      <section className="reviews-section">
        <h3>Lasă un review</h3>
        <p className="reviews-intro">
          Împărtășește experiența ta cu GrInfo! Review-urile tale ne ajută să îmbunătățim platforma.
        </p>
        
        <div className="review-form-container">
          <form className="review-form" onSubmit={handleSubmitReview}>
            <div className="review-name-field">
              <label htmlFor="review-name">Nume (opțional):</label>
              <input
                id="review-name"
                type="text"
                className="review-name-input"
                placeholder="Introdu numele tău"
                value={reviewName}
                onChange={(e) => setReviewName(e.target.value)}
              />
            </div>
            <div className="review-rating">
              <label>Rating:</label>
              <div className="rating-stars">
                {[1, 2, 3, 4, 5].map((star) => (
                  <button
                    key={star}
                    type="button"
                    className={`star-btn ${reviewRating >= star ? 'active' : ''}`}
                    onClick={() => setReviewRating(star)}
                    aria-label={`${star} stele`}
                  >
                    ★
                  </button>
                ))}
              </div>
            </div>
            <textarea
              className="review-textarea"
              placeholder="Scrie review-ul tău aici..."
              value={reviewContent}
              onChange={(e) => setReviewContent(e.target.value)}
              rows={4}
              required
            />
            <button
              type="submit"
              className="cta-primary"
              disabled={isSubmittingReview || !reviewContent.trim() || reviewRating === 0}//
            >
              {isSubmittingReview ? "Se trimite..." : "Trimite review"}
            </button>
          </form>
        </div>

        <div className="reviews-actions">
          <Link to="/reviews" className="cta-secondary">
            Vezi toate review-urile →
          </Link>
        </div>
      </section>
    </div>
  );
}

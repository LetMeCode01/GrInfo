import { useState, useEffect } from "react";
import { Link } from "react-router-dom";

export default function Reviews() {
  const [reviews, setReviews] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    fetchReviews();
  }, []);

  const fetchReviews = async () => {
    try {
      setIsLoading(true);
      setError(null);
      
      const response = await fetch("http://localhost:8000/api/reviews", {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
        },
      });
      
      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        throw new Error(errorData.error || `Eroare HTTP: ${response.status}`);
      }

      const data = await response.json();
      setReviews(Array.isArray(data) ? data : []);
      setError(null);
    } catch (err) {
      console.error("Error fetching reviews:", err);
      setError(err.message || "Eroare la încărcarea review-urilor");
      setReviews([]);
    } finally {
      setIsLoading(false);
    }
  };

  const formatDate = (dateString) => {
    const date = new Date(dateString);
    return date.toLocaleDateString("ro-RO", {
      year: "numeric",
      month: "long",
      day: "numeric",
    });
  };

  const renderStars = (rating) => {
    return "★".repeat(rating) + "☆".repeat(5 - rating);
  };

  if (isLoading) {
    return (
      <div className="reviews-page">
        <div className="loading-container">
          <p>Se încarcă review-urile...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="reviews-page">
        <div className="error-container">
          <p>Eroare: {error}</p>
          <button onClick={fetchReviews} className="cta-primary">
            Încearcă din nou
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="reviews-page">
      <div className="reviews-page-header">
        <h1>Toate review-urile</h1>
        <p className="reviews-count">
          {reviews.length} {reviews.length === 1 ? "review" : "review-uri"}
        </p>
        <Link to="/" className="back-link">
          ← Înapoi la pagina principală
        </Link>
      </div>

      {reviews.length === 0 ? (
        <div className="no-reviews">
          <p>Nu există review-uri încă. Fii primul care lasă un review!</p>
          <Link to="/" className="cta-primary">
            Lasă un review
          </Link>
        </div>
      ) : (
        <div className="reviews-list">
          {reviews.map((review) => (
            <div key={review.id} className="review-card">
              <div className="review-header">
                <div className="review-author">
                  <div className="author-avatar">
                    {review.username.charAt(0).toUpperCase()}
                  </div>
                  <div className="author-info">
                    <div className="author-name">{review.username}</div>
                    <div className="review-date">{formatDate(review.createdAt)}</div>
                  </div>
                </div>
                <div className="review-rating-display">
                  <span className="stars">{renderStars(review.rating)}</span>
                  <span className="rating-number">{review.rating}/5</span>
                </div>
              </div>
              <div className="review-content">{review.content}</div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

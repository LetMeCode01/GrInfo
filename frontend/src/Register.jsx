import React, { useState } from "react";
import { Link, useNavigate } from "react-router-dom";

export default function Register({ onRegister }) {
  const navigate = useNavigate();
  const [formData, setFormData] = useState({
    name: "",
    email: "",
    password: "",
  });
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
    setError("");
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError("");
    
    if (!formData.name || !formData.email || !formData.password) {
      setError("Toate câmpurile sunt obligatorii");
      return;
    }
    
    setLoading(true);
    
    try {
      const response = await fetch("/api/register", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          username: formData.name,
          email: formData.email,
          password: formData.password,
        }),
      });
      
      const data = await response.json();
      
      if (!response.ok) {
        throw new Error(data.error || "Înregistrare eșuată");
      }
      
      localStorage.setItem("token", data.token);
      localStorage.setItem("username", data.username);
      localStorage.setItem("userId", data.userId);
      
      if (onRegister) onRegister();
      
      alert(`Bun venit, ${data.username}! Contul a fost creat cu succes.`);
      localStorage.setItem("doOnboarding", "true");
      navigate("/");
      
    } catch (err) {
      setError(err.message || "Eroare la înregistrare");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="auth-page">
      <div className="auth-card">
        <div className="auth-head">
          <h1>Înregistrare</h1>
          <p>Creează un cont pentru a accesa GrInfo</p>
        </div>

        <form onSubmit={handleSubmit}>
          {error && (
            <div className="auth-error">
              {error}
            </div>
          )}
          
          <div className="auth-field">
            <label>Nume complet</label>
            <input
              type="text"
              name="name"
              value={formData.name}
              onChange={handleChange}
              placeholder="Ex. Ion Popescu"
              className="auth-input"
            />
          </div>

          <div className="auth-field">
            <label>Email</label>
            <input
              type="email"
              name="email"
              value={formData.email}
              onChange={handleChange}
              placeholder="email@exemplu.com"
              className="auth-input"
            />
          </div>

          <div className="auth-field auth-field-last">
            <label>Parolă</label>
            <input
              type="password"
              name="password"
              value={formData.password}
              onChange={handleChange}
              placeholder="••••••••"
              className="auth-input"
            />
          </div>

          <button
            type="submit"
            disabled={loading}
            className="auth-submit"
          >
            {loading ? "Se înregistrează..." : "Înregistrează-te"}
          </button>
        </form>
      </div>
    </div>
  );
}
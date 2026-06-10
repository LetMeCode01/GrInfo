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
    <div style={{
        minHeight: "calc(100vh - 120px)",
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        background: "linear-gradient(135deg, #667eea 0%, #764ba2 100%)",
        padding: "20px",
      }}>
      <div style={{
          background: "#fff",
          borderRadius: "12px",
          boxShadow: "0 10px 40px rgba(0, 0, 0, 0.15)",
          maxWidth: "450px",
          width: "100%",
          padding: "40px",
        }}>
        <div style={{ textAlign: "center", marginBottom: 32 }}>
          <h1 style={{ fontSize: "28px", fontWeight: 700, color: "#222", margin: "0 0 8px 0" }}>Înregistrare</h1>
          <p style={{ color: "#666", fontSize: "14px", margin: 0 }}>Creează un cont pentru a accesa GrInfo</p>
        </div>

        <form onSubmit={handleSubmit}>
          {error && (
            <div style={{ padding: "12px", marginBottom: "20px", backgroundColor: "#fee", color: "#c33", borderRadius: "8px", fontSize: "14px" }}>
              {error}
            </div>
          )}
          
          <div style={{ marginBottom: 20 }}>
            <label style={{ display: "block", marginBottom: "8px", fontWeight: 600, color: "#333", fontSize: "14px" }}>Nume complet</label>
            <input
              type="text"
              name="name"
              value={formData.name}
              onChange={handleChange}
              placeholder="Ex. Ion Popescu"
              style={{ width: "100%", padding: "12px 14px", border: "1px solid #ddd", borderRadius: "8px", fontSize: "14px", boxSizing: "border-box", outline: "none" }}
            />
          </div>

          <div style={{ marginBottom: 20 }}>
            <label style={{ display: "block", marginBottom: "8px", fontWeight: 600, color: "#333", fontSize: "14px" }}>Email</label>
            <input
              type="email"
              name="email"
              value={formData.email}
              onChange={handleChange}
              placeholder="email@exemplu.com"
              style={{ width: "100%", padding: "12px 14px", border: "1px solid #ddd", borderRadius: "8px", fontSize: "14px", boxSizing: "border-box", outline: "none" }}
            />
          </div>

          <div style={{ marginBottom: 28 }}>
            <label style={{ display: "block", marginBottom: "8px", fontWeight: 600, color: "#333", fontSize: "14px" }}>Parolă</label>
            <input
              type="password"
              name="password"
              value={formData.password}
              onChange={handleChange}
              placeholder="••••••••"
              style={{ width: "100%", padding: "12px 14px", border: "1px solid #ddd", borderRadius: "8px", fontSize: "14px", boxSizing: "border-box", outline: "none" }}
            />
          </div>

          <button
            type="submit"
            style={{
              width: "100%",
              padding: "12px 16px",
              background: "linear-gradient(135deg, #667eea 0%, #764ba2 100%)",
              color: "#fff",
              border: "none",
              borderRadius: "8px",
              fontSize: "15px",
              fontWeight: 600,
              cursor: loading ? "not-allowed" : "pointer",
            }}
          >
            {loading ? "Se înregistrează..." : "Înregistrează-te"}
          </button>
        </form>
      </div>
    </div>
  );
}
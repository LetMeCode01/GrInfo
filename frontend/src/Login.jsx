import React, { useState } from "react";
import { Link, useNavigate } from "react-router-dom";

export default function Login({ onLogin }) {
  const navigate = useNavigate();
  const [formData, setFormData] = useState({
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
    
    if (!formData.email || !formData.password) {
      setError("Email și parola sunt obligatorii");
      return;
    }
    
    setLoading(true);
    
    try {
      const response = await fetch("/api/login", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          email: formData.email,
          password: formData.password,
        }),
      });
      
      const data = await response.json();
      
      if (!response.ok) {
        throw new Error(data.error || "Autentificare eșuată");
      }
      
      localStorage.setItem("token", data.token);
      localStorage.setItem("username", data.username);
      localStorage.setItem("userId", data.userId);
      
      if (onLogin) onLogin();
      
      navigate("/profile");
      
    } catch (err) {
      setError(err.message || "Email sau parolă incorectă");
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
          <h1 style={{ fontSize: "28px", fontWeight: 700, color: "#222", margin: "0 0 8px 0" }}>
            Autentificare
          </h1>
          <p style={{ color: "#666", fontSize: "14px", margin: 0 }}>
            Conectează-te la contul tău GrInfo
          </p>
        </div>

        <form onSubmit={handleSubmit}>
          {error && (
            <div style={{ padding: "12px", marginBottom: "20px", backgroundColor: "#fee", color: "#c33", borderRadius: "8px", fontSize: "14px" }}>
              {error}
            </div>
          )}
          
          <div style={{ marginBottom: 20 }}>
            <label style={{ display: "block", marginBottom: "8px", fontWeight: 600, color: "#333", fontSize: "14px" }}>
              Email
            </label>
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
            <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center", marginBottom: "8px" }}>
              <label style={{ fontWeight: 600, color: "#333", fontSize: "14px" }}>Parolă</label>
              <Link to="/" style={{ fontSize: "13px", color: "#667eea", textDecoration: "none" }}>Ați uitat parola?</Link>
            </div>
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
            disabled={loading}
            style={{
              width: "100%",
              padding: "12px 16px",
              background: loading ? "#ccc" : "linear-gradient(135deg, #667eea 0%, #764ba2 100%)",
              color: "#fff",
              border: "none",
              borderRadius: "8px",
              fontSize: "15px",
              fontWeight: 600,
              cursor: loading ? "not-allowed" : "pointer",
            }}
          >
            {loading ? "Se autentifică..." : "Autentifică-te"}
          </button>
        </form>
      </div>
    </div>
  );
}
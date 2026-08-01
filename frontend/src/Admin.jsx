import React, { useEffect, useMemo, useState } from "react";
import { useNavigate } from "react-router-dom";
import "./styles/admin.css";

const EMPTY_FORM = {
  category: "orientate",
  dificultate: "usoara",
  eloRating: 1000,
  enunt: "",
  explicatieRaspuns: "",
  graphData: "{}",
  optiuni: ["", "", "", ""],
  raspunsCorect: 0,
  isActive: true,
};

export default function Admin() {
  const navigate = useNavigate();
  const [questions, setQuestions] = useState([]);
  const [categories, setCategories] = useState([]);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [message, setMessage] = useState("");
  const [error, setError] = useState("");
  const [filterCategory, setFilterCategory] = useState("all");
  const [includeInactive, setIncludeInactive] = useState(true);
  const [editingId, setEditingId] = useState(null);
  const [form, setForm] = useState(EMPTY_FORM);

  const authHeaders = useMemo(() => {
    const token = localStorage.getItem("token");
    return {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    };
  }, []);

  useEffect(() => {
    const token = localStorage.getItem("token");
    if (!token) {
      navigate("/login");
      return;
    }

    const load = async () => {
      setLoading(true);
      setError("");
      try {
        const [catsRes, qRes] = await Promise.all([
          fetch("/api/grinfo/categories"),
          fetch("/api/grinfo/admin/questions?includeInactive=1", {
            headers: authHeaders,
          }),
        ]);

        if (!catsRes.ok) {
          throw new Error("Nu am putut încărca categoriile");
        }
        if (!qRes.ok) {
          if (qRes.status === 401) {
            localStorage.clear();
            navigate("/login");
            return;
          }
          if (qRes.status === 403) {
            throw new Error("Nu ai permisiuni de administrator pentru această pagină.");
          }
          throw new Error("Nu am putut încărca întrebările");
        }

        const catsData = await catsRes.json();
        const qData = await qRes.json();
        setCategories(Array.isArray(catsData.categories) ? catsData.categories : []);
        setQuestions(Array.isArray(qData.questions) ? qData.questions : []);
      } catch (e) {
        setError(e.message || "Eroare la încărcarea datelor");
      } finally {
        setLoading(false);
      }
    };

    load();
  }, [authHeaders, navigate]);

  const loadQuestions = async (category = filterCategory, include = includeInactive) => {
    setLoading(true);
    setError("");

    const params = new URLSearchParams();
    params.set("category", category);
    if (include) {
      params.set("includeInactive", "1");
    }

    try {
      const res = await fetch(`/api/grinfo/admin/questions?${params.toString()}`, {
        headers: authHeaders,
      });
      if (!res.ok) {
        if (res.status === 403) {
          throw new Error("Nu ai permisiuni de administrator pentru această operație.");
        }
        throw new Error("Nu am putut încărca întrebările");
      }
      const data = await res.json();
      setQuestions(Array.isArray(data.questions) ? data.questions : []);
    } catch (e) {
      setError(e.message || "Eroare la încărcare");
    } finally {
      setLoading(false);
    }
  };

  const updateOption = (index, value) => {
    const next = [...form.optiuni];
    next[index] = value;
    setForm((prev) => ({ ...prev, optiuni: next }));
  };

  const resetForm = () => {
    setEditingId(null);
    setForm(EMPTY_FORM);
    setMessage("");
    setError("");
  };

  const onEdit = (question) => {
    setEditingId(question.id);
    setMessage("");
    setError("");
    setForm({
      category: question.category,
      dificultate: question.dificultate,
      eloRating: question.eloRating,
      enunt: question.enunt,
      explicatieRaspuns: question.explicatieRaspuns,
      graphData: question.graphData || "{}",
      optiuni: Array.isArray(question.optiuni) && question.optiuni.length === 4 ? question.optiuni : ["", "", "", ""],
      raspunsCorect: Number.isInteger(question.raspunsCorect) ? question.raspunsCorect : 0,
      isActive: Boolean(question.isActive),
    });
  };

  const onSubmit = async (e) => {
    e.preventDefault();
    setSaving(true);
    setMessage("");
    setError("");

    const endpoint = editingId
      ? `/api/grinfo/admin/questions?id=${editingId}`
      : "/api/grinfo/admin/questions";
    const method = editingId ? "PUT" : "POST";

    try {
      const res = await fetch(endpoint, {
        method,
        headers: authHeaders,
        body: JSON.stringify(form),
      });

      const data = await res.json();
      if (!res.ok) {
        throw new Error(data.error || "Operația a eșuat");
      }

      setMessage(editingId ? "Întrebare actualizată." : "Întrebare adăugată.");
      resetForm();
      await loadQuestions();
    } catch (e) {
      setError(e.message || "Eroare la salvare");
    } finally {
      setSaving(false);
    }
  };

  const onSoftDelete = async (id) => {
    const ok = window.confirm("Sigur vrei să dezactivezi întrebarea?");
    if (!ok) {
      return;
    }

    setMessage("");
    setError("");
    try {
      const res = await fetch(`/api/grinfo/admin/questions?id=${id}`, {
        method: "DELETE",
        headers: authHeaders,
      });
      const data = await res.json();
      if (!res.ok) {
        throw new Error(data.error || "Nu am putut dezactiva întrebarea");
      }

      if (editingId === id) {
        resetForm();
      }
      setMessage("Întrebarea a fost dezactivată.");
      await loadQuestions();
    } catch (e) {
      setError(e.message || "Eroare la dezactivare");
    }
  };

  const onHardDelete = async (id) => {
    const ok = window.confirm("Ștergere permanentă? Această acțiune nu poate fi anulată.");
    if (!ok) {
      return;
    }

    setMessage("");
    setError("");
    try {
      const res = await fetch(`/api/grinfo/admin/questions?id=${id}&hard=1`, {
        method: "DELETE",
        headers: authHeaders,
      });
      const data = await res.json();
      if (!res.ok) {
        throw new Error(data.error || "Nu am putut șterge întrebarea");
      }

      if (editingId === id) {
        resetForm();
      }
      setMessage("Întrebarea a fost ștearsă permanent.");
      await loadQuestions();
    } catch (e) {
      setError(e.message || "Eroare la ștergere");
    }
  };

  return (
    <div className="admin-page">
      <div className="admin-shell">
        <div className="admin-header">
          <h1>Admin GrInfo</h1>
          <p>Adaugă, editează, dezactivează sau șterge permanent întrebări.</p>
        </div>

        <div className="admin-controls">
          <label>
            Categorie
            <select
              value={filterCategory}
              onChange={(e) => setFilterCategory(e.target.value)}
            >
              <option value="all">Toate</option>
              {categories.map((cat) => (
                <option key={cat.slug} value={cat.slug}>
                  {cat.name}
                </option>
              ))}
            </select>
          </label>

          <label className="admin-checkbox">
            <input
              type="checkbox"
              checked={includeInactive}
              onChange={(e) => setIncludeInactive(e.target.checked)}
            />
            Include întrebări inactive
          </label>

          <button type="button" className="admin-btn" onClick={() => loadQuestions()}>
            Reîncarcă
          </button>
          <button type="button" className="admin-btn admin-btn-muted" onClick={resetForm}>
            Întrebare nouă
          </button>
        </div>

        {message && <div className="admin-message">{message}</div>}
        {error && <div className="admin-error">{error}</div>}

        <form className="admin-form" onSubmit={onSubmit}>
          <h2>{editingId ? `Editare întrebarea #${editingId}` : "Adaugă întrebare"}</h2>

          <div className="admin-grid">
            <label>
              Categorie
              <select
                value={form.category}
                onChange={(e) => setForm((prev) => ({ ...prev, category: e.target.value }))}
              >
                <option value="orientate">orientate</option>
                <option value="neorientate">neorientate</option>
              </select>
            </label>

            <label>
              Dificultate
              <select
                value={form.dificultate}
                onChange={(e) => setForm((prev) => ({ ...prev, dificultate: e.target.value }))}
              >
                <option value="usoara">usoara</option>
                <option value="medie">medie</option>
                <option value="grea">grea</option>
              </select>
            </label>

            <label>
              ELO
              <input
                type="number"
                min="1"
                value={form.eloRating}
                onChange={(e) =>
                  setForm((prev) => ({ ...prev, eloRating: Number(e.target.value || 0) }))
                }
              />
            </label>

            <label className="admin-checkbox">
              <input
                type="checkbox"
                checked={form.isActive}
                onChange={(e) => setForm((prev) => ({ ...prev, isActive: e.target.checked }))}
              />
              Activă
            </label>
          </div>

          <label>
            Enunț
            <textarea
              value={form.enunt}
              onChange={(e) => setForm((prev) => ({ ...prev, enunt: e.target.value }))}
              rows={3}
              required
            />
          </label>

          <label>
            Explicație răspuns
            <textarea
              value={form.explicatieRaspuns}
              onChange={(e) =>
                setForm((prev) => ({ ...prev, explicatieRaspuns: e.target.value }))
              }
              rows={3}
              required
            />
          </label>

          <label>
            Graph data (JSON)
            <textarea
              value={form.graphData}
              onChange={(e) => setForm((prev) => ({ ...prev, graphData: e.target.value }))}
              rows={4}
            />
          </label>

          <div className="admin-options">
            {form.optiuni.map((opt, idx) => (
              <label key={idx}>
                Opțiunea {idx + 1}
                <input
                  type="text"
                  value={opt}
                  onChange={(e) => updateOption(idx, e.target.value)}
                  required
                />
              </label>
            ))}
          </div>

          <label>
            Răspuns corect
            <select
              value={form.raspunsCorect}
              onChange={(e) =>
                setForm((prev) => ({ ...prev, raspunsCorect: Number(e.target.value) }))
              }
            >
              <option value={0}>Opțiunea 1</option>
              <option value={1}>Opțiunea 2</option>
              <option value={2}>Opțiunea 3</option>
              <option value={3}>Opțiunea 4</option>
            </select>
          </label>

          <div className="admin-form-actions">
            <button type="submit" className="admin-btn" disabled={saving}>
              {saving ? "Se salvează..." : editingId ? "Salvează modificările" : "Adaugă întrebare"}
            </button>
            {editingId && (
              <button type="button" className="admin-btn admin-btn-muted" onClick={resetForm}>
                Renunță
              </button>
            )}
          </div>
        </form>

        <div className="admin-list">
          <h2>Întrebări ({questions.length})</h2>
          {loading ? (
            <p>Se încarcă...</p>
          ) : questions.length === 0 ? (
            <p>Nu există întrebări pentru filtrul selectat.</p>
          ) : (
            questions.map((q) => (
              <div key={q.id} className={`admin-card ${q.isActive ? "" : "is-inactive"}`}>
                <div className="admin-card-top">
                  <strong>#{q.id} · {q.category} · {q.dificultate} · ELO {q.eloRating}</strong>
                  <span>{q.isActive ? "Activă" : "Inactivă"}</span>
                </div>
                <p className="admin-enunt">{q.enunt}</p>
                <div className="admin-card-actions">
                  <button type="button" className="admin-btn" onClick={() => onEdit(q)}>
                    Editează
                  </button>
                  <button type="button" className="admin-btn admin-btn-muted" onClick={() => onSoftDelete(q.id)}>
                    Dezactivează
                  </button>
                  <button type="button" className="admin-btn admin-btn-danger" onClick={() => onHardDelete(q.id)}>
                    Șterge permanent
                  </button>
                </div>
              </div>
            ))
          )}
        </div>
      </div>
    </div>
  );
}

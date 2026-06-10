import React, { useEffect, useMemo, useRef, useState } from "react";
import { useSearchParams } from "react-router-dom";
import "./styles/grinfo-quiz.css";

const INITIAL_ELO = 1000;
const DIFFICULTY_ORDER = ["usoara", "medie", "grea"];
const QUIZ_STATE_KEY = "grinfo_quiz_active_state";

function normalizeQuestionKey(question) {
  return `${String(question?.category || "").trim().toLowerCase()}::${String(question?.enunt || "").trim().toLowerCase()}`;
}

function dedupeQuestions(questions) {
  const seen = new Set();
  const unique = [];

  for (const question of questions) {
    const key = normalizeQuestionKey(question);
    if (seen.has(key)) {
      continue;
    }

    seen.add(key);
    unique.push(question);
  }

  return unique;
}

function shuffleQuestions(questions) {
  const copy = [...questions];
  for (let i = copy.length - 1; i > 0; i -= 1) {
    const j = Math.floor(Math.random() * (i + 1));
    [copy[i], copy[j]] = [copy[j], copy[i]];
  }
  return copy;
}

function shuffleQuestionOptions(question) {
  const options = Array.isArray(question?.optiuni) ? question.optiuni : [];
  if (options.length === 0) return question;

  const items = options.map((text, index) => ({ text, isCorrect: index === question.raspunsCorect }));
  for (let i = items.length - 1; i > 0; i -= 1) {
    const j = Math.floor(Math.random() * (i + 1));
    [items[i], items[j]] = [items[j], items[i]];
  }

  const shuffledOptions = items.map((item) => item.text);
  const shuffledCorrectIndex = items.findIndex((item) => item.isCorrect);

  return {
    ...question,
    optiuni: shuffledOptions,
    raspunsCorect: shuffledCorrectIndex,
  };
}

function getInitialDifficultyByElo(currentElo) {
  if (currentElo <= 1200) return "usoara";
  if (currentElo <= 1500) return "medie";
  return "grea";
}

function pickNextQuestion(allQuestions, usedIds, targetDifficulty) {
  const remaining = allQuestions.filter((question) => !usedIds.has(question.id));
  if (remaining.length === 0) return null;

  const preferred = remaining.filter((question) => question.dificultate === targetDifficulty);
  const pool = preferred.length > 0 ? preferred : remaining;

  const randomPool = shuffleQuestions(pool);
  return randomPool[0] || null;
}

function lowerDifficulty(difficulty) {
  const idx = DIFFICULTY_ORDER.indexOf(difficulty);
  if (idx <= 0) return "usoara";
  return DIFFICULTY_ORDER[idx - 1];
}

function promoteDifficulty(difficulty, streakForDifficulty) {
  if (difficulty === "usoara" && streakForDifficulty >= 2) return "medie";
  if (difficulty === "medie" && streakForDifficulty >= 2) return "grea";
  return difficulty;
}

function getQuestionSkillElo(difficulty) {
  if (difficulty === "usoara") return 1000;
  if (difficulty === "medie") return 1300;
  return 1600;
}

function parseRenderableGraphData(rawGraphData) {
  if (!rawGraphData) return null;

  let graphData = rawGraphData;
  if (typeof graphData === "string") {
    try {
      graphData = JSON.parse(graphData);
    } catch (_error) {
      return null;
    }
  }

  const nodes = Array.isArray(graphData?.nodes) ? graphData.nodes : [];
  const edges = Array.isArray(graphData?.edges) ? graphData.edges : [];
  if (nodes.length === 0 || edges.length === 0) return null;

  return { nodes, edges };
}

function readPersistedQuizState() {
  try {
    const raw = localStorage.getItem(QUIZ_STATE_KEY);
    if (!raw) return null;
    const parsed = JSON.parse(raw);
    return parsed && typeof parsed === "object" ? parsed : null;
  } catch (_error) {
    return null;
  }
}

function calculateElo(userElo, isCorrect, questionDifficulty, consecutiveCorrectOnDifficulty) {
  const questionSkillElo = getQuestionSkillElo(questionDifficulty);
  const expected = 1 / (1 + Math.pow(10, (questionSkillElo - userElo) / 400));
  const score = isCorrect ? 1 : 0;

  let k = 24;
  if (questionDifficulty === "medie") k = 32;
  if (questionDifficulty === "grea") k = 40;

  if (isCorrect) {
    const streakBonus = Math.min(16, Math.max(0, consecutiveCorrectOnDifficulty - 1) * 4);
    k += streakBonus;
  }

  let delta = k * (score - expected);
  if (!isCorrect) {
    delta *= 1.1;
  }

  delta = Math.max(-48, Math.min(52, delta));

  return {
    delta: Math.round(delta * 10) / 10,
    newElo: Math.max(0, Math.round((userElo + delta) * 10) / 10),
  };
}

export default function GrInfoQuiz() {
  const [searchParams, setSearchParams] = useSearchParams();
  const persistedQuizState = readPersistedQuizState();
  const initialCategory = searchParams.get("category") || persistedQuizState?.selectedCategory || "all";
  const initialDifficultyChoice = searchParams.get("difficulty") || persistedQuizState?.selectedDifficultyChoice || "toate";
  const hasPersistedDraft = !!(persistedQuizState && Array.isArray(persistedQuizState.history) && persistedQuizState.history.length > 0);

  const [questionsPool, setQuestionsPool] = useState([]);
  const [selectedCategory, setSelectedCategory] = useState(initialCategory);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  const [userElo, setUserElo] = useState(INITIAL_ELO);
  const [usedQuestionIds, setUsedQuestionIds] = useState(new Set());
  const [currentQuestion, setCurrentQuestion] = useState(persistedQuizState?.currentQuestion || null);
  const [selectedOption, setSelectedOption] = useState(null);
  const [showExplanation, setShowExplanation] = useState(false);
  const [feedback, setFeedback] = useState("");
  const [history, setHistory] = useState([]);
  const [securityWarning, setSecurityWarning] = useState("");
  const [sessionSaved, setSessionSaved] = useState(false);
  const [sessionId, setSessionId] = useState(null);
  const [startingElo, setStartingElo] = useState(INITIAL_ELO);
  const [targetDifficulty, setTargetDifficulty] = useState("usoara");
  const [difficultyStreaks, setDifficultyStreaks] = useState({ usoara: 0, medie: 0, grea: 0 });
  const [selectedDifficultyChoice, setSelectedDifficultyChoice] = useState(initialDifficultyChoice);

  const quizContainerRef = useRef(null);
  const answerLockRef = useRef(false);
  const graphRef = useRef(null);
  const networkRef = useRef(null);
  const renderableGraphData = useMemo(
    () => parseRenderableGraphData(currentQuestion?.graphData),
    [currentQuestion]
  );

  function resetQuiz(nextElo) {
    const initialDifficulty = getInitialDifficultyByElo(nextElo);
    setStartingElo(nextElo);
    setCurrentQuestion(pickNextQuestion(questionsPool, new Set(), initialDifficulty));
    setUserElo(nextElo);
    setTargetDifficulty(initialDifficulty);
    setDifficultyStreaks({ usoara: 0, medie: 0, grea: 0 });
    setUsedQuestionIds(new Set());
    setSelectedOption(null);
    setShowExplanation(false);
    setFeedback("");
    setHistory([]);
    setSessionSaved(false);
    setSessionId(null);
    answerLockRef.current = false;
    try {
      localStorage.removeItem(QUIZ_STATE_KEY);
    } catch (_error) {}
  }

  const answeredCount = history.length;
  const isFinished = answeredCount >= 10 || !currentQuestion;
  const sessionDifficultyLabel =
    selectedDifficultyChoice && selectedDifficultyChoice !== "toate"
      ? selectedDifficultyChoice
      : history[0]?.difficulty || currentQuestion?.dificultate || targetDifficulty || "all";

  useEffect(() => {
    const categoryFromUrl = searchParams.get("category");
    if (categoryFromUrl && categoryFromUrl !== selectedCategory) {
      setSelectedCategory(categoryFromUrl);
    }
  }, [searchParams, selectedCategory]);

  useEffect(() => {
    async function loadQuestions() {
      setLoading(true);
      setError("");
      try {
        const query = new URLSearchParams({ limit: "100" });
        query.set("limit", "100");
        if (selectedCategory !== "all") {
          query.set("category", selectedCategory);
        }
        if (selectedDifficultyChoice && selectedDifficultyChoice !== "toate") {
          query.set("difficulty", selectedDifficultyChoice);
        }

        const token = localStorage.getItem("token");
        let currentProfileElo = INITIAL_ELO;
        if (token) {
          try {
            const profileResponse = await fetch("/api/grinfo/profile", {
              headers: { Authorization: `Bearer ${token}` },
            });
            if (profileResponse.ok) {
              const profileData = await profileResponse.json();
              if (typeof profileData.currentElo === "number") {
                currentProfileElo = profileData.currentElo;
              }
            }
          } catch (profileError) {
            console.error("Nu am putut încărca ELO-ul curent:", profileError);
          }
        }

        setStartingElo(currentProfileElo);
        const initialDifficulty =
          selectedDifficultyChoice && selectedDifficultyChoice !== "toate"
            ? selectedDifficultyChoice
            : getInitialDifficultyByElo(currentProfileElo);

        const response = await fetch(`/api/grinfo/questions?${query.toString()}`);
        if (!response.ok) {
          throw new Error("Nu s-au putut încărca întrebările din baza de date.");
        }

        const data = await response.json();
        const loaded = dedupeQuestions(Array.isArray(data.questions) ? data.questions : []);
        const randomized = shuffleQuestions(loaded).map(shuffleQuestionOptions);
        setQuestionsPool(randomized);
        const saved = readPersistedQuizState();
        if (saved && Array.isArray(saved.history) && saved.history.length > 0) {
          const savedUsedIds = new Set(saved.usedQuestionIds || []);
          const restoredQuestion = saved.currentQuestion
            || randomized.find((q) => q.id === saved.currentQuestionId)
            || pickNextQuestion(randomized, savedUsedIds, saved.targetDifficulty || initialDifficulty)
            || pickNextQuestion(randomized, new Set(), initialDifficulty);

          setSelectedCategory(saved.selectedCategory || selectedCategory);
          setSelectedDifficultyChoice(saved.selectedDifficultyChoice || "toate");
          setStartingElo(typeof saved.startingElo === "number" ? saved.startingElo : currentProfileElo);
          setUserElo(typeof saved.userElo === "number" ? saved.userElo : currentProfileElo);
          setTargetDifficulty(saved.targetDifficulty || initialDifficulty);
          setDifficultyStreaks(saved.difficultyStreaks || { usoara: 0, medie: 0, grea: 0 });
          setUsedQuestionIds(savedUsedIds);
          setSelectedOption(typeof saved.selectedOption === "number" ? saved.selectedOption : null);
          setShowExplanation(!!saved.showExplanation);
          setFeedback(saved.feedback || "");
          setHistory(saved.history || []);
          setCurrentQuestion(restoredQuestion);
          setSessionId(saved.sessionId || null);
          setSessionSaved(!!saved.sessionSaved);
          answerLockRef.current = !!saved.showExplanation;
          return;
        }

        const first = pickNextQuestion(randomized, new Set(), initialDifficulty);
        setCurrentQuestion(first);
        setUserElo(currentProfileElo);
        setTargetDifficulty(initialDifficulty);
        setDifficultyStreaks({ usoara: 0, medie: 0, grea: 0 });
        setUsedQuestionIds(new Set());
        setSelectedOption(null);
        setShowExplanation(false);
        setFeedback("");
        setHistory([]);
        setSessionId(null);
        setSessionSaved(false);
        answerLockRef.current = false;
        setCurrentQuestion(first);

      } catch (err) {
        setError(err.message || "Eroare de încărcare.");
      } finally {
        setLoading(false);
      }
    }

    loadQuestions();
  }, [selectedCategory, selectedDifficultyChoice]);

  useEffect(() => {
    function penalize(reason) {
      setUserElo((prev) => {
        const next = Math.max(0, prev - 10);
        fetch("/api/grinfo/incident", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            reason,
            eloPenalty: 10,
            category: selectedCategory,
            currentElo: next,
          }),
        }).catch(() => {});
        return next;
      });
      setSecurityWarning(`${reason}: penalizare -10 ELO`);
      window.setTimeout(() => setSecurityWarning(""), 2200);
    }

    function onVisibilityChange() {
      if (document.hidden) {
        penalize("Ai părăsit tab-ul");
      }
    }

    function onContextMenu(event) {
      if (quizContainerRef.current && quizContainerRef.current.contains(event.target)) {
        event.preventDefault();
        penalize("Click dreapta blocat");
      }
    }

    function onCopy(event) {
      if (quizContainerRef.current && quizContainerRef.current.contains(event.target)) {
        event.preventDefault();
        penalize("Copiere blocată");
      }
    }

    document.addEventListener("visibilitychange", onVisibilityChange);
    document.addEventListener("contextmenu", onContextMenu);
    document.addEventListener("copy", onCopy);

    return () => {
      document.removeEventListener("visibilitychange", onVisibilityChange);
      document.removeEventListener("contextmenu", onContextMenu);
      document.removeEventListener("copy", onCopy);
    };
  }, [selectedCategory]);

  useEffect(() => {
    if (!isFinished || sessionSaved || answeredCount === 0) return;

    const token = localStorage.getItem("token");
    if (!token) {
      setSessionSaved(true);
      return;
    }

    const correctAnswers = history.filter((h) => h.isCorrect).length;
    const endpoint = "/api/grinfo/session";
    const payload = sessionId
      ? {
          sessionId,
          category: selectedCategory,
          difficulty: sessionDifficultyLabel,
          finalElo: userElo,
          correctAnswers,
          totalQuestions: answeredCount,
        }
      : {
          category: selectedCategory,
          difficulty: sessionDifficultyLabel,
          initialElo: startingElo,
          totalQuestions: answeredCount,
          correctAnswers,
          finalElo: userElo,
        };

    fetch(endpoint, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify(payload),
    })
      .then(() => {
        setSessionSaved(true);
        setStartingElo(userElo);
        window.dispatchEvent(new Event("grinfo:profile-updated"));
      })
      .catch(() => setSessionSaved(true));
  }, [isFinished, sessionSaved, answeredCount, history, selectedCategory, userElo, sessionId, sessionDifficultyLabel, startingElo]);

  // Persist in-progress quiz state to localStorage so refresh/accidental leave restores
  useEffect(() => {
    try {
      if (loading) return;
      const payload = {
        selectedCategory,
        selectedDifficultyChoice,
        userElo,
        startingElo,
        history,
        currentQuestion: currentQuestion || null,
        currentQuestionId: currentQuestion?.id || null,
        usedQuestionIds: Array.from(usedQuestionIds || []),
        targetDifficulty,
        difficultyStreaks,
        sessionId,
        sessionSaved,
        selectedOption,
        showExplanation,
        feedback,
        timestamp: Date.now(),
      };
      localStorage.setItem(QUIZ_STATE_KEY, JSON.stringify(payload));
    } catch (e) {
      // ignore localStorage errors
    }
  }, [userElo, startingElo, history, currentQuestion, usedQuestionIds, targetDifficulty, difficultyStreaks, sessionId, sessionSaved, selectedCategory, selectedDifficultyChoice, selectedOption, showExplanation, feedback]);

  // Attempt to send final session/state on unload (best-effort)
  useEffect(() => {
    const saveOnExit = () => {
      try {
        // local save is already done by other effect; try to notify server if session exists
        if (!sessionId) return;
        const payload = {
          sessionId,
          category: selectedCategory,
          difficulty: sessionDifficultyLabel,
          finalElo: userElo,
          correctAnswers: history.filter((h) => h.isCorrect).length,
          totalQuestions: Math.max(1, history.length),
        };
        if (navigator && navigator.sendBeacon) {
          const blob = new Blob([JSON.stringify(payload)], { type: 'application/json' });
          navigator.sendBeacon('/api/grinfo/session', blob);
        } else {
          fetch('/api/grinfo/session', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload),
            keepalive: true,
          }).catch(() => {});
        }
      } catch (_) {}
    };

    window.addEventListener('beforeunload', saveOnExit);
    return () => window.removeEventListener('beforeunload', saveOnExit);
  }, [sessionId, selectedCategory, sessionDifficultyLabel, userElo, history]);

  // Render graph visualization on the right when question has graphData
  useEffect(() => {
    let cancelled = false;
    async function renderGraph() {
      if (!graphRef.current) return;
      if (!renderableGraphData) {
        if (networkRef.current) {
          try {
            networkRef.current.destroy();
          } catch (_) {}
          networkRef.current = null;
        }
        graphRef.current.innerHTML = "";
        return;
      }

      graphRef.current.innerHTML = "";
      const nodes = renderableGraphData.nodes.map((n) => ({ id: n.id, label: String(n.id) }));
      const edges = renderableGraphData.edges.map((e) => ({
        from: e.from ?? e.source,
        to: e.to ?? e.target,
        label: e.label ?? e.weight ?? e.cost ?? null,
        arrows:
          e.arrows ??
          (e.type === "->" ? "to" : selectedCategory === "orientate" ? "to" : undefined),
      }));

      try {
        const vis = await import("vis-network/standalone");
        if (cancelled) return;
        if (networkRef.current) {
          try {
            networkRef.current.destroy();
          } catch (_) {}
        }
        const data = { nodes: new vis.DataSet(nodes), edges: new vis.DataSet(edges) };
        const options = {
          autoResize: true,
          edges: {
            color: { color: "#97c2fc" },
            smooth: { type: "continuous", forceDirection: "none" },
            font: { align: "top" },
          },
          nodes: { color: { background: "#e7f0ff", border: "#3a74d6" } },
          layout: { improvedLayout: true },
          physics: {
            enabled: true,
            barnesHut: { gravitationalConstant: -1800, springLength: 150, springConstant: 0.03 },
            stabilization: { iterations: 100, fit: true },
          },
          interaction: { hover: true, dragNodes: false, zoomView: true },
        };
        networkRef.current = new vis.Network(graphRef.current, data, options);
      } catch (err) {
        // dynamic import failed or render failed
        console.error("Could not render graph preview:", err);
      }
    }

    renderGraph();
    return () => {
      cancelled = true;
    };
  }, [renderableGraphData]);

  async function onSelectAnswer(index) {
    if (!currentQuestion || showExplanation || answerLockRef.current) return;

    answerLockRef.current = true;

    const token = localStorage.getItem("token");
    let activeSessionId = sessionId;
    if (!activeSessionId && token) {
      try {
        const startResponse = await fetch("/api/grinfo/session", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({
            category: selectedCategory,
            difficulty: sessionDifficultyLabel,
            initialElo: startingElo,
            totalQuestions: 10,
          }),
        });

        if (startResponse.ok) {
          const startData = await startResponse.json();
          activeSessionId = startData.sessionId || null;
          setSessionId(activeSessionId);
        }
      } catch (startError) {
        console.error("Nu am putut porni sesiunea GrInfo:", startError);
      }
    }

    const isCorrect = index === currentQuestion.raspunsCorect;
    const currentDifficulty = currentQuestion.dificultate || targetDifficulty;
    const currentDifficultyStreak = difficultyStreaks[currentDifficulty] || 0;
    const nextDifficultyStreak = isCorrect ? currentDifficultyStreak + 1 : 0;
    const eloResult = calculateElo(userElo, isCorrect, currentDifficulty, nextDifficultyStreak);

    const updatedStreaks = { ...difficultyStreaks };
    updatedStreaks[currentDifficulty] = nextDifficultyStreak;

    const isDifficultyFixed = selectedDifficultyChoice && selectedDifficultyChoice !== "toate";
    let nextTargetDifficulty = targetDifficulty;
    if (!isDifficultyFixed) {
      if (isCorrect) {
        nextTargetDifficulty = promoteDifficulty(currentDifficulty, nextDifficultyStreak);
        if (nextTargetDifficulty !== currentDifficulty) {
          updatedStreaks[currentDifficulty] = 0;
        }
      } else {
        nextTargetDifficulty = lowerDifficulty(currentDifficulty);
      }
    } else {
      // keep target difficulty fixed
      nextTargetDifficulty = targetDifficulty;
    }

    const logEntry = {
      questionId: currentQuestion.id,
      isCorrect,
      eloBefore: userElo,
      eloAfter: eloResult.newElo,
      delta: eloResult.delta,
      category: currentQuestion.category,
      difficulty: currentQuestion.dificultate,
      questionText: currentQuestion.enunt,
      selectedOptionText: currentQuestion.optiuni[index],
      correctAnswerText: currentQuestion.optiuni[currentQuestion.raspunsCorect],
    };

    setSelectedOption(index);
    setShowExplanation(true);
    setDifficultyStreaks(updatedStreaks);
    setTargetDifficulty(nextTargetDifficulty);
    setUserElo(eloResult.newElo);
    setHistory((prev) => [...prev, logEntry]);
    setUsedQuestionIds((prev) => {
      const next = new Set(prev);
      next.add(currentQuestion.id);
      return next;
    });

    const answerText = currentQuestion.optiuni[currentQuestion.raspunsCorect];
    if (isCorrect) {
      setFeedback(`Corect. ELO +${Math.max(0, eloResult.delta).toFixed(1)}`);
    } else {
      setFeedback(`Greșit. Răspuns corect: ${answerText}. ELO ${eloResult.delta.toFixed(1)}`);
    }
  }

  function onContinue() {
    if (!currentQuestion) return;

    const used = new Set(usedQuestionIds);
    used.add(currentQuestion.id);

    const next = pickNextQuestion(questionsPool, used, targetDifficulty);
    setCurrentQuestion(next);
    setSelectedOption(null);
    setShowExplanation(false);
    setFeedback("");
    setUsedQuestionIds(used);
    answerLockRef.current = false;
  }

  const finalStats = useMemo(() => {
    const correct = history.filter((h) => h.isCorrect).length;
    const wrongAnswers = history.filter((h) => !h.isCorrect);
    return {
      correct,
      total: history.length,
      percentage: history.length ? ((correct / history.length) * 100).toFixed(1) : "0.0",
      wrongAnswers,
    };
  }, [history]);

  return (
    <div className="grinfo-page">
      <div className="grinfo-header">
        <h1>GrInfo Quiz</h1>
        <p>10 întrebări de grafuri, secvențial, cu ELO adaptiv.</p>
      </div>

      <div className="grinfo-toolbar">
        <label>
          Categorie
          <select
            value={selectedCategory}
            onChange={(e) => {
              const nextCategory = e.target.value;
              setSelectedCategory(nextCategory);
              const params = new URLSearchParams(searchParams.toString());
              if (nextCategory === "all") {
                params.delete("category");
              } else {
                params.set("category", nextCategory);
              }
              setSearchParams(params, { replace: true });
            }}
            disabled={answeredCount > 0}
          >
            <option value="all">Toate</option>
            <option value="orientate">Grafuri orientate</option>
            <option value="neorientate">Grafuri neorientate</option>
          </select>
        </label>
        <label style={{ marginLeft: 12 }}>
          Dificultate
          <select
            value={selectedDifficultyChoice}
            onChange={(e) => setSelectedDifficultyChoice(e.target.value)}
            disabled={answeredCount > 0}
          >
            <option value="toate">Toate</option>
            <option value="usoara">Ușoară</option>
            <option value="medie">Medie</option>
            <option value="grea">Grea</option>
          </select>
        </label>
        <div className="grinfo-chip">ELO: {userElo.toFixed(1)}</div>
        <div className="grinfo-chip">
          {selectedDifficultyChoice === "toate"
            ? `Dificultate țintă: ${targetDifficulty}`
            : `Dificultate FIXATA: ${selectedDifficultyChoice}`}
        </div>
        <div className="grinfo-chip">Întrebări: {Math.min(answeredCount + (showExplanation ? 0 : 1), 10)}/10</div>
      </div>

      {securityWarning && <div className="grinfo-alert">{securityWarning}</div>}

      {loading && !hasPersistedDraft && <div className="grinfo-card">Se încarcă întrebările din baza de date...</div>}
      {error && !loading && <div className="grinfo-card grinfo-error">{error}</div>}

      {(!loading || hasPersistedDraft) && !error && isFinished && (
        <div className="grinfo-card">
          <h2>Quiz finalizat</h2>
          <p>Răspunsuri corecte: {finalStats.correct}/{finalStats.total}</p>
          <p>Procent: {finalStats.percentage}%</p>
          <p>ELO final: {userElo.toFixed(1)}</p>

          {finalStats.wrongAnswers.length > 0 && (
            <div className="grinfo-wrong-list">
              <h3>Răspunsuri greșite</h3>
              {finalStats.wrongAnswers.map((entry, idx) => (
                <div key={`${entry.questionId}-${idx}`} className="grinfo-wrong-item">
                  <div className="grinfo-wrong-question">{entry.questionText}</div>
                  <div className="grinfo-wrong-answer">Ai ales: {entry.selectedOptionText}</div>
                  <div className="grinfo-wrong-answer grinfo-wrong-correct">Corect era: {entry.correctAnswerText}</div>
                </div>
              ))}
            </div>
          )}

          <button
            className="grinfo-btn"
            onClick={() => {
              resetQuiz(userElo);
            }}
          >
            Reîncarcă quiz
          </button>
        </div>
      )}

      {(!loading || hasPersistedDraft) && !error && !isFinished && currentQuestion && (
        <div className={`grinfo-card ${renderableGraphData ? "has-graph" : ""}`} ref={quizContainerRef}>
          <div className="grinfo-question-area">
            <div className="grinfo-meta">
              <span>{currentQuestion.category}</span>
              <span>{currentQuestion.dificultate}</span>
              <span>ELO întrebare: {currentQuestion.eloRating}</span>
            </div>

            <h2>{currentQuestion.enunt}</h2>

            <div className="grinfo-options">
              {currentQuestion.optiuni.map((option, idx) => {
                const isCorrect = showExplanation && idx === currentQuestion.raspunsCorect;
                const isWrongChoice = showExplanation && selectedOption === idx && idx !== currentQuestion.raspunsCorect;
                return (
                  <button
                    key={idx}
                    className={`grinfo-option ${isCorrect ? "is-correct" : ""} ${isWrongChoice ? "is-wrong" : ""}`}
                    onClick={() => onSelectAnswer(idx)}
                    disabled={showExplanation}
                  >
                    {String.fromCharCode(65 + idx)}. {option}
                  </button>
                );
              })}
            </div>

            {showExplanation && (
              <div className="grinfo-explanation">
                <p>{feedback}</p>
                <h3>Explicație</h3>
                <p>{currentQuestion.explicatieRaspuns}</p>
                <button className="grinfo-btn" onClick={onContinue}>
                  Continuă
                </button>
              </div>
            )}
          </div>

          {renderableGraphData && (
            <div className="grinfo-graph-panel">
              <div ref={graphRef} id="grinfo-graph-canvas" />
            </div>
          )}
        </div>
      )}
    </div>
  );
}

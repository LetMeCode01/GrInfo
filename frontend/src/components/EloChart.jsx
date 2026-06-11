import React, { useMemo, useState } from "react";
import { Line } from "react-chartjs-2";
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Tooltip,
  Legend,
} from "chart.js";

ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, Tooltip, Legend);

export default function EloChart({ matchesData = [] }) {
  const [filter, setFilter] = useState("all");

  const normalizedData = useMemo(() => {
    if (!Array.isArray(matchesData)) return [];

    return matchesData
      .filter((item) => item && Number.isFinite(Number(item.elo)) && item.date)
      .map((item) => ({
        date: String(item.date),
        elo: Number(item.elo),
      }));
  }, [matchesData]);

  // Keep requested behavior: component reverses natural DB order for chart left->right.
  const chronologicalData = useMemo(() => [...normalizedData].reverse(), [normalizedData]);

  const filteredData = useMemo(() => {
    if (filter === "last10") return chronologicalData.slice(-10);
    return chronologicalData;
  }, [filter, chronologicalData]);

  const data = useMemo(
    () => ({
      labels: filteredData.map((match) => match.date),
      datasets: [
        {
          label: "Evolutie ELO",
          data: filteredData.map((match) => match.elo),
          borderColor: "#4CAF50",
          backgroundColor: "rgba(76, 175, 80, 0.1)",
          borderWidth: 3,
          tension: 0.2,
          pointRadius: 0,
          pointHoverRadius: 6,
          pointHoverBackgroundColor: "#4CAF50",
          pointHoverBorderColor: "#fff",
          pointHoverBorderWidth: 2,
        },
      ],
    }),
    [filteredData]
  );

  const options = useMemo(
    () => ({
      responsive: true,
      maintainAspectRatio: false,
      plugins: {
        legend: { display: false },
        tooltip: {
          enabled: true,
          mode: "index",
          intersect: false,
          backgroundColor: "#1e1e2f",
          titleColor: "#fff",
          bodyColor: "#4CAF50",
          bodyFont: { size: 14, weight: "bold" },
          callbacks: {
            label(context) {
              return ` ELO: ${Number(context.parsed.y || 0).toFixed(1)}`;
            },
          },
        },
      },
      scales: {
        x: {
          grid: { display: false },
          ticks: { maxTicksLimit: 8 },
        },
        y: {
          grid: { color: "rgba(200, 200, 200, 0.1)" },
          ticks: {
            callback(value) {
              return `${value} pts`;
            },
          },
        },
      },
    }),
    []
  );

  if (!normalizedData.length) return null;

  return (
    <div
      style={{
        background: "#fff",
        padding: "20px",
        borderRadius: "12px",
        boxShadow: "0 4px 6px rgba(0,0,0,0.05)",
        marginTop: "16px",
      }}
    >
      <div
        style={{
          display: "flex",
          justifyContent: "space-between",
          alignItems: "center",
          marginBottom: "15px",
          gap: "10px",
          flexWrap: "wrap",
        }}
      >
        <h3 style={{ margin: 0, color: "#333" }}>Progresul tau ELO</h3>

        <div>
          <button
            onClick={() => setFilter("last10")}
            style={{
              marginRight: "5px",
              padding: "5px 10px",
              background: filter === "last10" ? "#4CAF50" : "#eee",
              color: filter === "last10" ? "#fff" : "#333",
              border: "none",
              borderRadius: "4px",
              cursor: "pointer",
            }}
          >
            Ultimele 10
          </button>
          <button
            onClick={() => setFilter("all")}
            style={{
              padding: "5px 10px",
              background: filter === "all" ? "#4CAF50" : "#eee",
              color: filter === "all" ? "#fff" : "#333",
              border: "none",
              borderRadius: "4px",
              cursor: "pointer",
            }}
          >
            Tot istoricul ({normalizedData.length})
          </button>
        </div>
      </div>

      <div style={{ height: "250px", position: "relative" }}>
        <Line data={data} options={options} />
      </div>
    </div>
  );
}

import React, { useCallback, useEffect, useImperativeHandle, useRef, useState } from "react";

const DEFAULT_TOOL = "draw";
const DEFAULT_COLOR = "#111827";
const DEFAULT_LINE_WIDTH = 3;

const Scratchpad = React.forwardRef(function Scratchpad({ storageKey }, ref) {
  const canvasRef = useRef(null);
  const containerRef = useRef(null);
  const isDrawingRef = useRef(false);

  const [tool, setTool] = useState(DEFAULT_TOOL);
  const [color, setColor] = useState(DEFAULT_COLOR);
  const [lineWidth, setLineWidth] = useState(DEFAULT_LINE_WIDTH);
  const [textDraft, setTextDraft] = useState(null);

  useImperativeHandle(ref, () => ({
    reset: () => {
      clearCanvas();
      if (storageKey) {
        try {
          localStorage.removeItem(storageKey);
        } catch (_error) {}
      }
    },
  }), [storageKey]);

  const persistCanvas = useCallback(() => {
    if (!storageKey || !canvasRef.current) return;
    try {
      const dataUrl = canvasRef.current.toDataURL("image/png");
      const payload = {
        dataUrl,
        tool,
        color,
        lineWidth,
      };
      localStorage.setItem(storageKey, JSON.stringify(payload));
    } catch (_error) {
      // ignore localStorage/canvas persistence errors
    }
  }, [storageKey, tool, color, lineWidth]);

  useEffect(() => {
    const canvas = canvasRef.current;
    if (!canvas) return;

    const ctx = canvas.getContext("2d");
    if (!ctx) return;

    ctx.lineCap = "round";
    ctx.lineJoin = "round";

    if (!storageKey) return;

    try {
      const raw = localStorage.getItem(storageKey);
      if (!raw) return;
      const parsed = JSON.parse(raw);

      if (parsed && typeof parsed === "object") {
        if (parsed.tool === "draw" || parsed.tool === "text") setTool(parsed.tool);
        if (typeof parsed.color === "string") setColor(parsed.color);
        if (typeof parsed.lineWidth === "number") setLineWidth(parsed.lineWidth);

        if (typeof parsed.dataUrl === "string" && parsed.dataUrl) {
          const image = new Image();
          image.onload = () => {
            ctx.clearRect(0, 0, canvas.width, canvas.height);
            ctx.drawImage(image, 0, 0, canvas.width, canvas.height);
          };
          image.src = parsed.dataUrl;
        }
      }
    } catch (_error) {
      // ignore invalid persisted payload
    }
  }, [storageKey]);

  useEffect(() => {
    const canvas = canvasRef.current;
    if (!canvas) return;

    const ctx = canvas.getContext("2d");
    if (!ctx) return;

    function getPoint(event) {
      const rect = canvas.getBoundingClientRect();
      return {
        x: event.clientX - rect.left,
        y: event.clientY - rect.top,
      };
    }

    function startDrawing(event) {
      if (tool !== "draw") return;
      isDrawingRef.current = true;
      const point = getPoint(event);
      ctx.beginPath();
      ctx.moveTo(point.x, point.y);
    }

    function draw(event) {
      if (!isDrawingRef.current || tool !== "draw") return;
      const point = getPoint(event);
      ctx.strokeStyle = color;
      ctx.lineWidth = lineWidth;
      ctx.lineTo(point.x, point.y);
      ctx.stroke();
    }

    function stopDrawing() {
      if (!isDrawingRef.current) return;
      isDrawingRef.current = false;
      ctx.beginPath();
      persistCanvas();
    }

    function handleCanvasClick(event) {
      if (tool !== "text") return;
      const point = getPoint(event);
      setTextDraft({ x: point.x, y: point.y, value: "" });
    }

    canvas.addEventListener("mousedown", startDrawing);
    canvas.addEventListener("mousemove", draw);
    canvas.addEventListener("mouseup", stopDrawing);
    canvas.addEventListener("mouseleave", stopDrawing);
    canvas.addEventListener("click", handleCanvasClick);

    return () => {
      canvas.removeEventListener("mousedown", startDrawing);
      canvas.removeEventListener("mousemove", draw);
      canvas.removeEventListener("mouseup", stopDrawing);
      canvas.removeEventListener("mouseleave", stopDrawing);
      canvas.removeEventListener("click", handleCanvasClick);
    };
  }, [tool, color, lineWidth, persistCanvas]);

  function clearCanvas() {
    const canvas = canvasRef.current;
    const ctx = canvas?.getContext("2d");
    if (!canvas || !ctx) return;
    ctx.clearRect(0, 0, canvas.width, canvas.height);
    setTextDraft(null);
    persistCanvas();
  }

  function commitText() {
    if (!textDraft || !textDraft.value.trim()) {
      setTextDraft(null);
      return;
    }

    const canvas = canvasRef.current;
    const ctx = canvas?.getContext("2d");
    if (!ctx) return;

    ctx.fillStyle = color;
    ctx.font = "16px Arial";
    ctx.fillText(textDraft.value, textDraft.x, textDraft.y + 12);
    setTextDraft(null);
    persistCanvas();
  }

  return (
    <div className="grinfo-scratchpad" ref={containerRef}>
      <div className="grinfo-scratchpad-toolbar">
        <button
          type="button"
          className={`grinfo-scratchpad-btn ${tool === "draw" ? "is-active" : ""}`}
          onClick={() => setTool("draw")}
        >
          Pensula
        </button>
        <button
          type="button"
          className={`grinfo-scratchpad-btn ${tool === "text" ? "is-active" : ""}`}
          onClick={() => setTool("text")}
        >
          Text
        </button>
        <label className="grinfo-scratchpad-label">
          Culoare
          <input
            type="color"
            value={color}
            onChange={(event) => setColor(event.target.value)}
          />
        </label>
        {tool === "draw" && (
          <label className="grinfo-scratchpad-label">
            Grosime
            <input
              type="range"
              min="1"
              max="14"
              value={lineWidth}
              onChange={(event) => setLineWidth(Number(event.target.value))}
            />
          </label>
        )}
        <button type="button" className="grinfo-scratchpad-btn clear" onClick={clearCanvas}>
          Sterge
        </button>
      </div>

      <div className="grinfo-scratchpad-hint">
        {tool === "draw"
          ? "Desen: tine click apasat si misca mouse-ul."
          : "Text: click pe canvas, scrie si Enter."}
      </div>

      <div className="grinfo-scratchpad-canvas-wrap">
        <canvas
          ref={canvasRef}
          width={320}
          height={220}
          className={`grinfo-scratchpad-canvas ${tool === "draw" ? "cursor-draw" : "cursor-text"}`}
        />

        {textDraft && (
          <input
            autoFocus
            type="text"
            className="grinfo-scratchpad-input"
            value={textDraft.value}
            style={{ left: textDraft.x, top: textDraft.y }}
            onChange={(event) => setTextDraft((prev) => (prev ? { ...prev, value: event.target.value } : prev))}
            onKeyDown={(event) => {
              if (event.key === "Enter") {
                event.preventDefault();
                commitText();
              }
              if (event.key === "Escape") {
                event.preventDefault();
                setTextDraft(null);
              }
            }}
            onBlur={commitText}
          />
        )}
      </div>
    </div>
  );
});

export default Scratchpad;
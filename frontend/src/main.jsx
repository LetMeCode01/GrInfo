import React from "react";
import { createRoot } from "react-dom/client";
import App from "./App";
import "./assets/main.css";
import "./assets/enhanced-styles.css";
import "./assets/header-styles.css";
import "./assets/utility-styles.css";
import "./styles/roadmap.css";

const container = document.getElementById("root");
const root = createRoot(container);
root.render(<App />);

// Aici avem codul pentru a porni aplicația React în containerul HTML cu id "root".

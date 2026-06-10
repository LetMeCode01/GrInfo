import React from "react";
import { Link } from "react-router-dom";

export default function Footer() {
  return (
    <footer className="site-footer">
      <div className="footer-container">
        <div className="footer-column">
          <h4>Adresă</h4>
          <div>PBTower, București, România</div>
        </div>

        <div className="footer-column">
          <h4>Contact</h4>
          <div className="footer-contact">
            Telefon: <a href="tel:+40712345678">+40 762 104 169</a>
          </div>
          <div className="footer-contact">
            Email:{" "}
            <a href="mailto:contact@grinfo.example">
              mihaisec22@gmail.com
            </a>
          </div>
        </div>

        <div className="footer-column">
          <h4>Social</h4>
          <div className="footer-social">
            <a
              href="https://instagram.com/mihai_sima_"
              target="_blank"
              rel="noreferrer"
              className="social-link"
            >
              Instagram
            </a>
            <span className="dot">•</span>
            <a
              href="https://twitter.com/Let_Me_Tweeeet"
              target="_blank"
              rel="noreferrer"
              className="social-link"
            >
              Twitter
            </a>
            <span className="dot">•</span>
            <a
              href="https://facebook.com/mihai.sima.12"
              target="_blank"
              rel="noreferrer"
              className="social-link"
            >
              Facebook
            </a>
          </div>
        </div>
      </div>

      <div className="footer-bottom">
        <div className="footer-links">
          <Link to="/terms">Termeni și condiții</Link>
          <Link to="/">Harta Site</Link>
        </div>
        <div className="copyright">
          © {new Date().getFullYear()} GrInfo. Toate drepturile rezervate.
        </div>
      </div>
    </footer>
  );
}

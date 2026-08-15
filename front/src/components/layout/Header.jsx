import { Link } from "react-router-dom";
import "./Header.css";

export default function Header() {
  return (
    <header className="header">
      <Link to="/" className="header__brand">
        <svg
          className="header__logo"
          viewBox="0 0 24 24"
          width="24"
          height="24"
          aria-hidden="true"
        >
          <path
            d="M12 3 3 10h2v10h5v-6h4v6h5V10h2z"
            fill="var(--color-primary)"
          />
        </svg>
        <span>SmartHouse</span>
      </Link>
    </header>
  );
}

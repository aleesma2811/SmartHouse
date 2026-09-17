import { Link } from "react-router-dom";
import "./InmuebleCard.css";

const TIPO_LABEL = {
  casa: "Casa",
  departamento: "Departamento",
};

export default function InmuebleCard({ inmueble, roomCount, onEdit, onDelete }) {
  return (
    <div className="inmueble-card">
      <Link to={`/inmueble/${inmueble.ID}`} className="inmueble-card__body">
        <span className={`inmueble-card__badge inmueble-card__badge--${inmueble.Tipo}`}>
          {TIPO_LABEL[inmueble.Tipo] ?? inmueble.Tipo}
        </span>
        <h3 className="inmueble-card__name">{inmueble.Nombre}</h3>
        <p className="inmueble-card__address">
          {inmueble.Direccion}
          {inmueble.Ciudad ? `, ${inmueble.Ciudad}` : ""}
        </p>
        <p className="inmueble-card__meta">
          {roomCount} habitación{roomCount === 1 ? "" : "es"}
        </p>
      </Link>

      <div className="inmueble-card__actions">
        <button
          type="button"
          className="inmueble-card__edit"
          title="Editar inmueble"
          aria-label="Editar inmueble"
          onClick={() => onEdit(inmueble)}
        >
          &#9998;
        </button>
        <button
          type="button"
          className="inmueble-card__delete"
          title="Eliminar inmueble"
          aria-label="Eliminar inmueble"
          onClick={() => onDelete(inmueble)}
        >
          &times;
        </button>
      </div>
    </div>
  );
}

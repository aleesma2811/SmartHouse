import { getServiceMeta } from "./serviceTheme";
import { IconPencil, IconTrash } from "../common/Icons";
import "./PlugList.css";

export default function PlugList({ plugs, onToggle, onEdit, onDelete }) {
  if (plugs.length === 0) {
    return <p className="plug-list__empty">Todavía no hay servicios en esta habitación.</p>;
  }

  return (
    <ul className="plug-list">
      {plugs.map((plug) => {
        const meta = getServiceMeta(plug.Tipo);
        const Icon = meta.Icon;
        return (
          <li key={plug.ID} className="plug-list__item">
            <span
              className="plug-list__icon"
              style={{ color: meta.color, background: meta.soft }}
            >
              <Icon width={17} height={17} />
            </span>

            <div className="plug-list__info">
              <span className="plug-list__name">{plug.Name}</span>
              <span className="plug-list__meta">
                {meta.label} · {plug.Consumo ?? 0} {meta.unit} ·{" "}
                {plug.On ? "Encendido" : "Apagado"}
              </span>
            </div>

            <button
              type="button"
              className={plug.On ? "plug-list__toggle plug-list__toggle--on" : "plug-list__toggle"}
              style={{ "--toggle-color": meta.color }}
              onClick={() => onToggle(plug)}
              title={plug.On ? "Apagar" : "Encender"}
              aria-pressed={plug.On}
            >
              <span className="plug-list__toggle-dot" />
            </button>

            <div className="plug-list__actions">
              <button
                type="button"
                className="icon-btn"
                title="Editar servicio"
                aria-label="Editar servicio"
                onClick={() => onEdit(plug)}
              >
                <IconPencil width={16} height={16} />
              </button>
              <button
                type="button"
                className="icon-btn icon-btn--danger"
                title="Eliminar servicio"
                aria-label="Eliminar servicio"
                onClick={() => onDelete(plug)}
              >
                <IconTrash width={16} height={16} />
              </button>
            </div>
          </li>
        );
      })}
    </ul>
  );
}

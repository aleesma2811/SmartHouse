import "./PlugList.css";

export default function PlugList({ plugs, onToggle, onEdit, onDelete }) {
  if (plugs.length === 0) {
    return <p className="plug-list__empty">Todavía no hay enchufes en esta habitación.</p>;
  }

  return (
    <ul className="plug-list">
      {plugs.map((plug) => (
        <li key={plug.ID} className="plug-list__item">
          <button
            type="button"
            className={plug.On ? "plug-list__toggle plug-list__toggle--on" : "plug-list__toggle"}
            onClick={() => onToggle(plug)}
            title={plug.On ? "Apagar" : "Encender"}
            aria-pressed={plug.On}
          >
            <span className="plug-list__toggle-dot" />
          </button>

          <div className="plug-list__info">
            <span className="plug-list__name">{plug.Name}</span>
            <span className="plug-list__meta">
              {plug.KwhConsump ?? 0} kWh · {plug.On ? "Encendido" : "Apagado"}
            </span>
          </div>

          <div className="plug-list__actions">
            <button type="button" className="btn btn-ghost" onClick={() => onEdit(plug)}>
              Editar
            </button>
            <button type="button" className="btn btn-danger" onClick={() => onDelete(plug)}>
              Eliminar
            </button>
          </div>
        </li>
      ))}
    </ul>
  );
}

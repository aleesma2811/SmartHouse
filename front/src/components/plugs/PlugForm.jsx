import { useState } from "react";
import { getServiceMeta, SERVICE_META } from "./serviceTheme";
import "../common/forms.css";
import "./PlugForm.css";

export default function PlugForm({ initialPlug, onSubmit, onCancel }) {
  const isEditing = Boolean(initialPlug);
  const [name, setName] = useState(initialPlug?.Name ?? "");
  const [tipo, setTipo] = useState(initialPlug?.Tipo ?? "luz");
  const [kwhConsump, setKwhConsump] = useState(initialPlug?.Consumo ?? "");
  const [on, setOn] = useState(initialPlug?.On ?? false);
  const [error, setError] = useState("");
  const [saving, setSaving] = useState(false);

  const meta = getServiceMeta(tipo);
  const Icon = meta.Icon;

  async function handleSubmit(e) {
    e.preventDefault();
    if (!name.trim()) {
      setError("El nombre es obligatorio");
      return;
    }
    setSaving(true);
    setError("");
    try {
      await onSubmit({ name: name.trim(), tipo, kwhConsump, on });
    } catch (err) {
      setError(err.message || "No se pudo guardar el servicio");
      setSaving(false);
    }
  }

  return (
    <form onSubmit={handleSubmit}>
      {error && <p className="form-error">{error}</p>}

      <div className="field">
        <label htmlFor="plug-name">Nombre</label>
        <input
          id="plug-name"
          type="text"
          value={name}
          onChange={(e) => setName(e.target.value)}
          placeholder="Ej. Lámpara de mesa"
          autoFocus
        />
      </div>

      <div className="field">
        <label htmlFor="plug-tipo">Tipo de servicio</label>
        <div className="plug-form__tipo-picker">
          {Object.entries(SERVICE_META).map(([value, option]) => {
            const OptionIcon = option.Icon;
            const active = value === tipo;
            return (
              <button
                key={value}
                type="button"
                className={
                  active
                    ? "plug-form__tipo-option plug-form__tipo-option--active"
                    : "plug-form__tipo-option"
                }
                style={{ "--tipo-color": option.color, "--tipo-soft": option.soft }}
                onClick={() => setTipo(value)}
                aria-pressed={active}
              >
                <OptionIcon width={18} height={18} />
                {option.label}
              </button>
            );
          })}
        </div>
      </div>

      <div className="field">
        <label htmlFor="plug-kwh">
          Consumo ({meta.unit})
        </label>
        <div className="plug-form__consumo">
          <Icon width={16} height={16} style={{ color: meta.color }} />
          <input
            id="plug-kwh"
            type="number"
            step="0.01"
            min="0"
            value={kwhConsump}
            onChange={(e) => setKwhConsump(e.target.value)}
            placeholder={`0.00 ${meta.unit}`}
          />
        </div>
      </div>

      <div className="field field--checkbox">
        <input
          id="plug-on"
          type="checkbox"
          checked={on}
          onChange={(e) => setOn(e.target.checked)}
        />
        <label htmlFor="plug-on">Encendido</label>
      </div>

      <div className="form-actions">
        <button type="button" className="btn btn-ghost" onClick={onCancel}>
          Cancelar
        </button>
        <button type="submit" className="btn btn-primary" disabled={saving}>
          {saving ? "Guardando..." : isEditing ? "Guardar cambios" : "Añadir servicio"}
        </button>
      </div>
    </form>
  );
}

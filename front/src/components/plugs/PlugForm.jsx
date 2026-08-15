import { useState } from "react";
import "../common/forms.css";

export default function PlugForm({ initialPlug, onSubmit, onCancel }) {
  const isEditing = Boolean(initialPlug);
  const [name, setName] = useState(initialPlug?.Name ?? "");
  const [kwhConsump, setKwhConsump] = useState(initialPlug?.KwhConsump ?? "");
  const [on, setOn] = useState(initialPlug?.On ?? false);
  const [error, setError] = useState("");
  const [saving, setSaving] = useState(false);

  async function handleSubmit(e) {
    e.preventDefault();
    if (!name.trim()) {
      setError("El nombre es obligatorio");
      return;
    }
    setSaving(true);
    setError("");
    try {
      await onSubmit({ name: name.trim(), kwhConsump, on });
    } catch (err) {
      setError(err.message || "No se pudo guardar el enchufe");
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
        <label htmlFor="plug-kwh">Consumo (kWh)</label>
        <input
          id="plug-kwh"
          type="number"
          step="0.01"
          min="0"
          value={kwhConsump}
          onChange={(e) => setKwhConsump(e.target.value)}
          placeholder="0.00"
        />
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
          {saving ? "Guardando..." : isEditing ? "Guardar cambios" : "Añadir enchufe"}
        </button>
      </div>
    </form>
  );
}

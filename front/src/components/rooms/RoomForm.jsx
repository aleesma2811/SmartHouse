import { useState } from "react";
import "../common/forms.css";

export default function RoomForm({ onSubmit, onCancel }) {
  const [name, setName] = useState("");
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
      await onSubmit({ name: name.trim() });
    } catch (err) {
      setError(err.message || "No se pudo crear la habitación");
      setSaving(false);
    }
  }

  return (
    <form onSubmit={handleSubmit}>
      {error && <p className="form-error">{error}</p>}
      <div className="field">
        <label htmlFor="room-name">Nombre de la habitación</label>
        <input
          id="room-name"
          type="text"
          value={name}
          onChange={(e) => setName(e.target.value)}
          placeholder="Ej. Cocina"
          autoFocus
        />
      </div>
      <div className="form-actions">
        <button type="button" className="btn btn-ghost" onClick={onCancel}>
          Cancelar
        </button>
        <button type="submit" className="btn btn-primary" disabled={saving}>
          {saving ? "Creando..." : "Crear habitación"}
        </button>
      </div>
    </form>
  );
}

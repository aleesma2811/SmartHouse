import { useState } from "react";
import "../common/forms.css";

export default function InmuebleForm({ initialInmueble, onSubmit, onCancel }) {
  const isEditing = Boolean(initialInmueble);
  const [nombre, setNombre] = useState(initialInmueble?.Nombre ?? "");
  const [direccion, setDireccion] = useState(initialInmueble?.Direccion ?? "");
  const [ciudad, setCiudad] = useState(initialInmueble?.Ciudad ?? "");
  const [tipo, setTipo] = useState(initialInmueble?.Tipo ?? "casa");
  const [error, setError] = useState("");
  const [saving, setSaving] = useState(false);

  async function handleSubmit(e) {
    e.preventDefault();
    if (!nombre.trim()) {
      setError("El nombre es obligatorio");
      return;
    }
    setSaving(true);
    setError("");
    try {
      await onSubmit({
        nombre: nombre.trim(),
        direccion: direccion.trim(),
        ciudad: ciudad.trim(),
        tipo,
      });
    } catch (err) {
      setError(err.message || "No se pudo guardar el inmueble");
      setSaving(false);
    }
  }

  return (
    <form onSubmit={handleSubmit}>
      {error && <p className="form-error">{error}</p>}

      <div className="field">
        <label htmlFor="inmueble-nombre">Nombre</label>
        <input
          id="inmueble-nombre"
          type="text"
          value={nombre}
          onChange={(e) => setNombre(e.target.value)}
          placeholder="Ej. Casa de la playa"
          autoFocus
        />
      </div>

      <div className="field">
        <label htmlFor="inmueble-tipo">Tipo</label>
        <select id="inmueble-tipo" value={tipo} onChange={(e) => setTipo(e.target.value)}>
          <option value="casa">Casa</option>
          <option value="departamento">Departamento</option>
        </select>
      </div>

      <div className="field">
        <label htmlFor="inmueble-direccion">Dirección</label>
        <input
          id="inmueble-direccion"
          type="text"
          value={direccion}
          onChange={(e) => setDireccion(e.target.value)}
          placeholder="Ej. Av. Siempre Viva 123"
        />
      </div>

      <div className="field">
        <label htmlFor="inmueble-ciudad">Ciudad</label>
        <input
          id="inmueble-ciudad"
          type="text"
          value={ciudad}
          onChange={(e) => setCiudad(e.target.value)}
          placeholder="Ej. Santiago"
        />
      </div>

      <div className="form-actions">
        <button type="button" className="btn btn-ghost" onClick={onCancel}>
          Cancelar
        </button>
        <button type="submit" className="btn btn-primary" disabled={saving}>
          {saving ? "Guardando..." : isEditing ? "Guardar cambios" : "Crear inmueble"}
        </button>
      </div>
    </form>
  );
}

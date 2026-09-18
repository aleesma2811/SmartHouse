import { useEffect, useState } from "react";
import { getInmuebles, createInmueble, updateInmueble, deleteInmueble } from "../api/inmuebles";
import { getRooms } from "../api/rooms";
import InmuebleCard from "../components/inmuebles/InmuebleCard";
import InmuebleForm from "../components/inmuebles/InmuebleForm";
import Modal from "../components/common/Modal";
import { IconPlus } from "../components/common/Icons";
import "./InmueblesPage.css";

export default function InmueblesPage() {
  const [inmuebles, setInmuebles] = useState([]);
  const [rooms, setRooms] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [showForm, setShowForm] = useState(false);
  const [editingInmueble, setEditingInmueble] = useState(null);

  async function loadData() {
    setLoading(true);
    setError("");
    try {
      const [inmueblesData, roomsData] = await Promise.all([getInmuebles(), getRooms()]);
      setInmuebles(inmueblesData || []);
      setRooms(roomsData || []);
    } catch (err) {
      setError(err.message || "No se pudo cargar la información");
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    loadData();
  }, []);

  function openCreateForm() {
    setEditingInmueble(null);
    setShowForm(true);
  }

  function openEditForm(inmueble) {
    setEditingInmueble(inmueble);
    setShowForm(true);
  }

  function closeForm() {
    setShowForm(false);
    setEditingInmueble(null);
  }

  async function handleSubmitInmueble(values) {
    if (editingInmueble) {
      await updateInmueble(editingInmueble.ID, values);
    } else {
      await createInmueble(values);
    }
    closeForm();
    await loadData();
  }

  async function handleDeleteInmueble(inmueble) {
    const roomCount = rooms.filter((r) => r.InmuebleID === inmueble.ID).length;
    const warning =
      roomCount > 0
        ? `"${inmueble.Nombre}" tiene ${roomCount} habitación(es) registradas. ¿Eliminar el inmueble de todas formas?`
        : `¿Eliminar el inmueble "${inmueble.Nombre}"?`;
    if (!window.confirm(warning)) return;

    try {
      await deleteInmueble(inmueble.ID);
      await loadData();
    } catch (err) {
      alert(err.message || "No se pudo eliminar el inmueble");
    }
  }

  return (
    <div className="inmuebles-page">
      <div className="inmuebles-page__header">
        <div>
          <h1>Inmuebles</h1>
          <p className="inmuebles-page__subtitle">
            Casas y departamentos registrados en tu cuenta
          </p>
        </div>
        <button type="button" className="btn btn-primary" onClick={openCreateForm}>
          <IconPlus width={15} height={15} /> Nuevo inmueble
        </button>
      </div>

      {loading && <p className="inmuebles-page__status">Cargando...</p>}
      {error && <p className="inmuebles-page__status inmuebles-page__status--error">{error}</p>}

      {!loading && !error && inmuebles.length === 0 && (
        <p className="inmuebles-page__status">
          Todavía no hay inmuebles. Crea el primero para empezar.
        </p>
      )}

      <div className="inmuebles-page__grid">
        {inmuebles.map((inmueble) => (
          <InmuebleCard
            key={inmueble.ID}
            inmueble={inmueble}
            roomCount={rooms.filter((r) => r.InmuebleID === inmueble.ID).length}
            onEdit={openEditForm}
            onDelete={handleDeleteInmueble}
          />
        ))}
      </div>

      {showForm && (
        <Modal
          title={editingInmueble ? "Editar inmueble" : "Nuevo inmueble"}
          onClose={closeForm}
        >
          <InmuebleForm
            initialInmueble={editingInmueble}
            onSubmit={handleSubmitInmueble}
            onCancel={closeForm}
          />
        </Modal>
      )}
    </div>
  );
}

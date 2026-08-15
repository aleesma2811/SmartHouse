import { useEffect, useState } from "react";
import { Link, useNavigate, useParams } from "react-router-dom";
import { getRoom, deleteRoom } from "../api/rooms";
import { createPlug, updatePlug, deletePlug } from "../api/plugs";
import RoomFloorPlan from "../components/rooms/RoomFloorPlan";
import PlugList from "../components/plugs/PlugList";
import PlugForm from "../components/plugs/PlugForm";
import Modal from "../components/common/Modal";
import "./RoomDetailPage.css";

export default function RoomDetailPage() {
  const { id } = useParams();
  const navigate = useNavigate();

  const [room, setRoom] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [editingPlug, setEditingPlug] = useState(null);
  const [showForm, setShowForm] = useState(false);

  async function loadRoom() {
    setLoading(true);
    setError("");
    try {
      const data = await getRoom(id);
      setRoom(data);
    } catch (err) {
      setError(err.message || "No se pudo cargar la habitación");
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    loadRoom();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [id]);

  function openCreateForm() {
    setEditingPlug(null);
    setShowForm(true);
  }

  function openEditForm(plug) {
    setEditingPlug(plug);
    setShowForm(true);
  }

  async function handleSubmitPlug({ name, kwhConsump, on }) {
    if (editingPlug) {
      await updatePlug(editingPlug.ID, { name, kwhConsump, on });
    } else {
      await createPlug({ name, kwhConsump, on, roomId: id });
    }
    setShowForm(false);
    setEditingPlug(null);
    await loadRoom();
  }

  async function handleToggle(plug) {
    try {
      await updatePlug(plug.ID, {
        name: plug.Name,
        kwhConsump: plug.KwhConsump,
        on: !plug.On,
      });
      await loadRoom();
    } catch (err) {
      alert(err.message || "No se pudo actualizar el enchufe");
    }
  }

  async function handleDeletePlug(plug) {
    if (!window.confirm(`¿Eliminar el enchufe "${plug.Name}"?`)) return;
    try {
      await deletePlug(plug.ID);
      await loadRoom();
    } catch (err) {
      alert(err.message || "No se pudo eliminar el enchufe");
    }
  }

  async function handleDeleteRoom() {
    if (!window.confirm(`¿Eliminar la habitación "${room.Name}" y volver al listado?`)) return;
    try {
      await deleteRoom(room.ID);
      navigate("/");
    } catch (err) {
      alert(err.message || "No se pudo eliminar la habitación");
    }
  }

  if (loading) return <p className="room-detail__status">Cargando...</p>;
  if (error) return <p className="room-detail__status room-detail__status--error">{error}</p>;
  if (!room) return null;

  const plugs = room.Plugs || [];

  return (
    <div className="room-detail">
      <Link to="/" className="room-detail__back">
        &larr; Todas las habitaciones
      </Link>

      <div className="room-detail__header">
        <h1>{room.Name}</h1>
        <button type="button" className="btn btn-danger" onClick={handleDeleteRoom}>
          Eliminar habitación
        </button>
      </div>

      <div className="room-detail__content">
        <div className="room-detail__plan">
          <RoomFloorPlan plugs={plugs} colorSeed={room.ID} label={room.Name} />
        </div>

        <div className="room-detail__plugs">
          <div className="room-detail__plugs-header">
            <h2>Enchufes</h2>
            <button type="button" className="btn btn-primary" onClick={openCreateForm}>
              + Añadir enchufe
            </button>
          </div>

          <PlugList
            plugs={plugs}
            onToggle={handleToggle}
            onEdit={openEditForm}
            onDelete={handleDeletePlug}
          />
        </div>
      </div>

      {showForm && (
        <Modal
          title={editingPlug ? "Editar enchufe" : "Nuevo enchufe"}
          onClose={() => setShowForm(false)}
        >
          <PlugForm
            initialPlug={editingPlug}
            onSubmit={handleSubmitPlug}
            onCancel={() => setShowForm(false)}
          />
        </Modal>
      )}
    </div>
  );
}

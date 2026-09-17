import { useEffect, useState } from "react";
import { Link, useNavigate, useParams } from "react-router-dom";
import { getRoom, deleteRoom } from "../api/rooms";
import { getPlugs, createPlug, updatePlug, deletePlug } from "../api/plugs";
import RoomFloorPlan from "../components/rooms/RoomFloorPlan";
import PlugList from "../components/plugs/PlugList";
import PlugForm from "../components/plugs/PlugForm";
import Modal from "../components/common/Modal";
import "./RoomDetailPage.css";

export default function RoomDetailPage() {
  const { inmuebleId, id } = useParams();
  const navigate = useNavigate();

  const [room, setRoom] = useState(null);
  const [plugs, setPlugs] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [editingPlug, setEditingPlug] = useState(null);
  const [showForm, setShowForm] = useState(false);

  async function loadRoom() {
    setLoading(true);
    setError("");
    try {
      const [roomData, plugsData] = await Promise.all([getRoom(id), getPlugs(id)]);
      setRoom(roomData);
      setPlugs(plugsData || []);
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

  async function handleSubmitPlug({ name, tipo, kwhConsump, on }) {
    if (editingPlug) {
      await updatePlug(editingPlug.ID, { name, tipo, kwhConsump, on });
    } else {
      await createPlug({ name, tipo, kwhConsump, on, roomId: id });
    }
    setShowForm(false);
    setEditingPlug(null);
    await loadRoom();
  }

  async function handleToggle(plug) {
    try {
      await updatePlug(plug.ID, {
        name: plug.Name,
        tipo: plug.Tipo,
        kwhConsump: plug.Consumo,
        on: !plug.On,
      });
      await loadRoom();
    } catch (err) {
      alert(err.message || "No se pudo actualizar el servicio");
    }
  }

  async function handleDeletePlug(plug) {
    if (!window.confirm(`¿Eliminar el servicio "${plug.Name}"?`)) return;
    try {
      await deletePlug(plug.ID);
      await loadRoom();
    } catch (err) {
      alert(err.message || "No se pudo eliminar el servicio");
    }
  }

  async function handleDeleteRoom() {
    if (!window.confirm(`¿Eliminar la habitación "${room.Name}" y volver al listado?`)) return;
    try {
      await deleteRoom(room.ID);
      navigate(`/inmueble/${inmuebleId}`);
    } catch (err) {
      alert(err.message || "No se pudo eliminar la habitación");
    }
  }

  if (loading) return <p className="room-detail__status">Cargando...</p>;
  if (error) return <p className="room-detail__status room-detail__status--error">{error}</p>;
  if (!room) return null;

  return (
    <div className="room-detail">
      <Link to={`/inmueble/${inmuebleId}`} className="room-detail__back">
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
            <h2>Servicios</h2>
            <button type="button" className="btn btn-primary" onClick={openCreateForm}>
              + Añadir servicio
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
          title={editingPlug ? "Editar servicio" : "Nuevo servicio"}
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

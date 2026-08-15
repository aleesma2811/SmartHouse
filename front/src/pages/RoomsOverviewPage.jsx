import { useEffect, useState } from "react";
import { getRooms, createRoom, deleteRoom } from "../api/rooms";
import { getPlugs } from "../api/plugs";
import RoomCard from "../components/rooms/RoomCard";
import RoomForm from "../components/rooms/RoomForm";
import Modal from "../components/common/Modal";
import "./RoomsOverviewPage.css";

export default function RoomsOverviewPage() {
  const [rooms, setRooms] = useState([]);
  const [plugs, setPlugs] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [showForm, setShowForm] = useState(false);

  async function loadData() {
    setLoading(true);
    setError("");
    try {
      const [roomsData, plugsData] = await Promise.all([getRooms(), getPlugs()]);
      setRooms(roomsData || []);
      setPlugs(plugsData || []);
    } catch (err) {
      setError(err.message || "No se pudo cargar la información");
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    loadData();
  }, []);

  async function handleCreateRoom({ name }) {
    await createRoom({ name });
    setShowForm(false);
    await loadData();
  }

  async function handleDeleteRoom(room) {
    const plugCount = plugs.filter((p) => p.RoomID === room.ID).length;
    const warning =
      plugCount > 0
        ? `"${room.Name}" tiene ${plugCount} enchufe(s) registrados. ¿Eliminar la habitación de todas formas?`
        : `¿Eliminar la habitación "${room.Name}"?`;
    if (!window.confirm(warning)) return;

    try {
      await deleteRoom(room.ID);
      await loadData();
    } catch (err) {
      alert(err.message || "No se pudo eliminar la habitación");
    }
  }

  return (
    <div className="rooms-page">
      <div className="rooms-page__header">
        <div>
          <h1>Habitaciones</h1>
          <p className="rooms-page__subtitle">
            Plano general de la casa y sus enchufes inteligentes
          </p>
        </div>
        <button type="button" className="btn btn-primary" onClick={() => setShowForm(true)}>
          + Nueva habitación
        </button>
      </div>

      {loading && <p className="rooms-page__status">Cargando...</p>}
      {error && <p className="rooms-page__status rooms-page__status--error">{error}</p>}

      {!loading && !error && rooms.length === 0 && (
        <p className="rooms-page__status">
          Todavía no hay habitaciones. Crea la primera para empezar.
        </p>
      )}

      <div className="rooms-page__grid">
        {rooms.map((room) => (
          <RoomCard
            key={room.ID}
            room={room}
            plugs={plugs.filter((p) => p.RoomID === room.ID)}
            onDelete={handleDeleteRoom}
          />
        ))}
      </div>

      {showForm && (
        <Modal title="Nueva habitación" onClose={() => setShowForm(false)}>
          <RoomForm onSubmit={handleCreateRoom} onCancel={() => setShowForm(false)} />
        </Modal>
      )}
    </div>
  );
}

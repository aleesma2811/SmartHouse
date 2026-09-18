import { useEffect, useState } from "react";
import { Link, useNavigate, useParams } from "react-router-dom";
import { getInmueble, updateInmueble, deleteInmueble } from "../api/inmuebles";
import { getRooms, createRoom, deleteRoom } from "../api/rooms";
import { getPlugs } from "../api/plugs";
import HouseFloorPlan from "../components/rooms/HouseFloorPlan";
import RoomForm from "../components/rooms/RoomForm";
import InmuebleForm from "../components/inmuebles/InmuebleForm";
import Modal from "../components/common/Modal";
import { IconArrowLeft, IconPencil, IconTrash, IconPlus } from "../components/common/Icons";
import "./InmuebleDetailPage.css";

const TIPO_LABEL = {
  casa: "Casa",
  departamento: "Departamento",
};

export default function InmuebleDetailPage() {
  const { inmuebleId } = useParams();
  const navigate = useNavigate();

  const [inmueble, setInmueble] = useState(null);
  const [rooms, setRooms] = useState([]);
  const [plugs, setPlugs] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [showForm, setShowForm] = useState(false);
  const [showEditInmueble, setShowEditInmueble] = useState(false);

  async function loadData() {
    setLoading(true);
    setError("");
    try {
      const [inmuebleData, roomsData, plugsData] = await Promise.all([
        getInmueble(inmuebleId),
        getRooms(inmuebleId),
        getPlugs(),
      ]);
      setInmueble(inmuebleData);
      setRooms(roomsData || []);
      setPlugs(plugsData || []);
    } catch (err) {
      setError(err.message || "No se pudo cargar el inmueble");
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    loadData();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [inmuebleId]);

  async function handleCreateRoom({ name }) {
    await createRoom({ name, inmuebleId });
    setShowForm(false);
    await loadData();
  }

  async function handleDeleteRoom(room) {
    const plugCount = plugs.filter((p) => p.RoomID === room.ID).length;
    const warning =
      plugCount > 0
        ? `"${room.Name}" tiene ${plugCount} servicio(s) registrados. ¿Eliminar la habitación de todas formas?`
        : `¿Eliminar la habitación "${room.Name}"?`;
    if (!window.confirm(warning)) return;

    try {
      await deleteRoom(room.ID);
      await loadData();
    } catch (err) {
      alert(err.message || "No se pudo eliminar la habitación");
    }
  }

  async function handleUpdateInmueble(values) {
    await updateInmueble(inmueble.ID, values);
    setShowEditInmueble(false);
    await loadData();
  }

  async function handleDeleteInmueble() {
    if (!window.confirm(`¿Eliminar el inmueble "${inmueble.Nombre}" y volver al listado?`)) return;
    try {
      await deleteInmueble(inmueble.ID);
      navigate("/");
    } catch (err) {
      alert(err.message || "No se pudo eliminar el inmueble");
    }
  }

  if (loading) return <p className="inmueble-detail__status">Cargando...</p>;
  if (error) return <p className="inmueble-detail__status inmueble-detail__status--error">{error}</p>;
  if (!inmueble) return null;

  return (
    <div className="inmueble-detail">
      <Link to="/" className="inmueble-detail__back">
        <IconArrowLeft width={15} height={15} /> Todos los inmuebles
      </Link>

      <div className="inmueble-detail__header">
        <div>
          <span className={`inmueble-detail__badge inmueble-detail__badge--${inmueble.Tipo}`}>
            {TIPO_LABEL[inmueble.Tipo] ?? inmueble.Tipo}
          </span>
          <h1>{inmueble.Nombre}</h1>
          <p className="inmueble-detail__subtitle">
            {inmueble.Direccion}
            {inmueble.Ciudad ? `, ${inmueble.Ciudad}` : ""}
          </p>
        </div>
        <div className="inmueble-detail__header-actions">
          <button type="button" className="btn btn-ghost" onClick={() => setShowEditInmueble(true)}>
            <IconPencil width={15} height={15} /> Editar
          </button>
          <button type="button" className="btn btn-danger" onClick={handleDeleteInmueble}>
            <IconTrash width={15} height={15} /> Eliminar
          </button>
        </div>
      </div>

      <div className="inmueble-detail__rooms-header">
        <h2>Planta del inmueble</h2>
        <button type="button" className="btn btn-primary" onClick={() => setShowForm(true)}>
          <IconPlus width={15} height={15} /> Nueva habitación
        </button>
      </div>

      {rooms.length === 0 && (
        <p className="inmueble-detail__status">
          Todavía no hay habitaciones. Crea la primera para empezar.
        </p>
      )}

      {rooms.length > 0 && (
        <HouseFloorPlan
          rooms={rooms}
          plugs={plugs}
          inmuebleId={inmuebleId}
          onDeleteRoom={handleDeleteRoom}
        />
      )}

      {showForm && (
        <Modal title="Nueva habitación" onClose={() => setShowForm(false)}>
          <RoomForm onSubmit={handleCreateRoom} onCancel={() => setShowForm(false)} />
        </Modal>
      )}

      {showEditInmueble && (
        <Modal title="Editar inmueble" onClose={() => setShowEditInmueble(false)}>
          <InmuebleForm
            initialInmueble={inmueble}
            onSubmit={handleUpdateInmueble}
            onCancel={() => setShowEditInmueble(false)}
          />
        </Modal>
      )}
    </div>
  );
}

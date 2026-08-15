import { Link } from "react-router-dom";
import RoomFloorPlan from "./RoomFloorPlan";
import "./RoomCard.css";

export default function RoomCard({ room, plugs, onDelete }) {
  const onCount = plugs.filter((p) => p.On).length;

  return (
    <div className="room-card">
      <Link to={`/room/${room.ID}`} className="room-card__plan">
        <RoomFloorPlan plugs={plugs} colorSeed={room.ID} label={room.Name} />
      </Link>

      <div className="room-card__footer">
        <div>
          <Link to={`/room/${room.ID}`} className="room-card__name">
            {room.Name}
          </Link>
          <p className="room-card__meta">
            {plugs.length} enchufe{plugs.length === 1 ? "" : "s"}
            {onCount > 0 ? ` · ${onCount} encendido${onCount === 1 ? "" : "s"}` : ""}
          </p>
        </div>
        <button
          type="button"
          className="room-card__delete"
          title="Eliminar habitación"
          aria-label="Eliminar habitación"
          onClick={() => onDelete(room)}
        >
          &times;
        </button>
      </div>
    </div>
  );
}

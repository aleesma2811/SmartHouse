import { Link } from "react-router-dom";
import RoomFloorPlan from "./RoomFloorPlan";
import { getRoomTint } from "./roomTheme";
import { IconTrash } from "../common/Icons";
import "./HouseFloorPlan.css";

// Renders every room of an inmueble as one continuous floor plan instead of
// a grid of isolated cards — rooms sit edge-to-edge like real spaces, each
// tinted by its guessed type, sized a little differently to break the
// "identical repeated card" pattern.
export default function HouseFloorPlan({ rooms, plugs, inmuebleId, onDeleteRoom }) {
  return (
    <div className="house-plan">
      <div className="house-plan__grid">
        {rooms.map((room) => {
          const roomPlugs = plugs.filter((p) => p.RoomID === room.ID);
          const onCount = roomPlugs.filter((p) => p.On).length;
          const tint = getRoomTint(room);
          const big = roomPlugs.length >= 3;

          return (
            <div
              key={room.ID}
              className={big ? "house-plan__room house-plan__room--big" : "house-plan__room"}
              style={{ background: `var(--floor-${tint})` }}
            >
              <Link
                to={`/inmueble/${inmuebleId}/room/${room.ID}`}
                className="house-plan__room-link"
              >
                <div className="house-plan__room-plan">
                  <RoomFloorPlan plugs={roomPlugs} room={room} label={room.Name} />
                </div>
                <div className="house-plan__room-footer">
                  <span className="house-plan__room-name">{room.Name}</span>
                  <span className="house-plan__room-meta">
                    {roomPlugs.length} servicio{roomPlugs.length === 1 ? "" : "s"}
                    {onCount > 0 ? ` · ${onCount} encendido${onCount === 1 ? "" : "s"}` : ""}
                  </span>
                </div>
              </Link>

              <button
                type="button"
                className="icon-btn icon-btn--danger house-plan__room-delete"
                title="Eliminar habitación"
                aria-label="Eliminar habitación"
                onClick={() => onDeleteRoom(room)}
              >
                <IconTrash width={15} height={15} />
              </button>
            </div>
          );
        })}
      </div>
    </div>
  );
}

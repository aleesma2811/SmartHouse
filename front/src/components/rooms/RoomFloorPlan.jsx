import { PLAN, getPlugSlots } from "./floorPlanGeometry";
import { getServiceMeta } from "../plugs/serviceTheme";
import { getRoomTint } from "./roomTheme";
import "./RoomFloorPlan.css";

function PlugMarker({ x, y, rotation, on, tipo }) {
  const meta = getServiceMeta(tipo);
  const Icon = meta.Icon;

  return (
    <g transform={`translate(${x} ${y})`}>
      <rect
        x={-7}
        y={-7}
        width={14}
        height={14}
        rx={3.5}
        transform={`rotate(${rotation})`}
        className={on ? "floor-plan__plate floor-plan__plate--on" : "floor-plan__plate"}
        style={on ? { fill: meta.soft, stroke: meta.color } : undefined}
      />
      <g
        transform="translate(-6 -6)"
        style={{ color: on ? meta.color : "var(--color-wall)", opacity: on ? 1 : 0.7 }}
      >
        <Icon width={12} height={12} strokeWidth={2.4} />
      </g>
    </g>
  );
}

export default function RoomFloorPlan({ plugs = [], room, colorSeed = 0, label }) {
  const slots = getPlugSlots(plugs.length);
  const tint = getRoomTint(room ?? { ID: colorSeed });

  return (
    <svg
      className="floor-plan"
      viewBox="0 0 240 160"
      role="img"
      aria-label={label ? `Plano de ${label}` : "Plano de la habitación"}
    >
      <rect
        x={PLAN.wallLeft}
        y={PLAN.wallTop}
        width={PLAN.wallRight - PLAN.wallLeft}
        height={PLAN.wallBottom - PLAN.wallTop}
        className="floor-plan__floor"
        style={{ fill: `var(--floor-${tint})` }}
      />

      {/* Walls, drawn as individual segments so the door gap stays open */}
      <line x1={PLAN.wallLeft} y1={PLAN.wallTop} x2={PLAN.windowStart} y2={PLAN.wallTop} className="floor-plan__wall" />
      <line x1={PLAN.windowEnd} y1={PLAN.wallTop} x2={PLAN.wallRight} y2={PLAN.wallTop} className="floor-plan__wall" />
      <line x1={PLAN.wallRight} y1={PLAN.wallTop} x2={PLAN.wallRight} y2={PLAN.wallBottom} className="floor-plan__wall" />
      <line x1={PLAN.wallRight} y1={PLAN.wallBottom} x2={PLAN.doorEnd} y2={PLAN.wallBottom} className="floor-plan__wall" />
      <line x1={PLAN.doorStart} y1={PLAN.wallBottom} x2={PLAN.wallLeft} y2={PLAN.wallBottom} className="floor-plan__wall" />
      <line x1={PLAN.wallLeft} y1={PLAN.wallBottom} x2={PLAN.wallLeft} y2={PLAN.wallTop} className="floor-plan__wall" />

      {/* Window (double line) */}
      <line x1={PLAN.windowStart} y1={PLAN.wallTop - 2} x2={PLAN.windowEnd} y2={PLAN.wallTop - 2} className="floor-plan__window" />
      <line x1={PLAN.windowStart} y1={PLAN.wallTop + 2} x2={PLAN.windowEnd} y2={PLAN.wallTop + 2} className="floor-plan__window" />

      {/* Door leaf + swing arc */}
      <line x1={PLAN.doorStart} y1={PLAN.wallBottom} x2={PLAN.doorStart} y2={PLAN.wallBottom - 40} className="floor-plan__door" />
      <path
        d={`M ${PLAN.doorEnd} ${PLAN.wallBottom} A 40 40 0 0 0 ${PLAN.doorStart} ${PLAN.wallBottom - 40}`}
        className="floor-plan__door-arc"
      />

      {slots.map((slot, i) => (
        <PlugMarker
          key={i}
          x={slot.x}
          y={slot.y}
          rotation={slot.rotation}
          on={plugs[i]?.On}
          tipo={plugs[i]?.Tipo}
        />
      ))}
    </svg>
  );
}

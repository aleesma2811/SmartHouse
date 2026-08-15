import { PLAN, getPlugSlots } from "./floorPlanGeometry";
import "./RoomFloorPlan.css";

const FLOOR_TINTS = 6;

function PlugMarker({ x, y, rotation, on }) {
  return (
    <g transform={`translate(${x} ${y}) rotate(${rotation})`}>
      <rect
        x={-5}
        y={-4}
        width={10}
        height={8}
        rx={1.5}
        className={on ? "floor-plan__plug floor-plan__plug--on" : "floor-plan__plug"}
      />
      <line x1={-2} y1={-4} x2={-2} y2={-1.5} className="floor-plan__plug-prong" />
      <line x1={2} y1={-4} x2={2} y2={-1.5} className="floor-plan__plug-prong" />
    </g>
  );
}

export default function RoomFloorPlan({ plugs = [], colorSeed = 0, label }) {
  const slots = getPlugSlots(plugs.length);
  const floorClass = `floor-plan__floor floor-plan__floor--${colorSeed % FLOOR_TINTS}`;

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
        className={floorClass}
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
        <PlugMarker key={i} x={slot.x} y={slot.y} rotation={slot.rotation} on={plugs[i]?.On} />
      ))}
    </svg>
  );
}

// Geometry for the simplified architectural floor plan (viewBox 0 0 240 160).
// Walls form a rectangle with a door gap on the bottom wall and a window
// segment on the top wall. Plug markers are distributed evenly along the
// remaining wall length (the door gap is excluded so markers never land on
// the doorway).

export const PLAN = {
  wallLeft: 20,
  wallRight: 220,
  wallTop: 24,
  wallBottom: 136,
  doorStart: 60,
  doorEnd: 100,
  windowStart: 140,
  windowEnd: 180,
};

const { wallLeft, wallRight, wallTop, wallBottom, doorStart, doorEnd } = PLAN;

// Walkable wall segments, in clockwise order, with the door gap removed.
const SEGMENTS = [
  { side: "top", x1: wallLeft, y1: wallTop, x2: wallRight, y2: wallTop },
  { side: "right", x1: wallRight, y1: wallTop, x2: wallRight, y2: wallBottom },
  { side: "bottom", x1: wallRight, y1: wallBottom, x2: doorEnd, y2: wallBottom },
  { side: "bottom", x1: doorStart, y1: wallBottom, x2: wallLeft, y2: wallBottom },
  { side: "left", x1: wallLeft, y1: wallBottom, x2: wallLeft, y2: wallTop },
];

const INSET = 11; // how far inward a plug marker sits from the wall line

function segmentLength(seg) {
  return Math.hypot(seg.x2 - seg.x1, seg.y2 - seg.y1);
}

const TOTAL_LENGTH = SEGMENTS.reduce((sum, seg) => sum + segmentLength(seg), 0);

function inwardOffset(side) {
  switch (side) {
    case "top":
      return { dx: 0, dy: INSET };
    case "bottom":
      return { dx: 0, dy: -INSET };
    case "left":
      return { dx: INSET, dy: 0 };
    case "right":
      return { dx: -INSET, dy: 0 };
    default:
      return { dx: 0, dy: 0 };
  }
}

function pointAtDistance(distance) {
  let remaining = distance;
  for (const seg of SEGMENTS) {
    const len = segmentLength(seg);
    if (remaining <= len || seg === SEGMENTS[SEGMENTS.length - 1]) {
      const t = len === 0 ? 0 : remaining / len;
      const x = seg.x1 + (seg.x2 - seg.x1) * t;
      const y = seg.y1 + (seg.y2 - seg.y1) * t;
      const { dx, dy } = inwardOffset(seg.side);
      const rotation = seg.side === "top" || seg.side === "bottom" ? 0 : 90;
      return { x: x + dx, y: y + dy, side: seg.side, rotation };
    }
    remaining -= len;
  }
  return { x: wallLeft, y: wallTop, side: "top", rotation: 0 };
}

// Returns `count` evenly-spaced slots {x, y, side, rotation} along the walls.
export function getPlugSlots(count) {
  if (count <= 0) return [];
  const slots = [];
  for (let i = 0; i < count; i += 1) {
    const t = ((i + 0.5) / count) * TOTAL_LENGTH;
    slots.push(pointAtDistance(t));
  }
  return slots;
}

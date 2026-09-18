// Hand-rolled line icons, 24x24, stroke-based so color follows currentColor.

function base(props) {
  return {
    width: 18,
    height: 18,
    viewBox: "0 0 24 24",
    fill: "none",
    stroke: "currentColor",
    strokeWidth: 1.8,
    strokeLinecap: "round",
    strokeLinejoin: "round",
    "aria-hidden": true,
    ...props,
  };
}

export function IconPencil(props) {
  return (
    <svg {...base(props)}>
      <path d="M14.5 4.5 19.5 9.5 8 21H3v-5z" />
      <path d="M12.5 6.5 17.5 11.5" />
    </svg>
  );
}

export function IconTrash(props) {
  return (
    <svg {...base(props)}>
      <path d="M4 7h16" />
      <path d="M9 7V4.5A1.5 1.5 0 0 1 10.5 3h3A1.5 1.5 0 0 1 15 4.5V7" />
      <path d="M6 7l1 13.5A1.5 1.5 0 0 0 8.5 22h7a1.5 1.5 0 0 0 1.5-1.5L18 7" />
      <path d="M10 11v6M14 11v6" />
    </svg>
  );
}

export function IconPlus(props) {
  return (
    <svg {...base(props)}>
      <path d="M12 5v14M5 12h14" />
    </svg>
  );
}

export function IconArrowLeft(props) {
  return (
    <svg {...base(props)}>
      <path d="M19 12H5" />
      <path d="M11 6l-6 6 6 6" />
    </svg>
  );
}

export function IconClose(props) {
  return (
    <svg {...base(props)}>
      <path d="M6 6l12 12M18 6 6 18" />
    </svg>
  );
}

export function IconDroplet(props) {
  return (
    <svg {...base(props)}>
      <path d="M12 3c3.5 4.2 6 7.6 6 10.7A6 6 0 1 1 6 13.7C6 10.6 8.5 7.2 12 3Z" />
    </svg>
  );
}

export function IconBolt(props) {
  return (
    <svg {...base(props)}>
      <path d="M13 3 4 14h6l-1 7 9-11h-6z" />
    </svg>
  );
}

export function IconFlame(props) {
  return (
    <svg {...base(props)}>
      <path d="M12 2c1 3-2.5 4.2-2.5 7.5a2.5 2.5 0 0 0 5 0c0-1-.5-1.8-.5-1.8 2 1 3 3.3 3 5.3a5 5 0 0 1-10 0C6 8.8 10 7 12 2Z" />
    </svg>
  );
}

export function IconHome(props) {
  return (
    <svg {...base(props)}>
      <path d="M4 11 12 4l8 7" />
      <path d="M6 10v10h5v-6h2v6h5V10" />
    </svg>
  );
}

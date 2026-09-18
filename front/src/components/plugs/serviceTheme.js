import { IconDroplet, IconBolt, IconFlame } from "../common/Icons";

// Single source of truth for how each service type looks and is measured.
// Used by the plug form/list and the floor-plan illustrations so a service
// type is always the same color + icon + unit everywhere in the app.
export const SERVICE_META = {
  agua: {
    label: "Agua",
    unit: "Litros",
    color: "var(--color-agua)",
    colorDark: "var(--color-agua-dark)",
    soft: "var(--color-agua-soft)",
    Icon: IconDroplet,
  },
  luz: {
    label: "Luz",
    unit: "kWh",
    color: "var(--color-luz)",
    colorDark: "var(--color-luz-dark)",
    soft: "var(--color-luz-soft)",
    Icon: IconBolt,
  },
  gas: {
    label: "Gas",
    unit: "M3",
    color: "var(--color-gas)",
    colorDark: "var(--color-gas-dark)",
    soft: "var(--color-gas-soft)",
    Icon: IconFlame,
  },
};

export function getServiceMeta(tipo) {
  return SERVICE_META[tipo] ?? SERVICE_META.luz;
}

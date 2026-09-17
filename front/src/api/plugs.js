import { request } from "./http";

export function getPlugs(roomId) {
  return request(roomId ? `/servicios?roomId=${roomId}` : "/servicios");
}

export function getPlug(id) {
  return request(`/servicios/${id}`);
}

export function createPlug({ name, tipo, kwhConsump, roomId, on }) {
  return request("/servicios", {
    method: "POST",
    body: JSON.stringify({
      Name: name,
      Tipo: tipo,
      Consumo: Number(kwhConsump) || 0,
      RoomID: Number(roomId),
      On: Boolean(on),
    }),
  });
}

export function updatePlug(id, { name, tipo, kwhConsump, on }) {
  return request(`/servicios/${id}`, {
    method: "PUT",
    body: JSON.stringify({
      Name: name,
      Tipo: tipo,
      Consumo: Number(kwhConsump) || 0,
      On: Boolean(on),
    }),
  });
}

export function deletePlug(id) {
  return request(`/servicios/${id}`, { method: "DELETE" });
}

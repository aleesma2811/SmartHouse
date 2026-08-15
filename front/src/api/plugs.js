import { request } from "./http";

export function getPlugs() {
  return request("/plugs");
}

export function getPlug(id) {
  return request(`/plugs/${id}`);
}

export function createPlug({ name, kwhConsump, roomId, on }) {
  return request("/plugs", {
    method: "POST",
    body: JSON.stringify({
      Name: name,
      KwhConsump: Number(kwhConsump) || 0,
      RoomID: Number(roomId),
      On: Boolean(on),
    }),
  });
}

export function updatePlug(id, { name, kwhConsump, on }) {
  return request(`/plugs/${id}`, {
    method: "PUT",
    body: JSON.stringify({
      Name: name,
      KwhConsump: Number(kwhConsump) || 0,
      On: Boolean(on),
    }),
  });
}

export function deletePlug(id) {
  return request(`/plugs/${id}`, { method: "DELETE" });
}

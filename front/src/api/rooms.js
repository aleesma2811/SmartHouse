import { request } from "./http";

export function getRooms(inmuebleId) {
  return request(inmuebleId ? `/rooms?inmuebleId=${inmuebleId}` : "/rooms");
}

export function getRoom(id) {
  return request(`/rooms/${id}`);
}

export function createRoom({ name, inmuebleId }) {
  return request("/rooms", {
    method: "POST",
    body: JSON.stringify({ Name: name, InmuebleID: Number(inmuebleId) }),
  });
}

export function deleteRoom(id) {
  return request(`/rooms/${id}`, { method: "DELETE" });
}

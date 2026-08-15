import { request } from "./http";

export function getRooms() {
  return request("/rooms");
}

export function getRoom(id) {
  return request(`/rooms/${id}`);
}

export function createRoom({ name }) {
  return request("/rooms", {
    method: "POST",
    body: JSON.stringify({ Name: name }),
  });
}

export function deleteRoom(id) {
  return request(`/rooms/${id}`, { method: "DELETE" });
}

import { request } from "./http";

export function getInmuebles() {
  return request("/inmuebles");
}

export function getInmueble(id) {
  return request(`/inmuebles/${id}`);
}

export function createInmueble({ nombre, direccion, ciudad, tipo }) {
  return request("/inmuebles", {
    method: "POST",
    body: JSON.stringify({
      Nombre: nombre,
      Direccion: direccion,
      Ciudad: ciudad,
      Tipo: tipo,
    }),
  });
}

export function updateInmueble(id, { nombre, direccion, ciudad, tipo }) {
  return request(`/inmuebles/${id}`, {
    method: "PUT",
    body: JSON.stringify({
      Nombre: nombre,
      Direccion: direccion,
      Ciudad: ciudad,
      Tipo: tipo,
    }),
  });
}

export function deleteInmueble(id) {
  return request(`/inmuebles/${id}`, { method: "DELETE" });
}

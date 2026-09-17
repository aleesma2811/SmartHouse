import { BrowserRouter, Routes, Route } from "react-router-dom";
import Header from "./components/layout/Header";
import InmueblesPage from "./pages/InmueblesPage";
import InmuebleDetailPage from "./pages/InmuebleDetailPage";
import RoomDetailPage from "./pages/RoomDetailPage";

export default function App() {
  return (
    <BrowserRouter>
      <Header />
      <main>
        <Routes>
          <Route path="/" element={<InmueblesPage />} />
          <Route path="/inmueble/:inmuebleId" element={<InmuebleDetailPage />} />
          <Route path="/inmueble/:inmuebleId/room/:id" element={<RoomDetailPage />} />
        </Routes>
      </main>
    </BrowserRouter>
  );
}

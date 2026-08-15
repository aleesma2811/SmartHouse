import { BrowserRouter, Routes, Route } from "react-router-dom";
import Header from "./components/layout/Header";
import RoomsOverviewPage from "./pages/RoomsOverviewPage";
import RoomDetailPage from "./pages/RoomDetailPage";

export default function App() {
  return (
    <BrowserRouter>
      <Header />
      <main>
        <Routes>
          <Route path="/" element={<RoomsOverviewPage />} />
          <Route path="/room/:id" element={<RoomDetailPage />} />
        </Routes>
      </main>
    </BrowserRouter>
  );
}

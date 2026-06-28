import { BrowserRouter, Routes, Route } from "react-router-dom";
import Header from "./components/Header";
import Footer from "./components/Footer";
import HomePage from "./pages/HomePage";
import MenuPage from "./pages/MenuPage";

export default function App() {
  return (
    <BrowserRouter>
      <Header />
      <div className="min-h-screen pt-24 bg-background">
        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/menu" element={<MenuPage />} />
        </Routes>
      </div>
      <Footer />
    </BrowserRouter>
  );
}

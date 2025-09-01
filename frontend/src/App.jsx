import React from "react";
import { BrowserRouter, Routes, Route, Link, Navigate } from "react-router-dom";
import TypePage from "./pages/TypePage";

export default function App() {
  return (
    <BrowserRouter>
      <div className="min-h-screen bg-slate-50">
        <div className="max-w-3xl mx-auto px-4 pt-6">
          <header className="mb-4">
            <h1 className="text-3xl font-bold">Hacker News — Clone</h1>

            {/* Nav bar */}
            <nav className="flex gap-4 px-3 py-2" style={{backgroundColor:"rgba(134, 128, 128, 0.87)", borderRadius:"8px", padding:"6px"}}>
              <Link to="/" className="text-black font-medium text-sm" style={{margin:"5px"}}>| all</Link>
              <Link to="/type/story" className="text-black font-medium text-sm"style={{margin:"5px"}}>| new / story</Link>
              <Link to="/type/comment" className="text-black font-medium text-sm"style={{margin:"5px"}}>| comments</Link>
              <Link to="/type/job" className="text-black font-medium text-sm"style={{margin:"5px"}}>| jobs</Link>
              <Link to="/type/poll" className="text-black font-medium text-sm"style={{margin:"5px"}}>| polls</Link>
              <Link to="/type/pollopt" className="text-black font-medium text-sm"style={{margin:"5px"}}>| pollopts</Link>
            </nav>
          </header>

          <main>
            <Routes>
              {/* default route -> all types */}
              <Route path="/" element={<TypePage typeParam={null} />} />
              {/* route per type */}
              <Route path="/type/:type" element={<TypePage />} />
              {/* fallback */}
              <Route path="*" element={<Navigate to="/" replace />} />
            </Routes>
          </main>
        </div>
      </div>
    </BrowserRouter>
  );
}

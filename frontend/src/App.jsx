import React, { useEffect, useState } from "react";

const API_BASE = import.meta.env.VITE_API_BASE_URL || "http://localhost:8080";

export default function App() {
  const [page, setPage] = useState(1);
  const [limit] = useState(20);
  const [items, setItems] = useState([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  useEffect(() => {
    fetchPage(page);
  }, [page]);

  async function fetchPage(p) {
    setLoading(true);
    setError(null);
    try {
      const res = await fetch(`${API_BASE}/api/topstories?page=${p}&limit=${limit}`);
      if (!res.ok) throw new Error(`Server error: ${res.status}`);
      const body = await res.json();
      setItems(body.items || []);
      setTotal(body.total || 0);
    } catch (err) {
      setError(err.message || "Unknown error");
    } finally {
      setLoading(false);
    }
  }

  function formatTime(epochSec) {
    if (!epochSec) return "-";
    const d = new Date(epochSec * 1000);
    return d.toLocaleString();
  }

  const totalPages = Math.max(1, Math.ceil(total / limit));

  return (
    <div className="min-h-screen px-10 sm:px-12 lg:px-14 p-10 bg-blue-500">
      <div className="app-container max-w-3xl mx-auto pt-6">
        <header className="mb-6">
          <h1 className="text-3xl font-bold ">Hacker News — Clone</h1>
        </header>

        <main>
          <div className="mb-4 flex items-center justify-between">
            <div className="flex items-center gap-2">
              <button
                className="px-3 py-1 rounded bg-slate-200 hover:bg-slate-300"
                onClick={() => setPage((p) => Math.max(1, p - 1))}
                disabled={page === 1 || loading}
              >
                Prev
              </button>
              <button
                className="px-3 py-1 rounded bg-slate-200 hover:bg-slate-300 p-50"
                onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
                disabled={page >= totalPages || loading}
              >
                Next
              </button>
              <span className="text-sm text-slate-600">
                Page {page} / {totalPages}
              </span>
            </div>

            <div>
              <button
                className="px-3 py-1 rounded bg-indigo-600 text-white hover:bg-indigo-700"
                onClick={() => fetchPage(1)}
                disabled={loading}
              >
                Refresh
              </button>
            </div>
          </div>

          {loading ? (
            <div className="p-6 bg-white rounded shadow">Loading...</div>
          ) : error ? (
            <div className="p-4 bg-red-50 text-red-700 rounded">Error: {error}</div>
          ) : (
            <ul className="space-y-3">
              {items.map((it) => (
                <li
                  key={it.id}
                  className="p-4 bg-white rounded shadow-sm max-w-2xl mx-auto"
                >
                  <a
                    className="text-lg font-medium text-sky-700 hover:underline"
                    href={it.url || "#"}
                    target="_blank"
                    rel="noreferrer"
                  >
                    {it.title || `(${it.type || "item"})`}
                  </a>
                  <div className="text-sm text-slate-600">
                    by {it.by || "-"} • {formatTime(it.time)} • {it.score ?? 0} pts
                  </div>
                </li>
              ))}
            </ul>
          )}

          <footer className="mt-6 text-sm text-slate-500">
            Total stories: {total}
          </footer>
        </main>
      </div>
    </div>
  );
}

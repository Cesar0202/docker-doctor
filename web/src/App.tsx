import { useEffect, useState } from 'react';
import { Activity, Box, HardDrive, Network, Layers, ShieldCheck, ShieldAlert } from 'lucide-react';
import { PieChart, Pie, Cell, ResponsiveContainer, Tooltip, Legend } from 'recharts';

interface ReportData {
  System: { IsReachable: boolean };
  Containers: { Total: number; Stopped: number };
  Images: { Total: number; Dangling: number };
  Volumes: { Total: number; Orphaned: number };
  Networks: { Total: number; Unused: number };
  Ports: { TotalExposed: number; InUse: number[] };
}

const COLORS = ['#3b82f6', '#10b981', '#f59e0b', '#ef4444', '#8b5cf6'];

function App() {
  const [data, setData] = useState<ReportData | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // Polling every 5 seconds
    const fetchData = async () => {
      try {
        const res = await fetch('http://localhost:8080/api/report');
        if (res.ok) {
          const json = await res.json();
          setData(json);
        }
      } catch (e) {
        console.error("Error fetching report", e);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
    const interval = setInterval(fetchData, 5000);
    return () => clearInterval(interval);
  }, []);

  if (loading && !data) {
    return <div className="loading">Analizando entorno Docker...</div>;
  }

  // Fallback data for preview if backend is not running
  const d = data || {
    System: { IsReachable: true },
    Containers: { Total: 12, Stopped: 4 },
    Images: { Total: 25, Dangling: 8 },
    Volumes: { Total: 5, Orphaned: 2 },
    Networks: { Total: 4, Unused: 1 },
    Ports: { TotalExposed: 8, InUse: [] }
  };

  const chartData = [
    { name: 'Contenedores Activos', value: d.Containers.Total - d.Containers.Stopped },
    { name: 'Detenidos', value: d.Containers.Stopped },
  ];

  return (
    <div className="dashboard">
      <header className="header">
        <h1><Activity size={36} color="#60a5fa" /> Docker Doctor</h1>
        <div className={`header-status ${d.System.IsReachable ? 'online' : 'offline'}`}>
          {d.System.IsReachable ? <ShieldCheck size={20} /> : <ShieldAlert size={20} />}
          {d.System.IsReachable ? 'Sistema Operativo' : 'Daemon Inaccesible'}
        </div>
      </header>

      <div className="grid">
        <div className="card">
          <div className="card-header"><Box size={18} /> Contenedores</div>
          <div className="card-value">{d.Containers.Total}</div>
          <div className={`card-subtext ${d.Containers.Stopped > 0 ? 'warning' : ''}`}>
            <span>Detenidos:</span> <span>{d.Containers.Stopped}</span>
          </div>
        </div>

        <div className="card">
          <div className="card-header"><Layers size={18} /> Imágenes</div>
          <div className="card-value">{d.Images.Total}</div>
          <div className={`card-subtext ${d.Images.Dangling > 0 ? 'warning' : ''}`}>
            <span>Dangling (Basura):</span> <span>{d.Images.Dangling}</span>
          </div>
        </div>

        <div className="card">
          <div className="card-header"><HardDrive size={18} /> Volúmenes</div>
          <div className="card-value">{d.Volumes.Total}</div>
          <div className={`card-subtext ${d.Volumes.Orphaned > 0 ? 'danger' : ''}`}>
            <span>Huérfanos:</span> <span>{d.Volumes.Orphaned}</span>
          </div>
        </div>

        <div className="card">
          <div className="card-header"><Network size={18} /> Redes</div>
          <div className="card-value">{d.Networks.Total}</div>
          <div className="card-subtext">
            <span>Sin uso:</span> <span>{d.Networks.Unused}</span>
          </div>
        </div>
      </div>

      <div className="charts-section">
        <div className="chart-container">
          <h2 className="chart-title">Distribución de Contenedores</h2>
          <div style={{ height: 300 }}>
            <ResponsiveContainer width="100%" height="100%">
              <PieChart>
                <Pie
                  data={chartData}
                  cx="50%"
                  cy="50%"
                  innerRadius={60}
                  outerRadius={100}
                  paddingAngle={5}
                  dataKey="value"
                  stroke="none"
                >
                  {chartData.map((entry, index) => (
                    <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                  ))}
                </Pie>
                <Tooltip 
                  contentStyle={{ backgroundColor: 'rgba(30, 41, 59, 0.9)', border: 'none', borderRadius: '8px', color: '#fff' }}
                  itemStyle={{ color: '#fff' }}
                />
                <Legend />
              </PieChart>
            </ResponsiveContainer>
          </div>
        </div>
      </div>
    </div>
  );
}

export default App;

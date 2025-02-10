import React, { useEffect, useState } from "react";
import './App.css'

interface ContainerData {
  containerID: string;
  ip: string;
  status: string;
  timestamp: string;
  datestamp: string;
}

const API_URL = "http://localhost:8080/ContainerList"; 

const App: React.FC = () => {
  const [data, setData] = useState<ContainerData[]>([]);
  const [loading, setLoading] = useState<boolean>(true);

  const fetchData = async () => {
    try {
      const response = await fetch(API_URL);
      const result = await response.json();
      console.log(JSON.stringify(result))
      setData(result);
    } catch (error) {
      console.error("Ошибка загрузки данных:", error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
    const interval = setInterval(fetchData, 5000);
    return () => clearInterval(interval);
  }, []);

  return (
    <div className="p-4">
      <h1 className="text-xl font-bold mb-4">Мониторинг контейнеров</h1>
      {loading ? (
        <p>Загрузка...</p>
      ) : (
        <table className="table-auto w-full border-collapse border border-gray-300">
          <thead>
            <tr className="bg-gray-200">
              <th className="border p-2">IP</th>
              <th className="border p-2">Status</th>
              <th className="border p-2">Timestamp</th>
              <th className="border p-2">Datestamp</th>
            </tr>
          </thead>
          <tbody>
            {data.map((item, index) => (
              <tr key={index} className="border">
                <td className="border p-2">{item.ip}</td>
                <td className={`border p-2 ${item.status === "down" ? "text-red-500" : "text-green-500"}`}>{item.status}</td>
                <td className="border p-2">{new Date(item.timestamp).toLocaleString()}</td>
                <td className="border p-2">{new Date(item.datestamp).toLocaleString()}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
}
  export default App
  
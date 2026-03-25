import React, { useState, useEffect } from 'react';

const App = () => {
  const [tasks, setTasks] = useState([]);
  const [newTitle, setNewTitle] = useState('');
  const API_URL =`${import.meta.env.VITE_Backend_API_URL}/tasks`;

  // --- READ ---
  const fetchTasks = async () => {
    try {
      const res = await fetch(API_URL);
      const data = await res.json();
      setTasks(data || []);
    } catch (err) {
      console.error("Failed to fetch:", err);
    }
  };

  useEffect(() => { fetchTasks(); }, []);

  // --- CREATE ---
  const addTask = async (e) => {
    e.preventDefault();
    await fetch(API_URL, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ title: newTitle, status: 'Pending' }),
    });
    setNewTitle('');
    fetchTasks();
  };

  // --- DELETE ---
  const deleteTask = async (title) => {
    await fetch(`${API_URL}/${title}`, { method: 'DELETE' });
    fetchTasks();
  };

  return (
    <div style={{ padding: '20px', fontFamily: 'sans-serif' }}>
      <h1>Task Tracker (Go + React)</h1>
      
      <form onSubmit={addTask} style={{ marginBottom: '20px' }}>
        <input 
          value={newTitle} 
          onChange={(e) => setNewTitle(e.target.value)}
          placeholder="Enter new task..." 
          style={{ padding: '8px', marginRight: '10px' }}
        />
        <button type="submit" style={{ padding: '8px 16px', cursor: 'pointer' }}>Add Task</button>
      </form>

      <table border="1" cellPadding="10" style={{ width: '100%', borderCollapse: 'collapse' }}>
        <thead>
          <tr style={{ backgroundColor: '#f4f4f4' }}>
            <th>Title</th>
            <th>Status</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {tasks.map((task, index) => (
            <tr key={index}>
              <td>{task.title}</td>
              <td>{task.status}</td>
              <td>
                <button 
                  onClick={() => deleteTask(task.title)}
                  style={{ color: 'red', border: '1px solid red', background: 'none', cursor: 'pointer' }}
                >
                  Delete
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default App;
import React, { useEffect, useState } from "react";
import axios from "axios";
import Swal from 'sweetalert2';
export default function UsersPage() {
  const [users, setUsers] = useState([]);

  useEffect(() => {
    const fetchUsers = async () => {
      try {
        const res = await axios.get("http://localhost:8081/users", {
          headers: {
            Authorization: `Bearer ${localStorage.getItem("token")}`,
          },
        });
        console.log(res.data.users); // Make sure the data structure matches
        setUsers(res.data.users);
      } catch (error) {
        console.error("Failed to fetch users:", error);
        handleLogout();
      }
    };

    fetchUsers();
  }, []);

  const handleLogout = async () => {
    try {
      const token = localStorage.getItem('token');
      if (token) {
        await axios.post('http://localhost:8081/auth/logout', {}, {
          headers: {
            'Authorization': `Bearer ${token}`,
          }
        });
      }
      
      localStorage.removeItem('token');
       // Clear the token from localStorage or sessionStorage
    localStorage.removeItem('token');  // Assuming you stored it in localStorage
  
    // Optionally, you can clear sessionStorage too
    sessionStorage.removeItem('token');
  
    // Redirect to login or homepage
    window.location.href = '/'; // Redirect to the login page or home
      setUsers(null);
  
    } catch (err) {
      console.error('Logout error:', err);
      Swal.fire({
        title: 'Logout Error',
        text: err,
        icon: 'error'
      });
    }
  };
  return (
    <div className="p-6">
      <button
            className="text-blue-600 ml-2 underline"
            onClick={handleLogout}
          >
            Logout
          </button>
      <h1 className="text-2xl font-bold mb-4">All Users</h1>
      {users.length === 0 ? (
        <p>No users found.</p>
      ) : (
        <ul className="space-y-2">
          {users.map((user, index) => (
            <li key={index} className="p-2 bg-gray-100 rounded shadow">
              <p><strong>Username:</strong> {user.username}</p>
              <p><strong>Email:</strong> {user.email}</p>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}

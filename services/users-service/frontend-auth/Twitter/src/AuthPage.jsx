import React, { useState } from "react";
import axios from "axios";
import Swal from 'sweetalert2';

export default function AuthPage() {
  const [isLogin, setIsLogin] = useState(true);
  const [formData, setFormData] = useState({
    username: "",
    password: "",
    email: "",
  });

  const handleChange = (e) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };
  const handleSubmit = async (e) => {
    e.preventDefault();
    const url = isLogin
      ? "http://localhost:8081/login"
      : "http://localhost:8081/register";
  
    try {
      const payload = isLogin
        ? {
            username: formData.username,
            password: formData.password,
          }
        : {
          username: formData.username,
          firstName: formData.firstName,
          lastName: formData.lastName,
          email: formData.email,
          password: formData.password
      
      }
      
  
      const res = await axios.post(url, payload, {
        headers: {
          "Content-Type": "application/json",
        },
      });
  
      if (isLogin) {
        Swal.fire({
          title: 'Success!',
          text: 'You have successfully logged in!'+res.data.token,
          icon: 'success',
          confirmButtonText: 'OK',
          confirmButtonColor: '#3b82f6'
        });
      } else {
        Swal.fire({
          title: 'Registered!',
          text: 'Your account has been created',
          icon: 'success',
          confirmButtonText: 'OK',
          confirmButtonColor: '#3b82f6'
        });
        setIsLogin(true);
        setFormData({
          username: "",
          password: "",
          email: "",
          firstName: "",
          lastName: ""
        });
      }
    } catch (err) {
      console.error("Auth error:", err.response?.data || err.message);
      Swal.fire({
        title: 'Error!',
        text: err.response?.data?.error || 
             err.response?.data?.error_description || 
             err.message,
        icon: 'error',
        confirmButtonText: 'OK',
        confirmButtonColor: '#3b82f6'
      });
    }
  };
  

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-100">
      <div className="bg-white shadow-xl rounded-2xl p-8 w-96">
      <div className="flex justify-center mb-6">
          <img 
            src="/free-twitter.png" 
            alt="Twitter Logo" 
            className="h-12 w-25"
            width="100px"
            height="100px"
          />
        </div>
        <h2 className="text-xl font-bold mb-4">
          {isLogin ? "Login" : "Register"}
        </h2>
        <form onSubmit={handleSubmit} className="space-y-4">
          <input
            name="username"
            type="text"
            placeholder="Username"
            value={formData.username}
            onChange={handleChange}
            className="w-full p-2 border rounded"
            required
          /><br/>
         {!isLogin && (
  <>
    <input
      name="email"
      type="email"
      placeholder="Email"
      value={formData.email}
      onChange={handleChange}
      className="w-full p-2 border rounded"
      required
    /><br/>
    <input
      name="firstName"
      type="text"
      placeholder="First Name"
      value={formData.firstName || ""}
      onChange={handleChange}
      className="w-full p-2 border rounded"
    /><br/>
    <input
      name="lastName"
      type="text"
      placeholder="Last Name"
      value={formData.lastName || ""}
      onChange={handleChange}
      className="w-full p-2 border rounded"
    /><br/>
  </>
)}

          <input
            name="password"
            type="password"
            placeholder="Password"
            value={formData.password}
            onChange={handleChange}
            className="w-full p-2 border rounded"
            required
          /><br/>
          <button
            type="submit"
            className="w-full bg-blue-600 text-white py-2 rounded hover:bg-blue-700"
          >
            {isLogin ? "Login" : "Register"}
          </button>
        </form>
        <div className="text-sm text-center mt-4">
          {isLogin ? "Don't have an account?" : "Already have an account?"}
          <button
            className="text-blue-600 ml-2 underline"
            onClick={() => setIsLogin(!isLogin)}
          >
            {isLogin ? "Register" : "Login"}
          </button>
        </div>
      </div>
    </div>
  );
}

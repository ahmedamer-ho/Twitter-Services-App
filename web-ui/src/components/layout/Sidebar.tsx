import React from 'react';
import { NavLink } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import { FaHome, FaUser, FaHashtag, FaBell, FaEnvelope } from 'react-icons/fa';
import { FiLogOut } from 'react-icons/fi';
import '../../styles/layout.css';

export const Sidebar: React.FC = () => {
    const { user, logout } = useAuth();

    if (!user) return null;

    const navItems = [
        { icon: FaHome, label: 'Home', path: '/' },
        { icon: FaHashtag, label: 'Explore', path: '/explore' },
        { icon: FaBell, label: 'Notifications', path: '/notifications' },
        { icon: FaEnvelope, label: 'Messages', path: '/messages' },
        { icon: FaUser, label: 'Profile', path: `/profile/${user.id}` },
    ];

    return (
        <nav className="layout-sidebar">
            <div className="sidebar-top">
                <div className="logo-container">
                    <h1 className="logo-text">𝕏</h1>
                </div>

                {navItems.map((item) => (
                    <NavLink
                        key={item.path}
                        to={item.path}
                        className={({ isActive }) => `nav-link ${isActive ? 'active' : ''}`}
                    >
                        <item.icon className="nav-icon" />
                        <span className="nav-label">{item.label}</span>
                    </NavLink>
                ))}

                <button className="tweet-btn">Tweet</button>
            </div>

            <div className="user-profile-pill">
                <div className="user-info">
                    <div className="avatar">
                        {user.username[0]?.toUpperCase()}
                    </div>
                    <div className="user-details">
                        <div className="font-bold">{user.username}</div>
                        <div className="text-gray-500 text-sm">@{user.username}</div>
                    </div>
                </div>
                <button onClick={logout} className="logout-btn" title="Logout">
                    <FiLogOut size={20} />
                </button>
            </div>
        </nav>
    );
};

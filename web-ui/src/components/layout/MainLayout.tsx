import React from 'react';
import { Sidebar } from './Sidebar';
import { RightBar } from './RightBar';
import { Outlet } from 'react-router-dom';
import '../../styles/layout.css';

export const MainLayout: React.FC = () => {
    return (
        <div className="app-layout">
            <header>
                <Sidebar />
            </header>

            <main className="layout-main">
                <Outlet />
            </main>

            <div>
                <RightBar />
            </div>
        </div>
    );
};

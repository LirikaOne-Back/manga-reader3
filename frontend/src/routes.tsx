import React from 'react';
import { createBrowserRouter, Navigate } from 'react-router-dom';
import Layout from './components/layout/Layout';
import Home from './pages/Home';
import Catalog from './pages/Catalog';
import MangaDetails from './pages/MangaDetails';
import ChapterReader from './pages/ChapterReader';
import Login from './pages/Login';
import Register from './pages/Register';
import Profile from './pages/Profile';
import NotFound from './pages/NotFound';

// Простой Guard для защищенных маршрутов
const ProtectedRoute = ({ children }: { children: React.ReactNode }) => {
    // В реальном приложении проверка авторизации через состояние или хук
    const isAuthenticated = localStorage.getItem('accessToken') !== null;

    if (!isAuthenticated) {
        // Перенаправляем на страницу входа, если пользователь не авторизован
        return <Navigate to="/login" replace />;
    }

    return <>{children}</>;
};

const router = createBrowserRouter([
    {
        path: '/',
        element: <Layout />,
        children: [
            {
                index: true,
                element: <Home />,
            },
            {
                path: 'catalog',
                element: <Catalog />,
            },
            {
                path: 'manga/:id',
                element: <MangaDetails />,
            },
            {
                path: 'manga/:mangaId/chapter/:chapterId',
                element: <ChapterReader />,
            },
            {
                path: 'login',
                element: <Login />,
            },
            {
                path: 'register',
                element: <Register />,
            },
            {
                path: 'profile',
                element: (
                    <ProtectedRoute>
                        <Profile />
                    </ProtectedRoute>
                ),
            },
            {
                path: 'bookmarks',
                element: (
                    <ProtectedRoute>
                        <Profile activeTab="bookmarks" />
                    </ProtectedRoute>
                ),
            },
            {
                path: '*',
                element: <NotFound />,
            },
        ],
    },
]);

export default router;
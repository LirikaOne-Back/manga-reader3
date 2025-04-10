import React from 'react';
import { Link } from 'react-router-dom';

const NotFound: React.FC = () => {
    return (
        <div className="min-h-screen bg-gray-100 dark:bg-dark-900 flex items-center justify-center px-4">
            <div className="max-w-lg w-full bg-white dark:bg-dark-800 rounded-lg shadow-md p-8 text-center">
                <h1 className="text-6xl font-bold text-primary-600 dark:text-primary-400 mb-4">404</h1>
                <h2 className="text-2xl font-semibold mb-4">Страница не найдена</h2>
                <p className="text-gray-600 dark:text-gray-400 mb-8">
                    Извините, страница, которую вы ищете, не существует или была перемещена.
                </p>
                <div className="space-y-4">
                    <Link
                        to="/"
                        className="block w-full py-3 px-4 bg-primary-600 hover:bg-primary-700 text-white font-medium rounded-md transition-colors"
                    >
                        Вернуться на главную
                    </Link>
                    <Link
                        to="/catalog"
                        className="block w-full py-3 px-4 bg-gray-200 hover:bg-gray-300 dark:bg-dark-700 dark:hover:bg-dark-600 text-gray-800 dark:text-gray-200 font-medium rounded-md transition-colors"
                    >
                        Перейти в каталог
                    </Link>
                </div>
            </div>
        </div>
    );
};

export default NotFound;
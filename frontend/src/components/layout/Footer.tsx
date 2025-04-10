import React from 'react';
import { Link } from 'react-router-dom';
import { BookOpenIcon } from '@heroicons/react/24/outline';

const Footer: React.FC = () => {
    const currentYear = new Date().getFullYear();

    return (
        <footer className="bg-white dark:bg-dark-800 py-8 border-t border-gray-200 dark:border-dark-700">
            <div className="container mx-auto px-4">
                <div className="flex flex-col md:flex-row justify-between">
                    {/* Логотип и описание */}
                    <div className="mb-6 md:mb-0 md:w-1/3">
                        <Link to="/" className="flex items-center mb-4">
                            <BookOpenIcon className="h-8 w-8 text-primary-600" />
                            <span className="ml-2 text-xl font-bold text-dark-900 dark:text-white">MangaReader</span>
                        </Link>
                        <p className="text-gray-600 dark:text-gray-400 text-sm">
                            Читай любимую мангу онлайн в любое время и в любом месте. Тысячи манги и манхвы с регулярными обновлениями.
                        </p>
                    </div>

                    {/* Навигация */}
                    <div className="grid grid-cols-2 sm:grid-cols-3 gap-8">
                        <div>
                            <h3 className="text-base font-semibold text-gray-800 dark:text-gray-200 mb-3">Навигация</h3>
                            <ul className="space-y-2">
                                <li>
                                    <Link to="/" className="text-gray-600 dark:text-gray-400 hover:text-primary-600 dark:hover:text-primary-400 text-sm">
                                        Главная
                                    </Link>
                                </li>
                                <li>
                                    <Link to="/catalog" className="text-gray-600 dark:text-gray-400 hover:text-primary-600 dark:hover:text-primary-400 text-sm">
                                        Каталог
                                    </Link>
                                </li>
                                <li>
                                    <Link to="/catalog?sortBy=rating&sortDesc=true" className="text-gray-600 dark:text-gray-400 hover:text-primary-600 dark:hover:text-primary-400 text-sm">
                                        Популярное
                                    </Link>
                                </li>
                                <li>
                                    <Link to="/catalog?status=ongoing" className="text-gray-600 dark:text-gray-400 hover:text-primary-600 dark:hover:text-primary-400 text-sm">
                                        Онгоинги
                                    </Link>
                                </li>
                            </ul>
                        </div>

                        <div>
                            <h3 className="text-base font-semibold text-gray-800 dark:text-gray-200 mb-3">Жанры</h3>
                            <ul className="space-y-2">
                                <li>
                                    <Link to="/catalog?genre=Боевик" className="text-gray-600 dark:text-gray-400 hover:text-primary-600 dark:hover:text-primary-400 text-sm">
                                        Боевик
                                    </Link>
                                </li>
                                <li>
                                    <Link to="/catalog?genre=Романтика" className="text-gray-600 dark:text-gray-400 hover:text-primary-600 dark:hover:text-primary-400 text-sm">
                                        Романтика
                                    </Link>
                                </li>
                                <li>
                                    <Link to="/catalog?genre=Фэнтези" className="text-gray-600 dark:text-gray-400 hover:text-primary-600 dark:hover:text-primary-400 text-sm">
                                        Фэнтези
                                    </Link>
                                </li>
                                <li>
                                    <Link to="/catalog?genre=Комедия" className="text-gray-600 dark:text-gray-400 hover:text-primary-600 dark:hover:text-primary-400 text-sm">
                                        Комедия
                                    </Link>
                                </li>
                            </ul>
                        </div>

                        <div>
                            <h3 className="text-base font-semibold text-gray-800 dark:text-gray-200 mb-3">Аккаунт</h3>
                            <ul className="space-y-2">
                                <li>
                                    <Link to="/login" className="text-gray-600 dark:text-gray-400 hover:text-primary-600 dark:hover:text-primary-400 text-sm">
                                        Войти
                                    </Link>
                                </li>
                                <li>
                                    <Link to="/register" className="text-gray-600 dark:text-gray-400 hover:text-primary-600 dark:hover:text-primary-400 text-sm">
                                        Регистрация
                                    </Link>
                                </li>
                                <li>
                                    <Link to="/profile" className="text-gray-600 dark:text-gray-400 hover:text-primary-600 dark:hover:text-primary-400 text-sm">
                                        Профиль
                                    </Link>
                                </li>
                                <li>
                                    <Link to="/bookmarks" className="text-gray-600 dark:text-gray-400 hover:text-primary-600 dark:hover:text-primary-400 text-sm">
                                        Закладки
                                    </Link>
                                </li>
                            </ul>
                        </div>
                    </div>
                </div>

                <div className="mt-8 pt-6 border-t border-gray-200 dark:border-dark-700">
                    <div className="flex flex-col md:flex-row md:justify-between md:items-center">
                        <p className="text-sm text-gray-600 dark:text-gray-400">
                            &copy; {currentYear} MangaReader. Все права защищены.
                        </p>
                        <div className="mt-4 md:mt-0">
                            <p className="text-sm text-gray-600 dark:text-gray-400">
                                Разработано для чтения манги с любовью к аниме и манге.
                            </p>
                        </div>
                    </div>
                </div>
            </div>
        </footer>
    );
};

export default Footer;
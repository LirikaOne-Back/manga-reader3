import React, { useState, Fragment } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import {
    MagnifyingGlassIcon,
    BookOpenIcon,
    SunIcon,
    MoonIcon,
    UserIcon,
    BookmarkIcon,
    Bars3Icon,
    XMarkIcon
} from '@heroicons/react/24/outline';
import { Disclosure, Menu, Transition } from '@headlessui/react';

interface HeaderProps {
    isDarkMode: boolean;
    toggleDarkMode: () => void;
}

const Header: React.FC<HeaderProps> = ({ isDarkMode, toggleDarkMode }) => {
    const [searchQuery, setSearchQuery] = useState('');
    const navigate = useNavigate();

    // Имитация авторизации пользователя (в реальном приложении это будет из хука useAuth)
    const isAuthenticated = false; // localStorage.getItem('accessToken') !== null;

    // Обработчик отправки формы поиска
    const handleSearch = (e: React.FormEvent) => {
        e.preventDefault();
        if (searchQuery.trim()) {
            navigate(`/catalog?search=${encodeURIComponent(searchQuery.trim())}`);
        }
    };

    return (
        <header className="bg-white dark:bg-dark-800 shadow-md sticky top-0 z-50">
            <Disclosure as="nav" className="container mx-auto px-4">
                {({ open }) => (
                    <>
                        <div className="flex justify-between h-16">
                            {/* Логотип и название сайта */}
                            <div className="flex items-center">
                                <Link to="/" className="flex items-center">
                                    <BookOpenIcon className="h-8 w-8 text-primary-600" />
                                    <span className="ml-2 text-xl font-bold text-dark-900 dark:text-white">MangaReader</span>
                                </Link>
                            </div>

                            {/* Навигация для десктопов */}
                            <div className="hidden md:flex items-center space-x-4">
                                <Link
                                    to="/catalog"
                                    className="px-3 py-2 rounded-md text-gray-700 dark:text-gray-200 hover:bg-gray-100 dark:hover:bg-dark-700"
                                >
                                    Каталог
                                </Link>
                                <Link
                                    to="/catalog?sortBy=rating&sortDesc=true"
                                    className="px-3 py-2 rounded-md text-gray-700 dark:text-gray-200 hover:bg-gray-100 dark:hover:bg-dark-700"
                                >
                                    Популярное
                                </Link>
                                <Link
                                    to="/catalog?status=ongoing"
                                    className="px-3 py-2 rounded-md text-gray-700 dark:text-gray-200 hover:bg-gray-100 dark:hover:bg-dark-700"
                                >
                                    Онгоинги
                                </Link>
                            </div>

                            {/* Поиск для десктопов */}
                            <div className="hidden md:flex items-center">
                                <form onSubmit={handleSearch} className="relative mr-4">
                                    <input
                                        type="text"
                                        placeholder="Поиск манги..."
                                        className="w-48 lg:w-64 px-4 py-2 pl-10 rounded-full bg-gray-100 dark:bg-dark-700 focus:outline-none focus:ring-2 focus:ring-primary-500"
                                        value={searchQuery}
                                        onChange={(e) => setSearchQuery(e.target.value)}
                                    />
                                    <button
                                        type="submit"
                                        className="absolute left-0 top-0 mt-2.5 ml-3 text-gray-500 dark:text-gray-400"
                                    >
                                        <MagnifyingGlassIcon className="h-5 w-5" />
                                    </button>
                                </form>

                                {/* Переключатель темы */}
                                <button
                                    onClick={toggleDarkMode}
                                    className="p-2 rounded-full text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-dark-700"
                                    aria-label={isDarkMode ? 'Включить светлую тему' : 'Включить темную тему'}
                                >
                                    {isDarkMode ? (
                                        <SunIcon className="h-6 w-6" />
                                    ) : (
                                        <MoonIcon className="h-6 w-6" />
                                    )}
                                </button>

                                {/* Кнопки авторизации или меню пользователя */}
                                {isAuthenticated ? (
                                    <Menu as="div" className="relative ml-4">
                                        <Menu.Button className="flex rounded-full bg-gray-100 dark:bg-dark-700 p-1 text-gray-600 dark:text-gray-300 hover:ring-2 hover:ring-primary-500">
                                            <UserIcon className="h-6 w-6" />
                                        </Menu.Button>
                                        <Transition
                                            as={Fragment}
                                            enter="transition ease-out duration-100"
                                            enterFrom="transform opacity-0 scale-95"
                                            enterTo="transform opacity-100 scale-100"
                                            leave="transition ease-in duration-75"
                                            leaveFrom="transform opacity-100 scale-100"
                                            leaveTo="transform opacity-0 scale-95"
                                        >
                                            <Menu.Items className="absolute right-0 mt-2 w-48 origin-top-right rounded-md bg-white dark:bg-dark-700 py-1 shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none">
                                                <Menu.Item>
                                                    {({ active }) => (
                                                        <Link
                                                            to="/profile"
                                                            className={`${
                                                                active ? 'bg-gray-100 dark:bg-dark-600' : ''
                                                            } block px-4 py-2 text-sm text-gray-700 dark:text-gray-200`}
                                                        >
                                                            Профиль
                                                        </Link>
                                                    )}
                                                </Menu.Item>
                                                <Menu.Item>
                                                    {({ active }) => (
                                                        <Link
                                                            to="/bookmarks"
                                                            className={`${
                                                                active ? 'bg-gray-100 dark:bg-dark-600' : ''
                                                            } block px-4 py-2 text-sm text-gray-700 dark:text-gray-200`}
                                                        >
                                                            Закладки
                                                        </Link>
                                                    )}
                                                </Menu.Item>
                                                <Menu.Item>
                                                    {({ active }) => (
                                                        <button
                                                            className={`${
                                                                active ? 'bg-gray-100 dark:bg-dark-600' : ''
                                                            } block w-full text-left px-4 py-2 text-sm text-gray-700 dark:text-gray-200`}
                                                            // onClick={logout}
                                                        >
                                                            Выйти
                                                        </button>
                                                    )}
                                                </Menu.Item>
                                            </Menu.Items>
                                        </Transition>
                                    </Menu>
                                ) : (
                                    <div className="flex items-center ml-4 space-x-2">
                                        <Link
                                            to="/login"
                                            className="px-4 py-2 rounded-md text-primary-600 dark:text-primary-400 font-medium hover:bg-primary-50 dark:hover:bg-primary-900/20"
                                        >
                                            Войти
                                        </Link>
                                        <Link
                                            to="/register"
                                            className="px-4 py-2 rounded-md bg-primary-600 hover:bg-primary-700 text-white font-medium"
                                        >
                                            Регистрация
                                        </Link>
                                    </div>
                                )}
                            </div>

                            {/* Кнопка мобильного меню */}
                            <div className="flex items-center md:hidden">
                                <button
                                    onClick={toggleDarkMode}
                                    className="p-2 rounded-full text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-dark-700 mr-2"
                                >
                                    {isDarkMode ? (
                                        <SunIcon className="h-6 w-6" />
                                    ) : (
                                        <MoonIcon className="h-6 w-6" />
                                    )}
                                </button>

                                <Disclosure.Button className="p-2 rounded-md text-gray-600 dark:text-gray-200 hover:bg-gray-100 dark:hover:bg-dark-700">
                                    {open ? (
                                        <XMarkIcon className="h-6 w-6" />
                                    ) : (
                                        <Bars3Icon className="h-6 w-6" />
                                    )}
                                </Disclosure.Button>
                            </div>
                        </div>

                        {/* Мобильное меню */}
                        <Disclosure.Panel className="md:hidden">
                            <div className="px-2 pt-2 pb-3 space-y-1">
                                <form onSubmit={handleSearch} className="relative mb-3">
                                    <input
                                        type="text"
                                        placeholder="Поиск манги..."
                                        className="w-full px-4 py-2 pl-10 rounded-full bg-gray-100 dark:bg-dark-700 focus:outline-none focus:ring-2 focus:ring-primary-500"
                                        value={searchQuery}
                                        onChange={(e) => setSearchQuery(e.target.value)}
                                    />
                                    <button
                                        type="submit"
                                        className="absolute left-0 top-0 mt-2.5 ml-3 text-gray-500 dark:text-gray-400"
                                    >
                                        <MagnifyingGlassIcon className="h-5 w-5" />
                                    </button>
                                </form>

                                <Disclosure.Button
                                    as={Link}
                                    to="/catalog"
                                    className="block px-3 py-2 rounded-md text-gray-700 dark:text-gray-200 hover:bg-gray-100 dark:hover:bg-dark-700"
                                >
                                    Каталог
                                </Disclosure.Button>

                                <Disclosure.Button
                                    as={Link}
                                    to="/catalog?sortBy=rating&sortDesc=true"
                                    className="block px-3 py-2 rounded-md text-gray-700 dark:text-gray-200 hover:bg-gray-100 dark:hover:bg-dark-700"
                                >
                                    Популярное
                                </Disclosure.Button>

                                <Disclosure.Button
                                    as={Link}
                                    to="/catalog?status=ongoing"
                                    className="block px-3 py-2 rounded-md text-gray-700 dark:text-gray-200 hover:bg-gray-100 dark:hover:bg-dark-700"
                                >
                                    Онгоинги
                                </Disclosure.Button>

                                {isAuthenticated ? (
                                    <>
                                        <Disclosure.Button
                                            as={Link}
                                            to="/profile"
                                            className="block px-3 py-2 rounded-md text-gray-700 dark:text-gray-200 hover:bg-gray-100 dark:hover:bg-dark-700"
                                        >
                                            <UserIcon className="h-5 w-5 inline mr-2" />
                                            Профиль
                                        </Disclosure.Button>

                                        <Disclosure.Button
                                            as={Link}
                                            to="/bookmarks"
                                            className="block px-3 py-2 rounded-md text-gray-700 dark:text-gray-200 hover:bg-gray-100 dark:hover:bg-dark-700"
                                        >
                                            <BookmarkIcon className="h-5 w-5 inline mr-2" />
                                            Закладки
                                        </Disclosure.Button>

                                        <button
                                            className="block w-full text-left px-3 py-2 rounded-md text-gray-700 dark:text-gray-200 hover:bg-gray-100 dark:hover:bg-dark-700"
                                            // onClick={logout}
                                        >
                                            Выйти
                                        </button>
                                    </>
                                ) : (
                                    <div className="pt-4 pb-3 border-t border-gray-200 dark:border-dark-600">
                                        <div className="flex items-center px-3 space-x-3">
                                            <Link
                                                to="/login"
                                                className="flex-1 px-4 py-2 rounded-md text-center text-primary-600 dark:text-primary-400 border border-primary-600 dark:border-primary-400 font-medium"
                                            >
                                                Войти
                                            </Link>
                                            <Link
                                                to="/register"
                                                className="flex-1 px-4 py-2 rounded-md text-center bg-primary-600 hover:bg-primary-700 text-white font-medium"
                                            >
                                                Регистрация
                                            </Link>
                                        </div>
                                    </div>
                                )}
                            </div>
                        </Disclosure.Panel>
                    </>
                )}
            </Disclosure>
        </header>
    );
};

export default Header;
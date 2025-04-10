import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { Tab } from '@headlessui/react';
import { UserIcon, BookmarkIcon, ClockIcon, CogIcon } from '@heroicons/react/24/outline';
import MangaCard from '../components/manga/MangaCard';
import { Manga } from '../types/manga.types';
import MangaService from '../services/manga.service';

interface ProfileProps {
    activeTab?: string;
}

const Profile: React.FC<ProfileProps> = ({ activeTab = 'profile' }) => {
    const [selectedTab, setSelectedTab] = useState(0);
    const [bookmarks, setBookmarks] = useState<Manga[]>([]);
    const [readHistory, setReadHistory] = useState<any[]>([]);
    const [loading, setLoading] = useState(true);

    // Определяем индекс активной вкладки на основе пропса
    useEffect(() => {
        switch (activeTab) {
            case 'bookmarks':
                setSelectedTab(1);
                break;
            case 'history':
                setSelectedTab(2);
                break;
            case 'settings':
                setSelectedTab(3);
                break;
            default:
                setSelectedTab(0);
        }
    }, [activeTab]);

    // Загружаем данные для профиля
    useEffect(() => {
        const fetchUserData = async () => {
            setLoading(true);

            try {
                // Загружаем закладки
                const bookmarksData = await MangaService.getBookmarks();
                setBookmarks(bookmarksData);

                // Загружаем историю чтения
                const historyData = await MangaService.getReadHistory();
                setReadHistory(historyData);
            } catch (error) {
                console.error('Failed to fetch user data', error);
            } finally {
                setLoading(false);
            }
        };

        fetchUserData();
    }, []);

    // Заглушка для данных пользователя (в реальном приложении будет из API)
    const user = {
        username: 'Пользователь',
        email: 'user@example.com',
        avatarUrl: null,
        joinDate: '2023-01-15',
    };

    // Функция для отображения даты
    const formatDate = (dateString: string) => {
        return new Date(dateString).toLocaleDateString('ru-RU', {
            year: 'numeric',
            month: 'long',
            day: 'numeric',
        });
    };

    return (
        <div className="container mx-auto px-4 py-8">
            <div className="bg-white dark:bg-dark-800 rounded-lg shadow-md overflow-hidden">
                <div className="p-6 sm:p-8 bg-gradient-to-r from-primary-600 to-secondary-600 text-white">
                    <div className="flex flex-col sm:flex-row items-center sm:items-start gap-6">
                        <div className="w-24 h-24 sm:w-32 sm:h-32 bg-white dark:bg-dark-700 rounded-full flex items-center justify-center text-primary-600 dark:text-primary-400 text-4xl font-bold shadow-lg">
                            {user.avatarUrl ? (
                                <img
                                    src={user.avatarUrl}
                                    alt={user.username}
                                    className="w-full h-full rounded-full object-cover"
                                />
                            ) : (
                                user.username.charAt(0).toUpperCase()
                            )}
                        </div>

                        <div className="text-center sm:text-left">
                            <h1 className="text-2xl sm:text-3xl font-bold">{user.username}</h1>
                            <p className="text-primary-100 mt-1">{user.email}</p>
                            <p className="text-primary-200 mt-2">
                                На сайте с {formatDate(user.joinDate)}
                            </p>
                        </div>
                    </div>
                </div>

                <Tab.Group selectedIndex={selectedTab} onChange={setSelectedTab}>
                    <Tab.List className="flex border-b border-gray-200 dark:border-dark-700">
                        <Tab
                            className={({ selected }) =>
                                `flex items-center py-4 px-6 font-medium text-sm focus:outline-none whitespace-nowrap
                ${selected
                                    ? 'text-primary-600 dark:text-primary-400 border-b-2 border-primary-600 dark:border-primary-400'
                                    : 'text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300'
                                }`
                            }
                        >
                            <UserIcon className="h-5 w-5 mr-2" />
                            Профиль
                        </Tab>

                        <Tab
                            className={({ selected }) =>
                                `flex items-center py-4 px-6 font-medium text-sm focus:outline-none whitespace-nowrap
                ${selected
                                    ? 'text-primary-600 dark:text-primary-400 border-b-2 border-primary-600 dark:border-primary-400'
                                    : 'text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300'
                                }`
                            }
                        >
                            <BookmarkIcon className="h-5 w-5 mr-2" />
                            Закладки
                        </Tab>

                        <Tab
                            className={({ selected }) =>
                                `flex items-center py-4 px-6 font-medium text-sm focus:outline-none whitespace-nowrap
                ${selected
                                    ? 'text-primary-600 dark:text-primary-400 border-b-2 border-primary-600 dark:border-primary-400'
                                    : 'text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300'
                                }`
                            }
                        >
                            <ClockIcon className="h-5 w-5 mr-2" />
                            История
                        </Tab>

                        <Tab
                            className={({ selected }) =>
                                `flex items-center py-4 px-6 font-medium text-sm focus:outline-none whitespace-nowrap
                ${selected
                                    ? 'text-primary-600 dark:text-primary-400 border-b-2 border-primary-600 dark:border-primary-400'
                                    : 'text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300'
                                }`
                            }
                        >
                            <CogIcon className="h-5 w-5 mr-2" />
                            Настройки
                        </Tab>
                    </Tab.List>

                    <Tab.Panels className="p-6">
                        {/* Вкладка профиля */}
                        <Tab.Panel>
                            <div className="max-w-2xl mx-auto">
                                <h2 className="text-xl font-semibold mb-6">Информация о пользователе</h2>

                                <div className="bg-gray-50 dark:bg-dark-700 rounded-lg p-6">
                                    <div className="mb-4">
                                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Имя пользователя</label>
                                        <div className="text-gray-900 dark:text-gray-100">{user.username}</div>
                                    </div>

                                    <div className="mb-4">
                                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Email</label>
                                        <div className="text-gray-900 dark:text-gray-100">{user.email}</div>
                                    </div>

                                    <div className="mb-4">
                                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Дата регистрации</label>
                                        <div className="text-gray-900 dark:text-gray-100">{formatDate(user.joinDate)}</div>
                                    </div>

                                    <div className="pt-4 border-t border-gray-200 dark:border-dark-600">
                                        <button
                                            className="btn btn-primary w-full sm:w-auto"
                                        >
                                            Редактировать профиль
                                        </button>
                                    </div>
                                </div>
                            </div>
                        </Tab.Panel>

                        {/* Вкладка закладок */}
                        <Tab.Panel>
                            <h2 className="text-xl font-semibold mb-6">Мои закладки</h2>

                            {loading ? (
                                <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-4">
                                    {[...Array(10)].map((_, index) => (
                                        <div key={index} className="animate-pulse">
                                            <div className="bg-gray-200 dark:bg-dark-700 rounded-md aspect-[2/3]"></div>
                                            <div className="h-4 bg-gray-200 dark:bg-dark-700 rounded-md mt-2 w-3/4"></div>
                                            <div className="h-3 bg-gray-200 dark:bg-dark-700 rounded-md mt-2 w-1/2"></div>
                                        </div>
                                    ))}
                                </div>
                            ) : bookmarks.length > 0 ? (
                                <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-4">
                                    {bookmarks.map((manga) => (
                                        <MangaCard key={manga.id} manga={manga} />
                                    ))}
                                </div>
                            ) : (
                                <div className="text-center py-12 bg-gray-50 dark:bg-dark-700 rounded-lg">
                                    <BookmarkIcon className="mx-auto h-12 w-12 text-gray-400" />
                                    <h3 className="mt-2 text-lg font-medium text-gray-900 dark:text-gray-100">Нет закладок</h3>
                                    <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
                                        У вас пока нет сохраненных закладок. Добавьте мангу в закладки, чтобы быстро вернуться к ней позже.
                                    </p>
                                    <div className="mt-6">
                                        <Link to="/catalog" className="btn btn-primary">
                                            Перейти в каталог
                                        </Link>
                                    </div>
                                </div>
                            )}
                        </Tab.Panel>

                        {/* Вкладка истории */}
                        <Tab.Panel>
                            <h2 className="text-xl font-semibold mb-6">История чтения</h2>

                            {loading ? (
                                <div className="space-y-4">
                                    {[...Array(5)].map((_, index) => (
                                        <div key={index} className="animate-pulse bg-gray-50 dark:bg-dark-700 p-4 rounded-lg">
                                            <div className="flex items-center gap-4">
                                                <div className="bg-gray-200 dark:bg-dark-600 rounded-md h-16 w-12"></div>
                                                <div className="flex-1">
                                                    <div className="h-4 bg-gray-200 dark:bg-dark-600 rounded-md w-1/3 mb-2"></div>
                                                    <div className="h-3 bg-gray-200 dark:bg-dark-600 rounded-md w-1/4"></div>
                                                </div>
                                                <div className="h-3 bg-gray-200 dark:bg-dark-600 rounded-md w-20"></div>
                                            </div>
                                        </div>
                                    ))}
                                </div>
                            ) : readHistory.length > 0 ? (
                                <div className="space-y-4">
                                    {readHistory.map((item) => (
                                        <div key={item.id} className="bg-gray-50 dark:bg-dark-700 p-4 rounded-lg hover:bg-gray-100 dark:hover:bg-dark-600 transition-colors">
                                            <Link to={`/manga/${item.mangaId}/chapter/${item.chapterId}`} className="flex items-center gap-4">
                                                <div className="h-16 w-12 bg-gray-200 dark:bg-dark-600 rounded overflow-hidden">
                                                    <img
                                                        src={item.manga?.coverUrl}
                                                        alt={item.manga?.title}
                                                        className="h-full w-full object-cover"
                                                    />
                                                </div>
                                                <div className="flex-1">
                                                    <h3 className="font-medium text-gray-900 dark:text-gray-100">{item.manga?.title}</h3>
                                                    <p className="text-sm text-gray-500 dark:text-gray-400">
                                                        Глава {item.chapter?.number}: {item.chapter?.title}
                                                    </p>
                                                </div>
                                                <div className="text-sm text-gray-500 dark:text-gray-400">
                                                    {formatDate(item.readAt)}
                                                </div>
                                            </Link>
                                        </div>
                                    ))}
                                </div>
                            ) : (
                                <div className="text-center py-12 bg-gray-50 dark:bg-dark-700 rounded-lg">
                                    <ClockIcon className="mx-auto h-12 w-12 text-gray-400" />
                                    <h3 className="mt-2 text-lg font-medium text-gray-900 dark:text-gray-100">Нет истории чтения</h3>
                                    <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
                                        История чтения появится автоматически, когда вы начнете читать мангу.
                                    </p>
                                    <div className="mt-6">
                                        <Link to="/catalog" className="btn btn-primary">
                                            Перейти к чтению
                                        </Link>
                                    </div>
                                </div>
                            )}
                        </Tab.Panel>

                        {/* Вкладка настроек */}
                        <Tab.Panel>
                            <div className="max-w-2xl mx-auto">
                                <h2 className="text-xl font-semibold mb-6">Настройки профиля</h2>

                                <div className="space-y-8">
                                    {/* Общие настройки */}
                                    <div className="bg-gray-50 dark:bg-dark-700 rounded-lg p-6">
                                        <h3 className="text-lg font-medium mb-4">Общие настройки</h3>

                                        <div className="space-y-4">
                                            <div>
                                                <label htmlFor="username" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                                                    Имя пользователя
                                                </label>
                                                <input
                                                    type="text"
                                                    id="username"
                                                    className="input"
                                                    defaultValue={user.username}
                                                />
                                            </div>

                                            <div>
                                                <label htmlFor="email" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                                                    Email
                                                </label>
                                                <input
                                                    type="email"
                                                    id="email"
                                                    className="input"
                                                    defaultValue={user.email}
                                                />
                                            </div>

                                            <button className="btn btn-primary">
                                                Сохранить изменения
                                            </button>
                                        </div>
                                    </div>

                                    {/* Изменение пароля */}
                                    <div className="bg-gray-50 dark:bg-dark-700 rounded-lg p-6">
                                        <h3 className="text-lg font-medium mb-4">Изменение пароля</h3>

                                        <div className="space-y-4">
                                            <div>
                                                <label htmlFor="current-password" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                                                    Текущий пароль
                                                </label>
                                                <input
                                                    type="password"
                                                    id="current-password"
                                                    className="input"
                                                />
                                            </div>

                                            <div>
                                                <label htmlFor="new-password" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                                                    Новый пароль
                                                </label>
                                                <input
                                                    type="password"
                                                    id="new-password"
                                                    className="input"
                                                />
                                            </div>

                                            <div>
                                                <label htmlFor="confirm-password" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                                                    Подтверждение пароля
                                                </label>
                                                <input
                                                    type="password"
                                                    id="confirm-password"
                                                    className="input"
                                                />
                                            </div>

                                            <button className="btn btn-primary">
                                                Изменить пароль
                                            </button>
                                        </div>
                                    </div>

                                    {/* Удаление аккаунта */}
                                    <div className="bg-red-50 dark:bg-red-900/20 rounded-lg p-6 border border-red-200 dark:border-red-800">
                                        <h3 className="text-lg font-medium text-red-800 dark:text-red-300 mb-4">Удаление аккаунта</h3>

                                        <p className="text-red-600 dark:text-red-400 mb-4">
                                            Удаление аккаунта - это необратимое действие. Вся ваша информация, включая закладки и историю чтения, будет удалена.
                                        </p>

                                        <button className="bg-red-600 hover:bg-red-700 text-white font-medium py-2 px-4 rounded-md">
                                            Удалить аккаунт
                                        </button>
                                    </div>
                                </div>
                            </div>
                        </Tab.Panel>
                    </Tab.Panels>
                </Tab.Group>
            </div>
        </div>
    );
};

export default Profile;
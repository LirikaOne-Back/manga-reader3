import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { Manga, Genre } from '../types/manga.types';
import MangaCard from '../components/manga/MangaCard';
import MangaService from '../services/manga.service';
import { ArrowRightIcon, FireIcon, ClockIcon, SparklesIcon } from '@heroicons/react/24/solid';

const Home: React.FC = () => {
    const [popularManga, setPopularManga] = useState<Manga[]>([]);
    const [recentUpdates, setRecentUpdates] = useState<Manga[]>([]);
    const [genres, setGenres] = useState<Genre[]>([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const fetchHomeData = async () => {
            try {
                setLoading(true);
                // Загружаем популярную мангу (сортировка по рейтингу)
                const popularResponse = await MangaService.getAll({
                    sortBy: 'rating',
                    sortDesc: true,
                    page: 1,
                    pageSize: 6
                });
                setPopularManga(popularResponse.data);

                // Загружаем последние обновления (сортировка по дате)
                const recentResponse = await MangaService.getAll({
                    sortBy: 'date',
                    sortDesc: true,
                    page: 1,
                    pageSize: 6
                });
                setRecentUpdates(recentResponse.data);

                // Загружаем жанры
                const genresData = await MangaService.getGenres();
                setGenres(genresData);
            } catch (error) {
                console.error('Failed to fetch home data', error);
            } finally {
                setLoading(false);
            }
        };

        fetchHomeData();
    }, []);

    // Компонент заглушка при загрузке
    if (loading) {
        return (
            <div className="container mx-auto px-4 py-8">
                <div className="animate-pulse">
                    <div className="h-10 bg-gray-200 dark:bg-dark-700 rounded-md mb-8 w-1/3"></div>
                    <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-4">
                        {[...Array(6)].map((_, index) => (
                            <div key={index} className="bg-gray-200 dark:bg-dark-700 rounded-md aspect-[2/3]"></div>
                        ))}
                    </div>
                    <div className="h-10 bg-gray-200 dark:bg-dark-700 rounded-md my-8 w-1/3"></div>
                    <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-4">
                        {[...Array(6)].map((_, index) => (
                            <div key={index} className="bg-gray-200 dark:bg-dark-700 rounded-md aspect-[2/3]"></div>
                        ))}
                    </div>
                </div>
            </div>
        );
    }

    return (
        <div className="container mx-auto px-4 py-8">
            {/* Баннер / Герой секция */}
            <div className="relative rounded-xl overflow-hidden bg-gradient-to-r from-primary-600 to-secondary-600 mb-12 p-6 md:p-8 text-white">
                <div className="relative z-10 max-w-2xl">
                    <h1 className="text-4xl md:text-5xl font-bold mb-4">Читайте любимую мангу онлайн</h1>
                    <p className="text-lg md:text-xl mb-6 opacity-90">
                        Тысячи манги и манхвы с регулярными обновлениями. Находите новые истории и следите за любимыми сериями.
                    </p>
                    <Link
                        to="/catalog"
                        className="inline-flex items-center px-6 py-3 rounded-full bg-white text-primary-700 font-semibold shadow-lg hover:bg-gray-100 transition-colors"
                    >
                        Перейти в каталог
                        <ArrowRightIcon className="ml-2 h-5 w-5" />
                    </Link>
                </div>
                {/* Декоративный элемент */}
                <div className="absolute right-0 top-0 h-full w-1/3 opacity-20">
                    <svg viewBox="0 0 100 100" preserveAspectRatio="none" className="h-full w-full">
                        <path d="M0 0 L100 0 L100 100 L0 100 Z" fill="white" />
                    </svg>
                </div>
            </div>

            {/* Популярная манга */}
            <section className="mb-12">
                <div className="flex items-center justify-between mb-6">
                    <h2 className="text-2xl font-bold flex items-center">
                        <FireIcon className="h-6 w-6 text-red-500 mr-2" />
                        Популярное
                    </h2>
                    <Link
                        to="/catalog?sortBy=rating&sortDesc=true"
                        className="text-primary-600 hover:text-primary-700 dark:text-primary-400 dark:hover:text-primary-300 font-medium flex items-center"
                    >
                        Показать все
                        <ArrowRightIcon className="ml-1 h-4 w-4" />
                    </Link>
                </div>
                <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-4">
                    {popularManga.map((manga) => (
                        <MangaCard key={manga.id} manga={manga} />
                    ))}
                </div>
            </section>

            {/* Последние обновления */}
            <section className="mb-12">
                <div className="flex items-center justify-between mb-6">
                    <h2 className="text-2xl font-bold flex items-center">
                        <ClockIcon className="h-6 w-6 text-blue-500 mr-2" />
                        Последние обновления
                    </h2>
                    <Link
                        to="/catalog?sortBy=date&sortDesc=true"
                        className="text-primary-600 hover:text-primary-700 dark:text-primary-400 dark:hover:text-primary-300 font-medium flex items-center"
                    >
                        Показать все
                        <ArrowRightIcon className="ml-1 h-4 w-4" />
                    </Link>
                </div>
                <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-4">
                    {recentUpdates.map((manga) => (
                        <MangaCard key={manga.id} manga={manga} />
                    ))}
                </div>
            </section>

            {/* Популярные жанры */}
            <section>
                <div className="flex items-center mb-6">
                    <h2 className="text-2xl font-bold flex items-center">
                        <SparklesIcon className="h-6 w-6 text-yellow-500 mr-2" />
                        Популярные жанры
                    </h2>
                </div>
                <div className="flex flex-wrap gap-2">
                    {genres.slice(0, 12).map((genre) => (
                        <Link
                            key={genre.id}
                            to={`/catalog?genre=${genre.name}`}
                            className="inline-block rounded-full bg-gray-100 hover:bg-gray-200 dark:bg-dark-700 dark:hover:bg-dark-600 px-4 py-2 text-sm font-medium transition-colors"
                        >
                            {genre.name}
                        </Link>
                    ))}
                </div>
            </section>
        </div>
    );
};

export default Home;
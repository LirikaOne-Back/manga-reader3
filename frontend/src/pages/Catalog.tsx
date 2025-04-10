import React, { useEffect, useState } from 'react';
import { useSearchParams } from 'react-router-dom';
import MangaCard from '../components/manga/MangaCard';
import { Manga, Genre, MangaFilter } from '../types/manga.types';
import MangaService from '../services/manga.service';
import { AdjustmentsHorizontalIcon, MagnifyingGlassIcon, XMarkIcon } from '@heroicons/react/24/outline';

const Catalog: React.FC = () => {
    const [searchParams, setSearchParams] = useSearchParams();
    const [mangas, setMangas] = useState<Manga[]>([]);
    const [genres, setGenres] = useState<Genre[]>([]);
    const [totalItems, setTotalItems] = useState(0);
    const [loading, setLoading] = useState(true);
    const [showFilters, setShowFilters] = useState(false);

    // Получаем параметры из URL
    const genre = searchParams.get('genre') || '';
    const status = searchParams.get('status') || '';
    const sortBy = searchParams.get('sortBy') || 'title';
    const sortDesc = searchParams.get('sortDesc') === 'true';
    const search = searchParams.get('search') || '';
    const page = parseInt(searchParams.get('page') || '1', 10);
    const pageSize = parseInt(searchParams.get('pageSize') || '24', 10);

    useEffect(() => {
        // Загрузка жанров при первом рендере
        const fetchGenres = async () => {
            try {
                const genresData = await MangaService.getGenres();
                setGenres(genresData);
            } catch (error) {
                console.error('Failed to fetch genres', error);
            }
        };

        fetchGenres();
    }, []);

    useEffect(() => {
        // Загрузка манги при изменении параметров фильтрации
        const fetchMangas = async () => {
            setLoading(true);

            try {
                const filter: MangaFilter = {
                    genre,
                    status,
                    sortBy,
                    sortDesc,
                    search,
                    page,
                    pageSize
                };

                const response = await MangaService.getAll(filter);
                setMangas(response.data);
                setTotalItems(response.total);
            } catch (error) {
                console.error('Failed to fetch mangas', error);
            } finally {
                setLoading(false);
            }
        };

        fetchMangas();
    }, [genre, status, sortBy, sortDesc, search, page, pageSize]);

    // Функция для обновления параметров фильтрации
    const updateFilter = (name: string, value: string | boolean) => {
        // Возвращаемся на первую страницу при изменении фильтров
        if (name !== 'page') {
            searchParams.set('page', '1');
        }

        if (value === '' || value === false) {
            searchParams.delete(name);
        } else {
            searchParams.set(name, String(value));
        }

        setSearchParams(searchParams);
    };

    // Обработчик изменения поиска
    const handleSearchChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const value = e.target.value;
        if (value) {
            searchParams.set('search', value);
        } else {
            searchParams.delete('search');
        }
        searchParams.set('page', '1');
        setSearchParams(searchParams);
    };

    // Обработчик очистки всех фильтров
    const clearAllFilters = () => {
        setSearchParams(new URLSearchParams({ page: '1', pageSize: String(pageSize) }));
    };

    // Вычисляем общее количество страниц
    const totalPages = Math.ceil(totalItems / pageSize);

    return (
        <div className="container mx-auto px-4 py-8">
            <div className="flex justify-between items-center mb-6">
                <h1 className="text-3xl font-bold">Каталог манги</h1>

                <button
                    onClick={() => setShowFilters(!showFilters)}
                    className="md:hidden btn btn-outline"
                >
                    <AdjustmentsHorizontalIcon className="h-5 w-5 mr-2" />
                    Фильтры
                </button>
            </div>

            <div className="flex flex-col md:flex-row gap-6">
                {/* Сайдбар с фильтрами - на мобильных скрыт по умолчанию */}
                <div className={`md:w-1/4 lg:w-1/5 ${showFilters ? 'block' : 'hidden md:block'}`}>
                    <div className="bg-white dark:bg-dark-800 rounded-lg shadow-md p-4 sticky top-20">
                        <div className="flex justify-between items-center mb-4">
                            <h2 className="text-xl font-semibold">Фильтры</h2>
                            {(genre || status || sortBy !== 'title' || sortDesc || search) && (
                                <button
                                    onClick={clearAllFilters}
                                    className="text-sm text-primary-600 hover:text-primary-700 dark:text-primary-400 flex items-center"
                                >
                                    <XMarkIcon className="h-4 w-4 mr-1" />
                                    Сбросить все
                                </button>
                            )}
                        </div>

                        {/* Поиск */}
                        <div className="mb-6">
                            <label htmlFor="search" className="block text-sm font-medium mb-2">
                                Поиск
                            </label>
                            <div className="relative">
                                <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                                    <MagnifyingGlassIcon className="h-5 w-5 text-gray-400" />
                                </div>
                                <input
                                    type="text"
                                    id="search"
                                    className="input pl-10"
                                    placeholder="Название манги..."
                                    value={search}
                                    onChange={handleSearchChange}
                                />
                            </div>
                        </div>

                        {/* Статус */}
                        <div className="mb-6">
                            <label className="block text-sm font-medium mb-2">
                                Статус
                            </label>
                            <select
                                className="input"
                                value={status}
                                onChange={(e) => updateFilter('status', e.target.value)}
                            >
                                <option value="">Все статусы</option>
                                <option value="ongoing">Выходит</option>
                                <option value="completed">Завершён</option>
                                <option value="hiatus">На паузе</option>
                            </select>
                        </div>

                        {/* Сортировка */}
                        <div className="mb-6">
                            <label className="block text-sm font-medium mb-2">
                                Сортировка
                            </label>
                            <select
                                className="input"
                                value={sortBy}
                                onChange={(e) => updateFilter('sortBy', e.target.value)}
                            >
                                <option value="title">По названию</option>
                                <option value="rating">По рейтингу</option>
                                <option value="date">По дате</option>
                            </select>
                        </div>

                        <div className="mb-6">
                            <label className="flex items-center">
                                <input
                                    type="checkbox"
                                    className="rounded border-gray-300 text-primary-600 focus:ring-primary-500"
                                    checked={sortDesc}
                                    onChange={(e) => updateFilter('sortDesc', e.target.checked)}
                                />
                                <span className="ml-2 text-sm">По убыванию</span>
                            </label>
                        </div>

                        {/* Жанры */}
                        <div>
                            <label className="block text-sm font-medium mb-2">
                                Жанры
                            </label>
                            <div className="max-h-60 overflow-y-auto space-y-2 pr-2">
                                {genres.map((g) => (
                                    <label key={g.id} className="flex items-center">
                                        <input
                                            type="radio"
                                            name="genre"
                                            className="rounded-full border-gray-300 text-primary-600 focus:ring-primary-500"
                                            checked={genre === g.name}
                                            onChange={() => updateFilter('genre', g.name)}
                                        />
                                        <span className="ml-2 text-sm">{g.name}</span>
                                    </label>
                                ))}
                                {genre && (
                                    <button
                                        className="text-sm text-primary-600 hover:text-primary-700 dark:text-primary-400 mt-2"
                                        onClick={() => updateFilter('genre', '')}
                                    >
                                        Сбросить жанр
                                    </button>
                                )}
                            </div>
                        </div>
                    </div>
                </div>

                {/* Основная область с мангой */}
                <div className="md:w-3/4 lg:w-4/5">
                    {/* Результаты и количество манги */}
                    <div className="mb-6 flex justify-between items-center">
                        <div>
              <span className="text-gray-600 dark:text-gray-400">
                Найдено: <span className="font-semibold">{totalItems}</span>
                  {totalItems === 1 ? ' манга' :
                      totalItems > 1 && totalItems < 5 ? ' манги' : ' манг'}
              </span>
                        </div>

                        <div>
                            <select
                                className="input w-auto"
                                value={pageSize}
                                onChange={(e) => updateFilter('pageSize', e.target.value)}
                            >
                                <option value="12">12 на странице</option>
                                <option value="24">24 на странице</option>
                                <option value="48">48 на странице</option>
                            </select>
                        </div>
                    </div>

                    {/* Сетка манги */}
                    {loading ? (
                        <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6 gap-4">
                            {[...Array(pageSize)].map((_, index) => (
                                <div key={index} className="animate-pulse">
                                    <div className="bg-gray-200 dark:bg-dark-700 rounded-md aspect-[2/3]"></div>
                                    <div className="h-4 bg-gray-200 dark:bg-dark-700 rounded-md mt-2 w-3/4"></div>
                                    <div className="h-3 bg-gray-200 dark:bg-dark-700 rounded-md mt-2 w-1/2"></div>
                                </div>
                            ))}
                        </div>
                    ) : mangas.length > 0 ? (
                        <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6 gap-4">
                            {mangas.map((manga) => (
                                <MangaCard key={manga.id} manga={manga} />
                            ))}
                        </div>
                    ) : (
                        <div className="text-center py-12 bg-gray-50 dark:bg-dark-800 rounded-lg">
                            <h3 className="text-xl font-semibold mb-2">Ничего не найдено</h3>
                            <p className="text-gray-600 dark:text-gray-400 mb-6">
                                Попробуйте изменить параметры поиска или фильтрации
                            </p>
                            <button
                                onClick={clearAllFilters}
                                className="btn btn-primary"
                            >
                                Сбросить все фильтры
                            </button>
                        </div>
                    )}

                    {/* Пагинация */}
                    {totalPages > 1 && (
                        <div className="mt-8 flex justify-center">
                            <nav className="flex items-center justify-between">
                                <button
                                    onClick={() => updateFilter('page', String(page - 1))}
                                    disabled={page === 1}
                                    className={`btn ${page === 1 ? 'opacity-50 cursor-not-allowed' : 'btn-outline'}`}
                                >
                                    Предыдущая
                                </button>

                                <div className="mx-4 flex items-center">
                                    {/* Рендерим только ближайшие страницы для экономии места */}
                                    {[...Array(totalPages)].map((_, index) => {
                                        const pageNumber = index + 1;
                                        // Показываем первую, последнюю и страницы вокруг текущей
                                        if (
                                            pageNumber === 1 ||
                                            pageNumber === totalPages ||
                                            (pageNumber >= page - 1 && pageNumber <= page + 1)
                                        ) {
                                            return (
                                                <button
                                                    key={pageNumber}
                                                    onClick={() => updateFilter('page', String(pageNumber))}
                                                    className={`mx-1 px-3 py-1 rounded-md ${
                                                        pageNumber === page
                                                            ? 'bg-primary-600 text-white'
                                                            : 'bg-gray-100 dark:bg-dark-700 hover:bg-gray-200 dark:hover:bg-dark-600'
                                                    }`}
                                                >
                                                    {pageNumber}
                                                </button>
                                            );
                                        } else if (
                                            (pageNumber === 2 && page > 3) ||
                                            (pageNumber === totalPages - 1 && page < totalPages - 2)
                                        ) {
                                            // Показываем многоточие
                                            return <span key={pageNumber} className="mx-1">...</span>;
                                        }
                                        return null;
                                    })}
                                </div>

                                <button
                                    onClick={() => updateFilter('page', String(page + 1))}
                                    disabled={page === totalPages}
                                    className={`btn ${page === totalPages ? 'opacity-50 cursor-not-allowed' : 'btn-outline'}`}
                                >
                                    Следующая
                                </button>
                            </nav>
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
};

export default Catalog;
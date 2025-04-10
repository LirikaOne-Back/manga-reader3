import React, { useEffect, useState } from 'react';
import { useParams, Link } from 'react-router-dom';
import { Manga, Chapter } from '../types/manga.types';
import MangaService from '../services/manga.service';
import {
    StarIcon,
    BookmarkIcon,
    BookOpenIcon,
    ClockIcon,
    CalendarIcon,
    UserIcon,
    TagIcon,
    ChartBarIcon
} from '@heroicons/react/24/outline';
import { BookmarkIcon as BookmarkSolidIcon } from '@heroicons/react/24/solid';

// Компонент для информации о манге и списка глав
const MangaDetails: React.FC = () => {
    const { id } = useParams<{ id: string }>();
    const [manga, setManga] = useState<Manga | null>(null);
    const [chapters, setChapters] = useState<Chapter[]>([]);
    const [loading, setLoading] = useState(true);
    const [bookmarked, setBookmarked] = useState(false);

    useEffect(() => {
        const fetchMangaDetails = async () => {
            if (!id) return;

            try {
                setLoading(true);

                // Получаем информацию о манге
                const mangaId = parseInt(id, 10);
                const mangaData = await MangaService.getById(mangaId);
                setManga(mangaData);

                // Получаем список глав
                const chaptersData = await MangaService.getChapters(mangaId);
                setChapters(chaptersData);

                // Проверяем, добавлена ли манга в закладки
                try {
                    const bookmarks = await MangaService.getBookmarks();
                    setBookmarked(bookmarks.some(b => b.id === mangaId));
                } catch (error) {
                    // Если пользователь не авторизован, просто игнорируем ошибку
                    console.log('Failed to load bookmarks, user may not be logged in');
                }
            } catch (error) {
                console.error('Failed to fetch manga details', error);
            } finally {
                setLoading(false);
            }
        };

        fetchMangaDetails();
    }, [id]);

    // Функция для форматирования даты
    const formatDate = (dateString: string) => {
        return new Date(dateString).toLocaleDateString('ru-RU', {
            year: 'numeric',
            month: 'long',
            day: 'numeric',
        });
    };

    // Функция для добавления/удаления закладки
    const toggleBookmark = async () => {
        if (!manga) return;

        try {
            if (bookmarked) {
                await MangaService.removeBookmark(manga.id);
            } else {
                await MangaService.addBookmark(manga.id);
            }
            setBookmarked(!bookmarked);
        } catch (error) {
            console.error('Failed to toggle bookmark', error);
        }
    };

    // Функция для получения текста статуса
    const getStatusText = (status: string): string => {
        switch (status) {
            case 'ongoing':
                return 'Выходит';
            case 'completed':
                return 'Завершён';
            case 'hiatus':
                return 'На паузе';
            default:
                return status;
        }
    };

    // Компонент-заглушка при загрузке
    if (loading) {
        return (
            <div className="container mx-auto px-4 py-8">
                <div className="animate-pulse">
                    <div className="flex flex-col md:flex-row gap-8">
                        <div className="md:w-1/3 lg:w-1/4 h-96 bg-gray-200 dark:bg-dark-700 rounded-lg"></div>
                        <div className="md:w-2/3 lg:w-3/4 space-y-4">
                            <div className="h-12 bg-gray-200 dark:bg-dark-700 rounded-md w-3/4"></div>
                            <div className="h-8 bg-gray-200 dark:bg-dark-700 rounded-md w-1/2"></div>
                            <div className="h-24 bg-gray-200 dark:bg-dark-700 rounded-md"></div>
                            <div className="flex flex-wrap gap-2">
                                {[...Array(5)].map((_, i) => (
                                    <div key={i} className="h-8 w-20 bg-gray-200 dark:bg-dark-700 rounded-full"></div>
                                ))}
                            </div>
                        </div>
                    </div>
                    <div className="mt-12">
                        <div className="h-10 bg-gray-200 dark:bg-dark-700 rounded-md mb-6 w-1/4"></div>
                        <div className="space-y-3">
                            {[...Array(5)].map((_, i) => (
                                <div key={i} className="h-16 bg-gray-200 dark:bg-dark-700 rounded-md"></div>
                            ))}
                        </div>
                    </div>
                </div>
            </div>
        );
    }

    // Если манга не найдена
    if (!manga) {
        return (
            <div className="container mx-auto px-4 py-8">
                <div className="text-center py-12">
                    <h2 className="text-2xl font-bold text-gray-700 dark:text-gray-300 mb-4">Манга не найдена</h2>
                    <p className="mb-6">Возможно, манга была удалена или у вас неверная ссылка.</p>
                    <Link
                        to="/catalog"
                        className="btn btn-primary"
                    >
                        Вернуться в каталог
                    </Link>
                </div>
            </div>
        );
    }

    return (
        <div className="container mx-auto px-4 py-8">
            {/* Верхняя часть с информацией о манге */}
            <div className="flex flex-col md:flex-row gap-8 mb-12">
                {/* Левая колонка с обложкой */}
                <div className="md:w-1/3 lg:w-1/4 flex flex-col">
                    <div className="rounded-lg overflow-hidden shadow-lg mb-4">
                        <img
                            src={manga.coverUrl}
                            alt={manga.title}
                            className="w-full object-cover"
                        />
                    </div>

                    <div className="flex flex-col gap-2">
                        <button
                            onClick={toggleBookmark}
                            className={`btn ${bookmarked ? 'btn-primary' : 'btn-outline'} w-full flex items-center justify-center`}
                        >
                            {bookmarked ? (
                                <>
                                    <BookmarkSolidIcon className="h-5 w-5 mr-2" />
                                    В закладках
                                </>
                            ) : (
                                <>
                                    <BookmarkIcon className="h-5 w-5 mr-2" />
                                    Добавить в закладки
                                </>
                            )}
                        </button>

                        {chapters.length > 0 && (
                            <Link
                                to={`/manga/${manga.id}/chapter/${chapters[0].id}`}
                                className="btn btn-secondary w-full flex items-center justify-center"
                            >
                                <BookOpenIcon className="h-5 w-5 mr-2" />
                                Читать
                            </Link>
                        )}
                    </div>
                </div>

                {/* Правая колонка с информацией */}
                <div className="md:w-2/3 lg:w-3/4">
                    <h1 className="text-3xl font-bold mb-2">{manga.title}</h1>

                    {manga.alterTitle && (
                        <h2 className="text-xl text-gray-600 dark:text-gray-400 mb-4">{manga.alterTitle}</h2>
                    )}

                    <div className="flex items-center mb-4">
                        <StarIcon className="h-5 w-5 text-yellow-500 mr-1" />
                        <span className="text-lg font-semibold mr-4">{manga.rating.toFixed(1)}</span>

                        <span className={`inline-flex items-center rounded-full px-3 py-1 text-sm font-semibold ${
                            manga.status === 'ongoing' ? 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200' :
                                manga.status === 'completed' ? 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200' :
                                    'bg-amber-100 text-amber-800 dark:bg-amber-900 dark:text-amber-200'
                        }`}>
              {getStatusText(manga.status)}
            </span>
                    </div>

                    <p className="text-gray-700 dark:text-gray-300 mb-6 whitespace-pre-line">
                        {manga.description}
                    </p>

                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
                        <div className="flex items-center">
                            <CalendarIcon className="h-5 w-5 text-gray-500 dark:text-gray-400 mr-2" />
                            <span className="text-gray-700 dark:text-gray-300">Год выпуска: <span className="font-medium">{manga.year}</span></span>
                        </div>

                        <div className="flex items-center">
                            <UserIcon className="h-5 w-5 text-gray-500 dark:text-gray-400 mr-2" />
                            <span className="text-gray-700 dark:text-gray-300">Автор: <span className="font-medium">{manga.author}</span></span>
                        </div>

                        {manga.artist && (
                            <div className="flex items-center">
                                <UserIcon className="h-5 w-5 text-gray-500 dark:text-gray-400 mr-2" />
                                <span className="text-gray-700 dark:text-gray-300">Художник: <span className="font-medium">{manga.artist}</span></span>
                            </div>
                        )}

                        <div className="flex items-center">
                            <ChartBarIcon className="h-5 w-5 text-gray-500 dark:text-gray-400 mr-2" />
                            <span className="text-gray-700 dark:text-gray-300">Глав: <span className="font-medium">{chapters.length}</span></span>
                        </div>
                    </div>

                    <div className="flex items-center mb-1">
                        <TagIcon className="h-5 w-5 text-gray-500 dark:text-gray-400 mr-2" />
                        <span className="text-gray-700 dark:text-gray-300">Жанры:</span>
                    </div>

                    <div className="flex flex-wrap gap-2">
                        {manga.genres.map((genre) => (
                            <Link
                                key={genre.id}
                                to={`/catalog?genre=${genre.name}`}
                                className="inline-block rounded-full bg-gray-100 hover:bg-gray-200 dark:bg-dark-700 dark:hover:bg-dark-600 px-3 py-1 text-sm font-medium transition-colors"
                            >
                                {genre.name}
                            </Link>
                        ))}
                    </div>
                </div>
            </div>

            {/* Список глав */}
            <div>
                <h2 className="text-2xl font-bold mb-6 flex items-center">
                    <BookOpenIcon className="h-6 w-6 mr-2" />
                    Список глав
                </h2>

                {chapters.length === 0 ? (
                    <div className="text-center py-8 bg-gray-50 dark:bg-dark-800 rounded-lg">
                        <p className="text-gray-600 dark:text-gray-400">Главы этой манги еще не добавлены</p>
                    </div>
                ) : (
                    <div className="space-y-2 divide-y divide-gray-200 dark:divide-dark-700">
                        {chapters.map((chapter) => (
                            <Link
                                key={chapter.id}
                                to={`/manga/${manga.id}/chapter/${chapter.id}`}
                                className="block py-4 px-4 rounded-lg hover:bg-gray-50 dark:hover:bg-dark-800 transition-colors"
                            >
                                <div className="flex justify-between items-center">
                                    <div>
                                        <span className="font-semibold">Глава {chapter.number}</span>
                                        {chapter.title && <span className="ml-2 text-gray-600 dark:text-gray-400">{chapter.title}</span>}
                                    </div>

                                    <div className="flex items-center text-sm text-gray-500 dark:text-gray-400">
                                        <ClockIcon className="h-4 w-4 mr-1" />
                                        {formatDate(chapter.createdAt)}
                                    </div>
                                </div>
                            </Link>
                        ))}
                    </div>
                )}
            </div>
        </div>
    );
};

export default MangaDetails;
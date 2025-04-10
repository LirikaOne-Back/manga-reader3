import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import MangaReader from '../components/manga/MangaReader';
import MangaService from '../services/manga.service';
import { Chapter, Page } from '../types/manga.types';

const ChapterReader: React.FC = () => {
    const { mangaId, chapterId } = useParams<{ mangaId: string; chapterId: string }>();
    const [chapter, setChapter] = useState<Chapter | null>(null);
    const [pages, setPages] = useState<Page[]>([]);
    const [chapters, setChapters] = useState<Chapter[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchChapterData = async () => {
            if (!mangaId || !chapterId) {
                setError('Неверная ссылка');
                setLoading(false);
                return;
            }

            try {
                setLoading(true);
                setError(null);

                const mangaIdNum = parseInt(mangaId, 10);
                const chapterIdNum = parseInt(chapterId, 10);

                // Получаем информацию о главе
                const chapterData = await MangaService.getChapter(chapterIdNum);
                setChapter(chapterData);

                // Получаем список страниц
                const pagesData = await MangaService.getChapterPages(chapterIdNum);
                setPages(pagesData);

                // Получаем список всех глав для навигации
                const chaptersData = await MangaService.getChapters(mangaIdNum);
                // Сортируем главы по номеру в убывающем порядке (от новых к старым)
                setChapters(chaptersData.sort((a, b) => b.number - a.number));

            } catch (error) {
                console.error('Failed to fetch chapter data', error);
                setError('Не удалось загрузить главу. Попробуйте позже.');
            } finally {
                setLoading(false);
            }
        };

        fetchChapterData();
    }, [mangaId, chapterId]);

    // Заглушка при загрузке
    if (loading) {
        return (
            <div className="min-h-screen bg-gray-100 dark:bg-dark-900 flex justify-center items-center">
                <div className="text-center">
                    <div className="inline-block h-12 w-12 animate-spin rounded-full border-4 border-solid border-primary-500 border-r-transparent align-[-0.125em] motion-reduce:animate-[spin_1.5s_linear_infinite]" role="status">
                        <span className="sr-only">Загрузка...</span>
                    </div>
                    <p className="mt-2 text-gray-600 dark:text-gray-400">Загрузка главы...</p>
                </div>
            </div>
        );
    }

    // Обработка ошибок
    if (error || !chapter || !chapters.length || !pages.length) {
        return (
            <div className="min-h-screen bg-gray-100 dark:bg-dark-900 flex justify-center items-center">
                <div className="text-center max-w-md p-6 bg-white dark:bg-dark-800 rounded-lg shadow-md">
                    <h2 className="text-2xl font-bold text-red-600 mb-4">Ошибка</h2>
                    <p className="mb-6 text-gray-700 dark:text-gray-300">
                        {error || 'Не удалось загрузить главу. Возможно, глава была удалена или перемещена.'}
                    </p>
                    <button
                        onClick={() => window.history.back()}
                        className="btn btn-primary"
                    >
                        Вернуться назад
                    </button>
                </div>
            </div>
        );
    }

    return (
        <MangaReader
            mangaId={parseInt(mangaId!, 10)}
            chapterId={parseInt(chapterId!, 10)}
            chapter={chapter}
            chapters={chapters}
            pages={pages}
        />
    );
};

export default ChapterReader;
import React, { useState, useEffect, useRef, useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import {
    ArrowLeftIcon,
    ArrowRightIcon,
    Bars3Icon,
    ArrowsPointingOutIcon,
    ArrowsPointingInIcon,
    ChevronUpIcon,
    ChevronDownIcon,
    HomeIcon
} from '@heroicons/react/24/solid';
import { Page, Chapter } from '../../types/manga.types';
import MangaService from '../../services/manga.service';

interface MangaReaderProps {
    mangaId: number;
    chapterId: number;
    pages: Page[];
    chapter: Chapter;
    chapters: Chapter[];
}

// Режимы чтения
enum ReadingMode {
    PAGED = 'paged',
    VERTICAL = 'vertical',
}

// Компонент читалки манги
const MangaReader: React.FC<MangaReaderProps> = ({
                                                     mangaId,
                                                     chapterId,
                                                     pages,
                                                     chapter,
                                                     chapters
                                                 }) => {
    const navigate = useNavigate();
    const [currentPage, setCurrentPage] = useState(1);
    const [readingMode, setReadingMode] = useState<ReadingMode>(
        localStorage.getItem('readingMode') as ReadingMode || ReadingMode.VERTICAL
    );
    const [showControls, setShowControls] = useState(true);
    const [showChapters, setShowChapters] = useState(false);
    const [isFullscreen, setIsFullscreen] = useState(false);
    const readerRef = useRef<HTMLDivElement>(null);
    const controlsTimeoutRef = useRef<NodeJS.Timeout | null>(null);

    // Сохраняем прогресс чтения
    useEffect(() => {
        const saveProgress = async () => {
            try {
                await MangaService.saveReadProgress(chapterId, currentPage);
            } catch (error) {
                console.error('Failed to save reading progress', error);
            }
        };

        // Сохраняем прогресс с задержкой, чтобы не отправлять запросы при быстром пролистывании
        const timeoutId = setTimeout(saveProgress, 2000);
        return () => clearTimeout(timeoutId);
    }, [chapterId, currentPage]);

    // Сохраняем режим чтения
    useEffect(() => {
        localStorage.setItem('readingMode', readingMode);
    }, [readingMode]);

    // Обработка нажатий клавиш для навигации
    useEffect(() => {
        const handleKeyDown = (e: KeyboardEvent) => {
            if (readingMode === ReadingMode.PAGED) {
                if (e.key === 'ArrowRight' || e.key === 'd' || e.key === 'D') {
                    nextPage();
                } else if (e.key === 'ArrowLeft' || e.key === 'a' || e.key === 'A') {
                    prevPage();
                }
            }

            // Общие клавиши для обоих режимов
            if (e.key === 'f' || e.key === 'F') {
                toggleFullscreen();
            } else if (e.key === 'm' || e.key === 'M') {
                toggleReadingMode();
            }
        };

        window.addEventListener('keydown', handleKeyDown);
        return () => {
            window.removeEventListener('keydown', handleKeyDown);
        };
    }, [currentPage, readingMode, pages.length]);

    // Автоматически скрывать элементы управления после бездействия
    useEffect(() => {
        const handleMouseMove = () => {
            setShowControls(true);

            if (controlsTimeoutRef.current) {
                clearTimeout(controlsTimeoutRef.current);
            }

            controlsTimeoutRef.current = setTimeout(() => {
                setShowControls(false);
            }, 3000);
        };

        window.addEventListener('mousemove', handleMouseMove);

        // Начальный таймер
        controlsTimeoutRef.current = setTimeout(() => {
            setShowControls(false);
        }, 3000);

        return () => {
            window.removeEventListener('mousemove', handleMouseMove);
            if (controlsTimeoutRef.current) {
                clearTimeout(controlsTimeoutRef.current);
            }
        };
    }, []);

    // Переход к предыдущей странице или главе
    const prevPage = useCallback(() => {
        if (currentPage > 1) {
            setCurrentPage(currentPage - 1);
        } else {
            // Если мы на первой странице, переходим к предыдущей главе
            const currentChapterIndex = chapters.findIndex(ch => ch.id === chapterId);
            if (currentChapterIndex < chapters.length - 1) {
                const prevChapter = chapters[currentChapterIndex + 1];
                navigate(`/manga/${mangaId}/chapter/${prevChapter.id}`);
            }
        }
    }, [currentPage, chapterId, chapters, mangaId, navigate]);

    // Переход к следующей странице или главе
    const nextPage = useCallback(() => {
        if (currentPage < pages.length) {
            setCurrentPage(currentPage + 1);
        } else {
            // Если мы на последней странице, переходим к следующей главе
            const currentChapterIndex = chapters.findIndex(ch => ch.id === chapterId);
            if (currentChapterIndex > 0) {
                const nextChapter = chapters[currentChapterIndex - 1];
                navigate(`/manga/${mangaId}/chapter/${nextChapter.id}`);
            }
        }
    }, [currentPage, pages.length, chapterId, chapters, mangaId, navigate]);

    // Переключение полноэкранного режима
    const toggleFullscreen = useCallback(() => {
        if (!document.fullscreenElement) {
            readerRef.current?.requestFullscreen().catch(err => {
                console.error(`Ошибка при переходе в полноэкранный режим: ${err.message}`);
            });
            setIsFullscreen(true);
        } else {
            document.exitFullscreen();
            setIsFullscreen(false);
        }
    }, []);

    // Переключение режима чтения
    const toggleReadingMode = useCallback(() => {
        setReadingMode(prev =>
            prev === ReadingMode.PAGED ? ReadingMode.VERTICAL : ReadingMode.PAGED
        );
    }, []);

    // Переход к указанной главе
    const goToChapter = (chapterId: number) => {
        navigate(`/manga/${mangaId}/chapter/${chapterId}`);
        setShowChapters(false);
    };

    // Переход к списку манги
    const goToMangaDetails = () => {
        navigate(`/manga/${mangaId}`);
    };

    return (
        <div ref={readerRef} className="min-h-screen bg-gray-100 dark:bg-dark-900">
            {/* Верхняя панель навигации - отображается при наведении */}
            <div
                className={`fixed top-0 left-0 right-0 z-10 bg-dark-800 bg-opacity-70 backdrop-blur-sm transition-opacity duration-300 ${
                    showControls ? 'opacity-100' : 'opacity-0 pointer-events-none'
                }`}
            >
                <div className="container mx-auto p-4 flex justify-between items-center">
                    <button
                        onClick={goToMangaDetails}
                        className="text-white flex items-center space-x-2"
                    >
                        <HomeIcon className="h-5 w-5" />
                        <span>{chapter.title}</span>
                    </button>

                    <div className="flex items-center space-x-4">
                        <button
                            onClick={() => setShowChapters(!showChapters)}
                            className="text-white flex items-center space-x-1"
                        >
                            <Bars3Icon className="h-5 w-5" />
                            <span>Главы</span>
                        </button>

                        <button
                            onClick={toggleReadingMode}
                            className="text-white flex items-center space-x-1"
                        >
                            {readingMode === ReadingMode.PAGED ? (
                                <>
                                    <ChevronDownIcon className="h-5 w-5" />
                                    <span>Вертикально</span>
                                </>
                            ) : (
                                <>
                                    <ArrowLeftIcon className="h-5 w-5" />
                                    <ArrowRightIcon className="h-5 w-5" />
                                    <span>Постранично</span>
                                </>
                            )}
                        </button>

                        <button
                            onClick={toggleFullscreen}
                            className="text-white"
                        >
                            {isFullscreen ? (
                                <ArrowsPointingInIcon className="h-5 w-5" />
                            ) : (
                                <ArrowsPointingOutIcon className="h-5 w-5" />
                            )}
                        </button>
                    </div>
                </div>
            </div>

            {/* Выпадающий список глав */}
            {showChapters && (
                <div className="fixed top-16 right-4 z-20 w-64 max-h-96 overflow-y-auto bg-white dark:bg-dark-800 rounded-md shadow-lg">
                    <div className="p-2 border-b border-gray-200 dark:border-dark-700">
                        <h3 className="font-semibold text-lg">Главы</h3>
                    </div>
                    <div className="p-2">
                        {chapters.map((ch) => (
                            <button
                                key={ch.id}
                                className={`w-full text-left p-2 rounded-md hover:bg-gray-100 dark:hover:bg-dark-700 ${
                                    ch.id === chapterId ? 'bg-primary-100 dark:bg-primary-900' : ''
                                }`}
                                onClick={() => goToChapter(ch.id)}
                            >
                                Глава {ch.number}: {ch.title}
                            </button>
                        ))}
                    </div>
                </div>
            )}

            {/* Основная область чтения */}
            {readingMode === ReadingMode.VERTICAL ? (
                // Вертикальный режим (прокрутка)
                <div className="reader-container pt-16">
                    {pages.map((page) => (
                        <div key={page.id} className="mb-4">
                            <img
                                src={page.imageUrl}
                                alt={`Страница ${page.number}`}
                                className="reader-page"
                                loading="lazy"
                            />
                        </div>
                    ))}
                </div>
            ) : (
                // Постраничный режим
                <div className="pt-16 h-screen flex items-center justify-center">
                    <div className="relative max-w-4xl mx-auto">
                        <img
                            src={pages[currentPage - 1]?.imageUrl}
                            alt={`Страница ${currentPage}`}
                            className="reader-page max-h-[calc(100vh-100px)]"
                        />

                        {/* Навигация по сторонам экрана */}
                        <button
                            className="absolute top-0 left-0 w-1/2 h-full opacity-0"
                            onClick={prevPage}
                            disabled={currentPage === 1 && chapters.findIndex(ch => ch.id === chapterId) === chapters.length - 1}
                        />
                        <button
                            className="absolute top-0 right-0 w-1/2 h-full opacity-0"
                            onClick={nextPage}
                            disabled={currentPage === pages.length && chapters.findIndex(ch => ch.id === chapterId) === 0}
                        />

                        {/* Индикатор страниц */}
                        <div className={`fixed bottom-4 left-1/2 transform -translate-x-1/2 bg-dark-800 bg-opacity-70 text-white px-4 py-2 rounded-full transition-opacity duration-300 ${
                            showControls ? 'opacity-100' : 'opacity-0'
                        }`}>
                            {currentPage} / {pages.length}
                        </div>
                    </div>
                </div>
            )}

            {/* Панель управления (для постраничного режима) */}
            {readingMode === ReadingMode.PAGED && (
                <div className={`reader-controls flex space-x-4 transition-opacity duration-300 ${
                    showControls ? 'opacity-100' : 'opacity-0 pointer-events-none'
                }`}>
                    <button
                        onClick={prevPage}
                        className="p-2 rounded-full bg-white text-dark-800 dark:bg-dark-700 dark:text-white disabled:opacity-50"
                        disabled={currentPage === 1 && chapters.findIndex(ch => ch.id === chapterId) === chapters.length - 1}
                    >
                        <ArrowLeftIcon className="h-6 w-6" />
                    </button>
                    <button
                        onClick={nextPage}
                        className="p-2 rounded-full bg-white text-dark-800 dark:bg-dark-700 dark:text-white disabled:opacity-50"
                        disabled={currentPage === pages.length && chapters.findIndex(ch => ch.id === chapterId) === 0}
                    >
                        <ArrowRightIcon className="h-6 w-6" />
                    </button>
                </div>
            )}
        </div>
    );
};

export default MangaReader;
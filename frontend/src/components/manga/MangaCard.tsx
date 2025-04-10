import React from 'react';
import { Link } from 'react-router-dom';
import { Manga } from '../../types/manga.types';
import { StarIcon } from '@heroicons/react/24/solid';

interface MangaCardProps {
    manga: Manga;
}

const MangaCard: React.FC<MangaCardProps> = ({ manga }) => {
    return (
        <Link to={`/manga/${manga.id}`} className="manga-card block">
            <div className="relative h-full">
                {/* Обложка манги */}
                <div className="overflow-hidden">
                    <img
                        src={manga.coverUrl}
                        alt={manga.title}
                        className="manga-cover transition-transform duration-300 hover:scale-105"
                        loading="lazy"
                    />
                </div>

                {/* Статус манги */}
                <div className="absolute top-2 right-2">
          <span className={`inline-flex items-center rounded-full px-2 py-1 text-xs font-semibold ${getStatusColor(manga.status)}`}>
            {getStatusText(manga.status)}
          </span>
                </div>

                {/* Информация о манге */}
                <div className="p-3">
                    <h3 className="mb-1 text-base font-semibold line-clamp-2" title={manga.title}>
                        {manga.title}
                    </h3>

                    {/* Рейтинг */}
                    <div className="flex items-center mb-2">
                        <StarIcon className="h-4 w-4 text-yellow-400" />
                        <span className="ml-1 text-sm">{manga.rating.toFixed(1)}</span>
                    </div>

                    {/* Жанры */}
                    <div className="flex flex-wrap gap-1 mt-1">
                        {manga.genres.slice(0, 2).map((genre) => (
                            <span
                                key={genre.id}
                                className="inline-block rounded-full bg-primary-100 px-2 py-0.5 text-xs text-primary-800 dark:bg-primary-900 dark:text-primary-200"
                            >
                {genre.name}
              </span>
                        ))}
                        {manga.genres.length > 2 && (
                            <span className="inline-block rounded-full bg-gray-100 px-2 py-0.5 text-xs text-gray-600 dark:bg-gray-800 dark:text-gray-300">
                +{manga.genres.length - 2}
              </span>
                        )}
                    </div>
                </div>
            </div>
        </Link>
    );
};

// Вспомогательные функции для отображения статуса
const getStatusColor = (status: string): string => {
    switch (status) {
        case 'ongoing':
            return 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200';
        case 'completed':
            return 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200';
        case 'hiatus':
            return 'bg-amber-100 text-amber-800 dark:bg-amber-900 dark:text-amber-200';
        default:
            return 'bg-gray-100 text-gray-800 dark:bg-gray-800 dark:text-gray-200';
    }
};

const getStatusText = (status: string): string => {
    switch (status) {
        case 'ongoing':
            return 'Выходит';
        case 'completed':
            return 'Завершён';
        case 'hiatus':
            return 'Пауза';
        default:
            return status;
    }
};

export default MangaCard;
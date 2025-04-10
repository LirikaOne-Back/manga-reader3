import api from './api';
import {
    Manga,
    MangaFilter,
    MangaListResponse,
    Chapter,
    Page,
    Genre,
} from '../types/manga.types';

const MangaService = {
    // Получение списка манги с фильтрацией
    async getAll(filter: MangaFilter): Promise<MangaListResponse> {
        const params = new URLSearchParams();

        if (filter.genre) params.append('genre', filter.genre);
        if (filter.status) params.append('status', filter.status);
        if (filter.sortBy) params.append('sortBy', filter.sortBy);
        if (filter.sortDesc !== undefined) params.append('sortDesc', String(filter.sortDesc));
        if (filter.search) params.append('search', filter.search);

        params.append('page', String(filter.page));
        params.append('pageSize', String(filter.pageSize));

        const response = await api.get<MangaListResponse>(`/manga?${params.toString()}`);
        return response.data;
    },

    // Получение манги по ID
    async getById(id: number): Promise<Manga> {
        const response = await api.get<Manga>(`/manga/${id}`);
        return response.data;
    },

    // Получение списка глав манги
    async getChapters(mangaId: number): Promise<Chapter[]> {
        const response = await api.get<Chapter[]>(`/manga/${mangaId}/chapters`);
        return response.data;
    },

    // Получение страниц главы
    async getChapterPages(chapterId: number): Promise<Page[]> {
        const response = await api.get<Page[]>(`/chapters/${chapterId}/pages`);
        return response.data;
    },

    // Получение информации о главе
    async getChapter(chapterId: number): Promise<Chapter> {
        const response = await api.get<Chapter>(`/chapters/${chapterId}`);
        return response.data;
    },

    // Получение жанров
    async getGenres(): Promise<Genre[]> {
        const response = await api.get<Genre[]>('/manga/genres');
        return response.data;
    },

    // Добавление манги в закладки
    async addBookmark(mangaId: number): Promise<void> {
        await api.post('/bookmarks', { mangaId });
    },

    // Удаление манги из закладок
    async removeBookmark(mangaId: number): Promise<void> {
        await api.delete(`/bookmarks/${mangaId}`);
    },

    // Получение закладок пользователя
    async getBookmarks(): Promise<Manga[]> {
        const response = await api.get<Manga[]>('/bookmarks');
        return response.data;
    },

    // Сохранение прогресса чтения
    async saveReadProgress(chapterId: number, page: number): Promise<void> {
        await api.post('/history', { chapterId, page });
    },

    // Получение истории чтения
    async getReadHistory(): Promise<any[]> {
        const response = await api.get('/history');
        return response.data;
    },
};

// Экспорт по умолчанию для использования через import MangaService from '...'
export default MangaService;
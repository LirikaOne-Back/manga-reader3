export interface Genre {
    id: number;
    name: string;
}

export interface Manga {
    id: number;
    title: string;
    alterTitle?: string;
    description: string;
    coverUrl: string;
    year: number;
    status: 'ongoing' | 'completed' | 'hiatus';
    author: string;
    artist?: string;
    rating: number;
    genres: Genre[];
    createdAt: string;
    updatedAt: string;
}

export interface Chapter {
    id: number;
    mangaId: number;
    number: number;
    title: string;
    pageCount: number;
    createdAt: string;
    updatedAt: string;
}

export interface Page {
    id: number;
    chapterId: number;
    number: number;
    imageUrl: string;
}

export interface MangaFilter {
    genre?: string;
    status?: string;
    sortBy?: string;
    sortDesc?: boolean;
    search?: string;
    page: number;
    pageSize: number;
}

export interface MangaListResponse {
    data: Manga[];
    total: number;
    page: number;
    size: number;
}

export interface Bookmark {
    id: number;
    userId: number;
    mangaId: number;
    createdAt: string;
}

export interface ReadHistory {
    id: number;
    userId: number;
    mangaId: number;
    chapterId: number;
    page: number;
    readAt: string;
}
import { configureStore } from '@reduxjs/toolkit';
// Здесь импортируем слайсы (reducers), когда они будут созданы
// import { mangaReducer } from './manga.slice';
// import { authReducer } from './auth.slice';

export const store = configureStore({
    reducer: {
        // manga: mangaReducer,
        // auth: authReducer,
        // Пока что используем пустой reducer
        dummy: (state = {}) => state
    },
    middleware: (getDefaultMiddleware) =>
        getDefaultMiddleware({
            serializableCheck: false,
        }),
});

// Типы для хранилища
export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;
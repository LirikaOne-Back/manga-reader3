import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { BookOpenIcon } from '@heroicons/react/24/outline';
import { toast } from 'react-hot-toast';

const Register: React.FC = () => {
    const [username, setUsername] = useState('');
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [confirmPassword, setConfirmPassword] = useState('');
    const [loading, setLoading] = useState(false);

    const navigate = useNavigate();

    // Обработчик формы регистрации
    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();

        // Валидация формы
        if (!username.trim() || !email.trim() || !password.trim() || !confirmPassword.trim()) {
            toast.error('Пожалуйста, заполните все поля');
            return;
        }

        if (password !== confirmPassword) {
            toast.error('Пароли не совпадают');
            return;
        }

        if (password.length < 6) {
            toast.error('Пароль должен содержать не менее 6 символов');
            return;
        }

        // Проверка формата email
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        if (!emailRegex.test(email)) {
            toast.error('Введите корректный email адрес');
            return;
        }

        setLoading(true);

        try {
            // В реальном приложении здесь будет запрос на сервер для регистрации
            // const response = await authService.register(username, email, password);

            // Имитация успешной регистрации
            setTimeout(() => {
                toast.success('Регистрация выполнена успешно!');
                setLoading(false);
                navigate('/login');
            }, 1000);
        } catch (error) {
            setLoading(false);
            toast.error('Ошибка регистрации. Пожалуйста, попробуйте снова.');
            console.error('Registration error:', error);
        }
    };

    return (
        <div className="min-h-screen bg-gray-100 dark:bg-dark-900 flex flex-col justify-center py-12 sm:px-6 lg:px-8">
            <div className="sm:mx-auto sm:w-full sm:max-w-md">
                <div className="flex justify-center">
                    <Link to="/" className="flex items-center">
                        <BookOpenIcon className="h-12 w-12 text-primary-600" />
                    </Link>
                </div>
                <h2 className="mt-6 text-center text-3xl font-extrabold text-gray-900 dark:text-white">
                    Создайте аккаунт
                </h2>
                <p className="mt-2 text-center text-sm text-gray-600 dark:text-gray-400">
                    Или{' '}
                    <Link to="/login" className="font-medium text-primary-600 hover:text-primary-500 dark:text-primary-400">
                        войдите, если у вас уже есть аккаунт
                    </Link>
                </p>
            </div>

            <div className="mt-8 sm:mx-auto sm:w-full sm:max-w-md">
                <div className="bg-white dark:bg-dark-800 py-8 px-4 shadow sm:rounded-lg sm:px-10">
                    <form className="space-y-6" onSubmit={handleSubmit}>
                        <div>
                            <label htmlFor="username" className="block text-sm font-medium text-gray-700 dark:text-gray-300">
                                Имя пользователя
                            </label>
                            <div className="mt-1">
                                <input
                                    id="username"
                                    name="username"
                                    type="text"
                                    autoComplete="username"
                                    required
                                    className="input"
                                    value={username}
                                    onChange={(e) => setUsername(e.target.value)}
                                />
                            </div>
                        </div>

                        <div>
                            <label htmlFor="email" className="block text-sm font-medium text-gray-700 dark:text-gray-300">
                                Email
                            </label>
                            <div className="mt-1">
                                <input
                                    id="email"
                                    name="email"
                                    type="email"
                                    autoComplete="email"
                                    required
                                    className="input"
                                    value={email}
                                    onChange={(e) => setEmail(e.target.value)}
                                />
                            </div>
                        </div>

                        <div>
                            <label htmlFor="password" className="block text-sm font-medium text-gray-700 dark:text-gray-300">
                                Пароль
                            </label>
                            <div className="mt-1">
                                <input
                                    id="password"
                                    name="password"
                                    type="password"
                                    autoComplete="new-password"
                                    required
                                    className="input"
                                    value={password}
                                    onChange={(e) => setPassword(e.target.value)}
                                />
                            </div>
                        </div>

                        <div>
                            <label htmlFor="confirm-password" className="block text-sm font-medium text-gray-700 dark:text-gray-300">
                                Подтверждение пароля
                            </label>
                            <div className="mt-1">
                                <input
                                    id="confirm-password"
                                    name="confirm-password"
                                    type="password"
                                    autoComplete="new-password"
                                    required
                                    className="input"
                                    value={confirmPassword}
                                    onChange={(e) => setConfirmPassword(e.target.value)}
                                />
                            </div>
                        </div>

                        <div className="flex items-center">
                            <input
                                id="terms"
                                name="terms"
                                type="checkbox"
                                required
                                className="h-4 w-4 text-primary-600 focus:ring-primary-500 border-gray-300 rounded"
                            />
                            <label htmlFor="terms" className="ml-2 block text-sm text-gray-900 dark:text-gray-300">
                                Я согласен с{' '}
                                <a href="#" className="font-medium text-primary-600 hover:text-primary-500 dark:text-primary-400">
                                    условиями использования
                                </a>
                            </label>
                        </div>

                        <div>
                            <button
                                type="submit"
                                className="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 disabled:opacity-50 disabled:cursor-not-allowed"
                                disabled={loading}
                            >
                                {loading ? 'Регистрация...' : 'Зарегистрироваться'}
                            </button>
                        </div>
                    </form>

                    <div className="mt-6">
                        <div className="relative">
                            <div className="absolute inset-0 flex items-center">
                                <div className="w-full border-t border-gray-300 dark:border-dark-600"></div>
                            </div>
                            <div className="relative flex justify-center text-sm">
                <span className="px-2 bg-white dark:bg-dark-800 text-gray-500 dark:text-gray-400">
                  Или зарегистрируйтесь через
                </span>
                            </div>
                        </div>

                        <div className="mt-6 grid grid-cols-2 gap-3">
                            <div>
                                <a
                                    href="#"
                                    className="w-full inline-flex justify-center py-2 px-4 border border-gray-300 dark:border-dark-600 rounded-md shadow-sm bg-white dark:bg-dark-700 text-sm font-medium text-gray-700 dark:text-gray-200 hover:bg-gray-50 dark:hover:bg-dark-600"
                                >
                                    <span className="sr-only">Войти через Google</span>
                                    <svg className="h-5 w-5" fill="currentColor" viewBox="0 0 24 24" aria-hidden="true">
                                        <path
                                            d="M12.48 10.92v3.28h7.84c-.24 1.84-.853 3.187-1.787 4.133-1.147 1.147-2.933 2.4-6.053 2.4-4.827 0-8.6-3.893-8.6-8.72s3.773-8.72 8.6-8.72c2.6 0 4.507 1.027 5.907 2.347l2.307-2.307C18.747 1.44 16.133 0 12.48 0 5.867 0 .307 5.387.307 12s5.56 12 12.173 12c3.573 0 6.267-1.173 8.373-3.36 2.16-2.16 2.84-5.213 2.84-7.667 0-.76-.053-1.467-.173-2.053H12.48z"
                                        />
                                    </svg>
                                </a>
                            </div>

                            <div>
                                <a
                                    href="#"
                                    className="w-full inline-flex justify-center py-2 px-4 border border-gray-300 dark:border-dark-600 rounded-md shadow-sm bg-white dark:bg-dark-700 text-sm font-medium text-gray-700 dark:text-gray-200 hover:bg-gray-50 dark:hover:bg-dark-600"
                                >
                                    <span className="sr-only">Войти через VK</span>
                                    <svg className="h-5 w-5" fill="currentColor" viewBox="0 0 24 24" aria-hidden="true">
                                        <path
                                            d="M12.785 16.241s.288-.032.436-.194c.136-.148.132-.427.132-.427s-.02-1.304.587-1.496c.596-.19 1.362 1.259 2.173 1.818.613.422 1.077.33 1.077.33l2.163-.03s1.132-.07.594-.957c-.044-.073-.312-.658-1.605-1.861-1.353-1.258-1.172-1.055.458-3.233.994-1.328 1.39-2.138 1.267-2.486-.118-.33-.846-.244-.846-.244l-2.433.015s-.18-.025-.314.055c-.13.079-.214.262-.214.262s-.382 1.018-.893 1.882c-1.076 1.822-1.506 1.918-1.683 1.805-.41-.266-.308-1.07-.308-1.642 0-1.784.27-2.527-.525-2.72-.264-.065-.456-.108-1.128-.116-.865-.01-1.597.003-2.01.207-.277.134-.49.435-.36.452.16.022.524.098.716.362.249.342.24 1.107.24 1.107s.142 2.112-.332 2.374c-.327.18-.777-.187-1.74-1.838-.494-.853-.867-1.795-.867-1.795s-.071-.176-.2-.272c-.154-.115-.372-.152-.372-.152l-2.314.015s-.344.01-.47.16c-.113.135-.01.414-.01.414s1.79 4.185 3.812 6.293c1.857 1.937 3.965 1.81 3.965 1.81h.957z"
                                        />
                                    </svg>
                                </a>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default Register;
import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';

import './Registration.css'
import Header from '../components/layout/Header/Header'

function Registration({ onLogin }) {
    return(
        <div className="all_page">
            <Header />
            <RegistrationAndLogin onLogin={onLogin} />
        </div>
    );
}

function RegistrationAndLogin({ onLogin }) {
    const [isLogIn, setIsLogIn] = useState(false);
    const [isActive, setIsActive] = useState(false);
    const [formData, setFormData] = useState({
        email: '',
        phone: '+7',
        username: '',
        password: ''
    });
    const [loginData, setLoginData] = useState({
        username: '',
        password: ''
    });
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');
    const navigate = useNavigate();

    const handleToggleForm = () => {
        setIsLogIn(!isLogIn);
        setError('');
        setLoginData({ username: '', password: '' });
        setFormData({ email: '', phone: '+7', username: '', password: '' });
    };

    const handleLoginChange = (e) => {
        const { name, value } = e.target;
        setLoginData(prev => ({ ...prev, [name]: value }));
    };

    const handleRegisterChange = (e) => {
        const { name, value } = e.target;
        setFormData(prev => ({ ...prev, [name]: value }));
    };

    const handlePhoneChange = (e) => {
        let value = e.target.value;
        if (!value.startsWith('+7')) value = '+7';
        const digits = value.replace(/[^\d+]/g, '');
        if (digits.length <= 12) {
            setFormData(prev => ({ ...prev, phone: digits }));
        }
    };

    const handlePhoneFocus = (e) => {
        if (e.target.value === '') {
            setFormData(prev => ({ ...prev, phone: '+7' }));
        }
    };

    const handleRegister = async () => {
        if (!formData.email || !formData.phone || !formData.username || !formData.password) {
            setError('Пожалуйста, заполните все поля');
            return;
        }

        if (formData.phone.length < 12) {
            setError('Введите корректный номер телефона');
            return;
        }

        setLoading(true);
        setError('');

        try {
            const response = await axios.post('/registration', {
                email: formData.email,
                phone: formData.phone,
                username: formData.username,
                password: formData.password
            });

            if (response.data.success) {
                if (response.data.token) {
                    localStorage.setItem('authToken', response.data.token);
                }
                onLogin(response.data.token);
                setIsActive(true);
                setTimeout(() => navigate('/'), 1500);
            } else {
                setError(response.data.message || 'Ошибка при регистрации');
            }
        } catch (error) {
            console.error('Ошибка регистрации:', error);
            if (error.response?.data?.message) {
                setError(error.response.data.message);
            } else if (error.response?.data?.error) {
                setError(error.response.data.error);
            } else {
                setError('Ошибка при регистрации. Попробуйте снова.');
            }
        } finally {
            setLoading(false);
        }
    };

    const handleLogin = async () => {
        if (!loginData.username || !loginData.password) {
            setError('Пожалуйста, заполните все поля');
            return;
        }

        setLoading(true);
        setError('');

        try {
            const response = await axios.post('/login', {
                username: loginData.username,
                password: loginData.password
            });

            if (response.data.success) {
                if (response.data.token) {
                    localStorage.setItem('authToken', response.data.token);
                }
                onLogin(response.data.token);
                setIsActive(true);
                setTimeout(() => navigate('/'), 1500);
            } else {
                setError(response.data.message || 'Ошибка при входе');
            }
        } catch (error) {
            console.error('Ошибка входа:', error);
            if (error.response?.data?.message) {
                setError(error.response.data.message);
            } else if (error.response?.data?.error) {
                setError(error.response.data.error);
            } else {
                setError('Ошибка при входе. Проверьте логин и пароль.');
            }
        } finally {
            setLoading(false);
        }
    };

    return(
        <div className="container_login">  
            <div className="table">
                <div className="table-cell">
                    <div className={`container ${isLogIn ? 'log-in' : ''} ${isActive ? 'active' : ''}`}>
                        <div className="box"></div>
                        <div className="container-forms">
                            <div className="container-info">
                                <div className="info-item">
                                    <div className="table">
                                        <div className="table-cell">
                                            <p>Уже есть аккаунт?</p>
                                            <div className="btn" onClick={handleToggleForm}>Войти</div>
                                        </div>
                                    </div>
                                </div>
                                <div className="info-item">
                                    <div className="table">
                                        <div className="table-cell">
                                            <p>Нет аккаунта?</p>
                                            <div className="btn" onClick={handleToggleForm}>Регистрация</div>
                                        </div>
                                    </div>
                                </div>
                            </div>
                            <div className="container-form">
                                <div className="form-item log-in">
                                    <div className="table">
                                        <div className="table-cell">
                                            <input 
                                                type="text" 
                                                name="username" 
                                                placeholder="Логин" 
                                                value={loginData.username}
                                                onChange={handleLoginChange}
                                                disabled={loading}
                                                autoComplete="username"
                                            />
                                            <input 
                                                type="password" 
                                                name="password" 
                                                placeholder="Пароль" 
                                                value={loginData.password}
                                                onChange={handleLoginChange}
                                                disabled={loading}
                                                autoComplete="current-password"
                                            />
                                            {error && <div className="error-message">{error}</div>}
                                            <div 
                                                className={`btn ${loading ? 'loading' : ''}`} 
                                                onClick={handleLogin}
                                                disabled={loading}
                                            >
                                                {loading ? 'Загрузка...' : 'Войти'}
                                            </div>
                                        </div>
                                    </div>
                                </div>
                                
                                <div className="form-item sign-up">
                                    <div className="table">
                                        <div className="table-cell">
                                            <input 
                                                type="email" 
                                                name="email" 
                                                placeholder="Email" 
                                                value={formData.email}
                                                onChange={handleRegisterChange}
                                                disabled={loading}
                                                autoComplete="email"
                                            />
                                            <input 
                                                type="tel" 
                                                name="phone" 
                                                placeholder="+7 (XXX) XXX-XX-XX"
                                                value={formData.phone}
                                                onChange={handlePhoneChange}
                                                onFocus={handlePhoneFocus}
                                                disabled={loading}
                                                autoComplete="tel"
                                            />
                                            <input 
                                                type="text" 
                                                name="username" 
                                                placeholder="Логин" 
                                                value={formData.username}
                                                onChange={handleRegisterChange}
                                                disabled={loading}
                                                autoComplete="username"
                                            />
                                            <input 
                                                type="password" 
                                                name="password" 
                                                placeholder="Пароль" 
                                                value={formData.password}
                                                onChange={handleRegisterChange}
                                                disabled={loading}
                                                autoComplete="new-password"
                                            />
                                            {error && <div className="error-message">{error}</div>}
                                            <div 
                                                className={`btn ${loading ? 'loading' : ''}`} 
                                                onClick={handleRegister}
                                                disabled={loading}
                                            >
                                                {loading ? 'Загрузка...' : 'Регистрация'}
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}

export default Registration
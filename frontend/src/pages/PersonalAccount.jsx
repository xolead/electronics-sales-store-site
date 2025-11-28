import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import Header from "../components/layout/Header/Header";
import HeaderForLogined from "../components/layout/Header/HeaderForLogined";
import './PersonalAccount.css';
import axios from 'axios';

const PersonalAccount = () => {
    const [isAdmin, setIsAdmin] = useState(false);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        // Здесь будет запрос к серверу для проверки прав администратора
        checkAdminStatus();
    }, []);

    const checkAdminStatus = async () => {
        try {
            // Замените этот пример на реальный запрос к вашему API
            const response = await axios.get('/api/admin/check-status');
            setIsAdmin(response.data.isAdmin);
            
        } catch (error) {
            console.error('Ошибка при проверке прав администратора:', error);
            setIsAdmin(false);
        } finally {
            setLoading(false);
        }
    };

    // Если данные еще загружаются, можно показать заглушку
    if (loading) {
        return (
            <div className="personal-account-page">
                <HeaderForLogined />
                <div className="personal-account-container">
                    <div className="loading-container">
                        <p>Загрузка...</p>
                    </div>
                </div>
            </div>
        );
    }

    return (
        <div className="personal-account-page">
            <HeaderForLogined />
            
            <div className="personal-account-container">
                <div className="personal-account-header">
                    <h1>Личный кабинет</h1>
                    {isAdmin && (
                        <Link to="/admin" className="home-link admin-panel-link" style={{fontSize:"20px"}}>
                            Админ
                        </Link>
                    )}
                </div>
                
                <div className="promo-banner">
                    <div className="promo-content">
                        <h2 style={{fontSize: '30px'}}>
                            🎉 Получите по <span style={{color: '#ee73a3', fontSize: '55px'}}>е</span>-баллу за каждые 100 рублей в чеке!
                        </h2>
                    </div>
                </div>
                
                <div className="account-sections">
                    <div className="account-section">
                        <h3>Мои данные</h3>
                        <div className="section-content">
                            <p>Имя: Иван Иванов</p>
                            <p>Email: example@mail.com</p>
                            <p>Телефон: +7 (999) 999-99-99</p>
                        </div>
                    </div>
                    
                    <div className="account-section">
                        <h3>Мои заказы</h3>
                        <div className="section-content">
                            <p>Последние заказы будут отображаться здесь</p>
                        </div>
                    </div>
                    
                    <div className="account-section">
                        <h3>Е-баллы</h3>
                        <div className="section-content">
                            <div className="points-balance">
                                <span className="points-count">150 баллов</span>
                                <p>Доступно для использования</p>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    )
}

export default PersonalAccount;
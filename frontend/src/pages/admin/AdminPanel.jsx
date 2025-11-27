import React, { useState } from 'react';
import { Link, useLocation } from 'react-router-dom';
import Header from '../../components/layout/Header/Header';
import './AdminPanel.css';
import AdminProducts from './AdminProducts';

function AdminPanel() {
  const location = useLocation();
  
  const navItems = [
    { path: '/admin/dashboard', label: 'Дашборд', icon: '📊' },
    { path: '/admin/products', label: 'Товары', icon: '🛍️' },
    { path: '/admin/orders', label: 'Заказы', icon: '📦', badge: 5 },
    { path: '/admin/users', label: 'Пользователи', icon: '👥' },
    { 
      path: '/admin/settings', 
      label: 'Настройки', 
      icon: '⚙️',
      dropdown: [
        { path: '/admin/settings/general', label: 'Основные' },
        { path: '/admin/settings/payments', label: 'Оплата' },
        { path: '/admin/settings/shipping', label: 'Доставка' }
      ]
    },
  ];

  return (
    <div className="admin-panel-page">
      <Header />
      
      {/* Admin Navigation Bar */}
      <nav className="admin-navbar">
        <div className="admin-nav-container">
          <ul className="admin-nav-menu">
            {navItems.map((item) => (
              <li 
                key={item.path} 
                className={`admin-nav-item ${item.dropdown ? 'admin-nav-dropdown' : ''}`}
              >
                <Link 
                  to={item.path}
                  className={`admin-nav-link ${location.pathname === item.path ? 'active' : ''}`}
                >
                  <span className="admin-nav-icon">{item.icon}</span>
                  {item.label}
                  {item.badge && <span className="admin-nav-badge">{item.badge}</span>}
                </Link>
                
                {item.dropdown && (
                  <div className="dropdown-menu">
                    {item.dropdown.map(dropdownItem => (
                      <Link 
                        key={dropdownItem.path}
                        to={dropdownItem.path}
                        className="dropdown-item"
                      >
                        {dropdownItem.label}
                      </Link>
                    ))}
                  </div>
                )}
              </li>
            ))}
          </ul>
        </div>
      </nav>

      {/* Admin Content */}
      <div className="admin-content">
        <div className="admin-section">
          <h2>Панель управления</h2>
          <p>Добро пожаловать в панель администратора!</p>
          {/* Здесь будет контент конкретной страницы */}
        </div>
      </div>
    </div>
  );
}

export default AdminPanel;
import React, { useState, useEffect } from 'react';
import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom';

import './App.css';
import Header from './components/layout/Header/Header'
import HeaderForLogined from './components/layout/Header/HeaderForLogined'
import Create from './pages/admin/Create';
import Cart from './pages/Cart';
import Registration from './pages/Registration';
import ProductDetail from './pages/ProductDetail';
import AdminProducts from './pages/admin/AdminProducts'; 
import AdminPanel from './pages/admin/AdminPanel';
import { getCategoryFromParameters } from './utils/parameters';
import { 
  loadProducts, 
  getFullImageUrl, 
  deleteProduct 
} from './utils/loadProductsAndDelete';
import PersonalAccount from './pages/PersonalAccount';
import axios from 'axios';

const App = () => {
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [loading, setLoading] = useState(true);

  // Настройка интерцептора для axios
  useEffect(() => {
    const requestInterceptor = axios.interceptors.request.use(
      (config) => {
        const token = localStorage.getItem('authToken');
        if (token) {
          config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
      },
      (error) => {
        return Promise.reject(error);
      }
    );

    // Очистка интерцептора при размонтировании
    return () => {
      axios.interceptors.request.eject(requestInterceptor);
    };
  }, []);

  // Проверяем авторизацию при загрузке приложения
  useEffect(() => {
    checkAuthStatus();
  }, []);

  const checkAuthStatus = async () => {
    try {
      const token = localStorage.getItem('authToken');
      
      if (token) {
        // Если есть токен, проверяем его валидность на сервере
        // Замените на ваш реальный эндпоинт проверки токена
        const response = await axios.get('/api/auth/verify', {
          headers: {
            'Authorization': `Bearer ${token}`
          }
        });
        
        if (response.status === 200) {
          setIsLoggedIn(true);
        } else {
          // Если токен невалидный, удаляем его
          localStorage.removeItem('authToken');
          setIsLoggedIn(false);
        }
      } else {
        setIsLoggedIn(false);
      }
    } catch (error) {
      console.error('Ошибка при проверке авторизации:', error);
      // Если ошибка 401 (Unauthorized) или другая, считаем пользователя неавторизованным
      if (error.response && error.response.status === 401) {
        localStorage.removeItem('authToken');
      }
      setIsLoggedIn(false);
    } finally {
      setLoading(false);
    }
  };

  // Функции для управления авторизацией
  const handleLogin = (token) => {
    localStorage.setItem('authToken', token);
    setIsLoggedIn(true);
  };

  const handleLogout = async () => {
    try {
      // Опционально: отправляем запрос на logout на сервер
      await axios.post('/api/auth/logout');
    } catch (error) {
      console.error('Ошибка при выходе:', error);
    } finally {
      localStorage.removeItem('authToken');
      setIsLoggedIn(false);
    }
  };

  // Если проверка авторизации еще идет, показываем заглушку
  if (loading) {
    return (
      <div className="app-loading">
        <div className="loading-spinner"></div>
        <p>Загрузка...</p>
      </div>
    );
  }

  return (
    <Router>
      <Routes>
        <Route path="/" element={<HomePage isLoggedIn={isLoggedIn} onLogout={handleLogout} />} />
        <Route path="/create" element={<Create isLoggedIn={isLoggedIn} onLogout={handleLogout} />} />
        <Route path="/Cart" element={<Cart isLoggedIn={isLoggedIn} onLogout={handleLogout} />} />
        <Route path="/registration" element={<Registration onLogin={handleLogin} />} />
        <Route path="/personal" element={<PersonalAccount isLoggedIn={isLoggedIn} onLogout={handleLogout} />} />
        <Route path="/product/:id" element={<ProductDetail isLoggedIn={isLoggedIn} onLogout={handleLogout} />} />
        <Route path="/admin" element={<AdminPanel isLoggedIn={isLoggedIn} onLogout={handleLogout} />} />
        <Route path="/admin/products" element={<AdminProducts isLoggedIn={isLoggedIn} onLogout={handleLogout} />} />
        <Route path="/admin/create" element={<Create isLoggedIn={isLoggedIn} onLogout={handleLogout} />} />
      </Routes>
    </Router>
  );
}
function HomePage({ isLoggedIn, onLogout }) {
  return (
    <div className="App">
      {isLoggedIn ? (
        <HeaderForLogined onLogout={onLogout} />
      ) : (
        <Header />
      )}
      <ShoppingList_box>
        <ShoppingList isLoggedIn={isLoggedIn} />
      </ShoppingList_box>
    </div>
  );
}

function ShoppingList({ isLoggedIn }) {
  const [selectedProduct, setSelectedProduct] = useState(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [products, setProducts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  // Загрузка товаров при монтировании компонента
  useEffect(() => {
    loadProductsData();
  }, []);
 
  const loadProductsData = async () => {
    try {
      setLoading(true);
      setError(null);
      
      const result = await loadProducts();
      setProducts(result.data);
      
      if (!result.success && result.error) {
        setError(result.error);
      }
      
    } catch (err) {
      setError('Неожиданная ошибка при загрузке товаров');
      console.error('Ошибка загрузки:', err);
    } finally {
      setLoading(false);
    }
  };

  const handleBuyClick = (product) => {
    
    setSelectedProduct(product);
    setIsModalOpen(true);
  };
  
  const handleConfirm = () => {
    const cartItem = {
      ...selectedProduct,
      quantity: 1
    };
    
    const existingCart = JSON.parse(localStorage.getItem('electronic_cart') || '[]');
    const existingItemIndex = existingCart.findIndex(item => item.id === selectedProduct.id);
    
    if (existingItemIndex >= 0) {
      existingCart[existingItemIndex].quantity += 1;
    } else {
      existingCart.push(cartItem);
    }
    
    localStorage.setItem('electronic_cart', JSON.stringify(existingCart));
    window.dispatchEvent(new Event('cartUpdated'));
    
    setIsModalOpen(false);
    setSelectedProduct(null);
  };

  const handleCancel = () => {
    setIsModalOpen(false);
    setSelectedProduct(null);
  };

  const refreshProducts = () => {
    loadProductsData();
  };

  if (loading) {
    return (
      <div className="loading-container">
        <div className="loading-spinner"></div>
        <p>Загрузка товаров...</p>
      </div>
    );
  }

  if (error && products.length === 0) {
    return (
      <div className="error-container">
        <p>{error}</p>
        <button onClick={loadProductsData} className="retry-btn">
          Попробовать снова
        </button>
      </div>
    );
  }

  return (
    <>
      <div className="shop_products">
        <div className="products-header">
          <h2>Список товаров ({products.length})</h2>
          <button onClick={refreshProducts} className="refresh-btn">
            Обновить
          </button>
        </div>

        <div className="products-grid">
            {products.map(product => (
              <div key={product.id} className="product-card">
                <Link to={`/product/${product.id}`} className="product-card-link">
                  <img 
                    src={getFullImageUrl(product.images[0])} 
                    alt={product.name} 
                    className="product-image" 
                  />
                  <div className="product-details">
                    <span className="category">{getCategoryFromParameters(product.parameters)}</span>
                    <h3>{product.name}</h3>
                  </div>
                </Link>
            
                <div className="price-section">
                  <span className="price">{product.price.toLocaleString()} ₽</span>
                  <button 
                    className="buy-btn" 
                    onClick={(e) => {
                      e.stopPropagation();
                      handleBuyClick(product);
                    }}
                  >
                    Купить
                  </button>
                </div>
              </div>
            ))}
        </div>
          
        {products.length === 0 && !loading && (
          <div className="empty-state">
            <p>Товары не найдены</p>
            {isLoggedIn && (
              <Link to="/create" className="create-first-product">
                Добавить первый товар
              </Link>
            )}
          </div>
        )}
      </div>

      {/* Модальное окно */}
      {isModalOpen && (
        <div className="modal-overlay" onClick={handleCancel}>
          <div className="modal-content" onClick={(e) => e.stopPropagation()}>
            <h2>Подтверждение покупки</h2>
            <div className="modal-product-info">
              <img src={getFullImageUrl(selectedProduct.images[0])} alt={selectedProduct?.name} />
              <div>
                <h3>{selectedProduct?.name}</h3>
                <p className="price">{selectedProduct?.price.toLocaleString()} ₽</p>
              </div>
            </div>
            <p>Вы уверены, что хотите добавить этот товар в корзину?</p>
            <div className="modal-buttons">
              <button className="confirm-btn" onClick={handleConfirm}>
                Да, добавить в корзину
              </button>
              <button className="cancel-btn" onClick={handleCancel}>
                Отмена
              </button>
            </div>
          </div>
        </div>
      )}
    </>
  );
}

function ShoppingList_box({ children }) {
  return (
    <div className='ShoppingList_box'>
      {children}
    </div>
  );
}

export default App;
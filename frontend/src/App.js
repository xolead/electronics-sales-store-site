import './App.css';
import React, { useState, useEffect } from 'react';
import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom';
import Create from './pages/Create';

const App = () => {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/create" element={<Create />} />
      </Routes>
    </Router>
  );
}

// Функция для получения всех товаров
const getAll = async () => {
  try {
    const response = await fetch('http://localhost:8080/products', {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const data = await response.json();
    return data;
    
  } catch (error) {
    console.error('Ошибка при получении товаров:', error);
    throw error;
  }
};

// Альтернативный вариант с обработкой разных статусов
const getAllWithDetailedErrorHandling = async () => {
  try {
    const response = await fetch('http://localhost:8080/products');

    if (response.status === 404) {
      throw new Error('Сервер не найден (404)');
    }

    if (response.status === 500) {
      throw new Error('Ошибка сервера (500)');
    }

    if (!response.ok) {
      throw new Error(`Ошибка HTTP: ${response.status}`);
    }

    const data = await response.json();
    return data;

  } catch (error) {
    console.error('Ошибка при загрузке товаров:', error);
    
    // Можно показать уведомление пользователю
    if (error.message.includes('Failed to fetch')) {
      alert('Не удалось подключиться к серверу. Проверьте, запущен ли бэкенд.');
    }
    
    throw error;
  }
};

function HomePage() {
  return (
    <div className="App">
      <Header />
      <ShoppingList_box>
        <ShoppingList />
      </ShoppingList_box>
    </div>
  );
}

function Header() {
  return (
    <>
      <div className="header">
        <div className='header_box'>
          <img src="/img/cart.png" className='cart' alt="Cart" />
          <Link to="/create" className="create-link">
            Добавить  
          </Link>
        </div>
      </div>
    </>
  );
}

function ShoppingList() {
  const [selectedProduct, setSelectedProduct] = useState(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [products, setProducts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  // Загрузка товаров при монтировании компонента
  useEffect(() => {
    loadProducts();
  }, []);

  const loadProducts = async () => {
    try {
      setLoading(true);
      setError(null);
      
      // Используем функцию getAll для загрузки товаров
      const productsData = await getAll();
      setProducts(productsData);
      
    } catch (err) {
      setError('Не удалось загрузить товары');
      console.error('Ошибка загрузки:', err);
      
      // Если API не доступно, используем mock данные
      setProducts(getMockProducts());
    } finally {
      setLoading(false);
    }
  };

  // Mock данные на случай если бэкенд не доступен
  const getMockProducts = () => [
    { id: 1, name: "iPhone 14", price: 79999, image: "/img/iphone_14.jpg", category: "Смартфоны" },
    { id: 2, name: "MacBook Air", price: 99999, image: "/img/macbook_air.jpg", category: "Ноутбуки" },
    { id: 3, name: "AirPods", price: 12999, image: "/img/airpods.jpg", category: "Аксессуары" },
    { id: 4, name: "iPad", price: 39999, image: "/img/ipad.jpg", category: "Планшеты" },
    { id: 5, name: "Samsung Galaxy", price: 69999, image: "/img/iphone_14.jpg", category: "Смартфоны" },
    { id: 6, name: "Dell XPS", price: 89999, image: "/img/macbook_air.jpg", category: "Ноутбуки" },
  ];

  const handleBuyClick = (product) => {
    setSelectedProduct(product);
    setIsModalOpen(true);
  };

  const handleConfirm = () => {
    alert(`Товар "${selectedProduct.name}" добавлен в корзину!`);
    setIsModalOpen(false);
    setSelectedProduct(null);
  };

  const handleCancel = () => {
    setIsModalOpen(false);
    setSelectedProduct(null);
  };

  // Функция для обновления списка товаров (можно вызвать после добавления нового товара)
  const refreshProducts = () => {
    loadProducts();
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
        <button onClick={loadProducts} className="retry-btn">
          Попробовать снова
        </button>
      </div>
    );
  }

  return (
    <>
      <div className="shop_products">
        {/* Кнопка обновления (опционально) */}
        <div className="products-header">
          <h2>Список товаров ({products.length})</h2>
          <button onClick={refreshProducts} className="refresh-btn">
            Обновить
          </button>
        </div>

        <div className="products-grid">
          {products.map(product => (
            <div key={product.id} className="product-card">
              <img 
                src={product.image} 
                alt={product.name} 
                className="product-image" 
                onError={(e) => {
                  // Запасное изображение если основное не загрузилось
                  e.target.src = '/img/placeholder.jpg';
                }}
              />
              <div className="product-details">
                <span className="category">{product.category}</span>
                <h3>{product.name}</h3>
                <div className="price-section">
                  <span className="price">{product.price.toLocaleString()} ₽</span>
                  <button 
                    className="buy-btn" 
                    onClick={() => handleBuyClick(product)}
                  >
                    Купить
                  </button>
                </div>
              </div>
            </div>
          ))}
        </div>

        {products.length === 0 && !loading && (
          <div className="empty-state">
            <p>Товары не найдены</p>
            <Link to="/create" className="create-first-product">
              Добавить первый товар
            </Link>
          </div>
        )}
      </div>

      {/* Модальное окно */}
      {isModalOpen && (
        <div className="modal-overlay" onClick={handleCancel}>
          <div className="modal-content" onClick={(e) => e.stopPropagation()}>
            <h2>Подтверждение покупки</h2>
            <div className="modal-product-info">
              <img src={selectedProduct?.image} alt={selectedProduct?.name} />
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
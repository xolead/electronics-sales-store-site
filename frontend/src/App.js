import './App.css';
import React, { useState, useEffect } from 'react';
import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom';
import Create from './pages/Create';
import Cart from './pages/Cart';
import ProductDetail from './pages/ProductDetail';
import axios from 'axios';


const api = axios.create({
});


const DeleteProduct = async (id) => {
  await axios.delete('/product/' + id)
}

// базовый URL S3
const getFullImageUrl = (filename) => {
  const url = `https://electronic.s3.regru.cloud/products/${filename}`;
  return url;
};

const App = () => {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/create" element={<Create />} />
        <Route path="/Cart" element={<Cart />} />
        <Route path="/product/:id" element={<ProductDetail />} />
      </Routes>
    </Router>
  );
}


const getAll = async () => {
  try {
    console.log('🔄 Запрашиваем товары...');
    const response = await axios.get('/product');
    console.log('📦 Полный ответ:', response.data);
    
    if (response.data && response.data.Products) {
      console.log('✅ Товары найдены:', response.data.Products);
      return response.data.Products;
    } else if (response.data && Array.isArray(response.data)) {
      console.log('✅ Товары (массив):', response.data);
      return response.data;
    } else {
      console.warn('⚠️ Товары не найдены в ответе');
      return [];
    }
    
  } catch (error) {
    console.error('❌ Ошибка при получении товаров:', error);
    return [];
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
        <Link to="/cart" className="cart-link">
          <img src="/img/cart.png" className='cart' alt="Cart" />
          </Link>
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
 
  // for (let i = 1; i < 15; i ++){
  //   DeleteProduct(i)
  // }
  

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

  const getCategoryFromParameters = (parametersString) => {
    if (!parametersString) return '';
      // Разделяем строку по |
      const pairs = parametersString.split('|');
      
      // Ищем параметр с ключом "Категория"
      for (let pair of pairs) {
        const [key, value] = pair.split('=');
        if (key && value && key.trim() === 'Категория') {
          return value.trim();
        }
      }
    }

  const handleBuyClick = (product) => {
    // 1. Открываем модальное окно (оригинальный функционал)
    setSelectedProduct(product);
    setIsModalOpen(true);
  };
  
  const handleConfirm = () => {
    // 2. Добавляем товар в корзину при подтверждении
    const cartItem = {
      ...selectedProduct,
      quantity: 1
    };
    
    // Получаем текущую корзину из localStorage
    const existingCart = JSON.parse(localStorage.getItem('electronic_cart') || '[]');
    
    // Проверяем, есть ли уже такой товар в корзине
    const existingItemIndex = existingCart.findIndex(item => item.id === selectedProduct.id);
    
    if (existingItemIndex >= 0) {
      // Увеличиваем количество, если товар уже в корзине
      existingCart[existingItemIndex].quantity += 1;
    } else {
      // Добавляем новый товар
      existingCart.push(cartItem);
    }
    
    // Сохраняем обновленную корзину
    localStorage.setItem('electronic_cart', JSON.stringify(existingCart));
    setIsModalOpen(false);
    setSelectedProduct(null);
  };

  const handleCancel = () => {
    setIsModalOpen(false);
    setSelectedProduct(null);
  };

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
                {/* Обертка для кликабельной области */}
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
            
                {/* Секция с ценой и кнопкой - НЕ кликабельная */}
                <div className="price-section">
                  <span className="price">{product.price.toLocaleString()} ₽</span>
                  <button 
                    className="buy-btn" 
                    onClick={(e) => {
                      e.stopPropagation(); // Останавливаем всплытие события
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
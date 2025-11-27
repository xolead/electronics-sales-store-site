// pages/ProductDetail.js
import React, { useState, useEffect } from 'react';
import { useParams, Link } from 'react-router-dom';
import axios from 'axios';
import './ProductDetail.css';
import Header from '../components/layout/Header/Header'

// базовый URL S3
const getFullImageUrl = (filename) => {
  const url = `https://electronic.s3.regru.cloud/products/${filename}`;
  return url;
};

// Функция для парсинга всех параметров
const parseAllParameters = (parametersString) => {
  if (!parametersString) return [];
  
  const parameters = [];
  
  try {
    const pairs = parametersString.split('|');
    
    pairs.forEach(pair => {
      // Убираем лишние пробелы и разбиваем по первому знаку =
      const trimmedPair = pair.trim();
      const equalsIndex = trimmedPair.indexOf('=');
      
      if (equalsIndex > 0) {
        const key = trimmedPair.substring(0, equalsIndex).trim();
        const value = trimmedPair.substring(equalsIndex + 1).trim();
        
        if (key && value) {
          parameters.push({ key, value });
        }
      }
    });
    
    return parameters;
  } catch (error) {
    console.error('Ошибка парсинга параметров:', error);
    return [];
  }
};

// Хук для отслеживания корзины
const useCartCount = () => {
  const [cartCount, setCartCount] = useState(0);

  // Функция для обновления количества товаров в корзине
  const updateCartCount = () => {
    const cart = JSON.parse(localStorage.getItem('electronic_cart') || '[]');
    // Подсчитываем количество различных товаров (по id)
    const uniqueItemsCount = cart.length;
    setCartCount(uniqueItemsCount);
  };

  // Слушаем изменения в localStorage
  useEffect(() => {
    updateCartCount();
    
    // Функция для обработки событий storage
    const handleStorageChange = (e) => {
      if (e.key === 'electronic_cart') {
        updateCartCount();
      }
    };

    // Слушаем события storage (из других вкладок)
    window.addEventListener('storage', handleStorageChange);
    
    // Слушаем custom event (из этой же вкладки)
    window.addEventListener('cartUpdated', updateCartCount);

    return () => {
      window.removeEventListener('storage', handleStorageChange);
      window.removeEventListener('cartUpdated', updateCartCount);
    };
  }, []);

  return cartCount;
};

const ProductDetail = () => {
  const { id } = useParams();
  const [product, setProduct] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [currentImageIndex, setCurrentImageIndex] = useState(0);
  const [productParameters, setProductParameters] = useState([]);

  useEffect(() => {
    loadProduct();
  }, [id]);

  useEffect(() => {
    if (product && product.parameters) {
      const parsedParams = parseAllParameters(product.parameters);
      setProductParameters(parsedParams);
    }
  }, [product]);

  const loadProduct = async () => {
    try {
      setLoading(true);
      setError(null);
      
      const response = await axios.get(`/product/${id}`);
      console.log('📦 Получен товар:', response.data);
      
      if (response.data) {
        setProduct(response.data);
      } else {
        setError('Товар не найден');
      }
      
    } catch (err) {
      console.error('❌ Ошибка загрузки товара:', err);
      setError('Не удалось загрузить информацию о товаре');
      
      // Mock данные для демонстрации
      setProduct({
        id: parseInt(id),
        name: "iPhone 14 Pro",
        description: "Современный смартфон с передовыми технологиями. Оснащен мощным процессором, улучшенной камерой и длительным временем работы от батареи.",
        parameters: "Категория=Смартфоны|Цвет=Черный|Память=128ГБ|Материал=Алюминий|Экран=6.1 дюймов",
        count: 10,
        price: 79999,
        images: [
          "iphone_14_pro_1.jpg",
          "iphone_14_pro_2.jpg",
          "iphone_14_pro_3.jpg",
          "iphone_14_pro_4.jpg"
        ]
      });
    } finally {
      setLoading(false);
    }
  };

  // Функции для карусели
  const nextImage = () => {
    if (product && product.images && product.images.length > 0) {
      setCurrentImageIndex((prevIndex) => 
        prevIndex === product.images.length - 1 ? 0 : prevIndex + 1
      );
    }
  };

  const prevImage = () => {
    if (product && product.images && product.images.length > 0) {
      setCurrentImageIndex((prevIndex) => 
        prevIndex === 0 ? product.images.length - 1 : prevIndex - 1
      );
    }
  };

  const goToImage = (index) => {
    setCurrentImageIndex(index);
  };

  const handleBuyClick = () => {
    if (!product) return;

    const cartItem = {
      ...product,
      quantity: 1
    };
    
    const existingCart = JSON.parse(localStorage.getItem('electronic_cart') || '[]');
    const existingItemIndex = existingCart.findIndex(item => item.id === product.id);
    
    if (existingItemIndex >= 0) {
      existingCart[existingItemIndex].quantity += 1;
    } else {
      existingCart.push(cartItem);
    }
    
    localStorage.setItem('electronic_cart', JSON.stringify(existingCart));
    
    // Триггерим событие обновления корзины
    window.dispatchEvent(new Event('cartUpdated'));
    
    alert(`Товар "${product.name}" добавлен в корзину!`);
  };

  if (loading) {
    return (
      <div className="product-detail-page">
        <Header />
        <div className="product-detail-container">
          <div className="loading-container">
            <div className="loading-spinner"></div>
            <p>Загрузка товара...</p>
          </div>
        </div>
      </div>
    );
  }

  if (error && !product) {
    return (
      <div className="product-detail-page">
        <Header />
        <div className="product-detail-container">
          <div className="error-container">
            <p>{error}</p>
            <button onClick={loadProduct} className="retry-btn">
              Попробовать снова
            </button>
            <Link to="/" className="back-link">
              Вернуться на главную
            </Link>
          </div>
        </div>
      </div>
    );
  }

  if (!product) {
    return (
      <div className="product-detail-page">
        <Header />
        <div className="product-detail-container">
          <div className="empty-state">
            <p>Товар не найден</p>
            <Link to="/" className="back-link">
              Вернуться на главную
            </Link>
          </div>
        </div>
      </div>
    );
  }

  const hasImages = product.images && product.images.length > 0;
  const currentImage = hasImages ? product.images[currentImageIndex] : null;
  const category = productParameters.find(param => param.key === 'Категория')?.value || 'Без категории';

  return (
    <div className="product-detail-page">
      <Header />
      <div className="product-detail-container">
        <div className="breadcrumb">
          <Link to="/" className="breadcrumb-link">Главная</Link>
          <span className="breadcrumb-separator">/</span>
          <span className="breadcrumb-current">{product.name}</span>
        </div>

        <div className="product-detail-content">
          {/* Левая колонка - карусель изображений */}
          <div className="product-images">
            <div className="carousel-container">
              <div className="carousel-slide">
                {currentImage ? (
                  <img 
                    src={getFullImageUrl(currentImage)} 
                    alt={`${product.name} - фото ${currentImageIndex + 1}`}
                    onError={(e) => {
                      e.target.src = '/img/placeholder.jpg';
                    }}
                  />
                ) : (
                  <img 
                    src="/img/placeholder.jpg" 
                    alt="Нет изображения"
                  />
                )}
              </div>
              
              {/* Навигация карусели */}
              {hasImages && product.images.length > 1 && (
                <div className="carousel-nav">
                  <button 
                    className="carousel-btn" 
                    onClick={prevImage}
                  >
                    ‹
                  </button>
                  
                  <div className="carousel-indicators">
                    {product.images.map((_, index) => (
                      <div
                        key={index}
                        className={`carousel-dot ${index === currentImageIndex ? 'active' : ''}`}
                        onClick={() => goToImage(index)}
                      />
                    ))}
                  </div>
                  
                  <div className="carousel-counter">
                    {currentImageIndex + 1} / {product.images.length}
                  </div>
                  
                  <button 
                    className="carousel-btn" 
                    onClick={nextImage}
                  >
                    ›
                  </button>
                </div>
              )}
            </div>
            
            {/* Миниатюры */}
            {hasImages && product.images.length > 1 && (
              <div className="image-thumbnails">
                {product.images.map((image, index) => (
                  <div 
                    key={index}
                    className={`thumbnail ${index === currentImageIndex ? 'active' : ''}`}
                    onClick={() => goToImage(index)}
                  >
                    <img 
                      src={getFullImageUrl(image)} 
                      alt={`${product.name} - миниатюра ${index + 1}`}
                      onError={(e) => {
                        e.target.src = '/img/placeholder.jpg';
                      }}
                    />
                  </div>
                ))}
              </div>
            )}
          </div>

          {/* Правая колонка - информация о товаре */}
          <div className="product-info">
            <span className="product-category">{category}</span>
            <h1 className="product-title">{product.name}</h1>
            
            <div className="product-price-section">
              <span className="product-price">{product.price?.toLocaleString()} ₽</span>
              <div className="stock-info">
                {product.count > 0 ? (
                  <span className="in-stock">В наличии ({product.count} шт.)</span>
                ) : (
                  <span className="out-of-stock">Нет в наличии</span>
                )}
              </div>
            </div>

            <div className="product-actions">
              <button 
                className="buy-btn-large"
                onClick={handleBuyClick}
                disabled={product.count === 0}
              >
                {product.count > 0 ? 'Добавить в корзину' : 'Нет в наличии'}
              </button>
            </div>

            {/* Описание товара */}
            {product.description && (
              <div className="product-description">
                <h3>Описание</h3>
                <p>{product.description}</p>
              </div>
            )}

            {/* Параметры товара */}
            <div className="product-parameters">
              <h3>Характеристики</h3>
              <div className="parameters-list">
                {/* Отображаем все распарсенные параметры кроме категории */}
                {productParameters
                  .filter(param => param.key !== 'Категория')
                  .map((param, index) => (
                    <div key={index} className="parameter-item">
                      <span className="parameter-name">{param.key}:</span>
                      <span className="parameter-value">{param.value}</span>
                    </div>
                  ))
                }
                
                {/* Базовые параметры, которые всегда показываем */}
                <div className="parameter-item">
                  <span className="parameter-name">Количество на складе:</span>
                  <span className="parameter-value">{product.count} шт.</span>
                </div>
                
                {/* Если нет дополнительных параметров */}
                {productParameters.filter(param => param.key !== 'Категория').length === 0 && (
                  <div className="no-parameters-message">
                    <p>Дополнительные характеристики не указаны</p>
                  </div>
                )}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};


export default ProductDetail;
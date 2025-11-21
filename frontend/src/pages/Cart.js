import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import './Cart.css';
import axios from 'axios';

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

// Компонент Header
const Header = () => {
  const cartCount = useCartCount();

  return (
    <div className="header">
      <div className='header_box'>
        <Link to="/cart" className="cart-link">
          <div style={{ position: 'relative', display: 'inline-block' }}>
            <img src="/img/cart.png" className='cart' alt="Cart" />
            {cartCount > 0 && (
                <span 
                  style={{
                    position: 'absolute',
                    top: '-5px',
                    right: '-5px',
                    backgroundColor: '#ff4444',
                    color: 'white',
                    borderRadius: '50%',
                    width: '20px',
                    height: '20px',
                    fontSize: '12px',
                    fontWeight: 'bold',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    boxShadow: '0 2px 4px rgba(0,0,0,0.2)'
                  }}
                >
                  {cartCount}
                </span>
              )}
          </div>
        </Link>
        <Link to="/" className="create-link">
          Главная  
        </Link>
      </div>
    </div>
  );
};

const Cart = () => {
  const [cartItems, setCartItems] = useState([]);
  const [loading, setLoading] = useState(true);
  const [totalPrice, setTotalPrice] = useState(0);
  const [isCheckingOut, setIsCheckingOut] = useState(false);
  const [productStocks, setProductStocks] = useState({}); // Храним остатки товаров

  // Загрузка корзины из localStorage при монтировании
  useEffect(() => {
    loadCartItems();
  }, []);

  // Пересчет общей суммы при изменении корзины
  useEffect(() => {
    calculateTotalPrice();
  }, [cartItems]);

  // Загрузка актуальных остатков товаров
  useEffect(() => {
    if (cartItems.length > 0) {
      loadProductStocks();
    }
  }, [cartItems]);

  const loadCartItems = () => {
    try {
      const savedCart = localStorage.getItem('electronic_cart');
      if (savedCart) {
        const cartData = JSON.parse(savedCart);
        setCartItems(cartData);
      }
    } catch (error) {
      console.error('Ошибка загрузки корзины:', error);
    } finally {
      setLoading(false);
    }
  };

  // Загрузка актуальных остатков товаров с сервера
  const loadProductStocks = async () => {
    try {
      const stockPromises = cartItems.map(item =>
        axios.get(`/product/${item.id}`)
          .then(response => ({
            id: item.id,
            stock: response.data.count
          }))
          .catch(error => {
            console.error(`Ошибка загрузки товара ${item.id}:`, error);
            return {
              id: item.id,
              stock: item.count || 0 // Используем сохраненное значение как fallback
            };
          })
      );

      const stocks = await Promise.all(stockPromises);
      const stockMap = {};
      stocks.forEach(stock => {
        stockMap[stock.id] = stock.stock;
      });
      setProductStocks(stockMap);
    } catch (error) {
      console.error('Ошибка загрузки остатков товаров:', error);
    }
  };

  const saveCartToStorage = (items) => {
    localStorage.setItem('electronic_cart', JSON.stringify(items));
    // Триггерим событие обновления корзины
    window.dispatchEvent(new Event('cartUpdated'));
  };

  const calculateTotalPrice = () => {
    const total = cartItems.reduce((sum, item) => {
      return sum + (item.price * item.quantity);
    }, 0);
    setTotalPrice(total);
  };

  const updateQuantity = (id, newQuantity) => {
    if (newQuantity < 1) return;

    const availableStock = productStocks[id] || 0;
    const currentCartItem = cartItems.find(item => item.id === id);
    const currentQuantity = currentCartItem ? currentCartItem.quantity : 0;

    // Проверяем, не превышает ли новое количество доступный остаток
    if (newQuantity > availableStock) {
      alert(`❌ Нельзя добавить больше ${availableStock} шт. этого товара\nДоступно на складе: ${availableStock} шт.`);
      return;
    }

    const updatedCart = cartItems.map(item => 
      item.id === id ? { ...item, quantity: newQuantity } : item
    );
    
    setCartItems(updatedCart);
    saveCartToStorage(updatedCart);
  };

  const removeFromCart = (id) => {
    const updatedCart = cartItems.filter(item => item.id !== id);
    setCartItems(updatedCart);
    saveCartToStorage(updatedCart);
    
    // Удаляем из кэша остатков
    setProductStocks(prev => {
      const newStocks = { ...prev };
      delete newStocks[id];
      return newStocks;
    });
  };

  const clearCart = () => {
    setCartItems([]);
    setProductStocks({});
    localStorage.removeItem('electronic_cart');
    // Триггерим событие обновления корзины
    window.dispatchEvent(new Event('cartUpdated'));
  };

  const getFullImageUrl = (filename) => {
    return `https://electronic.s3.regru.cloud/products/${filename}`;
  };

  // Функция для получения категории из параметров
  const getCategoryFromParameters = (parametersString) => {
    if (!parametersString) return '';
    
    try {
      const pairs = parametersString.split('|');
      
      for (let pair of pairs) {
        const [key, value] = pair.split('=');
        if (key && value && key.trim() === 'Категория') {
          return value.trim();
        }
      }
      
      return '';
    } catch (error) {
      console.error('Ошибка парсинга категории:', error);
      return '';
    }
  };

  // Функция для отправки запроса на изменение количества товара
  const updateProductCountOnServer = async (productId, quantityChange) => {
    try {
      const response = await axios.put('/product/change', {
        ID: productId,
        Count: -quantityChange // Отправляем отрицательное значение
      }, {
        headers: {
          'Content-Type': 'application/json'
        }
      });

      console.log(`✅ Количество товара ${productId} уменьшено на ${quantityChange}`);
      return response.data;
    } catch (error) {
      console.error(`❌ Ошибка при обновлении количества товара ${productId}:`, error);
      throw error;
    }
  };

  // Проверка доступности всех товаров перед оформлением заказа
  const validateCartBeforeCheckout = () => {
    const errors = [];

    cartItems.forEach(item => {
      const availableStock = productStocks[item.id] || 0;
      if (item.quantity > availableStock) {
        errors.push({
          productName: item.name,
          requested: item.quantity,
          available: availableStock
        });
      }
    });

    return errors;
  };

  const handleCheckout = async () => {
    if (cartItems.length === 0) {
      alert('Корзина пуста!');
      return;
    }

    // Проверяем доступность товаров
    const validationErrors = validateCartBeforeCheckout();
    if (validationErrors.length > 0) {
      const errorMessage = validationErrors.map(error => 
        `• ${error.productName}: запрошено ${error.requested} шт., доступно ${error.available} шт.`
      ).join('\n');
      
      alert(`❌ Недостаточно товаров на складе:\n\n${errorMessage}\n\nПожалуйста, измените количество товаров в корзине.`);
      return;
    }

    setIsCheckingOut(true);

    try {
      // Отправляем запросы для каждого товара в корзине
      const updatePromises = cartItems.map(item => 
        updateProductCountOnServer(item.id, item.quantity)
      );

      // Ждем завершения всех запросов
      await Promise.all(updatePromises);

      alert(`✅ Заказ оформлен!\nОбщая сумма: ${totalPrice.toLocaleString()} ₽\nТовары: ${cartItems.reduce((sum, item) => sum + item.quantity, 0)} шт.\n\nСпасибо за покупку!`);
      
      // Очищаем корзину после успешного оформления
      clearCart();
      
    } catch (error) {
      console.error('Ошибка при оформлении заказа:', error);
      alert('❌ Произошла ошибка при оформлении заказа. Пожалуйста, попробуйте еще раз.');
    } finally {
      setIsCheckingOut(false);
    }
  };

  // Функция для отображения информации о доступном количестве
  const getStockInfo = (item) => {
    const availableStock = productStocks[item.id];
    
    if (availableStock === undefined) {
      return <span className="stock-loading">Загрузка...</span>;
    }
    
    if (availableStock === 0) {
      return <span className="stock-out">Нет в наличии</span>;
    }
    
    if (item.quantity > availableStock) {
      return <span className="stock-warning">Доступно: {availableStock} шт.</span>;
    }
    
    return <span className="stock-available">Доступно: {availableStock} шт.</span>;
  };

  if (loading) {
    return (
      <div className="cart-page">
        <Header />
        <div className="cart-container">
          <div className="loading-container">
            <div className="loading-spinner"></div>
            <p>Загрузка корзины...</p>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="cart-page">
      <Header />
      
      <div className="cart-container">
        <div className="cart-header">
          <h1>Корзина покупок</h1>
          {cartItems.length > 0 && (
            <button 
              className="clear-cart-btn"
              onClick={clearCart}
            >
              Очистить корзину
            </button>
          )}
        </div>

        {cartItems.length === 0 ? (
          <div className="empty-cart">
            <div className="empty-cart-icon">🛒</div>
            <h2>Ваша корзина пуста</h2>
            <p>Добавьте товары из каталога, чтобы сделать заказ</p>
            <Link to="/" className="continue-shopping-btn">
              Продолжить покупки
            </Link>
          </div>
        ) : (
          <div className="cart-content">
            <div className="cart-items">
              {cartItems.map(item => {
                const availableStock = productStocks[item.id] || 0;
                const isOutOfStock = availableStock === 0;
                const exceedsStock = item.quantity > availableStock;
                
                return (
                  <div key={item.id} className={`cart-item ${exceedsStock ? 'exceeds-stock' : ''}`}>
                    <div className="item-image">
                      <img 
                        src={getFullImageUrl(item.images[0])} 
                        alt={item.name}
                        onError={(e) => {
                          e.target.src = '/img/placeholder.jpg';
                        }}
                      />
                    </div>
                    
                    <div className="item-details">
                      <h3 className="item-name">{item.name}</h3>
                      <p className="item-category">
                        {getCategoryFromParameters(item.parameters)}
                      </p>
                      <div className="stock-info">
                        {getStockInfo(item)}
                      </div>
                      {item.description && (
                        <p className="item-description">{item.description}</p>
                      )}
                    </div>

                    <div className="item-controls">
                      <div className="quantity-controls">
                        <button 
                          className="quantity-btn"
                          onClick={() => updateQuantity(item.id, item.quantity - 1)}
                          disabled={item.quantity <= 1 || isOutOfStock}
                        >
                          -
                        </button>
                        <span className={`quantity ${exceedsStock ? 'exceeds' : ''}`}>
                          {item.quantity}
                        </span>
                        <button 
                          className="quantity-btn"
                          onClick={() => updateQuantity(item.id, item.quantity + 1)}
                          disabled={isOutOfStock || item.quantity >= availableStock}
                        >
                          +
                        </button>
                      </div>

                      <div className="item-price">
                        {(item.price * item.quantity).toLocaleString()} ₽
                      </div>

                      <button 
                        className="remove-btn"
                        onClick={() => removeFromCart(item.id)}
                      >
                        Удалить
                      </button>
                    </div>
                  </div>
                );
              })}
            </div>

            <div className="cart-summary">
              <div className="summary-card">
                <h3>Итого</h3>
                <div className="summary-row">
                  <span>Товары ({cartItems.reduce((sum, item) => sum + item.quantity, 0)} шт.)</span>
                  <span>{totalPrice.toLocaleString()} ₽</span>
                </div>
                <div className="summary-row">
                  <span>Доставка</span>
                  <span>Бесплатно</span>
                </div>
                <div className="summary-divider"></div>
                <div className="summary-total">
                  <span>Общая сумма</span>
                  <span className="total-price">{totalPrice.toLocaleString()} ₽</span>
                </div>
                
                {/* Предупреждение о недостатке товаров */}
                {cartItems.some(item => {
                  const availableStock = productStocks[item.id] || 0;
                  return item.quantity > availableStock;
                }) && (
                  <div className="checkout-warning">
                    ⚠️ Некоторые товары недоступны в запрошенном количестве
                  </div>
                )}
                
                <button 
                  className="checkout-btn"
                  onClick={handleCheckout}
                  disabled={isCheckingOut || cartItems.some(item => {
                    const availableStock = productStocks[item.id] || 0;
                    return item.quantity > availableStock;
                  })}
                >
                  {isCheckingOut ? (
                    <>
                      <div className="loading-spinner" style={{width: '20px', height: '20px', display: 'inline-block', marginRight: '10px'}}></div>
                      Оформление...
                    </>
                  ) : (
                    'Оформить заказ'
                  )}
                </button>
                <Link to="/" className="continue-shopping-link">
                  ← Продолжить покупки
                </Link>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default Cart;
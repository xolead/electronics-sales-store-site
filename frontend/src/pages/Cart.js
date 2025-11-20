import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import './Cart.css';
import axios from 'axios';

const Cart = () => {
  const [cartItems, setCartItems] = useState([]);
  const [loading, setLoading] = useState(true);
  const [totalPrice, setTotalPrice] = useState(0);

  // Загрузка корзины из localStorage при монтировании
  useEffect(() => {
    loadCartItems();
  }, []);

  // Пересчет общей суммы при изменении корзины
  useEffect(() => {
    calculateTotalPrice();
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

  const saveCartToStorage = (items) => {
    localStorage.setItem('electronic_cart', JSON.stringify(items));
  };

  const calculateTotalPrice = () => {
    const total = cartItems.reduce((sum, item) => {
      return sum + (item.price * item.quantity);
    }, 0);
    setTotalPrice(total);
  };

  const updateQuantity = (id, newQuantity) => {
    if (newQuantity < 1) return;

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
  };

  const clearCart = () => {
    setCartItems([]);
    localStorage.removeItem('electronic_cart');
  };

  const getFullImageUrl = (filename) => {
    return `https://electronic.s3.regru.cloud/products/${filename}`;
  };

  const handleCheckout = () => {
    if (cartItems.length === 0) {
      alert('Корзина пуста!');
      return;
    }

    alert(`Заказ оформлен! Общая сумма: ${totalPrice.toLocaleString()} ₽\nСпасибо за покупку!`);
    clearCart();
  };

  if (loading) {
    return (
      <div className="cart-page">
        <header className="header">
          <div className='header_box'>
            <img src="/img/cart.png" className='cart' alt="Cart" />
            <Link to="/" className="home-link">
              Главная
            </Link>
          </div>
        </header>
        <div className="loading-container">
          <div className="loading-spinner"></div>
          <p>Загрузка корзины...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="cart-page">
      <header className="header">
        <div className='header_box'>
          <img src="/img/cart.png" className='cart' alt="Cart" />
          <Link to="/" className="home-link">
            Главная
          </Link>
        </div>
      </header>

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
              {cartItems.map(item => (
                <div key={item.id} className="cart-item">
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
                    <p className="item-category">{item.parameters}</p>
                    <p className="item-description">{item.description}</p>
                  </div>

                  <div className="item-controls">
                    <div className="quantity-controls">
                      <button 
                        className="quantity-btn"
                        onClick={() => updateQuantity(item.id, item.quantity - 1)}
                        disabled={item.quantity <= 1}
                      >
                        -
                      </button>
                      <span className="quantity">{item.quantity}</span>
                      <button 
                        className="quantity-btn"
                        onClick={() => updateQuantity(item.id, item.quantity + 1)}
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
              ))}
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
                <button 
                  className="checkout-btn"
                  onClick={handleCheckout}
                >
                  Оформить заказ
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
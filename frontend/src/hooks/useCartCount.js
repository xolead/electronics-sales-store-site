import { useState, useEffect } from 'react';


export const useCartCount = () => {
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


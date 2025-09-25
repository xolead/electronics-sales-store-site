import './App.css';
import React, { useState } from 'react';


const App = () => {
  return (
    <div className="App">
        <Header />
          <ShoppingList_box>
            <ShoppingList />
          </ShoppingList_box>
    </div>
  );
}


function Header () {
  return(
    <>
    <div className="header">
      <div className='header_box'>
        <img src="/img/cart.png" className='cart'></img>
      </div>
    </div>
    </>
  );
}

function ShoppingList () {

  const [selectedProduct, setSelectedProduct] = useState(null);
  const [isModalOpen, setIsModalOpen] = useState(false);

  const products = [
    { id: 1, name: "iPhone 14", price: 79999, image: "/img/iphone_14.jpg", category: "Смартфоны" },
    { id: 2, name: "MacBook Air", price: 99999, image: "/img/macbook_air.jpg", category: "Ноутбуки" },
    { id: 3, name: "AirPods", price: 12999, image: "/img/airpods.jpg", category: "Аксессуары" },
    { id: 4, name: "iPad", price: 39999, image: "/img/ipad.jpg", category: "Планшеты" },
    { id: 5, name: "iPhone 14", price: 79999, image: "/img/iphone_14.jpg", category: "Смартфоны" },
    { id: 6, name: "MacBook Air", price: 99999, image: "/img/macbook_air.jpg", category: "Ноутбуки" },
    { id: 7, name: "AirPods", price: 12999, image: "/img/airpods.jpg", category: "Аксессуары" },
    { id: 8, name: "iPad", price: 39999, image: "/img/ipad.jpg", category: "Планшеты" },
    { id: 9, name: "iPad", price: 39999, image: "/img/ipad.jpg", category: "Планшеты" },
    { id: 10, name: "iPhone 14", price: 79999, image: "/img/iphone_14.jpg", category: "Смартфоны" },
    { id: 11, name: "MacBook Air", price: 99999, image: "/img/macbook_air.jpg", category: "Ноутбуки" },
    { id: 12, name: "AirPods", price: 12999, image: "/img/airpods.jpg", category: "Аксессуары" },
    { id: 13, name: "iPad", price: 39999, image: "/img/ipad.jpg", category: "Планшеты" }
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


  return(
    <>
    <div className="shop_products">
      <div className="products-grid">
        {products.map(product => (
          <div key={product.id} className="product-card">
            <img src={product.image} alt={product.name} className="product-image" />
            <div className="product-details">
              <span className="category">{product.category}</span>
              <h3>{product.name}</h3>
              <div className="price-section">
                <span className="price">{product.price.toLocaleString()} ₽</span>
                <button className="buy-btn" onClick={() => handleBuyClick(product)}>Купить</button>
              </div>
            </div>
          </div>
        ))}
      </div>
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
            <button className="confirm-btn" onClick={handleConfirm}>Да, добавить в корзину</button>
            <button className="cancel-btn" onClick={handleCancel}>Отмена</button>
          </div>
        </div>
      </div>
    )}
    
    </>

  );
}

function ShoppingList_box ({children}) {
  return(
    <div className='ShoppingList_box'>
      {children}
    </div>
  );
}

export default App;

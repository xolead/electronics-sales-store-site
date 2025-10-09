import React from 'react';
import './App.css';

const About = () => {
  return (
    <div className="App">
      <Header />
      <div className="about-page">
        <div className="about-content">
          <div>О нашем магазине</div>
          <p>Мы лучший магазин электроники с 2020 года!</p>
          <p>У нас вы найдёте самые современные гаджеты по выгодным ценам.</p>
          <button onClick={() => window.history.back()} className="back-btn">
            Вернуться в магазин
          </button>
        </div>
      </div>
    </div>
  );
}

// Header компонент (такой же как в App.js)
function Header() {
  return(
    <>
      <div className="header">
        <div className='header_box'>
          <img src="/img/cart.png" className='cart' alt="Cart"></img>
          <a href="/" className="home-link">🏠 Главная</a>
        </div>
      </div>
    </>
  );
}

export default About;
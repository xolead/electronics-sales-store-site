import React from 'react';
import { Link } from 'react-router-dom';
import { useCartCount } from '../../../hooks/useCartCount';
import './Header.css'

function HeaderForLogined() {
    const cartCount = useCartCount();
  
    return (
      <>
      
        <div className="header">
          <div className='header_box'>

          <Link to="/" className="main-link">
              Главная  
            </Link>

            <Link to="/cart" className="cart-link">
              <div style={{ position: 'relative', display: 'inline-block' }}>
                <img src="/img/cart3.png" className='cart' alt="Cart" />
                {cartCount > 0 && (
                  <span 
                    style={{
                      position: 'absolute',
                      top: '-2px',
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

            <Link to ="/Personal" className=''>
            <img src='/img/loggined.png' className='registration_and_login'></img>
            </Link>
            
          </div>
        </div>
      </>
    );
  }

export default HeaderForLogined
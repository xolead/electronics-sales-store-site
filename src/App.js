import './App.css';
import React from 'react';


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
    return(
      <>
      <div className="shop_products">
          
      </div>
      </>
    );
}

function ShoppingList_box () {
  return(
    <div className='ShoppingList_box'></div>
  );
}

export default App;

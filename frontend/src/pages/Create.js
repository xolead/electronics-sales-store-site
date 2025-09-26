import React, { useState, useRef } from 'react';
import { Link } from 'react-router-dom';
import './Create.css';

const Create = () => {
  const [productData, setProductData] = useState({
    name: '',
    price: '',
    category: '',
    description: ''
  });
  const [previewImage, setPreviewImage] = useState(null);
  const [isDragging, setIsDragging] = useState(false);
  const fileInputRef = useRef(null);

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setProductData(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const handleFileSelect = (file) => {
    if (file && file.type.startsWith('image/')) {
      const reader = new FileReader();
      reader.onload = (e) => {
        setPreviewImage(e.target.result);
      };
      reader.readAsDataURL(file);
    }
  };

  const handleDragOver = (e) => {
    e.preventDefault();
    setIsDragging(true);
  };

  const handleDragLeave = (e) => {
    e.preventDefault();
    setIsDragging(false);
  };

  const handleDrop = (e) => {
    e.preventDefault();
    setIsDragging(false);
    
    const files = e.dataTransfer.files;
    if (files.length > 0) {
      handleFileSelect(files[0]);
    }
  };

  const handleFileInputChange = (e) => {
    const file = e.target.files[0];
    if (file) {
      handleFileSelect(file);
    }
  };

  const handleDragAreaClick = () => {
    fileInputRef.current?.click();
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    
    if (!previewImage) {
      alert('Пожалуйста, добавьте изображение товара');
      return;
    }

    // Здесь будет логика сохранения товара
    console.log('Данные товара:', {
      ...productData,
      image: previewImage
    });

    alert('Товар успешно добавлен!');
    
    // Очистка формы
    setProductData({
      name: '',
      price: '',
      category: '',
      description: ''
    });
    setPreviewImage(null);
  };

  return (
    <div className="App">
      <Header />
      <div className="create-page">
        <div className="create-container">
          <h1>Добавить новый товар</h1>
          
          <form onSubmit={handleSubmit} className="product-form">
            {/* Drag & Drop область для изображения */}
            <div className="image-upload-section">
              <h3>Изображение товара</h3>
              <div 
                className={`drop-zone ${isDragging ? 'dragging' : ''} ${previewImage ? 'has-image' : ''}`}
                onDragOver={handleDragOver}
                onDragLeave={handleDragLeave}
                onDrop={handleDrop}
                onClick={handleDragAreaClick}
              >
                <input 
                  type="file" 
                  ref={fileInputRef}
                  onChange={handleFileInputChange}
                  accept="image/*"
                  style={{ display: 'none' }}
                />
                
                {previewImage ? (
                  <div className="image-preview">
                    <img src={previewImage} alt="Preview" />
                    <button 
                      type="button" 
                      className="remove-image-btn"
                      onClick={(e) => {
                        e.stopPropagation();
                        setPreviewImage(null);
                      }}
                    >
                      ×
                    </button>
                  </div>
                ) : (
                  <div className="drop-zone-content">
                    <div className="drop-icon">📁</div>
                    <p>Перетащите изображение сюда или кликните для выбора</p>
                    <span>PNG, JPG, JPEG (макс. 5MB)</span>
                  </div>
                )}
              </div>
            </div>

            {/* Общая информация о товаре */}
            <div className="product-info-section">
              <h3>Информация о товаре</h3>
              
              <div className="form-group">
                <label htmlFor="name">Название товара *</label>
                <input
                  type="text"
                  id="name"
                  name="name"
                  value={productData.name}
                  onChange={handleInputChange}
                  required
                  placeholder="Например: iPhone 14 Pro"
                />
              </div>

              <div className="form-group">
                <label htmlFor="price">Цена (₽) *</label>
                <input
                  type="number"
                  id="price"
                  name="price"
                  value={productData.price}
                  onChange={handleInputChange}
                  required
                  placeholder="79999"
                  min="0"
                />
              </div>

              <div className="form-group">
                <label htmlFor="category">Категория *</label>
                <select
                  id="category"
                  name="category"
                  value={productData.category}
                  onChange={handleInputChange}
                  required
                >
                  <option value="">Выберите категорию</option>
                  <option value="Смартфоны">Смартфоны</option>
                  <option value="Ноутбуки">Ноутбуки</option>
                  <option value="Планшеты">Планшеты</option>
                  <option value="Аксессуары">Аксессуары</option>
                  <option value="Техника">Техника</option>
                </select>
              </div>

              <div className="form-group">
                <label htmlFor="description">Описание товара</label>
                <textarea
                  id="description"
                  name="description"
                  value={productData.description}
                  onChange={handleInputChange}
                  rows="4"
                  placeholder="Подробное описание товара, характеристики, преимущества..."
                />
              </div>
            </div>

            <div className="form-actions">
              <Link to="/" className="cancel-btn">
                Отмена
              </Link>
              <button type="submit" className="submit-btn">
                Добавить товар
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
}

// Header компонент
function Header() {
  return(
    <div className="header">
      <div className='header_box'>
        <img src="/img/cart.png" className='cart' alt="Cart" />
        <Link to="/" className="home-link">
          Главная
        </Link>
      </div>
    </div>
  );
}

export default Create;
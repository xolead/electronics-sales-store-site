import React, { useState, useRef } from 'react';
import { Link } from 'react-router-dom';
import './Create.css';
import axios from 'axios';

const Create = () => {
  const [productData, setProductData] = useState({
    name: '',
    price: '',
    category: '',
    description: ''
  });
  const [previewImages, setPreviewImages] = useState([]);
  const [isDragging, setIsDragging] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const fileInputRef = useRef(null);

  // Функция создания продукта - получаем URLs для загрузки
  const createProductAndGetUrls = async (productData, fileNames) => {
    try {
      const response = await axios.post('/product', {
        name: productData.name,
        price: Number(productData.price),
        description: productData.description,
        parameters: productData.category || '',
        count: Number(productData.count) || 1,
        images: fileNames // МАССИВ НАЗВАНИЙ КАРТИНОК
      }, {
        headers: {
          'Content-Type': 'application/json',
        }, 
      });

      return response.data; // { URLs: ["s3-url1", "s3-url2", ...] }

    } catch (error) {
      console.error('Ошибка при создании товара:', error);
      throw error;
    }
  };

  // Функция загрузки файлов на S3
  const uploadFilesToS3 = async (files, s3Urls) => {
    for (let i = 0; i < files.length; i++) {
      try {
        const file = files[i];
        const s3Url = s3Urls[i];

        // Загружаем файл на S3 используя полученный URL
        await axios.put(s3Url, file, {
          headers: {
            'Content-Type': file.type,
          },
        });

        console.log(`✅ Файл "${file.name}" загружен на S3`);

      } catch (error) {
        console.error(`Ошибка при загрузке файла ${files[i].name}:`, error);
        throw new Error(`Не удалось загрузить файл: ${files[i].name}`);
      }
    }
  };

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setProductData(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const handleFileSelect = (files) => {
    const imageFiles = Array.from(files).filter(file => file.type.startsWith('image/'));
    
    if (imageFiles.length === 0) return;

    const newImages = imageFiles.map(file => {
      return new Promise((resolve) => {
        const reader = new FileReader();
        reader.onload = (e) => {
          resolve({
            id: Date.now() + Math.random(),
            url: e.target.result,
            file: file,
            fileName: file.name
          });
        };
        reader.readAsDataURL(file);
      });
    });

    Promise.all(newImages).then(images => {
      setPreviewImages(prev => [...prev, ...images]);
    });
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
      handleFileSelect(files);
    }
  };

  const handleFileInputChange = (e) => {
    const files = e.target.files;
    if (files.length > 0) {
      handleFileSelect(files);
    }
  };

  const handleDragAreaClick = () => {
    fileInputRef.current?.click();
  };

  const removeImage = (id) => {
    setPreviewImages(prev => prev.filter(img => img.id !== id));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    if (previewImages.length === 0) {
      alert('Пожалуйста, добавьте хотя бы одно изображение товара');
      return;
    }

    if (!productData.name || !productData.price || !productData.category) {
      alert('Пожалуйста, заполните все обязательные поля');
      return;
    }

    setIsSubmitting(true);

    try {
      const fileNames = previewImages.map(img => img.file.name);
      const files = previewImages.map(img => img.file);

      console.log('📤 Создаем товар и получаем URLs для загрузки...');

      const response = await createProductAndGetUrls(productData, fileNames);
      console.log('✅ Получены URLs для загрузки:', response.urls);

      console.log('🔄 Загружаем файлы на S3...');
      await uploadFilesToS3(files, response.urls);
      console.log('✅ Все файлы загружены на S3');

      alert('Товар успешно добавлен!');
      
      setProductData({
        name: '',
        price: '',
        category: '',
        description: ''
      });
      setPreviewImages([]);

    } catch (error) {
      console.error('❌ Ошибка при создании товара:', error);
      alert(`Ошибка при создании товара: ${error.message}`);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="create-page">
      <header className="header">
        <div className='header_box'>
          <img src="/img/cart.png" className='cart' alt="Cart" />
          <Link to="/" className="home-link">
            Главная
          </Link>
        </div>
      </header>

      <div className="create-container">
        <div className='text_add'>Добавить новый товар</div>
        
        <form onSubmit={handleSubmit} className="product-form">
          {/* Drag & Drop область для изображений */}
          <div className="image-upload-section">
            <h3>Изображения товара *</h3>
            <div 
              className={`drop-zone ${isDragging ? 'dragging' : ''} ${previewImages.length > 0 ? 'has-images' : ''}`}
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
                multiple
                style={{ display: 'none' }}
              />
              
              {previewImages.length > 0 ? (
                <div className="images-preview-container">
                  <div className="images-grid">
                    {previewImages.map((image) => (
                      <div key={image.id} className="image-preview-item">
                        <img src={image.url} alt="Preview" />
                        <button 
                          type="button" 
                          className="remove-image-btn"
                          onClick={(e) => {
                            e.stopPropagation();
                            removeImage(image.id);
                          }}
                        >
                          ×
                        </button>
                      </div>
                    ))}
                    <div className="add-more-images">
                      <div className="add-more-content">
                        <div className="add-icon">+</div>
                        <span>Добавить еще</span>
                      </div>
                    </div>
                  </div>
                </div>
              ) : (
                <div className="drop-zone-content">
                  <div className="drop-icon">📁</div>
                  <p>Перетащите изображения сюда или кликните для выбора</p>
                  <span>PNG, JPG, JPEG (макс. 5MB каждое)</span>
                  <span className="multiple-hint">Можно выбрать несколько файлов</span>
                </div>
              )}
            </div>
            {previewImages.length > 0 && (
              <div className="images-counter">
                Добавлено изображений: {previewImages.length}
              </div>
            )}
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
  );
};

export default Create;
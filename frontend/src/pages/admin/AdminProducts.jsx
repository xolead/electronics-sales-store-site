import React, { useState, useEffect, useRef } from 'react';
import { Link } from 'react-router-dom';
import axios from 'axios';
import './AdminProducts.css';
import { getCategoryFromParameters } from '../../utils/parameters';
import {
    addParameter,
    updateParameter,
    removeParameter,
    validateParameters,
    prepareParametersForSubmit
} from '../../utils/parameters';

const AdminProducts = () => {
  const [products, setProducts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [editingProduct, setEditingProduct] = useState(null);
  const [showDeleteModal, setShowDeleteModal] = useState(false);
  const [productToDelete, setProductToDelete] = useState(null);
  const [searchTerm, setSearchTerm] = useState('');
  const [saveLoading, setSaveLoading] = useState(false);
  const [deleteLoading, setDeleteLoading] = useState(false);
  const [previewImages, setPreviewImages] = useState([]);
  const [isDragging, setIsDragging] = useState(false);
  const [parameters, setParameters] = useState([]);
  const fileInputRef = useRef(null);

  useEffect(() => {
    loadProducts();
  }, []);

  const loadProducts = async () => {
    try {
      setLoading(true);
      setError(null);
      const response = await axios.get('/product');
      console.log('📦 Получены товары:', response.data);
      
      let productsData = [];
      if (response.data && response.data.Products) {
        productsData = response.data.Products;
      } else if (response.data && Array.isArray(response.data)) {
        productsData = response.data;
      }
      
      setProducts(productsData);
    } catch (err) {
      console.error('❌ Ошибка загрузки товаров:', err);
      setError('Не удалось загрузить товары');
    } finally {
      setLoading(false);
    }
  };

  const handleEdit = (product) => {
    // Загружаем существующие изображения
    const existingImages = product.images ? product.images.map((image, index) => ({
      id: `existing-${index}`,
      url: `https://electronic.s3.regru.cloud/products/${image}`,
      fileName: image,
      isExisting: true
    })) : [];
    
    setPreviewImages(existingImages);
    
    // Парсим параметры
    const parsedParameters = product.parameters ? 
      product.parameters.split('|').map(param => {
        const [key, value] = param.split('=');
        return { id: Date.now() + Math.random(), key: key || '', value: value || '' };
      }) : [];
    
    setParameters(parsedParameters);
    setEditingProduct({ ...product });
  };

  const handleCancelEdit = () => {
    setEditingProduct(null);
    setPreviewImages([]);
    setParameters([]);
  };

  // Функции для управления изображениями
  const handleFileSelect = (files) => {
    const imageFiles = Array.from(files).filter((file) =>
      file.type.startsWith("image/"),
    );

    if (imageFiles.length === 0) return;

    const newImages = imageFiles.map((file) => {
      return new Promise((resolve) => {
        const reader = new FileReader();
        reader.onload = (e) => {
          resolve({
            id: Date.now() + Math.random(),
            url: e.target.result,
            file: file,
            fileName: file.name,
            isNew: true
          });
        };
        reader.readAsDataURL(file);
      });
    });

    Promise.all(newImages).then((images) => {
      setPreviewImages((prev) => [...prev, ...images]);
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
    setPreviewImages((prev) => prev.filter((img) => img.id !== id));
  };

  // Функции для управления параметрами
  const handleAddParameter = () => {
    addParameter(parameters, setParameters);
  };

  const handleUpdateParameter = (id, field, value) => {
    updateParameter(parameters, setParameters, id, field, value);
  };

  const handleRemoveParameter = (id) => {
    removeParameter(parameters, setParameters, id);
  };

  // Функция загрузки новых файлов на S3
  const uploadFilesToS3 = async (files, s3Urls) => {
    for (let i = 0; i < files.length; i++) {
      try {
        const file = files[i];
        const s3Url = s3Urls[i];

        await axios.put(s3Url, file, {
          headers: {
            "Content-Type": file.type,
            "x-amz-acl": "public-read",
          },
        });

        console.log(`✅ Файл "${file.name}" загружен на S3`);
      } catch (error) {
        console.error(`Ошибка при загрузке файла ${files[i].name}:`, error);
        throw new Error(`Не удалось загрузить файл: ${files[i].name}`);
      }
    }
  };

  const handleSave = async () => {
    if (!editingProduct) return;

    try {
      setSaveLoading(true);
      
      // Проверяем параметры
      const validation = validateParameters(parameters);
      if (!validation.isValid) {
        alert(`Пожалуйста, заполните все добавленные параметры (${validation.incompleteCount} не заполнено)`);
        return;
      }

      // Подготавливаем данные для отправки
      const parametersString = prepareParametersForSubmit(parameters, editingProduct.category);
      
      // Разделяем изображения на существующие и новые
      const existingImages = previewImages.filter(img => img.isExisting).map(img => img.fileName);
      const newImages = previewImages.filter(img => img.isNew);
      const newImageFiles = newImages.map(img => img.file);
      const newImageNames = newImages.map(img => img.fileName);

      // Обновляем товар с существующими изображениями
      const productData = {
        name: editingProduct.name,
        price: Number(editingProduct.price),
        category: editingProduct.category,
        description: editingProduct.description,
        count: Number(editingProduct.count),
        parameters: parametersString,
        images: [...existingImages, ...newImageNames] // Объединяем старые и новые имена файлов
      };

      // Сначала обновляем товар
      await axios.put(`/product/${editingProduct.id}`, productData);
      
      // Если есть новые изображения, загружаем их на S3
      if (newImages.length > 0) {
        console.log("🔄 Загружаем новые файлы на S3...");
        
        // Получаем URLs для загрузки новых файлов
        const uploadResponse = await axios.post('/product/upload-urls', {
          fileNames: newImageNames
        });
        
        await uploadFilesToS3(newImageFiles, uploadResponse.data.urls);
        console.log("✅ Все новые файлы загружены на S3");
      }
      
      // Обновляем локальное состояние
      setProducts(prev => prev.map(p => 
        p.id === editingProduct.id ? { ...p, ...productData } : p
      ));
      
      setEditingProduct(null);
      setPreviewImages([]);
      setParameters([]);
      alert('Товар успешно обновлен!');
    } catch (err) {
      console.error('❌ Ошибка обновления товара:', err);
      alert('Ошибка при обновлении товара');
    } finally {
      setSaveLoading(false);
    }
  };

  const handleDeleteClick = (product) => {
    setProductToDelete(product);
    setShowDeleteModal(true);
  };

  const handleDeleteConfirm = async () => {
    if (!productToDelete) return;

    try {
      setDeleteLoading(true);
      await axios.delete(`/product/${productToDelete.id}`);
      
      // Удаляем из локального состояния
      setProducts(prev => prev.filter(p => p.id !== productToDelete.id));
      
      setShowDeleteModal(false);
      setProductToDelete(null);
      alert('Товар успешно удален!');
    } catch (err) {
      console.error('❌ Ошибка удаления товара:', err);
      alert('Ошибка при удалении товара');
    } finally {
      setDeleteLoading(false);
    }
  };

  const handleDeleteCancel = () => {
    setShowDeleteModal(false);
    setProductToDelete(null);
  };

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setEditingProduct(prev => ({
      ...prev,
      [name]: value
    }));
  };

  // Фильтрация товаров по поиску
  const filteredProducts = products.filter(product =>
    product.name?.toLowerCase().includes(searchTerm.toLowerCase()) ||
    product.category?.toLowerCase().includes(searchTerm.toLowerCase())
  );

  if (loading) {
    return (
      <div className="admin-products">
        <div className="admin-section">
          <div className="admin-loading-container">
            <div className="admin-loading-spinner"></div>
            <p>Загрузка товаров...</p>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="admin-products">
      <div className="admin-section">
        <div className="admin-products-header">
          <h2>Управление товарами</h2>
          <Link to="/admin" className='admin-add-product-btn' style={{marginRight: '130px'}} >
            Вернуться
          </Link>

          <Link to="/admin/create" className="admin-add-product-btn">
            + Добавить товар
          </Link>
        </div>

        {error && (
          <div className="admin-error-message">
            <p>{error}</p>
            <button onClick={loadProducts} className="admin-retry-btn">
              Попробовать снова
            </button>
          </div>
        )}

        <div className="admin-products-controls">
          <div className="admin-search-box">
            <input
              type="text"
              placeholder="Поиск по названию или категории..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="admin-search-input"
            />
          </div>
          <div className="admin-products-stats">
            Всего товаров: {products.length}
          </div>
        </div>

        <div className="admin-products-table-container">
          <table className="admin-products-table">
            <thead>
              <tr>
                <th>ID</th>
                <th>Изображение</th>
                <th>Название</th>
                <th>Категория</th>
                <th>Цена</th>
                <th>Количество</th>
                <th>Действия</th>
              </tr>
            </thead>
            <tbody>
              {filteredProducts.length === 0 ? (
                <tr>
                  <td colSpan="7" className="admin-no-products">
                    {searchTerm ? 'Товары не найдены' : 'Нет товаров'}
                  </td>
                </tr>
              ) : (
                filteredProducts.map(product => (
                  <tr key={product.id}>
                    <td className="admin-product-id">{product.id}</td>
                    <td className="admin-product-image">
                      {product.images && product.images.length > 0 ? (
                        <img 
                          src={`https://electronic.s3.regru.cloud/products/${product.images[0]}`}
                          alt={product.name}
                          onError={(e) => {
                            e.target.src = '/img/placeholder.jpg';
                          }}
                        />
                      ) : (
                        <div className="admin-no-image">Нет фото</div>
                      )}
                    </td>
                    <td className="admin-product-name">{product.name}</td>
                    <td className="admin-product-category">{getCategoryFromParameters(product.parameters)}</td>
                    <td className="admin-product-price">{product.price?.toLocaleString()} ₽</td>
                    <td className="admin-product-count">{product.count} шт.</td>
                    <td className="admin-product-actions">
                      <button
                        onClick={() => handleEdit(product)}
                        className="admin-edit-btn"
                        title="Редактировать"
                      >
                        ✏️
                      </button>
                      <button
                        onClick={() => handleDeleteClick(product)}
                        className="admin-delete-btn"
                        title="Удалить"
                      >
                        🗑️
                      </button>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </div>

      {/* Модальное окно редактирования */}
      {editingProduct && (
        <div className="admin-modal-overlay">
          <div className="admin-modal-content">
            <h3>Редактирование товара</h3>
            
            {/* Секция изображений */}
            <div className="image-upload-section">
              <h3>Изображения товара</h3>
              <div
                className={`drop-zone ${isDragging ? "dragging" : ""} ${previewImages.length > 0 ? "has-images" : ""}`}
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
                  style={{ display: "none" }}
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
                    <span className="multiple-hint">
                      Можно выбрать несколько файлов
                    </span>
                  </div>
                )}
              </div>
              {previewImages.length > 0 && (
                <div className="images-counter">
                  Изображений: {previewImages.length}
                </div>
              )}
            </div>

            <div className="form-group">
              <label>Название товара *</label>
              <input
                type="text"
                name="name"
                value={editingProduct.name || ''}
                onChange={handleInputChange}
                required
              />
            </div>

            <div className="form-row">
              <div className="form-group">
                <label>Цена (₽) *</label>
                <input
                  type="number"
                  name="price"
                  value={editingProduct.price || ''}
                  onChange={handleInputChange}
                  required
                  min="0"
                />
              </div>

              <div className="form-group">
                <label>Количество *</label>
                <input
                  type="number"
                  name="count"
                  value={editingProduct.count || ''}
                  onChange={handleInputChange}
                  required
                  min="0"
                />
              </div>
            </div>

            <div className="form-group">
              <label>Категория *</label>
              <select
                name="category"
                value={editingProduct.category || ''}
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
              <label>Описание</label>
              <textarea
                name="description"
                value={editingProduct.description || ''}
                onChange={handleInputChange}
                rows="4"
                placeholder="Описание товара..."
              />
            </div>

            {/* Секция характеристик */}
            <div className="parameters-section">
              <div className="parameters-header">
                <h4>Дополнительные характеристики</h4>
                <button 
                  type="button" 
                  className="add-parameter-btn"
                  onClick={handleAddParameter}
                >
                  + Добавить характеристику
                </button>
              </div>
              
              <div className="parameters-info">
                <span>Формат: "Название: Значение" (например: Цвет: Черный)</span>
              </div>
              
              {parameters.map((param, index) => (
                <div key={param.id} className="parameter-row">
                  <input
                    type="text"
                    placeholder="Название"
                    value={param.key}
                    onChange={(e) => handleUpdateParameter(param.id, 'key', e.target.value)}
                    className="parameter-key"
                  />
                  <span className="parameter-equals"></span>
                  <input
                    type="text"
                    placeholder="Значение"
                    value={param.value}
                    onChange={(e) => handleUpdateParameter(param.id, 'value', e.target.value)}
                    className="parameter-value"
                  />
                  <button
                    type="button"
                    className="remove-parameter-btn"
                    onClick={() => handleRemoveParameter(param.id)}
                  >
                    ×
                  </button>
                </div>
              ))}
              
              {parameters.length === 0 && (
                <div className="no-parameters">
                  <p>Пока не добавлено ни одной характеристики</p>
                  <span>Например: Цвет: Черный, Память: 128ГБ, Материал: Алюминий</span>
                </div>
              )}
            </div>

            <div className="admin-modal-actions">
              <button
                onClick={handleCancelEdit}
                className="admin-cancel-btn"
                disabled={saveLoading}
              >
                Отмена
              </button>
              <button
                onClick={handleSave}
                className="admin-save-btn"
                disabled={saveLoading}
              >
                {saveLoading ? 'Сохранение...' : 'Сохранить'}
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Модальное окно удаления */}
      {showDeleteModal && productToDelete && (
        <div className="admin-modal-overlay">
          <div className="admin-modal-content">
            <h3>Удаление товара</h3>
            <p>Вы уверены, что хотите удалить товар "<strong>{productToDelete.name}</strong>"?</p>
            <p>Это действие нельзя отменить.</p>
            
            <div className="admin-modal-actions">
              <button
                onClick={handleDeleteCancel}
                className="admin-cancel-btn"
                disabled={deleteLoading}
              >
                Отмена
              </button>
              <button
                onClick={handleDeleteConfirm}
                className="admin-delete-confirm-btn"
                disabled={deleteLoading}
              >
                {deleteLoading ? 'Удаление...' : 'Удалить'}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default AdminProducts;
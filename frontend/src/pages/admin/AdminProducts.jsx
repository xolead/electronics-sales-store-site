import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import axios from 'axios';
import './AdminProducts.css';
import { getCategoryFromParameters } from '../../utils/formattingCategory';
import Create from './Create';

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
    setEditingProduct({ ...product });
  };

  const handleCancelEdit = () => {
    setEditingProduct(null);
  };

  const handleSave = async () => {
    if (!editingProduct) return;

    try {
      setSaveLoading(true);
      
      // Подготавливаем данные для отправки
      const productData = {
        name: editingProduct.name,
        price: Number(editingProduct.price),
        category: editingProduct.category,
        description: editingProduct.description,
        count: Number(editingProduct.count),
        parameters: editingProduct.parameters
      };

      await axios.put(`/product/${editingProduct.id}`, productData);
      
      // Обновляем локальное состояние
      setProducts(prev => prev.map(p => 
        p.id === editingProduct.id ? { ...p, ...productData } : p
      ));
      
      setEditingProduct(null);
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
          <div className="loading-container">
            <div className="loading-spinner"></div>
            <p>Загрузка товаров...</p>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="admin-products">
      <div className="admin-section">
        <div className="products-header">
          <h2>Управление товарами</h2>
          <Link to="/admin/create" className="add-product-btn">
            + Добавить товар
          </Link>
        </div>

        {error && (
          <div className="error-message">
            <p>{error}</p>
            <button onClick={loadProducts} className="retry-btn">
              Попробовать снова
            </button>
          </div>
        )}

        <div className="products-controls">
          <div className="search-box">
            <input
              type="text"
              placeholder="Поиск по названию или категории..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="search-input"
            />
          </div>
          <div className="products-stats">
            Всего товаров: {products.length}
          </div>
        </div>

        <div className="products-table-container">
          <table className="products-table">
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
                  <td colSpan="7" className="no-products">
                    {searchTerm ? 'Товары не найдены' : 'Нет товаров'}
                  </td>
                </tr>
              ) : (
                filteredProducts.map(product => (
                  <tr key={product.id}>
                    <td className="product-id">{product.id}</td>
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
                        <div className="no-image">Нет фото</div>
                      )}
                    </td>
                    <td className="product-name">{product.name}</td>
                    <td className="admin-product-category">{getCategoryFromParameters(product.parameters)}</td>
                    <td className="product-price">{product.price?.toLocaleString()} ₽</td>
                    <td className="product-count">{product.count} шт.</td>
                    <td className="product-actions">
                      <button
                        onClick={() => handleEdit(product)}
                        className="edit-btn"
                        title="Редактировать"
                      >
                        ✏️
                      </button>
                      <button
                        onClick={() => handleDeleteClick(product)}
                        className="delete-btn"
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
        <div className="modal-overlay">
          <div className="modal-content">
            <h3>Редактирование товара</h3>
            
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

            <div className="form-group">
              <label>Характеристики</label>
              <textarea
                name="parameters"
                value={editingProduct.parameters || ''}
                onChange={handleInputChange}
                rows="3"
                placeholder="Формат: ключ=значение|ключ=значение"
              />
              <small>Формат: Цвет=Черный|Память=128ГБ</small>
            </div>

            <div className="modal-actions">
              <button
                onClick={handleCancelEdit}
                className="cancel-btn"
                disabled={saveLoading}
              >
                Отмена
              </button>
              <button
                onClick={handleSave}
                className="save-btn"
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
        <div className="modal-overlay">
          <div className="modal-content">
            <h3>Удаление товара</h3>
            <p>Вы уверены, что хотите удалить товар "<strong>{productToDelete.name}</strong>"?</p>
            <p>Это действие нельзя отменить.</p>
            
            <div className="modal-actions">
              <button
                onClick={handleDeleteCancel}
                className="cancel-btn"
                disabled={deleteLoading}
              >
                Отмена
              </button>
              <button
                onClick={handleDeleteConfirm}
                className="delete-confirm-btn"
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
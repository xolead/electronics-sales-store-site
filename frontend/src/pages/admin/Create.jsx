import React, { useState, useRef } from "react";
import { Link } from "react-router-dom";
import "./Create.css";
import axios from "axios";
import {
    addParameter,
    updateParameter,
    removeParameter,
    parseParameters,
    prepareParametersForSubmit,
    validateParameters
} from '../../utils/parameters'

const Create = () => {
  const [productData, setProductData] = useState({
    name: "",
    price: "",
    category: "",
    description: "",
    count: "1",
  });
  const [previewImages, setPreviewImages] = useState([]);
  const [isDragging, setIsDragging] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [parameters, setParameters] = useState([]);
  const fileInputRef = useRef(null);

  // Функция создания продукта
  const createProductAndGetUrls = async (productData, fileNames) => {
    try {
      // Формируем строку параметров с помощью новой функции
      const parametersString = prepareParametersForSubmit(parameters, productData.category);

      const response = await axios.post(
        "/product",
        {
          name: productData.name,
          price: Number(productData.price),
          description: productData.description,
          parameters: parametersString,
          count: Number(productData.count) || 1,
          images: fileNames,
        },
        {
          headers: {
            "Content-Type": "application/json",
          },
        },
      );

      return response.data;
    } catch (error) {
      console.error("Ошибка при создании товара:", error);
      throw error;
    }
  };

  // Функция загрузки файлов на S3
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

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setProductData((prev) => ({
      ...prev,
      [name]: value,
    }));
  };

  // Обработчики для параметров
  const handleAddParameter = () => {
    addParameter(parameters, setParameters);
  };

  const handleUpdateParameter = (id, field, value) => {
    updateParameter(parameters, setParameters, id, field, value);
  };

  const handleRemoveParameter = (id) => {
    removeParameter(parameters, setParameters, id);
  };

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

  const handleSubmit = async (e) => {
    e.preventDefault();

    if (previewImages.length === 0) {
      alert("Пожалуйста, добавьте хотя бы одно изображение товара");
      return;
    }

    if (!productData.name || !productData.price || !productData.category) {
      alert("Пожалуйста, заполните все обязательные поля");
      return;
    }

    // Проверяем параметры с помощью новой функции
    const validation = validateParameters(parameters);
    if (!validation.isValid) {
      alert(`Пожалуйста, заполните все добавленные параметры (${validation.incompleteCount} не заполнено)`);
      return;
    }

    setIsSubmitting(true);

    try {
      const fileNames = previewImages.map((img) => img.file.name);
      const files = previewImages.map((img) => img.file);

      console.log("📤 Создаем товар и получаем URLs для загрузки...");

      const response = await createProductAndGetUrls(productData, fileNames);
      console.log("✅ Получены URLs для загрузки:", response.urls);

      console.log("🔄 Загружаем файлы на S3...");
      await uploadFilesToS3(files, response.urls);
      console.log("✅ Все файлы загружены на S3");

      alert("Товар успешно добавлен!");

      // Сбрасываем форму
      setProductData({
        name: "",
        price: "",
        category: "",
        description: "",
        count: "1",
      });
      setPreviewImages([]);
      setParameters([]);
    } catch (error) {
      console.error("❌ Ошибка при создании товара:", error);
      alert(`Ошибка при создании товара: ${error.message}`);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="create-page">
      <header className="header">
        <div className='header_box'>
          <Link to="/admin/products" className="home-link">
            Вернуться
          </Link>
        </div>
      </header>

      <div className="create-container">
        <div className="text_add">Добавить новый товар</div>

        <form onSubmit={handleSubmit} className="product-form">
          {/* Drag & Drop область для изображений */}
          <div className="image-upload-section">
            <h3>Изображения товара *</h3>
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
                Добавлено изображений: {previewImages.length}
              </div>
            )}
          </div>

          {/* Общая информация о товаре */}
          <div className="product-info-section">
            <h3>Информация о товаре</h3>

            <div className="form-row">
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
            </div>

            <div className="form-row">
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
                <label htmlFor="count">Количество товара *</label>
                <input
                  type="number"
                  id="count"
                  name="count"
                  value={productData.count}
                  onChange={handleInputChange}
                  required
                  placeholder="1"
                  min="1"
                />
              </div>
            </div>

            {/* Динамические параметры */}
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

            {/* Описание товара */}
            <div className="form-group">
              <label htmlFor="description">Описание товара</label>
              <textarea
                id="description"
                name="description"
                value={productData.description}
                onChange={handleInputChange}
                rows="6"
                placeholder="Подробное описание товара, характеристики, преимущества, особенности использования..."
              />
            </div>
          </div>

          <div className="form-actions">
            <Link to="/" className="cancel-btn">
              Отмена
            </Link>
            <button 
              type="submit" 
              className="submit-btn"
              disabled={isSubmitting}
            >
              {isSubmitting ? "Добавление..." : "Добавить товар"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default Create;
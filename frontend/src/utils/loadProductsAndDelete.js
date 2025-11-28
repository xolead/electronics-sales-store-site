import axios from 'axios';
import { getAll } from '../services/api';

// Базовый URL S3
export const getFullImageUrl = (filename) => {
  const url = `https://electronic.s3.regru.cloud/products/${filename}`;
  return url;
};

// Функция загрузки товаров
export const loadProducts = async () => {    
  const productsData = await getAll();
  return {
    success: true,
    data: productsData,
    error: null
  };
};

// Функция удаления товара
export const deleteProduct = async (id) => {
  try {
    await axios.delete(`/product/${id}`);
    return { success: true, error: null };
  } catch (error) {
    console.error('Ошибка удаления товара:', error);
    return { success: false, error: 'Не удалось удалить товар' };
  }
};
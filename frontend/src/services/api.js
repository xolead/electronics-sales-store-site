import axios from 'axios';

export const api = axios.create({
});


export const getAll = async () => {
    try {
      console.log('🔄 Запрашиваем товары...');
      const response = await axios.get('/product');
      console.log('📦 Полный ответ:', response.data);
      
      if (response.data && response.data.Products) {
        console.log('✅ Товары найдены:', response.data.Products);
        return response.data.Products;
      } else if (response.data && Array.isArray(response.data)) {
        console.log('✅ Товары (массив):', response.data);
        return response.data;
      } else {
        console.warn('⚠️ Товары не найдены в ответе');
        return [];
      }
      
    } catch (error) {
      console.error('❌ Ошибка при получении товаров:', error);
      return [];
    }
  };
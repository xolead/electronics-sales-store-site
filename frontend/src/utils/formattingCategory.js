  // Функция для получения категории из параметров
  export const getCategoryFromParameters = (parametersString) => {
    if (!parametersString) return '';
    
    try {
      const pairs = parametersString.split('|');
      
      for (let pair of pairs) {
        const [key, value] = pair.split('=');
        if (key && value && key.trim() === 'Категория') {
          return value.trim();
        }
      }
      
      return '';
    } catch (error) {
      console.error('Ошибка парсинга категории:', error);
      return '';
    }
  };
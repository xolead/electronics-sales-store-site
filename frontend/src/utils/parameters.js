
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
 
 export const addParameter = (parameters, setParameters) => {
    setParameters([...parameters, { key: "", value: "", id: Date.now() }]);
  };
  
  // Функция обновления параметра
  export const updateParameter = (parameters, setParameters, id, field, value) => {
    setParameters(parameters.map(param => 
      param.id === id ? { ...param, [field]: value } : param
    ));
  };
  
  // Функция удаления параметра
  export const removeParameter = (parameters, setParameters, id) => {
    setParameters(parameters.filter(param => param.id !== id));
  };
  
  // Функция парсинга параметров из строки в массив объектов
  export const parseParameters = (parametersString) => {
    if (!parametersString) return [];
    
    const parameters = [];
    
    try {
      const pairs = parametersString.split('|');
      
      pairs.forEach(pair => {
        const trimmedPair = pair.trim();
        const equalsIndex = trimmedPair.indexOf('=');
        
        if (equalsIndex > 0) {
          const key = trimmedPair.substring(0, equalsIndex).trim();
          const value = trimmedPair.substring(equalsIndex + 1).trim();
          
          if (key && value) {
            parameters.push({ key, value, id: Date.now() + Math.random() });
          }
        }
      });
      
      return parameters;
    } catch (error) {
      console.error('Ошибка парсинга параметров:', error);
      return [];
    }
  };
  
  // Функция преобразования массива параметров обратно в строку
  export const formatParametersToString = (parametersArray) => {
    return parametersArray
      .filter(param => param.key && param.value)
      .map(param => `${param.key}=${param.value}`)
      .join('|');
  };
  
  // Функция подготовки параметров для отправки на сервер
  export const prepareParametersForSubmit = (parameters, category) => {
    const parametersArray = [];
    
    // Добавляем категорию как первый параметр
    if (category) {
      parametersArray.push(`Категория=${category}`);
    }
    
    // Добавляем остальные параметры
    parameters.forEach(param => {
      if (param.key && param.value) {
        parametersArray.push(`${param.key}=${param.value}`);
      }
    });
    
    // Объединяем все параметры через |
    return parametersArray.join('|');
  };
  
  // Функция проверки заполненности параметров
  export const validateParameters = (parameters) => {
    const incompleteParameters = parameters.filter(param => !param.key || !param.value);
    return {
      isValid: incompleteParameters.length === 0,
      incompleteCount: incompleteParameters.length
    };
  };
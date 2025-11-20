import React from 'react';
import { render, screen } from '@testing-library/react';
import App from './App';

// Простой тест без моков роутинга
test('рендерит главную страницу', () => {
  render(<App />);
  // Проверяем что приложение рендерится без ошибок
  expect(screen.getByText(/Добавить/i)).toBeInTheDocument();
});
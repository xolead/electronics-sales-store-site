import React from 'react';
import { render, screen } from '@testing-library/react';
import Create from './Create';

test('рендерит форму создания товара', () => {
  render(<Create />);
  expect(screen.getByText('Добавить новый товар')).toBeInTheDocument();
});
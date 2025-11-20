// src/setupTests.js
import '@testing-library/jest-dom';

// Optional: мокаем console.error чтобы тесты не засорялись предупреждениями
const originalError = console.error;
beforeAll(() => {
  console.error = (...args) => {
    if (/Warning.*not wrapped in act/.test(args[0])) {
      return;
    }
    originalError.call(console, ...args);
  };
});

afterAll(() => {
  console.error = originalError;
});
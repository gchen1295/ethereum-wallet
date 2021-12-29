/** @type {import('ts-jest/dist/types').InitialOptionsTsJest} */
module.exports = {
  preset: 'ts-jest',
  testEnvironment: 'node',
  transform: {
    '^.+\\.(ts|tsx)?$': 'ts-jest',
     // Use babel-jest to transpile tests with the next/babel preset
    // https://jestjs.io/docs/configuration#transform-objectstring-pathtotransformer--pathtotransformer-object
    '^.+\\.(js|jsx|ts|tsx)$': ['babel-jest', { presets: ['next/babel'] }],
  },
  transformIgnorePatterns: [
    '/node_modules/',
    '^.+\\.module\\.(css|sass|scss)$',
  ],
  moduleNameMapper: {
    '@components/(.*)': '<rootDir>/src/components/$1',
    '@assets/(.*)': '<rootDir>/src/assets/$1',
    '@store/(.*)': '<rootDir>/src/store/$1',
    '@services/(.*)': '<rootDir>/src/services/$1',
    '@api/(.*)': '<rootDir>/src/pages/api/$1',
    '@libs/(.*)': '<rootDir>/src/libs/$1',
    '@styles/(.*)': '<rootDir>/src/styles/$1',
    '@interfaces/(.*)': '<rootDir>/src/interfaces/$1',
  },
  roots: [
    '<rootDir>/src'
  ],
  setupFiles: ["./jest.setup.js"],
  testPathIgnorePatterns: ['<rootDir>/node_modules/', '<rootDir>/.next/'],
  collectCoverageFrom: [
    '**/*.{js,jsx,ts,tsx}',
    '!**/*.d.ts',
    '!**/node_modules/**',
  ],
};
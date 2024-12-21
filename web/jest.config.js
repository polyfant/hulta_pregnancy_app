export default {
  testEnvironment: 'jsdom',
  transform: {
    '^.+\\.js$': 'babel-jest',
  },
  moduleNameMapper: {
    '^@/(.*)$': '<rootDir>/static/js/$1',
  },
  setupFiles: ['fake-indexeddb/auto'],
};

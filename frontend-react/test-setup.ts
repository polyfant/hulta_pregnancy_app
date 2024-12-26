import { afterEach } from 'bun:test';
import { cleanup } from '@testing-library/react';
import '@testing-library/jest-dom';
import { JSDOM } from 'jsdom';

const dom = new JSDOM('<!DOCTYPE html><html><body></body></html>', {
    url: 'http://localhost',
    pretendToBeVisual: true,
    runScripts: 'dangerously',
});

// Set up a mock window environment
const window = dom.window;

// Copy properties from window to global
Object.defineProperty(global, 'window', {
    value: window,
    writable: true,
});

Object.defineProperty(global, 'document', {
    value: window.document,
    writable: true,
});

Object.defineProperty(global, 'navigator', {
    value: window.navigator,
    writable: true,
});

// Add other necessary DOM properties
global.HTMLElement = window.HTMLElement;
global.Element = window.Element;
global.Node = window.Node;
global.DocumentFragment = window.DocumentFragment;

// Mock window methods
global.getComputedStyle = window.getComputedStyle.bind(window);
global.requestAnimationFrame = (callback) => setTimeout(callback, 0);
global.cancelAnimationFrame = (id) => clearTimeout(id);

// Add event listener methods to the window
global.addEventListener = window.addEventListener.bind(window);
global.removeEventListener = window.removeEventListener.bind(window);
global.dispatchEvent = window.dispatchEvent.bind(window);

// Clean up after each test
afterEach(() => {
    cleanup();
    if (document.body) {
        document.body.innerHTML = '';
    }
});

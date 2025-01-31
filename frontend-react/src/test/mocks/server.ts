import { setupServer } from 'msw/node';
import { handlers } from '../setup';

export const server = setupServer(...handlers); 
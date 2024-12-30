import dayjs from 'dayjs';

/**
 * Format a date string or Date object to a human-readable format
 * @param date Date string or Date object
 * @returns Formatted date string (e.g., "2024-12-30")
 */
export const formatDate = (date: string | Date): string => {
  return dayjs(date).format('YYYY-MM-DD');
};

/**
 * Check if a date is in the future
 * @param date Date to check
 * @returns boolean
 */
export const isFutureDate = (date: string | Date): boolean => {
  return dayjs(date).isAfter(dayjs());
};

/**
 * Calculate the difference in days between two dates
 * @param date1 First date
 * @param date2 Second date
 * @returns number of days
 */
export const daysBetween = (date1: string | Date, date2: string | Date): number => {
  return dayjs(date2).diff(dayjs(date1), 'day');
};

/**
 * Add days to a date
 * @param date Starting date
 * @param days Number of days to add
 * @returns New date string
 */
export const addDays = (date: string | Date, days: number): string => {
  return dayjs(date).add(days, 'day').format('YYYY-MM-DD');
};

/**
 * Get today's date in YYYY-MM-DD format
 * @returns Today's date string
 */
export const getTodayFormatted = (): string => {
  return dayjs().format('YYYY-MM-DD');
};

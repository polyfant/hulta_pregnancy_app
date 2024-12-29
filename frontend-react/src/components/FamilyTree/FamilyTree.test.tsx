import { test, expect, beforeEach, describe, jest, afterEach } from 'bun:test';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { FamilyTree } from './FamilyTree';
import { MantineProvider } from '@mantine/core';
import React from 'react';

const mockFamilyData = {
  horse: {
    id: 1,
    name: "Main Horse",
    breed: "Arabian",
    gender: "MARE",
    dateOfBirth: "2020-01-01",
    age: "3 years",
  },
  mother: {
    id: 2,
    name: "Mother Horse",
    breed: "Arabian",
    gender: "MARE",
    dateOfBirth: "2015-01-01",
    age: "8 years",
    isExternal: false,
  },
  father: {
    id: 3,
    name: "Father Horse",
    breed: "Arabian",
    gender: "STALLION",
    dateOfBirth: "2014-01-01",
    age: "9 years",
    isExternal: false,
  },
  siblings: [
    {
      id: 4,
      name: "Sibling Horse",
      breed: "Arabian",
      gender: "STALLION",
      dateOfBirth: "2021-01-01",
      age: "2 years",
      isExternal: false,
    },
  ],
  offspring: [
    {
      id: 5,
      name: "Offspring Horse",
      breed: "Arabian",
      gender: "MARE",
      dateOfBirth: "2023-01-01",
      age: "0 years",
      isExternal: false,
    },
  ],
};

// Mock fetch
global.fetch = (() =>
  Promise.resolve({
    ok: true,
    json: () => Promise.resolve(mockFamilyData),
  })
) as unknown as typeof fetch;

describe('FamilyTree', () => {
  const onMemberClick = jest.fn();

  beforeEach(() => {
    onMemberClick.mockClear();
  });

  test('renders and fetches family data', async () => {
    render(
      <MantineProvider>
        <FamilyTree horseId={1} onMemberClick={onMemberClick} />
      </MantineProvider>
    );

    // Wait for data to load
    await waitFor(() => {
      expect(screen.getByText('Parents')).toBeTruthy();
    });

    // Check if parents are rendered
    expect(screen.getByText('Mother Horse')).toBeTruthy();
    expect(screen.getByText('Father Horse')).toBeTruthy();

    // Check if siblings section is rendered
    expect(screen.getByText('Siblings (1)')).toBeTruthy();
    expect(screen.getByText('Sibling Horse')).toBeTruthy();

    // Check if offspring section is rendered
    expect(screen.getByText('Offspring (1)')).toBeTruthy();
    expect(screen.getByText('Offspring Horse')).toBeTruthy();
  });

  test('handles member clicks', async () => {
    render(
      <MantineProvider>
        <FamilyTree horseId={1} onMemberClick={onMemberClick} />
      </MantineProvider>
    );

    await waitFor(() => {
      expect(screen.getByText('Mother Horse')).toBeTruthy();
    });

    // Click on a family member
    fireEvent.click(screen.getByText('Mother Horse'));
    expect(onMemberClick).toHaveBeenCalledWith(2);

    fireEvent.click(screen.getByText('Father Horse'));
    expect(onMemberClick).toHaveBeenCalledWith(3);

    fireEvent.click(screen.getByText('Sibling Horse'));
    expect(onMemberClick).toHaveBeenCalledWith(4);

    fireEvent.click(screen.getByText('Offspring Horse'));
    expect(onMemberClick).toHaveBeenCalledWith(5);
  });

  test('handles error state', async () => {
    // Mock fetch to return an error
    global.fetch = (() =>
      Promise.resolve({
        ok: false,
        status: 404,
        statusText: 'Not Found',
      })
    ) as unknown as typeof fetch;

    render(
      <MantineProvider>
        <FamilyTree horseId={1} onMemberClick={onMemberClick} />
      </MantineProvider>
    );

    await waitFor(() => {
      expect(screen.getByText(/Error loading family data/i)).toBeTruthy();
    });
  });

  test('handles section expansion', async () => {
    render(
      <MantineProvider>
        <FamilyTree horseId={1} onMemberClick={onMemberClick} />
      </MantineProvider>
    );

    await waitFor(() => {
      expect(screen.getByText('Siblings (1)')).toBeTruthy();
    });

    // Initially, siblings section should be collapsed
    expect(screen.queryByText('Sibling Horse')).not.toBeVisible();

    // Click to expand siblings section
    fireEvent.click(screen.getByText('Siblings (1)'));

    // Sibling should now be visible
    expect(screen.getByText('Sibling Horse')).toBeVisible();
  });

  test('handles external parents', async () => {
    const externalData = {
      ...mockFamilyData,
      mother: {
        name: "External Mare",
        isExternal: true,
        externalSource: "External Mare",
      },
      father: {
        name: "External Stallion",
        isExternal: true,
        externalSource: "External Stallion",
      },
    };

    global.fetch = (() =>
      Promise.resolve({
        ok: true,
        json: () => Promise.resolve(externalData),
      })
    ) as unknown as typeof fetch;

    render(
      <MantineProvider>
        <FamilyTree horseId={1} onMemberClick={onMemberClick} />
      </MantineProvider>
    );

    await waitFor(() => {
      expect(screen.getByText('External Mare (External)')).toBeTruthy();
      expect(screen.getByText('External Stallion (External)')).toBeTruthy();
    });
  });

  test('handles errors gracefully', async () => {
    const consoleError = jest.spyOn(console, 'error').mockImplementation(() => {});
    
    global.fetch = (() =>
      Promise.resolve({
        ok: false,
        status: 500,
      })
    ) as unknown as typeof fetch;

    render(
      <MantineProvider>
        <FamilyTree horseId={1} onMemberClick={onMemberClick} />
      </MantineProvider>
    );

    await waitFor(() => {
      expect(consoleError).toHaveBeenCalledWith(
        'Error fetching family tree:',
        expect.any(Error)
      );
    });

    consoleError.mockRestore();
  });
});

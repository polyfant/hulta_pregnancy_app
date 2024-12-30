import { test, expect, beforeEach, describe, jest, afterEach } from 'bun:test';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { FamilyTree } from './FamilyTree';
import { MantineProvider } from '@mantine/core';
import React from 'react';

const mockHorse = {
  id: 1,
  name: "Main Horse",
  breed: "Arabian",
  gender: "MARE",
  dateOfBirth: "2020-01-01",
  motherId: 2,
  fatherId: 3,
  mother: {
    id: 2,
    name: "Mother Horse",
    breed: "Arabian",
    gender: "MARE",
    dateOfBirth: "2015-01-01",
  },
  father: {
    id: 3,
    name: "Father Horse",
    breed: "Arabian",
    gender: "STALLION",
    dateOfBirth: "2014-01-01",
  }
};

describe('FamilyTree', () => {
  test('renders horse with parents', () => {
    render(
      <MantineProvider>
        <FamilyTree horse={mockHorse} />
      </MantineProvider>
    );

    expect(screen.getByText('Main Horse')).toBeTruthy();
    expect(screen.getByText('Mother Horse')).toBeTruthy();
    expect(screen.getByText('Father Horse')).toBeTruthy();
  });

  test('renders external parents correctly', () => {
    const horseWithExternalParents = {
      ...mockHorse,
      motherId: undefined,
      fatherId: undefined,
      mother: undefined,
      father: undefined,
      externalMother: "External Mare",
      externalFather: "External Stallion"
    };

    render(
      <MantineProvider>
        <FamilyTree horse={horseWithExternalParents} />
      </MantineProvider>
    );

    expect(screen.getByText('External Mare')).toBeTruthy();
    expect(screen.getByText('External Stallion')).toBeTruthy();
  });

  test('handles missing parents gracefully', () => {
    const horseWithoutParents = {
      ...mockHorse,
      motherId: undefined,
      fatherId: undefined,
      mother: undefined,
      father: undefined
    };

    render(
      <MantineProvider>
        <FamilyTree horse={horseWithoutParents} />
      </MantineProvider>
    );

    expect(screen.getByText('Main Horse')).toBeTruthy();
    expect(screen.queryByText('Mother Horse')).toBeNull();
    expect(screen.queryByText('Father Horse')).toBeNull();
  });
});

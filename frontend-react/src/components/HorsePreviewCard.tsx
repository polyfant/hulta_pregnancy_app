import {
  Card,
  Group,
  Text,
  Badge,
  Stack
} from '@mantine/core';
import { IconMars, IconVenus } from '@tabler/icons-react';
import { Horse } from '../types/horse';

interface HorsePreviewCardProps {
  horse: Horse | null;
  externalName?: string;
  role: 'mother' | 'father';
}

export function HorsePreviewCard({ horse, externalName, role }: HorsePreviewCardProps) {
  if (!horse && !externalName) return null;

  const isExternal = !horse && externalName;
  const gender = role === 'mother' ? 'MARE' : 'STALLION';
  const color = role === 'mother' ? 'pink' : 'blue';

  return (
    <Card
      withBorder
      mt="sm"
      bg={role === 'mother' ? 'pink.0' : 'blue.0'}
    >
      <Stack gap="xs">
        <Group>
          {gender === 'MARE' ? (
            <IconVenus size="1.2rem" color="var(--mantine-color-pink-6)" />
          ) : (
            <IconMars size="1.2rem" color="var(--mantine-color-blue-6)" />
          )}
          <Text fw={500}>
            {isExternal ? externalName : horse?.name}
          </Text>
          {isExternal && (
            <Badge variant="outline" color={color}>
              External
            </Badge>
          )}
        </Group>

        {!isExternal && horse && (
          <Stack gap={4}>
            <Text size="sm" c="dimmed">
              Breed: {horse.breed}
            </Text>
            {horse.age && (
              <Text size="sm" c="dimmed">
                Age: {horse.age}
              </Text>
            )}
          </Stack>
        )}
      </Stack>
    </Card>
  );
}

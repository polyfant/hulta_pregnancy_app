import {
  Badge,
  Card,
  Group,
  Stack,
  Text
} from '@mantine/core';
import { GenderFemale, GenderMale } from '@phosphor-icons/react';
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

  const isMale = horse?.gender === 'STALLION' || horse?.gender === 'GELDING';

  return (
    <Card
      withBorder
      mt="sm"
      bg="dark.7"
    >
      <Stack gap="xs">
        <Group>
          {isMale ? (
            <GenderMale size="1.2rem" color="var(--mantine-color-blue-6)" />
          ) : (
            <GenderFemale size="1.2rem" color="var(--mantine-color-pink-6)" />
          )}
          <Text fw={500} c="white">
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

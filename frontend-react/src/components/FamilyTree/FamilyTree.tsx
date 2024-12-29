import { useState } from 'react';
import {
  Paper,
  Title,
  Text,
  Card,
  Group,
  Stack,
  Button,
  ActionIcon,
  Badge,
  Tooltip
} from '@mantine/core';
import {
  IconChevronRight,
  IconChevronDown,
  IconHorse,
  IconMars,
  IconVenus
} from '@tabler/icons-react';

interface Horse {
  id: string;
  name: string;
  gender: 'male' | 'female';
  birthDate?: string;
  color?: string;
  breed?: string;
  sire?: Horse;
  dam?: Horse;
}

interface FamilyTreeProps {
  horse: Horse;
}

interface TreeNodeProps {
  horse: Horse;
  level: number;
  maxLevel?: number;
}

function TreeNode({ horse, level, maxLevel = 3 }: TreeNodeProps) {
  const [expanded, setExpanded] = useState(level < 2);
  const hasParents = horse.sire || horse.dam;
  const canExpand = level < maxLevel && hasParents;

  const handleToggle = () => {
    if (canExpand) {
      setExpanded(!expanded);
    }
  };

  return (
    <div style={{ marginLeft: level > 0 ? 40 : 0 }}>
      <Card
        withBorder
        shadow="sm"
        padding="sm"
        radius="md"
        style={{ marginBottom: 8 }}
      >
        <Group justify="space-between">
          <Group>
            {canExpand && (
              <ActionIcon
                variant="subtle"
                onClick={handleToggle}
                aria-label={expanded ? 'Collapse' : 'Expand'}
              >
                {expanded ? <IconChevronDown size="1rem" /> : <IconChevronRight size="1rem" />}
              </ActionIcon>
            )}
            <IconHorse
              size="1.2rem"
              style={{ color: horse.gender === 'male' ? 'var(--mantine-color-blue-6)' : 'var(--mantine-color-pink-6)' }}
            />
            <Stack spacing={0}>
              <Text fw={500}>{horse.name}</Text>
              <Group spacing="xs">
                {horse.breed && (
                  <Badge size="sm" variant="light">
                    {horse.breed}
                  </Badge>
                )}
                <Badge
                  size="sm"
                  variant="light"
                  color={horse.gender === 'male' ? 'blue' : 'pink'}
                >
                  {horse.gender === 'male' ? <IconMars size="0.8rem" /> : <IconVenus size="0.8rem" />}
                </Badge>
              </Group>
            </Stack>
          </Group>
          {horse.birthDate && (
            <Tooltip label="Birth Date">
              <Text size="sm" c="dimmed">
                {new Date(horse.birthDate).getFullYear()}
              </Text>
            </Tooltip>
          )}
        </Group>
      </Card>

      {expanded && hasParents && (
        <div>
          {horse.sire && <TreeNode horse={horse.sire} level={level + 1} maxLevel={maxLevel} />}
          {horse.dam && <TreeNode horse={horse.dam} level={level + 1} maxLevel={maxLevel} />}
        </div>
      )}
    </div>
  );
}

export function FamilyTree({ horse }: FamilyTreeProps) {
  return (
    <Paper p="md">
      <Title order={2} mb="md">Family Tree</Title>
      <TreeNode horse={horse} level={0} />
    </Paper>
  );
}

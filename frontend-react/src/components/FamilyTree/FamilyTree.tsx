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
import { Horse } from '../../types/horse';

interface FamilyTreeHorse extends Horse {
  father?: FamilyTreeHorse;
  mother?: FamilyTreeHorse;
}

interface FamilyTreeProps {
  horse: FamilyTreeHorse;
}

interface TreeNodeProps {
  horse: FamilyTreeHorse;
  level: number;
  maxLevel?: number;
}

function TreeNode({ horse, level, maxLevel = 3 }: TreeNodeProps) {
  const [expanded, setExpanded] = useState(level < 2);
  const hasParents = horse.fatherId || horse.motherId || horse.externalFather || horse.externalMother;
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
              style={{ color: horse.gender === 'STALLION' ? 'var(--mantine-color-blue-6)' : 'var(--mantine-color-pink-6)' }}
            />
            <Stack gap={0}>
              <Group>
                {horse.gender === 'STALLION' && <IconMars size={16} />}
                {horse.gender === 'MARE' && <IconVenus size={16} />}
                <Text fw={500}>{horse.name}</Text>
              </Group>
              <Group gap="xs">
                {horse.breed && (
                  <Badge size="sm" variant="light">
                    {horse.breed}
                  </Badge>
                )}
                <Badge
                  size="sm"
                  variant="light"
                  color={horse.gender === 'STALLION' || horse.gender === 'GELDING' ? 'blue' : 'pink'}
                >
                  {horse.gender === 'STALLION' || horse.gender === 'GELDING' ? <IconMars size="0.8rem" /> : <IconVenus size="0.8rem" />}
                </Badge>
              </Group>
            </Stack>
          </Group>
          {horse.dateOfBirth && (
            <Tooltip label="Birth Date">
              <Text size="sm" c="dimmed">
                {new Date(horse.dateOfBirth).getFullYear()}
              </Text>
            </Tooltip>
          )}
        </Group>
      </Card>

      {expanded && (
        <div>
          {(horse.mother || horse.externalMother) && (
            <TreeNode
              horse={horse.mother || { 
                id: -1,
                name: horse.externalMother || 'Unknown Mother',
                gender: 'MARE',
                breed: 'External',
                dateOfBirth: '',
              }}
              level={level + 1}
              maxLevel={maxLevel}
            />
          )}
          {(horse.father || horse.externalFather) && (
            <TreeNode
              horse={horse.father || { 
                id: -1,
                name: horse.externalFather || 'Unknown Father',
                gender: 'STALLION',
                breed: 'External',
                dateOfBirth: '',
              }}
              level={level + 1}
              maxLevel={maxLevel}
            />
          )}
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

import React, { useState, useEffect } from 'react';
import {
  Box,
  Title,
  Accordion,
  List,
  ThemeIcon,
  Text,
  Alert,
  Loader,
  Stack,
  Badge,
  Group
} from '@mantine/core';
import {
  IconCircleCheck,
  IconAlertTriangle,
  IconList
} from '@tabler/icons-react';
import { PregnancyGuideline } from '../../types/pregnancy';

interface PregnancyGuidelinesProps {
  horseId: string;
}

const PregnancyGuidelines: React.FC<PregnancyGuidelinesProps> = ({ horseId }) => {
  const [guidelines, setGuidelines] = useState<PregnancyGuideline[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchGuidelines = async () => {
      try {
        const response = await fetch(`/api/horses/${horseId}/pregnancy/guidelines`);
        if (!response.ok) throw new Error('Failed to fetch guidelines');
        const data = await response.json();
        setGuidelines(data);
      } catch (err) {
        setError('Failed to load pregnancy guidelines');
      } finally {
        setLoading(false);
      }
    };

    fetchGuidelines();
  }, [horseId]);

  if (loading) {
    return (
      <Box style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: '200px' }}>
        <Loader />
      </Box>
    );
  }

  if (error) {
    return (
      <Alert color="red" title="Error">
        {error}
      </Alert>
    );
  }

  return (
    <Box>
      <Title order={3} mb="md">Pregnancy Guidelines</Title>

      <Accordion>
        {guidelines.map((guideline) => (
          <Accordion.Item key={guideline.stage} value={guideline.stage}>
            <Accordion.Control>
              <Group>
                <Text fw={500}>{guideline.title}</Text>
                <Badge>{guideline.stage}</Badge>
              </Group>
            </Accordion.Control>
            <Accordion.Panel>
              <Stack spacing="md">
                <Text size="sm">{guideline.description}</Text>

                {guideline.recommendations.length > 0 && (
                  <Box>
                    <Text fw={500} mb="xs">Recommendations</Text>
                    <List
                      spacing="xs"
                      size="sm"
                      icon={
                        <ThemeIcon color="teal" size={20} radius="xl">
                          <IconCircleCheck size={12} />
                        </ThemeIcon>
                      }
                    >
                      {guideline.recommendations.map((rec, index) => (
                        <List.Item key={index}>{rec}</List.Item>
                      ))}
                    </List>
                  </Box>
                )}

                {guideline.warnings.length > 0 && (
                  <Box>
                    <Text fw={500} mb="xs">Warnings</Text>
                    <List
                      spacing="xs"
                      size="sm"
                      icon={
                        <ThemeIcon color="red" size={20} radius="xl">
                          <IconAlertTriangle size={12} />
                        </ThemeIcon>
                      }
                    >
                      {guideline.warnings.map((warning, index) => (
                        <List.Item key={index}>{warning}</List.Item>
                      ))}
                    </List>
                  </Box>
                )}

                {guideline.checkpoints.length > 0 && (
                  <Box>
                    <Text fw={500} mb="xs">Checkpoints</Text>
                    <List
                      spacing="xs"
                      size="sm"
                      icon={
                        <ThemeIcon color="blue" size={20} radius="xl">
                          <IconList size={12} />
                        </ThemeIcon>
                      }
                    >
                      {guideline.checkpoints.map((checkpoint, index) => (
                        <List.Item key={index}>{checkpoint}</List.Item>
                      ))}
                    </List>
                  </Box>
                )}
              </Stack>
            </Accordion.Panel>
          </Accordion.Item>
        ))}
      </Accordion>
    </Box>
  );
};

export default PregnancyGuidelines;

import { MdCheckCircle } from 'react-icons/md';
import { FiPlus, FiSearch, FiEdit, FiTrash2 } from 'react-icons/fi';

import { Card, Title, Text, Stack, List, ThemeIcon, Paper } from '@mantine/core';


export default function PregnancyGuidelines() {
  const guidelines = {
    nutrition: [
      'Provide high-quality forage (hay or pasture)',
      'Feed grain according to body condition',
      'Ensure access to fresh, clean water',
      'Provide salt and mineral supplements',
      'Monitor body condition score regularly'
    ],
    exercise: [
      'Maintain regular light exercise until late pregnancy',
      'Avoid strenuous activities',
      'Allow daily turnout for movement',
      'Monitor for signs of discomfort'
    ],
    healthcare: [
      'Schedule regular veterinary check-ups',
      'Keep vaccinations up to date',
      'Monitor for signs of complications',
      'Maintain proper hoof care',
      'Keep deworming schedule current'
    ],
    environment: [
      'Provide clean, safe housing',
      'Ensure adequate ventilation',
      'Maintain consistent routine',
      'Prepare foaling area in advance',
      'Monitor environmental temperatures'
    ]
  };

  return (
    <Stack spacing="lg">
      <Title order={2}>Pregnancy Care Guidelines</Title>

      <Card withBorder>
        <Title order={3} mb="md">Nutrition</Title>
        <List
          spacing="xs"
          size="sm"
          center
          icon={
            <ThemeIcon color="teal" size={24} radius="xl">
              <MdCheckCircle size="1rem" />
            </ThemeIcon>
          }
        >
          {guidelines.nutrition.map((item, index) => (
            <List.Item key={index}>{item}</List.Item>
          ))}
        </List>
      </Card>

      <Card withBorder>
        <Title order={3} mb="md">Exercise</Title>
        <List
          spacing="xs"
          size="sm"
          center
          icon={
            <ThemeIcon color="blue" size={24} radius="xl">
              <MdCheckCircle size="1rem" />
            </ThemeIcon>
          }
        >
          {guidelines.exercise.map((item, index) => (
            <List.Item key={index}>{item}</List.Item>
          ))}
        </List>
      </Card>

      <Card withBorder>
        <Title order={3} mb="md">Healthcare</Title>
        <List
          spacing="xs"
          size="sm"
          center
          icon={
            <ThemeIcon color="green" size={24} radius="xl">
              <MdCheckCircle size="1rem" />
            </ThemeIcon>
          }
        >
          {guidelines.healthcare.map((item, index) => (
            <List.Item key={index}>{item}</List.Item>
          ))}
        </List>
      </Card>

      <Card withBorder>
        <Title order={3} mb="md">Environment</Title>
        <List
          spacing="xs"
          size="sm"
          center
          icon={
            <ThemeIcon color="grape" size={24} radius="xl">
              <MdCheckCircle size="1rem" />
            </ThemeIcon>
          }
        >
          {guidelines.environment.map((item, index) => (
            <List.Item key={index}>{item}</List.Item>
          ))}
        </List>
      </Card>

      <Card withBorder>
        <Title order={3} mb="md">Important Notes</Title>
        <Text size="sm" c="dimmed">
          These guidelines are general recommendations. Always consult with your veterinarian for
          specific advice tailored to your mare's needs. Pay attention to any changes in behavior
          or condition and report concerns to your vet promptly.
        </Text>
      </Card>
    </Stack>
  );
}


import { Card, Stack, Text, SimpleGrid, ActionIcon, Group } from '@mantine/core';
import { CaretLeft, CaretRight } from '@phosphor-icons/react';
import { FC, useState } from 'react';

export const CalendarView: FC = () => {
  const [currentMonth, setCurrentMonth] = useState(new Date());
  
  const handlePreviousMonth = () => {
    setCurrentMonth(prev => {
      const newDate = new Date(prev);
      newDate.setMonth(prev.getMonth() - 1);
      return newDate;
    });
  };

  const handleNextMonth = () => {
    setCurrentMonth(prev => {
      const newDate = new Date(prev);
      newDate.setMonth(prev.getMonth() + 1);
      return newDate;
    });
  };

  return (
    <Card shadow="sm" padding="lg" radius="md">
      <Stack gap="md">
        <Group justify="space-between">
          <Text fw={700} size="xl">Due Date Calendar</Text>
          <Group>
            <ActionIcon 
              variant="subtle" 
              onClick={handlePreviousMonth}
              aria-label="Previous month"
            >
              <CaretLeft size={20} />
            </ActionIcon>
            <Text>{currentMonth.toLocaleDateString('en-US', { month: 'long', year: 'numeric' })}</Text>
            <ActionIcon 
              variant="subtle" 
              onClick={handleNextMonth}
              aria-label="Next month"
            >
              <CaretRight size={20} />
            </ActionIcon>
          </Group>
        </Group>
        
        <SimpleGrid cols={7}>
          {['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'].map((day) => (
            <Text key={day} ta="center" fw={500} c="dimmed">
              {day}
            </Text>
          ))}
          {/* We'll implement the actual calendar days in the next iteration */}
        </SimpleGrid>
      </Stack>
    </Card>
  );
}
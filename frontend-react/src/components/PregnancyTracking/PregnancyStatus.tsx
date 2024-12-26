import React from 'react';
import {
  Box,
  Typography,
  LinearProgress,
  List,
  ListItem,
  ListItemText,
  Divider,
  Chip
} from '@mui/material';
import { format, differenceInDays } from 'date-fns';

interface PregnancyStatusProps {
  status: {
    conceptionDate: string;
    currentStage: string;
    daysInPregnancy: number;
    expectedDueDate: string;
    lastEvent?: {
      date: string;
      eventType: string;
      description: string;
    };
  };
}

const PregnancyStatus: React.FC<PregnancyStatusProps> = ({ status }) => {
  const totalPregnancyDays = 340; // ~11 months
  const progress = (status.daysInPregnancy / totalPregnancyDays) * 100;
  const daysUntilDue = differenceInDays(new Date(status.expectedDueDate), new Date());

  const getStageColor = (stage: string) => {
    switch (stage) {
      case 'EARLY':
        return '#4CAF50';
      case 'MIDDLE':
        return '#2196F3';
      case 'LATE':
        return '#FF9800';
      case 'NEARTERM':
        return '#F44336';
      case 'FOALING':
        return '#9C27B0';
      default:
        return '#757575';
    }
  };

  return (
    <Box>
      <Typography variant="h6" gutterBottom>
        Pregnancy Status
      </Typography>

      <Box sx={{ mb: 3 }}>
        <Box display="flex" justifyContent="space-between" alignItems="center" mb={1}>
          <Typography variant="body2" color="textSecondary">
            Progress ({status.daysInPregnancy} days)
          </Typography>
          <Typography variant="body2" color="textSecondary">
            {Math.round(progress)}%
          </Typography>
        </Box>
        <LinearProgress 
          variant="determinate" 
          value={progress} 
          sx={{ height: 10, borderRadius: 5 }}
        />
      </Box>

      <List>
        <ListItem>
          <ListItemText
            primary="Current Stage"
            secondary={
              <Chip
                label={status.currentStage}
                sx={{
                  bgcolor: getStageColor(status.currentStage),
                  color: 'white',
                  mt: 1
                }}
              />
            }
          />
        </ListItem>
        
        <Divider />
        
        <ListItem>
          <ListItemText
            primary="Conception Date"
            secondary={format(new Date(status.conceptionDate), 'MMMM d, yyyy')}
          />
        </ListItem>
        
        <Divider />
        
        <ListItem>
          <ListItemText
            primary="Expected Due Date"
            secondary={
              <>
                {format(new Date(status.expectedDueDate), 'MMMM d, yyyy')}
                <Typography variant="body2" color="textSecondary">
                  ({daysUntilDue} days remaining)
                </Typography>
              </>
            }
          />
        </ListItem>

        {status.lastEvent && (
          <>
            <Divider />
            <ListItem>
              <ListItemText
                primary="Last Event"
                secondary={
                  <>
                    <Typography variant="body2">
                      {status.lastEvent.eventType}
                    </Typography>
                    <Typography variant="body2" color="textSecondary">
                      {format(new Date(status.lastEvent.date), 'MMMM d, yyyy')}
                    </Typography>
                    <Typography variant="body2">
                      {status.lastEvent.description}
                    </Typography>
                  </>
                }
              />
            </ListItem>
          </>
        )}
      </List>
    </Box>
  );
};

export default PregnancyStatus;

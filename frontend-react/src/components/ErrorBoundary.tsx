import React, { ErrorInfo, ReactNode } from 'react';
import { Alert, Container, Stack, Text, Button } from '@mantine/core';
import { IconAlertCircle } from '@tabler/icons-react';

interface ErrorBoundaryProps {
  children: ReactNode;
  fallback?: (error: Error) => ReactNode;
}

interface ErrorBoundaryState {
  hasError: boolean;
  error?: Error;
}

class ErrorBoundary extends React.Component<ErrorBoundaryProps, ErrorBoundaryState> {
  constructor(props: ErrorBoundaryProps) {
    super(props);
    this.state = { hasError: false };
  }

  static getDerivedStateFromError(error: Error): ErrorBoundaryState {
    return { hasError: true, error };
  }

  componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    // Log the error to an error reporting service like Sentry
    console.error('Uncaught error:', error, errorInfo);
  }

  handleReset = () => {
    this.setState({ hasError: false, error: undefined });
  }

  render() {
    if (this.state.hasError) {
      // Custom fallback UI or use provided fallback
      const defaultFallback = (
        <Container size="xs" my="xl">
          <Alert 
            icon={<IconAlertCircle size="1rem" />} 
            title="Something went wrong" 
            color="red" 
            variant="filled"
          >
            <Stack>
              <Text>An unexpected error occurred. Please try again.</Text>
              <Text size="xs" color="dimmed">
                Error Details: {this.state.error?.message}
              </Text>
              <Button 
                color="white" 
                variant="outline" 
                onClick={this.handleReset}
              >
                Try Again
              </Button>
            </Stack>
          </Alert>
        </Container>
      );

      return this.props.fallback 
        ? this.props.fallback(this.state.error!) 
        : defaultFallback;
    }

    return this.props.children;
  }
}

export default ErrorBoundary;

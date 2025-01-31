import { Alert, Button, Stack, Text, Title } from '@mantine/core';
import { Warning } from '@phosphor-icons/react';
import React, { Component, ErrorInfo, ReactNode } from 'react';

interface Props {
	children: ReactNode;
	fallback?: ReactNode;
}

interface State {
	hasError: boolean;
	error: Error | null;
	errorInfo: ErrorInfo | null;
}

export class ErrorBoundary extends Component<Props, State> {
	public state: State = {
		hasError: false,
		error: null,
		errorInfo: null,
	};

	public static getDerivedStateFromError(error: Error): State {
		return { hasError: true, error, errorInfo: null };
	}

	public componentDidCatch(error: Error, errorInfo: ErrorInfo) {
		console.error('Error caught by boundary:', error, errorInfo);
		this.setState({
			error,
			errorInfo,
		});
	}

	private handleReset = () => {
		this.setState({
			hasError: false,
			error: null,
			errorInfo: null,
		});
	};

	public render() {
		if (this.state.hasError) {
			return (
				<Stack align="center" justify="center" h="100%" spacing="lg" p="xl">
					<Alert 
						icon={<Warning size="1.5rem" />} 
						title="Something went wrong" 
						color="red"
						variant="filled"
					>
						<Stack spacing="md">
							<Text size="sm">
								{this.state.error?.message || 'An unexpected error occurred'}
							</Text>
							<Button 
								onClick={this.handleReset}
								variant="white"
								color="red"
							>
								Try Again
							</Button>
						</Stack>
					</Alert>
				</Stack>
			);
		}

		return this.props.children;
	}
}

export default ErrorBoundary;

import { Alert } from '@mantine/core';
import { Warning } from '@phosphor-icons/react';
import { Component, ReactNode } from 'react';

interface ErrorBoundaryProps {
	children: ReactNode;
}

interface ErrorBoundaryState {
	hasError: boolean;
	error: Error | null;
}

class ErrorBoundary extends Component<ErrorBoundaryProps, ErrorBoundaryState> {
	constructor(props: ErrorBoundaryProps) {
		super(props);
		this.state = {
			hasError: false,
			error: null,
		};
	}

	static getDerivedStateFromError(error: Error): ErrorBoundaryState {
		return {
			hasError: true,
			error,
		};
	}

	override componentDidCatch(error: Error, errorInfo: React.ErrorInfo): void {
		console.error('Error caught by boundary:', error, errorInfo);
	}

	override render(): ReactNode {
		if (this.state.hasError) {
			return (
				<Alert icon={<Warning size='1rem' />} title='Error' color='red'>
					{this.state.error?.message || 'Something went wrong'}
				</Alert>
			);
		}

		return this.props.children;
	}
}

export default ErrorBoundary;

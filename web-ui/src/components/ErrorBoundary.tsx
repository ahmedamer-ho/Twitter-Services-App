import { Component } from 'react';
import type { ErrorInfo, ReactNode } from 'react';

interface Props {
    children: ReactNode;
}

interface State {
    hasError: boolean;
    error: Error | null;
}

export class ErrorBoundary extends Component<Props, State> {
    public state: State = {
        hasError: false,
        error: null,
    };

    public static getDerivedStateFromError(error: Error): State {
        return { hasError: true, error };
    }

    public componentDidCatch(error: Error, errorInfo: ErrorInfo) {
        console.error('Uncaught error:', error, errorInfo);
    }

    public render() {
        if (this.state.hasError) {
            return (
                <div style={{
                    display: 'flex',
                    flexDirection: 'column',
                    alignItems: 'center',
                    justifyContent: 'center',
                    minHeight: '100vh',
                    padding: '20px',
                    textAlign: 'center'
                }}>
                    <div className="premium-card" style={{ maxWidth: '500px' }}>
                        <h2 style={{ color: 'var(--error)', marginBottom: '16px' }}>
                            Oops! Something went wrong
                        </h2>
                        <p style={{ color: 'var(--text-dim)', marginBottom: '24px' }}>
                            We're sorry for the inconvenience. Please try refreshing the page.
                        </p>
                        {this.state.error && (
                            <details style={{
                                marginBottom: '24px',
                                textAlign: 'left',
                                background: 'var(--bg-dark)',
                                padding: '12px',
                                borderRadius: '8px',
                                fontSize: '0.85rem'
                            }}>
                                <summary style={{ cursor: 'pointer', marginBottom: '8px' }}>
                                    Error details
                                </summary>
                                <code style={{ color: 'var(--error)' }}>
                                    {this.state.error.toString()}
                                </code>
                            </details>
                        )}
                        <button
                            className="btn-primary"
                            onClick={() => window.location.reload()}
                        >
                            Refresh Page
                        </button>
                    </div>
                </div>
            );
        }

        return this.props.children;
    }
}

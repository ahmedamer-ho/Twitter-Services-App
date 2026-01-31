import React, { useState } from 'react';
import { useAuth } from '../context/AuthContext.tsx';
import { api } from '../services/api.ts';

export const AuthFeature: React.FC = () => {
    const [isLogin, setIsLogin] = useState(true);
    const [formData, setFormData] = useState({
        username: '',
        password: '',
        email: '',
        firstname: '',
        lastname: '',
    });
    const [error, setError] = useState('');
    const [loading, setLoading] = useState(false);
    const { login } = useAuth();

    const handleToggle = () => {
        setIsLogin(!isLogin);
        setError('');
    };

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setFormData({ ...formData, [e.target.name]: e.target.value });
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError('');
        setLoading(true);

        try {
            if (isLogin) {
                const data = await api.post('/auth/login', {
                    username: formData.username,
                    password: formData.password,
                });

                // In a real app, you'd decode the JWT to get user info or fetch profile.
                // For now, let's mock the user object from the response or just use the username.
                login(data.token, {
                    id: 'placeholder-id',
                    username: formData.username,
                    email: formData.email || `${formData.username}@example.com`
                });
            } else {
                await api.post('/auth/register', {
                    username: formData.username,
                    password: formData.password,
                    email: formData.email,
                    firstname: formData.firstname,
                    lastname: formData.lastname,
                });
                setIsLogin(true);
                setError('Registered successfully! Please login.');
            }
        } catch (err: any) {
            setError(err.message || 'Authentication failed');
        } finally {
            setLoading(false);
        }
    };

    return (
        <div style={{
            display: 'flex',
            justifyContent: 'center',
            alignItems: 'center',
            minHeight: '100vh',
            padding: '20px'
        }}>
            <div className="premium-card" style={{ width: '100%', maxWidth: '400px' }}>
                <h2 style={{ marginBottom: '24px', textAlign: 'center' }}>
                    {isLogin ? 'Welcome Back' : 'Join Twitter'}
                </h2>

                {error && (
                    <div style={{
                        color: error.includes('successfully') ? 'var(--success)' : 'var(--error)',
                        marginBottom: '16px',
                        fontSize: '0.9rem',
                        textAlign: 'center'
                    }}>
                        {error}
                    </div>
                )}

                <form onSubmit={handleSubmit} style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
                    {!isLogin && (
                        <>
                            <div style={{ display: 'flex', gap: '12px' }}>
                                <input
                                    name="firstname"
                                    placeholder="First Name"
                                    required
                                    onChange={handleChange}
                                    style={{ width: '50%' }}
                                />
                                <input
                                    name="lastname"
                                    placeholder="Last Name"
                                    required
                                    onChange={handleChange}
                                    style={{ width: '50%' }}
                                />
                            </div>
                            <input
                                name="email"
                                type="email"
                                placeholder="Email Address"
                                required
                                onChange={handleChange}
                            />
                        </>
                    )}
                    <input
                        name="username"
                        placeholder="Username"
                        required
                        onChange={handleChange}
                    />
                    <input
                        name="password"
                        type="password"
                        placeholder="Password"
                        required
                        onChange={handleChange}
                    />

                    <button
                        type="submit"
                        className="btn-primary"
                        disabled={loading}
                        style={{ marginTop: '8px' }}
                    >
                        {loading ? 'Processing...' : (isLogin ? 'Log In' : 'Sign Up')}
                    </button>
                </form>

                <div style={{ marginTop: '24px', textAlign: 'center', color: 'var(--text-dim)', fontSize: '0.9rem' }}>
                    {isLogin ? "Don't have an account? " : "Already have an account? "}
                    <span
                        onClick={handleToggle}
                        style={{ color: 'var(--primary)', cursor: 'pointer', fontWeight: '600' }}
                    >
                        {isLogin ? 'Sign Up' : 'Log In'}
                    </span>
                </div>
            </div>
        </div>
    );
};

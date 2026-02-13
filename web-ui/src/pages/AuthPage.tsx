import React, { useState } from 'react';
import { useAuth } from '../context/AuthContext';
import { api, endpoints } from '../services/api';
import '../styles/auth.css';

export const AuthPage: React.FC = () => {
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
                const data = await api.post(endpoints.auth.login, {
                    username: formData.username,
                    password: formData.password,
                });

                // In a real app the token should contain user ID and payload
                // For now we assume the backend returns what we need or we decode it
                login(data.token, {
                    id: data.id || 'placeholder-id',
                    username: formData.username,
                    email: formData.email || `${formData.username}@example.com`
                });
            } else {
                await api.post(endpoints.auth.register, {
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
        <div className="auth-container">
            <div className="auth-box">
                <div className="auth-logo">
                    <h1>𝕏</h1>
                </div>

                <h2 className="auth-title">
                    {isLogin ? 'Sign in to X' : 'Join X today'}
                </h2>

                {error && (
                    <div className={`auth-error ${error.includes('successfully') ? 'auth-success' : ''}`}>
                        {error}
                    </div>
                )}

                <form onSubmit={handleSubmit} className="auth-form">
                    {!isLogin && (
                        <>
                            <div className="auth-input-group">
                                <input
                                    name="firstname"
                                    placeholder="First Name"
                                    required
                                    onChange={handleChange}
                                    className="auth-input"
                                />
                                <input
                                    name="lastname"
                                    placeholder="Last Name"
                                    required
                                    onChange={handleChange}
                                    className="auth-input"
                                />
                            </div>
                            <input
                                name="email"
                                type="email"
                                placeholder="Email Address"
                                required
                                onChange={handleChange}
                                className="auth-input"
                            />
                        </>
                    )}
                    <input
                        name="username"
                        placeholder="Username"
                        required
                        onChange={handleChange}
                        className="auth-input"
                    />
                    <input
                        name="password"
                        type="password"
                        placeholder="Password"
                        required
                        onChange={handleChange}
                        className="auth-input"
                    />

                    <button
                        type="submit"
                        disabled={loading}
                        className="auth-submit-btn"
                    >
                        {loading ? 'Processing...' : (isLogin ? 'Next' : 'Create account')}
                    </button>
                </form>

                <div className="auth-footer">
                    {isLogin ? "Don't have an account? " : "Already have an account? "}
                    <button
                        onClick={handleToggle}
                        className="auth-toggle-link"
                    >
                        {isLogin ? 'Sign up' : 'Sign in'}
                    </button>
                </div>
            </div>
        </div>
    );
};

import React, { useState, useEffect } from 'react';
import { useAuth } from '../context/AuthContext.tsx';
import { api } from '../services/api.ts';

interface Tweet {
    id: string;
    authorId: string;
    content: string;
    createdAt: string;
}

export const MainTimeline: React.FC = () => {
    const { user, logout } = useAuth();
    const [tweets, setTweets] = useState<Tweet[]>([]);
    const [content, setContent] = useState('');
    const [loading, setLoading] = useState(false);
    const [fetching, setFetching] = useState(false);

    const fetchTimeline = async () => {
        setFetching(true);
        try {
            // Fetch global timeline for demonstration
            const data = await api.get('/timeline/global');
            setTweets(data.tweets || []);
        } catch (err) {
            console.error('Failed to fetch timeline', err);
        } finally {
            setFetching(false);
        }
    };

    useEffect(() => {
        fetchTimeline();
        const interval = setInterval(fetchTimeline, 10000); // Poll every 10s
        return () => clearInterval(interval);
    }, []);

    const handlePostTweet = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!content.trim()) return;

        setLoading(true);
        try {
            await api.post('/tweets/', { content }, {
                'Idempotency-Key': `web-${Date.now()}`
            });
            setContent('');
            await fetchTimeline();
        } catch (err) {
            console.error('Failed to post tweet', err);
        } finally {
            setLoading(false);
        }
    };

    return (
        <div style={{ display: 'flex', minHeight: '100vh', background: 'var(--bg-dark)' }}>
            {/* Sidebar */}
            <nav style={{
                width: '275px',
                padding: '20px',
                borderRight: '1px solid var(--border)',
                display: 'flex',
                flexDirection: 'column',
                position: 'fixed',
                height: '100vh'
            }}>
                <h1 style={{ marginBottom: '32px', color: 'var(--primary)', fontSize: '2rem' }}>𝕏</h1>
                <div style={{ display: 'flex', flexDirection: 'column', gap: '20px', flex: 1 }}>
                    <div style={{ fontWeight: 'bold', fontSize: '1.2rem', cursor: 'pointer' }}>Home</div>
                    <div style={{ color: 'var(--text-dim)', cursor: 'not-allowed' }}>Explore</div>
                    <div style={{ color: 'var(--text-dim)', cursor: 'not-allowed' }}>Notifications</div>
                    <div style={{ color: 'var(--text-dim)', cursor: 'not-allowed' }}>Messages</div>
                </div>
                <div style={{ borderTop: '1px solid var(--border)', paddingTop: '20px' }}>
                    <div style={{ marginBottom: '12px', fontWeight: 'bold' }}>@{user?.username}</div>
                    <button onClick={logout} className="btn-outline" style={{ width: '100%' }}>Log Out</button>
                </div>
            </nav>

            {/* Main Content */}
            <main style={{ marginLeft: '275px', width: '600px', borderRight: '1px solid var(--border)', minHeight: '100vh' }}>
                <header className="glass" style={{
                    padding: '16px',
                    position: 'sticky',
                    top: 0,
                    zIndex: 10,
                    borderBottom: '1px solid var(--border)'
                }}>
                    <h3>Home</h3>
                </header>

                {/* Tweet Box */}
                <div style={{ padding: '16px', borderBottom: '1px solid var(--border)' }}>
                    <form onSubmit={handlePostTweet}>
                        <textarea
                            placeholder="What is happening?!"
                            value={content}
                            onChange={(e) => setContent(e.target.value)}
                            style={{
                                width: '100%',
                                background: 'transparent',
                                border: 'none',
                                fontSize: '1.2rem',
                                resize: 'none',
                                minHeight: '100px',
                                padding: '0'
                            }}
                        />
                        <div style={{ display: 'flex', justifyContent: 'flex-end', paddingTop: '12px', borderTop: '1px solid var(--border)' }}>
                            <button
                                type="submit"
                                className="btn-primary"
                                disabled={loading || !content.trim()}
                                style={{ borderRadius: '9999px' }}
                            >
                                {loading ? 'Posting...' : 'Tweet'}
                            </button>
                        </div>
                    </form>
                </div>

                {/* Feed */}
                <div>
                    {fetching && tweets.length === 0 ? (
                        <div style={{ padding: '40px', textAlign: 'center', color: 'var(--text-dim)' }}>Loading tweets...</div>
                    ) : (
                        tweets.map((tweet) => (
                            <div key={tweet.id} className="premium-card" style={{
                                border: 'none',
                                borderBottom: '1px solid var(--border)',
                                borderRadius: 0,
                                padding: '16px'
                            }}>
                                <div style={{ fontWeight: 'bold', marginBottom: '4px' }}>
                                    User {tweet.authorId.substring(0, 8)}... <span style={{ fontWeight: 'normal', color: 'var(--text-dim)', fontSize: '0.9rem' }}>· {new Date(tweet.createdAt).toLocaleTimeString()}</span>
                                </div>
                                <div style={{ fontSize: '1rem', whiteSpace: 'pre-wrap' }}>{tweet.content}</div>
                            </div>
                        ))
                    )}
                    {tweets.length === 0 && !fetching && (
                        <div style={{ padding: '40px', textAlign: 'center', color: 'var(--text-dim)' }}>No tweets yet! Follow someone or post one.</div>
                    )}
                </div>
            </main>

            {/* Widgets (Placeholder) */}
            <aside style={{ padding: '20px', flex: 1 }}>
                <div className="premium-card" style={{ background: 'var(--surface)', borderRadius: '16px' }}>
                    <h3 style={{ marginBottom: '16px' }}>What's happening</h3>
                    <div style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
                        <div>
                            <div style={{ fontSize: '0.8rem', color: 'var(--text-dim)' }}>Trending in Technology</div>
                            <div style={{ fontWeight: 'bold' }}>#Microservices</div>
                        </div>
                        <div>
                            <div style={{ fontSize: '0.8rem', color: 'var(--text-dim)' }}>Trending in Software</div>
                            <div style={{ fontWeight: 'bold' }}>#Golang</div>
                        </div>
                        <div>
                            <div style={{ fontSize: '0.8rem', color: 'var(--text-dim)' }}>Trending in Web</div>
                            <div style={{ fontWeight: 'bold' }}>#ReactTS</div>
                        </div>
                    </div>
                </div>
            </aside>
        </div>
    );
};

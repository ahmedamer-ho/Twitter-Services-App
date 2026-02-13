import React, { useState } from 'react';
import { useAuth } from '../../context/AuthContext';
import { api, endpoints } from '../../services/api';
import '../../styles/tweet.css';

interface TweetBoxProps {
    onTweetPosted: () => void;
}

export const TweetBox: React.FC<TweetBoxProps> = ({ onTweetPosted }) => {
    const { user } = useAuth();
    const [content, setContent] = useState('');
    const [loading, setLoading] = useState(false);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!content.trim()) return;

        setLoading(true);
        try {
            await api.post(endpoints.tweets.create, { content }, {
                headers: {
                    'Idempotency-Key': `web-${Date.now()}`
                }
            });
            setContent('');
            onTweetPosted();
        } catch (error) {
            console.error('Failed to post tweet:', error);
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="tweet-box">
            <div className="tweet-avatar">
                {user?.username[0]?.toUpperCase()}
            </div>
            <div className="tweet-box-content">
                <form onSubmit={handleSubmit}>
                    <textarea
                        className="tweet-input"
                        placeholder="What is happening?!"
                        value={content}
                        onChange={(e) => setContent(e.target.value)}
                    />
                    <div className="tweet-box-footer">
                        <div className="text-primary text-sm font-bold">
                            {/* Icons for media could go here */}
                        </div>
                        <button
                            type="submit"
                            className="tweet-submit-btn"
                            disabled={loading || !content.trim()}
                        >
                            {loading ? 'Posting...' : 'Tweet'}
                        </button>
                    </div>
                </form>
            </div>
        </div>
    );
};

import React, { useState, useEffect } from 'react';
import { api, endpoints } from '../services/api';
import { Tweet } from '../components/tweets/Tweet';
import { TweetBox } from '../components/tweets/TweetBox';
import '../styles/home.css';

interface TweetData {
    id: string;
    authorId: string;
    content: string;
    createdAt: string;
}

export const HomePage: React.FC = () => {
    const [tweets, setTweets] = useState<TweetData[]>([]);
    const [loading, setLoading] = useState(false);

    const fetchTimeline = async () => {
        setLoading(true);
        try {
            const data = await api.get(endpoints.tweets.feed);
            setTweets(data.tweets || []);
        } catch (err) {
            console.error('Failed to fetch timeline', err);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchTimeline();
    }, []);

    return (
        <div>
            <div className="home-header">
                <h2>Home</h2>
            </div>

            <TweetBox onTweetPosted={fetchTimeline} />

            <div className="timeline">
                {loading && tweets.length === 0 ? (
                    <div className="timeline-loading">Loading tweets...</div>
                ) : (
                    tweets.map((tweet) => (
                        <Tweet key={tweet.id} tweet={tweet} />
                    ))
                )}

                {!loading && tweets.length === 0 && (
                    <div className="timeline-empty">
                        No tweets yet. Be the first to post!
                    </div>
                )}
            </div>
        </div>
    );
};

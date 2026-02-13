import React from 'react';
import { formatDistanceToNow } from 'date-fns';
import { FaRegComment, FaRetweet, FaRegHeart, FaShare } from 'react-icons/fa';
import '../../styles/tweet.css';

interface TweetProps {
    tweet: {
        id: string;
        authorId: string;
        content: string;
        createdAt: string;
        likes?: number;
        retweets?: number;
        replies?: number;
    };
}

export const Tweet: React.FC<TweetProps> = ({ tweet }) => {
    return (
        <article className="tweet">
            <div className="tweet-avatar">
                {tweet.authorId[0]?.toUpperCase() || 'U'}
            </div>
            <div className="tweet-content">
                <div className="tweet-header">
                    <span className="tweet-author-name">User {tweet.authorId.substring(0, 8)}</span>
                    <span className="tweet-author-handle">@{tweet.authorId.substring(0, 8)}</span>
                    <span className="tweet-time">· {formatDistanceToNow(new Date(tweet.createdAt), { addSuffix: true })}</span>
                </div>
                <div className="tweet-text">{tweet.content}</div>
                <div className="tweet-actions">
                    <button className="tweet-action-btn">
                        <FaRegComment />
                        <span>{tweet.replies || 0}</span>
                    </button>
                    <button className="tweet-action-btn retweet">
                        <FaRetweet />
                        <span>{tweet.retweets || 0}</span>
                    </button>
                    <button className="tweet-action-btn like">
                        <FaRegHeart />
                        <span>{tweet.likes || 0}</span>
                    </button>
                    <button className="tweet-action-btn">
                        <FaShare />
                    </button>
                </div>
            </div>
        </article>
    );
};
